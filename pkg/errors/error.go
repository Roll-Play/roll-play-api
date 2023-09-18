package api_error

const InternalServerErrorMessage = "something went wrong: %v"

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
