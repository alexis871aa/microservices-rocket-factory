package part

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
)

func (s *service) GetPart(ctx context.Context, uuid string) (model.PartInfo, error) {
	part, err := s.partRepository.GetPart(ctx, uuid)
	if err != nil {
		return model.PartInfo{}, err
	}

	return part, nil
}
