package env

import "github.com/caarlos0/env/v11"

type paymentClientEnvConfig struct {
	Host string `env:"PAYMENT_GRPC_HOST,required"`
	Port string `env:"PAYMENT_GRPC_PORT,required"`
}

type paymentClientConfig struct {
	raw paymentClientEnvConfig
}

func NewPaymentClientConfig() (*paymentClientConfig, error) {
	var raw paymentClientEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &paymentClientConfig{raw: raw}, nil
}

func (cfg *paymentClientConfig) Address() string {
	return cfg.raw.Host + ":" + cfg.raw.Port
}
