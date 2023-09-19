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
	NotFound                   = "Not found with %s: %s"
	DtoError                   = "Error with Dto"
	SavingError                = "Error saving %s with values %s"
	ParseError                 = "Error parsing %s"
	QueryParamError            = "Query param '%s' error"
	DbError                    = "Error with %s"
	InternalServerErrorMessage = "something went wrong: %v"
)

func CustomError(c echo.Context, https int, message string, args ...any) error {
	errs := Error{
		Error:   http.StatusText(https),
		Message: fmt.Sprintf(message, args...),
	}
	return c.JSON(https, errs)
}
