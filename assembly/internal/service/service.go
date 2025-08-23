package service

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/assembly/internal/model"
)

type OrderConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type OrderProducerService interface {
	ProduceShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}

type AssemblyService interface{}

type service struct {
	assemblyProducerService OrderProducerService
}

func NewService(assemblyProducerService OrderProducerService) *service {
	return &service{
		assemblyProducerService: assemblyProducerService,
	}
}
