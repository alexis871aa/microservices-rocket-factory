package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"

	"github.com/alexis871aa/microservices-rocket-factory/assembly/internal/config"
	kafkaConverter "github.com/alexis871aa/microservices-rocket-factory/assembly/internal/converter/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/assembly/internal/converter/kafka/decoder"
	"github.com/alexis871aa/microservices-rocket-factory/assembly/internal/service"
	orderConsumer "github.com/alexis871aa/microservices-rocket-factory/assembly/internal/service/consumer/order_consumer"
	orderProducer "github.com/alexis871aa/microservices-rocket-factory/assembly/internal/service/producer/order_producer"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka/producer"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderProducerService service.OrderProducerService
	orderConsumerService service.ConsumerService

	orderPaidConsumer wrappedKafka.Consumer
	consumerGroup     sarama.ConsumerGroup

	orderPaidDecoder      kafkaConverter.OrderPaidDecoder
	syncProducer          sarama.SyncProducer
	orderAssemblyProducer wrappedKafka.Producer
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderAssemblyProducer())
	}

	return d.orderProducerService
}

func (d *diContainer) OrderConsumerService() service.ConsumerService {
	if d.orderConsumerService == nil {
		d.orderConsumerService = orderConsumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder())
	}

	return d.orderConsumerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}

	return d.consumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderAssemblyProducer() wrappedKafka.Producer {
	if d.orderAssemblyProducer == nil {
		d.orderAssemblyProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderAssembledProducer.Topic(),
			logger.Logger(),
		)
	}

	return d.orderAssemblyProducer
}
