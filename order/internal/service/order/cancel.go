package order

import (
	"context"
	"time"

	"github.com/samber/lo"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	order, err := s.Get(ctx, orderUUID)
	if err != nil {
		return model.ErrOrderNotFound
	}

	switch order.Status {
	case model.StatusPaid:
		return model.ErrOrderAlreadyPaid
	case model.StatusCancelled:
		return model.ErrOrderCancelled
	case model.StatusPendingPayment:
		order.Status = model.StatusCancelled
		order.UpdatedAt = lo.ToPtr(time.Now())
		err = s.orderRepository.Update(ctx, orderUUID, *order)
		if err != nil {
			return err
		}
		return nil
	default:
		return model.ErrInvalidOrderStatus
	}
}
