package handler

import (
	"log"
	"net/http"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	api_error "github.com/Roll-play/roll-play-backend/pkg/errors"
	repository "github.com/Roll-play/roll-play-backend/pkg/repositories"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type GroupHandler struct {
	storage *sqlx.DB
}

func NewGroupHandler(storage *sqlx.DB) *GroupHandler {
	return &GroupHandler{
		storage,
	}
}

func (groupHandler *GroupHandler) CreateGroupHandler(context echo.Context) error {
	groupDto := new(entities.GroupDto)

	if err := context.Bind(groupDto); err != nil {
		log.Println("Error binding body: ", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.DTO_ERROR)
	}

	userId := context.Get("user").(uuid.UUID)
	groupRepository := repository.NewGroupRepository(groupHandler.storage)
	savedGroup, err := groupRepository.Create(groupDto, userId)
	if err != nil {
		return api_error.CustomError(context, http.StatusInternalServerError, api_error.SAVING_ERROR, "group", groupDto)
	}
	return context.JSON(http.StatusCreated, savedGroup)

}

func (groupHandler *GroupHandler) GetGroupHandler(context echo.Context) error {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		log.Println("Error parsing id as uuid", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.PARSE_ERROR, "id")
	}

	userId := context.Get("user").(uuid.UUID)
	groupRepository := repository.NewGroupRepository(groupHandler.storage)
	group, err := groupRepository.FindById(id)
	if err != nil {
		return api_error.CustomError(context, http.StatusBadRequest, api_error.NOT_FOUND, "id", id)
	}

	if group.UserId != userId {
		return api_error.CustomError(context, http.StatusBadRequest, api_error.NOT_FOUND, "id", userId)
	}

	return context.JSON(http.StatusOK, group)

}

func (groupHandler *GroupHandler) DeleteGroupHandler(context echo.Context) error {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		log.Println("Error parsing id as uuid", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.PARSE_ERROR, "id")
	}

	userId := context.Get("user").(uuid.UUID)
	groupRepository := repository.NewGroupRepository(groupHandler.storage)
	group, err := groupRepository.FindById(id)
	if err != nil {
		return api_error.CustomError(context, http.StatusBadRequest, api_error.NOT_FOUND, "id", id)
	}

	if group.UserId != userId {
		return api_error.CustomError(context, http.StatusBadRequest, api_error.NOT_FOUND, "id", userId)
	}

	err = groupRepository.DeleteGroup(id)
	if err != nil {
		return api_error.CustomError(context, http.StatusInternalServerError, api_error.DB_ERROR, "delete")
	}

	return context.NoContent(http.StatusNoContent)
}
