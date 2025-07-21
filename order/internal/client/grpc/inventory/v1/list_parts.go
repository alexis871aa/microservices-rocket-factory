package v1

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/client/converter"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	geteratedInventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &geteratedInventoryV1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	})
	if err != nil {
		return []model.Part{}, err
	}

	return converter.PartListToModel(parts), nil
}
