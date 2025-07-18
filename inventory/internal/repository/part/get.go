package part

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
	repoConverter "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, uuid string) (model.PartInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoPart, ok := r.data[uuid]
	if !ok {
		return model.PartInfo{}, model.ErrPartNotFound
	}

	return model.PartInfo{Part: repoConverter.PartToModel(repoPart)}, nil
}
