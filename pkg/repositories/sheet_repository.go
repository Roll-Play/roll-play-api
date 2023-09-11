package repository

import (
	"fmt"
	"reflect"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SheetRepository struct {
	sheet *entities.Sheet
}

func (sr *SheetRepository) Create(db *sqlx.DB) error {
	err := db.Get(sr.sheet, "INSERT INTO sheets (name, description, properties, background) VALUES ($1, $2, $3, $4) RETURNING *",
		sr.sheet.Name, sr.sheet.Description, sr.sheet.Properties, sr.sheet.Background)

	if err != nil {
		return err
	}

	return nil

}

func NewSheetRepository(entity *entities.Sheet) *SheetRepository {
	return &SheetRepository{
		sheet: entity,
	}
}

func (sr *SheetRepository) FindById(db *sqlx.DB, id uuid.UUID) error {
	err := db.Get(sr.sheet, "SELECT s.* FROM sheets s WHERE s.id = $1", id)

	if err != nil {
		return err
	}

	return nil
}

func (sr *SheetRepository) FindAll(db *sqlx.DB, page int, size int) (*[]entities.Sheet, error) {
	r := []entities.Sheet{}

	err := db.Select(&r, "SELECT s.* FROM sheets s ORDER BY s.name LIMIT $1 OFFSET $2", size, page*size)

	if err != nil {
		return &r, err
	}

	return &r, nil
}

func (sr *SheetRepository) Update(db *sqlx.DB, os *entities.SheetUpdate, id uuid.UUID) (*entities.Sheet, error) {
	v := reflect.ValueOf(os)
	rv := reflect.Indirect(v)
	types := rv.Type()
	sqlQuery := "UPDATE sheets SET "
	var args []interface{}
	count := 0
	for i := 0; i < rv.NumField(); i++ {
		if rv.Field(i).String() != "" {
			count += 1
			sqlQuery += types.Field(i).Name + fmt.Sprintf("=$%d, ", count)
			args = append(args, rv.Field(i).String())
		}
	}

	sqlQuery = sqlQuery[:len(sqlQuery)-2] + fmt.Sprintf(" WHERE id=$%d RETURNING *", count+1)
	args = append(args, id.String())

	us := new(entities.Sheet)
	err := db.Get(us, sqlQuery, args...)

	if err != nil {
		return us, err
	}

	return us, nil
}

func (sr *SheetRepository) Delete(db *sqlx.DB, id uuid.UUID) error {
	_, err := db.Exec("DELETE FROM sheets WHERE id=$1", id)

	if err != nil {
		return err
	}

	return nil
}
