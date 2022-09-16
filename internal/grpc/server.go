package grpc

import (
	"context"
	"fmt"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/apispecdoc"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/logger"
	"github.com/rog-golang-buddies/api_hub_common/apispecproto"
	"google.golang.org/grpc"
	"net"
)

type ApiSpecDocServerImpl struct {
	apispecproto.UnimplementedApiSpecDocServer
	service apispecdoc.Service
	log     logger.Logger
}

func (asds *ApiSpecDocServerImpl) Search(ctx context.Context, req *apispecproto.SearchRequest) (*apispecproto.SearchResponse, error) {
	asds.log.Info("Search: ", req)
	return asds.service.Search(ctx, req)
}

func (asds *ApiSpecDocServerImpl) Get(ctx context.Context, req *apispecproto.GetRequest) (*apispecproto.GetResponse, error) {
	asds.log.Info("Get: ", req)
	return asds.service.Get(ctx, req)
}

func NewASDServer(log logger.Logger, service apispecdoc.Service) apispecproto.ApiSpecDocServer {
	return &ApiSpecDocServerImpl{
		log:     log,
		service: service,
	}
}

func StartServer(ctx context.Context, log logger.Logger, server apispecproto.ApiSpecDocServer, listener net.Listener) chan error {
	grpcServer := grpc.NewServer()
	resCh := make(chan error, 1)
	go func() {
		defer close(resCh)
		defer grpcServer.GracefulStop()
		errCh := startServerInternal(grpcServer, server, listener)
		select {
		case <-ctx.Done():
			log.Info("context done, stopping grpc server...")
		case err, ok := <-errCh:
			if ok {
				log.Error("received error from the grpc server stopping notify and stop...")
				resCh <- err
			} else {
				log.Info("grpc server channel closed")
			}
		}
	}()
	return resCh
}

func startServerInternal(grpcServer *grpc.Server, server apispecproto.ApiSpecDocServer, listener net.Listener) chan error {
	errCh := make(chan error)
	go func() {
		defer close(errCh)
		apispecproto.RegisterApiSpecDocServer(grpcServer, server)
		if err := grpcServer.Serve(listener); err != nil {
			errCh <- err
		}
	}()
	return errCh
}

func NewGRPCListener(conf *config.GRPCConfig) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%s", conf.Host, conf.Port))
}
