package models

type Payment struct {
	OrderID uint    `json:"orderId"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
	Status  string  `json:"status"`
}
