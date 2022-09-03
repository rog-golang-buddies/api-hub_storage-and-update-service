package apispecdoc

type ExternalDoc struct {
	ID          int `gorm:"primaryKey"`
	Description string
	URL         string
	ApiMethodID *uint
}
