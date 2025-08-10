package config

import "time"

type InventoryClientConfig interface {
	Address() string
}

type PaymentClientConfig interface {
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
