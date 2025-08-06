package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	partV1API "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/api/inventory/v1"
	partRepository "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/part"
	partService "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/service/part"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v", cerr)
		}
	}()

	s := grpc.NewServer()

	ctx := context.Background()
	lerr := godotenv.Load(".env")
	if lerr != nil {
		log.Printf("failed to load .env file: %v\n", lerr)
		return
	}

	dbUri := os.Getenv("MONGO_URI")

	client, connerr := mongo.Connect(ctx, options.Client().ApplyURI(dbUri))
	if connerr != nil {
		log.Printf("failed to connect to database: %v\n", connerr)
		return
	}
	defer func() {
		derr := client.Disconnect(ctx)
		if derr != nil {
			log.Printf("failed to disconnect: %v\n", derr)
		}
	}()

	perr := client.Ping(ctx, nil)
	if perr != nil {
		log.Printf("failed to ping database %v\n", perr)
	}

	db := client.Database("inventory")

	repo := partRepository.NewRepository(db)
	ierr := repo.InitParts(ctx)
	if ierr != nil {
		log.Printf("failed to init parts: %v\n", ierr)
	}

	service := partService.NewService(repo)
	api := partV1API.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("starting InventoryService server on port %d", grpcPort)
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
