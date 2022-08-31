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
	Groups      []*GroupEntity
	ApiMethods  []*ApiMethodEntity
	Md5sum      string
	FetchedAt   time.Time
}
