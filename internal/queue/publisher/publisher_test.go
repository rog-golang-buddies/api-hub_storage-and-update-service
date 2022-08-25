package publisher

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_logger "github.com/rog-golang-buddies/internal/logger/mocks"
	publisher "github.com/rog-golang-buddies/queue/publisher/mocks"
)

func TestClosePublisher(t *testing.T) {
	ctrl := gomock.NewController(t)
	pub := publisher.NewMockPublisher(ctrl)
	log := mock_logger.NewMockLogger(ctrl)
	log.EXPECT().Info(gomock.Any())

	pub.EXPECT().Close().Return(nil)
	ClosePublisher(pub, log)
}
