package card

import (
	"card/internal/dto"
	"card/internal/repository"
	"context"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetCards(ctx context.Context, req dto.CardsRequestDTO) ([]dto.CardDTO, int64, error) {
	entities, total, err := s.repo.GetCards(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.CardDTO, len(entities))
	for i, entity := range entities {
		result[i] = dto.CardDTO{
			ID:   entity.ID,
			Name: entity.Name,
		}
	}

	return result, total, nil
}
