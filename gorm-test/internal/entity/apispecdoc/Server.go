package apispecdoc

type Server struct {
	ID          int `gorm:"primaryKey"`
	URL         string
	Description string
	ApiMethodID uint
}
