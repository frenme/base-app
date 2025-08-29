package v1

import (
	"context"
	"net/http"
	sharedconfig "shared/pkg/config"
	shareddto "shared/pkg/dto"
	"shared/pkg/logger"
	sharedutils "shared/pkg/utils"
	"strconv"
	"user/internal/dto"
	"user/internal/modules/user"
	"user/internal/utils"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *user.Service
	logger  *logger.Logger
}

func NewHandler(service *user.Service, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// GetUsers
// @Summary     Get users
// @Tags        Users
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Param       take query int false "Take" default(10)
// @Param       skip query int false "Skip" default(0)
// @Param       query query string false "Query" default()
// @Success     200  {object}  dto.UsersResponseDTO "Users response"
// @Failure     400  {object}  shareddto.ErrorResponse "Bad request"
// @Router      /v1/users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	h.logger.Info("get all users")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	var req dto.UsersRequestDTO
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

	users, total, err := h.service.GetUsers(ctx, req)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	response := dto.UsersResponseDTO{
		PaginationResponse: shareddto.PaginationResponse{
			Total: total,
			Take:  req.Take,
			Skip:  req.Skip,
		},
		Items: utils.MapUsersToDTO(users),
	}
	c.JSON(http.StatusOK, response)
}

// GetCurrentUser
// @Summary     Get current user
// @Tags        Users
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Success     200  {object}  dto.UserDTO "User response"
// @Failure     400  {object}  shareddto.ErrorResponse "Bad request"
// @Router      /v1/users/current [get]
func (h *Handler) GetCurrentUser(c *gin.Context) {
	h.logger.Info("get current user")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	user, err := h.service.GetCurrentUser(ctx)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, utils.MapUserDTO(user))
}

// UpdateUser
// @Summary     Update a user
// @Tags        Users
// @Security    Bearer
// @Accept      json
// @Produce     json
// @Param       id path int true "User ID"
// @Param       request body dto.UpdateUserDTO false "User data"
// @Success     200  {object}  dto.UserDTO "User response"
// @Failure     400  {object}  shareddto.ErrorResponse "Bad request"
// @Router      /v1/users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	h.logger.Info("update a user")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	id := c.Param("id")
	userId, err := strconv.Atoi(id)
	if err != nil {
		sharedutils.RespondWithError(c, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var req dto.UpdateUserDTO
	sharedutils.HandleBodyRequestData(c, &req)

	sharedutils.TrimStrings(&req)

	user, err := h.service.UpdateUser(ctx, userId, req)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	response := utils.MapUserDTO(user)
	h.logger.Info(response)
	c.JSON(http.StatusOK, response)
}
