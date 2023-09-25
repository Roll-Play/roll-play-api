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
	contex := suite.app.Server.NewContext(request, record)
	contex.Set("user", savedId)

	var jsonRes entities.Sheet

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.CreateSheetHandler(contex)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	var sheet entities.Sheet
	suite.db.Get(&sheet, "SELECT name, description, properties, background, user_id FROM sheets WHERE id=$1", jsonRes.Id)

	assert.Equal(t, jsonRes.Name, sheet.Name)
	assert.Equal(t, jsonRes.Description, sheet.Description)
	assert.Equal(t, jsonRes.Properties, sheet.Properties)
	assert.Equal(t, jsonRes.Background, sheet.Background)
	assert.Equal(t, savedId, sheet.UserId)
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
	err = suite.db.Get(&sheet, `INSERT INTO sheets (name, description, properties, background, user_id) VALUES ($1, $2, $3, $4, $5) 
								RETURNING id, name, description, properties, background, user_id`,
		sheet.Name, sheet.Description, sheet.Properties, sheet.Background, savedId)

	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	contex := suite.app.Server.NewContext(request, record)
	contex.SetPath("/sheet/:id")
	contex.SetParamNames("id")
	contex.SetParamValues(sheet.Id.String())
	contex.Set("user", savedId)

	var jsonRes entities.Sheet

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.GetSheetHandler(contex)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes.Name, sheet.Name)
	assert.Equal(t, jsonRes.Description, sheet.Description)
	assert.Equal(t, jsonRes.Properties, sheet.Properties)
	assert.Equal(t, jsonRes.Background, sheet.Background)
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
		UserId:      savedId,
	}
	err = createSheet(suite.db, &sheet)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	contex := suite.app.Server.NewContext(request, record)
	contex.SetPath("/sheet/:id")
	contex.SetParamNames("id")
	contex.SetParamValues(sheet.Id.String())

	randId, err := uuid.NewRandom()
	assert.NoError(t, err)
	contex.Set("user", randId)

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.GetSheetHandler(contex)
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
		UserId:      savedId,
	}
	err = createSheet(suite.db, &sheet)
	assert.NoError(t, err)

	savedId2, err := setupUser(suite.db, "test2", t)
	assert.NoError(t, err)
	sheet2 := entities.Sheet{
		Name:        "B Other Name",
		Description: "Not a description",
		Properties:  "This should look like a json",
		Background:  "Not a background",
		UserId:      savedId2,
	}
	err = createSheet(suite.db, &sheet2)
	assert.NoError(t, err)

	sheet3 := entities.Sheet{
		Name:        "contex Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
		UserId:      savedId,
	}
	err = createSheet(suite.db, &sheet3)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	urlq := request.URL.Query()
	urlq.Add("page", "0")
	urlq.Add("size", "2")
	request.URL.RawQuery = urlq.Encode()
	record := httptest.NewRecorder()
	contex := suite.app.Server.NewContext(request, record)
	contex.Set("user", savedId)

	var jsonRes []entities.Sheet

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.GetSheetListHandler(contex)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, record.Code)

	json.Unmarshal(record.Body.Bytes(), &jsonRes)

	assert.Equal(t, jsonRes[0].Name, sheet.Name)
	assert.Equal(t, jsonRes[0].Description, sheet.Description)
	assert.Equal(t, jsonRes[0].Properties, sheet.Properties)
	assert.Equal(t, jsonRes[0].Background, sheet.Background)

	assert.Equal(t, jsonRes[1].Name, sheet3.Name)
	assert.Equal(t, jsonRes[1].Description, sheet3.Description)
	assert.Equal(t, jsonRes[1].Properties, sheet3.Properties)
	assert.Equal(t, jsonRes[1].Background, sheet3.Background)
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
		UserId:      savedId,
	}
	err = createSheet(suite.db, &sheet)
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
	contex := suite.app.Server.NewContext(request, record)
	contex.SetPath("/sheet/:id")
	contex.SetParamNames("id")
	contex.SetParamValues(sheet.Id.String())
	contex.Set("user", savedId)

	var jsonRes entities.Sheet

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.UpdateSheetHandler(contex)

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

func (suite *SheetHandlersSuite) TestDeleteSheetHandlerSuccess() {
	t := suite.T()

	savedId, err := setupUser(suite.db, "test", t)
	assert.NoError(t, err)

	sheet := entities.Sheet{
		Name:        "Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
		UserId:      savedId,
	}
	err = createSheet(suite.db, &sheet)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()
	contex := suite.app.Server.NewContext(request, record)
	contex.SetPath("/sheet/:id")
	contex.SetParamNames("id")
	contex.SetParamValues(sheet.Id.String())
	contex.Set("user", savedId)

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.DeleteSheetHandler(contex)

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

	contex := suite.app.Server.NewContext(request, record)
	contex.SetPath("/sheet/:id")
	contex.SetParamNames("id")
	contex.SetParamValues(ruuid.String())
	contex.Set("user", savedId)

	var jsonRes api_error.Error

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.GetSheetHandler(contex)
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

	contex := suite.app.Server.NewContext(request, record)
	contex.SetPath("/sheet/:id")
	contex.SetParamNames("id")
	contex.SetParamValues(ruuid.String())
	contex.Set("user", savedId)

	var jsonRes api_error.Error

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.DeleteSheetHandler(contex)
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
		UserId:      savedId,
	}
	err = createSheet(suite.db, &sheet)
	assert.NoError(t, err)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	record := httptest.NewRecorder()

	ruuid, err := uuid.NewRandom()
	assert.NoError(t, err)

	contex := suite.app.Server.NewContext(request, record)
	contex.SetPath("/sheet/:id")
	contex.SetParamNames("id")
	contex.SetParamValues(ruuid.String())
	contex.Set("user", ruuid)

	var jsonRes api_error.Error

	sheetHandler := handler.NewSheetHandler(suite.db)
	err = sheetHandler.DeleteSheetHandler(contex)
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

func createSheet(db *sqlx.DB, sheet *entities.Sheet) error {
	err := db.Get(sheet, `INSERT INTO sheets (name, description, properties, background, user_id) VALUES ($1, $2, $3, $4, $5) 
	RETURNING id, name, description, properties, background, user_id`,
		sheet.Name, sheet.Description, sheet.Properties, sheet.Background, sheet.UserId)
	return err
}
