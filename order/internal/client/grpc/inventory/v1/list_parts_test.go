package v1

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
	"github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1/mocks"
)

func TestClient_ListParts(t *testing.T) {
	t.Run("success_empty_filter", func(t *testing.T) {
		ctx := context.Background()
		mockClient := mocks.NewInventoryServiceClient(t)
		client := NewClient(mockClient)

		filter := model.PartsFilter{}
		expectedResponse := &inventoryV1.ListPartsResponse{
			Parts: []*inventoryV1.Part{},
		}

		mockClient.EXPECT().
			ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{},
			}).
			Return(expectedResponse, nil).
			Once()

		result, err := client.ListParts(ctx, filter)

		require.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("grpc_error", func(t *testing.T) {
		ctx := context.Background()
		mockClient := mocks.NewInventoryServiceClient(t)
		client := NewClient(mockClient)

		filter := model.PartsFilter{}
		grpcErr := gofakeit.Error()

		mockClient.EXPECT().
			ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{},
			}).
			Return(nil, grpcErr).
			Once()

		result, err := client.ListParts(ctx, filter)

		require.Error(t, err)
		assert.ErrorIs(t, err, grpcErr)
		assert.Empty(t, result)
	})

	t.Run("success_with_uuids_filter", func(t *testing.T) {
		ctx := context.Background()
		mockClient := mocks.NewInventoryServiceClient(t)
		client := NewClient(mockClient)

		uuids := []string{gofakeit.UUID(), gofakeit.UUID()}
		filter := model.PartsFilter{
			Uuids: &uuids,
		}

		expectedResponse := &inventoryV1.ListPartsResponse{
			Parts: []*inventoryV1.Part{},
		}

		mockClient.EXPECT().
			ListParts(ctx, &inventoryV1.ListPartsRequest{
				Filter: &inventoryV1.PartsFilter{
					Uuids: uuids,
				},
			}).
			Return(expectedResponse, nil).
			Once()

		result, err := client.ListParts(ctx, filter)

		require.NoError(t, err)
		assert.Empty(t, result)
	})
}
