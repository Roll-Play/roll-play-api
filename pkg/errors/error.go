package errors

type Error struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}
