package db

import (
	"os"
	"shared/pkg/db"
	"shared/pkg/logger"
)

var PostgresDB db.Postgres
var MongoDB db.Mongo
var log *logger.Logger

func init() {
	log = logger.New()
}

func InitConnections() {
	PostgresDB = db.GetPostgresClient(
		os.Getenv("POSTGRES_MASTER_CONNECTION"),
		os.Getenv("POSTGRES_REPLICAS_CONNECTION"),
		"user-service",
	)

	MongoDB = db.CreateMongoClient(os.Getenv("MONGO_CONNECTION"), "user-service")
}
