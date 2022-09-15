package service

import (
	"context"
	"encoding/json"
	"errors"
	asdentity "github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/apispecdoc"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/dto"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/logger"
	"github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	"github.com/rog-golang-buddies/api_hub_common/apispecproto"
	"time"
)

func NewService(log logger.Logger, repo asdentity.AsdRepository) asdentity.Service {
	return &ServiceImpl{log: log, asdRepo: repo}
}

var asdTypeMap = map[string]apispecproto.Type{
	string(apispecdoc.TypeOpenApi): apispecproto.Type_OPEN_API,
}

var methodTypeMap = map[string]apispecproto.MethodType{
	string(apispecdoc.MethodConnect): apispecproto.MethodType_CONNECT,
	string(apispecdoc.MethodGet):     apispecproto.MethodType_GET,
	string(apispecdoc.MethodPut):     apispecproto.MethodType_PUT,
	string(apispecdoc.MethodPost):    apispecproto.MethodType_POST,
	string(apispecdoc.MethodDelete):  apispecproto.MethodType_DELETE,
	string(apispecdoc.MethodOptions): apispecproto.MethodType_OPTIONS,
	string(apispecdoc.MethodHead):    apispecproto.MethodType_HEAD,
	string(apispecdoc.MethodPatch):   apispecproto.MethodType_PATCH,
	string(apispecdoc.MethodTrace):   apispecproto.MethodType_TRACE,
}

var schemaTypeMap = map[apispecdoc.SchemaType]apispecproto.SchemaType{
	apispecdoc.Unknown:    apispecproto.SchemaType_UNKNOWN_SCHEMA,
	apispecdoc.NotDefined: apispecproto.SchemaType_NOT_DEFINED,
	apispecdoc.Integer:    apispecproto.SchemaType_INTEGER,
	apispecdoc.Boolean:    apispecproto.SchemaType_BOOLEAN,
	apispecdoc.Number:     apispecproto.SchemaType_NUMBER,
	apispecdoc.String:     apispecproto.SchemaType_STRING,
	apispecdoc.Date:       apispecproto.SchemaType_DATE,
	apispecdoc.Array:      apispecproto.SchemaType_ARRAY,
	apispecdoc.Map:        apispecproto.SchemaType_MAP,
	apispecdoc.OneOf:      apispecproto.SchemaType_ONE_OF,
	apispecdoc.AnyOf:      apispecproto.SchemaType_ANY_OF,
	apispecdoc.AllOf:      apispecproto.SchemaType_ALL_OF,
	apispecdoc.Not:        apispecproto.SchemaType_NOT,
	apispecdoc.Object:     apispecproto.SchemaType_OBJECT,
}

var parameterTypeMap = map[apispecdoc.ParameterType]apispecproto.ParameterType{
	apispecdoc.ParameterQuery:  apispecproto.ParameterType_QUERY,
	apispecdoc.ParameterHeader: apispecproto.ParameterType_HEADER,
	apispecdoc.ParameterPath:   apispecproto.ParameterType_PATH,
	apispecdoc.ParameterCookie: apispecproto.ParameterType_COOKIE,
}

type ServiceImpl struct {
	log     logger.Logger
	asdRepo asdentity.AsdRepository
}

func (s *ServiceImpl) Search(ctx context.Context, req *apispecproto.SearchRequest) (*apispecproto.SearchResponse, error) {
	if req == nil {
		s.log.Error("nil request body received")
		return nil, errors.New("request body must not be nil")
	}
	pageReq := dto.PageRequest{}
	if req.Page != nil {
		pageReq.Page = int(*req.Page) - 1
	}
	if req.PerPage != nil {
		pageReq.PerPage = int(*req.PerPage)
	} else {
		pageReq.PerPage = 10
	}
	asdPage, err := s.asdRepo.SearchShort(ctx, req.Search, pageReq)
	if err != nil {
		return nil, err
	}
	res := new(apispecproto.SearchResponse)
	resDocs := make([]*apispecproto.ShortASD, 0)
	for _, asd := range asdPage.Data {
		resDocs = append(resDocs, &apispecproto.ShortASD{
			Id:          uint32(asd.ID),
			Name:        asd.Title,
			Description: asd.Description,
		})
	}
	res.ShortSpecDocs = resDocs
	res.Page = &apispecproto.Page{
		Total:   int32(asdPage.Total),
		Current: int32(asdPage.Page),
		PerPage: int32(asdPage.PerPage),
	}
	return res, nil
}

func (s *ServiceImpl) Get(ctx context.Context, req *apispecproto.GetRequest) (*apispecproto.GetResponse, error) {
	if req == nil {
		s.log.Error("nil request body received")
		return nil, errors.New("request body must not be nil")
	}
	apiSpecDoc, err := s.asdRepo.FindById(ctx, uint(req.Id))
	if err != nil {
		s.log.Error("error while find ASD by ID: ", err)
		return nil, err
	}
	if apiSpecDoc == nil {
		s.log.Infof("API spec document not found by id %d", req.Id)
		return nil, nil
	}
	resAsd, err := entityToFullAsd(apiSpecDoc)
	if err != nil {
		return nil, err
	}
	return &apispecproto.GetResponse{ApiSpecDoc: resAsd}, nil
}

func (s *ServiceImpl) Save(ctx context.Context, asd *apispecdoc.ApiSpecDoc) (uint, error) {
	if asd == nil {
		return 0, errors.New("nil asd model received")
	}
	asdEntity, err := asdToEntity(asd)
	if err != nil {
		return 0, err
	}
	//Check records by md5 hash sum - if exists than all methods the same and update not required
	asdByHash, err := s.asdRepo.FindByHash(ctx, asd.Md5Sum)
	if err != nil {
		return 0, err
	}
	if asdByHash != nil {
		s.log.Infof("record '%s' hash '%s' no changes", asd.Title, asd.Md5Sum)
		return asdByHash.ID, nil
	}
	//Check records by file url - if exists than need to update ASD in db (prev step didn't find matched hash - so hash changed)
	asdByUrl, err := s.asdRepo.FindByUrl(ctx, asd.Url)
	if err != nil {
		return 0, err
	}
	if asdByUrl != nil {
		asdByUrl.Title = asdEntity.Title
		asdByUrl.Description = asdEntity.Description
		asdByUrl.Type = asdEntity.Type
		asdByUrl.Groups = asdEntity.Groups
		asdByUrl.ApiMethods = asdEntity.ApiMethods
		asdByUrl.Md5sum = asdEntity.Md5sum
		asdByUrl.Url = asdEntity.Url
		//clear and reattach all dependencies
		err = s.asdRepo.Update(ctx, asdByUrl)
		if err != nil {
			return 0, err
		}
		s.log.Infof("record '%s' with hash '%s' updated", asd.Title, asd.Md5Sum)
		return asdByUrl.ID, nil
	}
	s.log.Infof("create new record for '%s' hash '%s'", asd.Title, asd.Md5Sum)
	return s.asdRepo.Save(ctx, asdEntity)
}

func asdToEntity(dto *apispecdoc.ApiSpecDoc) (*asdentity.ApiSpecDoc, error) {
	groups := make([]*asdentity.Group, 0, len(dto.Groups))
	for _, group := range dto.Groups {
		methods, err := methodsToEntities(group.Methods)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &asdentity.Group{
			Name:        group.Name,
			Description: group.Description,
			ApiMethods:  methods,
		})
	}
	methods, err := methodsToEntities(dto.Methods)
	if err != nil {
		return nil, err
	}
	return &asdentity.ApiSpecDoc{
		Title:       dto.Title,
		Description: dto.Description,
		Type:        string(dto.Type),
		Groups:      groups,
		ApiMethods:  methods,
		Md5sum:      dto.Md5Sum,
		Url:         dto.Url,
		FetchedAt:   time.Now(),
	}, nil
}

func methodToEntity(method *apispecdoc.ApiMethod) (*asdentity.ApiMethod, error) {
	if method == nil {
		return &asdentity.ApiMethod{}, nil
	}
	var params []byte
	var err error
	if method.Parameters != nil {
		params, err = json.Marshal(method.Parameters)
		if err != nil {
			return nil, err
		}
	}
	var body []byte
	if method.RequestBody != nil {
		body, err = json.Marshal(method.RequestBody)
		if err != nil {
			return nil, err
		}
	}
	var servers []*asdentity.Server
	if method.Servers != nil {
		servers = make([]*asdentity.Server, 0, len(method.Servers))
		for _, server := range method.Servers {
			servers = append(servers, &asdentity.Server{
				URL:         server.Url,
				Description: server.Description,
			})
		}
	}
	var extDoc asdentity.ExternalDoc
	if method.ExternalDoc != nil {
		extDoc = asdentity.ExternalDoc{
			Description: method.ExternalDoc.Description,
			URL:         method.ExternalDoc.Url,
		}
	}
	return &asdentity.ApiMethod{
		Path:        method.Path,
		Name:        method.Name,
		Description: method.Description,
		Type:        string(method.Type),
		Parameters:  string(params),
		Servers:     servers,
		RequestBody: string(body),
		ExternalDoc: &extDoc,
	}, nil
}

func methodsToEntities(methods []*apispecdoc.ApiMethod) ([]*asdentity.ApiMethod, error) {
	if methods == nil {
		return make([]*asdentity.ApiMethod, 0), nil
	}
	resMeth := make([]*asdentity.ApiMethod, 0, len(methods))
	for _, method := range methods {
		methEntity, err := methodToEntity(method)
		if err != nil {
			return nil, err
		}
		resMeth = append(resMeth, methEntity)
	}
	return resMeth, nil
}

func entityToFullAsd(asd *asdentity.ApiSpecDoc) (*apispecproto.FullASD, error) {
	convertMethods := func(methods []*asdentity.ApiMethod) ([]*apispecproto.ApiMethod, error) {
		if methods == nil {
			return nil, nil
		}
		resMethods := make([]*apispecproto.ApiMethod, 0, len(methods))
		for _, method := range methods {
			asdMethod, err := entityToFullASDMethod(method)
			if err != nil {
				return nil, err
			}
			resMethods = append(resMethods, asdMethod)
		}
		return resMethods, nil
	}
	rootMethods, err := convertMethods(asd.ApiMethods)
	if err != nil {
		return nil, err
	}

	groups := make([]*apispecproto.Group, 0, len(asd.Groups))
	for _, group := range asd.Groups {
		methods, err := convertMethods(group.ApiMethods)
		if err != nil {
			return nil, err
		}
		groups = append(groups, &apispecproto.Group{
			Name:        group.Name,
			Description: group.Description,
			Methods:     methods,
		})
	}
	return &apispecproto.FullASD{
		Id:          uint32(asd.ID),
		Title:       asd.Title,
		Description: asd.Description,
		Url:         asd.Url,
		Type:        apiTypeToResponse(asd.Type),
		Groups:      groups,
		Methods:     rootMethods,
	}, nil
}

func entityToFullASDMethod(method *asdentity.ApiMethod) (*apispecproto.ApiMethod, error) {
	if method == nil {
		return &apispecproto.ApiMethod{}, nil
	}
	var params []*apispecdoc.Parameter
	var err error
	if len(method.Parameters) != 0 {
		if err = json.Unmarshal([]byte(method.Parameters), &params); err != nil {
			return nil, err
		}
	}

	var resParams []*apispecproto.Parameter
	if params != nil {
		resParams = make([]*apispecproto.Parameter, 0, len(params))
		for _, param := range params {
			resParams = append(resParams, &apispecproto.Parameter{
				Name:        param.Name,
				In:          paramTypeToResponse(param.In),
				Description: param.Description,
				Required:    param.Required,
				Schema:      entitySchemaToResponse(param.Schema),
			})
		}
	}

	var body apispecdoc.RequestBody
	var resBody *apispecproto.RequestBody
	if len(method.RequestBody) != 0 {
		if err = json.Unmarshal([]byte(method.RequestBody), &body); err != nil {
			return nil, err
		}
		resBody = new(apispecproto.RequestBody)
		resBody.Description = body.Description
		resBody.Required = body.Required
		if body.Content != nil {
			resContent := make([]*apispecproto.RequestBody_MediaTypeObject, 0, len(body.Content))
			for _, mediaTypeObj := range body.Content {
				resContent = append(resContent, &apispecproto.RequestBody_MediaTypeObject{
					MediaType: mediaTypeObj.MediaType,
					Schema:    entitySchemaToResponse(mediaTypeObj.Schema),
				})
			}
			resBody.Content = resContent
		}
	}

	var resServers []*apispecproto.Server
	if method.Servers != nil {
		resServers = make([]*apispecproto.Server, 0, len(method.Servers))
		for _, server := range method.Servers {
			resServers = append(resServers, &apispecproto.Server{
				Url:         server.URL,
				Description: server.Description,
			})
		}
	}
	var resExtDoc = &apispecproto.ExternalDoc{}
	if method.ExternalDoc != nil {
		resExtDoc.Description = method.ExternalDoc.Description
		resExtDoc.Url = method.ExternalDoc.URL
	}
	return &apispecproto.ApiMethod{
		Path:        method.Path,
		Name:        method.Name,
		Description: method.Description,
		Type:        methodTypeToResponse(method.Type),
		Parameters:  resParams,
		Servers:     resServers,
		RequestBody: resBody,
		ExternalDoc: resExtDoc,
	}, nil
}

func entitySchemaToResponse(schema *apispecdoc.Schema) *apispecproto.Schema {
	if schema == nil {
		return nil
	}
	return &apispecproto.Schema{
		Key:         schema.Key,
		Type:        schemaTypeToResponse(schema.Type),
		Description: schema.Description,
		Fields:      entitySchemasToResponses(schema.Fields),
	}
}

func entitySchemasToResponses(schemas []*apispecdoc.Schema) []*apispecproto.Schema {
	resSchemas := make([]*apispecproto.Schema, 0, len(schemas))
	for _, schema := range schemas {
		resSchemas = append(resSchemas, entitySchemaToResponse(schema))
	}
	return resSchemas
}

func paramTypeToResponse(pt apispecdoc.ParameterType) apispecproto.ParameterType {
	res, ok := parameterTypeMap[pt]
	if ok {
		return res
	}
	return apispecproto.ParameterType_UNKNOWN_PARAM
}

func methodTypeToResponse(mt string) apispecproto.MethodType {
	res, ok := methodTypeMap[mt]
	if ok {
		return res
	}
	return apispecproto.MethodType_UNKNOWN_METHOD
}

func schemaTypeToResponse(st apispecdoc.SchemaType) apispecproto.SchemaType {
	res, ok := schemaTypeMap[st]
	if ok {
		return res
	}
	return apispecproto.SchemaType_UNKNOWN_SCHEMA
}

func apiTypeToResponse(asdT string) apispecproto.Type {
	res, ok := asdTypeMap[asdT]
	if ok {
		return res
	}
	return apispecproto.Type_UNKNOWN_API
}
