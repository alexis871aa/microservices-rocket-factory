package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/model"
	repoConverter "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/converter"
	repoModel "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/model"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (model.Part, error) {
	var repoPart repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"_id": uuid}).Decode(&repoPart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}

		return model.Part{}, err
	}

	return repoConverter.PartToModel(repoPart), nil
}
