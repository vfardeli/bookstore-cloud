package handlers

import (
	"bookstore-cloud/internal/db"
	"bookstore-cloud/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateOrder(c *gin.Context) {
	var order models.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&order)
	c.JSON(http.StatusOK, order)
}
