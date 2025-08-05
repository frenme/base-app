package agency

import (
	"agency/internal/dto"
	"agency/internal/repository"
	"context"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetAgencies(ctx context.Context, req dto.AgenciesRequestDTO) ([]dto.AgencyDTO, int64, error) {
	entities, total, err := s.repo.GetAgencies(ctx, req)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.AgencyDTO, len(entities))
	for i, entity := range entities {
		result[i] = dto.AgencyDTO{
			ID:       entity.ID,
			Name:     entity.Name,
			ImageUrl: entity.ImageUrl,
			Status:   entity.Status,
		}
	}

	return result, total, nil
}

func (s *Service) CreateAgency(ctx context.Context, req dto.CreateAgencyDTO) (*dto.AgencyDTO, error) {
	entity, err := s.repo.CreateAgency(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.AgencyDTO{
		ID:       entity.ID,
		Name:     entity.Name,
		ImageUrl: entity.ImageUrl,
		Status:   entity.Status,
	}, nil
}

func (s *Service) UpdateAgency(ctx context.Context, id int, req dto.UpdateAgencyDTO) (*dto.AgencyDTO, error) {
	entity, err := s.repo.UpdateAgency(ctx, id, req)
	if err != nil {
		return nil, err
	}

	return &dto.AgencyDTO{
		ID:       entity.ID,
		Name:     entity.Name,
		ImageUrl: entity.ImageUrl,
		Status:   entity.Status,
	}, nil
}
