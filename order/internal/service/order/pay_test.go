package order

import (
	"errors"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestPay() {
	s.Run("success", func() {
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

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(existingOrder, nil).Once()
		s.paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, paymentMethod).Return(transactionUUID, nil).Once()
		s.orderRepository.On("Update", s.ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(nil).Once()

		result, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), transactionUUID, result)
	})

	s.Run("order_not_found", func() {
		orderUUID := gofakeit.UUID()
		paymentMethod := model.PaymentMethodCard

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

		result, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
		assert.Empty(s.T(), result)
	})

	s.Run("order_already_paid", func() {
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

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(existingOrder, nil).Once()

		result, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrOrderAlreadyPaid)
		assert.Empty(s.T(), result)
	})

	s.Run("order_cancelled", func() {
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := model.PaymentMethodCard

		existingOrder := &model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.StatusCancelled,
		}

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(existingOrder, nil).Once()

		result, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrOrderCancelled)
		assert.Empty(s.T(), result)
	})

	s.Run("payment_client_error", func() {
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		paymentMethod := model.PaymentMethodCard
		paymentErr := errors.New("payment service unavailable")

		existingOrder := &model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.StatusPendingPayment,
		}

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(existingOrder, nil).Once()
		s.paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, paymentMethod).Return("", paymentErr).Once()

		result, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, paymentErr)
		assert.Empty(s.T(), result)
	})

	s.Run("repository_update_error", func() {
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

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(existingOrder, nil).Once()
		s.paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, paymentMethod).Return(transactionUUID, nil).Once()
		s.orderRepository.On("Update", s.ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(repoErr).Once()

		result, err := s.service.Pay(s.ctx, orderUUID, paymentMethod)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, repoErr)
		assert.Empty(s.T(), result)
	})

	s.Run("success_different_payment_methods", func() {
		testCases := []struct {
			name          string
			paymentMethod model.PaymentMethod
		}{
			{"sbp", model.PaymentMethodSBP},
			{"credit_card", model.PaymentMethodCreditCard},
			{"investor_money", model.PaymentMethodInvestorMoney},
		}

		for _, tc := range testCases {
			s.Run(tc.name, func() {
				orderUUID := gofakeit.UUID()
				userUUID := gofakeit.UUID()
				transactionUUID := gofakeit.UUID()

				existingOrder := &model.Order{
					OrderUUID: orderUUID,
					UserUUID:  userUUID,
					Status:    model.StatusPendingPayment,
				}

				s.orderRepository.On("Get", s.ctx, orderUUID).Return(existingOrder, nil).Once()
				s.paymentClient.On("PayOrder", mock.Anything, orderUUID, userUUID, tc.paymentMethod).Return(transactionUUID, nil).Once()
				s.orderRepository.On("Update", s.ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(nil).Once()

				result, err := s.service.Pay(s.ctx, orderUUID, tc.paymentMethod)

				require.NoError(s.T(), err)
				assert.Equal(s.T(), transactionUUID, result)
			})
		}
	})
}
