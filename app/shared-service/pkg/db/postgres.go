// Package db provides a PostgreSQL client.
package db

import (
	"shared/pkg/logger"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type Postgres = *gorm.DB

func GetPostgresClient(masterDSN string, replicaDSN string, serviceName string) *gorm.DB {
	var PostgresDB *gorm.DB
	log := logger.New()
	maxAttempts := 100
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		log.Info("Attempt to initialize connections: ", "attempt", attempt, "service", serviceName)

		if PostgresDB == nil {
			PostgresDB = connectPostgres(masterDSN, replicaDSN)
		}

		if PostgresDB != nil {
			break
		}

		time.Sleep(10 * time.Second)
	}

	if PostgresDB == nil {
		log.Error("Failed to initialize all connections after maximum attempts", "service", serviceName)
		return nil
	}

	return PostgresDB
}

func connectPostgres(masterDSN string, replicaDSN string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{})
	if err != nil {
		return nil
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if err := sqlDB.Ping(); err != nil {
		return nil
	}

	if err := db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{postgres.Open(masterDSN)},
		Replicas: []gorm.Dialector{postgres.Open(replicaDSN)},
		Policy:   dbresolver.RandomPolicy{},
	})); err != nil {
		return nil
	}

	return db
}
