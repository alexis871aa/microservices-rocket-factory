package order

import (
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/repository"
	def "github.com/alexis871aa/microservices-rocket-factory/order/internal/service"
)

var _ def.OrderService = (*service)(nil)

type service struct {
	orderRepository repository.OrderRepository

	inventoryClient grpc.InventoryClient
	paymentClient   grpc.PaymentClient
}

func NewService(orderRepository repository.OrderRepository, inventoryClient grpc.InventoryClient, paymentClient grpc.PaymentClient) *service {
	return &service{
		orderRepository: orderRepository,
		inventoryClient: inventoryClient,
		paymentClient:   paymentClient,
	}
}
