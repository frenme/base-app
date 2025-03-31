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

// OrderHandler
// @Summary     Get order data - v2
// @Description Some description for this route
// @Tags        order
// @Accept      json
// @Produce     json
// @Success     200  {object}  string "Info about order"
// @Failure     400  {object}  string "Bad request"
// @Router      /v2 [get]
func OrderHandler(c *gin.Context) {
	logger.Info(" log main order handler v2")
	services.GetOrderAnotherData(c)
}
