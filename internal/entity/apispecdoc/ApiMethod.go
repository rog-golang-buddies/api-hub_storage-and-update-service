package apispecdoc

type ApiMethodEntity struct {
	ID           int `gorm:"primaryKey"`
	Path         string
	Name         string
	Description  string
	Type         string
	Parameters   []Parameter
	Servers      []Server
	RequestBody  RequestBody
	ExternalDoc  ExternalDoc
	GroupID      uint
	ApiSpecDocID uint
}
