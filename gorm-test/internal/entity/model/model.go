package model

import (
	"time"
)

type ApiSpecDocEntity struct {
	ID          int       `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Title       string    `gorm:"column:title"`
	Description string    `gorm:"column:description"`
	Type        int       `gorm:"column:type"`
	Md5Sum      string    `gorm:"column:md5sum"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
	FetchedAt   time.Time `gorm:"column:fetched_at"`
}

func (m *ApiSpecDocEntity) TableName() string {
	return "api_spec_doc_entity"
}

type Group struct {
	ID           int    `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Name         string `gorm:"column:name"`
	Description  string `gorm:"column:description"`
	ApiSpecDocID int    `gorm:"column:api_spec_doc_id"`
}

func (m *Group) TableName() string {
	return "group"
}

type ApiMethod struct {
	ID           int    `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Path         string `gorm:"column:path"`
	Name         string `gorm:"column:name"`
	Description  string `gorm:"column:description"`
	Type         string `gorm:"column:type"`
	ApiSpecDocID int    `gorm:"column:api_spec_doc_id"`
	GroupID      int    `gorm:"column:group_id"`
}

func (m *ApiMethod) TableName() string {
	return "api_method"
}

type RequestBody struct {
	ID          int    `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Description string `gorm:"column:description"`
	Required    int    `gorm:"column:required"`
	ApiMethodID int    `gorm:"column:api_method_id"`
}

func (m *RequestBody) TableName() string {
	return "request_body"
}

type MediaTypeObject struct {
	ID            int `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	RequestBodyID int `gorm:"column:request_body_id"`
	SchemaID      int `gorm:"column:schema_id"`
}

func (m *MediaTypeObject) TableName() string {
	return "media_type_object"
}

type Schema struct {
	ID          int    `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Key         string `gorm:"column:key"`
	Type        string `gorm:"column:type"`
	Description string `gorm:"column:description"`
	ParentID    int    `gorm:"column:parent_id"`
}

func (m *Schema) TableName() string {
	return "schema"
}

type Parameter struct {
	ID          int    `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Name        string `gorm:"column:name"`
	In          string `gorm:"column:in"`
	Description string `gorm:"column:description"`
	SchemaID    int    `gorm:"column:schema_id"`
	Required    int    `gorm:"column:required"`
	ApiMethodID int    `gorm:"column:api_method_id"`
}

func (m *Parameter) TableName() string {
	return "parameter"
}

type Server struct {
	ID          int    `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Url         string `gorm:"column:url"`
	Description string `gorm:"column:description"`
	ApiMethodID int    `gorm:"column:api_method_id"`
}

func (m *Server) TableName() string {
	return "server"
}

type ExternalDoc struct {
	ID          int    `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	Description string `gorm:"column:description"`
	Url         string `gorm:"column:url"`
	ApiMethodID int    `gorm:"column:api_method_id"`
}

func (m *ExternalDoc) TableName() string {
	return "external_doc"
}
