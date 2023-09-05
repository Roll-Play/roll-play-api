package repository

import (
	"github.com/Roll-play/roll-play-backend/pkg/entities"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	user *entities.User
}

func (ur *UserRepository) Create(db *sqlx.DB) error {

	err := db.QueryRowx(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, created_at, updated_at",
		ur.user.Username,
		ur.user.Email,
		ur.user.Password,
	).Scan(&ur.user.Id, &ur.user.CreatedAt, &ur.user.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

func NewUserRepository(entity *entities.User) *UserRepository {
	return &UserRepository{
		user: entity,
	}
}
