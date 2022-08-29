package apispecdoc

type Group struct {
	ID                 uint `gorm:"primaryKey"`
	Name               string
	Description        string
	ApiSpecDocEntityID uint
}
