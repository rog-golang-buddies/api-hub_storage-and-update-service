package apispecdoc

import (
	"context"
	"testing"

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
		Md5sum:      "d1092341234",
		Groups:      groups,
		ApiMethods:  apiMeth,
	}
	rep := RepositoryImpl{db: gDb}
	err := rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
}

func TestSaveErrorOnMultipleApiMethodRelations(t *testing.T) {
	t.Skip("Currently skipping because of bug https://github.com/go-gorm/gorm/issues/5673")
	meth := &ApiMethod{Path: "test/path", Name: "test name"}
	apiMethodsG := []*ApiMethod{meth}
	apiMethods := []*ApiMethod{meth}
	groups := []*Group{{Name: "test group", Description: "test description", ApiMethods: apiMethodsG}}
	asd := ApiSpecDoc{Title: "test ASD", Groups: groups, ApiMethods: apiMethods}
	rep := RepositoryImpl{db: gDb}
	err := rep.Save(context.Background(), &asd)
	assert.NotNil(t, err)
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
		Md5sum:      "d1092341234",
		Groups:      groups,
		ApiMethods:  apiMeth,
	}
	rep := RepositoryImpl{db: gDb}
	err := rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	result, err := rep.FindById(context.Background(), entity.ID)
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, result)
	err = rep.Delete(context.Background(), &entity)
	assert.Nil(t, err)
	result, err = rep.FindById(context.Background(), entity.ID)
	if err != nil {
		t.Error(err)
	}
	assert.Nil(t, result)
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
		Md5sum:      "d1092341234",
		Groups:      groups,
		ApiMethods:  apiMeth,
	}
	rep := RepositoryImpl{db: gDb}
	err := rep.Save(context.Background(), &entity)
	if err != nil {
		t.Error(err)
	}
	servG = []*Server{{URL: "test google url", Description: "test description Google"}}
	apiMethG = []*ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups = []*Group{{Name: "test google", ApiMethods: apiMethG}}
	servs = []*Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth = []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity = ApiSpecDoc{
		Title:       "Google API",
		Description: "API for Google",
		Type:        "2",
		Md5sum:      "i1oj234981",
		Groups:      groups,
		ApiMethods:  apiMeth,
	}
	err = rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	result, _ := rep.FindById(context.Background(), entity.ID)
	assert.NotNil(t, result)
}

func TestSearchShort(t *testing.T) {
	servG := []*Server{{URL: "google.com", Description: "test description Google"}}
	apiMethG := []*ApiMethod{{Path: "test/path", Name: "test Google", Servers: servG}}
	groups := []*Group{{Name: "test google", ApiMethods: apiMethG}}
	servs := []*Server{{URL: "test servG", Description: "test Goggle 2"}}
	apiMeth := []*ApiMethod{{Path: "test2/path", Name: "second test method", Servers: servs}}
	entity := ApiSpecDoc{
		Title:       "Google API",
		Description: "API for Google",
		Type:        "2",
		Md5sum:      "i1oj234981",
		Groups:      groups,
		ApiMethods:  apiMeth,
	}
	rep := RepositoryImpl{db: gDb}
	err := rep.Save(context.Background(), &entity)
	assert.Nil(t, err)
	result, err := rep.SearchShort(context.Background(), "Google")
	if err != nil {
		t.Error(err)
	}
	assert.NotNil(t, result)
}
