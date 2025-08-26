package v1

import (
	"card/internal/dto"
	"card/internal/modules/card"
	"context"
	"net/http"
	sharedConfig "shared/pkg/config"
	sharedDTO "shared/pkg/dto"
	"shared/pkg/logger"
	sharedUtils "shared/pkg/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *card.Service
	logger  *logger.Logger
}

func NewHandler(service *card.Service, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GetCards
// @Summary     Get cards
// @Tags        Cards
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Param       take query int false "Take" default(10)
// @Param       skip query int false "Skip" default(0)
// @Param       query query string false "Query" default()
// @Success     200  {object}  dto.CardsResponseDTO "Cards response"
// @Failure     400  {object}  sharedDTO.ErrorResponse "Bad request"
// @Router      /v1/cards [get]
func (h *Handler) GetCards(c *gin.Context) {
	h.logger.Info("get all cards")

	ctx, cancel := context.WithTimeout(c, sharedConfig.Timeout)
	defer cancel()

	var req dto.CardsRequestDTO
	err := sharedUtils.HandleQueryRequestData(c, &req)
	if err != nil {
		return
	}
	if req.Take == 0 {
		req.Take = 10
	}
	if req.Skip == 0 {
		req.Skip = 0
	}

	cards, total, err := h.service.GetCards(ctx, req)
	if err != nil {
		sharedUtils.HandleError(c, err)
		return
	}

	response := dto.CardsResponseDTO{
		PaginationResponse: sharedDTO.PaginationResponse{
			Total: total,
			Take:  req.Take,
			Skip:  req.Skip,
		},
		Items: cards,
	}
	c.JSON(http.StatusOK, response)
}
