package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	partV1API "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/api/inventory/v1"
	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/config"
	partRepository "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/part"
	partService "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/service/part"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

const configPath = "./deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	lis, err := net.Listen("tcp", config.AppConfig().InventoryGRPC.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("failed to close listener: %v", err)
		}
	}()

	s := grpc.NewServer()

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
	if err != nil {
		log.Printf("failed to connect to database: %v\n", err)
		return
	}
	defer func() {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Printf("failed to disconnect: %v\n", err)
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("failed to ping database %v\n", err)
	}

	db := client.Database("inventory")

	repo := partRepository.NewRepository(db)
	err = repo.InitParts(ctx)
	if err != nil {
		log.Printf("failed to init parts: %v\n", err)
	}

	service := partService.NewService(repo)
	api := partV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("starting InventoryService server on port %s", config.AppConfig().InventoryGRPC.Port())
		err = s.Serve(lis)
		if err != nil {
			log.Fatalf("failed to serve: %v", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down InventoryService server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
