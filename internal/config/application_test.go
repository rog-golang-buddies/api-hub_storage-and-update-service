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

	expScrResQueue := "scraping_res_queue"
	err = os.Setenv("QUEUE_SCRAPING_RESULT_QUEUE", expScrResQueue)
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

	conf, err := ReadConfig()
	assert.Nil(t, err)
	assert.Equal(t, expScrResQueue, conf.Queue.SaveASDRequestQueue)
	assert.Equal(t, expNotificationQueue, conf.Queue.NotificationQueue)
	assert.Equal(t, expQueueUrl, conf.Queue.Url)
	assert.Equal(t, expQueueConcurrency, conf.Queue.Concurrency)
}
