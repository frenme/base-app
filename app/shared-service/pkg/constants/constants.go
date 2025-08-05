package constants

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
