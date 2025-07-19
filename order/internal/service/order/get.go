package order

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID string) (*model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return nil, model.ErrOrderNotFound
	}

	return order, nil
}
