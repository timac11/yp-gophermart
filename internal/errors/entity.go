package errors

import (
	_errors "errors"
	"fmt"
)

const (
	EntityNotFound      EntityType = "EntityNotFound"
	EntityAlreadyExists EntityType = "EntityAlreadyExists"
)

type EntityType string

type EntityError struct {
	Inner  error
	Type   EntityType
	Params any
}

func (e *EntityError) Error() string {
	return fmt.Sprintf("Error caused due to %v; Error type: %v; additional context: %v", e.Inner, e.Type, e.Params)
}

func (e *EntityError) Unwrap() error {
	return e.Inner
}

func IsEntityNotFoundErr(err error) bool {
	var entityError *EntityError
	return _errors.As(err, &entityError) && entityError.Type == EntityNotFound
}

func IsEntityAlreadyExistsErr(err error) bool {
	var entityError *EntityError
	return _errors.As(err, &entityError) && entityError.Type == EntityAlreadyExists
}

func NewEntityError(err error, errorType EntityType, params any) error {
	return &EntityError{
		Inner:  err,
		Type:   errorType,
		Params: params,
	}
}
