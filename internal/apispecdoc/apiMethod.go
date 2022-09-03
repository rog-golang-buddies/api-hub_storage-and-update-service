package apispecdoc

type ApiMethod struct {
	ID           int `gorm:"primaryKey"`
	Path         string
	Name         string
	Description  string
	Type         string
	Parameters   string
	Servers      []*Server
	RequestBody  string
	ExternalDoc  *ExternalDoc
	GroupID      *uint
	ApiSpecDocID *uint
}
