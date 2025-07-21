package order

import (
	"context"
	"errors"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID string) (*model.Order, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return nil, model.ErrOrderNotFound
		}
		return nil, err
	}

	return order, nil
}
