package order

import (
	"context"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func (r *repository) Delete(_ context.Context, orderUUID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.data[orderUUID]; !ok {
		return model.ErrOrderNotFound
	}
	delete(r.data, orderUUID)
	return nil
}
