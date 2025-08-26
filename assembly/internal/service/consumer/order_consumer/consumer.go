package order_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/alexis871aa/microservices-rocket-factory/assembly/internal/converter/kafka"
	def "github.com/alexis871aa/microservices-rocket-factory/assembly/internal/service"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
)

var _ def.ConsumerService = (*service)(nil)

type service struct {
	orderConsumer    kafka.Consumer
	orderPaidDecoder kafkaConverter.OrderPaidDecoder
}

func NewService(orderConsumer kafka.Consumer, orderPaidDecoder kafkaConverter.OrderPaidDecoder) *service {
	return &service{
		orderConsumer:    orderConsumer,
		orderPaidDecoder: orderPaidDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting order consumer service")

	err := s.orderConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Failed to start order consumer service", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order consumer service successfully started")
	return nil
}
