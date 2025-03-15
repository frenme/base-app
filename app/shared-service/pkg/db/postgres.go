package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

func CreatePostgresPool(connectionString string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connectionString)
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute
	config.MaxConns = 20
	pool, err := pgxpool.ConnectConfig(ctx, config)

	if err != nil {
		log.Println("Pool of connections to PostgreSQL is failed: ", err)
	}

	return pool
}
