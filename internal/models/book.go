package models

type Book struct {
	ID     uint `gorm:"primaryKey"`
	Title  string
	Author string
	Price  float64
}
