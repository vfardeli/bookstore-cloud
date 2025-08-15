package handlers

import (
	"book-service/internal/db"
	"book-service/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddBook(c *gin.Context) {
	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&book)
	c.JSON(http.StatusOK, book)
}

func ListBooks(c *gin.Context) {
	var books []models.Book
	db.DB.Find(&books)
	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	id := c.Param("id")
	var book models.Book
	result := db.DB.First(&book, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}
