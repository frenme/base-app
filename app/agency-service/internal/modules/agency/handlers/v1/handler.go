package v1

import (
	"agency/internal/dto"
	"agency/internal/modules/agency"
	"context"
	"net/http"
	sharedconfig "shared/pkg/config"
	shareddto "shared/pkg/dto"
	"shared/pkg/logger"
	sharedutils "shared/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *agency.Service
	logger  *logger.Logger
}

func NewHandler(service *agency.Service, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GetAgencies
// @Summary     Get agencies
// @Tags        Agencies
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Param       take query int false "Take" default(10)
// @Param       skip query int false "Skip" default(0)
// @Param       query query string false "Query" default()
// @Success     200  {object}  dto.AgenciesResponseDTO "Agencies response"
// @Failure     400  {object}  shareddto.ErrorResponse "Bad request"
// @Router      /v1/agencies [get]
func (h *Handler) GetAgencies(c *gin.Context) {
	h.logger.Info("get all agencies")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	var req dto.AgenciesRequestDTO
	err := sharedutils.HandleQueryRequestData(c, &req)
	if err != nil {
		return
	}
	if req.Take == 0 {
		req.Take = 10
	}
	if req.Skip == 0 {
		req.Skip = 0
	}

	agencies, total, err := h.service.GetAgencies(ctx, req)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	response := dto.AgenciesResponseDTO{
		PaginationResponse: shareddto.PaginationResponse{
			Total: total,
			Take:  req.Take,
			Skip:  req.Skip,
		},
		Items: agencies,
	}
	c.JSON(http.StatusOK, response)
}

// CreateAgency
// @Summary     Create an agency
// @Tags        Agencies
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Param       request body dto.CreateAgencyDTO true "Agency data"
// @Success     201  {object}  dto.AgencyDTO "Agency response"
// @Failure     400  {object}  shareddto.ErrorResponse "Bad request"
// @Router      /v1/agencies [post]
func (h *Handler) CreateAgency(c *gin.Context) {
	h.logger.Info("create an agency")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	var req dto.CreateAgencyDTO
	sharedutils.HandleBodyRequestData(c, &req)

	sharedutils.TrimStrings(&req)

	agency, err := h.service.CreateAgency(ctx, req)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, agency)
}

// UpdateAgency
// @Summary     Update an agency
// @Tags        Agencies
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Param       id path int true "Agency ID"
// @Param       request body dto.UpdateAgencyDTO false "Agency data"
// @Success     200  {object}  dto.AgencyDTO "Agency response"
// @Failure     400  {object}  shareddto.ErrorResponse "Bad request"
// @Router      /v1/agencies/{id} [put]
func (h *Handler) UpdateAgency(c *gin.Context) {
	h.logger.Info("update an agency")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		sharedutils.RespondWithError(c, http.StatusBadRequest, "Invalid agency ID")
		return
	}

	var req dto.UpdateAgencyDTO
	sharedutils.HandleBodyRequestData(c, &req)

	sharedutils.TrimStrings(&req)

	agency, err := h.service.UpdateAgency(ctx, id, req)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, agency)
}
