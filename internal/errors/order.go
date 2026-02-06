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

type OrderAlreadyProcessedError struct {
	OrderNum string
}

func (e *OrderAlreadyProcessedError) Error() string {
	return fmt.Sprintf("Already processed order num: %v", e.OrderNum)
}

func IsOrderAlreadyProcessedErr(err error) bool {
	var orderError *OrderAlreadyProcessedError
	return _errors.As(err, &orderError)
}

func OrderAlreadyProcessedErr(orderNum string) error {
	return &InvalidOrderNumError{
		OrderNum: orderNum,
	}
}
