package handlers

import (
	"bookstore-cloud/internal/db"
	"bookstore-cloud/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
