package entities

import "github.com/google/uuid"

type SheetUser struct {
	SheetId    uuid.UUID `db:"sheet_id" json:"sheet_id"`
	UserId     uuid.UUID `db:"user_id" json:"user_id"`
	Permission int       `db:"permission" json:"permission"`
	Owner      bool      `db:"owner" json:"owner"`
}
