package order

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

const inventoryTimeout = 5 * time.Second

func (s *service) Create(ctx context.Context, userUUID string, partUUIDs []string) (*model.Order, error) {
	if len(partUUIDs) == 0 {
		return nil, model.ErrPartsNotFound
	}

	ctxWithTimeout, cancel := context.WithTimeout(ctx, inventoryTimeout)
	defer cancel()

	parts, err := s.inventoryClient.ListParts(ctxWithTimeout, model.PartsFilter{Uuids: &partUUIDs})
	if err != nil {
		return nil, err
	}

	if len(parts) != len(partUUIDs) {
		return nil, model.ErrPartsNotFound
	}

	var totalPrice float32
	for _, part := range parts {
		totalPrice += float32(part.Price)
	}

	order := model.Order{
		OrderUUID:  uuid.NewString(),
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: totalPrice,
		Status:     model.StatusPendingPayment,
		CreatedAt:  lo.ToPtr(time.Now()),
		UpdatedAt:  nil,
	}

	oerr := s.orderRepository.Create(ctx, order)
	if oerr != nil {
		return nil, oerr
	}
	return &order, nil
}
