package errors

import "fmt"

type InvalidPasswordError struct {
	Password string
}

func (e *InvalidPasswordError) Error() string {
	return fmt.Sprintf("Invalid password: %v", e.Password)
}

func NewInvalidPasswordError(password string) error {
	return &InvalidPasswordError{
		Password: password,
	}
}
