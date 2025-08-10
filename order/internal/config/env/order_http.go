package env

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type orderHttpEnvConfig struct {
	Host            string        `env:"HTTP_HOST,required"`
	Port            string        `env:"HTTP_PORT,required"`
	ReadTimeout     time.Duration `env:"HTTP_READ_TIMEOUT,required"`
	ShutdownTimeout time.Duration `env:"HTTP_SHUT_DOWN_TIMEOUT,required"`
}

type orderHttpConfig struct {
	raw orderHttpEnvConfig
}

func NewOrderHTTPConfig() (*orderHttpConfig, error) {
	var raw orderHttpEnvConfig
	if err := env.Parse(&raw); err != nil {
		return nil, err
	}

	return &orderHttpConfig{raw: raw}, nil
}

func (cfg *orderHttpConfig) Address() string {
	return cfg.raw.Host + ":" + cfg.raw.Port
}

func (cfg *orderHttpConfig) Port() string {
	return cfg.raw.Port
}

func (cfg *orderHttpConfig) ReadTimeout() time.Duration {
	return cfg.raw.ReadTimeout
}

func (cfg *orderHttpConfig) ShutdownTimeout() time.Duration {
	return cfg.raw.ShutdownTimeout
}
