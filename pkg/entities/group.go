package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Group struct {
	Id        uuid.UUID    `db:"id" json:"id"`
	UserId    uuid.UUID    `db:"user_id" json:"user_id"`
	Name      string       `db:"name" json:"name"`
	Public    bool         `db:"public" json:"public"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type GroupDto struct {
	Name   string `json:"name"`
	Public bool   `json:"public"`
}
