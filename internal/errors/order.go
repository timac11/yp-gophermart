package errors

import (
	_errors "errors"
	"fmt"
)

type InvalidOrderNumError struct {
	OrderNum string
}

func (e *InvalidOrderNumError) Error() string {
	return fmt.Sprintf("Invalid order num: %v", e.OrderNum)
}

func IsInvalidOrderNumErr(err error) bool {
	var orderError *InvalidOrderNumError
	return _errors.As(err, &orderError)
}

func NewInvalidOrderNumErr(orderNum string) error {
	return &InvalidOrderNumError{
		OrderNum: orderNum,
	}
}
