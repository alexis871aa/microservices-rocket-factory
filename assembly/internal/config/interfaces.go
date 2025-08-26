package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderAssembledProducerConfig interface {
	Topic() string
	Config() *sarama.Config
}

type OrderPaidConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}

type KafkaConfig interface {
	Brokers() []string
}
