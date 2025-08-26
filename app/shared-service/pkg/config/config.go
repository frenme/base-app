package config

import (
	"os"
	"shared/pkg/models"
	"time"
)

const (
	Timeout = 5 * time.Second
)

var JwtConfig = models.JWTConfig{
	SecretKey:           os.Getenv("JWT_SECRET_KEY"),
	AccessTokenTTL:      time.Hour,
	RefreshTokenTTL:     7 * 24 * time.Hour,
	RefreshTokenCleanup: 24 * time.Hour,
}

func GetPostgresConfig() (masterDSN, replicaDSN string) {
	masterDSN = os.Getenv("POSTGRES_MASTER_CONNECTION")
	replicaDSN = os.Getenv("POSTGRES_REPLICA_CONNECTION")

	if replicaDSN == "" {
		replicaDSN = masterDSN
	}

	return masterDSN, replicaDSN
}
