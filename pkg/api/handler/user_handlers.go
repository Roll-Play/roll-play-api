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

func (userHandler *UserHandler) SignUpHandler(context echo.Context) error {
	user := new(entities.User)

	if err := context.Bind(user); err != nil {
		return context.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.INTERNAL_SERVER_ERROR, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	userRepository := repository.NewUserRepository(userHandler.storage)

	emailInUse, err := userRepository.FindByEmail(user.Email)

	if err != nil && !(err.Error() == sql.ErrNoRows.Error()) {
		return context.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.INTERNAL_SERVER_ERROR, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}
	if emailInUse != nil {
		return context.JSON(http.StatusConflict, api_error.Error{
			Error:   "e-mail already in use",
			Message: http.StatusText(http.StatusConflict),
		})
	}
	usernameInUse, err := userRepository.FindByUsername(user.Username)

	if err != nil && !(err.Error() == sql.ErrNoRows.Error()) {
		return context.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.INTERNAL_SERVER_ERROR, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	if usernameInUse != nil {
		return context.JSON(http.StatusConflict, api_error.Error{
			Error:   "username already in use",
			Message: http.StatusText(http.StatusConflict),
		})
	}
	password, err := utils.HashPassword(user.Password)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   "something went wrong",
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	user.Password = password

	createdUser, err := userRepository.Create(*user)

	if err != nil {
		return err
	}

	return context.JSON(http.StatusCreated, &UserResponse{
		Id:        createdUser.Id,
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt,
		UpdatedAt: createdUser.UpdatedAt,
	})
}

func (userHandler *UserHandler) LoginHandler(context echo.Context) error {
	requestUser := new(entities.User)

	if err := context.Bind(requestUser); err != nil {
		context.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.INTERNAL_SERVER_ERROR, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	userRepository := repository.NewUserRepository(userHandler.storage)

	plaintext := requestUser.Password

	user, err := userRepository.FindByEmail(requestUser.Email)

	if err != nil && !(err.Error() == sql.ErrNoRows.Error()) {
		return context.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.INTERNAL_SERVER_ERROR, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	if user == nil {
		return context.JSON(http.StatusNotFound, api_error.Error{
			Error:   "user not found",
			Message: http.StatusText(http.StatusNotFound),
		})
	}

	err = utils.ComparePassword(plaintext, user.Password)

	if err != nil {
		return context.JSON(http.StatusUnauthorized, api_error.Error{
			Error:   "credentials don't match",
			Message: http.StatusText(http.StatusUnauthorized),
		})
	}

	token, err := utils.CreateJWT(user.Id, 60*60*1000*24)

	if err != nil {
		return context.JSON(http.StatusInternalServerError, api_error.Error{
			Error:   fmt.Sprintf(api_error.INTERNAL_SERVER_ERROR, err.Error()),
			Message: http.StatusText(http.StatusInternalServerError),
		})
	}

	return context.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}
