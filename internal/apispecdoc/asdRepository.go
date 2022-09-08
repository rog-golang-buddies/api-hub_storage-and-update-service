package apispecdoc

import (
	"context"
	"errors"
	"fmt"

	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/dto"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

//go:generate mockgen -source=asdRepository.go -destination=./mocks/asdRepository.go
type AsdRepository interface {
	//Save saves new ApiSpecDoc entity to the database
	Save(ctx context.Context, asd *ApiSpecDoc) (uint, error)
	//Delete ApiSpecDoc soft, i.e. update deleted_at field and prevent the record from appearing in the requests
	Delete(ctx context.Context, asd *ApiSpecDoc) error
	//Update ApiSpecDoc by replacing all old nested elements with new ones
	Update(ctx context.Context, asd *ApiSpecDoc) error
	//FindById returns full ApiSpecDoc with all nested elements or nil if a such record does not exist
	FindById(ctx context.Context, id uint) (*ApiSpecDoc, error)
	//FindByHash returns ApiSpecDoc without nested elements or nil if nothing is found
	FindByHash(ctx context.Context, hash string) (*ApiSpecDoc, error)
	//FindByUrl returns ApiSpecDoc without nested elements or nil if nothing is found
	FindByUrl(ctx context.Context, url string) (*ApiSpecDoc, error)
	//SearchShort returns a slice of ApiSpecDoc without nested elements that match search string
	//The search goes by title and url fields
	SearchShort(ctx context.Context, search string, page dto.PageRequest) (dto.Page[*ApiSpecDoc], error)
}

type AsdRepositoryImpl struct {
	db *gorm.DB
}

func (r *AsdRepositoryImpl) Save(ctx context.Context, asd *ApiSpecDoc) (uint, error) {
	result := r.db.WithContext(ctx).Create(&asd)
	return asd.ID, result.Error
}

func (r *AsdRepositoryImpl) Delete(ctx context.Context, asd *ApiSpecDoc) error {
	result := r.db.WithContext(ctx).Delete(&asd)
	return result.Error
}

func (r *AsdRepositoryImpl) Update(ctx context.Context, asd *ApiSpecDoc) error {
	if asd == nil {
		return errors.New("old ASD model must not be null")
	}
	if asd == nil {
		return errors.New("asd model must not be null")
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("api_spec_doc_id = ?", asd.ID).Delete(&ApiMethod{}).Error; err != nil {
			return err
		}
		if err := tx.Where("api_spec_doc_id = ?", asd.ID).Delete(&Group{}).Error; err != nil {
			return err
		}
		return tx.Save(&asd).Error
	})
}

func (r *AsdRepositoryImpl) FindById(ctx context.Context, id uint) (*ApiSpecDoc, error) {
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

func (r *AsdRepositoryImpl) FindByUrl(ctx context.Context, url string) (*ApiSpecDoc, error) {
	var specDocs []*ApiSpecDoc
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

func (r *AsdRepositoryImpl) SearchShort(ctx context.Context, search string, page dto.PageRequest) (dto.Page[*ApiSpecDoc], error) {
	var specDocs dto.Page[*ApiSpecDoc]
	err := r.db.WithContext(ctx).Limit(page.Page).Where("title LIKE ?", "%"+search+"%").Or("url LIKE ?", "%"+search+"%").Find(&specDocs).Error
	if err != nil {
		return dto.Page[*ApiSpecDoc]{}, err
	}
	return specDocs, nil
}

func NewASDRepository(db *gorm.DB) AsdRepository {
	return &AsdRepositoryImpl{db: db}
}
