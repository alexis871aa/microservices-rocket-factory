package converter

import (
	"github.com/samber/lo"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

func PartInfoToProto(info *model.PartInfo) *inventoryV1.GetPartResponse {
	var tags []string
	if info.Part.Tags != nil {
		tags = *info.Part.Tags
	}

	var meta map[string]*inventoryV1.Value
	if info.Part.Metadata != nil {
		meta = lo.MapEntries(info.Part.Metadata, func(k string, v model.Value) (string, *inventoryV1.Value) {
			return k, ModelValueToProto(&v)
		})
	}

	var createdAtProto *timestamppb.Timestamp
	if info.Part.CreatedAt != nil {
		createdAtProto = timestamppb.New(*info.Part.CreatedAt)
	}

	var updatedAtProto *timestamppb.Timestamp
	if info.Part.UpdatedAt != nil {
		updatedAtProto = timestamppb.New(*info.Part.UpdatedAt)
	}

	return &inventoryV1.GetPartResponse{
		Part: &inventoryV1.Part{
			Uuid:          info.Part.Uuid,
			Name:          info.Part.Name,
			Description:   info.Part.Description,
			Price:         info.Part.Price,
			StockQuantity: info.Part.StockQuantity,
			Category:      inventoryV1.Category(info.Part.Category),
			Dimensions: &inventoryV1.Dimensions{
				Length: info.Part.Dimensions.Length,
				Width:  info.Part.Dimensions.Width,
				Height: info.Part.Dimensions.Height,
				Weight: info.Part.Dimensions.Weight,
			},
			Manufacturer: &inventoryV1.Manufacturer{
				Name:    info.Part.Manufacturer.Name,
				Country: info.Part.Manufacturer.Country,
				Website: info.Part.Manufacturer.Website,
			},
			Tags:      tags,
			Metadata:  meta,
			CreatedAt: createdAtProto,
			UpdatedAt: updatedAtProto,
		},
	}
}

func PartsInfoFilterToProto(info *model.PartsInfoFilter) *inventoryV1.ListPartsResponse {
	if info == nil {
		return &inventoryV1.ListPartsResponse{}
	}

	var parts []*inventoryV1.Part
	for _, part := range parts {
		parts = append(parts, &inventoryV1.Part{
			Uuid:          part.Uuid,
			Name:          part.Name,
			Description:   part.Description,
			Price:         part.Price,
			StockQuantity: part.StockQuantity,
			Category:      part.Category,
			Dimensions:    part.Dimensions,
			Metadata:      part.Metadata,
			CreatedAt:     part.CreatedAt,
			UpdatedAt:     part.UpdatedAt,
		})
	}

	return &inventoryV1.ListPartsResponse{
		Parts: parts,
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

	if len(p.Uuids) < 0 {
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
