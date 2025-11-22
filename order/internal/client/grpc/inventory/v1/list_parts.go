package v1

import (
	"context"

	clientConverter "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/converter"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	grpcAuth "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/middleware/grpc"
	geteratedInventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	// Добавляем session UUID в gRPC metadata для передачи в Inventory сервис
	ctx = grpcAuth.ForwardSessionUUIDToGRPC(ctx)

	parts, err := c.generatedClient.ListParts(ctx, &geteratedInventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})
	if err != nil {
		return []model.Part{}, err
	}

	return clientConverter.PartListToModel(parts), nil
}
