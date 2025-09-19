package v1

import (
	"context"
	"net/http"
	"os"
	sharedconfig "shared/pkg/config"
	"shared/pkg/logger"
	sharedutils "shared/pkg/utils"
	rediscache "temp/internal/modules/redis-cache"
	"time"

	echopb "shared/pkg/grps/proto/echo"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Handler описывает зависимости обработчиков temp-service
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

	result, err := h.service.TestCachePerformance(ctx)
	if err != nil {
		sharedutils.HandleError(c, err)
		return
	}

	c.JSON(http.StatusOK, result)
}

// HandleKafkaMessage
// @Summary     Handle kafka message
// @Tags        Temp
// @Accept      json
// @Produce     json
// @Router      /v1/kafka-message [get]
func (h *Handler) HandleKafkaMessage(c *gin.Context) {
	h.logger.Info("send message to kafka")

	ctx, cancel := context.WithTimeout(c, sharedconfig.Timeout)
	defer cancel()

	h.service.TestKafkaProducer(ctx)

	c.JSON(http.StatusOK, "Message sent to Kafka")
}

// EchoUser вызывает gRPC Echo у user-service
// @Summary     Call user-service Echo via gRPC
// @Tags        Temp
// @Produce     json
// @Param       message query string false "message" default(hi)
// @Success     200  {object}  map[string]string
// @Router      /v1/echo-user [get]
func (h *Handler) EchoUser(c *gin.Context) {
	addr := getenv("USER_SERVICE_GRPC_ADDR", "user-service-golang:50051")
	msg := c.Query("message")
	if msg == "" {
		msg = "hi"
	}

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

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
