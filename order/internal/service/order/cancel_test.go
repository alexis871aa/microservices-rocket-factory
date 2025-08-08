package order

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	clientMocks "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc/mocks"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	serviceMocks "github.com/alexis871aa/microservices-rocket-factory/order/internal/service/mocks"
)

func Test_SuccessCancelPendingOrder(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusPendingPayment,
		CreatedAt: lo.ToPtr(time.Now()),
		UpdatedAt: nil,
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()
	orderRepository.On("Update", ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(nil).Once()

	err := service.Cancel(ctx, orderUUID)

	require.NoError(t, err)
}

func Test_CancelErrorWhenOrderNotFound(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()

	orderRepository.On("Get", ctx, orderUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

	err := service.Cancel(ctx, orderUUID)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrOrderNotFound)
}

func Test_CancelErrorWhenOrderAlreadyPaid(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	transactionUUID := gofakeit.UUID()

	existingOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		Status:          model.StatusPaid,
		TransactionUUID: &transactionUUID,
		PaymentMethod:   lo.ToPtr(model.PaymentMethodCard),
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()

	err := service.Cancel(ctx, orderUUID)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrOrderAlreadyPaid)
}

func Test_CancelErrorWhenOrderAlreadyCancelled(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusCancelled,
		UpdatedAt: lo.ToPtr(time.Now()),
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()

	err := service.Cancel(ctx, orderUUID)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrOrderCancelled)
}

func Test_CancelErrorWhenRepositoryGetFails(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	repoErr := errors.New("database connection failed")

	orderRepository.On("Get", ctx, orderUUID).Return(&model.Order{}, repoErr).Once()

	err := service.Cancel(ctx, orderUUID)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}

func Test_CancelErrorWhenRepositoryUpdateFails(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	repoErr := errors.New("database update failed")

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusPendingPayment,
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()
	orderRepository.On("Update", ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(repoErr).Once()

	err := service.Cancel(ctx, orderUUID)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
}

func Test_CancelErrorWhenEmptyUUID(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	emptyUUID := ""

	orderRepository.On("Get", ctx, emptyUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

	err := service.Cancel(ctx, emptyUUID)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrOrderNotFound)
}
