package internal

import (
	"context"
	"fmt"

	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/logger"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/queue"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/queue/handler"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/queue/publisher"
)

func Start() int {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf, err := config.ReadConfig() //read configuration from env
	if err != nil {
		fmt.Println("error while reading configuration")
		return 1
	}
	log, err := logger.NewLogger(conf)
	if err != nil {
		fmt.Println("error creating logger: ", err)
		return 1
	}

	//initialize publisher connection to the queue
	//this library assumes using one publisher and one consumer per application
	//https://github.com/wagslane/go-rabbitmq/issues/79
	pub, err := publisher.NewPublisher(conf.Queue, log)
	if err != nil {
		log.Error("error while starting publisher: ", err)
		return 1
	}
	defer publisher.ClosePublisher(pub, log)
	//initialize consumer connection to the queue
	consumer, err := queue.NewConsumer(conf.Queue, log)
	if err != nil {
		log.Error("error while connecting to the queue: ", err)
		return 1
	}
	defer queue.CloseConsumer(consumer, log)

	handl := handler.NewApiSpecDocHandler(pub, conf.Queue, log)
	listener := queue.NewListener()
	err = listener.Start(ctx, consumer, &conf.Queue, handl)
	if err != nil {
		log.Error("error while listening queue ", err)
		return 1
	}

	<-ctx.Done()

	log.Info("application stopped gracefully (not)")
	return 0
}
