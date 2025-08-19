package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	r := gin.Default()
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// POST /logs -> receive logs from any service
	r.POST("/logs", func(c *gin.Context) {
		var payload map[string]interface{}
		if err := c.BindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		// Simply print logs here, can later push to file / ELK / Loki
		logrus.WithFields(payload).Info("Received log")

		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.Run(":9000") // run log-service
}
