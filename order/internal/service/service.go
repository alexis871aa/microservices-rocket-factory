package service

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error)
	Get(ctx context.Context, orderUUID string) (*model.Order, error)
	Cancel(ctx context.Context, orderUUID string) error
	Pay(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error)
}

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, orderUUID string) (*model.Order, error)
	Update(ctx context.Context, orderUUID string, newOrder model.Order) error
	Delete(ctx context.Context, orderUUID string) error
}

type OrderConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceOrderPaid(ctx context.Context, event model.OrderPaid) error
}
