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

func (s *Service) GetCards(ctx context.Context) (*PerformanceResult, error) {
	err := s.repo.SetDataMongo(ctx)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	mongoResults, err := s.repo.GetAllDataMongo(ctx)
	if err != nil {
		return nil, err
	}
	mongoDuration := time.Since(start)

	redisData, err := s.repo.GetDataRedis(ctx)
	var redisDuration time.Duration

	if errors.Is(err, redis.Nil) {
		err = s.repo.SetDataRedis(ctx, mongoResults)
		if err != nil {
			return nil, err
		}

		start = time.Now()
		_, err = s.repo.GetDataRedis(ctx)
		if err != nil {
			return nil, err
		}
		redisDuration = time.Since(start)
	} else if err != nil {
		return nil, err
	} else {
		start = time.Now()
		var cachedResults []map[string]interface{}
		err = json.Unmarshal([]byte(redisData), &cachedResults)
		if err != nil {
			return nil, err
		}
		redisDuration = time.Since(start)
	}

	return &PerformanceResult{
		MongoDuration: mongoDuration,
		RedisDuration: redisDuration,
	}, nil
}
