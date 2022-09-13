package apispecdoc

import (
	"context"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/dto"
)

// AsdPage here represents a fixed type.
// It's just workaround for gomock i.e. it can't generate a mock of interface with generics used in it.
// Issue closed https://github.com/golang/mock/issues/621 - awaiting for gomock version 1.7.0.
// TODO delete on gomock 1.7.0 version released
type AsdPage = dto.Page[*ApiSpecDoc]

//go:generate mockgen -source=repository.go -destination=./mock/repository.go -package=apispecdoc
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
	SearchShort(ctx context.Context, search string, page dto.PageRequest) (AsdPage, error)
}
