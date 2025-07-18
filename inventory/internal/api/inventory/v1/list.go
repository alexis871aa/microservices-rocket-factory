package v1

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/converter"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filter := req.GetFilter()

	parts, err := a.partService.ListParts(ctx, *converter.ProtoToPartsFilter(filter))
	if err != nil {
		return &inventoryV1.ListPartsResponse{}, err
	}

	return converter.PartsInfoFilterToProto(&parts), nil
}
