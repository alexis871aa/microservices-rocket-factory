package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

const testsTimeout = 5 * time.Minute

var (
	env             *TestEnvironment
	inventoryClient inventoryV1.InventoryServiceClient
	suiteCtx        context.Context
	suiteCtxCancel  context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail) // хендлер, который будет выполняться, когда что-то зафейлится
	RunSpecs(t, "Inventory Service Integration Test Suite")
}

// setupInventoryClient - настраивает gRPC клиент для подключения к inventory сервису
func setupInventoryClient(ctx context.Context) {
	// получаем адрес приложения из контейнера
	inventoryAddr := env.App.Address()

	// Создаем gRPC подключение
	conn, err := grpc.NewClient(inventoryAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(ctx, "не удалось подключиться к inventory сервису", zap.Error(err))
	}

	inventoryClient = inventoryV1.NewInventoryServiceClient(conn)
	logger.Info(ctx, "🚀 gRPC клиент inventory сервиса создан", zap.String("address", inventoryAddr))
}

var _ = BeforeSuite(func() {
	// иницируем логгер
	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCtxCancel = context.WithTimeout(context.Background(), testsTimeout)

	// загружаем .env файл и устанавливаем переменные в окружении
	envVars, err := godotenv.Read(filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env"))
	if err != nil {
		logger.Fatal(suiteCtx, "не удалось загрузить .env файл", zap.Error(err))
	}

	// устанавливаем переменные в окружение процесса
	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	logger.Info(suiteCtx, "Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)

	// настраиваем gRPC клиент
	setupInventoryClient(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCtxCancel()
})
