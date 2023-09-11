package api

import (
	"time"

	"github.com/Roll-play/roll-play-backend/pkg/api/handler"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Application struct {
	Server  *echo.Echo
	Storage *sqlx.DB
}

func NewApp(dbConnString string) (*Application, error) {
	storage, err := newDB(dbConnString)

	if err != nil {
		return nil, err
	}
	server := echo.New()

	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	setRoutes(server, storage)
	return &Application{
		Server:  server,
		Storage: storage,
	}, nil
}

func setRoutes(server *echo.Echo, storage *sqlx.DB) {
	server.GET("/healthz", handler.HealthHandler)
	setUserRoutes(server, storage)
	setSheetRoutes(server, storage)
}

func setUserRoutes(server *echo.Echo, storage *sqlx.DB) {
	uh := handler.NewUserHandler(storage)
	server.POST("/user", uh.SignUpHandler)
}

func setSheetRoutes(server *echo.Echo, storage *sqlx.DB) {
	sh := handler.NewSheetHandler(storage)
	server.POST("/sheet", sh.CreateSheetHandler)
	server.GET("/sheet", sh.GetSheetListHandler)
	server.GET("/sheet/:id", sh.GetSheetHandler)
	server.PATCH("/sheet/:id", sh.UpdateSheetHandler)
	server.DELETE("/sheet/:id", sh.DeleteSheetHandler)
}

func newDB(connString string) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", connString)

	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Minute * time.Duration(5))

	return db, err
}
