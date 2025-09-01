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
	for attempt := range maxAttempts {
		log.Info("Attempt to initialize connections: ", attempt, serviceName)

		if PostgresDB == nil {
			PostgresDB = connectPostgres(masterDSN, replicaDSN)
		}

		if PostgresDB != nil {
			break
		}

		time.Sleep(10 * time.Second)
	}

	if PostgresDB == nil {
		log.Error("Failed to initialize all connections after maximum attempts", serviceName)
		return nil
	}

	return PostgresDB
}

func connectPostgres(masterDSN string, replicaDSN string) *gorm.DB {
	log := logger.New()

	db, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{})
	if err != nil {
		log.Error("Failed to open PostgreSQL connection", " error: ", err)
		return nil
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Error("Failed to get underlying sql.DB", " error: ", err)
		return nil
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	if err := sqlDB.Ping(); err != nil {
		log.Error("Failed to ping PostgreSQL", " error: ", err)
		return nil
	}

	if err := db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{postgres.Open(masterDSN)},
		Replicas: []gorm.Dialector{postgres.Open(replicaDSN)},
		Policy:   dbresolver.RandomPolicy{},
	})); err != nil {
		log.Error("Failed to configure db resolver", " error: ", err)
		return nil
	}

	log.Info("Successfully connected to PostgreSQL")
	return db
}
