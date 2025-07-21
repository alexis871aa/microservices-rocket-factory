package converter

import (
	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
	repoModel "github.com/alexis871aa/microservices-rocket-factory/order/internal/repository/model"
)

func OrderToModel(order repoModel.Order) model.Order {
	return model.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   convertPaymentMethodToModel(order.PaymentMethod),
		Status:          model.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func convertPaymentMethodToModel(repoMethod *repoModel.PaymentMethod) *model.PaymentMethod {
	if repoMethod == nil {
		return nil
	}

	domainMethod := model.PaymentMethod(*repoMethod)
	return &domainMethod
}

func ModelToOrder(order model.Order) repoModel.Order {
	return repoModel.Order{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   convertPaymentMethodToRepo(order.PaymentMethod),
		Status:          repoModel.OrderStatus(order.Status),
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}

func convertPaymentMethodToRepo(paymentMethod *model.PaymentMethod) *repoModel.PaymentMethod {
	if paymentMethod == nil {
		return nil
	}

	repoMethod := repoModel.PaymentMethod(*paymentMethod)
	return &repoMethod
}
