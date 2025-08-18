package integration

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
)

func (env *TestEnvironment) InsertTestPart() {}

// ClearInventoryCollection - удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearInventoryCollection(ctx context.Context) error {
	// используем бд из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
