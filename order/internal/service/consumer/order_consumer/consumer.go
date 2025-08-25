package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/alexis871aa/microservices-rocket-factory/order/internal/converter/kafka"
	def "github.com/alexis871aa/microservices-rocket-factory/order/internal/service"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
)

var _ def.OrderConsumerService = (*service)(nil)

type service struct {
	orderConsumer             kafka.Consumer
	orderShipAssembledDecoder kafkaConverter.OrderAssembledDecoder
	orderRepository           def.OrderRepository
}

func NewService(orderConsumer kafka.Consumer, orderShipAssembledDecoder kafkaConverter.OrderAssembledDecoder, orderRepository def.OrderRepository) *service {
	return &service{
		orderConsumer:             orderConsumer,
		orderShipAssembledDecoder: orderShipAssembledDecoder,
		orderRepository:           orderRepository,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer service")

	err := s.orderConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Consume from order.assembled topic error", zap.Error(err))
		return err
	}

	return nil
}
