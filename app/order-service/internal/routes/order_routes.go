package routes

import (
	"github.com/gin-gonic/gin"
	"order/internal/handlers"
	"os"
)

func RegisterOrderRoutes(router *gin.Engine) {
	router.GET("/", handlers.OrderHandler)
	router.GET("/redis", handlers.OrderCachingHandler)

	router.Run(":" + os.Getenv("ORDER_SERVICE_PORT"))
}
