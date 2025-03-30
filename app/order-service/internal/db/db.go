package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"shared/pkg/db"
	"time"
)

var PostgresMaster *pgxpool.Pool
var PostgresReplicas *pgxpool.Pool
var MongoClient *mongo.Client

func InitConnections() {
	time.Sleep(120 * time.Second)
	PostgresMaster = db.CreatePostgresPool(os.Getenv("POSTGRES_MASTER_CONNECTION"))
	PostgresReplicas = db.CreatePostgresPool(os.Getenv("POSTGRES_REPLICAS_CONNECTION"))
	MongoClient = db.CreateMongoClient(os.Getenv("MONGO_CONNECTION"))
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
