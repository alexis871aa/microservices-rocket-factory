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

func Test_SuccessPayOrder(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard
	transactionUUID := gofakeit.UUID()

	existingOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		PartUuids:       []string{gofakeit.UUID()},
		TotalPrice:      100.0,
		Status:          model.StatusPendingPayment,
		TransactionUUID: nil,
		PaymentMethod:   nil,
		CreatedAt:       lo.ToPtr(time.Now()),
		UpdatedAt:       nil,
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, paymentMethod).Return(transactionUUID, nil).Once()
	orderRepository.On("Update", ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(nil).Once()

	result, err := service.Pay(ctx, orderUUID, paymentMethod)

	require.NoError(t, err)
	assert.Equal(t, transactionUUID, result)
}

func Test_PayErrorWhenOrderNotFound(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard

	orderRepository.On("Get", ctx, orderUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

	result, err := service.Pay(ctx, orderUUID, paymentMethod)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrOrderNotFound)
	assert.Empty(t, result)
}

func Test_PayErrorWhenOrderAlreadyPaid(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard
	existingTransactionUUID := gofakeit.UUID()

	existingOrder := &model.Order{
		OrderUUID:       orderUUID,
		UserUUID:        userUUID,
		Status:          model.StatusPaid,
		TransactionUUID: &existingTransactionUUID,
		PaymentMethod:   lo.ToPtr(model.PaymentMethodSBP),
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()

	result, err := service.Pay(ctx, orderUUID, paymentMethod)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrOrderAlreadyPaid)
	assert.Empty(t, result)
}

func Test_PayErrorWhenOrderCancelled(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusCancelled,
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()

	result, err := service.Pay(ctx, orderUUID, paymentMethod)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrOrderCancelled)
	assert.Empty(t, result)
}

func Test_PayErrorWhenPaymentClientFails(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard
	paymentErr := errors.New("payment service unavailable")

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusPendingPayment,
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, paymentMethod).Return("", paymentErr).Once()

	result, err := service.Pay(ctx, orderUUID, paymentMethod)

	require.Error(t, err)
	assert.ErrorIs(t, err, paymentErr)
	assert.Empty(t, result)
}

func Test_PayErrorWhenRepositoryUpdateFails(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCard
	transactionUUID := gofakeit.UUID()
	repoErr := errors.New("database update failed")

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusPendingPayment,
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, paymentMethod).Return(transactionUUID, nil).Once()
	orderRepository.On("Update", ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(repoErr).Once()

	result, err := service.Pay(ctx, orderUUID, paymentMethod)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Empty(t, result)
}

func Test_SuccessPayOrderWithSBP(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodSBP
	transactionUUID := gofakeit.UUID()

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusPendingPayment,
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, paymentMethod).Return(transactionUUID, nil).Once()
	orderRepository.On("Update", ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(nil).Once()

	result, err := service.Pay(ctx, orderUUID, paymentMethod)

	require.NoError(t, err)
	assert.Equal(t, transactionUUID, result)
}

func Test_SuccessPayOrderWithCreditCard(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodCreditCard
	transactionUUID := gofakeit.UUID()

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusPendingPayment,
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, paymentMethod).Return(transactionUUID, nil).Once()
	orderRepository.On("Update", ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(nil).Once()

	result, err := service.Pay(ctx, orderUUID, paymentMethod)

	require.NoError(t, err)
	assert.Equal(t, transactionUUID, result)
}

func Test_SuccessPayOrderWithInvestorMoney(t *testing.T) {
	ctx := context.Background()
	orderRepository := serviceMocks.NewOrderRepository(t)
	inventoryClient := clientMocks.NewInventoryClient(t)
	paymentClient := clientMocks.NewPaymentClient(t)
	service := NewService(orderRepository, inventoryClient, paymentClient)

	orderUUID := gofakeit.UUID()
	userUUID := gofakeit.UUID()
	paymentMethod := model.PaymentMethodInvestorMoney
	transactionUUID := gofakeit.UUID()

	existingOrder := &model.Order{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.StatusPendingPayment,
	}

	orderRepository.On("Get", ctx, orderUUID).Return(existingOrder, nil).Once()
	paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, paymentMethod).Return(transactionUUID, nil).Once()
	orderRepository.On("Update", ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(nil).Once()

	result, err := service.Pay(ctx, orderUUID, paymentMethod)

	require.NoError(t, err)
	assert.Equal(t, transactionUUID, result)
}
