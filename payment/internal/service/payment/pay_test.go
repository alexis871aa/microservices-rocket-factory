package payment

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexis871aa/microservices-rocket-factory/payment/internal/service"
)

func TestService_PayOrder(t *testing.T) {
	s := NewService()

	t.Run("success_card_payment", func(t *testing.T) {
		ctx := context.Background()
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := service.PaymentMethodCard

		transactionUUID, err := s.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.NotEmpty(t, transactionUUID)
		assert.Len(t, transactionUUID, 36)
		assert.Contains(t, transactionUUID, "-")
	})

	t.Run("success_sbp_payment", func(t *testing.T) {
		ctx := context.Background()
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := service.PaymentMethodSBP

		transactionUUID, err := s.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.NotEmpty(t, transactionUUID)
		assert.Len(t, transactionUUID, 36)
		assert.Contains(t, transactionUUID, "-")
	})

	t.Run("success_credit_card_payment", func(t *testing.T) {
		ctx := context.Background()
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := service.PaymentMethodCreditCard

		transactionUUID, err := s.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.NotEmpty(t, transactionUUID)
		assert.Len(t, transactionUUID, 36)
		assert.Contains(t, transactionUUID, "-")
	})

	t.Run("success_investor_money_payment", func(t *testing.T) {
		ctx := context.Background()
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := service.PaymentMethodInvestorMoney

		transactionUUID, err := s.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.NotEmpty(t, transactionUUID)
		assert.Len(t, transactionUUID, 36)
		assert.Contains(t, transactionUUID, "-")
	})

	t.Run("success_unknown_payment_method", func(t *testing.T) {
		ctx := context.Background()
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := service.PaymentMethodUnknown

		transactionUUID, err := s.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.NotEmpty(t, transactionUUID)
		assert.Len(t, transactionUUID, 36)
		assert.Contains(t, transactionUUID, "-")
	})

	t.Run("success_with_empty_order_uuid", func(t *testing.T) {
		ctx := context.Background()
		orderUUID := ""
		userUUID := gofakeit.UUID()
		paymentMethod := service.PaymentMethodCard

		transactionUUID, err := s.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.NotEmpty(t, transactionUUID)
		assert.Len(t, transactionUUID, 36)
		assert.Contains(t, transactionUUID, "-")
	})

	t.Run("success_with_empty_user_uuid", func(t *testing.T) {
		ctx := context.Background()
		orderUUID := gofakeit.UUID()
		userUUID := ""
		paymentMethod := service.PaymentMethodCard

		transactionUUID, err := s.PayOrder(ctx, orderUUID, userUUID, paymentMethod)

		require.NoError(t, err)
		assert.NotEmpty(t, transactionUUID)
		assert.Len(t, transactionUUID, 36)
		assert.Contains(t, transactionUUID, "-")
	})

	t.Run("unique_transaction_uuids", func(t *testing.T) {
		ctx := context.Background()
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := service.PaymentMethodCard

		var transactionUUIDs []string
		for i := 0; i < 5; i++ {
			transactionUUID, err := s.PayOrder(ctx, orderUUID, userUUID, paymentMethod)
			require.NoError(t, err)
			transactionUUIDs = append(transactionUUIDs, transactionUUID)
		}

		uniqueUUIDs := make(map[string]bool)
		for _, uuid := range transactionUUIDs {
			assert.False(t, uniqueUUIDs[uuid], "UUID должен быть уникальным: %s", uuid)
			uniqueUUIDs[uuid] = true
		}
	})
}
