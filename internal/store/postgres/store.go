package postgres

import (
	"database/sql"
	"fmt"

	"twitter/internal/config"
	"twitter/internal/store"

	_ "github.com/lib/pq"
)

type PostgresStore struct {
	db     *sql.DB
	config *config.Config
}

func NewPostgresStore(c *config.Config) store.Store {
	return &PostgresStore{
		config: c,
	}
}

func (ps *PostgresStore) Open() error {
	db, err := sql.Open(
		"postgres",
		fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", ps.config.Postgres.User, ps.config.Postgres.PWD, ps.config.Postgres.Host, ps.config.Postgres.Port, ps.config.Postgres.DB, ps.config.Postgres.SSLMode),
	)
	if err != nil {
		return fmt.Errorf("error opening database: %v", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("error pinging database: %v", err)
	}
	ps.db = db
	return nil
}
func (ps *PostgresStore) Close() error {
	return ps.db.Close()
}
