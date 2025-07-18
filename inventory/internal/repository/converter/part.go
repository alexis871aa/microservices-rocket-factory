package converter

import (
	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
	repoModel "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/model"
)

func PartToModel(part repoModel.Part) model.Part {
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

	return model.Part{
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
		UpdatedAt:     part.UpdatedAt,
	}
}
