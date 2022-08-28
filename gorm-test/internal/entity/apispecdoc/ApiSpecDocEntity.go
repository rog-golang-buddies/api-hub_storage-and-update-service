package apispecdoc

import (
	"time"

	"gorm.io/gorm"
)

type ApiSpecDocEntity struct {
	gorm.Model
	Title       string
	Description string
	Type        int
	Groups      []Group     `gorm:"many2many:entity_groups"`
	ApiMethods  []ApiMethod `gorm:"many2many:entity_apimethods"`
	Md5sum      string
	FetchedAt   time.Time
}
