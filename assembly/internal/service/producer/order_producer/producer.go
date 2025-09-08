package order_producer

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	"github.com/alexis871aa/microservices-rocket-factory/assembly/internal/model"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/kafka"
	"github.com/alexis871aa/microservices-rocket-factory/platform/pkg/logger"
	eventsV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/proto/events/v1"
)

type service struct {
	orderProducer kafka.Producer
}

func NewService(orderProducer kafka.Producer) *service {
	return &service{
		orderProducer: orderProducer,
	}
}

func (p *service) ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error {
	msg := &eventsV1.ShipAssembled{
		EventUuid:    event.EventUUID,
		OrderUuid:    event.OrderUUID,
		UserUuid:     event.UserUUID,
		BuildTimeSec: event.BuildTimeSec,
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		logger.Error(ctx, "failed to marshal OrderAssembled", zap.Error(err))
		return err
	}

	err = p.orderProducer.Send(ctx, []byte(event.EventUUID), payload)
	if err != nil {
		logger.Error(ctx, "failed to publish order assembled", zap.Error(err))
		return err
	}

	return nil
}
