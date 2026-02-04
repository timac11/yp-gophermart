package errors

import (
	"fmt"

	"github.com/timac11/yp-gophermart/internal/model"
)

type TooManyRequestsError struct {
	RetryAfter int
}

func (e *TooManyRequestsError) Error() string {
	return fmt.Sprintf("Retry after: %v", e.RetryAfter)
}

func NewTooManyRequestsError(retryAfter int) error {
	return &TooManyRequestsError{
		RetryAfter: retryAfter,
	}
}

type InvalidAccrualStatusError struct {
	Order  string
	Status model.AccrualResultStatus
}

func (e *InvalidAccrualStatusError) Error() string {
	return fmt.Sprintf("Invalid order: %v status: %v", e.Order, e.Status)
}

func NewInvalidAccrualStatusError(order string, status model.AccrualResultStatus) error {
	return &InvalidAccrualStatusError{
		Order:  order,
		Status: status,
	}
}

type AccrualNotRegisteredError struct {
	Order string
}

func (e *AccrualNotRegisteredError) Error() string {
	return fmt.Sprintf("Accrual not registered for order: %v", e.Order)
}

func NewAccrualNotRegisteredError(order string) error {
	return &AccrualNotRegisteredError{
		Order: order,
	}
}
