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
	authGroup := group.Group("/auth")
	{
		authGroup.POST("/register", r.handler.Register)
		authGroup.POST("/login", r.handler.Login)
		authGroup.POST("/refresh", r.handler.RefreshToken)
	}
}
