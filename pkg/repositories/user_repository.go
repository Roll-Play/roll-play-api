package repository

import (
	"errors"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func (userRepository *UserRepository) Create(user entities.User) (*entities.User, error) {

	row := userRepository.db.QueryRowx(
		"INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id, username, email, created_at, updated_at",
		user.Username,
		user.Email,
		user.Password,
	)

	if row.Err() != nil {
		return nil, row.Err()
	}

	u := new(entities.User)

	row.Scan(&u.Id, &u.Username, &u.Email, &u.CreatedAt, &u.UpdatedAt)

	return u, nil
}

func (userRepository *UserRepository) FindByEmail(email string) (*entities.User, error) {
	u := new(entities.User)
	err := userRepository.db.Get(u, "SELECT * FROM users WHERE email=$1", email)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (userRepository *UserRepository) FindByUsername(username string) (*entities.User, error) {
	u := new(entities.User)
	err := userRepository.db.Get(u, "SELECT * FROM users where username=$1", username)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (userRepository *UserRepository) FindById(id uuid.UUID) (*entities.User, error) {
	u := new(entities.User)
	err := userRepository.db.Get(u, "SELECT * FROM users where id=$1", id)

	if err != nil {
		return nil, err
	}

	return u, nil
}

func (userRepository *UserRepository) IsActive(id uuid.UUID) error {
	var count int

	err := userRepository.db.Get(&count, "SELECT COUNT(id) FROM users WHERE id=$1 AND is_active=true", id)

	if err != nil {
		return err
	}

	if count < 1 {
		return errors.New("user is not active")
	}

	return nil
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}
