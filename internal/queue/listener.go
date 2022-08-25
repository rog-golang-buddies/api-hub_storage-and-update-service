package queue

import (
	"context"
	"errors"

	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/queue/handler"
	"github.com/wagslane/go-rabbitmq"
)

// Listener represents the consumer wrapper with the method to start listening for all events for this service
//
//go:generate mockgen -source=listener.go -destination=./mocks/listener.go
type Listener interface {
	//Start listening queues
	Start(
		ctx context.Context,
		consumer Consumer,
		config *config.QueueConfig,
		handler handler.Handler,
	) error
}

type ListenerImpl struct {
}

func (listener *ListenerImpl) Start(
	ctx context.Context,
	consumer Consumer,
	config *config.QueueConfig,
	handler handler.Handler,
) error {
	if consumer == nil {
		return errors.New("queue consumer must not be nil")
	}
	if config == nil {
		return errors.New("configuration must not be nil")
	}

	handl := func(delivery rabbitmq.Delivery) rabbitmq.Action {
		return handler.Handle(ctx, delivery)
	}

	err := consumer.StartConsuming(
		handl,
		config.SaveASDRequestQueue,
		[]string{}, //No binding, consuming with the default exchange directly by queue name
		rabbitmq.WithConsumeOptionsConcurrency(config.Concurrency),
		rabbitmq.WithConsumeOptionsQueueDurable,
	)
	if err != nil {
		return err
	}

	return nil
}

func NewListener() Listener {
	return &ListenerImpl{}
}
