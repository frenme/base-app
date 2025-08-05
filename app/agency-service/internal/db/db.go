package db

import (
	"log/slog"
	"os"
	"shared/pkg/db"
	"shared/pkg/utils"
	"time"

	"gorm.io/gorm"
)

var PostgresDB *gorm.DB
var logger *utils.Logger

func init() {
	logger = utils.CreateLogger()
}

func InitConnections() {
	maxAttempts := 100
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		logger.Info("Attempt to initialize connections: ", slog.Int("attempt", attempt))

		if PostgresDB == nil {
			PostgresDB = db.CreatePostgresClient(os.Getenv("POSTGRES_MASTER_CONNECTION"), os.Getenv("POSTGRES_REPLICAS_CONNECTION"))
		}

		if PostgresDB != nil {
			break
		}

		time.Sleep(10 * time.Second)
	}

	if PostgresDB == nil {
		logger.Error("Failed to initialize all connections after maximum attempts")
	}
}
