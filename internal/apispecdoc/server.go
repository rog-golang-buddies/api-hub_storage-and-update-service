package apispecdoc

type ServerEntity struct {
	ID           int `gorm:"primaryKey"`
	URL          string
	Description  string
	ApiMethodsID uint
}

func (ServerEntity) TableName() string {
	return "servers"
}
