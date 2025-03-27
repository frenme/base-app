package v1

import (
	"github.com/gin-gonic/gin"
	"log/slog"
	"order/internal/services"
	"os"
	"shared/pkg/utils"
)

// Логирование http запросов и трассировка
// Читаем входные параметры запроса, выдаем ошибку если они некорректны
// Декордируем входные json параметры в модель OrderRequest, выдаем ошибку если не парсится
// Проверяем обязательные поля в запросе, выдаем ошибку если чего то нет
// Прокидываем данные в сервисный слой, возвращаем

var logger *slog.Logger

func init() {
	handler := utils.LoggerHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
	logger = slog.New(handler)
}

// OrderHandler
// @Summary     Get order data
// @Description Some description for this route
// @Tags        order
// @Accept      json
// @Produce     json
// @Success     200  {object}  string "Info about order"
// @Failure     400  {object}  string "Bad request"
// @Router      /v1 [get]
func OrderHandler(c *gin.Context) {
	logger.Info("log main order handler v1")
	services.GetOrderData(c)
}

// OrderCachingHandler
// @Summary     Cache with Redis
// @Description Some description for this route
// @Tags        caching
// @Accept      json
// @Produce     json
// @Success     200  {object}  string "Info about order"
// @Failure     400  {object}  string "Bad request"
// @Router      /v1/redis [get]
func OrderCachingHandler(c *gin.Context) {
	logger.Info("log caching order handler")
	services.OrderCache(c)
}
