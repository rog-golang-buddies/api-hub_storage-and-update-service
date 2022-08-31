package apispecdoc

type ServerEntity struct {
	ID          int `gorm:"primaryKey"`
	URL         string
	Description string
	ApiMethodID uint
}
