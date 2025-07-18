package v1

import (
	def "github.com/alexis871aa/microservices-rocket-factory/order/client/grpc"
	geteratedInventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient geteratedInventoryV1.InventoryServiceClient
}

func NewClient(generatedClient geteratedInventoryV1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
