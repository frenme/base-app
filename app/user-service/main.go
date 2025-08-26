package main

import (
	"os"
	"shared/pkg/middleware"
	_ "user/docs"
	"user/internal/db"
	"user/internal/repository"
	"user/internal/utils"

	sharedConfig "shared/pkg/config"
	"shared/pkg/logger"
	auth "user/internal/modules/auth"
	authRoutes "user/internal/modules/auth/handlers/v1"
	user "user/internal/modules/user"
	userRoutes "user/internal/modules/user/handlers/v1"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title User service
// @version 1.0
// @BasePath /api/user-service
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	db.InitConnections()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET(utils.APIBasePath+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger := logger.New()

	authMiddleware := middleware.AuthMiddleware(sharedConfig.JwtConfig.SecretKey)
	reqIDMiddleware := middleware.RequestIDMiddleware()

	userRepository := repository.NewRepository()

	userService := user.NewService(userRepository)
	authService := auth.NewService(sharedConfig.JwtConfig)

	userHandler := userRoutes.NewHandler(userService, logger)
	authHandler := authRoutes.NewHandler(authService, logger)

	userRoutes := userRoutes.NewRoutes(userHandler)
	authRoutes := authRoutes.NewRoutes(authHandler)

	groupV1 := router.Group(utils.APIBasePath + "/v1")
	{
		authRoutes.RegisterRoutes(groupV1)

		protectedGroup := groupV1.Group("")
		protectedGroup.Use(authMiddleware)
		protectedGroup.Use(reqIDMiddleware)
		{
			userRoutes.RegisterRoutes(protectedGroup)
		}
	}

	router.Run(":" + os.Getenv("USER_SERVICE_PORT"))
}
