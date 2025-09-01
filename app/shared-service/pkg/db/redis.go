package db

import (
	"context"
	"shared/pkg/logger"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis = *redis.ClusterClient

func CreateRedisClient(connectionString string, serviceName string) *redis.ClusterClient {
	var RedisDB *redis.ClusterClient
	log := logger.New()

	redisNodes := strings.Split(connectionString, ",")

	maxAttempts := 30
	for attempt := range maxAttempts {
		log.Info("Attempting to connect to Redis cluster: ", attempt, serviceName)

		RedisDB = redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:       redisNodes,
			DialTimeout: 10 * time.Second,
			MaxRetries:  3,
			ClusterSlots: func(ctx context.Context) ([]redis.ClusterSlot, error) {
				return nil, nil
			},
		})

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		_, err := RedisDB.Ping(ctx).Result()
		cancel()

		if err == nil {
			log.Info("Redis cluster is ready ", "service: ", serviceName)
			return RedisDB
		}

		log.Warn("Failed to ping Redis cluster", "error: ", err.Error(), "service: ", serviceName)
		RedisDB.Close()
		RedisDB = nil

		if attempt < maxAttempts {
			time.Sleep(5 * time.Second)
		}
	}

	log.Error("Failed to connect to Redis cluster after all attempts", "service: ", serviceName)
	return nil
}
