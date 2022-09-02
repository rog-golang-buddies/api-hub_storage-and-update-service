package apispecdoc

type ApiMethodEntity struct {
	ID           int `gorm:"primaryKey"`
	Path         string
	Name         string
	Description  string
	Type         string
	Parameters   string
	Servers      []*ServerEntity `gorm:"foreignKey:ApiMethodID"`
	RequestBody  string
	ExternalDoc  *ExternalDocEntity `gorm:"foreignKey:ApiMethodID"`
	GroupID      uint
	ApiSpecDocID uint
}

func (ApiMethodEntity) TableName() string {
	return "api_methods"
}
