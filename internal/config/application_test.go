package config

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadConfig(t *testing.T) {
	expUrlRequestQueue := "url_request_queue"
	err := os.Setenv("QUEUE_URL_REQUEST_QUEUE", expUrlRequestQueue)
	assert.Nil(t, err)

	expSaveASDReqQueue := "test_save_from_queue"
	err = os.Setenv("QUEUE_SAVE_ASD_REQUEST_QUEUE", expSaveASDReqQueue)
	assert.Nil(t, err)

	expNotificationQueue := "test_notifications"
	err = os.Setenv("QUEUE_NOTIFICATION_QUEUE", expNotificationQueue)
	assert.Nil(t, err)

	expQueueUrl := "test_url"
	err = os.Setenv("QUEUE_URL", expQueueUrl)
	assert.Nil(t, err)

	expQueueConcurrency := 50
	err = os.Setenv("QUEUE_CONCURRENCY", strconv.Itoa(expQueueConcurrency))
	assert.Nil(t, err)

	expGRPCHost := "test gRPC host"
	err = os.Setenv("GRPC_HOST", expGRPCHost)
	assert.Nil(t, err)

	expGRPCPort := "test gRPC port"
	err = os.Setenv("GRPC_PORT", expGRPCPort)
	assert.Nil(t, err)

	conf, err := ReadConfig()
	assert.Nil(t, err)
	assert.Equal(t, expSaveASDReqQueue, conf.Queue.SaveASDRequestQueue)
	assert.Equal(t, expNotificationQueue, conf.Queue.NotificationQueue)
	assert.Equal(t, expQueueUrl, conf.Queue.Url)
	assert.Equal(t, expQueueConcurrency, conf.Queue.Concurrency)
	assert.Equal(t, expGRPCHost, conf.GRPC.Host)
	assert.Equal(t, expGRPCPort, conf.GRPC.Port)
}
