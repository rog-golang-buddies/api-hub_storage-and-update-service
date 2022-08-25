package openapi

import (
	"context"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	"github.com/rog-golang-buddies/internal/logger"
)

func parseOpenAPI(ctx context.Context, content []byte) (*openapi3.T, error) {
	loader := openapi3.Loader{Context: ctx, IsExternalRefsAllowed: false}
	doc, err := loader.LoadFromData(content)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

func openapiToApiSpec(log logger.Logger, openapi *openapi3.T) *apispecdoc.ApiSpecDoc {
	asd := apispecdoc.ApiSpecDoc{
		Title:       openapi.Info.Title,
		Description: openapi.Info.Description,
		Type:        apispecdoc.TypeOpenApi,
		Methods:     make([]*apispecdoc.ApiMethod, 0),
	}

	groups := tagToGroup(openapi.Tags)
	groupMap := make(map[string]*apispecdoc.Group)
	for _, group := range groups {
		groupMap[group.Name] = group
	}

	asd.Groups = groups

	populateMethods(log, &asd, openapi.Paths, openapi.Servers)
	return &asd
}

// tagToGroup creates group with empty methods
func tagToGroup(tags []*openapi3.Tag) []*apispecdoc.Group {
	groups := make([]*apispecdoc.Group, 0, len(tags))
	if tags == nil {
		return groups
	}
	for _, tag := range tags {
		group := new(apispecdoc.Group)
		group.Name = tag.Name
		group.Description = tag.Description
		group.Methods = make([]*apispecdoc.ApiMethod, 0)
		groups = append(groups, group)
	}
	return groups
}

func populateMethods(log logger.Logger, asd *apispecdoc.ApiSpecDoc, paths openapi3.Paths, rootServers openapi3.Servers) {
	groupMap := make(map[string]*apispecdoc.Group)
	for _, group := range asd.Groups {
		groupMap[group.Name] = group
	}
	for url, path := range paths {
		for httpMethod, operation := range path.Operations() {
			method := new(apispecdoc.ApiMethod)
			method.Path = url
			method.Description = operation.Description
			method.ExternalDoc = convertExternalDoc(operation.ExternalDocs)
			method.Type = apispecdoc.MethodType(httpMethod)
			method.Parameters = convertParameters(operation.Parameters)
			if operation.Servers != nil {
				method.Servers = convertServers(*operation.Servers)
			} else {
				method.Servers = convertServers(rootServers)
			}
			if operation.RequestBody != nil {
				method.RequestBody = convertBody(operation.RequestBody.Value)
			}
			if operation.Tags != nil && len(operation.Tags) > 0 {
				addedToAnyGroup := false
				for _, tag := range operation.Tags {
					if group, ok := groupMap[tag]; ok {
						group.Methods = append(group.Methods, method)
						addedToAnyGroup = true
					} else {
						log.Warnf("inconsistent state found; tag %s not mentioned in the tags section", tag)
					}
				}
				if !addedToAnyGroup {
					asd.Methods = append(asd.Methods, method)
				}
			} else {
				asd.Methods = append(asd.Methods, method)
			}
		}
	}
}

func convertExternalDoc(oEDocs *openapi3.ExternalDocs) *apispecdoc.ExternalDoc {
	if oEDocs == nil {
		return nil
	}
	return &apispecdoc.ExternalDoc{
		Description: oEDocs.Description,
		Url:         oEDocs.URL,
	}
}

func convertServers(oServers openapi3.Servers) []*apispecdoc.Server {
	servers := make([]*apispecdoc.Server, 0, len(oServers))
	for _, oServ := range oServers {
		server := apispecdoc.Server{
			Url:         oServ.URL,
			Description: oServ.Description,
		}
		servers = append(servers, &server)
	}
	return servers
}

func convertParameters(oParams openapi3.Parameters) []*apispecdoc.Parameter {
	resParams := make([]*apispecdoc.Parameter, 0, len(oParams))
	for _, oParRef := range oParams {
		oPar := oParRef.Value
		if oPar == nil || oPar.Schema == nil {
			continue
		}

		param := apispecdoc.Parameter{
			Name:        oPar.Name,
			In:          apispecdoc.ParameterType(oPar.In),
			Description: oPar.Description,
			Schema:      convertSchema("", oParRef.Value.Schema.Value),
			Required:    oPar.Required,
		}
		resParams = append(resParams, &param)
	}
	return resParams
}

func convertBody(body *openapi3.RequestBody) *apispecdoc.RequestBody {
	specBody := new(apispecdoc.RequestBody)
	specBody.Description = body.Description
	specBody.Required = body.Required
	specContent := make([]*apispecdoc.MediaTypeObject, 0, len(body.Content))
	for cType, content := range body.Content {
		if content.Schema == nil || content.Schema.Value == nil {
			continue
		}
		specContent = append(specContent,
			&apispecdoc.MediaTypeObject{MediaType: cType, Schema: convertSchema("", content.Schema.Value)})
	}
	specBody.Content = specContent
	return specBody
}

func convertSchema(key string, schema *openapi3.Schema) *apispecdoc.Schema {
	resSchema := new(apispecdoc.Schema)
	if schema == nil {
		return resSchema
	}
	resSchema.Key = key
	resSchema.Description = schema.Description
	resSchema.Type = apispecdoc.ResolveSchemaType(schema.Type)
	resSchema.Fields = make([]*apispecdoc.Schema, 0)
	switch resSchema.Type {
	case apispecdoc.Object:
		//If the type is an Object it can be an Object or Map. The map represents additional properties - can be only one of Object/Map
		if schema.Properties != nil {
			for pKey, prop := range schema.Properties {
				resSchema.Fields = append(resSchema.Fields, convertSchema(pKey, prop.Value))
			}
		} else if schema.AdditionalProperties != nil {
			resSchema.Type = apispecdoc.Map
			resSchema.Fields = append(resSchema.Fields, convertSchema("", schema.AdditionalProperties.Value))
		}
	case apispecdoc.Array:
		if schema.Items != nil {
			resSchema.Fields = append(resSchema.Fields, convertSchema("", schema.Items.Value))
		}
	case apispecdoc.NotDefined:
		//If type is not defined it means that here one of the "combine" types is used. So need to check them all
		switch true {
		case schema.OneOf != nil && len(schema.OneOf) > 0:
			resSchema.Type = apispecdoc.OneOf
			for _, sch := range schema.OneOf {
				resSchema.Fields = append(resSchema.Fields, convertSchema("", sch.Value))
			}
		case schema.AnyOf != nil && len(schema.AnyOf) > 0:
			resSchema.Type = apispecdoc.AnyOf
			for _, sch := range schema.AnyOf {
				resSchema.Fields = append(resSchema.Fields, convertSchema("", sch.Value))
			}
		case schema.AllOf != nil && len(schema.AllOf) > 0:
			resSchema.Type = apispecdoc.AllOf
			for _, sch := range schema.AllOf {
				resSchema.Fields = append(resSchema.Fields, convertSchema("", sch.Value))
			}
		case schema.Not != nil:
			resSchema.Type = apispecdoc.Not
			resSchema.Fields = append(resSchema.Fields, convertSchema("", schema.Not.Value))
		}
	}

	return resSchema
}
