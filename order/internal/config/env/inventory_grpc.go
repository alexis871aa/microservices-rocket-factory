package env

import "github.com/caarlos0/env/v11"

type inventoryClientEnvConfig struct {
	Host string `env:"INVENTORY_GRPC_HOST,required"`
	Port string `env:"INVENTORY_GRPC_PORT,required"`
}

type inventoryClientConfig struct {
	raw inventoryClientEnvConfig
}

func NewInventoryClientConfig() (*inventoryClientConfig, error) {
	var raw inventoryClientEnvConfig
	err := env.Parse(&raw)
	if err != nil {
		return nil, err
	}

	return &inventoryClientConfig{raw: raw}, nil
}

func (cfg *inventoryClientConfig) Address() string {
	return cfg.raw.Host + ":" + cfg.raw.Port
}
