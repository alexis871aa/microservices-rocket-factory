package order

import (
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc"
	def "github.com/alexis871aa/microservices-rocket-factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository     def.OrderRepository
	orderProduceService def.OrderProducerService

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(orderRepository def.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient, orderProduceService def.OrderProducerService) *service {
	return &service{
		orderRepository:     orderRepository,
		inventoryClient:     inventoryClient,
		paymentClient:       paymentClient,
		orderProduceService: orderProduceService,
	}
}
