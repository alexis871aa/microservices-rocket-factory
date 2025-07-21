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

func (s *ServiceSuite) TestCreate() {
	s.Run("success", func() {
		userUUID := gofakeit.UUID()
		partUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}

		parts := []model.Part{
			{
				Uuid:  partUUIDs[0],
				Name:  "Engine",
				Price: 100.0,
				Dimensions: model.Dimensions{
					Length: 10.0,
					Width:  5.0,
					Height: 3.0,
					Weight: 50.0,
				},
				Manufacturer: model.Manufacturer{
					Name:    "SpaceX",
					Country: "USA",
					Website: "spacex.com",
				},
				CreatedAt: lo.ToPtr(time.Now()),
			},
			{
				Uuid:  partUUIDs[1],
				Name:  "Fuel Tank",
				Price: 200.0,
				Dimensions: model.Dimensions{
					Length: 20.0,
					Width:  10.0,
					Height: 15.0,
					Weight: 100.0,
				},
				Manufacturer: model.Manufacturer{
					Name:    "Boeing",
					Country: "USA",
					Website: "boeing.com",
				},
				CreatedAt: lo.ToPtr(time.Now()),
			},
		}

		filter := model.PartsFilter{
			Uuids: &partUUIDs,
		}

		s.inventoryClient.On("ListParts", mock.Anything, filter).Return(parts, nil).Once()
		s.orderRepository.On("Create", mock.Anything, mock.AnythingOfType("model.Order")).Return(nil).Once()

		result, err := s.service.Create(s.ctx, userUUID, partUUIDs)

		require.NoError(s.T(), err)
		assert.NotNil(s.T(), result)
		assert.Equal(s.T(), userUUID, result.UserUUID)
		assert.Equal(s.T(), partUUIDs, result.PartUuids)
		assert.Equal(s.T(), float32(300.0), result.TotalPrice)
		assert.Equal(s.T(), model.StatusPendingPayment, result.Status)
		assert.NotEmpty(s.T(), result.OrderUUID)
		assert.NotNil(s.T(), result.CreatedAt)
		assert.Nil(s.T(), result.UpdatedAt)
		assert.Nil(s.T(), result.TransactionUUID)
		assert.Nil(s.T(), result.PaymentMethod)
	})

	s.Run("parts_not_found", func() {
		userUUID := gofakeit.UUID()
		partUUIDs := []string{gofakeit.UUID(), gofakeit.UUID()}

		filter := model.PartsFilter{
			Uuids: &partUUIDs,
		}

		s.inventoryClient.On("ListParts", mock.Anything, filter).Return([]model.Part{}, nil).Once()

		result, err := s.service.Create(s.ctx, userUUID, partUUIDs)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrPartsNotFound)
		assert.Nil(s.T(), result)
	})

	s.Run("partial_parts_found", func() {
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

		s.inventoryClient.On("ListParts", mock.Anything, filter).Return(parts, nil).Once()

		result, err := s.service.Create(s.ctx, userUUID, partUUIDs)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrPartsNotFound)
		assert.Nil(s.T(), result)
	})

	s.Run("inventory_client_error", func() {
		userUUID := gofakeit.UUID()
		partUUIDs := []string{gofakeit.UUID()}
		clientErr := errors.New("inventory service unavailable")

		filter := model.PartsFilter{
			Uuids: &partUUIDs,
		}

		s.inventoryClient.On("ListParts", mock.Anything, filter).Return([]model.Part{}, clientErr).Once()

		result, err := s.service.Create(s.ctx, userUUID, partUUIDs)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, clientErr)
		assert.Nil(s.T(), result)
	})

	s.Run("repository_error", func() {
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

		s.inventoryClient.On("ListParts", mock.Anything, filter).Return(parts, nil).Once()
		s.orderRepository.On("Create", mock.Anything, mock.AnythingOfType("model.Order")).Return(repoErr).Once()

		result, err := s.service.Create(s.ctx, userUUID, partUUIDs)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, repoErr)
		assert.Nil(s.T(), result)
	})

	s.Run("empty_part_uuids", func() {
		userUUID := gofakeit.UUID()
		partUUIDs := []string{}

		// Не устанавливаем моки, так как сервис должен вернуть ошибку до вызова клиента

		result, err := s.service.Create(s.ctx, userUUID, partUUIDs)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrPartsNotFound)
		assert.Nil(s.T(), result)
	})
}
