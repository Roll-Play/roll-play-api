package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Roll-play/roll-play-backend/pkg/entities"
	api_error "github.com/Roll-play/roll-play-backend/pkg/errors"
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
		return c.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.InternalServerErrorMessage, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	ur := repository.NewUserRepository(uh.storage)

	emailInUse, err := ur.FindByEmail(u.Email)

	if err != nil && !(err.Error() == sql.ErrNoRows.Error()) {
		return c.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.InternalServerErrorMessage, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}
	if emailInUse != nil {
		return c.JSON(http.StatusConflict, api_error.Error{
			Error:   "e-mail already in use",
			Message: http.StatusText(http.StatusConflict),
		})
	}
	usernameInUse, err := ur.FindByUsername(u.Username)

	if err != nil && !(err.Error() == sql.ErrNoRows.Error()) {
		return c.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.InternalServerErrorMessage, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	if usernameInUse != nil {
		return c.JSON(http.StatusConflict, api_error.Error{
			Error:   "username already in use",
			Message: http.StatusText(http.StatusConflict),
		})
	}
	password, err := utils.HashPassword(u.Password)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   "something went wrong",
			Message: http.StatusText(http.StatusInternalServerError),
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
		c.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.InternalServerErrorMessage, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	ur := repository.NewUserRepository(uh.storage)

	plaintext := r.Password

	u, err := ur.FindByEmail(r.Email)

	if err != nil && !(err.Error() == sql.ErrNoRows.Error()) {
		return c.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.InternalServerErrorMessage, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	if u == nil {
		return c.JSON(http.StatusNotFound, api_error.Error{
			Error:   "user not found",
			Message: http.StatusText(http.StatusNotFound),
		})
	}

	err = utils.ComparePassword(plaintext, u.Password)

	if err != nil {
		return c.JSON(http.StatusUnauthorized, api_error.Error{
			Error:   "credentials don't match",
			Message: http.StatusText(http.StatusUnauthorized),
		})
	}

	token, err := utils.CreateJWT(u.Id, 60*60*1000*24)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.InternalServerErrorMessage, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
