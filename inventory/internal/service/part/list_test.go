package part

import (
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
)

func (s *ServiceSuite) TestListParts() {
	s.Run("success_empty_filter", func() {
		filter := model.PartsFilter{}
		expectedParts := []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Word(),
				Price:         gofakeit.Float64(),
				StockQuantity: gofakeit.Int64(),
				Category:      model.Category(gofakeit.IntRange(1, 4)),
				Dimensions: model.Dimensions{
					Length: gofakeit.Float64(),
					Width:  gofakeit.Float64(),
					Height: gofakeit.Float64(),
					Weight: gofakeit.Float64(),
				},
				Manufacturer: model.Manufacturer{
					Name:    gofakeit.Name(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags: lo.ToPtr([]string{gofakeit.Word(), gofakeit.Word()}),
				Metadata: map[string]model.Value{
					"color": {Str: lo.ToPtr(gofakeit.Color())},
					"count": {Int: lo.ToPtr(gofakeit.Int64())},
				},
				CreatedAt: lo.ToPtr(time.Now()),
				UpdatedAt: lo.ToPtr(time.Now()),
			},
			{
				Uuid:          gofakeit.UUID(),
				Name:          gofakeit.Name(),
				Description:   gofakeit.Word(),
				Price:         gofakeit.Float64(),
				StockQuantity: gofakeit.Int64(),
				Category:      model.Category(gofakeit.IntRange(1, 4)),
				Dimensions: model.Dimensions{
					Length: gofakeit.Float64(),
					Width:  gofakeit.Float64(),
					Height: gofakeit.Float64(),
					Weight: gofakeit.Float64(),
				},
				Manufacturer: model.Manufacturer{
					Name:    gofakeit.Name(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags: lo.ToPtr([]string{gofakeit.Word()}),
				Metadata: map[string]model.Value{
					"material": {Str: lo.ToPtr(gofakeit.Word())},
					"premium":  {Bool: lo.ToPtr(gofakeit.Bool())},
				},
				CreatedAt: lo.ToPtr(time.Now()),
				UpdatedAt: nil,
			},
		}

		s.partRepository.On("ListParts", s.ctx, filter).Return(expectedParts, nil).Once()

		actualParts, err := s.service.ListParts(s.ctx, filter)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), expectedParts, actualParts)
		assert.Len(s.T(), actualParts, 2)
	})

	s.Run("success_with_filter", func() {
		uuids := []string{gofakeit.UUID(), gofakeit.UUID()}
		filter := model.PartsFilter{
			Uuids: &uuids,
		}
		expectedParts := []model.Part{
			{
				Uuid:          uuids[0],
				Name:          gofakeit.Name(),
				Description:   gofakeit.Word(),
				Price:         gofakeit.Float64(),
				StockQuantity: gofakeit.Int64(),
				Category:      model.CategoryEngine,
				Dimensions: model.Dimensions{
					Length: gofakeit.Float64(),
					Width:  gofakeit.Float64(),
					Height: gofakeit.Float64(),
					Weight: gofakeit.Float64(),
				},
				Manufacturer: model.Manufacturer{
					Name:    gofakeit.Name(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				Tags:      lo.ToPtr([]string{"engine", "main"}),
				Metadata:  map[string]model.Value{"type": {Str: lo.ToPtr("rocket")}},
				CreatedAt: lo.ToPtr(time.Now()),
			},
		}

		s.partRepository.On("ListParts", s.ctx, filter).Return(expectedParts, nil).Once()

		actualParts, err := s.service.ListParts(s.ctx, filter)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), expectedParts, actualParts)
		assert.Len(s.T(), actualParts, 1)
		assert.Equal(s.T(), uuids[0], actualParts[0].Uuid)
	})

	s.Run("success_empty_result", func() {
		filter := model.PartsFilter{
			Names: lo.ToPtr([]string{"nonexistent"}),
		}
		expectedParts := []model.Part{}

		s.partRepository.On("ListParts", s.ctx, filter).Return(expectedParts, nil).Once()

		actualParts, err := s.service.ListParts(s.ctx, filter)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), expectedParts, actualParts)
		assert.Empty(s.T(), actualParts)
	})

	s.Run("repository_error", func() {
		filter := model.PartsFilter{}
		repoErr := gofakeit.Error()

		s.partRepository.On("ListParts", s.ctx, filter).Return([]model.Part{}, repoErr).Once()

		actualParts, err := s.service.ListParts(s.ctx, filter)

		require.Error(s.T(), err)
		assert.ErrorIs(s.T(), err, repoErr)
		assert.Empty(s.T(), actualParts)
	})

	s.Run("success_category_filter", func() {
		categories := []model.Category{model.CategoryEngine, model.CategoryFuel}
		filter := model.PartsFilter{
			Categories: &categories,
		}
		expectedParts := []model.Part{
			{
				Uuid:          gofakeit.UUID(),
				Name:          "Engine Part",
				Category:      model.CategoryEngine,
				Price:         gofakeit.Float64(),
				StockQuantity: gofakeit.Int64(),
				Dimensions: model.Dimensions{
					Length: gofakeit.Float64(),
					Width:  gofakeit.Float64(),
					Height: gofakeit.Float64(),
					Weight: gofakeit.Float64(),
				},
				Manufacturer: model.Manufacturer{
					Name:    gofakeit.Name(),
					Country: gofakeit.Country(),
					Website: gofakeit.URL(),
				},
				CreatedAt: lo.ToPtr(time.Now()),
			},
		}

		s.partRepository.On("ListParts", s.ctx, filter).Return(expectedParts, nil).Once()

		actualParts, err := s.service.ListParts(s.ctx, filter)

		require.NoError(s.T(), err)
		assert.Equal(s.T(), expectedParts, actualParts)
		assert.Len(s.T(), actualParts, 1)
		assert.Equal(s.T(), model.CategoryEngine, actualParts[0].Category)
	})
}
