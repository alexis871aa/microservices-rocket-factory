package repository

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.Order) error
	Get(ctx context.Context, orderUUID string) (*model.Order, error)
	Update(ctx context.Context, orderUUID string, newOrder model.Order) error
	Delete(ctx context.Context, orderUUID string) error
}
