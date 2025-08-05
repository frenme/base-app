package dto

import (
	"shared/pkg/dto"
)

type AgencyDTO struct {
	ID       int     `json:"id"`
	Name     string  `json:"name" example:"agency name"`
	ImageUrl *string `json:"imageUrl" example:"https://image-service.com/image.png"`
	Status   string  `json:"status" example:"active"`
}

type CreateAgencyDTO struct {
	Name     string  `json:"name" binding:"required" example:"agency name"`
	ImageUrl *string `json:"imageUrl" example:"https://image-service.com/image.png"`
}

type UpdateAgencyDTO struct {
	Name     *string `json:"name" example:"agency name"`
	ImageUrl *string `json:"imageUrl" example:"https://image-service.com/image.png"`
	Status   *string `json:"status" example:"active"`
}

type AgenciesRequestDTO struct {
	dto.PaginationRequest
	Query *string `form:"query" example:"some search query"`
}

type AgenciesResponseDTO struct {
	dto.PaginationResponse
	Items []AgencyDTO `json:"items"`
}
