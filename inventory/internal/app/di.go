package app

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	inventoryV1API "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/api/inventory/v1"
	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/config"
	inventoryRepository "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/repository/part"
	"github.com/alexis871aa/microservices-rocket-factory/inventory/internal/service"
	inventoryService "github.com/alexis871aa/microservices-rocket-factory/inventory/internal/service/part"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/closer"
	authV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/auth/v1"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API      inventoryV1.InventoryServiceServer
	inventoryService    service.InventoryService
	inventoryRepository service.InventoryRepository
	iamClient           authV1.AuthServiceClient

	mongoDBClient *mongo.Client
	mongoDBHandle *mongo.Database
	iamGRPCConn   *grpc.ClientConn
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.InventoryService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.InventoryService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewService(d.InventoryRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) service.InventoryRepository {
	if d.inventoryRepository == nil {
		repoCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		repo := inventoryRepository.NewRepository(repoCtx, d.MongoDBHandle(ctx))

		err := repo.InitParts(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to init parts: %v\n", err))
		}

		d.inventoryRepository = repo
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBHandle(ctx context.Context) *mongo.Database {
	if d.mongoDBHandle == nil {
		d.mongoDBHandle = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DatabaseName())
	}

	return d.mongoDBHandle
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.AppConfig().Mongo.URI()))
		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to MongoDB: %s\n", err.Error()))
		}

		err = client.Ping(ctx, readpref.Primary())
		if err != nil {
			panic(fmt.Sprintf("💥 failed to ping MongoDB: %v\n", err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) IAMClient(ctx context.Context) authV1.AuthServiceClient {
	if d.iamClient == nil {
		d.iamClient = authV1.NewAuthServiceClient(d.IAMGRPCConn(ctx))
	}
	return d.iamClient
}

func (d *diContainer) IAMGRPCConn(_ context.Context) *grpc.ClientConn {
	if d.iamGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().IAMClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("💥 failed to connect to IAM service: %v", err))
		}

		closer.AddNamed("IAM gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.iamGRPCConn = conn
	}
	return d.iamGRPCConn
}
