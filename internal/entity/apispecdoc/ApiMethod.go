package apispecdoc

type ApiMethodEntity struct {
	ID           int `gorm:"primaryKey"`
	Path         string
	Name         string
	Description  string
	Type         string
	Parameters   string
	Servers      []*ServerEntity
	RequestBody  string
	ExternalDoc  *ExternalDocEntity
	GroupID      uint
	ApiSpecDocID uint
}
