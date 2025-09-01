// Package db provides a MongoDB client.
package db

import (
	"context"
	"shared/pkg/logger"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Mongo = *mongo.Client

func CreateMongoClient(connectionString string, serviceName string) *mongo.Client {
	log := logger.New()

	clientOptions := options.Client().
		SetMaxPoolSize(200).
		SetMaxConnIdleTime(30 * time.Second).
		ApplyURI(connectionString)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Info("Connection to MongoDB is failed: ", err, "service: ", serviceName)
	}
	if err = client.Ping(ctx, nil); err != nil {
		log.Info("Connection to MongoDB is failed: ", err, "service: ", serviceName)
	}
	return client
}
