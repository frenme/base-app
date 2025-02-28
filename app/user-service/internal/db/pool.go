package db

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"
	"time"
)

var PoolMaster *pgxpool.Pool
var PoolReplicas *pgxpool.Pool

func InitPools() {
	PoolMaster = createPool(os.Getenv("POSTGRES_MASTER_CONNECTION"))
	PoolReplicas = createPool(os.Getenv("POSTGRES_REPLICAS_CONNECTION"))
}

func ClosePools() {
	if PoolMaster != nil {
		PoolMaster.Close()
	}
	if PoolReplicas != nil {
		PoolReplicas.Close()
	}
}

func createPool(connectionString string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(connectionString)
	config.MaxConnLifetime = 30 * time.Minute
	config.MaxConnIdleTime = 5 * time.Minute
	config.HealthCheckPeriod = 1 * time.Minute
	config.MaxConns = 20
	pool, err := pgxpool.ConnectConfig(ctx, config)

	if err != nil {
		log.Fatalf("Pool of connections to PostgreSQL is failed: %v", err)
	}

	return pool
}
