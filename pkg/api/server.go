package api

import (
	"github.com/Roll-play/roll-play-backend/pkg/api/handler"
	"github.com/Roll-play/roll-play-backend/pkg/storage"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Application struct {
	Server *echo.Echo
	DB *sqlx.DB
}

func NewApp(dbConnString string) (*Application, error) {

	db, err := storage.NewDB(dbConnString, 10, 5, 5)

	if err != nil {
		return nil, err
	}
	server := echo.New()

	setRoutes(server)
	return &Application{
		Server: server,
		DB: db,
	}, nil
}

func setRoutes(server *echo.Echo) {
	server.GET("/healthz", handler.HealthHandler)
}