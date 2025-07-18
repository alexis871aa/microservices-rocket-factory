package grpc

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
