package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{
			"service": "Order",
		})
	})

	r.Run(":" + os.Getenv("ORDER_SERVICE_PORT"))
}
