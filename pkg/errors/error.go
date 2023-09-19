package api_error

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

const (
	NOT_FOUND             = "Not found with %s: %s"
	DTO_ERROR             = "Error with Dto"
	SAVING_ERROR          = "Error saving %s with values %s"
	PARSE_ERROR           = "Error parsing %s"
	QUERY_PARAM_ERROR     = "Query param '%s' error"
	DB_ERROR              = "Error with %s"
	INTERNAL_SERVER_ERROR = "Something went wrong: %v"
)

func CustomError(context echo.Context, https int, message string, args ...any) error {
	newError := Error{
		Error:   http.StatusText(https),
		Message: fmt.Sprintf(message, args...),
	}
	return context.JSON(https, newError)
}
