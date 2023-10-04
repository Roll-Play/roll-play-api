package repository

import (
	"log"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SheetUserRepository struct {
	db *sqlx.DB
}

func NewSheetUserRepository(connection *sqlx.DB) *SheetUserRepository {
	return &SheetUserRepository{
		db: connection,
	}
}

func (sur *SheetUserRepository) CreateSheetUserRelation(sheetId uuid.UUID, userId uuid.UUID, permission int, owner bool) (*entities.SheetUser, error) {
	sheetUser := new(entities.SheetUser)
	err := sur.db.Get(sheetUser,
		`INSERT INTO sheet_user (sheet_id, user_id, permission, owner) 
			VALUES ($1, $2, $3, $4) 
			RETURNING *`,
		sheetId.String(), userId.String(), permission, owner)

	if err != nil {
		log.Println("Error creating record for sheet: ", sheetId)
		return nil, err
	}

	return sheetUser, nil
}
