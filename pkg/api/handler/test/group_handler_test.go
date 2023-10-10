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

type GroupHandlersSuite struct {
	suite.Suite
	app *api.Application
	db  *sqlx.DB
}

func (suite *GroupHandlersSuite) SetupTest() {
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

func (suite *GroupHandlersSuite) TearDownTest() {
	suite.db.Close()
}

func TestGroupHandlersSuite(t *testing.T) {
	suite.Run(t, new(GroupHandlersSuite))
}

func (suite *GroupHandlersSuite) TestPostGroupHandlerSuccess() {
	t := suite.T()
	savedId, err := utils.SetupUser(suite.db, "test", t)
	assert.NoError(t, err)

	groupBody := entities.GroupDto{
		Name:   "Testing Group",
		Public: true,
	}

	requestBody, err := json.Marshal(groupBody)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.Set("user", savedId)

	var jsonRes entities.Group

	groupHandler := handler.NewGroupHandler(suite.db)
	err = groupHandler.CreateGroupHandler(context)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	var group entities.Group
	suite.db.Get(&group, "SELECT user_id, name, public FROM player_group s WHERE s.id=$1", jsonRes.Id)

	assert.Equal(t, jsonRes.Name, group.Name)
	assert.Equal(t, jsonRes.Public, group.Public)
	assert.Equal(t, jsonRes.UserId, group.UserId)
}

func (suite *GroupHandlersSuite) TestGetGroupHandlerSuccess() {
	t := suite.T()

	savedId, err := utils.SetupUser(suite.db, "test", t)
	assert.NoError(t, err)

	group := entities.Group{
		Name:   "Test Name",
		Public: true,
	}
	createGroup(suite.db, &group, savedId)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/group/:id")
	context.SetParamNames("id")
	context.SetParamValues(group.Id.String())
	context.Set("user", savedId)

	groupHandler := handler.NewGroupHandler(suite.db)
	err = groupHandler.GetGroupHandler(context)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, record.Code)

	var jsonRes entities.Group
	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Name, group.Name)
	assert.Equal(t, jsonRes.Public, group.Public)
	assert.Equal(t, jsonRes.UserId, savedId)
}

func (suite *SheetHandlersSuite) TestGetGroupHandlerFail() {
	t := suite.T()

	savedId, err := utils.SetupUser(suite.db, "test", t)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()

	ruuid, err := uuid.NewRandom()
	assert.NoError(t, err)

	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/group/:id")
	context.SetParamNames("id")
	context.SetParamValues(ruuid.String())
	context.Set("user", savedId)

	var jsonRes api_error.Error

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.GetSheetHandler(context)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Error, http.StatusText(record.Code))
	assert.Equal(t, jsonRes.Message, fmt.Sprintf(api_error.NOT_FOUND, "id", ruuid))
}

func (suite *SheetHandlersSuite) TestGetGroupHandlerFailWrongUser() {
	t := suite.T()

	savedId, err := utils.SetupUser(suite.db, "test", t)
	assert.NoError(t, err)

	group := entities.Group{
		Name:   "Test Name",
		Public: true,
	}
	createGroup(suite.db, &group, savedId)
	assert.NoError(t, err)

	randId, err := uuid.NewRandom()
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/group/:id")
	context.SetParamNames("id")
	context.SetParamValues(group.Id.String())
	context.Set("user", randId)

	var jsonRes api_error.Error

	groupHandler := handler.NewGroupHandler(suite.db)
	err = groupHandler.GetGroupHandler(context)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Error, http.StatusText(record.Code))
	assert.Equal(t, jsonRes.Message, fmt.Sprintf(api_error.NOT_FOUND, "id", randId))
}

func (suite *GroupHandlersSuite) TestDeleteGroupHandlerSuccess() {
	t := suite.T()

	savedId, err := utils.SetupUser(suite.db, "test", t)
	assert.NoError(t, err)

	group := entities.Group{
		Name:   "Test Name",
		Public: true,
	}
	createGroup(suite.db, &group, savedId)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodDelete, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/group/:id")
	context.SetParamNames("id")
	context.SetParamValues(group.Id.String())
	context.Set("user", savedId)

	groupHandler := handler.NewGroupHandler(suite.db)
	err = groupHandler.DeleteGroupHandler(context)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, record.Code)

	test := new(entities.Sheet)
	err = suite.db.Get(test, "SELECT * FROM sheets WHERE id=$1", group.Id)

	assert.Error(t, err)
}

func (suite *SheetHandlersSuite) TestDeleteGroupHandlerFail() {
	t := suite.T()

	savedId, err := utils.SetupUser(suite.db, "test", t)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodDelete, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()

	ruuid, err := uuid.NewRandom()
	assert.NoError(t, err)

	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/group/:id")
	context.SetParamNames("id")
	context.SetParamValues(ruuid.String())
	context.Set("user", savedId)

	var jsonRes api_error.Error

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.DeleteSheetHandler(context)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Error, http.StatusText(record.Code))
	assert.Equal(t, jsonRes.Message, fmt.Sprintf(api_error.NOT_FOUND, "id", ruuid))
}

func (suite *SheetHandlersSuite) TestDeleteGroupHandlerFailWrongUser() {
	t := suite.T()

	savedId, err := utils.SetupUser(suite.db, "test", t)
	assert.NoError(t, err)

	group := entities.Group{
		Name:   "Test Name",
		Public: true,
	}
	createGroup(suite.db, &group, savedId)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodDelete, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/group/:id")
	context.SetParamNames("id")
	context.SetParamValues(group.Id.String())

	randId, err := uuid.NewRandom()
	assert.NoError(t, err)
	context.Set("user", randId)

	var jsonRes api_error.Error

	groupHandler := handler.NewGroupHandler(suite.db)
	err = groupHandler.DeleteGroupHandler(context)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Error, http.StatusText(record.Code))
	assert.Equal(t, jsonRes.Message, fmt.Sprintf(api_error.NOT_FOUND, "id", randId))
}

func createGroup(db *sqlx.DB, group *entities.Group, userId uuid.UUID) error {
	err := db.Get(group, `INSERT INTO player_group (name, public, user_id) VALUES ($1, $2, $3)
	RETURNING id ,name, public`,
		group.Name, group.Public, userId)
	if err != nil {
		return err
	}
	return nil
}
