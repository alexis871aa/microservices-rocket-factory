package repository

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
)

type PartRepository interface {
	GetPart(ctx context.Context, uuid string) (model.Part, error)
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
