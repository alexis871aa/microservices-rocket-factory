package order_assembled_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.shipAssembledDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode order assembled event", zap.Error(err))
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

	err = s.telegramService.SendShipAssembledNotification(ctx, event)
	if err != nil {
		logger.Error(ctx, "Failed to send ship assembled notification", zap.Error(err))
		return err
	}

	return nil
}
