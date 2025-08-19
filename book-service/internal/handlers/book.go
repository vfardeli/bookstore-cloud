package handlers

import (
	"net/http"

	"book-service/internal/db"
	"book-service/internal/models"
	"book-service/internal/utils"

	"github.com/gin-gonic/gin"
)

func AddBook(c *gin.Context) {
	reqID := c.MustGet("RequestID").(string)
	utils.SendLog("book-service", reqID, "info", "Creating new book", nil)

	var book models.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&book)
	c.JSON(http.StatusOK, book)
}

func ListBooks(c *gin.Context) {
	reqID := c.MustGet("RequestID").(string)
	utils.SendLog("book-service", reqID, "info", "Fetching book list", nil)

	var books []models.Book
	db.DB.Find(&books)
	c.JSON(http.StatusOK, books)
}

func GetBook(c *gin.Context) {
	reqID := c.MustGet("RequestID").(string)
	id := c.Param("id")

	utils.SendLog("book-service", reqID, "info", "Fetching book details", map[string]interface{}{
		"book_id": id,
	})

	var book models.Book
	result := db.DB.First(&book, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	c.JSON(http.StatusOK, book)
}
