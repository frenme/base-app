package utils

import (
	"agency/internal/dto"
	"shared/pkg/models"
)

// Mappers for Agencies
func MapAgenciesToDTO(agencies []models.Agency) []dto.AgencyDTO {
	result := make([]dto.AgencyDTO, len(agencies))
	for i, agency := range agencies {
		result[i] = MapAgencyDTO(&agency)
	}
	return result
}

func MapAgencyDTO(model *models.Agency) dto.AgencyDTO {
	return dto.AgencyDTO{
		ID:       model.ID,
		Name:     model.Name,
		ImageUrl: model.ImageUrl,
		Status:   string(models.Active),
	}
}

// Mappers for Artists
func MapArtistsToDTO(artists []models.Artist) []dto.ArtistDTO {
	result := make([]dto.ArtistDTO, len(artists))
	for i, artist := range artists {
		result[i] = MapArtistDTO(&artist)
	}
	return result
}

func MapArtistDTO(model *models.Artist) dto.ArtistDTO {
	return dto.ArtistDTO{
		ID:        model.ID,
		AgencyID:  model.AgencyID,
		Name:      model.Name,
		DebutDate: model.DebutDate,
		ImageUrl:  model.ImageUrl,
		Status:    string(models.Active),
	}
}
