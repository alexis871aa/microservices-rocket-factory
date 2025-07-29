package v1

import (
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/service"
	orderV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/openapi/order/v1"
)

var _ orderV1.Handler = (*api)(nil)

type api struct {
	orderService service.OrderService
}

func NewAPI(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
