package dto

import "shared/pkg/dto"

type CardDTO struct {
	ID            int    `json:"id"`
	Name          string `json:"name" example:"Card Name"`
	Description   string `json:"description" example:"Card Description"`
	FrontImageURL string `json:"front_image_url" example:"Card Front Image URL"`
	BackImageURL  string `json:"back_image_url" example:"Card Back Image URL"`
	NumUsersOwn   int    `json:"num_users_own" example:"5"`
	NumUsersWish  int    `json:"num_users_wish" example:"10"`
	Status        string `json:"status" example:"Card Status"`
}

type CardsRequestDTO struct {
	dto.PaginationRequest
	Query *string `form:"query" example:"some search query"`
}

type CardsResponseDTO struct {
	dto.PaginationResponse
	Items []CardDTO `json:"items"`
}
