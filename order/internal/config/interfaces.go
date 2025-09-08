package config

import (
	"time"

	"github.com/IBM/sarama"
)

type InventoryGRPCClientConfig interface {
	Address() string
}

type PaymentGRPCClientConfig interface {
	Address() string
}

type OrderHTTPConfig interface {
	Address() string
	Port() string
	ReadTimeout() time.Duration
	ShutdownTimeout() time.Duration
}

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
	DatabaseName() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderPaidProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderAssembledConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
