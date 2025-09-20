package v1

import (
	"context"
	"net/http"
	"shared/pkg/logger"
	sharedutils "shared/pkg/utils"
	rediscache "temp/internal/modules/redis-cache"
	"time"

	echopb "shared/pkg/grps/proto/echo"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Handler struct {
	service *rediscache.Service
	logger  *logger.Logger
}

func NewHandler(service *rediscache.Service, logger *logger.Logger) *Handler {
	return &Handler{service: service, logger: logger}
}

// PingUserServiceWithGRPS call with gRPC from user-service
// @Summary     Call user-service via gRPC
// @Tags        Temp
// @Produce     json
// @Param       message query string false "message" default(hi)
// @Success     200  {object}  map[string]string
// @Router      /v1/call-user-service-with-grpc [get]
func (h *Handler) PingUserServiceWithGRPS(c *gin.Context) {
	addr := "user-service-golang:50051"
	msg := c.Query("message")

	ctx, cancel := context.WithTimeout(c, 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		sharedutils.RespondWithError(c, http.StatusBadGateway, "grpc dial error: "+err.Error())
		return
	}
	defer conn.Close()

	client := echopb.NewEchoServiceClient(conn)
	resp, err := client.Echo(ctx, &echopb.EchoRequest{Message: msg})
	if err != nil {
		sharedutils.RespondWithError(c, http.StatusBadGateway, "grpc call error: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]string{"message": resp.GetMessage()})
}
