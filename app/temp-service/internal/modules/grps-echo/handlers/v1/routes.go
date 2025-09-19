package v1

import (
	"github.com/gin-gonic/gin"
)

type Routes struct {
	handler *Handler
}

func NewRoutes(handler *Handler) *Routes {
	return &Routes{handler: handler}
}

func (r *Routes) RegisterRoutes(group *gin.RouterGroup) {
	group.GET("/echo-user", r.handler.EchoUser)
}
