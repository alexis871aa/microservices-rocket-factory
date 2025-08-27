package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	kafkaConverter "github.com/alexis871aa/microservices-rocket-factory/notification/internal/converter/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
)

type service struct {
	orderConsumer        kafka.Consumer
	shipAssembledDecoder kafkaConverter.ShipAssembledDecoder
}

func NewService(orderConsumer kafka.Consumer, shipAssembledDecoder kafkaConverter.ShipAssembledDecoder) *service {
	return &service{
		orderConsumer:        orderConsumer,
		shipAssembledDecoder: shipAssembledDecoder,
	}
}

func (s *service) RunConsumer(ctx context.Context) error {
	logger.Info(ctx, "Starting ship assembled consumer service")

	err := s.orderConsumer.Consume(ctx, s.OrderHandler)
	if err != nil {
		logger.Error(ctx, "Failed to start ship assembled consumer service", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Ship assembly consumer service successfully started")
	return nil
}
