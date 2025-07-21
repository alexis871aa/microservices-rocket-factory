package v1

import (
	"context"
	"errors"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/converter"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	orderV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/openapi/order/v1"
)

func (a *api) PaymentOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PaymentOrderParams) (orderV1.PaymentOrderRes, error) {
	transactionUUID, err := a.orderService.Pay(ctx, params.OrderUUID, converter.PaymentMethodToModel(req.PaymentMethod))
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
				Message: "Заказ уже оплачен",
			}, nil
		case errors.Is(err, model.ErrOrderCancelled):
			return &orderV1.ConflictError{
				Code:    409,
				Message: "Нельзя оплатить отмененный заказ",
			}, nil
		default:
			return nil, err
		}
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}, nil
}
