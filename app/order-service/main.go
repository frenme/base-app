package main

import (
	"github.com/gin-gonic/gin"
	"order/internal/db"
	"order/internal/routes"
)

func main() {
	db.InitConnections()
	defer db.CloseConnections()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	routes.RegisterOrderRoutes(router)
}
