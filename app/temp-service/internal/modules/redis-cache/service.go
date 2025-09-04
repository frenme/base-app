package rediscache

import (
	"context"
	"encoding/json"
	"errors"
	"temp/internal/repository"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	repo *repository.Repository
}

type PerformanceResult struct {
	MongoDuration time.Duration `json:"mongo_duration"`
	RedisDuration time.Duration `json:"redis_duration"`
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) TestCachePerformance(ctx context.Context) (*PerformanceResult, error) {
	if err := s.repo.SetDataMongo(ctx); err != nil {
		return nil, err
	}

	mongoStart := time.Now()
	mongoResults, err := s.repo.GetAllDataMongo(ctx)
	if err != nil {
		return nil, err
	}
	mongoDuration := time.Since(mongoStart)

	// ensure key exists (warm cache)
	if _, err := s.repo.GetDataRedis(ctx); errors.Is(err, redis.Nil) {
		if err := s.repo.SetDataRedis(ctx, mongoResults); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// measure only redis read (without json unmarshal)
	redisStart := time.Now()
	redisData, err := s.repo.GetDataRedis(ctx)
	if err != nil {
		return nil, err
	}
	redisDuration := time.Since(redisStart)
	var cachedResults []map[string]interface{}
	if err := json.Unmarshal([]byte(redisData), &cachedResults); err != nil {
		return nil, err
	}

	return &PerformanceResult{
		MongoDuration: mongoDuration,
		RedisDuration: redisDuration,
	}, nil
}
