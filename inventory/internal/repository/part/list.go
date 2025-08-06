package part

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	mongoFilter := createMongoFilter(filter)

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		return nil, err
	}
	defer func() {
		cerr := cursor.Close(ctx)
		if cerr != nil {
			log.Printf("failed to close cursor: %v\n", err)
		}
	}()

	var repoParts []repoModel.Part
	err = cursor.All(ctx, &repoParts)
	if err != nil {
		return nil, err
	}

	parts := make([]model.Part, len(repoParts))
	for i, repoPart := range repoParts {
		parts[i] = converter.PartToModel(repoPart)
	}

	return parts, nil
}

func createMongoFilter(filter model.PartsFilter) bson.M {
	mongoFilter := bson.M{}

	if filter.Uuids != nil && len(*filter.Uuids) > 0 {
		mongoFilter["_id"] = bson.M{"$in": *filter.Uuids}
	}

	if filter.Names != nil && len(*filter.Names) > 0 {
		mongoFilter["name"] = bson.M{"$in": *filter.Names}
	}

	if filter.Categories != nil && len(*filter.Categories) > 0 {
		mongoFilter["category"] = bson.M{"$in": *filter.Categories}
	}

	if filter.ManufacturerCountries != nil && len(*filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer.country"] = bson.M{"$in": *filter.ManufacturerCountries}
	}

	if filter.Tags != nil && len(*filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": *filter.Tags}
	}

	return mongoFilter
}
