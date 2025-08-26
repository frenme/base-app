package main

import (
	_ "card/docs"
	"card/internal/db"
	"card/internal/modules/card"
	cardHandlers "card/internal/modules/card/handlers/v1"
	"card/internal/repository"
	"card/internal/utils"
	"os"
	"shared/pkg/logger"
	"shared/pkg/middleware"

	sharedConfig "shared/pkg/config"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Card service
// @version 1.0
// @BasePath /api/card-service
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

	repository := repository.NewRepository(db.PostgresDB)

	cardService := card.NewService(repository)
	cardHandler := cardHandlers.NewHandler(cardService, logger)
	cardRoutes := cardHandlers.NewRoutes(cardHandler)

	groupV1 := router.Group(utils.APIBasePath + "/v1")
	{
		protectedGroup := groupV1.Group("")
		protectedGroup.Use(authMiddleware)
		protectedGroup.Use(reqIDMiddleware)
		{
			cardRoutes.RegisterRoutes(protectedGroup)
		}
	}

	router.Run(":" + os.Getenv("CARD_SERVICE_PORT"))
}
