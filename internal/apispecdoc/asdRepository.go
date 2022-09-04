package apispecdoc

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

//go:generate mockgen -source=asdRepository.go -destination=./mocks/asdRepository.go
type AsdRepository interface {
	Save(ctx context.Context, asd *ApiSpecDoc) (uint, error)
	Delete(ctx context.Context, asd *ApiSpecDoc) error
	Update(ctx context.Context, asdOld *ApiSpecDoc, asdNew *ApiSpecDoc) error
	FindById(ctx context.Context, id uint) (*ApiSpecDoc, error)
	FindByHash(ctx context.Context, hash string) (*ApiSpecDoc, error)
	SearchShort(ctx context.Context, search string) ([]*ApiSpecDoc, error)
}

type AsdRepositoryImpl struct {
	db *gorm.DB
}

func (r *AsdRepositoryImpl) Save(ctx context.Context, asd *ApiSpecDoc) (uint, error) {
	result := r.db.WithContext(ctx).Create(&asd)
	return asd.ID, result.Error
}

func (*AsdRepositoryImpl) Delete(ctx context.Context, asd *ApiSpecDoc) error {
	return errors.New("not implemented")
}

func (r *AsdRepositoryImpl) Update(ctx context.Context, asd *ApiSpecDoc, asdNew *ApiSpecDoc) error {
	if asd == nil {
		return errors.New("old ASD model must not be null")
	}
	if asdNew == nil {
		return errors.New("new ASD model must not be null")
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("api_spec_doc_id = ?", asd.ID).Delete(&ApiMethod{}).Error; err != nil {
			return err
		}
		if err := tx.Where("api_spec_doc_id = ?", asd.ID).Delete(&Group{}).Error; err != nil {
			return err
		}
		asd.Groups = asdNew.Groups
		asd.ApiMethods = asdNew.ApiMethods
		return tx.Save(&asd).Error
	})
}

func (*AsdRepositoryImpl) FindById(ctx context.Context, id uint) (*ApiSpecDoc, error) {
	return nil, errors.New("not implemented")
}

func (r *AsdRepositoryImpl) FindByHash(ctx context.Context, hash string) (*ApiSpecDoc, error) {
	var specDocs []*ApiSpecDoc
	err := r.db.WithContext(ctx).Where("md5sum = ?", hash).Find(&specDocs).Error
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

func (*AsdRepositoryImpl) SearchShort(ctx context.Context, search string) ([]*ApiSpecDoc, error) {
	return nil, errors.New("not implemented")
}

func NewASDRepository(db *gorm.DB) AsdRepository {
	return &AsdRepositoryImpl{db: db}
}
