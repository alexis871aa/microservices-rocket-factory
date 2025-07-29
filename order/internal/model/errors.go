package model

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrPartsNotFound      = errors.New("not all parts found")
	ErrOrderAlreadyPaid   = errors.New("order already paid")
	ErrOrderCancelled     = errors.New("order cancelled")
	ErrInvalidOrderStatus = errors.New("invalid order status for operation")
)
