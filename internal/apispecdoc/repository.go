package apispecdoc

import (
	"context"
	"errors"
	"gorm.io/gorm"
)

type Repository interface {
	Save(ctx context.Context, asd *ApiSpecDoc) (uint, error)
	Delete(ctx context.Context, asd *ApiSpecDoc) error
	FindById(ctx context.Context, id uint) (*ApiSpecDoc, error)
	SearchShort(ctx context.Context, search string) ([]*ApiSpecDoc, error)
}

type RepositoryImpl struct {
	db *gorm.DB
}

func (r *RepositoryImpl) Save(ctx context.Context, asd *ApiSpecDoc) (uint, error) {
	result := r.db.WithContext(ctx).Create(&asd)
	return asd.ID, result.Error
}

func (*RepositoryImpl) Delete(ctx context.Context, asd *ApiSpecDoc) error {
	return errors.New("not implemented")
}

func (*RepositoryImpl) FindById(ctx context.Context, id uint) (*ApiSpecDoc, error) {
	return nil, errors.New("not implemented")
}

func (*RepositoryImpl) SearchShort(ctx context.Context, search string) ([]*ApiSpecDoc, error) {
	return nil, errors.New("not implemented")
}

func NewRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{db: db}
}
