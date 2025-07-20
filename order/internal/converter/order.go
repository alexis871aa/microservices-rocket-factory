package converter

import (
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	orderV1 "github.com/alexis871aa/microservices-rocket-factory/shared/pkg/openapi/order/v1"
)

func PaymentMethodToModel(apiMethod orderV1.PaymentMethod) model.PaymentMethod {
	switch apiMethod {
	case orderV1.PaymentMethodCARD:
		return model.PaymentMethodCard
	case orderV1.PaymentMethodSBP:
		return model.PaymentMethodSBP
	case orderV1.PaymentMethodCREDITCARD:
		return model.PaymentMethodCreditCard
	case orderV1.PaymentMethodINVESTORMONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnknown
	}
}

func OrderStatusToAPI(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.StatusPendingPayment:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.StatusPaid:
		return orderV1.OrderStatusPAID
	case model.StatusCancelled:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusUNKNOWN
	}
}

func OrderToDTO(order *model.Order) orderV1.OrderDto {
	dto := orderV1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		Status:     OrderStatusToAPI(order.Status),
	}

	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderV1.NewOptNilString(*order.TransactionUUID)
	}

	if order.PaymentMethod != nil {
		apiPaymentMethod := PaymentMethodToOrderDtoPaymentMethod(*order.PaymentMethod)
		dto.PaymentMethod = &orderV1.NilOrderDtoPaymentMethod{
			Value: apiPaymentMethod,
			Null:  false,
		}
	}

	return dto
}

func PaymentMethodToOrderDtoPaymentMethod(domainMethod model.PaymentMethod) orderV1.OrderDtoPaymentMethod {
	switch domainMethod {
	case model.PaymentMethodCard:
		return orderV1.OrderDtoPaymentMethodCARD
	case model.PaymentMethodSBP:
		return orderV1.OrderDtoPaymentMethodSBP
	case model.PaymentMethodCreditCard:
		return orderV1.OrderDtoPaymentMethodCREDITCARD
	case model.PaymentMethodInvestorMoney:
		return orderV1.OrderDtoPaymentMethodINVESTORMONEY
	default:
		return orderV1.OrderDtoPaymentMethodUNKNOWN
	}
}
