package service

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/assembly/internal/model"
)

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembled) error
}
