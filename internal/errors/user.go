package errors

import (
	_errors "errors"
	"fmt"
)

type InvalidPasswordError struct {
	Password string
}

func (e *InvalidPasswordError) Error() string {
	return fmt.Sprintf("Invalid password: %v", e.Password)
}

func IsInvalidPasswordError(err error) bool {
	var passError *InvalidPasswordError
	return _errors.As(err, &passError)
}

func NewInvalidPasswordError(password string) error {
	return &InvalidPasswordError{
		Password: password,
	}
}
