package order

import (
	"context"
	"log"

	sq "github.com/Masterminds/squirrel"

	"github.com/alexis871aa/microservices-rocket-factory/order/internal/model"
)

func (r *repository) Delete(ctx context.Context, orderUUID string) error {
	builder := sq.Delete("orders").
		Where(sq.Eq{"order_uuid": orderUUID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		log.Printf("failed to build query: %v\n", err)
		return err
	}

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		log.Printf("failed to delete order: %v\n", err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return model.ErrOrderNotFound
	}

	return nil
}
