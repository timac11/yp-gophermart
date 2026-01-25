package model

type OrderStatus string

const (
	New        OrderStatus = "NEW"
	Processing OrderStatus = "PROCESSING"
	Invalid    OrderStatus = "INVALID"
	Processed  OrderStatus = "PROCESSED"
)

type OrderModel struct {
	Id        string
	OrderNum  string
	UserId    string
	CreatedAt string
	UpdatedAt string
}

type OrderInfo struct {
	Number     string   `json:"number"`
	Status     string   `json:"status"`
	Accrual    *float64 `json:"accrual,omitempty"`
	UploadedAt string   `json:"uploaded_at"`
}
