package dto

import (
	"shared/pkg/dto"
	"time"
)

type ArtistDTO struct {
	ID           int       `json:"id"`
	AgencyID     int       `json:"agencyId" example:"1"`
	Name         string    `json:"name" example:"Artist John Doe"`
	DebutDate    time.Time `json:"debutDate" example:"2025-01-01T00:00:00Z"`
	NumFollowers int       `json:"numFollowers" example:"100"`
	ImageUrl     *string   `json:"imageUrl" example:"https://image-service.com/image.png"`
	Status       string    `json:"status" example:"disband"`
}

type CreateArtistDTO struct {
	AgencyID  int       `json:"agencyId" binding:"required" example:"1"`
	Name      string    `json:"name" binding:"required" example:"Artist John Doe"`
	DebutDate time.Time `json:"debutDate" binding:"required" example:"2025-01-01T00:00:00Z"`
	ImageUrl  *string   `json:"imageUrl" example:"https://image-service.com/image.png"`
}

type UpdateArtistDTO struct {
	AgencyID  *int       `json:"agencyId" example:"1"`
	Name      *string    `json:"name" example:"Artist John Doe"`
	DebutDate *time.Time `json:"debutDate" example:"2025-01-01T00:00:00Z"`
	ImageUrl  *string    `json:"imageUrl" example:"https://image-service.com/image.png"`
	Status    *string    `json:"status" example:"disband"`
}

type ArtistsRequestDTO struct {
	dto.PaginationRequest
	Query          *string `form:"query" example:"some search query"`
	AgencyID       *int    `form:"agencyId" example:"1"`
	IsUserFollowed *bool   `form:"isUserFollowed" example:"true"`
}

type ArtistsResponseDTO struct {
	dto.PaginationResponse
	Items []ArtistDTO `json:"items"`
}

type SubscriptionResponseDTO struct {
	ArtistId   int  `json:"artistId" example:"10"`
	IsFollowed bool `json:"isFollowed" example:"true"`
}
