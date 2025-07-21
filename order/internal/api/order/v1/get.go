package v1

import (
	"context"
	"errors"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/converter"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	orderV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrderById(ctx context.Context, params orderV1.GetOrderByIdParams) (orderV1.GetOrderByIdRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.NotFoundError{
				Code:    404,
				Message: "Заказ не найден",
			}, nil
		}
		return nil, err
	}

	return &orderV1.GetOrderResponse{
		OrderDto: converter.OrderToDTO(order),
	}, nil
}
