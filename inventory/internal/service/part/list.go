package part

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
)

func (s *service) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := s.partRepository.ListParts(ctx, filter)
	if err != nil {
		return []model.Part{}, err
	}

	return parts, nil
}
