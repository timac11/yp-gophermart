package model

type AccrualResultStatus string

const (
	AccrualRegistered AccrualResultStatus = "REGISTERED"
	AccrualInvalid    AccrualResultStatus = "INVALID"
	AccrualProcessing AccrualResultStatus = "PROCESSING"
	AccrualProcessed  AccrualResultStatus = "PROCESSED"
)

type AccrualResult struct {
	Order   string              `json:"order"`
	Status  AccrualResultStatus `json:"status"`
	Accrual *float64            `json:"accrual"`
}
