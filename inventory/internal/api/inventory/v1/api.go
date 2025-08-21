package v1

import (
	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/service"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	partService service.InventoryService
}

func NewAPI(partService service.InventoryService) *api {
	return &api{
		partService: partService,
	}
}
