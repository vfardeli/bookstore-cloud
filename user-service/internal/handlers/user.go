package handlers

import (
	"net/http"

	"user-service/internal/db"
	"user-service/internal/models"
	"user-service/internal/utils"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	reqID := c.MustGet("RequestID").(string)
	utils.SendLog("user-service", reqID, "info", "Registering new user", nil)

	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Create(&user)
	c.JSON(http.StatusOK, user)
}

func LoginUser(c *gin.Context) {
	reqID := c.MustGet("RequestID").(string)
	utils.SendLog("user-service", reqID, "info", "Logging in user", nil)

	var req models.User
	var user models.User

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := db.DB.Where("username = ? AND password = ?", req.Username, req.Password).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user})
}
