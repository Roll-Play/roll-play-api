package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Roll-play/roll-play-backend/pkg/api"
	"github.com/Roll-play/roll-play-backend/pkg/api/handler"
	"github.com/Roll-play/roll-play-backend/pkg/entities"
	api_error "github.com/Roll-play/roll-play-backend/pkg/errors"
	"github.com/Roll-play/roll-play-backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SheetUserHandlersSuite struct {
	suite.Suite
	app *api.Application
	db  *sqlx.DB
}

func (suite *SheetUserHandlersSuite) SetupTest() {
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

func (suite *SheetUserHandlersSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *SheetHandlersSuite) TestPostSheetUserHandlerSuccess() {
	t := suite.T()
	savedId, err := utils.SetupUser(suite.db, "test", t)
	assert.NoError(t, err)

	savedId2, err := utils.SetupUser(suite.db, "test2", t)
	assert.NoError(t, err)

	sheet := entities.Sheet{
		Name:        "A Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	utils.CreateSheetAndSheetUser(suite.db, &sheet, savedId, entities.WRITE, true)
	assert.NoError(t, err)

	us := handler.PermissionList{
		Permission: 0,
		UserIds:    []uuid.UUID{savedId2},
	}

	requestBody, errm := json.Marshal(us)
	assert.NoError(t, errm)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.Set("user", savedId)
	context.SetPath("/sheet/:id/permission")
	context.SetParamNames("id")
	context.SetParamValues(sheet.Id.String())

	sheetUserHandler := handler.NewSheetUserHandler(suite.db)
	err = sheetUserHandler.UpdatePermissions(context)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, record.Code)

	var sheetUser entities.SheetUser
	suite.db.Get(&sheetUser, "SELECT * FROM sheet_user WHERE sheet_id=$1 AND user_id=$2", sheet.Id, savedId2)

	assert.Equal(t, sheetUser.SheetId, sheet.Id)
	assert.Equal(t, sheetUser.UserId, savedId2)
	assert.Equal(t, sheetUser.Owner, false)
	assert.Equal(t, sheetUser.Permission, entities.WRITE)
}

func (suite *SheetHandlersSuite) TestPostSheetUserHandlerFailWithNotOwner() {
	t := suite.T()
	savedId, err := utils.SetupUser(suite.db, "test", t)
	assert.NoError(t, err)

	savedId2, err := utils.SetupUser(suite.db, "test2", t)
	assert.NoError(t, err)

	sheet := entities.Sheet{
		Name:        "A Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	utils.CreateSheetAndSheetUser(suite.db, &sheet, savedId, entities.WRITE, false)
	assert.NoError(t, err)

	us := handler.PermissionList{
		Permission: 0,
		UserIds:    []uuid.UUID{savedId2},
	}

	requestBody, errm := json.Marshal(us)
	assert.NoError(t, errm)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.Set("user", savedId)
	context.SetPath("/sheet/:id/permission")
	context.SetParamNames("id")
	context.SetParamValues(sheet.Id.String())

	sheetUserHandler := handler.NewSheetUserHandler(suite.db)
	err = sheetUserHandler.UpdatePermissions(context)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, record.Code)

	var jsonRes api_error.Error
	json.Unmarshal(record.Body.Bytes(), &jsonRes)
	assert.Equal(t, jsonRes.Error, http.StatusText(record.Code))
	assert.Equal(t, jsonRes.Message, fmt.Sprintf(api_error.PERMISSION_ERROR, savedId))
}

func TestSheetUserHandlersSuite(t *testing.T) {
	suite.Run(t, new(SheetHandlersSuite))
}
