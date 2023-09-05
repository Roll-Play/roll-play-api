package handler

import (
	"net/http"
	"time"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	repository "github.com/Roll-play/roll-play-backend/pkg/repositories"
	"github.com/Roll-play/roll-play-backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	storage *sqlx.DB
}

type UserResponse struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewUserHandler(storage *sqlx.DB) *UserHandler {
	return &UserHandler{
		storage,
	}
}

func (uh *UserHandler) SignUpHandler(c echo.Context) error {
	u := new(entities.User)

	if err := c.Bind(u); err != nil {
		return err
	}

	password, err := utils.HashPassword(u.Password)

	if err != nil {
		return err
	}

	u.Password = password
	ur := repository.NewUserRepository(u)

	err = ur.Create(uh.storage)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, &UserResponse{
		Id:        u.Id,
		Username:  u.Username,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	})
}
