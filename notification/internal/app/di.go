package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"

	httpClient "github.com/alexis871aa/microservices-rocket-factory/notification/internal/client/http"
	telegramClient "github.com/alexis871aa/microservices-rocket-factory/notification/internal/client/http/telegram"
	"github.com/alexis871aa/microservices-rocket-factory/notification/internal/config"
	kafkaConverter "github.com/alexis871aa/microservices-rocket-factory/notification/internal/converter/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/notification/internal/converter/kafka/decoder"
	"github.com/alexis871aa/microservices-rocket-factory/notification/internal/service"
	orderAssembledConsumer "github.com/alexis871aa/microservices-rocket-factory/notification/internal/service/consumer/order_assembled_consumer"
	orderPaidConsumer "github.com/alexis871aa/microservices-rocket-factory/notification/internal/service/consumer/order_paid_consumer"
	telegramService "github.com/alexis871aa/microservices-rocket-factory/notification/internal/service/telegram"
	wrappedKafka "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka/consumer"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
	kafkaMiddleware "github.com/alexis871aa/microservices-rocket-factory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	orderPaidConsumerService      service.ConsumerService
	orderAssembledConsumerService service.ConsumerService
	telegramService               service.TelegramService

	telegramClient httpClient.TelegramClient
	telegramBot    *bot.Bot

	orderPaidConsumer      wrappedKafka.Consumer
	orderPaidConsumerGroup sarama.ConsumerGroup
	orderPaidDecoder       kafkaConverter.OrderPaidDecoder

	orderAssembledConsumer      wrappedKafka.Consumer
	orderAssembledConsumerGroup sarama.ConsumerGroup
	orderAssembledDecoder       kafkaConverter.ShipAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderPaidConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(d.OrderPaidConsumer(), d.OrderPaidDecoder(), d.TelegramService(ctx))
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) OrderAssembledConsumerService(ctx context.Context) service.ConsumerService {
	if d.orderAssembledConsumerService == nil {
		d.orderAssembledConsumerService = orderAssembledConsumer.NewService(d.OrderAssembledConsumer(), d.OrderAssembledDecoder(), d.TelegramService(ctx))
	}

	return d.orderAssembledConsumerService
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{
				config.AppConfig().OrderPaidConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderPaidConsumer
}

func (d *diContainer) OrderAssembledConsumer() wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderAssembledConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidConsumer.GroupID(),
			config.AppConfig().OrderPaidConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order paid consumer group: %s\n", err))
		}

		d.orderPaidConsumerGroup = consumerGroup
	}
	return d.orderPaidConsumerGroup
}

func (d *diContainer) OrderAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.orderAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order assembled consumer group: %s\n", err))
		}

		d.orderAssembledConsumerGroup = consumerGroup
	}
	return d.orderAssembledConsumerGroup
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderPaidDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderPaidDecoder()
	}

	return d.orderPaidDecoder
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.ShipAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewShipAssembledDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) TelegramService(ctx context.Context) service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegramService.NewService(d.TelegramClient(ctx))
	}

	return d.telegramService
}

func (d *diContainer) TelegramClient(ctx context.Context) httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot(ctx))
	}

	return d.telegramClient
}

func (d *diContainer) TelegramBot(_ context.Context) *bot.Bot {
	if d.telegramBot == nil {
		b, err := bot.New(config.AppConfig().TelegramBot.Token())
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s\n", err))
		}

		d.telegramBot = b
	}

	return d.telegramBot
}
