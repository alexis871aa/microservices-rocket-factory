package kafka

import "github.com/alexis871aa/microservices-rocket-factory/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembled, error)
}
