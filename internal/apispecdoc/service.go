package apispecdoc

import (
	"context"
	"github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	"github.com/rog-golang-buddies/api_hub_common/apispecproto"
)

//go:generate mockgen -source=service.go -destination=./mock/service.go -package=apispecdoc
type Service interface {
	Search(context.Context, *apispecproto.SearchRequest) (*apispecproto.SearchResponse, error)
	Get(context.Context, *apispecproto.GetRequest) (*apispecproto.GetResponse, error)
	Save(context.Context, *apispecdoc.ApiSpecDoc) (uint, error)
}
