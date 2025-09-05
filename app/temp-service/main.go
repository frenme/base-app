package main

import (
	"context"
	"os"
	"shared/pkg/logger"
	"shared/pkg/middleware"
	_ "temp/docs"
	"temp/internal/db"
	rediscache "temp/internal/modules/redis-cache"
	handlers "temp/internal/modules/redis-cache/handlers/v1"
	"temp/internal/repository"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

// @title Temp service
// @version 1.0
// @BasePath /api/temp-service
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	db.InitConnections()

	db.StartKafkaConsumers(context.Background())

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET(os.Getenv("TEMP_SERVICE_PATH")+"/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	logger := logger.New()

	reqIDMiddleware := middleware.RequestIDMiddleware()

	repository := repository.NewRepository(db.PostgresDB)
	service := rediscache.NewService(repository)
	handler := handlers.NewHandler(service, logger)
	routes := handlers.NewRoutes(handler)

	groupV1 := router.Group(os.Getenv("TEMP_SERVICE_PATH") + "/v1")
	{
		protectedGroup := groupV1.Group("")
		protectedGroup.Use(reqIDMiddleware)
		{
			routes.RegisterRoutes(protectedGroup)
		}
	}

	router.Run(":" + os.Getenv("TEMP_SERVICE_PORT"))
}
