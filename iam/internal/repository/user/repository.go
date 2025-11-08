package user

import (
	"database/sql"

	def "github.com/alexis871aa/microservices-rocket-factory/iam/internal/repository"
)

var _ def.UserRepository = (*repository)(nil)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {
	return &repository{
		db: db,
	}
}
