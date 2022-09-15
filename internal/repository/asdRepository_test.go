package repository

import (
	"context"
	"testing"

	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/apispecdoc"

	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	servG := []*apispecdoc.Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*apispecdoc.Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*apispecdoc.Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := apispecdoc.ApiSpecDoc{
		Title:       "Trello API",
		Description: "API for Trello",
		Type:        "1",
		Md5sum:      "981734bf",
		Url:         "test_url",
		Groups:      groups,
		ApiMethods:  apiMeth,
	}
	rep := AsdRepositoryImpl{db: gDb}
	id, err := rep.Save(context.Background(), &entity)
	assert.Equal(t, id, entity.ID)
	assert.Nil(t, err)
}

func TestSaveErrorOnMultipleApiMethodRelations(t *testing.T) {
	t.Skip("Currently skipping because of bug https://github.com/go-gorm/gorm/issues/5673")
	meth := &apispecdoc.ApiMethod{Path: "test/path", Name: "test name"}
	apiMethodsG := []*apispecdoc.ApiMethod{meth}
	apiMethods := []*apispecdoc.ApiMethod{meth}
	groups := []*apispecdoc.Group{{Name: "test group", Description: "test description", ApiMethods: apiMethodsG}}
	asd := apispecdoc.ApiSpecDoc{Title: "test ASD", Groups: groups, ApiMethods: apiMethods}
	rep := AsdRepositoryImpl{db: gDb}
	id, err := rep.Save(context.Background(), &asd)
	assert.NotNil(t, err)
	assert.Equal(t, id, 0)
}

func TestDelete(t *testing.T) {
	servG := []*apispecdoc.Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*apispecdoc.Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*apispecdoc.Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := apispecdoc.ApiSpecDoc{
		Title:       "Trello API",
		Description: "API for Trello",
		Type:        "1",
		Md5sum:      "pook943",
		Groups:      groups,
		Url:         "www.trello.com",
		ApiMethods:  apiMeth,
	}
	rep := AsdRepositoryImpl{db: gDb}
	id, err := rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	result, err := rep.FindById(context.Background(), entity.ID)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.ID, entity.ID)
	err = rep.Delete(context.Background(), &entity)
	assert.Nil(t, err)
	result, err = rep.FindById(context.Background(), entity.ID)
	assert.Nil(t, err)
	assert.Nil(t, result)
}

func TestUpdate(t *testing.T) {
	servG := []*apispecdoc.Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*apispecdoc.Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*apispecdoc.Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := apispecdoc.ApiSpecDoc{
		Title:       "Trello API",
		Description: "API for Trello",
		Type:        "1",
		Md5sum:      "jjujwadk2",
		Groups:      groups,
		Url:         "wwww.trello.com",
		ApiMethods:  apiMeth,
	}
	rep := AsdRepositoryImpl{db: gDb}
	id, err := rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	servG = []*apispecdoc.Server{{URL: "test google url", Description: "test description Google"}}
	apiMethG = []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups = []*apispecdoc.Group{{Name: "test google", ApiMethods: apiMethG}}
	servs = []*apispecdoc.Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth = []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity = apispecdoc.ApiSpecDoc{
		Title:       "Google API",
		Description: "API for Google",
		Type:        "2",
		Md5sum:      "290384hrfi",
		Groups:      groups,
		Url:         "wwww.google.com",
		ApiMethods:  apiMeth,
	}
	err = rep.Update(context.Background(), &entity)
	assert.Nil(t, err)
	result, err := rep.FindById(context.Background(), entity.ID)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity.ID, result.ID)
	assert.Equal(t, entity.Title, result.Title)
}

func TestUpdateErrorOnNil(t *testing.T) {
	rep := AsdRepositoryImpl{db: gDb}
	err := rep.Update(context.Background(), nil)
	assert.NotNil(t, err)
}

func TestFindById(t *testing.T) {
	servG := []*apispecdoc.Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*apispecdoc.Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*apispecdoc.Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := apispecdoc.ApiSpecDoc{
		Title:       "Trello API",
		Description: "API for Trello",
		Type:        "1",
		Md5sum:      "lkasjhdl343125",
		Groups:      groups,
		Url:         "wwwww.trello.com",
		ApiMethods:  apiMeth,
	}
	rep := AsdRepositoryImpl{db: gDb}
	id, err := rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	servG = []*apispecdoc.Server{{URL: "test google url", Description: "test description Google"}}
	apiMethG = []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups = []*apispecdoc.Group{{Name: "test google", ApiMethods: apiMethG}}
	servs = []*apispecdoc.Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth = []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity = apispecdoc.ApiSpecDoc{
		Title:       "Google API",
		Description: "API for Google",
		Type:        "2",
		Md5sum:      "109238hrfeaslfuh",
		Groups:      groups,
		Url:         "www.google.com",
		ApiMethods:  apiMeth,
	}
	id, err = rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	result, err := rep.FindById(context.Background(), entity.ID)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.ID, entity.ID)
	assert.Equal(t, result.Type, entity.Type)
}

func TestFindByHash(t *testing.T) {
	servG := []*apispecdoc.Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*apispecdoc.Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*apispecdoc.Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := apispecdoc.ApiSpecDoc{
		Title:       "Trello API",
		Description: "API for Trello",
		Type:        "1",
		Md5sum:      "595f44fec1e92a71d3e9e77456ba80d1",
		Groups:      groups,
		Url:         "wwwwww.trello.com",
		ApiMethods:  apiMeth,
	}
	rep := AsdRepositoryImpl{db: gDb}
	id, err := rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	result, err := rep.FindByHash(context.Background(), "595f44fec1e92a71d3e9e77456ba80d1")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Md5sum, entity.Md5sum)
	assert.Equal(t, result.ID, entity.ID)
}

func TestFindByUrl(t *testing.T) {
	servG := []*apispecdoc.Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*apispecdoc.Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*apispecdoc.Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := apispecdoc.ApiSpecDoc{
		Title:       "Trello API",
		Description: "API for Trello",
		Type:        "1",
		Md5sum:      "95f44fec1e92a71d3e9e77456ba80d1",
		Groups:      groups,
		Url:         "wwwwwww.trello.com",
		ApiMethods:  apiMeth,
	}
	rep := AsdRepositoryImpl{db: gDb}
	id, err := rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	result, err := rep.FindByUrl(context.Background(), "wwwwwww.trello.com")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Url, entity.Url)
	assert.Equal(t, result.ID, entity.ID)
}

func TestSearchShort(t *testing.T) {
	servG := []*apispecdoc.Server{{URL: "google.com", Description: "test description Google"}}
	apiMethG := []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups := []*apispecdoc.Group{{Name: "test google", ApiMethods: apiMethG}}
	servs := []*apispecdoc.Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth := []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := apispecdoc.ApiSpecDoc{
		Title:       "Google API",
		Description: "API for Google",
		Type:        "2",
		Md5sum:      "lkjafs871324r",
		Groups:      groups,
		Url:         "wwwww.google.com",
		ApiMethods:  apiMeth,
	}
	rep := AsdRepositoryImpl{db: gDb}
	id, err := rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	servG = []*apispecdoc.Server{{URL: "test google url", Description: "test description Google"}}
	apiMethG = []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups = []*apispecdoc.Group{{Name: "test google", ApiMethods: apiMethG}}
	servs = []*apispecdoc.Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth = []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity = apispecdoc.ApiSpecDoc{
		Title:       "microsoft API",
		Description: "API for microsoft",
		Type:        "2",
		Md5sum:      "asdf422423123jkj",
		Groups:      groups,
		Url:         "www.microsoft.com",
		ApiMethods:  apiMeth,
	}
	id, err = rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	servG = []*apispecdoc.Server{{URL: "test google url", Description: "test description Google"}}
	apiMethG = []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups = []*apispecdoc.Group{{Name: "test google", ApiMethods: apiMethG}}
	servs = []*apispecdoc.Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth = []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity = apispecdoc.ApiSpecDoc{
		Title:       "amazon API",
		Description: "API for amazon",
		Type:        "2",
		Md5sum:      "asdfoqwefjipqwef00",
		Groups:      groups,
		Url:         "www.amazon.com",
		ApiMethods:  apiMeth,
	}
	id, err = rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	servG = []*apispecdoc.Server{{URL: "test google url", Description: "test description Google"}}
	apiMethG = []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups = []*apispecdoc.Group{{Name: "test google", ApiMethods: apiMethG}}
	servs = []*apispecdoc.Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth = []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity = apispecdoc.ApiSpecDoc{
		Title:       "netflix API",
		Description: "API for netflix",
		Type:        "2",
		Md5sum:      "afqweqweqwe11123",
		Groups:      groups,
		Url:         "www.netflix.com",
		ApiMethods:  apiMeth,
	}
	id, err = rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	servG = []*apispecdoc.Server{{URL: "test google url", Description: "test description Google"}}
	apiMethG = []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups = []*apispecdoc.Group{{Name: "test google", ApiMethods: apiMethG}}
	servs = []*apispecdoc.Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth = []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity = apispecdoc.ApiSpecDoc{
		Title:       "apple API",
		Description: "API for apple",
		Type:        "2",
		Md5sum:      "vmmvmvmvfs89304",
		Groups:      groups,
		Url:         "www.apple.com",
		ApiMethods:  apiMeth,
	}
	id, err = rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	servG = []*apispecdoc.Server{{URL: "test google url", Description: "test description Google"}}
	apiMethG = []*apispecdoc.ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups = []*apispecdoc.Group{{Name: "test google", ApiMethods: apiMethG}}
	servs = []*apispecdoc.Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth = []*apispecdoc.ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity = apispecdoc.ApiSpecDoc{
		Title:       "Google 2 API",
		Description: "API for Google 2",
		Type:        "2",
		Md5sum:      "bbbbb6b7bb77b",
		Groups:      groups,
		Url:         "www.Google-Google.com",
		ApiMethods:  apiMeth,
	}
	id, err = rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	assert.Equal(t, id, entity.ID)
	number := dto.PageRequest{Page: 4}
	result, err := rep.SearchShort(context.Background(), "Google", number)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, result.Data[3].Title, entity.Title)
	assert.Equal(t, number.Page, result.Page)
	assert.Equal(t, number.PerPage, result.PerPage)
	assert.GreaterOrEqual(t, result.Total, int64(len(result.Data)))
}
