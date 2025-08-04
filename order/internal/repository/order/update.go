package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func (r *repository) Update(ctx context.Context, orderUUID string, newOrder model.Order) error {
	builder := sq.Update("orders").
		Where(sq.Eq{"order_uuid": orderUUID}).
		PlaceholderFormat(sq.Dollar).
		Set("user_uuid", newOrder.UserUUID).
		Set("part_uuids", pq.Array(newOrder.PartUuids)).
		Set("total_price", newOrder.TotalPrice).
		Set("transaction_uuid", newOrder.TransactionUUID).
		Set("payment_method", newOrder.PaymentMethod).
		Set("status", newOrder.Status)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build update query: %v", err)
		return err
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("failed to update order: %v", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrOrderNotFound
	}

	log.Printf("updated %d rows for order: %s", rowsAffected, orderUUID)
	return nil
}
