package models

type Order struct {
	ID       uint    `json:"id"`
	UserID   uint    `json:"user_id"`
	BookID   uint    `json:"book_id"`
	Quantity int     `json:"quantity"`
	Status   string  `json:"status"`
	Amount   float64 `json:"amount"`
}
