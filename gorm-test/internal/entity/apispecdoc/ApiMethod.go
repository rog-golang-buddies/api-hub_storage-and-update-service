package apispecdoc

type ApiMethod struct {
	ID               int `gorm:"primaryKey"`
	Path             string
	Name             string
	Description      string
	Type             string
	ApiSpecDocID     uint
	ApiSpecDocEntity ApiSpecDocEntity
	GroupID          uint
	Group            Group
}
