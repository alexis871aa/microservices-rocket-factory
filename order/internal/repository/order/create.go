package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func (r *repository) Create(ctx context.Context, order model.Order) error {
	builder := sq.Insert("orders").
		PlaceholderFormat(sq.Dollar).
		Columns("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status").
		Values(order.OrderUUID, order.UserUUID, pq.Array(order.PartUuids), order.TotalPrice, order.TransactionUUID, order.PaymentMethod, order.Status)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return err
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("failed to execute query: %v\n", err)
	}
	return err
}
