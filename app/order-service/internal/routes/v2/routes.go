package v2

import (
	"github.com/gin-gonic/gin"
)

func RegisterOrderRoutes(group *gin.RouterGroup) {
	group.GET("/", OrderHandler)
}
