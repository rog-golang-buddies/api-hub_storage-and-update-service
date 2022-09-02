package apispecdoc

type ExternalDocEntity struct {
	ID          int `gorm:"primaryKey"`
	Description string
	URL         string
	ApiMethodID uint
}

func (ExternalDocEntity) TableName() string {
	return "external_docs"
}
