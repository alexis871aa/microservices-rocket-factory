package order

import (
	"context"

	"github.com/samber/lo"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, orderUUID string) (*model.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.data[orderUUID]
	if !ok {
		return &model.Order{}, model.ErrOrderNotFound
	}

	return lo.ToPtr(converter.OrderToModel(order)), nil
}
