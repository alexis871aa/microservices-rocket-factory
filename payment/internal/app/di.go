package app

import (
	"context"

	paymentV1API "github.com/alexis871aa/microservices-rocket-factory/payment/internal/api/payment/v1"
	"github.com/alexis871aa/microservices-rocket-factory/payment/internal/service"
	paymentService "github.com/alexis871aa/microservices-rocket-factory/payment/internal/service/payment"
	paymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	paymentV1API   paymentV1.PaymentServiceServer
	paymentService service.PaymentService
}

func NewDiContainer() *diContainer { return &diContainer{} }

func (d *diContainer) PaymentV1API(ctx context.Context) paymentV1.PaymentServiceServer {
	if d.paymentV1API == nil {
		d.paymentV1API = paymentV1API.NewAPI(d.PaymentService(ctx))
	}

	return d.paymentV1API
}

func (d *diContainer) PaymentService(_ context.Context) service.PaymentService {
	if d.paymentService == nil {
		d.paymentService = paymentService.NewService()
	}

	return d.paymentService
}
