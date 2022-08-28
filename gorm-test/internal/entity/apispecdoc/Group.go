package apispecdoc

type Group struct {
	ID               uint `gorm:"primaryKey"`
	Name             string
	Description      string
	ApiSpecDocID     uint
	ApiSpecDocEntity ApiSpecDocEntity
	ApiMethods       []ApiMethod `gorm:"many2many:apimethod_groups"`
}
