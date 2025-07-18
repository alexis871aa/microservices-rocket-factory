package v1

import (
	"context"

	generatedPaymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod generatedPaymentV1.PaymentMethod) (string, error) {
	resp, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: paymentMethod,
	})
	if err != nil {
		return "", err
	}

	return resp.TransactionUuid, nil
}
