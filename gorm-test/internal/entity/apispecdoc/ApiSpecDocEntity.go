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
	Groups      []Group
	ApiMethods  []ApiMethod
	Md5sum      string
	FetchedAt   time.Time
}
