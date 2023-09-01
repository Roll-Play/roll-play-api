package user

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id uuid.UUID `db:"id"`
	Username string `db:"username"`
	Email string `db:"email"`
	Password string `db:"password"`
	IsActive bool `db:"is_active"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time `db:"deleted_at"`
}

