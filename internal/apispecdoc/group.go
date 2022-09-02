package apispecdoc

type GroupEntity struct {
	ID            uint `gorm:"primaryKey"`
	Name          string
	Description   string
	ApiSpecDocsID uint
	ApiMethods    []*ApiMethodEntity `gorm:"foreignKey:GroupsID"`
}

func (GroupEntity) TableName() string {
	return "groups"
}
