package handler

import (
	"context"

	"github.com/wagslane/go-rabbitmq"
)

type Handler interface {
	Handle(ctx context.Context, delivery rabbitmq.Delivery) rabbitmq.Action
}
