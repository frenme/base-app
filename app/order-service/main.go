package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"order/internal/db"
	"os"
	"shared/pkg/models"
)

func main() {
	db.InitPool(os.Getenv("POSTGRES_CONNECTION_STRING"))
	defer db.ClosePool()

	router := gin.Default()
	router.GET("/", pingHandler)
	router.Run(":" + os.Getenv("ORDER_SERVICE_PORT"))
}

func pingHandler(c *gin.Context) {
	user := models.User{Name: "egor1"}

	var nameDb string
	ctx := context.Background()
	err := db.Pool.QueryRow(ctx, "SELECT name FROM users").Scan(&nameDb)
	if err != nil {
		fmt.Println("Error in `SELECT name FROM users`")
		fmt.Println(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"object from another package": user,
		"name from postgres":          nameDb,
	})
}
