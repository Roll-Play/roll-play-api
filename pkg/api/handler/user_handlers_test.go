package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Roll-play/roll-play-backend/pkg/api"
	"github.com/Roll-play/roll-play-backend/pkg/api/handler"
	"github.com/Roll-play/roll-play-backend/pkg/utils"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestUserHandlers(t *testing.T) {
	t.Run("registers a user", func(t *testing.T) {
		e := echo.New()

		db, err := utils.SetupTestDB("../../../.env")

		assert.NoError(t, err)

		err = db.Ping()

		assert.NoError(t, err)

		schema := `
				CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

				CREATE TABLE IF NOT EXISTS users (
						id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
						username VARCHAR(50) NOT NULL,
						email VARCHAR(255) NOT NULL,
						password VARCHAR(255) NOT NULL,
						is_active boolean DEFAULT FALSE,
						created_at TIMESTAMP DEFAULT now(),
						updated_at TIMESTAMP DEFAULT now(),
						deleted_at TIMESTAMP
				);
		`

		err = utils.ExecSchema(db, schema)

		assert.NoError(t, err)

		app := api.Application{
			Server:  e,
			Storage: db,
		}

		requestBody := []byte(`{
				"username": "fizi",
				"email": "fizi@gmail.com",
				"password": "123123"
			}`)

		req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(requestBody))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()

		c := app.Server.NewContext(req, rec)

		var jsonRes handler.UserResponse

		uh := handler.NewUserHandler(db)

		assert.NoError(t, uh.SignUpHandler(c))
		json.Unmarshal(rec.Body.Bytes(), &jsonRes)

		assert.Equal(t, "fizi", jsonRes.Username)
	})
}
