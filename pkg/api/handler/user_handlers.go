package handler

import (
	"database/sql"
	"fmt"
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
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	ur := repository.NewUserRepository(uh.storage)

	emailInUse, err := ur.FindByEmail(u.Email)

	if err != nil && !(err.Error() == sql.ErrNoRows.Error()) {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("something went wrong: %v", err),
		})
	}
	if emailInUse != nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "e-mail already in use",
		})
	}
	usernameInUse, err := ur.FindByUsername(u.Username)

	if err != nil && !(err.Error() == sql.ErrNoRows.Error()) {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("something went wrong: %v", err),
		})
	}

	if usernameInUse != nil {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "username already in use",
		})
	}
	password, err := utils.HashPassword(u.Password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "something went wrong",
		})
	}

	u.Password = password

	createdUser, err := ur.Create(*u)

	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, &UserResponse{
		Id:        createdUser.Id,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	})
}

func (uh *UserHandler) LoginHandler(c echo.Context) error {
	r := new(entities.User)

	if err := c.Bind(r); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	ur := repository.NewUserRepository(uh.storage)

	plaintext := r.Password

	u, err := ur.FindByEmail(r.Email)

	if err != nil && !(err.Error() == sql.ErrNoRows.Error()) {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("something went wrong: %v", err),
		})
	}

	if u == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "user not found",
		})
	}

	err = utils.ComparePassword(plaintext, u.Password)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "credentials don't match",
		})
	}

	token, err := utils.CreateJWT(u.Id, 60*60*1000*24)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
