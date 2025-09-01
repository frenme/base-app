package main

import (
	"log/slog"
	"os"
	"shared/pkg/db"
	"shared/pkg/logger"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

var log *logger.Logger

func init() {
	log = logger.New()
}

func main() {
	log.Info("Shared service started")

	makeMigration()
}

func makeMigration() {
	var gormDB *gorm.DB
	const maxAttempts = 100
	for attempt := range maxAttempts {
		gormDB = db.GetPostgresClient(
			os.Getenv("POSTGRES_MASTER_CONNECTION"),
			os.Getenv("POSTGRES_REPLICAS_CONNECTION"),
			"shared-service",
		)
		if gormDB != nil {
			sqlDB, err := gormDB.DB()
			if err == nil && sqlDB.Ping() == nil {
				break
			}
		}
		log.Info("Attempt to initialize connections: ", slog.Int("attempt", attempt))
		time.Sleep(2 * time.Second)
		if attempt == maxAttempts {
			log.Error("Failed to initialize PostgreSQL connection after maximum attempts")
			return
		}
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		log.Errorf("Failed to get sql.DB from GORM: %v", err)
		return
	}
	defer sqlDB.Close()

	driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		log.Errorf("postgres.WithInstance: %v", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/shared-service/migrations", "postgres", driver,
	)
	if err != nil {
		log.Errorf("migrate.NewWithDatabaseInstance: %v", err)
		return
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Errorf("m.Up: %v", err)
		return
	}

	log.Info("Migrations applied successfully")
}
