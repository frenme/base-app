package main

import (
	_ "agency/docs"
	"agency/internal/db"
	"agency/internal/modules/agency"
	agencyHandlers "agency/internal/modules/agency/handlers/v1"
	artist "agency/internal/modules/artist"
	artistHandlers "agency/internal/modules/artist/handlers/v1"
	"agency/internal/repository"
	"agency/internal/utils"
	"os"
	"shared/pkg/middleware"
	sharedUtils "shared/pkg/utils"

	sharedConstants "shared/pkg/constants"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
	router.GET(utils.APIBasePath+"/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logger := sharedUtils.CreateLogger()

	authMiddleware := middleware.AuthMiddleware(sharedConstants.JwtConfig.SecretKey)

	repository := repository.NewRepository(db.PostgresDB)

	artistService := artist.NewService(repository)
	artistHandler := artistHandlers.NewHandler(artistService, logger)
	artistRoutes := artistHandlers.NewRoutes(artistHandler)

	agencyService := agency.NewService(repository)
	agencyHandler := agencyHandlers.NewHandler(agencyService, logger)
	agencyRoutes := agencyHandlers.NewRoutes(agencyHandler)

	groupV1 := router.Group(utils.APIBasePath + "/v1")
	{
		protectedGroup := groupV1.Group("")
		protectedGroup.Use(authMiddleware)
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
