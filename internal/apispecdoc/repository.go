package apispecdoc

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Repository interface {
	Save(ctx context.Context, asd *ApiSpecDoc) error
	Delete(ctx context.Context, asd *ApiSpecDoc) error
	FindById(ctx context.Context, id uint) (*ApiSpecDoc, error)
	SearchShort(ctx context.Context, search string) ([]*ApiSpecDoc, error)
}

type RepositoryImpl struct {
	db *gorm.DB
}

func (r *RepositoryImpl) Save(ctx context.Context, asd *ApiSpecDoc) error {
	result := r.db.WithContext(ctx).Create(&asd)
	return result.Error
}

func (r *RepositoryImpl) Delete(ctx context.Context, asd *ApiSpecDoc) error {
	result := r.db.WithContext(ctx).Delete(&asd)
	return result.Error
}

func (r *RepositoryImpl) FindById(ctx context.Context, id uint) (*ApiSpecDoc, error) {
	var specDocs []*ApiSpecDoc
	err := r.db.WithContext(ctx).Where("id = ?", id).Preload(clause.Associations).Find(&specDocs).Error
	if err != nil {
		return nil, err
	}
	switch len(specDocs) {
	case 0:
		return nil, nil
	case 1:
		return specDocs[0], nil
	default:
		return nil, fmt.Errorf("incorrect number of results, retrieved: %d", len(specDocs))
	}
}

func (r *RepositoryImpl) SearchShort(ctx context.Context, search string) ([]*ApiSpecDoc, error) {
	var specDocs []*ApiSpecDoc
	err := r.db.WithContext(ctx).Where("title LIKE ?", "%"+search+"%").Find(&specDocs).Error
	//err := r.db.WithContext(ctx).Where("title LIKE ?", "%"+search+"%").Or("url LIKE ?", "%"+search+"%").Find(&specDocs).Error
	if err != nil {
		return nil, err
	}
	return specDocs, nil
}

func NewASDRepository(db *gorm.DB) Repository {
	return &RepositoryImpl{db: db}
}
