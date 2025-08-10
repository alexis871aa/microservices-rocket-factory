package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/alexis871aa/microservices-rocket-factory/order/internal/api/order/v1"
	inventoryClient "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/config"
	customMiddleware "github.com/alexis871aa/microservices-rocket-factory/order/internal/middleware"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/migrator"
	orderRepository "github.com/alexis871aa/microservices-rocket-factory/order/internal/repository/order"
	orderService "github.com/alexis871aa/microservices-rocket-factory/order/internal/service/order"
	orderV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

const configPath = "./deploy/compose/order/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	inventoryConn, err := grpc.NewClient(
		config.AppConfig().InventoryClient.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("failed to close inventory connection: %v\n", err)
		}
	}()

	paymentConn, err := grpc.NewClient(
		config.AppConfig().PaymentClient.Address(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to connect: %v\n", err)
		return
	}
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("failed to close payment connection: %v\n", err)
		}
	}()

	ctx := context.Background()
	dbConn, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		err := dbConn.Close(ctx)
		if err != nil {
			log.Printf("failed to close database connection: %v\n", err)
		}
	}()

	err = dbConn.Ping(ctx)
	if err != nil {
		log.Printf("failed to ping database connection: %v\n", err)
		return
	}

	migratorRunner := migrator.NewMigrator(stdlib.OpenDB(*dbConn.Config().Copy()), config.AppConfig().Postgres.MigrationsDir())
	err = migratorRunner.Up()
	if err != nil {
		log.Printf("failed to run migrations: %v\n", err)
		return
	}

	repo := orderRepository.NewRepository(stdlib.OpenDB(*dbConn.Config().Copy()))
	service := orderService.NewService(
		repo,
		inventoryClient.NewClient(inventoryV1.NewInventoryServiceClient(inventoryConn)),
		paymentClient.NewClient(paymentV1.NewPaymentServiceClient(paymentConn)),
	)

	api := orderV1API.NewAPI(service)
	orderServer, err := orderV1.NewServer(api)
	if err != nil {
		log.Printf("ошибка создания сервера OpenAPI: %v", err)
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(2 * time.Second))
	r.Use(customMiddleware.RequestLogger)

	r.Mount("/", orderServer)

	server := &http.Server{
		Addr:              config.AppConfig().OrderHTTP.Address(),
		Handler:           r,
		ReadHeaderTimeout: config.AppConfig().OrderHTTP.ReadTimeout(),
	}

	go func() {
		log.Printf("🚀 HTTP-сервер запущен на порту %s\n", config.AppConfig().OrderHTTP.Port())
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("❌ Ошибка запуска http сервера: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Завершение работы http сервера...")

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig().OrderHTTP.ShutdownTimeout())
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("❌ Ошибка при остановке http сервера: %v\n", err)
	}

	log.Println("✅ Http сервер остановлен")
}
