package db

import (
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func CreatePostgresClient(masterDSN string, replicaDSN string) *gorm.DB {
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
