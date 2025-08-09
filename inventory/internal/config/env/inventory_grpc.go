package env

import "github.com/caarlos0/env/v11"

type inventoryGRPCEnvConfig struct {
	Host string `env:"GRPC_HOST,required"`
	Port string `env:"GRPC_PORT,required"`
}

type inventoryGRPCConfig struct {
	raw inventoryGRPCEnvConfig
}

func NewInventoryGRPCConfig() (*inventoryGRPCConfig, error) {
	var raw inventoryGRPCEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &inventoryGRPCConfig{raw: raw}, nil
}

func (cfg *inventoryGRPCConfig) Address() string {
	return cfg.raw.Host + ":" + cfg.raw.Port
}
