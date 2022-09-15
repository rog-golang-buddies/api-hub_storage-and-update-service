package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/apispecdoc"

	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AsdRepositoryImpl struct {
	db *gorm.DB
}

func (r *AsdRepositoryImpl) Save(ctx context.Context, asd *apispecdoc.ApiSpecDoc) (uint, error) {
	result := r.db.WithContext(ctx).Create(&asd)
	return asd.ID, result.Error
}

func (r *AsdRepositoryImpl) Delete(ctx context.Context, asd *apispecdoc.ApiSpecDoc) error {
	result := r.db.WithContext(ctx).Delete(&asd)
	return result.Error
}

func (r *AsdRepositoryImpl) Update(ctx context.Context, asd *apispecdoc.ApiSpecDoc) error {
	if asd == nil {
		return errors.New("asd model must not be null")
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("api_spec_doc_id = ?", asd.ID).Delete(&apispecdoc.ApiMethod{}).Error; err != nil {
			return err
		}
		if err := tx.Where("api_spec_doc_id = ?", asd.ID).Delete(&apispecdoc.Group{}).Error; err != nil {
			return err
		}
		return tx.Save(&asd).Error
	})
}

func (r *AsdRepositoryImpl) FindById(ctx context.Context, id uint) (*apispecdoc.ApiSpecDoc, error) {
	var specDocs []*apispecdoc.ApiSpecDoc
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

func (r *AsdRepositoryImpl) FindByHash(ctx context.Context, hash string) (*apispecdoc.ApiSpecDoc, error) {
	var specDocs []*apispecdoc.ApiSpecDoc
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

func (r *AsdRepositoryImpl) FindByUrl(ctx context.Context, url string) (*apispecdoc.ApiSpecDoc, error) {
	var specDocs []*apispecdoc.ApiSpecDoc
	err := r.db.WithContext(ctx).Where("url = ?", url).Find(&specDocs).Error
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

func (r *AsdRepositoryImpl) SearchShort(ctx context.Context, search string, page dto.PageRequest) (dto.Page[*apispecdoc.ApiSpecDoc], error) {
	var specDocs []*apispecdoc.ApiSpecDoc
	var count int64
	err := r.db.WithContext(ctx).
		Where("title LIKE ?", "%"+search+"%").Or("url LIKE ?", "%"+search+"%").
		Offset(page.Page * page.PerPage).Limit(page.PerPage).
		Order("id").Find(&specDocs).Error
	if err != nil {
		return dto.Page[*apispecdoc.ApiSpecDoc]{}, err
	}
	err = r.db.WithContext(ctx).Model(&apispecdoc.ApiSpecDoc{}).
		Where("title LIKE ?", "%"+search+"%").Or("url LIKE ?", "%"+search+"%").
		Count(&count).Error
	if err != nil {
		return dto.Page[*apispecdoc.ApiSpecDoc]{}, err
	}
	return dto.Page[*apispecdoc.ApiSpecDoc]{Data: specDocs, Page: page.Page, PerPage: page.PerPage, Total: int(count)}, nil
}

func NewASDRepository(db *gorm.DB) apispecdoc.AsdRepository {
	return &AsdRepositoryImpl{db: db}
}
