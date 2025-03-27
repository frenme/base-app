package v2

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"order/internal/services"
	"os"
	"shared/pkg/utils"
)

var logger *slog.Logger

func init() {
	handler := utils.LoggerHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
	logger = slog.New(handler)
}

func OrderHandler(c *gin.Context) {
	logger.Info("log main order handler v2")
	services.GetOrderAnotherData(c)
}
