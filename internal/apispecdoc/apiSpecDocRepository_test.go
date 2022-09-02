package apispecdoc

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSave(t *testing.T) {
	apiMeth := []*ApiMethodEntity{{Path: "test/path", Name: "test name"}}
	groups := []*GroupEntity{{Name: "test name", ApiMethods: apiMeth}}
	entity := ApiSpecDocEntity{Title: "Trello API", Description: "API for Trello", Type: "1", Md5sum: "d1092341234", Groups: groups}
	rep := ApiSpecDocRepositoryImpl{db: gDb}
	err := rep.Save(&entity)
	assert.Nil(t, err)
}
