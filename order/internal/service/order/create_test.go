package order

import (
	"context"
	"errors"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	clientMocks "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc/mocks"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	serviceMocks "github.com/alexis871aa/microservices-rocket-factory/order/internal/service/mocks"
)

func Test_SuccessCreateOrder(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	userUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}

	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Name:  "Engine",
			Price: 100.0,
		},
		{
			Uuid:  partUUIDs[1],
			Name:  "Fuel Tank",
			Price: 200.0,
		},
	}

	filter := model.PartsFilter{
		Uuids: &partUUIDs,
	}

	inventoryClient.On("ListParts", mock.Anything, filter).Return(parts, nil).Once()
	orderRepository.On("Create", mock.Anything, mock.AnythingOfType("model.Order")).Return(nil).Once()

	result, err := service.Create(ctx, userUUID, partUUIDs)

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, userUUID, result.UserUUID)
	assert.Equal(t, partUUIDs, result.PartUuids)
	assert.Equal(t, float32(300.0), result.TotalPrice)
	assert.Equal(t, model.StatusPendingPayment, result.Status)
	assert.NotEmpty(t, result.OrderUUID)
	assert.NotNil(t, result.CreatedAt)
}

func Test_CreateErrorWhenPartsNotFound(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	userUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}

	filter := model.PartsFilter{
		Uuids: &partUUIDs,
	}

	inventoryClient.On("ListParts", mock.Anything, filter).Return([]model.Part{}, nil).Once()

	result, err := service.Create(ctx, userUUID, partUUIDs)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrPartsNotFound)
	assert.Nil(t, result)
}

func Test_CreateErrorWhenPartialPartsFound(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	userUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}

	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Name:  "Engine",
			Price: 100.0,
		},
	}

	filter := model.PartsFilter{
		Uuids: &partUUIDs,
	}

	inventoryClient.On("ListParts", mock.Anything, filter).Return(parts, nil).Once()

	result, err := service.Create(ctx, userUUID, partUUIDs)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrPartsNotFound)
	assert.Nil(t, result)
}

func Test_CreateErrorWhenInventoryClientFails(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	userUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID()}
	clientErr := errors.New("inventory service unavailable")

	filter := model.PartsFilter{
		Uuids: &partUUIDs,
	}

	inventoryClient.On("ListParts", mock.Anything, filter).Return([]model.Part{}, clientErr).Once()

	result, err := service.Create(ctx, userUUID, partUUIDs)

	require.Error(t, err)
	assert.ErrorIs(t, err, clientErr)
	assert.Nil(t, result)
}

func Test_CreateErrorWhenRepositoryCreateFails(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	userUUID := gofakeit.UUID()
	partUUIDs := []string{gofakeit.UUID()}
	repoErr := errors.New("database connection failed")

	parts := []model.Part{
		{
			Uuid:  partUUIDs[0],
			Name:  "Engine",
			Price: 100.0,
		},
	}

	filter := model.PartsFilter{
		Uuids: &partUUIDs,
	}

	inventoryClient.On("ListParts", mock.Anything, filter).Return(parts, nil).Once()
	orderRepository.On("Create", mock.Anything, mock.AnythingOfType("model.Order")).Return(repoErr).Once()

	result, err := service.Create(ctx, userUUID, partUUIDs)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, result)
}

func Test_CreateErrorWhenEmptyPartUUIDs(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	userUUID := gofakeit.UUID()
	partUUIDs := []string{}

	result, err := service.Create(ctx, userUUID, partUUIDs)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrPartsNotFound)
	assert.Nil(t, result)
}
