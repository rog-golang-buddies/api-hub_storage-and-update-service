package apispecdoc

type GroupEntity struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Description  string
	ApiSpecDocID uint
	ApiMethods   []*ApiMethodEntity
}
