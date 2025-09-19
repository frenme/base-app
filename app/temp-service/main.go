package main

import (
	"context"
	"os"
	grpcserver "shared/pkg/grps/grpcserver"
	echopb "shared/pkg/grps/proto/echo"
	"shared/pkg/logger"
	"shared/pkg/middleware"
	_ "temp/docs"
	"temp/internal/db"
	echoimpl "temp/internal/modules/grps-echo/handlers/v1"
	rediscache "temp/internal/modules/redis-cache"
	handlers "temp/internal/modules/redis-cache/handlers/v1"
	"temp/internal/repository"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
)

// @title Temp service
// @version 1.0
// @BasePath /api/temp-service
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	logger := logger.New()

	db.InitConnections()
	db.StartKafkaConsumers(context.Background())
	grpcserver.Start(context.Background(), logger, ":"+os.Getenv("GRPC_PORT"), func(s *grpc.Server) {
		echopb.RegisterEchoServiceServer(s, echoimpl.NewEchoServer())
	})

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET(os.Getenv("TEMP_SERVICE_PATH")+"/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	reqIDMiddleware := middleware.RequestIDMiddleware()

	repository := repository.NewRepository(db.PostgresDB)
	service := rediscache.NewService(repository)
	handler := handlers.NewHandler(service, logger)
	routes := handlers.NewRoutes(handler)

	// grps-echo http routes
	grpcEchoHandler := echoimpl.NewHandler(service, logger)
	grpcEchoRoutes := echoimpl.NewRoutes(grpcEchoHandler)

	groupV1 := router.Group(os.Getenv("TEMP_SERVICE_PATH") + "/v1")
	{
		protectedGroup := groupV1.Group("")
		protectedGroup.Use(reqIDMiddleware)
		{
			routes.RegisterRoutes(protectedGroup)
			grpcEchoRoutes.RegisterRoutes(protectedGroup)
		}
	}

	router.Run(":" + os.Getenv("TEMP_SERVICE_PORT"))
}
