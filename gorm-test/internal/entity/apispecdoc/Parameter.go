package apispecdoc

import (
	"gorm.io/gorm"
)

type Parameter struct {
	gorm.Model
	Name        string
	In          string
	Description string
	SchemaID    []Schema
	Required    bool
	ApiMethodID []ApiMethod
}
