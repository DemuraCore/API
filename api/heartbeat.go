package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func Heartbeat(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"status":  "OK",
		"message": "Heartbeat successful",
		"time":    time.Now().Format(time.RFC3339),
	})

}
