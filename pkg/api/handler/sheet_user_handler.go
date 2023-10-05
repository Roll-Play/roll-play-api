package handler

import (
	"log"
	"net/http"

	api_error "github.com/Roll-play/roll-play-backend/pkg/errors"
	repository "github.com/Roll-play/roll-play-backend/pkg/repositories"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type SheetUserHandler struct {
	storage *sqlx.DB
}

type PermissionList struct {
	Permission int         `json:"permission"`
	UserIds    []uuid.UUID `json:"user_ids"`
}

func NewSheetUserHandler(storage *sqlx.DB) *SheetUserHandler {
	return &SheetUserHandler{
		storage,
	}
}

func (sheetUserHandler *SheetUserHandler) UpdatePermissions(context echo.Context) error {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		log.Println("Error parsing id as uuid", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.PARSE_ERROR, "id")
	}

	userId := context.Get("user").(uuid.UUID)
	sheetRepository := repository.NewSheetRepository(sheetUserHandler.storage)
	savedSheet, err := sheetRepository.FindByIdAndUserId(id, userId)
	if err != nil {
		return api_error.CustomError(context, http.StatusBadRequest, api_error.NOT_FOUND, "id", id)
	}

	if !savedSheet.Owner {
		return api_error.CustomError(context, http.StatusForbidden, api_error.PERMISSION_ERROR, userId)
	}

	permissionList := new(PermissionList)
	if err := context.Bind(permissionList); err != nil {
		log.Println("Error binding body: ", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.DTO_ERROR)
	}

	sheetUserRepository := repository.NewSheetUserRepository(sheetUserHandler.storage)

	for _, user := range permissionList.UserIds {
		_, err = sheetUserRepository.CreateSheetUserRelation(id, user, permissionList.Permission, false)
		if err != nil {
			return api_error.CustomError(context, http.StatusInternalServerError, api_error.SAVING_ERROR, "sheetUser", user)
		}
	}

	return context.JSON(http.StatusCreated, savedSheet)
}
