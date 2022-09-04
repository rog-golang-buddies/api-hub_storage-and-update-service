package internal

import (
	"context"
	"fmt"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/apispecdoc"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/db"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/grpc"
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
	DB, err := db.ConnectAndMigrate(log, &conf.DB)
	if err != nil {
		log.Error("error while db setup: ", err)
		return 1
	}
	asdRepo := apispecdoc.NewRepository(DB)
	asdServ := apispecdoc.NewService(log, asdRepo)

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

	handl := handler.NewApiSpecDocHandler(pub, conf.Queue, log, asdServ)
	listener := queue.NewListener()
	err = listener.Start(ctx, consumer, &conf.Queue, handl)
	if err != nil {
		log.Error("error while listening queue ", err)
		return 1
	}

	//creating grpc server
	lst, err := grpc.NewGRPCListener(&conf.GRPC)
	if err != nil {
		log.Error("error creating grpc listener: ", err)
		return 1
	}
	asdSrv := grpc.NewASDServer(log, asdServ)
	errCh := grpc.StartServer(ctx, log, asdSrv, lst)

	<-errCh

	log.Info("application stopped gracefully (not)")
	return 0
}
