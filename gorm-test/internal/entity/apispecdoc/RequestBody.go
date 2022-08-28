package apispecdoc

import (
	"gorm.io/gorm"
)

type RequestBody struct {
	gorm.Model
	Description string
	Required    bool
	ApiMethodID ApiMethod
}
