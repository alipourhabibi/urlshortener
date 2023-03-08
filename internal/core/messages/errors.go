package messages

import "fmt"

var (
	ErrNotFound            = fmt.Errorf("Not Found")
	ErrInternalServerError = fmt.Errorf("Internal Server Error")
	ErrUnauthorized        = fmt.Errorf("Unauthorized")
	ErrAlreadyExists       = fmt.Errorf("Already Exists")
)
