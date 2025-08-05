package entity

import (
	"time"
)

type TokenEntity struct {
	ID           int
	UserID       int
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

func (TokenEntity) TableName() string {
	return "refresh_tokens"
}
