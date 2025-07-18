package converter

import (
	"time"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartsFilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(*filter.Categories))
	for _, category := range *filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 *filter.Uuids,
		Names:                 *filter.Names,
		Categories:            categories,
		ManufacturerCountries: *filter.ManufacturerCountries,
		Tags:                  *filter.Tags,
	}
}

func PartListToModel(r *inventoryV1.ListPartsResponse) model.PartsInfoFilter {
	if r == nil {
		return model.PartsInfoFilter{}
	}

	var parts []model.Part
	for _, part := range r.Parts {
		var metadata map[string]model.Value
		if part.Metadata != nil {
			metadata = make(map[string]model.Value, len(part.Metadata))
			for k, v := range part.Metadata {
				metadata[k] = protoValueToModel(v)
			}
		}

		var tags *[]string
		if len(part.Tags) > 0 {
			tags = &part.Tags
		}

		parts = append(parts, model.Part{
			Uuid:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category:      model.Category(part.Category),
			Dimensions: model.Dimensions{
				Length: part.Dimensions.Length,
				Width:  part.Dimensions.Width,
				Height: part.Dimensions.Height,
				Weight: part.Dimensions.Weight,
			},
			Manufacturer: model.Manufacturer{
				Name:    part.Manufacturer.Name,
				Country: part.Manufacturer.Country,
				Website: part.Manufacturer.Website,
			},
			Tags:      tags,
			Metadata:  metadata,
			CreatedAt: func() *time.Time { t := part.CreatedAt.AsTime(); return &t }(),
			UpdatedAt: func() *time.Time { t := part.UpdatedAt.AsTime(); return &t }(),
		})
	}

	return model.PartsInfoFilter{
		Parts: parts,
	}
}

func protoValueToModel(v *inventoryV1.Value) model.Value {
	if v == nil {
		return model.Value{}
	}

	switch kind := v.Kind.(type) {
	case *inventoryV1.Value_StringValue:
		return model.Value{Str: &kind.StringValue}
	case *inventoryV1.Value_Int64Value:
		return model.Value{Int: &kind.Int64Value}
	case *inventoryV1.Value_DoubleValue:
		return model.Value{Float: &kind.DoubleValue}
	case *inventoryV1.Value_BoolValue:
		return model.Value{Bool: &kind.BoolValue}
	default:
		return model.Value{}
	}
}
