package env

import "github.com/caarlos0/env/v11"

type iamGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type iamGRPCConfig struct {
	raw iamGRPCEnvConfig
}

func NewIAMGRPCConfig() (*iamGRPCConfig, error) {
	var raw iamGRPCEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &iamGRPCConfig{raw: raw}, nil
}

func (cfg *iamGRPCConfig) Address() string {
	return cfg.raw.Host + ":" + cfg.raw.Port
}

func (cfg *iamGRPCConfig) Port() string {
	return cfg.raw.Port
}
