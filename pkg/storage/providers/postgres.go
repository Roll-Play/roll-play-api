package storage_providers

import (
	"time"

	repository_interfaces "github.com/Roll-play/roll-play-backend/pkg/repositories/interfaces"
	"github.com/jmoiron/sqlx"
)

type PostgresProvider struct {
	DB *sqlx.DB
}

type PostgresProviderAdapter struct {
	*PostgresProvider
}

func (pp *PostgresProvider) Connect(connectionString string) error {
	db, err := newPostgresDB(connectionString)

	if err != nil {
		return err
	}
	pp.DB = db
	return nil
}

func (pp *PostgresProvider) Create(repository repository_interfaces.Repository) error {
	err := repository.Create(pp.DB)

	if err != nil {
		return err
	}

	return nil
}

func (pp *PostgresProvider) Db() *sqlx.DB {
	return pp.DB
}

func newPostgresDB(connString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", connString)
	
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * time.Duration(5))

	return db, err
} 
