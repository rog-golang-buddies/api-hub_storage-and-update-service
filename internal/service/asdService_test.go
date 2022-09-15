package service

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/apispecdoc"
	asdmock "github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/apispecdoc/mock"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/dto"
	mock_logger "github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/logger/mocks"
	apispecdoc2 "github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	"github.com/rog-golang-buddies/api_hub_common/apispecproto"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
	"time"
)

func TestServiceImpl_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()
	search := "search"
	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
		conf:    &config.PageConfig{MinPerPage: 2},
	}
	var page, perPage int32 = 1, 10
	req := apispecproto.SearchRequest{
		Search:  "search",
		Page:    &page,
		PerPage: &perPage,
	}
	pageReq := dto.PageRequest{
		PerPage: int(perPage),
		Page:    int(page) - 1,
	}
	expAsd := apispecdoc.ApiSpecDoc{
		Model: gorm.Model{
			ID: 12,
		},
		Title:       "test title",
		Description: "test description",
		FetchedAt:   time.Now(),
	}
	pageRes := dto.Page[*apispecdoc.ApiSpecDoc]{
		Data:    []*apispecdoc.ApiSpecDoc{&expAsd},
		Page:    int(page),
		PerPage: int(perPage),
		Total:   5,
	}
	repo.EXPECT().SearchShort(ctx, search, pageReq).Return(pageRes, nil)
	result, err := service.Search(ctx, &req)
	assert.Nil(t, err)
	assert.NotNil(t, result)

	//Check pages
	assert.Equal(t, pageRes.Page, int(result.Page.Current))
	assert.Equal(t, pageRes.PerPage, int(result.Page.PerPage))
	assert.Equal(t, pageRes.Total, int(result.Page.Total))

	//Check that documents transferred
	assert.Equal(t, 1, len(result.ShortSpecDocs))
	resAsd := result.ShortSpecDocs[0]
	assert.Equal(t, expAsd.ID, uint(resAsd.Id))
	assert.Equal(t, expAsd.Title, resAsd.Name)
	assert.Equal(t, expAsd.Description, resAsd.Description)
}

func TestServiceImpl_SearchIncorrectPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()
	search := "search"
	minPage := 2
	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
		conf:    &config.PageConfig{MinPerPage: minPage},
	}
	var page, perPage int32 = 0, 0
	req := apispecproto.SearchRequest{
		Search:  "search",
		Page:    &page,
		PerPage: &perPage,
	}
	pageRes := dto.Page[*apispecdoc.ApiSpecDoc]{
		Data:    []*apispecdoc.ApiSpecDoc{},
		Page:    int(page),
		PerPage: int(perPage),
		Total:   5,
	}
	var pageReq *dto.PageRequest
	repo.EXPECT().SearchShort(ctx, search, gomock.Any()).
		Do(func(ctx context.Context, search string, pr dto.PageRequest) {
			pageReq = &pr
		}).
		Return(pageRes, nil)
	log.EXPECT().Warnf(gomock.Any(), gomock.Any()).Times(2)
	result, err := service.Search(ctx, &req)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, pageReq)
	assert.Equal(t, 0, pageReq.Page)
	assert.Equal(t, minPage, pageReq.PerPage)
}

func TestServiceImpl_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()

	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
	}

	req := apispecproto.GetRequest{Id: 5}
	var docId uint = 54
	var methodId uint = 5
	var groupId uint = 4
	expExtDoc := apispecdoc.ExternalDoc{
		ID:          2,
		Description: "test description 1",
		URL:         "test url 1",
		ApiMethodID: &methodId,
	}
	expMethod := apispecdoc.ApiMethod{
		ID:          methodId,
		Path:        "test path",
		Name:        "meth name",
		Description: "method description",
		Type:        "OPEN_API",
		Parameters:  `[{"Name":"par1","In":"HEADER","Description":"par1 d", "Required":false}]`,
		RequestBody: `{"Description":"description","Required":false,"Content":[{"MediaType":"application/json","Schema":{"Key":"","Type":"INTEGER","Description":"body description"}}]}`,
		Servers: []*apispecdoc.Server{
			{
				ID:          2,
				URL:         "url ",
				Description: "description",
				ApiMethodID: &methodId,
			},
		},
		ExternalDoc:  &expExtDoc,
		ApiSpecDocID: &docId,
	}
	expGrExtDoc := apispecdoc.ExternalDoc{
		ID:          3,
		Description: "test description 2",
		URL:         "test url 3",
		ApiMethodID: &methodId,
	}
	expGrMethod := apispecdoc.ApiMethod{
		ID:          methodId,
		Path:        "test path",
		Name:        "meth name",
		Description: "method description",
		Type:        "test type",
		Servers: []*apispecdoc.Server{
			{
				ID:          2,
				URL:         "url",
				Description: "description",
				ApiMethodID: &methodId,
			},
		},
		ExternalDoc: &expGrExtDoc,
		GroupID:     &groupId,
	}
	expGroup := apispecdoc.Group{
		ID:           4,
		Name:         "test name",
		Description:  "test description",
		ApiSpecDocID: &docId,
		ApiMethods:   []*apispecdoc.ApiMethod{&expGrMethod},
	}
	expAsd := apispecdoc.ApiSpecDoc{
		Model: gorm.Model{
			ID: docId,
		},
		Title:       "test title",
		Description: "test description",
		Type:        "OPEN_API",
		Groups:      []*apispecdoc.Group{&expGroup},
		ApiMethods:  []*apispecdoc.ApiMethod{&expMethod},
		Md5sum:      "test sum",
		Url:         "url",
		FetchedAt:   time.Now(),
	}

	repo.EXPECT().FindById(ctx, uint(req.Id)).Return(&expAsd, nil)
	resAsd, err := service.Get(ctx, &req)
	assert.Nil(t, err)
	assert.NotNil(t, resAsd)

	//Check result root
	fullAsd := resAsd.GetApiSpecDoc()
	assert.Equal(t, expAsd.ID, uint(fullAsd.Id))
	assert.Equal(t, apispecproto.Type_OPEN_API, fullAsd.Type)
	assert.Equal(t, expAsd.Title, fullAsd.Title)
	assert.Equal(t, expAsd.Description, fullAsd.Description)
	assert.Equal(t, expAsd.Url, fullAsd.Url)
	assert.Equal(t, 1, len(fullAsd.Methods))
	assert.Equal(t, 1, len(fullAsd.Groups))
	//Check root method
	resRootMeth := fullAsd.Methods[0]
	assert.Equal(t, expMethod.Path, resRootMeth.Path)
	assert.Equal(t, expMethod.Name, resRootMeth.Name)
	assert.Equal(t, expMethod.Description, resRootMeth.Description)
	assert.Equal(t, 1, len(resRootMeth.Parameters))
	assert.NotNil(t, resRootMeth.RequestBody)
	assert.NotNil(t, resRootMeth.RequestBody.Content)
	assert.Equal(t, 1, len(resRootMeth.RequestBody.Content))
	assert.NotNil(t, resRootMeth.ExternalDoc)
	resExtDoc := resRootMeth.ExternalDoc
	assert.Equal(t, expExtDoc.URL, resExtDoc.Url)
	assert.Equal(t, expExtDoc.Description, resExtDoc.Description)

	resGroup := fullAsd.Groups[0]
	assert.Equal(t, expGroup.Name, resGroup.Name)

	//Check group method
	resGrMeth := resGroup.Methods[0]
	assert.Equal(t, expGrMethod.Name, resGrMeth.Name)
	assert.Equal(t, expGrMethod.Description, resGrMeth.Description)
	assert.Nil(t, resGrMeth.Parameters)
	assert.Nil(t, resGrMeth.RequestBody)
	assert.NotNil(t, resGrMeth.ExternalDoc)
	resGrExtDoc := resGrMeth.ExternalDoc
	assert.Equal(t, expGrExtDoc.URL, resGrExtDoc.Url)
	assert.Equal(t, expGrExtDoc.Description, resGrExtDoc.Description)
}

func TestServiceImpl_GetNilArgError(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()

	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
		conf:    &config.PageConfig{MinPerPage: 2},
	}
	log.EXPECT().Error(gomock.Any())
	res, err := service.Get(ctx, nil)
	assert.Nil(t, res)
	assert.NotNil(t, err)
}

func TestServiceImpl_GetFindByIdError(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()

	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
	}
	log.EXPECT().Error(gomock.Any())
	var reqId uint = 54
	expErr := errors.New("some error on find query")
	repo.EXPECT().FindById(ctx, reqId).Return(nil, expErr)
	res, err := service.Get(ctx, &apispecproto.GetRequest{Id: uint32(reqId)})
	assert.Equal(t, expErr, err)
	assert.Nil(t, res)
}

func TestServiceImpl_GetNotFoundNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()

	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
	}
	log.EXPECT().Infof(gomock.Any(), gomock.Any())
	var reqId uint = 54
	repo.EXPECT().FindById(ctx, reqId).Return(nil, nil)
	res, err := service.Get(ctx, &apispecproto.GetRequest{Id: uint32(reqId)})
	assert.Nil(t, err)
	assert.Nil(t, res)
}

func TestServiceImpl_SaveDuplicateByHashReturnExistingId(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()

	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
	}
	hash := "duplicate hash"
	asdReq := apispecdoc2.ApiSpecDoc{
		Title:       "test title",
		Description: "test description",
		Type:        "OPEN_API",
		Groups:      nil,
		Methods:     nil,
		Md5Sum:      hash,
		Url:         "test url",
	}
	var expId uint = 54
	asdRes := apispecdoc.ApiSpecDoc{
		Model:  gorm.Model{ID: expId},
		Md5sum: hash,
	}
	repo.EXPECT().FindByHash(ctx, hash).Return(&asdRes, nil)
	log.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)
	id, err := service.Save(ctx, &asdReq)
	assert.Equal(t, expId, id)
	assert.Nil(t, err)
}

func TestServiceImpl_SaveRecordHashChangedCallUpdate(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()

	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
	}
	group1 := apispecdoc2.Group{
		Name:        "Group 1",
		Description: "Group 1 d",
	}
	group2 := apispecdoc2.Group{
		Name:        "Group 2",
		Description: "Group 2 d",
	}
	groups := []*apispecdoc2.Group{&group1, &group2}
	method1 := apispecdoc2.ApiMethod{Name: "method 1", Description: "Description 1"}
	method2 := apispecdoc2.ApiMethod{Name: "method 2", Description: "Description 2"}
	methods := []*apispecdoc2.ApiMethod{&method1, &method2}
	newHash := "new hash"
	asdReq := apispecdoc2.ApiSpecDoc{
		Title:       "test title",
		Description: "test description",
		Type:        "OPEN_API",
		Groups:      groups,
		Methods:     methods,
		Md5Sum:      newHash,
		Url:         "test url",
	}
	var expId uint = 54
	oldHash := "old hash"
	expAsdRes := apispecdoc.ApiSpecDoc{
		Model:       gorm.Model{ID: expId},
		Title:       "some title",
		Description: "some description",
		Type:        "UNKNOWN",
		Md5sum:      oldHash,
	}
	repo.EXPECT().FindByHash(ctx, newHash).Return(nil, nil)
	repo.EXPECT().FindByUrl(ctx, asdReq.Url).Return(&expAsdRes, nil)
	repo.EXPECT().Update(ctx, &expAsdRes).Return(nil)
	log.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)
	resId, err := service.Save(ctx, &asdReq)
	assert.Nil(t, err)
	assert.Equal(t, expId, resId)
	assert.Equal(t, expAsdRes.Title, asdReq.Title)
	assert.Equal(t, expAsdRes.Url, asdReq.Url)
	assert.Equal(t, newHash, expAsdRes.Md5sum)
	assert.Equal(t, expAsdRes.Description, asdReq.Description)
	assert.Equal(t, 2, len(expAsdRes.Groups))
	assert.Equal(t, 2, len(expAsdRes.ApiMethods))
}

func TestServiceImpl_Save(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()

	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
	}
	asdReq := apispecdoc2.ApiSpecDoc{
		Title:       "test title",
		Description: "test description",
		Type:        "OPEN_API",
		Md5Sum:      "some request hash",
		Url:         "test url",
		Methods: []*apispecdoc2.ApiMethod{
			{
				Path:        "test path",
				Name:        "meth name",
				Description: "method description",
				Type:        "test type",
				Servers: []*apispecdoc2.Server{
					{
						Url:         "url",
						Description: "description",
					},
				},
			},
		},
		Groups: []*apispecdoc2.Group{
			{
				Name:        "group",
				Description: "gr description",
				Methods: []*apispecdoc2.ApiMethod{
					{
						Path:        "test path",
						Name:        "meth name",
						Description: "method description",
						RequestBody: &apispecdoc2.RequestBody{
							Description: "Body description",
							Content: []*apispecdoc2.MediaTypeObject{
								{
									MediaType: "application/json",
									Schema: &apispecdoc2.Schema{
										Key:         "key",
										Type:        "INTEGER",
										Description: "Field description",
									},
								},
							},
							Required: false,
						},
						Type: "test type",
						Servers: []*apispecdoc2.Server{
							{
								Url:         "url",
								Description: "description",
							},
						},
						ExternalDoc: &apispecdoc2.ExternalDoc{
							Description: "Ext doc description",
							Url:         "ext doc url",
						},
					},
				},
			},
		},
	}
	var expId uint = 78
	var saveArg *apispecdoc.ApiSpecDoc
	repo.EXPECT().FindByHash(ctx, asdReq.Md5Sum).Return(nil, nil)
	repo.EXPECT().FindByUrl(ctx, asdReq.Url).Return(nil, nil)
	repo.EXPECT().Save(ctx, gomock.Any()).Do(func(ctx context.Context, arg *apispecdoc.ApiSpecDoc) {
		saveArg = arg
	}).Return(expId, nil)
	log.EXPECT().Infof(gomock.Any(), gomock.Any()).Times(1)

	asdRes, err := service.Save(ctx, &asdReq)
	assert.Nil(t, err)
	assert.NotNil(t, asdRes)
	assert.NotNil(t, saveArg)
	assert.Equal(t, asdReq.Title, saveArg.Title)
	assert.Equal(t, asdReq.Description, saveArg.Description)
	assert.Equal(t, asdReq.Url, saveArg.Url)
	assert.Equal(t, asdReq.Md5Sum, saveArg.Md5sum)
}

func TestServiceImpl_SaveFindByHashError(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()

	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
	}
	asdReq := apispecdoc2.ApiSpecDoc{
		Title:       "test title",
		Description: "test description",
		Type:        "OPEN_API",
		Md5Sum:      "some request hash",
		Url:         "test url",
	}
	expErr := errors.New("test error")
	repo.EXPECT().FindByHash(ctx, asdReq.Md5Sum).Return(nil, expErr)

	asdRes, err := service.Save(ctx, &asdReq)
	assert.NotNil(t, asdRes)
	assert.Equal(t, expErr, err)
}

func TestServiceImpl_SaveErrorOnNil(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)

	repo := asdmock.NewMockAsdRepository(ctrl)
	ctx := context.Background()

	service := ServiceImpl{
		log:     log,
		asdRepo: repo,
	}
	id, err := service.Save(ctx, nil)
	assert.NotNil(t, err)
	assert.Equal(t, uint(0), id)
}
