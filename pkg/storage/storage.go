package storage

import (
	repository_interfaces "github.com/Roll-play/roll-play-backend/pkg/repositories/interfaces"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Provider interface {
	Connect(connString string) error
	Create(repository_interfaces.Repository) error
}



type Storage struct {
	Provider *Provider
} 

func NewStorage(connString string, provider Provider) (*Storage, error) {
	err := provider.Connect(connString)

	if err != nil {
		return nil, err
	}

	return &Storage{
		Provider: &provider,
	}, nil
}