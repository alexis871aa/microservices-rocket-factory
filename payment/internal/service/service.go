package service

import "context"

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

type PaymentService interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod PaymentMethod) (string, error)
}
