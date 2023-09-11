package entities

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Sheet struct {
	Id          uuid.UUID    `db:"id" json:"id"`
	Name        string       `db:"name" json:"name"`
	Description string       `db:"description" json:"description"`
	Properties  string       `db:"properties" json:"properties"`
	Background  string       `db:"background" json:"background"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}

type SheetUpdate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Properties  string `json:"properties"`
	Background  string `json:"background"`
}
