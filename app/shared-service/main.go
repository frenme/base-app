package main

import (
	"database/sql"
	"log/slog"
	"os"
	"shared/pkg/utils"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var logger *utils.Logger

func init() {
	logger = utils.CreateLogger()
}

func main() {
	logger.Info("Shared service started")

	makeMigration()
}

func makeMigration() {
	connectionString := os.Getenv("POSTGRES_MASTER_CONNECTION")

	const attempts = 100
	for i := 1; i <= attempts; i++ {
		db, err := sql.Open("pgx", connectionString)
		if err == nil && db.Ping() == nil {
			db.Close()
			break
		}
		logger.Info("Attempt to initialize connections: ", slog.Int("attempt", i))
		time.Sleep(2 * time.Second)
		if i == attempts {
			logger.Error("Failed to initialize PostgreSQL connection after maximum attempts")
		}
	}

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		logger.Errorf("sql.Open: %v", err)
	}
	defer db.Close()

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		logger.Errorf("pgx.WithInstance: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", "pgx", driver,
	)
	if err != nil {
		logger.Errorf("migrate.NewWithDatabaseInstance: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Errorf("m.Up: %v", err)
	}

	logger.Info("Migrations applied successfully")
}
