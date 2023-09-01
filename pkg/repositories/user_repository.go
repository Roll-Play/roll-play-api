package repository

import (
	"github.com/Roll-play/roll-play-backend/pkg/entities"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	User *entities.User
}

func (ur *UserRepository) Create(db *sqlx.DB) error {

	_, err := db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", "fizi", "fizi@gmail.com", "123123")
	
	if err != nil {
		return err
	}

	return nil
}

func NewUserRepository(entity *entities.User) *UserRepository {
	return &UserRepository{
		User: entity,
	}
}