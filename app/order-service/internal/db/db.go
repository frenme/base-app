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
	maxAttempts := 50
	attempt := 0
	for attempt < maxAttempts {
		if PostgresMaster == nil || PostgresReplicas == nil || MongoClient == nil {
			log.Printf("Попытка подключения к БД %d из %d\n", attempt+1, maxAttempts)
			time.Sleep(10 * time.Second)
			attempt++
			continue
		}
		break
	}

	if PostgresMaster == nil || PostgresReplicas == nil || MongoClient == nil {
		log.Fatal("Не удалось установить подключение к базам данных после нескольких попыток")
	}
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
