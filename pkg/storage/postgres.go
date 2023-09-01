package storage

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type PostgresProvider struct {
	DB *sqlx.DB
}

func (pp *PostgresProvider) Connect(connectionString string) error {
	db, err := newPostgresDB(connectionString)

	if err != nil {
		return err
	}
	pp.DB = db
	return nil
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
