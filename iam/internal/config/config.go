package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/alexis871aa/microservices-rocket-factory/iam/internal/config/env"
)

var appConfig *config

type config struct {
	IAMGRPC  IAMGRPCConfig
	Logger   LoggerConfig
	Postgres PostgresConfig
	Redis    RedisConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	iamGRPCCfg, err := env.NewIAMGRPCConfig()
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

	redisCfg, err := env.NewRedisConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		IAMGRPC:  iamGRPCCfg,
		Logger:   loggerCfg,
		Postgres: postgresCfg,
		Redis:    redisCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
