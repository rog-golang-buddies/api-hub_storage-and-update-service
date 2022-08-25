package openapi

import (
	"context"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	mock_logger "github.com/rog-golang-buddies/internal/logger/mocks"
	"github.com/stretchr/testify/assert"
)

func TestParseOpenAPI(t *testing.T) {
	ctx := context.Background()
	content, err := os.ReadFile("./mocks/github_mock.yml")
	if err != nil {
		return
	}
	openAPI, err := parseOpenAPI(ctx, content)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	assert.NotNil(t, openAPI)
}

func TestOpenapiToApiSpec(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	content, err := os.ReadFile("./mocks/github_mock.yml")
	if err != nil {
		return
	}
	openAPI, err := parseOpenAPI(ctx, content)
	assert.Nil(t, err)
	if err != nil {
		return
	}
	asd := openapiToApiSpec(log, openAPI)
	assert.NotNil(t, asd)

	groupMap := make(map[string]*apispecdoc.Group, len(asd.Groups))
	for _, group := range asd.Groups {
		groupMap[group.Name] = group
	}

	//***Check "/admin/hooks" path element data ***************************************
	entAdmin, ok := groupMap["enterprise-admin"]
	assert.True(t, ok)
	assert.NotNil(t, entAdmin)
	assert.Equal(t, 2, len(entAdmin.Methods))
	//check get request
	getM := entAdmin.FindMethod(apispecdoc.MethodGet)
	assert.NotNil(t, getM)
	assert.Equal(t, apispecdoc.MethodGet, getM.Type)
	assert.Equal(t, "/admin/hooks", getM.Path)
	assert.NotNil(t, getM.Parameters)
	assert.Equal(t, 2, len(getM.Parameters))
	var perPageParam, pageParam *apispecdoc.Parameter
	for _, par := range getM.Parameters {
		switch par.Name {
		case "per_page":
			perPageParam = par
		case "page":
			pageParam = par
		}
	}
	assert.NotNil(t, perPageParam)
	assert.Equal(t, apispecdoc.Integer, perPageParam.Schema.Type)
	assert.Equal(t, apispecdoc.ParameterQuery, perPageParam.In)

	assert.NotNil(t, pageParam)
	assert.Equal(t, apispecdoc.Integer, pageParam.Schema.Type)
	assert.Equal(t, apispecdoc.ParameterQuery, perPageParam.In)

	assert.Nil(t, getM.RequestBody)
	assert.NotNil(t, getM.ExternalDoc)
	assert.NotNil(t, getM.Servers)

	//check post request
	postM := entAdmin.FindMethod(apispecdoc.MethodPost)
	assert.NotNil(t, postM)
	assert.Equal(t, apispecdoc.MethodPost, postM.Type)
	assert.Equal(t, "/admin/hooks", postM.Path)
	assert.Equal(t, 0, len(postM.Parameters))
	assert.NotNil(t, postM.RequestBody)

	pmBody := postM.RequestBody
	assert.NotNil(t, pmBody.Content)
	assert.True(t, pmBody.Required)
	pmContent := pmBody.FindContentByMediaType("application/json")
	assert.NotNil(t, pmContent)
	assert.NotNil(t, pmContent.Schema)
	pmSchema := pmContent.Schema
	assert.Equal(t, apispecdoc.Object, pmSchema.Type)

	assert.Equal(t, 4, len(pmSchema.Fields))
	nameField := pmSchema.FindField("name")
	assert.NotNil(t, nameField)
	assert.Equal(t, apispecdoc.String, nameField.Type)
	assert.NotEmpty(t, nameField.Description)

	configField := pmSchema.FindField("config")
	assert.NotNil(t, configField)
	assert.Equal(t, apispecdoc.Object, configField.Type)
	assert.NotEmpty(t, configField.Description)
	assert.NotNil(t, configField.FindField("url"))
	assert.NotNil(t, configField.FindField("content_type"))
	assert.NotNil(t, configField.FindField("secret"))
	assert.NotNil(t, configField.FindField("insecure_ssl"))

	eventsField := pmSchema.FindField("events")
	assert.NotNil(t, eventsField)
	assert.Equal(t, apispecdoc.Array, eventsField.Type)
	assert.NotEmpty(t, eventsField.Description)
	assert.NotNil(t, eventsField.Fields)
	assert.Equal(t, 1, len(eventsField.Fields))
	assert.Equal(t, apispecdoc.String, eventsField.Fields[0].Type)

	activeField := pmSchema.FindField("active")
	assert.NotNil(t, activeField)
	assert.Equal(t, apispecdoc.Boolean, activeField.Type)
	assert.NotEmpty(t, activeField.Description)
	//*************finish "/admin/hooks" checks *****************************

	//*************start some "/gists" checks *******************************
	gists, ok := groupMap["gists"]
	assert.True(t, ok)
	assert.NotNil(t, gists)
	assert.Equal(t, 2, len(gists.Methods))

	//post gists some checks
	postG := gists.FindMethod(apispecdoc.MethodPost)
	assert.NotNil(t, postG)
	assert.Equal(t, apispecdoc.MethodPost, postG.Type)
	assert.Equal(t, "/gists", postG.Path)

	gBody := postG.RequestBody
	assert.NotNil(t, gBody.Content)
	assert.True(t, gBody.Required)
	gContent := gBody.FindContentByMediaType("application/json")
	assert.NotNil(t, gContent)
	assert.NotNil(t, gContent.Schema)
	gSchema := gContent.Schema
	assert.Equal(t, apispecdoc.Object, gSchema.Type)

	assert.Equal(t, 3, len(gSchema.Fields))
	publicField := gSchema.FindField("public")
	assert.NotNil(t, publicField)
	assert.Equal(t, apispecdoc.OneOf, publicField.Type)
	assert.Equal(t, 2, len(publicField.Fields))

	//this anyOf doesn't have names (kin-openapi implementation), so just iterate them
	var boolField, strField *apispecdoc.Schema
	for _, field := range publicField.Fields {
		if field.Type == apispecdoc.String {
			strField = field
		} else if field.Type == apispecdoc.Boolean {
			boolField = field
		}
	}

	//Both types of fields presented
	assert.NotNil(t, boolField)
	assert.NotNil(t, strField)

	//*************finish some "/gists" checks ******************************
}
