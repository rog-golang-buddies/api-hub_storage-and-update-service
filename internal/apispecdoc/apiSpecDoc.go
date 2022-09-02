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
	Groups      []*GroupEntity     `gorm:"foreignKey:ApiSpecDocsID"`
	ApiMethods  []*ApiMethodEntity `gorm:"foreignKey:ApiSpecDocsID"`
	Md5sum      string
	FetchedAt   time.Time
}

func (ApiSpecDocEntity) TableName() string {
	return "api_spec_docs"
}
