package models

type Book struct {
	ID     uint   `gorm:"primaryKey"`
	ISBN   string `gorm:"unique"`
	Title  string
	Author string
	Price  float64
}
