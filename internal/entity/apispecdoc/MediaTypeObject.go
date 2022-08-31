package apispecdoc

type MediaTypeObject struct {
	ID            int `gorm:"primaryKey"`
	RequestBodyID uint
	RequestBody   RequestBody
	SchemaID      uint
	Schema        Schema
}
