package apispecdoc

type Parameter struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	In          string
	Description string
	Required    bool
	ApiMethodID uint
	SchemaID    uint
}
