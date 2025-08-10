package config

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type InventoryGRPCConfig interface {
	Address() string
	Port() string
}

type MongoConfig interface {
	URI() string
	DatabaseName() string
}
