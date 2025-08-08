package order

import (
	"context"
	"database/sql"
	"errors"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/lib/pq"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func (r *repository) Get(ctx context.Context, orderUUID string) (*model.Order, error) {
	builder := sq.Select("order_uuid", "user_uuid", "part_uuids", "total_price", "transaction_uuid", "payment_method", "status", "created_at", "updated_at").
		From("orders").
		Where(sq.Eq{"order_uuid": orderUUID}).
		PlaceholderFormat(sq.Dollar).
		Limit(1)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return nil, err
	}

	var order model.Order
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&order.OrderUUID,
		&order.UserUUID,
		pq.Array(&order.PartUuids),
		&order.TotalPrice,
		&order.TransactionUUID,
		&order.PaymentMethod,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrOrderNotFound
		}

		log.Printf("failed to select order: %v\n", err)
		return nil, err
	}

	return &order, nil
}
