package apispecdoc

import (
	"time"

	"gorm.io/gorm"
)

type ApiSpecDocEntity struct {
	gorm.Model
	Title       string
	Description string
	Type        string
	Groups      []*GroupEntity     `gorm:"foreignKey:ApiSpecDocID"`
	ApiMethods  []*ApiMethodEntity `gorm:"foreignKey:ApiSpecDocID"`
	Md5sum      string
	FetchedAt   time.Time
}

func (ApiSpecDocEntity) TableName() string {
	return "api_spec_docs"
}
