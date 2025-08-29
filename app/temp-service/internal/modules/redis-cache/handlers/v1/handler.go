package v1

import (
	"context"
	"net/http"
	sharedconfig "shared/pkg/config"
	"shared/pkg/logger"
	sharedutils "shared/pkg/utils"
	rediscache "temp/internal/modules/redis-cache"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *rediscache.Service
	logger  *logger.Logger
}

func NewHandler(service *rediscache.Service, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// HandleRedisCache
// @Summary     Handle redis cache and mongo db
// @Tags        Temp
// @Accept      json
// @Produce     json
// @Router      /v1/redis-cache [get]
func (h *Handler) HandleRedisCache(c *gin.Context) {
	h.logger.Info("get data from redis cache and mongo db")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	result, err := h.service.GetCards(ctx)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}
