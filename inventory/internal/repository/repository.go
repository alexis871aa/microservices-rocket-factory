package repository

import (
	"context"
	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
)

type PartRepository interface {
	Get(ctx context.Context, uuid string) (model.PartInfo, error)
	List(ctx context.Context, filter model.PartsFilter) (model.PartsInfoFilter, error)
}
