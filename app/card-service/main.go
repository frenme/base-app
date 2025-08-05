package main

import (
	_ "card/docs"
	"card/internal/db"
	"card/internal/routes/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
)

// @title Card service
// @version 1.0
// @description Golang project
// @BasePath /api/card-service
// @securityDefinitions.apikey CustomToken
// @in header
// @name X-Custom-Token
func main() {
	db.InitConnections()
	defer db.CloseConnections()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("api/card-service/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	groupV1 := router.Group("/api/card-service/v1")
	{
		v1.RegisterOrderRoutes(groupV1)
	}
	router.Run(":" + os.Getenv("CARD_SERVICE_PORT"))
}
