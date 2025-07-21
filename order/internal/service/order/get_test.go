package order

import (
	"errors"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func (s *ServiceSuite) TestGet() {
	s.Run("success", func() {
		orderUUID := gofakeit.UUID()
		userUUID := gofakeit.UUID()
		partUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}
		totalPrice := gofakeit.Float32()
		transactionUUID := gofakeit.UUID()
		paymentMethod := model.PaymentMethodCard
		createdAt := lo.ToPtr(time.Now())
		updatedAt := lo.ToPtr(time.Now())

		expectedOrder := &model.Order{
			OrderUUID:       orderUUID,
			UserUUID:        userUUID,
			PartUuids:       partUUIDs,
			TotalPrice:      totalPrice,
			TransactionUUID: &transactionUUID,
			PaymentMethod:   &paymentMethod,
			Status:          model.StatusPaid,
			CreatedAt:       createdAt,
			UpdatedAt:       updatedAt,
		}

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(expectedOrder, nil).Once()

		result, err := s.service.Get(s.ctx, orderUUID)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), expectedOrder, result)
		assert.Equal(s.T(), orderUUID, result.OrderUUID)
		assert.Equal(s.T(), userUUID, result.UserUUID)
		assert.Equal(s.T(), model.StatusPaid, result.Status)
		assert.NotNil(s.T(), result.PaymentMethod)
		assert.Equal(s.T(), paymentMethod, *result.PaymentMethod)
	})

	s.Run("order_not_found", func() {
		orderUUID := gofakeit.UUID()

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

		result, err := s.service.Get(s.ctx, orderUUID)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
		assert.Nil(s.T(), result)
	})

	s.Run("repository_error", func() {
		orderUUID := gofakeit.UUID()
		repoErr := errors.New("database connection failed")

		s.orderRepository.On("Get", s.ctx, orderUUID).Return(&model.Order{}, repoErr).Once()

		result, err := s.service.Get(s.ctx, orderUUID)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, repoErr)
		assert.Nil(s.T(), result)
	})

	s.Run("empty_uuid", func() {
		emptyUUID := ""

		s.orderRepository.On("Get", s.ctx, emptyUUID).Return(&model.Order{}, model.ErrOrderNotFound).Once()

		result, err := s.service.Get(s.ctx, emptyUUID)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrOrderNotFound)
		assert.Nil(s.T(), result)
	})
}
