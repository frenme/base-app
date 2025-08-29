package main

import (
	_ "agency/docs"
	"agency/internal/db"
	"agency/internal/modules/agency"
	agencyhandlers "agency/internal/modules/agency/handlers/v1"
	"agency/internal/modules/artist"
	artisthandlers "agency/internal/modules/artist/handlers/v1"
	"agency/internal/repository"
	"os"
	"shared/pkg/logger"
	"shared/pkg/middleware"

	sharedconfig "shared/pkg/config"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginswagger "github.com/swaggo/gin-swagger"
)

// @title Agency service
// @version 1.0
// @BasePath /api/agency-service
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	db.InitConnections()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET(os.Getenv("AGENCY_SERVICE_PATH")+"/swagger/*any", ginswagger.WrapHandler(swaggerfiles.Handler))

	logger := logger.New()

	authMiddleware := middleware.AuthMiddleware(sharedconfig.JwtConfig.SecretKey)
	reqIDMiddleware := middleware.RequestIDMiddleware()

	repository := repository.NewRepository(db.PostgresDB)

	artistService := artist.NewService(repository)
	artistHandler := artisthandlers.NewHandler(artistService, logger)
	artistRoutes := artisthandlers.NewRoutes(artistHandler)

	agencyService := agency.NewService(repository)
	agencyHandler := agencyhandlers.NewHandler(agencyService, logger)
	agencyRoutes := agencyhandlers.NewRoutes(agencyHandler)

	groupV1 := router.Group(os.Getenv("AGENCY_SERVICE_PATH") + "/v1")
	{
		protectedGroup := groupV1.Group("")
		protectedGroup.Use(authMiddleware)
		protectedGroup.Use(reqIDMiddleware)
		{
			artistRoutes.RegisterRoutes(protectedGroup)
			agencyRoutes.RegisterRoutes(protectedGroup)
		}
	}
	// groupV1.Use(authMiddleware)
	// artistRoutes.RegisterRoutes(groupV1)
	// agencyRoutes.RegisterRoutes(groupV1)

	router.Run(":" + os.Getenv("AGENCY_SERVICE_PORT"))
}
