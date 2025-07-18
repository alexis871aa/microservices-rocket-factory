package part

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) ListParts(_ context.Context, filter model.PartsFilter) (model.PartsInfoFilter, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var parts []model.Part

	for _, part := range r.data {
		var meta map[string]model.Value
		if part.Metadata != nil {
			meta = make(map[string]model.Value, len(part.Metadata))
			for k, v := range part.Metadata {
				meta[k] = model.Value{
					Str:   v.Str,
					Int:   v.Int,
					Float: v.Float,
					Bool:  v.Bool,
				}
			}
		}

		if matchesFilter(converter.PartToModel(part), filter) {
			parts = append(parts, model.Part{
				Uuid:          part.Uuid,
				Name:          part.Name,
				Description:   part.Description,
				Price:         part.Price,
				StockQuantity: part.StockQuantity,
				Category:      model.Category(part.Category),
				Dimensions:    model.Dimensions(part.Dimensions),
				Manufacturer:  model.Manufacturer(part.Manufacturer),
				Tags:          part.Tags,
				Metadata:      meta,
				CreatedAt:     part.CreatedAt,
			})
		}
	}

	return model.PartsInfoFilter{
		Parts: parts,
	}, nil
}

func matchesFilter(part model.Part, filter model.PartsFilter) bool {
	return matchesBy(part, filter.Uuids, func(part model.Part) string { return part.Uuid }) &&
		matchesBy(part, filter.Names, func(part model.Part) string { return part.Name }) &&
		matchesBy(part, filter.Categories, func(part model.Part) model.Category { return part.Category }) &&
		matchesBy(part, filter.ManufacturerCountries, func(part model.Part) string { return part.Manufacturer.Country })
}

func matchesBy[T comparable](part model.Part, filter *[]T, get func(part model.Part) T) bool {
	if filter == nil || len(*filter) == 0 {
		return true
	}
	for _, filterValue := range *filter {
		if get(part) == filterValue {
			return true
		}
	}
	return false
}
