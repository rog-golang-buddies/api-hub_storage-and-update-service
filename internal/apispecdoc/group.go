package apispecdoc

type Group struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Description  string
	ApiSpecDocID *uint
	ApiMethods   []*ApiMethod
}
