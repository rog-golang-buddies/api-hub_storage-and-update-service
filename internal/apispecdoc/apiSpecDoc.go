package apispecdoc

import (
	"time"

	"gorm.io/gorm"
)

type ApiSpecDoc struct {
	gorm.Model
	Title       string
	Description string
	Type        string
	Groups      []*Group
	ApiMethods  []*ApiMethod
	Md5sum      string
	FetchedAt   time.Time
}
