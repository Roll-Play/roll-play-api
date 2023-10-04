package handler_test

import (
	// Import your packages here

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
	repository "github.com/Roll-play/roll-play-backend/pkg/repositories"
	"github.com/Roll-play/roll-play-backend/pkg/utils"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SheetHandlersSuite struct {
	suite.Suite
	app *api.Application
	db  *sqlx.DB
}

func (suite *SheetHandlersSuite) SetupTest() {
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

func (suite *SheetHandlersSuite) TearDownTest() {
	suite.db.Close()
}

func (suite *SheetHandlersSuite) TestPostSheetHandlerSuccess() {
	t := suite.T()
	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)

	requestBody := []byte(`{
				"name": "Test Name",
				"description": "Not a lengthy description",
				"properties": "This should look like a json",
				"background": "Not a lenghty background"
			}`)

	request := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.Set("user", savedId)

	var jsonRes entities.Sheet

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.CreateSheetHandler(context)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	var sheet entities.Sheet
	suite.db.Get(&sheet, "SELECT s.name, s.description, s.properties, s.background FROM sheets s WHERE s.id=$1", jsonRes.Id)

	assert.Equal(t, jsonRes.Name, sheet.Name)
	assert.Equal(t, jsonRes.Description, sheet.Description)
	assert.Equal(t, jsonRes.Properties, sheet.Properties)
	assert.Equal(t, jsonRes.Background, sheet.Background)

	var sheetUser entities.SheetUser
	suite.db.Get(&sheetUser, "SELECT * FROM sheet_user WHERE sheet_id=$1 AND user_id=$2", jsonRes.Id, savedId)

	assert.Equal(t, sheetUser.Owner, true)
	assert.Equal(t, sheetUser.Permission, entities.READ)
}

func (suite *SheetHandlersSuite) TestGetSheetHandlerSuccess() {
	t := suite.T()
	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)

	sheet := entities.Sheet{
		Name:        "Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	createSheetAndSheetUser(suite.db, &sheet, savedId, entities.READ, true)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/sheet/:id")
	context.SetParamNames("id")
	context.SetParamValues(sheet.Id.String())
	context.Set("user", savedId)

	var jsonRes entities.Sheet

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.GetSheetHandler(context)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Name, sheet.Name)
	assert.Equal(t, jsonRes.Description, sheet.Description)
	assert.Equal(t, jsonRes.Properties, sheet.Properties)
	assert.Equal(t, jsonRes.Background, sheet.Background)
	assert.Equal(t, jsonRes.Permission, entities.READ)
	assert.Equal(t, jsonRes.Owner, true)
}

func (suite *SheetHandlersSuite) TestGetSheetHandlerFailWithWrongUser() {
	t := suite.T()
	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)

	sheet := entities.Sheet{
		Name:        "Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	createSheetAndSheetUser(suite.db, &sheet, savedId, entities.READ, true)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/sheet/:id")
	context.SetParamNames("id")
	context.SetParamValues(sheet.Id.String())

	randId, err := uuid.NewRandom()
	assert.NoError(t, err)
	context.Set("user", randId)

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.GetSheetHandler(context)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusBadRequest, record.Code)

	var jsonRes api_error.Error

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Error, http.StatusText(record.Code))
	assert.Equal(t, jsonRes.Message, fmt.Sprintf(api_error.NOT_FOUND, "id", sheet.Id.String()))
}

func (suite *SheetHandlersSuite) TestGetSheetListHandlerSuccess() {
	t := suite.T()

	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)
	sheet := entities.Sheet{
		Name:        "A Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	createSheetAndSheetUser(suite.db, &sheet, savedId, entities.READ, true)
	assert.NoError(t, err)

	savedId2, err := setupUser(suite.db, "test2", t)
	assert.NoError(t, err)
	sheet2 := entities.Sheet{
		Name:        "B Other Name",
		Description: "Not a description",
		Properties:  "This should look like a json",
		Background:  "Not a background",
	}
	createSheetAndSheetUser(suite.db, &sheet2, savedId2, entities.READ, true)
	assert.NoError(t, err)

	sheet3 := entities.Sheet{
		Name:        "context Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	createSheetAndSheetUser(suite.db, &sheet3, savedId, entities.READ, true)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	urlq := request.URL.Query()
	urlq.Add("page", "0")
	urlq.Add("size", "2")
	request.URL.RawQuery = urlq.Encode()
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.Set("user", savedId)

	var jsonRes []entities.Sheet

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.GetSheetListHandler(context)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes[0].Name, sheet.Name)
	assert.Equal(t, jsonRes[0].Description, sheet.Description)
	assert.Equal(t, jsonRes[0].Properties, sheet.Properties)
	assert.Equal(t, jsonRes[0].Background, sheet.Background)
	assert.Equal(t, jsonRes[0].Permission, entities.READ)
	assert.Equal(t, jsonRes[0].Owner, true)

	assert.Equal(t, jsonRes[1].Name, sheet3.Name)
	assert.Equal(t, jsonRes[1].Description, sheet3.Description)
	assert.Equal(t, jsonRes[1].Properties, sheet3.Properties)
	assert.Equal(t, jsonRes[1].Background, sheet3.Background)
	assert.Equal(t, jsonRes[1].Permission, entities.READ)
	assert.Equal(t, jsonRes[1].Owner, true)
}

func (suite *SheetHandlersSuite) TestPatchSheetHandlerSuccess() {
	t := suite.T()

	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)

	sheet := entities.Sheet{
		Name:        "Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	createSheetAndSheetUser(suite.db, &sheet, savedId, entities.READ, true)
	assert.NoError(t, err)

	us := entities.Sheet{
		Name:        "New Name",
		Description: "New description",
		Properties:  "New json",
	}

	requestBody, errm := json.Marshal(us)
	assert.NoError(t, errm)

	request := httptest.NewRequest(http.MethodPatch, "/", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/sheet/:id")
	context.SetParamNames("id")
	context.SetParamValues(sheet.Id.String())
	context.Set("user", savedId)

	var jsonRes entities.Sheet

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.UpdateSheetHandler(context)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Name, us.Name)
	assert.Equal(t, jsonRes.Description, us.Description)
	assert.Equal(t, jsonRes.Properties, us.Properties)
	assert.Equal(t, jsonRes.Background, sheet.Background)

	suite.db.Get(&jsonRes, "SELECT name, description, properties, background FROM sheets WHERE id=$1", jsonRes.Id)

	assert.Equal(t, jsonRes.Name, us.Name)
	assert.Equal(t, jsonRes.Description, us.Description)
	assert.Equal(t, jsonRes.Properties, us.Properties)
	assert.Equal(t, jsonRes.Background, sheet.Background)
}

func (suite *SheetHandlersSuite) TestPatchSheetHandlerFailWithForbiden() {
	t := suite.T()

	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)

	sheet := entities.Sheet{
		Name:        "Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	createSheetAndSheetUser(suite.db, &sheet, savedId, entities.READ, false)
	assert.NoError(t, err)

	us := entities.Sheet{
		Name:        "New Name",
		Description: "New description",
		Properties:  "New json",
	}

	requestBody, errm := json.Marshal(us)
	assert.NoError(t, errm)

	request := httptest.NewRequest(http.MethodPatch, "/", bytes.NewBuffer(requestBody))
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/sheet/:id")
	context.SetParamNames("id")
	context.SetParamValues(sheet.Id.String())
	context.Set("user", savedId)

	var jsonRes api_error.Error

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.UpdateSheetHandler(context)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusForbidden, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Error, http.StatusText(record.Code))
	assert.Equal(t, jsonRes.Message, fmt.Sprintf(api_error.PERMISSION_ERROR, savedId))
}

func (suite *SheetHandlersSuite) TestDeleteSheetHandlerSuccess() {
	t := suite.T()

	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)

	sheet := entities.Sheet{
		Name:        "Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	createSheetAndSheetUser(suite.db, &sheet, savedId, entities.READ, true)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/sheet/:id")
	context.SetParamNames("id")
	context.SetParamValues(sheet.Id.String())
	context.Set("user", savedId)

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.DeleteSheetHandler(context)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, record.Code)

	test := new(entities.Sheet)
	errex := suite.db.Get(test, "SELECT name, description, properties, background FROM sheets WHERE id=$1", sheet.Id)

	assert.Error(t, errex)
}

func (suite *SheetHandlersSuite) TestGetSheetHandlerFail() {
	t := suite.T()

	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()

	ruuid, err := uuid.NewRandom()
	assert.NoError(t, err)

	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/sheet/:id")
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

func (suite *SheetHandlersSuite) TestDeleteSheetHandlerFail() {
	t := suite.T()

	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()

	ruuid, err := uuid.NewRandom()
	assert.NoError(t, err)

	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/sheet/:id")
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

func (suite *SheetHandlersSuite) TestDeleteSheetHandlerFailWrongUser() {
	t := suite.T()

	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)

	sheet := entities.Sheet{
		Name:        "Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	err = createSheetAndSheetUser(suite.db, &sheet, savedId, entities.READ, true)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()

	ruuid, err := uuid.NewRandom()
	assert.NoError(t, err)

	context := suite.app.Server.NewContext(request, record)
	context.SetPath("/sheet/:id")
	context.SetParamNames("id")
	context.SetParamValues(ruuid.String())
	context.Set("user", ruuid)

	var jsonRes api_error.Error

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.DeleteSheetHandler(context)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Error, http.StatusText(record.Code))
	assert.Equal(t, jsonRes.Message, fmt.Sprintf(api_error.NOT_FOUND, "id", ruuid))
}

func TestSheetHandlersSuite(t *testing.T) {
	suite.Run(t, new(SheetHandlersSuite))
}

func setupUser(db *sqlx.DB, username string, t *testing.T) (uuid.UUID, error) {
	userRepository := repository.NewUserRepository(db)
	savedUser, err := userRepository.Create(entities.User{
		Username: username,
		Email:    username + "@test",
		Password: "test",
	})

	if err != nil {
		return uuid.New(), err
	}

	return savedUser.Id, nil
}

func createSheetAndSheetUser(db *sqlx.DB, sheet *entities.Sheet, userId uuid.UUID, permission int, owner bool) error {
	err := db.Get(sheet, `INSERT INTO sheets (name, description, properties, background) VALUES ($1, $2, $3, $4)
	RETURNING id, name, description, properties, background`,
		sheet.Name, sheet.Description, sheet.Properties, sheet.Background)
	if err != nil {
		return err
	}

	_, err = db.Exec(`INSERT INTO sheet_user (sheet_id, user_id, permission, owner) VALUES ($1, $2, $3, $4)`, sheet.Id, userId, permission, owner)
	return err
}
