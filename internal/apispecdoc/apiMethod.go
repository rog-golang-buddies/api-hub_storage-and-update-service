package apispecdoc

type ApiMethodEntity struct {
	ID            int `gorm:"primaryKey"`
	Path          string
	Name          string
	Description   string
	Type          string
	Parameters    string
	Servers       []*ServerEntity `gorm:"foreignKey:ApiMethodsID"`
	RequestBody   string
	ExternalDoc   *ExternalDocEntity `gorm:"foreignKey:ApiMethodsID"`
	GroupsID      uint
	ApiSpecDocsID uint
}

func (ApiMethodEntity) TableName() string {
	return "api_methods"
}
