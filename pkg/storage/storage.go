package storage

import (
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Provider interface {
	Connect(connString string) error
}

type Storage struct {
	DB *Provider  // *sqlx.DB
} 

func NewStorage(connString string, provider Provider) (*Storage, error) {
	err := provider.Connect(connString)

	if err != nil {
		return nil, err
	}

	return &Storage{
		DB: &provider,
	}, nil
}