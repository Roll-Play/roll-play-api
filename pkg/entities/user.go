package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID    `db:"id" json:"id"`
	Username  string       `db:"username" json:"username"`
	Email     string       `db:"email" json:"email"`
	Password  string       `db:"password" json:"password"`
	IsActive  bool         `db:"is_active" json:"is_active"`
	CreatedAt time.Time    `db:"created_at" json:"created_at"`
	UpdatedAt time.Time    `db:"updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at" json:"deleted_at"`
}
