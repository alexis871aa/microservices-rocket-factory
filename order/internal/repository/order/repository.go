package order

import (
	"database/sql"

	def "github.com/alexis871aa/microservices-rocket-factory/order/internal/repository"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}
