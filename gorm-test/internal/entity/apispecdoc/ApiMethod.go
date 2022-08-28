package apispecdoc

import (
	"gorm.io/gorm"
)

type ApiMethod struct {
	gorm.Model
	Path         string
	Name         string
	Description  string
	Type         string
	ApiSpecDocID *ApiSpecDocEntity
	GroupID      *Group
}
