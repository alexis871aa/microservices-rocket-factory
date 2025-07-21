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

func (s *ServiceSuite) TestCancel() {
	s.Run("success_pending_payment", func() {
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()

		existingOrder := &model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.StatusPendingPayment,
			CreatedAt: lo.ToPtr(time.Now()),
			UpdatedAt: nil,
		}

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(existingOrder, nil).Once()
		s.orderRepository.On("Update", s.ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(nil).Once()

		err := s.service.Cancel(s.ctx, orderUUID)

		require.NoError(s.T(), err)
	})

	s.Run("order_not_found", func() {
		orderUUID := gofakeit.UUID()

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

		err := s.service.Cancel(s.ctx, orderUUID)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
	})

	s.Run("order_already_paid", func() {
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

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(existingOrder, nil).Once()

		err := s.service.Cancel(s.ctx, orderUUID)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrOrderAlreadyPaid)
	})

	s.Run("order_already_cancelled", func() {
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()

		existingOrder := &model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.StatusCancelled,
			UpdatedAt: lo.ToPtr(time.Now()),
		}

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(existingOrder, nil).Once()

		err := s.service.Cancel(s.ctx, orderUUID)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrOrderCancelled)
	})

	s.Run("repository_get_error", func() {
		orderUUID := gofakeit.UUID()
		repoErr := errors.New("database connection failed")

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(&model.Order{}, repoErr).Once()

		err := s.service.Cancel(s.ctx, orderUUID)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, repoErr)
	})

	s.Run("repository_update_error", func() {
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		repoErr := errors.New("database update failed")

		existingOrder := &model.Order{
			OrderUUID: orderUUID,
			UserUUID:  userUUID,
			Status:    model.StatusPendingPayment,
		}

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(existingOrder, nil).Once()
		s.orderRepository.On("Update", s.ctx, orderUUID, mock.AnythingOfType("model.Order")).Return(repoErr).Once()

		err := s.service.Cancel(s.ctx, orderUUID)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, repoErr)
	})

	s.Run("empty_uuid", func() {
		emptyUUID := ""

		s.orderRepository.On("Get", s.ctx, emptyUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

		err := s.service.Cancel(s.ctx, emptyUUID)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
	})
}
