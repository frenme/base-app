package main

import (
	"log/slog"
	"shared/pkg/config"
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
	const attempts = 100
	for i := 1; i <= attempts; i++ {
		gormDB = db.CreatePostgresClient(config.GetPostgresConfig())
		if gormDB != nil {
			sqlDB, err := gormDB.DB()
			if err == nil && sqlDB.Ping() == nil {
				break
			}
		}
		log.Info("Attempt to initialize connections: ", slog.Int("attempt", i))
		time.Sleep(2 * time.Second)
		if i == attempts {
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
		"file://migrations", "postgres", driver,
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
