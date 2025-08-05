package auth

import (
	"context"
	"errors"
	"time"

	"user/internal/db"
	"user/internal/dto"
	"user/internal/repository"
	"user/internal/repository/entity"
	"user/internal/utils"

	"github.com/golang-jwt/jwt/v5"

	sharedModels "shared/pkg/models"
)

type Service struct {
	repository *repository.Repository
	jwtConfig  sharedModels.JWTConfig
}

func NewService(jwtConfig sharedModels.JWTConfig) *Service {
	return &Service{
		repository: repository.NewRepository(),
		jwtConfig:  jwtConfig,
	}
}

func (s *Service) Register(ctx context.Context, req dto.AuthRequestDTO) (*dto.TokenResponseDTO, error) {
	if len(req.Password) < 4 || len(req.Password) > 16 {
		return nil, ErrorInvalidPassword
	}

	user, _ := s.repository.GetUserByUsername(ctx, req.Username)
	if user != nil {
		return nil, ErrorUserAlreadyExists
	}

	user, err := s.repository.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	tokens, err := s.generateTokens(ctx, user.ID, user.Username)
	return tokens, err
}

func (s *Service) Login(ctx context.Context, req dto.AuthRequestDTO) (*dto.TokenResponseDTO, error) {
	user, err := s.repository.GetUserByUsername(ctx, req.Username)
	if err != nil {
		return nil, ErrorInvalidCredentials
	}

	if err := utils.ComparePasswords(user.Password, req.Password); err != nil {
		return nil, ErrorInvalidCredentials
	}

	return s.generateTokens(ctx, user.ID, user.Username)
}

func (s *Service) RefreshToken(ctx context.Context, req dto.RefreshTokenRequestDTO) (*dto.TokenResponseDTO, error) {
	var tokenRecord entity.TokenEntity
	if err := db.PostgresDB.
		WithContext(ctx).
		Where("refresh_token = ?", req.RefreshToken).
		Where("expires_at > ?", time.Now()).
		First(&tokenRecord).
		Error; err != nil {
		return nil, ErrorInvalidRefreshToken
	}

	var user entity.UserEntity
	if err := db.PostgresDB.WithContext(ctx).First(&user, tokenRecord.UserID).Error; err != nil {
		return nil, ErrorInvalidRefreshToken
	}

	if err := db.PostgresDB.WithContext(ctx).Delete(&tokenRecord).Error; err != nil {
		return nil, err
	}

	return s.generateTokens(ctx, user.ID, user.Username)
}

func (s *Service) generateTokens(ctx context.Context, userID int, username string) (*dto.TokenResponseDTO, error) {
	now := time.Now()

	accessClaims := s.buildTokenClaims(userID, username)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(s.jwtConfig.SecretKey))
	if err != nil {
		return nil, errors.New("error signing access token")
	}

	refreshClaims := s.buildTokenClaims(userID, username)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(s.jwtConfig.SecretKey))
	if err != nil {
		return nil, errors.New("error signing refresh token")
	}

	tokenRecord := entity.TokenEntity{
		UserID:       userID,
		RefreshToken: refreshTokenString,
		ExpiresAt:    now.Add(s.jwtConfig.RefreshTokenTTL),
		CreatedAt:    now,
	}

	if err := db.PostgresDB.WithContext(ctx).Create(&tokenRecord).Error; err != nil {
		return nil, errors.New("error saving refresh token")
	}

	go s.repository.CleanupExpiredTokens()

	return &dto.TokenResponseDTO{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		ExpiresIn:    int64(s.jwtConfig.AccessTokenTTL.Seconds()),
	}, nil
}

func (s *Service) buildTokenClaims(userID int, username string) *sharedModels.TokenClaims {
	now := time.Now()

	return &sharedModels.TokenClaims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtConfig.AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
}
