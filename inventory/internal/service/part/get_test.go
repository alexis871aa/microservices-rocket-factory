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

func Test_SuccessGetPart(t *testing.T) {
	ctx := context.Background()
	partRepository := serviceMocks.NewPartRepository(t)
	service := NewService(partRepository)

	uuid := gofakeit.UUID()
	expectedPart := model.Part{
		Uuid:          uuid,
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
			"color":  {Str: lo.ToPtr(gofakeit.Color())},
			"weight": {Float: lo.ToPtr(gofakeit.Float64())},
		},
		CreatedAt: lo.ToPtr(time.Now()),
		UpdatedAt: lo.ToPtr(time.Now()),
	}

	partRepository.On("GetPart", ctx, uuid).Return(expectedPart, nil).Once()

	actualPart, err := service.GetPart(ctx, uuid)

	require.NoError(t, err)
	assert.Equal(t, expectedPart, actualPart)
}

func Test_ErrorWhenCantGetPartFromRepository(t *testing.T) {
	ctx := context.Background()
	partRepository := serviceMocks.NewPartRepository(t)
	service := NewService(partRepository)

	uuid := gofakeit.UUID()
	repoErr := gofakeit.Error()

	partRepository.On("GetPart", ctx, uuid).Return(model.Part{}, repoErr).Once()

	actualPart, err := service.GetPart(ctx, uuid)

	require.Error(t, err)
	assert.ErrorIs(t, err, repoErr)
	assert.Empty(t, actualPart)
}

func Test_ErrorWhenPartNotFound(t *testing.T) {
	ctx := context.Background()
	partRepository := serviceMocks.NewPartRepository(t)
	service := NewService(partRepository)

	uuid := gofakeit.UUID()

	partRepository.On("GetPart", ctx, uuid).Return(model.Part{}, model.ErrPartNotFound).Once()

	actualPart, err := service.GetPart(ctx, uuid)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrPartNotFound)
	assert.Empty(t, actualPart)
}

func Test_ErrorWhenEmptyUUID(t *testing.T) {
	ctx := context.Background()
	partRepository := serviceMocks.NewPartRepository(t)
	service := NewService(partRepository)

	emptyUUID := ""

	partRepository.On("GetPart", ctx, emptyUUID).Return(model.Part{}, model.ErrPartNotFound).Once()

	actualPart, err := service.GetPart(ctx, emptyUUID)

	require.Error(t, err)
	assert.ErrorIs(t, err, model.ErrPartNotFound)
	assert.Empty(t, actualPart)
}
