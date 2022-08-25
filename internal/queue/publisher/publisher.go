package publisher

import (
	"io"

	"github.com/rog-golang-buddies/internal/config"
	"github.com/rog-golang-buddies/internal/logger"
	"github.com/wagslane/go-rabbitmq"
)

// Publisher is just an interface for the library publisher which doesn't have one.
//
//go:generate mockgen -source=publisher.go -destination=./mocks/publisher.go -package=publisher
type Publisher interface {
	io.Closer
	Publish(
		data []byte,
		routingKeys []string,
		optionFuncs ...func(*rabbitmq.PublishOptions),
	) error
}

// NewPublisher creates a publisher and connects to the rabbit under the hood.
// This method appears to be not testable cause it combines 2 responsibilities: create an instance and connect to a queue.
// I think we may rely on NewPublisher has been already tested in the library.
func NewPublisher(conf config.QueueConfig, log logger.Logger) (Publisher, error) {
	return rabbitmq.NewPublisher(
		conf.Url,
		rabbitmq.Config{},
		rabbitmq.WithPublisherOptionsLogger(log),
	)
}

func ClosePublisher(publisher Publisher, log logger.Logger) {
	log.Info("closing publisher")
	err := publisher.Close()
	if err != nil {
		log.Error("error while closing publisher: ", err)
	}
}
