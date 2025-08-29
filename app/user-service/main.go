package main

import (
	"os"
	"shared/pkg/middleware"
	_ "user/docs"
	"user/internal/db"
	"user/internal/repository"

	sharedconfig "shared/pkg/config"
	"shared/pkg/logger"
	auth "user/internal/modules/auth"
	authhandlers "user/internal/modules/auth/handlers/v1"
	user "user/internal/modules/user"
	userhandlers "user/internal/modules/user/handlers/v1"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
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
	router.GET(os.Getenv("USER_SERVICE_PATH")+"/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	logger := logger.New()

	authMiddleware := middleware.AuthMiddleware(sharedconfig.JwtConfig.SecretKey)
	reqIDMiddleware := middleware.RequestIDMiddleware()

	userRepository := repository.NewRepository(db.PostgresDB)

	userService := user.NewService(userRepository)
	authService := auth.NewService(userRepository, sharedconfig.JwtConfig)

	userHandler := userhandlers.NewHandler(userService, logger)
	authHandler := authhandlers.NewHandler(authService, logger)

	userRoutes := userhandlers.NewRoutes(userHandler)
	authRoutes := authhandlers.NewRoutes(authHandler)

	groupV1 := router.Group(os.Getenv("USER_SERVICE_PATH") + "/v1")
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
