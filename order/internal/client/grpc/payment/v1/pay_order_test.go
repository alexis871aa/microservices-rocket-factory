package v1

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	paymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
	"github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1/mocks"
)

func TestClient_PayOrder(t *testing.T) {
	t.Run("success_card_payment", func(t *testing.T) {
		ctx := context.Background()
		mockClient := mocks.NewPaymentServiceClient(t)
		client := NewClient(mockClient)

		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := model.PaymentMethodCard
		expectedTransactionUUID := gofakeit.UUID()

		expectedRequest := &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_CARD,
		}

		expectedResponse := &paymentV1.PayOrderResponse{
			TransactionUuid: expectedTransactionUUID,
		}

		mockClient.EXPECT().
			PayOrder(ctx, expectedRequest).
			Return(expectedResponse, nil).
			Once()

		result, err := client.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.Equal(t, expectedTransactionUUID, result)
	})

	t.Run("success_sbp_payment", func(t *testing.T) {
		ctx := context.Background()
		mockClient := mocks.NewPaymentServiceClient(t)
		client := NewClient(mockClient)

		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := model.PaymentMethodSBP
		expectedTransactionUUID := gofakeit.UUID()

		expectedRequest := &paymentV1.PayOrderRequest{
			OrderUuid:     orderUUID,
			UserUuid:      userUUID,
			PaymentMethod: paymentV1.PaymentMethod_SBP,
		}

		expectedResponse := &paymentV1.PayOrderResponse{
			TransactionUuid: expectedTransactionUUID,
		}

		mockClient.EXPECT().
			PayOrder(ctx, expectedRequest).
			Return(expectedResponse, nil).
			Once()

		result, err := client.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.Equal(t, expectedTransactionUUID, result)
	})

	t.Run("success_credit_card_payment", func(t *testing.T) {
		ctx := context.Background()
		mockClient := mocks.NewPaymentServiceClient(t)
		client := NewClient(mockClient)

		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := model.PaymentMethodCreditCard
		expectedTransactionUUID := gofakeit.UUID()

		mockClient.EXPECT().
			PayOrder(ctx, &paymentV1.PayOrderRequest{
				OrderUuid:     orderUUID,
				UserUuid:      userUUID,
				PaymentMethod: paymentV1.PaymentMethod_CREDIT_CARD,
			}).
			Return(&paymentV1.PayOrderResponse{
				TransactionUuid: expectedTransactionUUID,
			}, nil).
			Once()

		result, err := client.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.Equal(t, expectedTransactionUUID, result)
	})

	t.Run("success_investor_money_payment", func(t *testing.T) {
		ctx := context.Background()
		mockClient := mocks.NewPaymentServiceClient(t)
		client := NewClient(mockClient)

		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := model.PaymentMethodInvestorMoney
		expectedTransactionUUID := gofakeit.UUID()

		mockClient.EXPECT().
			PayOrder(ctx, &paymentV1.PayOrderRequest{
				OrderUuid:     orderUUID,
				UserUuid:      userUUID,
				PaymentMethod: paymentV1.PaymentMethod_INVESTOR_MONEY,
			}).
			Return(&paymentV1.PayOrderResponse{
				TransactionUuid: expectedTransactionUUID,
			}, nil).
			Once()

		result, err := client.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.Equal(t, expectedTransactionUUID, result)
	})

	t.Run("unknown_payment_method", func(t *testing.T) {
		ctx := context.Background()
		mockClient := mocks.NewPaymentServiceClient(t)
		client := NewClient(mockClient)

		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := model.PaymentMethodUnknown
		expectedTransactionUUID := gofakeit.UUID()

		mockClient.EXPECT().
			PayOrder(ctx, &paymentV1.PayOrderRequest{
				OrderUuid:     orderUUID,
				UserUuid:      userUUID,
				PaymentMethod: paymentV1.PaymentMethod_UNKNOWN,
			}).
			Return(&paymentV1.PayOrderResponse{
				TransactionUuid: expectedTransactionUUID,
			}, nil).
			Once()

		result, err := client.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.Equal(t, expectedTransactionUUID, result)
	})

	t.Run("grpc_error", func(t *testing.T) {
		ctx := context.Background()
		mockClient := mocks.NewPaymentServiceClient(t)
		client := NewClient(mockClient)

		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := model.PaymentMethodCard
		grpcErr := gofakeit.Error()

		mockClient.EXPECT().
			PayOrder(ctx, &paymentV1.PayOrderRequest{
				OrderUuid:     orderUUID,
				UserUuid:      userUUID,
				PaymentMethod: paymentV1.PaymentMethod_CARD,
			}).
			Return(nil, grpcErr).
			Once()

		result, err := client.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.Error(t, err)
		assert.ErrorIs(t, err, grpcErr)
		assert.Empty(t, result)
	})
}
