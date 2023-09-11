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
	db, err := utils.SetupTestDB("../../../../.env")
	assert.NoError(suite.T(), err)

	utils.RunMigrations("file://../../../../migrations")

	suite.app = &api.Application{
		Server:  e,
		Storage: db,
	}
	suite.db = db
}

func (suite *UserHandlersSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *UserHandlersSuite) TestUserHandlerSuccess() {
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

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	json.Unmarshal(rec.Body.Bytes(), &jsonRes)

	var user entities.User
	suite.db.Get(&user, "SELECT id, password, email, username FROM users WHERE id=$1", jsonRes.Id)

	assert.Equal(suite.T(), jsonRes.Username, user.Username)
	assert.Equal(suite.T(), jsonRes.Email, user.Email)

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("123123"))
	assert.NoError(suite.T(), err)
}

func TestUserHandlersSuite(t *testing.T) {
	suite.Run(t, new(UserHandlersSuite))
}
