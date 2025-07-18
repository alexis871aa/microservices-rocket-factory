package part

import (
	"sync"

	def "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository"
	repoModel "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/model"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	mu   sync.RWMutex
	data map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		data: make(map[string]repoModel.Part),
	}
}
