package repository

import (
	"agency/internal/db"
	"agency/internal/dto"
	"agency/internal/repository/entity"
	"agency/internal/utils"
	"context"
	"fmt"
	shareddb "shared/pkg/db"
	"shared/pkg/models"
	sharedutils "shared/pkg/utils"
	"strings"
)

type Repository struct {
	db shareddb.Postgres
}

func NewRepository(db shareddb.Postgres) *Repository {
	return &Repository{db: db}
}

// agency
func (r *Repository) GetAgencies(ctx context.Context, req dto.AgenciesRequestDTO) ([]entity.AgencyEntity, int64, error) {
	var agencyEntities []entity.AgencyEntity
	var total int64

	// todo: handle query parameter
	query := r.db.WithContext(ctx)
	if req.Query != nil && *req.Query != "" {
		search := "%" + strings.ToLower(*req.Query) + "%"
		query = query.Where("LOWER(name) LIKE ?", search)
	}

	err := query.Model(&entity.AgencyEntity{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count agencies: %w", err)
	}

	err = query.Limit(req.Take).Offset(req.Skip).Find(&agencyEntities).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find agencies: %w", err)
	}

	return agencyEntities, total, nil
}

func (r *Repository) GetAgencyByID(ctx context.Context, id int) (*entity.AgencyEntity, error) {
	var agencyEntity entity.AgencyEntity
	err := r.db.WithContext(ctx).First(&agencyEntity, id).Error
	if err != nil {
		return nil, &sharedutils.ErrorStatus{
			Base: sharedutils.ErrorNotFound,
			Msg:  "agency not found",
		}
	}

	return &agencyEntity, nil
}

func (r *Repository) CreateAgency(ctx context.Context, req dto.CreateAgencyDTO) (*entity.AgencyEntity, error) {
	agencyEntity := entity.AgencyEntity{
		Name:     req.Name,
		ImageUrl: req.ImageUrl,
		Status:   utils.AgencyStatusActive,
	}

	err := r.db.WithContext(ctx).Create(&agencyEntity).Error
	if err != nil {
		return nil, &sharedutils.ErrorStatus{
			Base: err,
			Msg:  "failed to create agency",
		}
	}

	return &agencyEntity, nil
}

func (r *Repository) UpdateAgency(ctx context.Context, id int, req dto.UpdateAgencyDTO) (*entity.AgencyEntity, error) {
	var agencyEntity entity.AgencyEntity

	_, err := r.GetAgencyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]any)
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.ImageUrl != nil {
		updates["image_url"] = *req.ImageUrl
	}
	if req.Status != nil {
		if *req.Status != utils.AgencyStatusActive && *req.Status != utils.AgencyStatusDisable {
			return nil, fmt.Errorf("invalid status: %s", *req.Status)
		}
		updates["status"] = *req.Status
	}

	err = r.db.WithContext(ctx).Model(&agencyEntity).Updates(updates).Error
	if err != nil {
		return nil, fmt.Errorf("failed to update agency: %w", err)
	}

	err = r.db.WithContext(ctx).First(&agencyEntity, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve updated agency: %w", err)
	}

	return &agencyEntity, nil
}

// artist
func (r *Repository) CreateArtist(ctx context.Context, req dto.CreateArtistDTO) (*models.Artist, error) {
	artistRecord := entity.ArtistEntity{
		AgencyID:  req.AgencyID,
		Name:      req.Name,
		DebutDate: req.DebutDate,
		ImageUrl:  req.ImageUrl,
		Status:    string(models.Active),
	}
	if err := db.PostgresDB.WithContext(ctx).Create(&artistRecord).Error; err != nil {
		return nil, err
	}

	artist := r.mapArtist(artistRecord)
	return &artist, nil
}

func (r *Repository) UpdateArtist(ctx context.Context, id int, req dto.UpdateArtistDTO) (*models.Artist, error) {
	var artistRecord entity.ArtistEntity
	if err := db.PostgresDB.WithContext(ctx).First(&artistRecord, id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]any)
	if req.AgencyID != nil {
		updates["agency_id"] = *req.AgencyID
	}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.ImageUrl != nil {
		updates["image_url"] = *req.ImageUrl
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.DebutDate != nil {
		updates["debut_date"] = *req.DebutDate
	}

	if err := db.PostgresDB.WithContext(ctx).Model(&artistRecord).Updates(updates).Error; err != nil {
		return nil, err
	}
	return nil, nil
}

func (r *Repository) GetArtists(ctx context.Context, req dto.ArtistsRequestDTO, userID int) ([]models.Artist, int64, error) {
	var records []entity.ArtistEntity
	var total int64

	query := db.PostgresDB.
		WithContext(ctx).
		Limit(req.Take).
		Offset(req.Skip)

	if req.AgencyID != nil {
		query = query.Where("agency_id = ?", *req.AgencyID)
	}

	if req.Query != nil && *req.Query != "" {
		search := "%" + strings.ToLower(*req.Query) + "%"
		query = query.Where("LOWER(name) LIKE ?", search)
	}

	if req.IsUserFollowed != nil && *req.IsUserFollowed {
		query = query.Joins("LEFT JOIN subscriptions ON artists.id = subscriptions.artist_id AND subscriptions.user_id = ?", userID)
		query = query.Where("subscriptions.id IS NOT NULL")
	}

	err := query.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return r.mapArtists(records), total, nil
}

func (r *Repository) GetArtistByID(ctx context.Context, id int) (*models.Artist, error) {
	var artistRecord entity.ArtistEntity
	if err := db.PostgresDB.WithContext(ctx).First(&artistRecord, id).Error; err != nil {
		return nil, &sharedutils.ErrorStatus{
			Base: sharedutils.ErrorNotFound,
			Msg:  "artist not found",
		}
	}
	artist := r.mapArtist(artistRecord)
	return &artist, nil
}

func (r *Repository) HandleSubscription(ctx context.Context, userID int, artistID int) (bool, error) {
	var subscription entity.SubscriptionEntity
	result := db.PostgresDB.WithContext(ctx).Where("artist_id = ? AND user_id = ?", artistID, userID).First(&subscription)
	if result.Error == nil {
		if err := db.PostgresDB.WithContext(ctx).Delete(&subscription).Error; err != nil {
			return false, err
		}
		return false, nil
	}

	newSubscription := entity.SubscriptionEntity{UserID: userID, ArtistID: artistID}
	if err := db.PostgresDB.WithContext(ctx).Create(&newSubscription).Error; err != nil {
		return false, err
	}

	return true, nil
}

// mappers
func (r *Repository) mapArtists(records []entity.ArtistEntity) []models.Artist {
	artists := make([]models.Artist, 0, len(records))
	for _, rec := range records {
		artists = append(artists, r.mapArtist(rec))
	}

	return artists
}

func (r *Repository) mapArtist(record entity.ArtistEntity) models.Artist {
	return models.Artist{
		ID:        record.ID,
		AgencyID:  record.AgencyID,
		Name:      record.Name,
		ImageUrl:  record.ImageUrl,
		DebutDate: record.DebutDate,
		Status:    models.ArtistStatus(record.Status),
	}
}
