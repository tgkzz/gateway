package handler

const (
	UsernameIsEmpty = "username is empty"
	PasswordIsEmpty = "password is empty"
)

type ErrorHandler struct {
	Op   string `json:"op"`
	Code string `json:"code"`
	Err  string `json:"err,omitempty"`
}
