package v1

import (
	"agency/internal/dto"
	"agency/internal/modules/artist"
	"agency/internal/utils"
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
	service *artist.Service
	logger  *logger.Logger
}

func NewHandler(service *artist.Service, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// CreateArtist
// @Summary     Create an artist
// @Tags        Artists
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Param       request body dto.CreateArtistDTO false "Create artist request"
// @Success     200  {object}  dto.ArtistDTO "Artist response"
// @Failure     400  {object}  shareddto.ErrorResponse "Bad request"
// @Router      /v1/artists [post]
func (h *Handler) CreateArtist(c *gin.Context) {
	h.logger.Info("create an artist")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	var req dto.CreateArtistDTO
	sharedutils.HandleBodyRequestData(c, &req)

	sharedutils.TrimStrings(&req)

	artist, err := h.service.CreateArtist(ctx, req)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.MapArtistDTO(artist))
}

// UpdateArtist
// @Summary     Update an artist
// @Tags        Artists
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Param       id path int true "Artist ID"
// @Param       request body dto.UpdateArtistDTO false "Update artist request"
// @Success     200  {object}  dto.ArtistDTO "Artist response"
// @Failure     400  {object}  shareddto.ErrorResponse "Bad request"
// @Router      /v1/artists/{id} [put]
func (h *Handler) UpdateArtist(c *gin.Context) {
	h.logger.Info("update an artist")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	id := c.Param("id")
	artistId, err := strconv.Atoi(id)
	if err != nil {
		sharedutils.RespondWithError(c, http.StatusBadRequest, "Invalid artist ID")
		return
	}

	var req dto.UpdateArtistDTO
	sharedutils.HandleBodyRequestData(c, &req)

	artist, err := h.service.UpdateArtist(ctx, artistId, req)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.MapArtistDTO(artist))
}

// GetArtists
// @Summary     Get artists
// @Tags        Artists
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Param       take query int false "Take" default(10)
// @Param       skip query int false "Skip" default(0)
// @Param       query query string false "Query" default()
// @Param       agencyId query int false "Agency ID" default()
// @Param       isUserFollowed query bool false "Is User Followed" default(false)
// @Success     200  {object}  dto.ArtistsResponseDTO "Artists response"
// @Failure     400  {object}  shareddto.ErrorResponse "Bad request"
// @Router      /v1/artists [get]
func (h *Handler) GetArtists(c *gin.Context) {
	h.logger.Info("get all artists")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	var req dto.ArtistsRequestDTO
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

	artists, total, err := h.service.GetArtists(ctx, req)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	response := dto.ArtistsResponseDTO{
		PaginationResponse: shareddto.PaginationResponse{
			Total: total,
			Take:  req.Take,
			Skip:  req.Skip,
		},
		Items: utils.MapArtistsToDTO(artists),
	}
	c.JSON(http.StatusOK, response)
}

// HandleSubscription
// @Summary     Handle a subscription to Artist
// @Tags        Artists
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Param       id path int true "Artist ID"
// @Success     200  {object}  dto.SubscriptionResponseDTO "Subscription response"
// @Failure     400  {object}  shareddto.ErrorResponse "Bad request"
// @Router      /v1/artists/{id}/subscription [post]
func (h *Handler) HandleSubscription(c *gin.Context) {
	h.logger.Info("handle a subscription")
	h.logger.Info(shareddto.ErrorResponse{})

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	id := c.Param("id")
	artistId, err := strconv.Atoi(id)
	if err != nil {
		sharedutils.RespondWithError(c, http.StatusBadRequest, "Invalid artist ID")
		return
	}

	isFollowed, err := h.service.HandleSubscription(ctx, artistId)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	response := dto.SubscriptionResponseDTO{IsFollowed: isFollowed, ArtistId: artistId}
	c.JSON(http.StatusOK, response)
}
