package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/alexis871aa/microservices-rocket-factory/notification/internal/converter/kafka"
	def "github.com/alexis871aa/microservices-rocket-factory/notification/internal/service"
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
	logger.Info(ctx, "Starting order paid consumer service")

	err := s.orderConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Failed to start order paid consumer service", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Order paid consumer service successfully started")
	return nil
}
