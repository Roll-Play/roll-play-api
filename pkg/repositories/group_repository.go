package repository

import (
	"log"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type GroupRepository struct {
	db *sqlx.DB
}

func NewGroupRepository(connection *sqlx.DB) *GroupRepository {
	return &GroupRepository{
		db: connection,
	}
}

func (gr *GroupRepository) Create(groupDto *entities.GroupDto, userId uuid.UUID) (*entities.Group, error) {
	group := new(entities.Group)
	err := gr.db.Get(group,
		`INSERT INTO player_group (name, public, user_id) 
			VALUES ($1, $2, $3) 
			RETURNING *`,
		groupDto.Name, groupDto.Public, userId)

	if err != nil {
		log.Println("Error creating record with dto: ", groupDto)
		return nil, err
	}

	return group, nil
}

func (gr *GroupRepository) FindById(groupId uuid.UUID) (*entities.Group, error) {
	group := new(entities.Group)
	err := gr.db.Get(group, `SELECT * FROM player_group WHERE id=$1`, groupId)

	if err != nil {
		log.Println("Error finding record with id: ", groupId)
		return nil, err
	}

	return group, nil
}

func (gr *GroupRepository) DeleteGroup(id uuid.UUID) error {
	_, err := gr.db.Exec("DELETE FROM player_group WHERE id=$1", id)
	if err != nil {
		log.Println("Error deleting record with id:", id)
		return err
	}

	return nil
}
