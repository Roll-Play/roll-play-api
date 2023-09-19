package errors

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

const (
	NotFound        = "Not found with %s: %s"
	DtoError        = "Error with Dto"
	SavingError     = "Error saving %s with values %s"
	ParseError      = "Error parsing %s"
	QueryParamError = "Query param '%s' error"
	DbError         = "Error with %s"
)
