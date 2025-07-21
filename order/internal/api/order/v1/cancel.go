package v1

import (
	"context"
	"errors"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	orderV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Заказ не найден",
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Заказ уже оплачен и не может быть отменён",
			}, nil
		case errors.Is(err, model.ErrOrderCancelled):
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Заказ уже отменён",
			}, nil
		default:
			return nil, err
		}
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
