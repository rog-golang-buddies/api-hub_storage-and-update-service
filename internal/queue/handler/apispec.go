package handler

import (
	"context"
	"encoding/json"

	_ "github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	"github.com/rog-golang-buddies/dto"
	"github.com/rog-golang-buddies/internal/config"
	"github.com/rog-golang-buddies/internal/logger"
	"github.com/rog-golang-buddies/internal/process"
	"github.com/rog-golang-buddies/internal/queue/publisher"
	"github.com/wagslane/go-rabbitmq"
)

type ApiSpecDocHandler struct {
	publisher publisher.Publisher
	config    config.QueueConfig
	processor process.UrlProcessor
	log       logger.Logger
}

func (asdh *ApiSpecDocHandler) Handle(ctx context.Context, delivery rabbitmq.Delivery) rabbitmq.Action {
	asdh.log.Infof("consumed: %v", string(delivery.Body))
	//call process here
	var req dto.ScrapingResult
	err := json.Unmarshal(delivery.Body, &req) // general change here
	if err != nil {
		asdh.log.Errorf("error unmarshalling message: '%v', err: %s", string(delivery.Body), err)
		if req.IsNotifyUser {
			procErr := dto.NewProcessingError(
				"internal unmarshalling problem occurred; probably incompatible model versions", err.Error())
			err = asdh.publish(&delivery, dto.NewUserNotification(&procErr), asdh.config.NotificationQueue)
			if err != nil {
				asdh.log.Error("error while notifying user")
			}
		}
		return rabbitmq.NackDiscard
	}

	//here processing of the request happens...

	// asd, err := asdh.processor.Process(ctx, req.FileUrl)
	// if err != nil {
	// 	asdh.log.Error("error while processing url: ", err)
	// 	if req.IsNotifyUser {
	// 		procErr := dto.NewProcessingError("error while processing url", err.Error())
	// 		err = asdh.publish(&delivery, dto.NewUserNotification(&procErr), asdh.config.NotificationQueue)
	// 		if err != nil {
	// 			asdh.log.Error("error while notifying user")
	// 		}
	// 	}
	// 	return rabbitmq.NackDiscard
	// }

	//publish to the required queue success or error
	result := dto.ScrapingResult{IsNotifyUser: req.IsNotifyUser, ApiSpecDoc: req.ApiSpecDoc}
	err = asdh.publish(&delivery, result, asdh.config.ScrapingResultQueue)
	if err != nil {
		asdh.log.Error("error while publishing: ", err)
		//Here is some error while publishing happened - probably something wrong with the queue
		return rabbitmq.NackDiscard
	}
	if req.IsNotifyUser {
		err = asdh.publish(&delivery, dto.NewUserNotification(nil), asdh.config.NotificationQueue)
		if err != nil {
			asdh.log.Error("error while notifying user")
			//don't discard this message because it was published to the storage service successfully
		}
	}
	asdh.log.Info("url scraped successfully")
	return rabbitmq.Ack
}

func (asdh *ApiSpecDocHandler) publish(delivery *rabbitmq.Delivery, message any, queue string) error {
	content, err := json.Marshal(message)
	if err != nil {
		asdh.log.Info("error while marshalling: ", err)
		return err
	}
	return asdh.publisher.Publish(content,
		[]string{queue},
		rabbitmq.WithPublishOptionsCorrelationID(delivery.CorrelationId),
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
}

func NewApiSpecDocHandler(publisher publisher.Publisher,
	config config.QueueConfig,
	processor process.UrlProcessor,
	log logger.Logger) Handler {
	return &ApiSpecDocHandler{
		publisher: publisher,
		config:    config,
		processor: processor,
		log:       log,
	}
}
