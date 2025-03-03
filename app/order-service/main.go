package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"order/internal/db"
	"os"
	"shared/pkg/models"
)

func main() {
	db.InitConnections()
	defer db.CloseConnections()

	router := gin.Default()
	router.GET("/", pingHandler)
	router.Run(":" + os.Getenv("ORDER_SERVICE_PORT"))
}

func pingHandler(c *gin.Context) {
	user := models.User{Name: "Alice"}
	ctx := context.Background()

	// postgresql example
	var postgresData string
	db.PostgresMaster.QueryRow(ctx, "SELECT current_database()").Scan(&postgresData)

	// mongodb example
	mongoDb := db.MongoClient.Database("exampleDb")
	collection := mongoDb.Collection("users")
	collection.InsertOne(ctx, user)
	mongoData, _ := collection.CountDocuments(ctx, bson.M{})

	c.JSON(http.StatusOK, gin.H{
		"object from another package": user,
		"postgresql data":             postgresData,
		"mongodb data":                mongoData,
	})
}
