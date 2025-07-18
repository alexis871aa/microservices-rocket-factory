package service

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
)

type PartService interface {
	GetPart(ctx context.Context, uuid string) (model.PartInfo, error)
	ListParts(ctx context.Context, filter model.PartsFilter) (model.PartsInfoFilter, error)
}
