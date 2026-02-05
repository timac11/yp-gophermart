package util

import (
	"github.com/timac11/yp-gophermart/internal/model"
)

func MapAccrualStatusToOrderStatus(status model.AccrualStatus) model.OrderStatus {
	if status == model.AccrualInvalid {
		return model.Invalid
	}
	if status == model.AccrualProcessing {
		return model.Processing
	}
	if status == model.AccrualProcessed {
		return model.Processed
	}

	return model.Processing
}
