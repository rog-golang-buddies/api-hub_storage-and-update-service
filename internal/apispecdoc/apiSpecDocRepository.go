package apispecdoc

import (
	"errors"
	"gorm.io/gorm"
)

type ApiSpecDocRepository interface {
	Save(asd *ApiSpecDocEntity) error
	Delete(asd *ApiSpecDocEntity) error
	FindById(id uint) (*ApiSpecDocEntity, error)
	SearchShort(search string) ([]*ApiSpecDocEntity, error)
}

type ApiSpecDocRepositoryImpl struct {
	db *gorm.DB
}

func (r *ApiSpecDocRepositoryImpl) Save(asd *ApiSpecDocEntity) error {
	result := r.db.Create(&asd)
	return result.Error
}

func (*ApiSpecDocRepositoryImpl) Delete(asd *ApiSpecDocEntity) error {
	return errors.New("not implemented")
}

func (*ApiSpecDocRepositoryImpl) FindById(id uint) (*ApiSpecDocEntity, error) {
	return nil, errors.New("not implemented")
}

func (*ApiSpecDocRepositoryImpl) SearchShort(search string) ([]*ApiSpecDocEntity, error) {
	return nil, errors.New("not implemented")
}

func NewASDRepository(db *gorm.DB) ApiSpecDocRepository {
	return &ApiSpecDocRepositoryImpl{db: db}
}
