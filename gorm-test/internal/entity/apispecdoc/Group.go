package apispecdoc

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	Name         string
	Description  string
	ApiSpecDocID ApiSpecDocEntity
}
