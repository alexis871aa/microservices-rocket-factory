package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/alexis871aa/microservices-rocket-factory/payment/internal/service"
)

func (s *Service) PayOrder(_ context.Context, orderUUID, userUUID string, paymentMethod service.PaymentMethod) (string, error) {
	transactionUUID := uuid.NewString()
	log.Printf("Оплата прошла успешно, transactionUUID: %s, orderUUID: %s, userUUID: %s, paymentMethod: %d",
		transactionUUID, orderUUID, userUUID, paymentMethod)

	return transactionUUID, nil
}
