package payment

import def "github.com/alexis871aa/microservices-rocket-factory/payment/internal/service"

var _ def.PaymentService = (*Service)(nil)

type Service struct{}

func NewService() *Service {
	return &Service{}
}
