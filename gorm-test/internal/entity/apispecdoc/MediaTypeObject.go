package apispecdoc

import (
	"gorm.io/gorm"
)

type MediaTypeObject struct {
	gorm.Model
	RequestBodyID []RequestBody
	SchemaID      []Schema
}
