package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderV1API "github.com/alexis871aa/microservices-rocket-factory/order/internal/api/order/v1"
	client "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc"
	inventoryClientV1 "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc/inventory/v1"
	paymentClientV1 "github.com/alexis871aa/microservices-rocket-factory/order/internal/client/grpc/payment/v1"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/config"
	kafkaConverter "github.com/alexis871aa/microservices-rocket-factory/order/internal/converter/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/converter/kafka/decoder"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/migrator"
	orderRepository "github.com/alexis871aa/microservices-rocket-factory/order/internal/repository/order"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/service"
	orderConsumer "github.com/alexis871aa/microservices-rocket-factory/order/internal/service/consumer/order_consumer"
	orderService "github.com/alexis871aa/microservices-rocket-factory/order/internal/service/order"
	orderProducer "github.com/alexis871aa/microservices-rocket-factory/order/internal/service/producer/order_producer"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka/producer"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/middleware/kafka"
	orderV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	orderV1API    orderV1.Handler
	orderV1Server *orderV1.Server

	orderService         service.OrderService
	orderProducerService service.OrderProducerService
	orderConsumerService service.OrderConsumerService

	orderRepository service.OrderRepository

	inventoryClient client.InventoryClient
	paymentClient   client.PaymentClient

	sqlDb   *sql.DB
	pgxConn *pgx.Conn

	inventoryGRPCConn *grpc.ClientConn
	paymentGRPCConn   *grpc.ClientConn

	migratorRunner *migrator.Migrator

	consumerGroup sarama.ConsumerGroup
	orderConsumer wrappedKafka.Consumer

	orderAssemblyDecoder kafkaConverter.OrderAssembledDecoder
	syncProducer         sarama.SyncProducer
	orderProducer        wrappedKafka.Producer
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
			panic(fmt.Sprintf("💥 failed to create OrderV1 server: %v", err))
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

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderProducer())
	}

	return d.orderProducerService
}

func (d *diContainer) OrderConsumerService() service.OrderConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(d.OrderConsumer(), d.OrderAssemblyDecoder())
	}

	return d.orderConsumerService
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
			panic(fmt.Sprintf("💥 failed to connect to database: %v", err))
		}

		err = conn.Ping(ctx)
		if err != nil {
			panic(fmt.Sprintf("💥 failed to ping database: %v", err))
		}

		closer.AddNamed("PostgresSQL connection", func(ctx context.Context) error {
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
			panic(fmt.Sprintf("💥 failed to connect to inventory service: %v", err))
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
			panic(fmt.Sprintf("💥 failed to connect to payment service: %v", err))
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

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumerConfig.GroupID(),
			config.AppConfig().OrderAssembledConsumerConfig.Config(),
		)
		if err != nil {
			panic(fmt.Errorf("failed to create consumer group: %s\n", err))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderConsumer() wrappedKafka.Consumer {
	if d.orderConsumer == nil {
		d.orderConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumerConfig.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderConsumer
}

func (d *diContainer) OrderAssemblyDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssemblyDecoder == nil {
		d.orderAssemblyDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssemblyDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducerConfig.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderProducer() wrappedKafka.Producer {
	if d.orderProducer == nil {
		d.orderProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducerConfig.Topic(),
			logger.Logger(),
		)
	}

	return d.orderProducer
}
