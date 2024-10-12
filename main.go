package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	// hello world
	router.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "hello world",
		})
	})

	router.Run("localhost:8080")
}
