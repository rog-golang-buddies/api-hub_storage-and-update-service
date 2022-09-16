package grpc

import (
	"context"
	"github.com/golang/mock/gomock"
	asdmock "github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/apispecdoc/mock"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	mock_logger "github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/logger/mocks"
	"github.com/rog-golang-buddies/api_hub_common/apispecproto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestServerStartsAndStops(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	asdService := asdmock.NewMockService(ctrl)

	server := NewASDServer(log, asdService)
	assert.NotNil(t, server)

	conf, err := config.ReadConfig()
	assert.Nil(t, err)
	assert.NotNil(t, conf)

	list, err := NewGRPCListener(&conf.GRPC)
	assert.Nil(t, err)
	assert.NotNil(t, list)

	grpcServer := grpc.NewServer()
	errCh := startServerInternal(grpcServer, server, list)
	assert.NotNil(t, errCh)
	select {
	case err, ok := <-errCh:
		if ok {
			t.Error("channel returns error ", err)
		} else {
			t.Error("chanel was closed")
		}
	case <-time.Tick(time.Millisecond * 100):
	}
	grpcServer.GracefulStop()
	select {
	case err, ok := <-errCh:
		if ok {
			t.Error("server stopped with error: ", err)
		}
	case <-time.Tick(time.Second * 1):
		t.Error("error channel wasn't closed on server shutdown")
	}
}

func TestServerStopsOnContextCancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	log.EXPECT().Info("context done, stopping grpc server...")
	asdService := asdmock.NewMockService(ctrl)

	server := NewASDServer(log, asdService)
	assert.NotNil(t, server)

	conf, err := config.ReadConfig()
	assert.Nil(t, err)
	assert.NotNil(t, conf)

	list, err := NewGRPCListener(&conf.GRPC)
	assert.Nil(t, err)
	assert.NotNil(t, list)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := StartServer(ctx, log, server, list)
	assert.NotNil(t, errCh)
	select {
	case err, ok := <-errCh:
		if ok {
			t.Error("channel returns error ", err)
		} else {
			t.Error("chanel was closed")
		}
	case <-time.Tick(time.Millisecond * 100):
	}
	cancel()
	select {
	case err, ok := <-errCh:
		if ok {
			t.Error("server stopped with error: ", err)
		}
	case <-time.Tick(time.Second * 1):
		t.Error("error channel wasn't closed on server shutdown")
	}
}

func TestApiSpecDocServerImpl_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	asdService := asdmock.NewMockService(ctrl)
	log := mock_logger.NewMockLogger(ctrl)
	log.EXPECT().Info(gomock.Any()).Times(1)

	server := NewASDServer(log, asdService)
	ctx := context.Background()
	var id uint32 = 54
	expResp := &apispecproto.GetResponse{ApiSpecDoc: &apispecproto.FullASD{Id: id}}
	expReq := &apispecproto.GetRequest{Id: id}
	asdService.EXPECT().Get(ctx, expReq).Return(expResp, nil)
	assert.NotNil(t, server)
	resp, err := server.Get(ctx, expReq)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.ApiSpecDoc)
	assert.Equal(t, expResp, resp)
}

func TestApiSpecDocServerImpl_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	asdService := asdmock.NewMockService(ctrl)
	log := mock_logger.NewMockLogger(ctrl)
	log.EXPECT().Info(gomock.Any()).Times(1)
	server := NewASDServer(log, asdService)
	ctx := context.Background()
	expResp := &apispecproto.SearchResponse{ShortSpecDocs: []*apispecproto.ShortASD{
		{
			Id:   54,
			Name: "test search name",
		},
	}}
	expReq := &apispecproto.SearchRequest{Search: "test search"}
	asdService.EXPECT().Search(ctx, expReq).Return(expResp, nil)
	assert.NotNil(t, server)
	resp, err := server.Search(ctx, expReq)
	assert.Nil(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.ShortSpecDocs)
	assert.Equal(t, expResp, resp)
}
