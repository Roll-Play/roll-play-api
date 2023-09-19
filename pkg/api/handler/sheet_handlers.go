package handler

import (
	"fmt"
	"log"
	"net/http"

	"strconv"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	"github.com/Roll-play/roll-play-backend/pkg/errors"
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

func (sh *SheetHandler) CreateSheetHandler(c echo.Context) error {
	sd := new(entities.SheetDto)
	if err := c.Bind(sd); err != nil {
		log.Println("Error binding body: ", err)
		return customError(c, http.StatusBadRequest, errors.DtoError)
	}

	sr := repository.NewSheetRepository(sh.storage)
	ns, err := sr.Create(sd)
	if err != nil {
		return customError(c, http.StatusInternalServerError, errors.SavingError, "sheet", sd)
	}
	return c.JSON(http.StatusCreated, ns)
}

func (sh *SheetHandler) GetSheetHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("Error parsing id as uuid", err)
		return customError(c, http.StatusBadRequest, errors.ParseError, "id")
	}

	sr := repository.NewSheetRepository(sh.storage)
	s, err := sr.FindById(id)
	if err != nil {
		return customError(c, http.StatusBadRequest, errors.NotFound, "id", id)
	}

	return c.JSON(http.StatusOK, s)
}

func (sh *SheetHandler) GetSheetListHandler(c echo.Context) error {
	p, err := strconv.Atoi(c.QueryParams().Get("page"))
	if err != nil {
		log.Println("Error converting page from string to int", err)
		return customError(c, http.StatusBadRequest, errors.QueryParamError, "page")
	}

	sz, err := strconv.Atoi(c.QueryParams().Get("size"))
	if err != nil {
		log.Println("Error converting size from string to int", err)
		return customError(c, http.StatusBadRequest, errors.QueryParamError, "size")
	}

	sr := repository.NewSheetRepository(sh.storage)
	sl, err := sr.FindAll(p, sz)
	if err != nil {
		return customError(c, http.StatusInternalServerError, errors.DbError, "findAll")
	}

	return c.JSON(http.StatusOK, sl)
}

func (sh *SheetHandler) UpdateSheetHandler(c echo.Context) error {
	os := new(entities.SheetDto)

	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("Error parsing id as uuid", err)
		return customError(c, http.StatusBadRequest, errors.ParseError, "id")
	}

	if err := c.Bind(os); err != nil {
		log.Println("Error binding body: ", err)
		return customError(c, http.StatusBadRequest, errors.DtoError)
	}

	sr := repository.NewSheetRepository(sh.storage)
	_, err = sr.FindById(id)
	if err != nil {
		return customError(c, http.StatusBadRequest, errors.NotFound, "id", id)
	}

	su, err := sr.Update(os, id)
	if err != nil {
		return customError(c, http.StatusInternalServerError, errors.DbError, "update")
	}

	return c.JSON(http.StatusOK, su)
}

func (sh *SheetHandler) DeleteSheetHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println("Error parsing id as uuid", err)
		return customError(c, http.StatusBadRequest, errors.ParseError, "id")
	}

	sr := repository.NewSheetRepository(sh.storage)

	_, err = sr.FindById(id)
	if err != nil {
		return customError(c, http.StatusBadRequest, errors.NotFound, "id", id)
	}

	err = sr.Delete(id)
	if err != nil {
		return customError(c, http.StatusInternalServerError, errors.DbError, "delete")
	}

	return c.NoContent(http.StatusNoContent)
}

func customError(c echo.Context, https int, message string, args ...any) error {
	errs := errors.Error{
		Error:   http.StatusText(https),
		Message: fmt.Sprintf(message, args...),
	}
	return c.JSON(https, errs)
}
