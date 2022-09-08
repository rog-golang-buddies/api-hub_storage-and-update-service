package apispecdoc

import (
	"context"
	"testing"

	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/dto"
	"github.com/stretchr/testify/assert"
)

func TestSave(t *testing.T) {
	servG := []*Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := ApiSpecDoc{
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
	assert.True(t, id == entity.ID)
	assert.Nil(t, err)
}

func TestSaveErrorOnMultipleApiMethodRelations(t *testing.T) {
	t.Skip("Currently skipping because of bug https://github.com/go-gorm/gorm/issues/5673")
	meth := &ApiMethod{Path: "test/path", Name: "test name"}
	apiMethodsG := []*ApiMethod{meth}
	apiMethods := []*ApiMethod{meth}
	groups := []*Group{{Name: "test group", Description: "test description", ApiMethods: apiMethodsG}}
	asd := ApiSpecDoc{Title: "test ASD", Groups: groups, ApiMethods: apiMethods}
	rep := AsdRepositoryImpl{db: gDb}
	id, err := rep.Save(context.Background(), &asd)
	assert.NotNil(t, err)
	assert.True(t, id == 0)
}

func TestDelete(t *testing.T) {
	servG := []*Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := ApiSpecDoc{
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
	assert.True(t, id == entity.ID)
	result, err := rep.FindById(context.Background(), entity.ID)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.ID == entity.ID)
	err = rep.Delete(context.Background(), &entity)
	assert.Nil(t, err)
	result, err = rep.FindById(context.Background(), entity.ID)
	assert.Nil(t, err)
	assert.Nil(t, result)
}

func TestUpdate(t *testing.T) {
	servG := []*Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := ApiSpecDoc{
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
	assert.True(t, id == entity.ID)
	servG = []*Server{{URL: "test google url", Description: "test description Google"}}
	apiMethG = []*ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups = []*Group{{Name: "test google", ApiMethods: apiMethG}}
	servs = []*Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth = []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity = ApiSpecDoc{
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
	assert.True(t, entity.ID == result.ID)
	assert.True(t, entity.Title == result.Title)
}

func TestFindById(t *testing.T) {
	servG := []*Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := ApiSpecDoc{
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
	assert.True(t, id == entity.ID)
	servG = []*Server{{URL: "test google url", Description: "test description Google"}}
	apiMethG = []*ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups = []*Group{{Name: "test google", ApiMethods: apiMethG}}
	servs = []*Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth = []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity = ApiSpecDoc{
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
	assert.True(t, id == entity.ID)
	result, err := rep.FindById(context.Background(), entity.ID)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.ID == entity.ID)
	assert.True(t, result.Type == entity.Type)
}

func TestFindByHash(t *testing.T) {
	servG := []*Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := ApiSpecDoc{
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
	assert.True(t, id == entity.ID)
	result, err := rep.FindByHash(context.Background(), "595f44fec1e92a71d3e9e77456ba80d1")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Md5sum == entity.Md5sum)
	assert.True(t, result.ID == entity.ID)
}

func TestFindByUrl(t *testing.T) {
	servG := []*Server{{URL: "test gr url", Description: "test description G"}}
	apiMethG := []*ApiMethod{{Path: "test/path", Name: "test name", Servers: servG}}
	groups := []*Group{{Name: "test name", ApiMethods: apiMethG}}
	servs := []*Server{{URL: "test servG", Description: "test description 2"}}
	apiMeth := []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := ApiSpecDoc{
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
	assert.True(t, id == entity.ID)
	result, err := rep.FindByUrl(context.Background(), "wwwwwww.trello.com")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.True(t, result.Url == entity.Url)
	assert.True(t, result.ID == entity.ID)
}

func TestSearchShort(t *testing.T) {
	t.Skip("generics")
	servG := []*Server{{URL: "google.com", Description: "test description Google"}}
	apiMethG := []*ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups := []*Group{{Name: "test google", ApiMethods: apiMethG}}
	servs := []*Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth := []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := ApiSpecDoc{
		Title:       "Google API",
		Description: "API for Google",
		Type:        "2",
		Md5sum:      "lkjafs871324r",
		Groups:      groups,
		Url:         "wwww.google.com",
		ApiMethods:  apiMeth,
	}
	rep := AsdRepositoryImpl{db: gDb}
	id, err := rep.Save(context.Background(), &entity)
	number := dto.PageRequest{Page: 1}
	assert.Nil(t, err)
	assert.True(t, id == entity.ID)
	result, err := rep.SearchShort(context.Background(), "Google", number)
	assert.Nil(t, err)
	assert.NotNil(t, result)
}
