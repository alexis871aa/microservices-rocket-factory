package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/alexis871aa/microservices-rocket-factory/order/internal/api/order/v1"
	client "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc"
	inventoryClientV1 "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/config"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/migrator"
	orderRepository "github.com/alexis871aa/microservices-rocket-factory/order/internal/repository/order"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/service"
	orderService "github.com/alexis871aa/microservices-rocket-factory/order/internal/service/order"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/closer"
	orderV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API    orderV1.Handler
	orderV1Server *orderV1.Server

	orderService service.OrderService

	orderRepository service.OrderRepository

	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient

	sqlDb   *sql.DB
	pgxConn *pgx.Conn

	inventoryGRPCConn *grpc.ClientConn
	paymentGRPCConn   *grpc.ClientConn

	migratorRunner *migrator.Migrator
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderV1API.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}

func (d *diContainer) OrderV1Server(ctx context.Context) *orderV1.Server {
	if d.orderV1Server == nil {
		server, err := orderV1.NewServer(d.OrderV1API(ctx))
		if err != nil {
			panic(fmt.Sprintf("failed to create OrderV1 server: %v", err))
		}
		d.orderV1Server = server
	}
	return d.orderV1Server
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(d.OrderRepository(ctx), d.InventoryClient(ctx), d.PaymentClient(ctx))
	}

	return d.orderService
}

func (d *diContainer) OrderRepository(ctx context.Context) service.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.SqlDB(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) InventoryClient(ctx context.Context) client.InventoryClient {
	if d.inventoryClient == nil {
		d.inventoryClient = inventoryClientV1.NewClient(inventoryV1.NewInventoryServiceClient(d.InventoryGRPCConn(ctx)))
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient(ctx context.Context) client.PaymentClient {
	if d.paymentClient == nil {
		d.paymentClient = paymentClientV1.NewClient(
			paymentV1.NewPaymentServiceClient(d.PaymentGRPCConn(ctx)),
		)
	}
	return d.paymentClient
}

func (d *diContainer) SqlDB(ctx context.Context) *sql.DB {
	if d.sqlDb == nil {
		d.sqlDb = stdlib.OpenDB(*d.PgxConn(ctx).Config().Copy())
	}

	return d.sqlDb
}

func (d *diContainer) PgxConn(ctx context.Context) *pgx.Conn {
	if d.pgxConn == nil {
		conn, err := pgx.Connect(ctx, config.AppConfig().Postgres.URI())
		if err != nil {
			panic(fmt.Sprintf("failed to connect to database: %v", err))
		}

		err = conn.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("failed to ping database: %v", err))
		}

		closer.AddNamed("PostgreSQL connection", func(ctx context.Context) error {
			return conn.Close(ctx)
		})

		d.pgxConn = conn
	}
	return d.pgxConn
}

func (d *diContainer) InventoryGRPCConn(_ context.Context) *grpc.ClientConn {
	if d.inventoryGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().InventoryClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to inventory service: %v", err))
		}

		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryGRPCConn = conn
	}
	return d.inventoryGRPCConn
}

func (d *diContainer) PaymentGRPCConn(_ context.Context) *grpc.ClientConn {
	if d.paymentGRPCConn == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().PaymentClient.Address(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to connect to payment service: %v", err))
		}

		closer.AddNamed("Payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentGRPCConn = conn
	}
	return d.paymentGRPCConn
}

func (d *diContainer) MigratorRunner(ctx context.Context) *migrator.Migrator {
	if d.migratorRunner == nil {
		d.migratorRunner = migrator.NewMigrator(
			d.SqlDB(ctx),
			config.AppConfig().Postgres.MigrationsDir(),
		)
	}
	return d.migratorRunner
}
