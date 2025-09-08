package kafka

import "github.com/alexis871aa/microservices-rocket-factory/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaid, error)
}
