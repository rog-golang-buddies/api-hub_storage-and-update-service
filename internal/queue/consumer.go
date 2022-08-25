package queue

import (
	"io"

	"github.com/rog-golang-buddies/internal/config"
	"github.com/rog-golang-buddies/internal/logger"
	"github.com/wagslane/go-rabbitmq"
)

// Consumer is just an interface for the library consumer which doesn't have one.
//
//go:generate mockgen -source=consumer.go -destination=./mocks/consumer.go
type Consumer interface {
	io.Closer
	StartConsuming(
		handler rabbitmq.Handler,
		queue string,
		routingKeys []string,
		optionFuncs ...func(*rabbitmq.ConsumeOptions),
	) error
}

func NewConsumer(conf config.QueueConfig, log logger.Logger) (Consumer, error) {
	consumer, err := rabbitmq.NewConsumer(
		conf.Url,
		rabbitmq.Config{},
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsLogger(log),
	)
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

func CloseConsumer(consumer Consumer, log logger.Logger) {
	log.Info("closing consumer")
	err := consumer.Close()
	if err != nil {
		log.Error("error while closing consumer: ", err)
	}
}
