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
	group.GET("/artists", r.handler.GetArtists)
	group.POST("/artists", r.handler.CreateArtist)
	group.PUT("/artists/:id", r.handler.UpdateArtist)
	group.POST("/artists/:id/subscription", r.handler.HandleSubscription)
}
