package errors

import "fmt"

const (
	EntityNotFound      = "EntityNotFound"
	EntityAlreadyExists = "EntityAlreadyExists"
)

type EntityError struct {
	Inner  error
	Type   string
	Params any
}

func (e *EntityError) Error() string {
	return fmt.Sprintf("Error caused due to %v; Error type: %v; additional context: %v", e.Inner, e.Type, e.Params)
}

func (e *EntityError) Unwrap() error {
	return e.Inner
}

func NewEntityError(err error, errorType string, params any) error {
	return &EntityError{
		Inner:  err,
		Type:   errorType,
		Params: params,
	}
}
