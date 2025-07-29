package part

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/converter"
)

func (r *repository) ListParts(_ context.Context, filter model.PartsFilter) ([]model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var parts []model.Part

	for _, repoPart := range r.data {
		modelPart := converter.PartToModel(repoPart)
		if matchesFilter(modelPart, filter) {
			parts = append(parts, modelPart)
		}
	}

	return parts, nil
}

func matchesFilter(part model.Part, filter model.PartsFilter) bool {
	return matchesBy(part, filter.Uuids, func(part model.Part) string { return part.Uuid }) &&
		matchesBy(part, filter.Names, func(part model.Part) string { return part.Name }) &&
		matchesBy(part, filter.Categories, func(part model.Part) model.Category { return part.Category }) &&
		matchesBy(part, filter.ManufacturerCountries, func(part model.Part) string { return part.Manufacturer.Country }) &&
		matchesByTags(part, filter.Tags)
}

func matchesByTags(part model.Part, filterTags *[]string) bool {
	if filterTags == nil || len(*filterTags) == 0 {
		return true
	}
	if part.Tags == nil {
		return false
	}

	partTags := *part.Tags
	for _, filterTag := range *filterTags {
		for _, partTag := range partTags {
			if partTag == filterTag {
				return true
			}
		}
	}
	return false
}

func matchesBy[T comparable](part model.Part, filter *[]T, get func(part model.Part) T) bool {
	if filter == nil || len(*filter) == 0 {
		return true
	}

	val := get(part)
	for _, filterVal := range *filter {
		if val == filterVal {
			return true
		}
	}
	return false
}
