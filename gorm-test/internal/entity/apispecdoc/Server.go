package apispecdoc

import (
	"gorm.io/gorm"
)

type Server struct {
	gorm.Model
	URL         string
	Description string
	ApiMethodID []ApiMethod
}
