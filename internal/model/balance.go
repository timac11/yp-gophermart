package model

import "time"

type BalanceInfo struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawModel struct {
	ID      string
	OrderID string
	Value   int64
}

type CreateWithdraw struct {
	Order string  `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawHistory struct {
	Order       string     `json:"order"`
	Sum         *float64   `json:"sum"`
	ProcessedAt *time.Time `json:"processed_at"`
}
