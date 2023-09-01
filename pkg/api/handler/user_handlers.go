package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type User struct {
	Name string `json:"name"`
}

func SignUpHandler(c echo.Context) error {
	u := new(User)

	if err := c.Bind(u); err != nil {
		return err
	}

	
	return c.JSON(http.StatusOK, u)
}