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

func (sh *SheetHandler) CreateSheetHandler(c echo.Context) error {
	sd := new(entities.SheetDto)
	if err := c.Bind(sd); err != nil {
		return err
	}

	sr := repository.NewSheetRepository(sh.storage)
	ns, err := sr.Create(sd)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, ns)
}

func (sh *SheetHandler) GetSheetHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		return err
	}

	sr := repository.NewSheetRepository(sh.storage)
	s, err := sr.FindById(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, s)
}

func (sh *SheetHandler) GetSheetListHandler(c echo.Context) error {
	p, err := strconv.Atoi(c.QueryParams().Get("page"))
	if err != nil {
		return err
	}

	sz, err := strconv.Atoi(c.QueryParams().Get("size"))
	if err != nil {
		return err
	}

	sr := repository.NewSheetRepository(sh.storage)
	sl, err := sr.FindAll(p, sz)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, sl)
}

func (sh *SheetHandler) UpdateSheetHandler(c echo.Context) error {
	os := new(entities.SheetDto)
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		return err
	}

	if err := c.Bind(os); err != nil {
		log.Println(err)
		return err
	}

	sr := repository.NewSheetRepository(sh.storage)
	_, err = sr.FindById(id)

	if err != nil {
		return err
	}

	su, err := sr.Update(os, id)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, su)
}

func (sh *SheetHandler) DeleteSheetHandler(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))

	if err != nil {
		return err
	}

	sr := repository.NewSheetRepository(sh.storage)
	sr.Delete(id)

	return c.NoContent(http.StatusNoContent)
}
