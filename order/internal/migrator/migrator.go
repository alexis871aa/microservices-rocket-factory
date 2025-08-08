package migrator

import (
	"database/sql"

	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db           *sql.DB
	migrationDir string
}

func NewMigrator(db *sql.DB, migrationDir string) *Migrator {
	return &Migrator{
		db:           db,
		migrationDir: migrationDir,
	}
}

func (m *Migrator) Up() error {
	err := goose.Up(m.db, m.migrationDir)
	if err != nil {
		return err
	}

	return nil
}
