package dto

import (
	"shared/pkg/dto"
	"time"
)

type UserDTO struct {
	ID          int        `json:"id"`
	Username    string     `json:"username" example:"john_doe"`
	Name        *string    `json:"name" example:"John Doe"`
	Nationality *string    `json:"nationality" example:"USA"`
	BirthDate   *time.Time `json:"birthDate" example:"2025-01-01T00:00:00Z"`
}

type UpdateUserDTO struct {
	Name        *string    `json:"name" example:"John Doe"`
	Username    *string    `json:"username" example:"john_doe"`
	Password    *string    `json:"password"  example:"password"`
	Nationality *string    `json:"nationality" example:"USA"`
	BirthDate   *time.Time `json:"birthDate" example:"2025-01-01T00:00:00Z"`
}

type UsersRequestDTO struct {
	dto.PaginationRequest
	Query *string `form:"query" example:"some search query"`
}

type UsersResponseDTO struct {
	dto.PaginationResponse
	Items []UserDTO `json:"items"`
}
