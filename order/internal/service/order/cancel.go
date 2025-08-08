package order

import (
	"context"
	"time"

	"github.com/samber/lo"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func statusToError(status model.OrderStatus) error {
	statusToError := map[model.OrderStatus]error{
		model.StatusPaid:      model.ErrOrderAlreadyPaid,
		model.StatusCancelled: model.ErrOrderCancelled,
	}

	if err, ok := statusToError[status]; ok {
		return err
	}

	if status != model.StatusPendingPayment {
		return model.ErrInvalidOrderStatus
	}

	return nil
}

func (s *service) Cancel(ctx context.Context, orderUUID string) error {
	order, err := s.Get(ctx, orderUUID)
	if err != nil {
		return err
	}

	if err := statusToError(order.Status); err != nil {
		return err
	}

	order.Status = model.StatusCancelled
	order.UpdatedAt = lo.ToPtr(time.Now())
	return s.orderRepository.Update(ctx, orderUUID, *order)
}
