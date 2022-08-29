package apispecdoc

type MediaTypeObject struct {
	ID            int `gorm:"primaryKey"`
	RequestBodyID uint
	SchemaID      uint
}
