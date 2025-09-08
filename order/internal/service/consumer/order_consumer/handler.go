package order_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderShipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode order event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("event_uuid", event.EventUUID),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("user_uuid", event.UserUUID),
		zap.Int64("build_time_sec", event.BuildTimeSec),
	)

	order, err := s.orderRepository.Get(ctx, event.OrderUUID)
	if err != nil {
		logger.Error(ctx, "Failed to get order", zap.String("order_uuid", event.OrderUUID), zap.Error(err))
		return err
	}
	order.Status = model.StatusCompleted

	err = s.orderRepository.Update(ctx, event.OrderUUID, *order)
	if err != nil {
		logger.Error(ctx, "Failed to update order", zap.String("order_uuid", event.OrderUUID), zap.Error(err))
		return err
	}

	return nil
}
