package errors

import (
	_errors "errors"
)

type InsufficientBalance struct{}

func (e *InsufficientBalance) Error() string {
	return "Insufficient balance error"
}

func IsInsufficientBalanceErr(err error) bool {
	var orderError *InsufficientBalance
	return _errors.As(err, &orderError)
}

func NewInsufficientBalance() error {
	return &InsufficientBalance{}
}
