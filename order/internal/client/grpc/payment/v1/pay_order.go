package v1

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	generatedPaymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

// convertPaymentMethodToProto конвертирует domain PaymentMethod в proto PaymentMethod
func convertPaymentMethodToProto(method model.PaymentMethod) generatedPaymentV1.PaymentMethod {
	switch method {
	case model.PaymentMethodCard:
		return generatedPaymentV1.PaymentMethod_CARD
	case model.PaymentMethodSBP:
		return generatedPaymentV1.PaymentMethod_SBP
	case model.PaymentMethodCreditCard:
		return generatedPaymentV1.PaymentMethod_CREDIT_CARD
	case model.PaymentMethodInvestorMoney:
		return generatedPaymentV1.PaymentMethod_INVESTOR_MONEY
	default:
		return generatedPaymentV1.PaymentMethod_UNKNOWN
	}
}

func (c *client) PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod model.PaymentMethod) (string, error) {
	resp, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{
		OrderUuid:     orderUUID,
		UserUuid:      userUUID,
		PaymentMethod: convertPaymentMethodToProto(paymentMethod),
	})
	if err != nil {
		return "", err
	}

	return resp.TransactionUuid, nil
}
