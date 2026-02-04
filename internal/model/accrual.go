package model

type AccrualStatus string

const (
	AccrualRegistered AccrualStatus = "REGISTERED"
	AccrualInvalid    AccrualStatus = "INVALID"
	AccrualProcessing AccrualStatus = "PROCESSING"
	AccrualProcessed  AccrualStatus = "PROCESSED"
)

type Accrual struct {
	Order   string        `json:"order"`
	Status  AccrualStatus `json:"status"`
	Accrual *float64      `json:"accrual"`
}
