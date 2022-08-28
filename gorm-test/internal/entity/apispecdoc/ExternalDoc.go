package apispecdoc

import "gorm.io/gorm"

type ExternalDoc struct {
	gorm.Model
	Description string
	URL         string
	ApiMethodID ApiMethod
}
