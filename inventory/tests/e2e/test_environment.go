package integration

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

// TestPart - структура для тестовых данных деталей космических кораблей
type TestPart struct {
	UUID          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      inventoryV1.Category
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// Предустановленные тестовые детали ракет 🚀
var (
	RocketEngine = TestPart{
		UUID:          "engine-001",
		Name:          "Ракетный двигатель V8",
		Description:   "Мощный двигатель для космических ракет",
		Price:         15000.50,
		StockQuantity: 5,
		Category:      inventoryV1.Category_ENGINE,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	FuelTank = TestPart{
		UUID:          "fuel-001",
		Name:          "Топливный бак 2000л",
		Description:   "Большой топливный бак для ракетного топлива",
		Price:         8000.00,
		StockQuantity: 12,
		Category:      inventoryV1.Category_FUEL,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	Porthole = TestPart{
		UUID:          "porthole-001",
		Name:          "Иллюминатор космический",
		Description:   "Прочный иллюминатор для космических кораблей",
		Price:         2500.75,
		StockQuantity: 8,
		Category:      inventoryV1.Category_PORTHOLE,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
)

// InsertTestPart - вставляет тестовую деталь в MongoDB
func (env *TestEnvironment) InsertTestPart(ctx context.Context, part TestPart) error {
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	collection := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName)

	// конвертируем TestPart в MongoDB документ
	doc := bson.M{
		"_id":            part.UUID,
		"name":           part.Name,
		"description":    part.Description,
		"price":          part.Price,
		"stock_quantity": part.StockQuantity,
		"category":       int(part.Category),
		"created_at":     part.CreatedAt,
		"updated_at":     part.UpdatedAt,
	}

	_, err := collection.InsertOne(ctx, doc)
	return err
}

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

// PartExists - проверяет существование детали в БД по UUID
func (env *TestEnvironment) PartExists(ctx context.Context, uuid string) (bool, error) {
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	collection := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName)

	count, err := collection.CountDocuments(ctx, bson.M{"_id": uuid})
	return count > 0, err
}

// GetPartsCount - возвращает общее количество деталей в коллекции
func (env *TestEnvironment) GetPartsCount(ctx context.Context) (int64, error) {
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	collection := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName)

	return collection.CountDocuments(ctx, bson.M{})
}

// GetPartsByCategory - возвращает количество деталей определенной категории
func (env *TestEnvironment) GetPartsByCategory(ctx context.Context, category inventoryV1.Category) (int64, error) {
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory-service" // fallback значение
	}

	collection := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName)

	return collection.CountDocuments(ctx, bson.M{"category": int(category)})
}
