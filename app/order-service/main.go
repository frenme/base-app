package main

import (
	"github.com/gin-gonic/gin"
	"order/internal/db"
	"order/internal/routes/v1"
	v2 "order/internal/routes/v2"
	"os"
)

func main() {
	db.InitConnections()
	defer db.CloseConnections()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	groupV1 := router.Group("/api/v1")
	{
		v1.RegisterOrderRoutes(groupV1)
	}
	groupV2 := router.Group("/api/v2")
	{
		v2.RegisterOrderRoutes(groupV2)
	}
	router.Run(":" + os.Getenv("ORDER_SERVICE_PORT"))
}
