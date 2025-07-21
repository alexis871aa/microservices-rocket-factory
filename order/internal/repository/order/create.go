package order

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/repository/converter"
)

func (r *repository) Create(_ context.Context, order model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[order.OrderUUID] = converter.ModelToOrder(order)
	return nil
}
