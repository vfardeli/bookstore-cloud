package models

type Payment struct {
	OrderID uint    `json:"order_id"`
	Amount  float64 `json:"amount"`
	Method  string  `json:"method"`
	Status  string  `json:"status"`
}
