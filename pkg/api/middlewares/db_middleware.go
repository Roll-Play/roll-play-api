package app_middleware

import (
	"github.com/Roll-play/roll-play-backend/pkg/storage"
	"github.com/labstack/echo/v4"
)

func DBMiddleware(storage *storage.Storage) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
					c.Set("db", storage.Provider)
					return next(c)
			}
	}
}