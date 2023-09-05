package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Roll-play/roll-play-backend/pkg/api/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type HandlersSuite struct {
	suite.Suite
	server *echo.Echo
}

func (suite *HandlersSuite) SetupTest() {
	e := echo.New()
	suite.server = e
}

func (suite *HandlersSuite) TestHealthHanlderHealthy() {

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := suite.server.NewContext(req, rec)
	var jsonRes handler.HealthResponse

	assert.NoError(suite.T(), handler.HealthHandler(c))
	json.Unmarshal(rec.Body.Bytes(), &jsonRes)
	assert.Equal(suite.T(), http.StatusOK, rec.Code)
	assert.Equal(suite.T(), handler.HealthResponse{
		Alive: true,
	}, jsonRes)
}

func TestHandlers(t *testing.T) {
	suite.Run(t, new(HandlersSuite))
}
