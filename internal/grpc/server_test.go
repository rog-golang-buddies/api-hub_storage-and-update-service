package grpc

import (
	"context"
	"github.com/golang/mock/gomock"
	mock_apispecdoc "github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/apispecdoc/mocks"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	mock_logger "github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/logger/mocks"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestServerStartsAndStops(t *testing.T) {
	ctrl := gomock.NewController(t)
	log := mock_logger.NewMockLogger(ctrl)
	asdService := mock_apispecdoc.NewMockService(ctrl)

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
	asdService := mock_apispecdoc.NewMockService(ctrl)

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
