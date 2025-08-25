package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/converter"
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

const paymentTimeout = 3 * time.Second

func (s *service) Pay(ctx context.Context, orderUUID string, paymentMethod model.PaymentMethod) (string, error) {
	ctxWithTimeout, cancel := context.WithTimeout(ctx, paymentTimeout)
	defer cancel()

	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return "", model.ErrOrderNotFound
	}

	if order.Status != model.StatusPendingPayment {
		switch order.Status {
		case model.StatusPaid:
			return "", model.ErrOrderAlreadyPaid
		case model.StatusCancelled:
			return "", model.ErrOrderCancelled
		default:
			return "", model.ErrInvalidOrderStatus
		}
	}

	transactionUUID, cerr := s.paymentClient.PayOrder(ctxWithTimeout, orderUUID, order.UserUUID, paymentMethod)
	if cerr != nil {
		return "", cerr
	}

	order.Status = model.StatusPaid
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &paymentMethod
	order.UpdatedAt = lo.ToPtr(time.Now())

	err = s.orderRepository.Update(ctx, orderUUID, *order)
	if err != nil {
		return "", err
	}

	err = s.orderProduceService.ProduceOrderPaid(ctx, model.OrderPaid{
		EventUUID:       uuid.NewString(),
		OrderUUID:       orderUUID,
		UserUUID:        order.UserUUID,
		PaymentMethod:   converter.PaymentMethodToString(paymentMethod),
		TransactionUUID: transactionUUID,
	})
	if err != nil {
		return "", err
	}

	return transactionUUID, nil
}
