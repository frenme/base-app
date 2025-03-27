package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "order/docs"
	"order/internal/db"
	"order/internal/routes/v1"
	"order/internal/routes/v2"
	"os"
)

// @title Base API
// @version 1.0
// @description Golang project
// @BasePath /api
// @securityDefinitions.apikey CustomToken
// @in header
// @name X-Custom-Token
func main() {
	db.InitConnections()
	defer db.CloseConnections()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
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
