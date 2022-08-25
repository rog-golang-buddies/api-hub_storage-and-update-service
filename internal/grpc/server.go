package grpc

import (
	"context"
	"errors"
	"github.com/rog-golang-buddies/api_hub_common/apispecproto"
)

type ApiSpecDocServerImpl struct {
	apispecproto.UnimplementedApiSpecDocServer
}

func (asds *ApiSpecDocServerImpl) Search(ctx context.Context, req *apispecproto.SearchRequest) (*apispecproto.SearchResponse, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func (asds *ApiSpecDocServerImpl) Get(ctx context.Context, req *apispecproto.GetRequest) (*apispecproto.GetResponse, error) {
	//TODO implement me
	return nil, errors.New("not implemented")
}

func NewASDServer() apispecproto.ApiSpecDocServer {
	return &ApiSpecDocServerImpl{}
}
