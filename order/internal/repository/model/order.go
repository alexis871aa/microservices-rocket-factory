package model

import "time"

type Order struct {
	OrderUUID       string         `json:"order_uuid"`
	UserUUID        string         `json:"user_uuid"`
	PartUuids       []string       `json:"part_uuids"`
	TotalPrice      float32        `json:"total_price"`
	TransactionUUID *string        `json:"transaction_uuid"`
	PaymentMethod   *PaymentMethod `json:"payment_method"`
	Status          OrderStatus    `json:"status"`
	CreatedAt       *time.Time     `json:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at"`
}

type OrderStatus int

const (
	StatusPendingPayment OrderStatus = iota
	StatusPaid
	StatusCancelled
)

type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)
