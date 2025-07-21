package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestGetPart() {
	s.Run("success", func() {
		var (
			uuid          = gofakeit.UUID()
			name          = gofakeit.Name()
			description   = gofakeit.Word()
			price         = gofakeit.Float64()
			stockQuantity = gofakeit.Int64()
			category      = model.Category(gofakeit.IntRange(1, 4))
			dimensions    = model.Dimensions{
				Length: gofakeit.Float64(),
				Width:  gofakeit.Float64(),
				Height: gofakeit.Float64(),
				Weight: gofakeit.Float64(),
			}
			manufacturer = model.Manufacturer{
				Name:    gofakeit.Name(),
				Country: gofakeit.Country(),
				Website: gofakeit.URL(),
			}
			tags = lo.ToPtr([]string{
				gofakeit.Word(),
				gofakeit.Word(),
				gofakeit.Word(),
			})
			metadata = map[string]model.Value{
				"color":         {Str: lo.ToPtr(gofakeit.Color())},
				"weight":        {Float: lo.ToPtr(gofakeit.Float64())},
				"is_premium":    {Bool: lo.ToPtr(gofakeit.Bool())},
				"serial_number": {Int: lo.ToPtr(gofakeit.Int64())},
			}
			createdAt = lo.ToPtr(time.Now())
			updatedAt = lo.ToPtr(time.Now())
		)

		expectedPart := model.Part{
			Uuid:          uuid,
			Name:          name,
			Description:   description,
			Price:         price,
			StockQuantity: stockQuantity,
			Category:      category,
			Dimensions:    dimensions,
			Manufacturer:  manufacturer,
			Tags:          tags,
			Metadata:      metadata,
			CreatedAt:     createdAt,
			UpdatedAt:     updatedAt,
		}

		s.partRepository.On("GetPart", s.ctx, uuid).Return(expectedPart, nil).Once()

		actualPart, err := s.service.GetPart(s.ctx, uuid)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), expectedPart, actualPart)
		assert.Equal(s.T(), uuid, actualPart.Uuid)
		assert.Equal(s.T(), name, actualPart.Name)
		assert.NotNil(s.T(), actualPart.Tags)
		assert.Len(s.T(), actualPart.Metadata, 4)
	})

	s.Run("repository_error", func() {
		var (
			uuid    = gofakeit.UUID()
			repoErr = gofakeit.Error()
		)

		s.partRepository.On("GetPart", s.ctx, uuid).Return(model.Part{}, repoErr).Once()

		actualPart, err := s.service.GetPart(s.ctx, uuid)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, repoErr)
		assert.Empty(s.T(), actualPart)
	})

	s.Run("part_not_found", func() {
		uuid := gofakeit.UUID()

		s.partRepository.On("GetPart", s.ctx, uuid).Return(model.Part{}, model.ErrPartNotFound).Once()

		actualPart, err := s.service.GetPart(s.ctx, uuid)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrPartNotFound)
		assert.Empty(s.T(), actualPart)
	})

	s.Run("empty_uuid", func() {
		emptyUUID := ""

		s.partRepository.On("GetPart", s.ctx, emptyUUID).Return(model.Part{}, model.ErrPartNotFound).Once()

		actualPart, err := s.service.GetPart(s.ctx, emptyUUID)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, model.ErrPartNotFound)
		assert.Empty(s.T(), actualPart)
	})
}
