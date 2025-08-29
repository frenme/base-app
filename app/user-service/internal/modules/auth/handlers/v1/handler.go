package v1

import (
	"context"
	"net/http"
	"user/internal/modules/auth"

	sharedconfig "shared/pkg/config"
	shareddto "shared/pkg/dto"
	"shared/pkg/logger"
	sharedutils "shared/pkg/utils"
	"user/internal/dto"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *auth.Service
	logger  *logger.Logger
}

func NewHandler(service *auth.Service, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// Register
// @Summary     Registration
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body dto.AuthRequestDTO true "Registration data"
// @Success     200 {object} dto.TokenResponseDTO "Access tokens"
// @Failure     400 {object} shareddto.ErrorResponse "Bad request"
// @Router      /v1/auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	h.logger.Info("register request")
	h.logger.Info(shareddto.ErrorResponse{})

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	req := h.getRequestBody(c)
	if req == nil {
		sharedutils.RespondWithError(c, http.StatusInternalServerError, "Server error")
		return
	}

	tokens, err := h.service.Register(ctx, *req)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// Login
// @Summary     Login
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body dto.AuthRequestDTO true "Login data"
// @Success     200 {object} dto.TokenResponseDTO "Access tokens"
// @Failure     400 {object} shareddto.ErrorResponse "Bad request"
// @Failure     401 {object} shareddto.ErrorResponse "Invalid credentials"
// @Router      /v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	h.logger.Infof("login request")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	req := h.getRequestBody(c)
	if req == nil {
		sharedutils.RespondWithError(c, http.StatusInternalServerError, "Server error")
		return
	}

	tokens, err := h.service.Login(ctx, *req)
	if err != nil {
		h.respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, tokens)
}

// RefreshToken
// @Summary     Refresh token
// @Tags        Auth
// @Accept      json
// @Produce     json
// @Param       request body dto.RefreshTokenRequestDTO true "Refresh token"
// @Success     200 {object} dto.TokenResponseDTO "New access tokens"
// @Failure     400 {object} shareddto.ErrorResponse "Bad request"
// @Failure     401 {object} shareddto.ErrorResponse "Invalid token"
// @Router      /v1/auth/refresh [post]
func (h *Handler) RefreshToken(c *gin.Context) {
	h.logger.Info("refresh token request")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	var req dto.RefreshTokenRequestDTO
	sharedutils.HandleBodyRequestData(c, &req)

	tokens, err := h.service.RefreshToken(ctx, req)
	if err != nil {
		if err == auth.ErrorInvalidRefreshToken {
			sharedutils.RespondWithError(c, http.StatusUnauthorized, err.Error())
			return
		}
		sharedutils.RespondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, tokens)
}

func (h *Handler) getRequestBody(c *gin.Context) *dto.AuthRequestDTO {
	var req dto.AuthRequestDTO
	sharedutils.HandleBodyRequestData(c, &req)
	sharedutils.TrimStrings(&req)
	return &req
}

func (h *Handler) respondWithError(c *gin.Context, err error) {
	switch err {
	case auth.ErrorInvalidCredentials:
		sharedutils.RespondWithError(c, http.StatusUnauthorized, err.Error())
	case auth.ErrorInvalidPassword, auth.ErrorUserAlreadyExists:
		sharedutils.RespondWithError(c, http.StatusBadRequest, err.Error())
	default:
		sharedutils.RespondWithError(c, http.StatusInternalServerError, err.Error())
	}
}
