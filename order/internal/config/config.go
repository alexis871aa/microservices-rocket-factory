package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/config/env"
)

var appConfig *config

type config struct {
	InventoryClient InventoryGRPCClientConfig
	PaymentClient   PaymentGRPCClientConfig
	OrderHTTP       OrderHTTPConfig
	Logger          LoggerConfig
	Postgres        PostgresConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	inventoryClientCfg, err := env.NewInventoryClientConfig()
	if err != nil {
		return err
	}

	paymentClientCfg, err := env.NewPaymentClientConfig()
	if err != nil {
		return err
	}

	orderHTTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		InventoryClient: inventoryClientCfg,
		PaymentClient:   paymentClientCfg,
		OrderHTTP:       orderHTTPCfg,
		Logger:          loggerCfg,
		Postgres:        postgresCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
