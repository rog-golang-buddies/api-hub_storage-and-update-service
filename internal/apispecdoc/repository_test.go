package apispecdoc

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
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
	id, err := rep.Save(context.Background(), &entity)
	assert.False(t, id == 0)
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
	id, err := rep.Save(context.Background(), &asd)
	assert.NotNil(t, err)
	assert.True(t, id == 0)
}
