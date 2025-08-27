package order_paid_consumer

import (
	"context"

	"go.uber.org/zap"

	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
)

func (s *service) OrderHandler(ctx context.Context, msg kafka.Message) error {
	event, err := s.orderPaidDecoder.Decode(msg.Value)
	if err != nil {
		logger.Error(ctx, "Failed to decode order paid event", zap.Error(err))
		return err
	}

	logger.Info(ctx, "Processing message",
		zap.String("topic", msg.Topic),
		zap.Any("partition", msg.Partition),
		zap.Any("offset", msg.Offset),
		zap.String("order_uuid", event.OrderUUID),
		zap.String("event_uuid", event.EventUUID),
		zap.String("payment_method", event.PaymentMethod),
		zap.String("transaction_uuid", event.TransactionUUID),
	)

	// дальше добавить логику по отправке сообщение в ТГ (Что заказ оплачен)

	return nil
}
