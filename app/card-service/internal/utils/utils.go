package utils

import (
	"card/internal/dto"
	"shared/pkg/models"
)

func MapCardsToDTO(cards []models.Card) []dto.CardDTO {
	result := make([]dto.CardDTO, len(cards))
	for i, card := range cards {
		result[i] = MapCardDTO(&card)
	}
	return result
}

func MapCardDTO(model *models.Card) dto.CardDTO {
	return dto.CardDTO{
		ID:            model.ID,
		Name:          model.Name,
		Description:   model.Description,
		FrontImageURL: model.FrontImageURL,
		BackImageURL:  model.BackImageURL,
		NumUsersOwn:   model.NumUsersOwn,
		NumUsersWish:  model.NumUsersWish,
		Status:        string(model.Status),
	}
}
