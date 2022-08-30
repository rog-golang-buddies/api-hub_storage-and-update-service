package apispecdoc

type RequestBody struct {
	ID               int `gorm:"primaryKey"`
	Description      string
	Required         bool
	MediaTypeObjects []MediaTypeObject
	ApiMethodID      uint
}
