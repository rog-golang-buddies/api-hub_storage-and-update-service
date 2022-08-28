package apispecdoc

import (
	"gorm.io/gorm"
)

type Schema struct {
	gorm.Model
	Key         string
	Type        string
	Description string
	ParentID    *Schema
}
