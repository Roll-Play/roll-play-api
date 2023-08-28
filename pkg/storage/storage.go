package storage

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/jmoiron/sqlx"
)

func NewDB(connString string, maxOpenConns int, maxIdleConns int, maxConnLifeTime int) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", connString)
	
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(time.Minute * time.Duration(maxConnLifeTime))

	return db, err
} 