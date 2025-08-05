package utils

import (
	"shared/pkg/models"
	"user/internal/dto"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ComparePasswords(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// mappers
func MapUsersToDTO(users []models.User) []dto.UserDTO {
	result := make([]dto.UserDTO, len(users))
	for i, user := range users {
		result[i] = MapUserDTO(&user)
	}
	return result
}

func MapUserDTO(model *models.User) dto.UserDTO {
	return dto.UserDTO{
		ID:          model.ID,
		Name:        model.Name,
		Username:    model.Username,
		Nationality: model.Nationality,
		BirthDate:   model.BirthDate,
	}
}
