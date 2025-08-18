package models

type Order struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	UserID   uint   `json:"user_id"`
	BookID   uint   `json:"book_id"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

type Book struct {
	ID     uint    `json:"id"`
	ISBN   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}
