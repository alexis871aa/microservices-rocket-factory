package part

import (
	def "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	partRepository def.InventoryRepository
}

func NewService(partRepository def.InventoryRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}
