package repository

import (
	"fmt"
	"log"
	"reflect"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type SheetRepository struct {
	db *sqlx.DB
}

func NewSheetRepository(connection *sqlx.DB) *SheetRepository {
	return &SheetRepository{
		db: connection,
	}
}

func (sr *SheetRepository) Create(sheetDto *entities.SheetDto, userId uuid.UUID) (*entities.Sheet, error) {
	sheet := new(entities.Sheet)
	err := sr.db.Get(sheet,
		`INSERT INTO sheets (name, description, properties, background, user_id) 
			VALUES ($1, $2, $3, $4, $5) 
			RETURNING *`,
		sheetDto.Name, sheetDto.Description, sheetDto.Properties, sheetDto.Background, userId)

	if err != nil {
		log.Println("Error creating record with dto: ", sheetDto)
		return nil, err
	}

	return sheet, nil

}

func (sr *SheetRepository) FindByIdAndUserId(id uuid.UUID, userId uuid.UUID) (*entities.Sheet, error) {
	sheet := new(entities.Sheet)
	err := sr.db.Get(sheet, "SELECT s.* FROM sheets s WHERE s.id = $1 AND s.user_id = $2", id, userId)

	if err != nil {
		log.Println("Error finding record with id: ", id)
		return nil, err
	}

	return sheet, nil
}

func (sr *SheetRepository) FindAllByUserId(page int, size int, userId uuid.UUID) (*[]entities.Sheet, error) {
	r := []entities.Sheet{}

	err := sr.db.Select(&r,
		`SELECT s.* FROM sheets s 
			WHERE s.user_id = $1 
			ORDER BY s.name 
			LIMIT $2 OFFSET $3`,
		userId, size, page*size)
	if err != nil {
		log.Println("Error finding records for page with page and size:", page, size)
		return &r, err
	}

	return &r, nil
}

func (sr *SheetRepository) Update(sheetDto *entities.SheetDto, id uuid.UUID) (*entities.Sheet, error) {
	v := reflect.ValueOf(sheetDto)
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
	err := sr.db.Get(us, sqlQuery, args...)
	if err != nil {
		log.Println("Error updating record with id:", id)
		log.Println("and dto:", sheetDto)
		return nil, err
	}

	return us, nil
}

func (sr *SheetRepository) Delete(id uuid.UUID) error {
	_, err := sr.db.Exec("DELETE FROM sheets WHERE id=$1", id)

	if err != nil {
		log.Println("Error deleting record with id:", id)
		return err
	}

	return nil
}
