package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func CreateMongoClient(connectionString string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(connectionString)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Println("Connection to MongoDB is failed: ", err)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Println("Connection to MongoDB is failed: ", err)
	}
	return client
}
