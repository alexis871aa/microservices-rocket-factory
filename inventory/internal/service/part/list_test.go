package part

import (
	"context"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
	serviceMocks "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/service/mocks"
)

func Test_SuccessListPartsWithEmptyFilter(t *testing.T) {
	ctx := context.Background()
	partRepository := serviceMocks.NewPartRepository(t)
	service := NewService(partRepository)

	filter := model.PartsFilter{}
	expectedParts := []model.Part{
		{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Word(),
			Price:         gofakeit.Float64(),
			StockQuantity: gofakeit.Int64(),
			Category:      model.Category(gofakeit.IntRange(1, 4)),
			CreatedAt:     lo.ToPtr(time.Now()),
		},
		{
			Uuid:          gofakeit.UUID(),
			Name:          gofakeit.Name(),
			Description:   gofakeit.Word(),
			Price:         gofakeit.Float64(),
			StockQuantity: gofakeit.Int64(),
			Category:      model.Category(gofakeit.IntRange(1, 4)),
			CreatedAt:     lo.ToPtr(time.Now()),
		},
	}

	partRepository.On("ListParts", ctx, filter).Return(expectedParts, nil).Once()

	actualParts, err := service.ListParts(ctx, filter)

	require.NoError(t, err)
	assert.Equal(t, expectedParts, actualParts)
	assert.Len(t, actualParts, 2)
}

func Test_SuccessListPartsWithUUIDFilter(t *testing.T) {
	ctx := context.Background()
	partRepository := serviceMocks.NewPartRepository(t)
	service := NewService(partRepository)

	uuids := []string{gofakeit.UUID(), gofakeit.UUID()}
	filter := model.PartsFilter{
		Uuids: &uuids,
	}
	expectedParts := []model.Part{
		{
			Uuid:          uuids[0],
			Name:          gofakeit.Name(),
			Price:         gofakeit.Float64(),
			StockQuantity: gofakeit.Int64(),
			Category:      model.CategoryEngine,
			CreatedAt:     lo.ToPtr(time.Now()),
		},
	}

	partRepository.On("ListParts", ctx, filter).Return(expectedParts, nil).Once()

	actualParts, err := service.ListParts(ctx, filter)

	require.NoError(t, err)
	assert.Equal(t, expectedParts, actualParts)
	assert.Len(t, actualParts, 1)
	assert.Equal(t, uuids[0], actualParts[0].Uuid)
}

func Test_SuccessListPartsWithEmptyResult(t *testing.T) {
	ctx := context.Background()
	partRepository := serviceMocks.NewPartRepository(t)
	service := NewService(partRepository)

	filter := model.PartsFilter{
		Names: lo.ToPtr([]string{"nonexistent"}),
	}
	expectedParts := []model.Part{}

	partRepository.On("ListParts", ctx, filter).Return(expectedParts, nil).Once()

	actualParts, err := service.ListParts(ctx, filter)

	require.NoError(t, err)
	assert.Equal(t, expectedParts, actualParts)
	assert.Empty(t, actualParts)
}

func Test_ErrorWhenCantListPartsFromRepository(t *testing.T) {
	ctx := context.Background()
	partRepository := serviceMocks.NewPartRepository(t)
	service := NewService(partRepository)

	filter := model.PartsFilter{}
	repoErr := gofakeit.Error()

	partRepository.On("ListParts", ctx, filter).Return([]model.Part{}, repoErr).Once()

	actualParts, err := service.ListParts(ctx, filter)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Empty(t, actualParts)
}

func Test_SuccessListPartsWithCategoryFilter(t *testing.T) {
	ctx := context.Background()
	partRepository := serviceMocks.NewPartRepository(t)
	service := NewService(partRepository)

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
			CreatedAt:     lo.ToPtr(time.Now()),
		},
	}

	partRepository.On("ListParts", ctx, filter).Return(expectedParts, nil).Once()

	actualParts, err := service.ListParts(ctx, filter)

	require.NoError(t, err)
	assert.Equal(t, expectedParts, actualParts)
	assert.Len(t, actualParts, 1)
	assert.Equal(t, model.CategoryEngine, actualParts[0].Category)
}
