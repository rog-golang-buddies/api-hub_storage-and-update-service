package handler

import (
	"context"
	"encoding/json"

	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/dto"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/logger"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/queue/publisher"
	_ "github.com/rog-golang-buddies/api_hub_common/apispecdoc"
	"github.com/wagslane/go-rabbitmq"
)

type ApiSpecDocHandler struct {
	publisher publisher.Publisher
	config    config.QueueConfig
	log       logger.Logger
}

func (asdh *ApiSpecDocHandler) Handle(ctx context.Context, delivery rabbitmq.Delivery) rabbitmq.Action {
	asdh.log.Infof("consumed: %v", string(delivery.Body))
	//get link to API from queue and unmarshal json response to req
	var req dto.ScrapingResult
	err := json.Unmarshal(delivery.Body, &req)
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
	log logger.Logger) Handler {
	return &ApiSpecDocHandler{
		publisher: publisher,
		config:    config,
		log:       log,
	}
}
