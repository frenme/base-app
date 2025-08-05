package user

import (
	"context"
	"errors"
	"shared/pkg/models"
	"user/internal/dto"
	"user/internal/repository"
)

type ContextKey string

const UserID ContextKey = "userID"

type Service struct {
	repo repository.RepositoryInterface
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUsers(ctx context.Context, req dto.UsersRequestDTO) ([]models.User, int64, error) {
	return s.repo.GetUsers(ctx, req)
}

func (s *Service) GetCurrentUser(ctx context.Context) (*models.User, error) {
	userID, ok := ctx.Value(UserID).(int)
	if !ok {
		return nil, errors.New("userID is not an integer")
	}

	user, err := s.repo.GetUserByID(ctx, userID)
	return user, err
}

func (s *Service) UpdateUser(ctx context.Context, id int, req dto.UpdateUserDTO) (*models.User, error) {
	return s.repo.UpdateUser(ctx, id, req)
}
