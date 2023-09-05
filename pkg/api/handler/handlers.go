package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HealthResponse struct {
	Alive bool `json:"alive"`
}

func HealthHandler(c echo.Context) error {
	r := HealthResponse{Alive: true}
	if err := c.Bind(r); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, r)
}
