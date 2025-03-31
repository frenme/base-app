package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"log/slog"
	"os"
	"shared/pkg/db"
	"shared/pkg/utils"
	"time"
)

var PostgresMaster *pgxpool.Pool
var PostgresReplicas *pgxpool.Pool
var MongoClient *mongo.Client
var logger *slog.Logger

func init() {
	handler := utils.LoggerHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)}
	logger = slog.New(handler)
}

func InitConnections() {
	maxAttempts := 100
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		logger.Info("Attempt %d to initialize connections: ", attempt)

		if PostgresMaster == nil || PostgresReplicas == nil {
			PostgresMaster = db.CreatePostgresPool(os.Getenv("POSTGRES_MASTER_CONNECTION"))
			PostgresReplicas = db.CreatePostgresPool(os.Getenv("POSTGRES_REPLICAS_CONNECTION"))
		}
		if MongoClient == nil {
			MongoClient = db.CreateMongoClient(os.Getenv("MONGO_CONNECTION"))
		}

		if PostgresMaster != nil && PostgresReplicas != nil && MongoClient != nil {
			break
		}

		time.Sleep(10 * time.Second)
	}

	if PostgresMaster == nil || PostgresReplicas == nil || MongoClient == nil {
		logger.Error("Failed to initialize all connections after maximum attempts")
	}
}

func CloseConnections() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if PostgresMaster != nil {
		PostgresMaster.Close()
	}
	if PostgresReplicas != nil {
		PostgresReplicas.Close()
	}
	if MongoClient != nil {
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Fatal("MongoDB error disconnecting", err)
		}
	}
}
