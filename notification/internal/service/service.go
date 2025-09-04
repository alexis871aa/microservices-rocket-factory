package service

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/notification/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type TelegramService interface {
	SendShipAssembledNotification(ctx context.Context, event model.ShipAssembled) error
	SendOrderPaidNotification(ctx context.Context, event model.OrderPaid) error
}
