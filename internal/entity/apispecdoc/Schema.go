package apispecdoc

type Schema struct {
	ID               int `gorm:"primaryKey"`
	Key              string
	Type             string
	Description      string
	ParentID         *Schema
	Parameters       []Parameter
	MediaTypeObjects []MediaTypeObject
}
