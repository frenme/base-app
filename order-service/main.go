package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"shared/pkg/models"
	"shared/pkg/utils"
)

func main() {
	router := gin.Default()
	router.GET("/test", pingHandler)
	router.Run(":" + os.Getenv("ORDER_SERVICE_PORT"))
}

func pingHandler(c *gin.Context) {
	user := models.User{Name: "egor1"}
	fmt.Println(user)
	utils.HelperCalc(123)

	c.JSON(http.StatusOK, gin.H{
		"message": user,
	})
}
