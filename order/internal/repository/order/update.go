package order

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/repository/converter"
)

func (r *repository) Update(_ context.Context, orderUUID string, newOrder model.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[orderUUID]; !ok {
		return model.ErrOrderNotFound
	}
	r.data[orderUUID] = converter.ModelToOrder(newOrder)
	return nil
}
