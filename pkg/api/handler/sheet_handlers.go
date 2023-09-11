package handler

import (
	"log"
	"net/http"

	"strconv"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
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

func (sh *SheetHandler) PostSheetHandler(c echo.Context) error {
	s := new(entities.Sheet)
	if err := c.Bind(s); err != nil {
		return err
	}

	sr := repository.NewSheetRepository(s)
	err := sr.Create(sh.storage)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, s)
}

func (sh *SheetHandler) GetSheetHandler(c echo.Context) error {
	s := new(entities.Sheet)
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		return err
	}

	sr := repository.NewSheetRepository(s)
	errf := sr.FindById(sh.storage, id)
	if errf != nil {
		return errf
	}

	return c.JSON(http.StatusOK, s)
}

func (sh *SheetHandler) GetSheetListHandler(c echo.Context) error {
	s := new(entities.Sheet)

	if err := c.Bind(s); err != nil {
		return err
	}

	p, errp := strconv.Atoi(c.QueryParams().Get("page"))
	sz, errs := strconv.Atoi(c.QueryParams().Get("size"))

	if errp != nil {
		return errp
	}

	if errs != nil {
		return errs
	}

	sr := repository.NewSheetRepository(s)
	sl, err := sr.FindAll(sh.storage, p, sz)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sl)
}

func (sh *SheetHandler) PatchSheetHandler(c echo.Context) error {
	s := new(entities.Sheet)
	os := new(entities.SheetUpdate)
	id, errp := uuid.Parse(c.Param("id"))

	if errp != nil {
		return errp
	}

	if err := c.Bind(os); err != nil {
		log.Println(err)
		return err
	}

	sr := repository.NewSheetRepository(s)
	err := sr.FindById(sh.storage, id)

	if err != nil {
		return err
	}

	su, erru := sr.Update(sh.storage, os, id)

	if erru != nil {
		return erru
	}

	return c.JSON(http.StatusOK, su)
}

func (sh *SheetHandler) DeleteSheetHandler(c echo.Context) error {
	s := new(entities.Sheet)
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		return err
	}

	sr := repository.NewSheetRepository(s)
	sr.Delete(sh.storage, id)

	return c.NoContent(http.StatusOK)
}
