package converter

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartToProto(part *model.Part) *inventoryV1.GetPartResponse {
	return &inventoryV1.GetPartResponse{
		Part: modelPartToProtoPart(part),
	}
}

func PartsToProto(parts []model.Part) *inventoryV1.ListPartsResponse {
	protoParts := make([]*inventoryV1.Part, 0, len(parts))
	for _, part := range parts {
		protoParts = append(protoParts, modelPartToProtoPart(&part))
	}

	return &inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}
}

func modelPartToProtoPart(part *model.Part) *inventoryV1.Part {
	var tags []string
	if part.Tags != nil {
		tags = *part.Tags
	}

	var meta map[string]*inventoryV1.Value
	if part.Metadata != nil {
		meta = lo.MapEntries(part.Metadata, func(k string, v model.Value) (string, *inventoryV1.Value) {
			return k, ModelValueToProto(&v)
		})
	}

	var createdAtProto *timestamppb.Timestamp
	if part.CreatedAt != nil {
		createdAtProto = timestamppb.New(*part.CreatedAt)
	}

	var updatedAtProto *timestamppb.Timestamp
	if part.UpdatedAt != nil {
		updatedAtProto = timestamppb.New(*part.UpdatedAt)
	}

	return &inventoryV1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      inventoryV1.Category(part.Category),
		Dimensions: &inventoryV1.Dimensions{
			Length: part.Dimensions.Length,
			Width:  part.Dimensions.Width,
			Height: part.Dimensions.Height,
			Weight: part.Dimensions.Weight,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    part.Manufacturer.Name,
			Country: part.Manufacturer.Country,
			Website: part.Manufacturer.Website,
		},
		Tags:      tags,
		Metadata:  meta,
		CreatedAt: createdAtProto,
		UpdatedAt: updatedAtProto,
	}
}

func ModelValueToProto(v *model.Value) *inventoryV1.Value {
	if v == nil {
		return &inventoryV1.Value{}
	}
	switch {
	case v.Str != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_StringValue{
				StringValue: *v.Str,
			},
		}
	case v.Int != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_Int64Value{
				Int64Value: *v.Int,
			},
		}
	case v.Float != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_DoubleValue{
				DoubleValue: *v.Float,
			},
		}
	case v.Bool != nil:
		return &inventoryV1.Value{
			Kind: &inventoryV1.Value_BoolValue{
				BoolValue: *v.Bool,
			},
		}
	default:
		return &inventoryV1.Value{}
	}
}

func ProtoToPartsFilter(p *inventoryV1.PartsFilter) *model.PartsFilter {
	if p == nil {
		return &model.PartsFilter{}
	}

	var uuids, names, manufacturerCountries, tags *[]string
	var categories *[]model.Category

	if len(p.Uuids) > 0 {
		uuids = &p.Uuids
	}
	if len(p.Names) > 0 {
		names = &p.Names
	}
	if len(p.ManufacturerCountries) > 0 {
		manufacturerCountries = &p.ManufacturerCountries
	}
	if len(p.Tags) > 0 {
		tags = &p.Tags
	}
	if len(p.Categories) > 0 {
		cat := make([]model.Category, len(p.Categories))
		for i, c := range p.Categories {
			cat[i] = model.Category(c)
		}
		categories = &cat
	}
	return &model.PartsFilter{
		Uuids:                 uuids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: manufacturerCountries,
		Tags:                  tags,
	}
}
