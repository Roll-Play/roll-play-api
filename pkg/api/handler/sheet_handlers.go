package handler

import (
	"log"
	"net/http"

	"strconv"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	api_error "github.com/Roll-play/roll-play-backend/pkg/errors"
	repository "github.com/Roll-play/roll-play-backend/pkg/repositories"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type SheetHandler struct {
	storage *sqlx.DB
}

func NewSheetHandler(storage *sqlx.DB) *SheetHandler {
	return &SheetHandler{
		storage,
	}
}

func (sheetHandler *SheetHandler) CreateSheetHandler(context echo.Context) error {
	sheetDto := new(entities.SheetDto)
	if err := context.Bind(sheetDto); err != nil {
		log.Println("Error binding body: ", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.DTO_ERROR)
	}

	userId := context.Get("user").(uuid.UUID)
	sheetRepository := repository.NewSheetRepository(sheetHandler.storage)
	savedSheet, err := sheetRepository.Create(sheetDto)
	if err != nil {
		return api_error.CustomError(context, http.StatusInternalServerError, api_error.SAVING_ERROR, "sheet", sheetDto)
	}

	sheetUserRepository := repository.NewSheetUserRepository(sheetHandler.storage)
	_, err = sheetUserRepository.CreateSheetUserRelation(savedSheet.Id, userId, entities.READ, true)
	if err != nil {
		return api_error.CustomError(context, http.StatusInternalServerError, api_error.SAVING_ERROR, "sheet", sheetDto)
	}
	return context.JSON(http.StatusCreated, savedSheet)
}

func (sheetHandler *SheetHandler) GetSheetHandler(context echo.Context) error {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		log.Println("Error parsing id as uuid", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.PARSE_ERROR, "id")
	}

	userId := context.Get("user").(uuid.UUID)
	sheetRepository := repository.NewSheetRepository(sheetHandler.storage)
	sheet, err := sheetRepository.FindByIdAndUserId(id, userId)
	if err != nil {
		return api_error.CustomError(context, http.StatusBadRequest, api_error.NOT_FOUND, "id", id)
	}

	return context.JSON(http.StatusOK, sheet)
}

func (sheetHandler *SheetHandler) GetSheetListHandler(context echo.Context) error {
	page, err := strconv.Atoi(context.QueryParams().Get("page"))
	if err != nil {
		log.Println("Error converting page from string to int", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.QUERY_PARAM_ERROR, "page")
	}

	size, err := strconv.Atoi(context.QueryParams().Get("size"))
	if err != nil {
		log.Println("Error converting size from string to int", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.QUERY_PARAM_ERROR, "size")
	}

	userId := context.Get("user").(uuid.UUID)
	sheetRepository := repository.NewSheetRepository(sheetHandler.storage)
	sheetList, err := sheetRepository.FindAllByUserId(page, size, userId)
	if err != nil {
		return api_error.CustomError(context, http.StatusInternalServerError, api_error.DB_ERROR, "findAll")
	}

	return context.JSON(http.StatusOK, sheetList)
}

func (sheetHandler *SheetHandler) UpdateSheetHandler(context echo.Context) error {
	sheetDto := new(entities.SheetDto)

	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		log.Println("Error parsing id as uuid", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.PARSE_ERROR, "id")
	}

	if err := context.Bind(sheetDto); err != nil {
		log.Println("Error binding body: ", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.DTO_ERROR)
	}

	userId := context.Get("user").(uuid.UUID)
	sheetRepository := repository.NewSheetRepository(sheetHandler.storage)
	sheet, err := sheetRepository.FindByIdAndUserId(id, userId)
	if err != nil {
		return api_error.CustomError(context, http.StatusBadRequest, api_error.NOT_FOUND, "id", id)
	}

	if !sheet.Owner && sheet.Permission != entities.WRITE {
		return api_error.CustomError(context, http.StatusForbidden, api_error.PERMISSION_ERROR, userId)
	}

	updatedSheet, err := sheetRepository.Update(sheetDto, id)
	if err != nil {
		return api_error.CustomError(context, http.StatusInternalServerError, api_error.DB_ERROR, "update")
	}

	return context.JSON(http.StatusOK, updatedSheet)
}

func (sheetHandler *SheetHandler) DeleteSheetHandler(context echo.Context) error {
	id, err := uuid.Parse(context.Param("id"))
	if err != nil {
		log.Println("Error parsing id as uuid", err)
		return api_error.CustomError(context, http.StatusBadRequest, api_error.PARSE_ERROR, "id")
	}

	userId := context.Get("user").(uuid.UUID)
	sheetRepository := repository.NewSheetRepository(sheetHandler.storage)
	sheet, err := sheetRepository.FindByIdAndUserId(id, userId)
	if err != nil {
		return api_error.CustomError(context, http.StatusBadRequest, api_error.NOT_FOUND, "id", id)
	}

	if !sheet.Owner {
		return api_error.CustomError(context, http.StatusForbidden, api_error.PERMISSION_ERROR, userId)
	}

	err = sheetRepository.Delete(id)
	if err != nil {
		return api_error.CustomError(context, http.StatusInternalServerError, api_error.DB_ERROR, "delete")
	}

	return context.NoContent(http.StatusNoContent)
}
