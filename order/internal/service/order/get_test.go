package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	clientMocks "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc/mocks"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	serviceMocks "github.com/alexis871aa/microservices-rocket-factory/order/internal/service/mocks"
)

func Test_SuccessGetOrder(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	orderProducerService := serviceMocks.NewOrderProducerService(t)
	service := NewService(orderRepository, inventoryClient, paymentClient, orderProducerService)

	orderUUID := gofakeit.UUID()
	expectedOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        gofakeit.UUID(),
		PartUuids:       []string{gofakeit.UUID(), gofakeit.UUID()},
		TotalPrice:      gofakeit.Float32(),
		TransactionUUID: lo.ToPtr(gofakeit.UUID()),
		PaymentMethod:   lo.ToPtr(model.PaymentMethodCard),
		Status:          model.StatusPaid,
		CreatedAt:       lo.ToPtr(time.Now()),
		UpdatedAt:       lo.ToPtr(time.Now()),
	}

	orderRepository.On("Get", ctx, orderUUID).Return(expectedOrder, nil).Once()

	result, err := service.Get(ctx, orderUUID)

	require.NoError(t, err)
	assert.Equal(t, expectedOrder, result)
}

func Test_GetErrorWhenOrderNotFound(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	orderProducerService := serviceMocks.NewOrderProducerService(t)
	service := NewService(orderRepository, inventoryClient, paymentClient, orderProducerService)

	orderUUID := gofakeit.UUID()

	orderRepository.On("Get", ctx, orderUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

	result, err := service.Get(ctx, orderUUID)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrOrderNotFound)
	assert.Nil(t, result)
}

func Test_GetErrorWhenRepositoryGetFails(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	orderProducerService := serviceMocks.NewOrderProducerService(t)
	service := NewService(orderRepository, inventoryClient, paymentClient, orderProducerService)

	orderUUID := gofakeit.UUID()
	repoErr := errors.New("database connection failed")

	orderRepository.On("Get", ctx, orderUUID).Return(&model.Order{}, repoErr).Once()

	result, err := service.Get(ctx, orderUUID)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Nil(t, result)
}

func Test_GetErrorWhenEmptyOrderUUID(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	orderProducerService := serviceMocks.NewOrderProducerService(t)
	service := NewService(orderRepository, inventoryClient, paymentClient, orderProducerService)

	emptyUUID := ""

	orderRepository.On("Get", ctx, emptyUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

	result, err := service.Get(ctx, emptyUUID)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrOrderNotFound)
	assert.Nil(t, result)
}
