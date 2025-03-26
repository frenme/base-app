package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"order/internal/db"
	"shared/pkg/models"
)

func GetMongoData(user models.User) int64 {
	ctx := context.Background()
	mongoDb := db.MongoClient.Database("exampleDb")
	collection := mongoDb.Collection("users")
	collection.InsertOne(ctx, user)
	mongoData, _ := collection.CountDocuments(ctx, bson.M{})
	return mongoData
}

func GetPostgresData() string {
	var postgresData string
	ctx := context.Background()
	db.PostgresMaster.QueryRow(ctx, "SELECT current_database()").Scan(&postgresData)
	return postgresData
}
