package artist

import (
	"agency/internal/dto"
	"agency/internal/repository"
	"context"
	"errors"
	"shared/pkg/models"
)

type Service struct {
	repo *repository.Repository
}

func NewService(repo *repository.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateArtist(ctx context.Context, req dto.CreateArtistDTO) (*models.Artist, error) {
	_, err := s.repo.GetAgencyByID(ctx, req.AgencyID)
	if err != nil {
		return nil, err
	}

	return s.repo.CreateArtist(ctx, req)
}

func (s *Service) UpdateArtist(ctx context.Context, id int, req dto.UpdateArtistDTO) (*models.Artist, error) {
	_, err := s.repo.GetArtistByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return s.repo.UpdateArtist(ctx, id, req)
}

func (s *Service) GetArtists(ctx context.Context, req dto.ArtistsRequestDTO) ([]models.Artist, int64, error) {
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return []models.Artist{}, 0, errors.New("userID not found")
	}

	return s.repo.GetArtists(ctx, req, userID)
}

func (s *Service) HandleSubscription(ctx context.Context, artistID int) (bool, error) {
	userID, ok := ctx.Value("userID").(int)
	if !ok {
		return false, errors.New("userID not found")
	}

	_, err := s.repo.GetArtistByID(ctx, artistID)
	if err != nil {
		return false, err
	}

	return s.repo.HandleSubscription(ctx, userID, artistID)
}
