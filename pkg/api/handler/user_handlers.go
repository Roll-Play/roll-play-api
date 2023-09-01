package handler

import (
	"net/http"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	repository "github.com/Roll-play/roll-play-backend/pkg/repositories"
	"github.com/Roll-play/roll-play-backend/pkg/storage"
	"github.com/labstack/echo/v4"
)

type User struct {
	Username string `json:"username"`
}

func SignUpHandler(c echo.Context) error {
	u := new(entities.User)

	storage := *c.Get("db").(*storage.Provider)

	if err := c.Bind(u); err != nil {
		return err
	}

	ur := repository.NewUserRepository(u)	
	
	err := storage.Create(ur)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, &User{
		Username: u.Username,
	})
}