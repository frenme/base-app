package db

import (
	"log/slog"
	"shared/pkg/config"
	"shared/pkg/db"
	"shared/pkg/logger"
	"time"

	"gorm.io/gorm"
)

var PostgresDB *gorm.DB
var log *logger.Logger

func init() {
	log = logger.New()
}

func InitConnections() {
	maxAttempts := 100
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Info("Attempt to initialize connections: ", slog.Int("attempt", attempt))

		if PostgresDB == nil {
			PostgresDB = db.CreatePostgresClient(config.GetPostgresConfig())
		}

		if PostgresDB != nil {
			break
		}

		time.Sleep(10 * time.Second)
	}

	if PostgresDB == nil {
		log.Error("Failed to initialize all connections after maximum attempts")
	}
}
