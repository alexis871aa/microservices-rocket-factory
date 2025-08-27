package kafka

import "github.com/alexis871aa/microservices-rocket-factory/notification/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaid, error)
}

type ShipAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembled, error)
}
