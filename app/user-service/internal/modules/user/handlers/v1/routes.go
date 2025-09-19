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
	group.GET("/users", r.handler.GetUsers)
	group.GET("/users/current", r.handler.GetCurrentUser)
	group.PUT("/users/:id", r.handler.UpdateUser)
	group.GET("/ping-temp", r.handler.PingTemp)
}
