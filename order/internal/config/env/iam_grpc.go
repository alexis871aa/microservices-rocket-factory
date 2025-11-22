package env

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type IAMGRPCConfig struct {
	Host string `env:"IAM_GRPC_HOST" envDefault:"localhost"`
	Port string `env:"IAM_GRPC_PORT" envDefault:"50053"`
}

func NewIAMGRPCConfig() (*IAMGRPCConfig, error) {
	cfg := &IAMGRPCConfig{}

	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse IAM gRPC config: %w", err)
	}

	return cfg, nil
}

func (c *IAMGRPCConfig) Address() string {
	return fmt.Sprintf("%s:%s", c.Host, c.Port)
}
