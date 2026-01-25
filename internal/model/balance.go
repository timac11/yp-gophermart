package model

type BalanceInfo struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawModel struct {
	Id        string
	OrderId   string
	Value     int64
	CreatedAt string
	UpdatedAt string
}

type CreateWithdraw struct {
	Order float64 `json:"order"`
	Sum   float64 `json:"sum"`
}

type WithdrawHistory struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}
