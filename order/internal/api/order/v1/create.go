package v1

import (
	"context"
	"errors"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	orderV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest, params orderV1.CreateOrderParams) (orderV1.CreateOrderRes, error) {
	order, err := a.orderService.Create(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrPartsNotFound):
			return &orderV1.BadRequestError{
				Code:    400,
				Message: "Не все необходимые детали найдены",
			}, nil
		default:
			return nil, err
		}
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  order.OrderUUID,
		TotalPrice: order.TotalPrice,
	}, nil
}
