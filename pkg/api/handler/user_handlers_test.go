package handler_test

import (
	// Import your packages here

	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Roll-play/roll-play-backend/pkg/api"
	"github.com/Roll-play/roll-play-backend/pkg/api/handler"
	"github.com/Roll-play/roll-play-backend/pkg/entities"
	"github.com/Roll-play/roll-play-backend/pkg/utils"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"golang.org/x/crypto/bcrypt"
)

type UserHandlersSuite struct {
	suite.Suite
	app *api.Application
	db  *sqlx.DB
}

func (suite *UserHandlersSuite) SetupTest() {
	// Initialize your test environment here.
	e := echo.New()
	db, err := utils.SetupTestDB("../../../.env")
	assert.NoError(suite.T(), err)

	utils.RunMigrations("file://../../../migrations")

	suite.app = &api.Application{
		Server:  e,
		Storage: db,
	}
	suite.db = db
}

func (suite *UserHandlersSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *UserHandlersSuite) TestUserHandlerSignUpSuccess() {
	requestBody := []byte(`{
				"username": "fizi",
				"email": "fizi@gmail.com",
				"password": "123123"
			}`)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := suite.app.Server.NewContext(req, rec)
	var jsonRes handler.UserResponse

	uh := handler.NewUserHandler(suite.db)
	err := uh.SignUpHandler(c)

	t := suite.T()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	json.Unmarshal(rec.Body.Bytes(), &jsonRes)

	var user entities.User
	err = suite.db.Get(&user, "SELECT id, password, email, username FROM users WHERE id=$1", jsonRes.Id)

	assert.NoError(t, err)
	assert.Equal(t, jsonRes.Username, user.Username)
	assert.Equal(t, jsonRes.Email, user.Email)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("123123"))

	assert.NoError(t, err)
}

func (suite *UserHandlersSuite) TestUserHandlerSingUpEmailInUse() {
	requestBody := []byte(`{
		"username": "fizi",
		"email": "fizi@gmail.com",
		"password": "123123"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := suite.app.Server.NewContext(req, rec)
	var jsonRes map[string]string

	uh := handler.NewUserHandler(suite.db)
	_, err := suite.db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", "fizi2", "fizi@gmail.com", "123123")

	t := suite.T()

	assert.NoError(t, err)

	err = uh.SignUpHandler(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)

	json.Unmarshal(rec.Body.Bytes(), &jsonRes)

	assert.Equal(t, "e-mail already in use", jsonRes["error"])
}

func (suite *UserHandlersSuite) TestUserHandlerSingUpUsernameInUse() {
	requestBody := []byte(`{
		"username": "fizi",
		"email": "fizi@gmail.com",
		"password": "123123"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := suite.app.Server.NewContext(req, rec)
	var jsonRes map[string]string

	uh := handler.NewUserHandler(suite.db)
	suite.db.Exec("INSERT INTO users (username, email, password) VALUES ($1, $2, $3)", "fizi", "fizi2@gmail.com", "123123")

	err := uh.SignUpHandler(c)

	t := suite.T()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusConflict, rec.Code)

	json.Unmarshal(rec.Body.Bytes(), &jsonRes)

	assert.Equal(t, "username already in use", jsonRes["error"])
}

func (suite *UserHandlersSuite) TestUserHandlerLoginSuccess() {
	requestBody := []byte(`{
		"email": "fizi@gmail.com",
		"password": "123123"
	}`)

	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := suite.app.Server.NewContext(req, rec)
	var jsonRes map[string]string

	hash, _ := utils.HashPassword("123123")
	uh := handler.NewUserHandler(suite.db)
	suite.db.Exec("INSERT INTO users (username, email, password, is_active) VALUES ($1, $2, $3, $4)", "fizi", "fizi@gmail.com", hash, true)

	err := uh.LoginHandler(c)

	t := suite.T()

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	json.Unmarshal(rec.Body.Bytes(), &jsonRes)

	_, ok := jsonRes["token"]

	assert.Equal(t, true, ok)
}

func TestUserHandlersSuite(t *testing.T) {
	suite.Run(t, new(UserHandlersSuite))
}
