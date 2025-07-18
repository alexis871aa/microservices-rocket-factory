package grpc

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	generatedPaymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, orderUUID, userUUID string, paymentMethod generatedPaymentV1.PaymentMethod) (string, error)
}
