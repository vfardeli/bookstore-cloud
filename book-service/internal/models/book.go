package models

type Book struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	ISBN   string  `json:"isbn" gorm:"unique"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}
