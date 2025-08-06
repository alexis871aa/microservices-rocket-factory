package part

import (
	def "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/service"
)

var _ def.PartService = (*service)(nil)

type service struct {
	partRepository def.PartRepository
}

func NewService(partRepository def.PartRepository) *service {
	return &service{
		partRepository: partRepository,
	}
}
