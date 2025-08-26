package repository

import (
	"card/internal/dto"
	"card/internal/repository/entity"
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// card
func (r *Repository) GetCards(ctx context.Context, req dto.CardsRequestDTO) ([]entity.CardEntity, int64, error) {
	var cardEntities []entity.CardEntity
	var total int64

	// todo: handle query parameter
	query := r.db.WithContext(ctx)
	if req.Query != nil && *req.Query != "" {
		search := "%" + strings.ToLower(*req.Query) + "%"
		query = query.Where("LOWER(name) LIKE ?", search)
	}

	err := query.Model(&entity.CardEntity{}).Count(&total).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count cards: %w", err)
	}

	err = query.Limit(req.Take).Offset(req.Skip).Find(&cardEntities).Error
	if err != nil {
		return nil, 0, fmt.Errorf("failed to find cards: %w", err)
	}

	return cardEntities, total, nil
}
