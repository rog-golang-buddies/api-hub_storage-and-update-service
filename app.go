package internal

import (
	"context"
	"fmt"

	"github.com/rog-golang-buddies/internal/config"
	"github.com/rog-golang-buddies/internal/load"
	"github.com/rog-golang-buddies/internal/logger"
	"github.com/rog-golang-buddies/internal/parse"
	"github.com/rog-golang-buddies/internal/parse/openapi"
	"github.com/rog-golang-buddies/internal/process"
	"github.com/rog-golang-buddies/internal/queue"
	"github.com/rog-golang-buddies/internal/queue/handler"
	"github.com/rog-golang-buddies/internal/queue/publisher"
	"github.com/rog-golang-buddies/internal/recognize"
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

	proc, err := createDefaultProcessor(log, conf)
	if err != nil {
		log.Error("error while creating processor: ", err)
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

	handl := handler.NewApiSpecDocHandler(pub, conf.Queue, proc, log)
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
func createDefaultProcessor(log logger.Logger, config *config.ApplicationConfig) (process.UrlProcessor, error) {
	recognizer := recognize.NewRecognizer(log)
	parsers := []parse.Parser{openapi.NewOpenApi(log)}
	converter := parse.NewConverter(log, parsers)
	loader := load.NewContentLoader(log, &config.Web)

	return process.NewProcessor(recognizer, converter, loader)
}
