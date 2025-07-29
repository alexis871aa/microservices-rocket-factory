package v1

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/payment/internal/service"
	paymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	orderUUID := req.GetOrderUuid()
	userUUID := req.GetUserUuid()
	paymentMethod := req.GetPaymentMethod()

	transactionUUID, err := a.paymentService.PayOrder(ctx, orderUUID, userUUID, service.PaymentMethod(paymentMethod))
	if err != nil {
		return nil, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
