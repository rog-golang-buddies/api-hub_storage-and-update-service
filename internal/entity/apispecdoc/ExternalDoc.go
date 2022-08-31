package apispecdoc

type ExternalDocEntity struct {
	ID          int `gorm:"primaryKey"`
	Description string
	URL         string
	ApiMethodID uint
}
