package api

import (
	"github.com/Roll-play/roll-play-backend/pkg/api/handler"
	"github.com/Roll-play/roll-play-backend/pkg/storage"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	Server *echo.Echo
	Storage *storage.Storage
}

func NewApp(dbConnString string) (*Application, error) {
	provider := new(storage.PostgresProvider)

	storage, err := storage.NewStorage(dbConnString, provider)

	if err != nil {
		return nil, err
	}
	server := echo.New()
	
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	setRoutes(server)
	return &Application{
		Server: server,
		Storage: storage,
	}, nil
}

func setRoutes(server *echo.Echo) {
	server.GET("/healthz", handler.HealthHandler)
	server.POST("/user", handler.SignUpHandler)
}