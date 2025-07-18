package converter

import (
	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
	repoModel "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/model"
)

func PartInfoToRepoModel(info model.PartInfo) repoModel.PartInfo {
	var meta map[string]repoModel.Value
	if info.Part.Metadata != nil {
		meta = make(map[string]repoModel.Value, len(info.Part.Metadata))
		for k, v := range info.Part.Metadata {
			meta[k] = repoModel.Value{
				Str:   v.Str,
				Int:   v.Int,
				Float: v.Float,
				Bool:  v.Bool,
			}
		}
	}

	return repoModel.PartInfo{
		Part: repoModel.Part{
			Uuid:          info.Part.Uuid,
			Name:          info.Part.Name,
			Description:   info.Part.Description,
			Price:         info.Part.Price,
			StockQuantity: info.Part.StockQuantity,
			Category:      repoModel.Category(info.Part.Category),
			Dimensions:    repoModel.Dimensions(info.Part.Dimensions),
			Manufacturer:  repoModel.Manufacturer(info.Part.Manufacturer),
			Tags:          info.Part.Tags,
			Metadata:      meta,
			CreatedAt:     info.Part.CreatedAt,
			UpdatedAt:     info.Part.UpdatedAt,
		},
	}
}

func PartInfoToModel(info repoModel.PartInfo) model.PartInfo {
	var meta map[string]model.Value
	if info.Part.Metadata != nil {
		meta = make(map[string]model.Value, len(info.Part.Metadata))
		for k, v := range info.Part.Metadata {
			meta[k] = model.Value{
				Str:   v.Str,
				Int:   v.Int,
				Float: v.Float,
				Bool:  v.Bool,
			}
		}
	}

	return model.PartInfo{
		Part: model.Part{
			Uuid:          info.Part.Uuid,
			Name:          info.Part.Name,
			Description:   info.Part.Description,
			Price:         info.Part.Price,
			StockQuantity: info.Part.StockQuantity,
			Category:      model.Category(info.Part.Category),
			Dimensions:    model.Dimensions(info.Part.Dimensions),
			Manufacturer:  model.Manufacturer(info.Part.Manufacturer),
			Tags:          info.Part.Tags,
			Metadata:      meta,
			CreatedAt:     info.Part.CreatedAt,
			UpdatedAt:     info.Part.UpdatedAt,
		},
	}
}

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
