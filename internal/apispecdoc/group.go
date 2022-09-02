package apispecdoc

type GroupEntity struct {
	ID           uint `gorm:"primaryKey"`
	Name         string
	Description  string
	ApiSpecDocID uint
	ApiMethods   []*ApiMethodEntity `gorm:"foreignKey:GroupID"`
}

func (GroupEntity) TableName() string {
	return "groups"
}
