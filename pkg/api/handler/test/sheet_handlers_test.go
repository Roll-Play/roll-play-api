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
	requestBody := []byte(`{
				"name": "Test Name",
				"description": "Not a lengthy description",
				"properties": "This should look like a json",
				"background": "Not a lenghty background"
			}`)

	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.app.Server.NewContext(req, rec)

	var jsonRes entities.Sheet

	sh := handler.NewSheetHandler(suite.db)
	err := sh.PostSheetHandler(c)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), http.StatusCreated, rec.Code)

	json.Unmarshal(rec.Body.Bytes(), &jsonRes)

	var sheet entities.Sheet
	suite.db.Get(&sheet, "SELECT name, description, properties, background FROM sheets WHERE id=$1", jsonRes.Id)

	assert.Equal(suite.T(), jsonRes.Name, sheet.Name)
	assert.Equal(suite.T(), jsonRes.Description, sheet.Description)
	assert.Equal(suite.T(), jsonRes.Properties, sheet.Properties)
	assert.Equal(suite.T(), jsonRes.Background, sheet.Background)
}

func (suite *SheetHandlersSuite) TestGetSheetHandlerSuccess() {
	sheet := entities.Sheet{
		Name:        "Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	err := suite.db.Get(&sheet, `INSERT INTO sheets (name, description, properties, background) VALUES ($1, $2, $3, $4) 
								RETURNING id, name, description, properties, background`,
		sheet.Name, sheet.Description, sheet.Properties, sheet.Background)

	assert.NoError(suite.T(), err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.app.Server.NewContext(req, rec)
	c.SetPath("/sheet/:id")
	c.SetParamNames("id")
	c.SetParamValues(sheet.Id.String())

	var jsonRes entities.Sheet

	sh := handler.NewSheetHandler(suite.db)
	errg := sh.GetSheetHandler(c)

	assert.NoError(suite.T(), errg)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	json.Unmarshal(rec.Body.Bytes(), &jsonRes)

	assert.Equal(suite.T(), jsonRes.Name, sheet.Name)
	assert.Equal(suite.T(), jsonRes.Description, sheet.Description)
	assert.Equal(suite.T(), jsonRes.Properties, sheet.Properties)
	assert.Equal(suite.T(), jsonRes.Background, sheet.Background)
}

func (suite *SheetHandlersSuite) TestGetSheetListHandlerSuccess() {
	sheet := entities.Sheet{
		Name:        "A Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	sheet2 := entities.Sheet{
		Name:        "B Other Name",
		Description: "Not a description",
		Properties:  "This should look like a json",
		Background:  "Not a background",
	}
	err := suite.db.Get(&sheet, `INSERT INTO sheets (name, description, properties, background) VALUES ($1, $2, $3, $4) 
								RETURNING id, name, description, properties, background`,
		sheet.Name, sheet.Description, sheet.Properties, sheet.Background)
	err2 := suite.db.Get(&sheet2, `INSERT INTO sheets (name, description, properties, background) VALUES ($1, $2, $3, $4) 
								RETURNING id, name, description, properties, background`,
		sheet2.Name, sheet2.Description, sheet2.Properties, sheet2.Background)

	assert.NoError(suite.T(), err)
	assert.NoError(suite.T(), err2)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	urlq := req.URL.Query()
	urlq.Add("page", "0")
	urlq.Add("size", "2")
	req.URL.RawQuery = urlq.Encode()
	rec := httptest.NewRecorder()
	c := suite.app.Server.NewContext(req, rec)

	var jsonRes []entities.Sheet

	sh := handler.NewSheetHandler(suite.db)
	errg := sh.GetSheetListHandler(c)

	assert.NoError(suite.T(), errg)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	json.Unmarshal(rec.Body.Bytes(), &jsonRes)

	assert.Equal(suite.T(), jsonRes[0].Name, sheet.Name)
	assert.Equal(suite.T(), jsonRes[0].Description, sheet.Description)
	assert.Equal(suite.T(), jsonRes[0].Properties, sheet.Properties)
	assert.Equal(suite.T(), jsonRes[0].Background, sheet.Background)

	assert.Equal(suite.T(), jsonRes[1].Name, sheet2.Name)
	assert.Equal(suite.T(), jsonRes[1].Description, sheet2.Description)
	assert.Equal(suite.T(), jsonRes[1].Properties, sheet2.Properties)
	assert.Equal(suite.T(), jsonRes[1].Background, sheet2.Background)
}

func (suite *SheetHandlersSuite) TestPatchSheetHandlerSuccess() {
	sheet := entities.Sheet{
		Name:        "Test Name",
		Description: "Not a lengthy description",
		Properties:  "This should look like a json",
		Background:  "Not a lenghty background",
	}
	err := suite.db.Get(&sheet, `INSERT INTO sheets (name, description, properties, background) VALUES ($1, $2, $3, $4) 
								RETURNING id`,
		sheet.Name, sheet.Description, sheet.Properties, sheet.Background)

	assert.NoError(suite.T(), err)

	us := entities.Sheet{
		Name:        "New Name",
		Description: "New description",
		Properties:  "New json",
	}

	requestBody, errm := json.Marshal(us)
	assert.NoError(suite.T(), errm)

	req := httptest.NewRequest(http.MethodPatch, "/", bytes.NewBuffer(requestBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.app.Server.NewContext(req, rec)
	c.SetPath("/sheet/:id")
	c.SetParamNames("id")
	c.SetParamValues(sheet.Id.String())

	var jsonRes entities.Sheet

	sh := handler.NewSheetHandler(suite.db)
	errg := sh.PatchSheetHandler(c)

	json.Unmarshal(rec.Body.Bytes(), &jsonRes)

	assert.NoError(suite.T(), errg)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)

	json.Unmarshal(rec.Body.Bytes(), &jsonRes)

	assert.Equal(suite.T(), jsonRes.Name, us.Name)
	assert.Equal(suite.T(), jsonRes.Description, us.Description)
	assert.Equal(suite.T(), jsonRes.Properties, us.Properties)
	assert.Equal(suite.T(), jsonRes.Background, sheet.Background)

	suite.db.Get(&jsonRes, "SELECT name, description, properties, background FROM sheets WHERE id=$1", jsonRes.Id)

	assert.Equal(suite.T(), jsonRes.Name, us.Name)
	assert.Equal(suite.T(), jsonRes.Description, us.Description)
	assert.Equal(suite.T(), jsonRes.Properties, us.Properties)
	assert.Equal(suite.T(), jsonRes.Background, sheet.Background)
}

func TestSheetHandlersSuite(t *testing.T) {
	suite.Run(t, new(SheetHandlersSuite))
}
