package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Roll-play/roll-play-backend/pkg/api/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)



func TestHandlers(t *testing.T) {
	t.Run("successfully requests /healthz endpoint", func(t *testing.T) {
		e := echo.New()

		req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := e.NewContext(req, rec)
		var jsonRes handler.HealthResponse
		
		if assert.NoError(t, handler.HealthHandler(c)) {
			json.Unmarshal(rec.Body.Bytes(), &jsonRes)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.Equal(t, handler.HealthResponse{
				Alive: true,
			}, jsonRes)
		}
	})
}