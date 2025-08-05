package repository

import (
	"context"
	"shared/pkg/models"
	"time"
	"user/internal/db"
	"user/internal/dto"
	"user/internal/repository/entity"
	"user/internal/utils"
)

type RepositoryInterface interface {
	GetUserByID(ctx context.Context, id int) (*models.User, error)
	UpdateUser(ctx context.Context, id int, req dto.UpdateUserDTO) (*models.User, error)
	GetUsers(ctx context.Context, req dto.UsersRequestDTO) ([]models.User, int64, error)
}

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

func (r *Repository) GetUsers(ctx context.Context, req dto.UsersRequestDTO) ([]models.User, int64, error) {
	var records []entity.UserEntity
	var total int64

	// todo: handle query parameter
	query := db.PostgresDB.
		WithContext(ctx).
		Limit(req.Take).
		Offset(req.Skip)

	err := query.Find(&records).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	return r.mapUsers(records), total, nil
}

func (r *Repository) CreateUser(ctx context.Context, req dto.AuthRequestDTO) (*models.User, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	userRecord := entity.UserEntity{Username: req.Username, Password: hashedPassword}
	if err := db.PostgresDB.WithContext(ctx).Create(&userRecord).Error; err != nil {
		return nil, err
	}

	user := r.mapUser(userRecord)
	return &user, nil
}

func (r *Repository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	var record entity.UserEntity
	if err := db.PostgresDB.WithContext(ctx).First(&record, id).Error; err != nil {
		return nil, err
	}
	user := r.mapUser(record)
	return &user, nil
}

func (r *Repository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var record entity.UserEntity
	if err := db.PostgresDB.WithContext(ctx).Where("username = ?", username).First(&record).Error; err != nil {
		return nil, err
	}
	user := r.mapUser(record)
	return &user, nil
}

func (r *Repository) UpdateUser(ctx context.Context, id int, req dto.UpdateUserDTO) (*models.User, error) {
	var userRecord entity.UserEntity
	if err := db.PostgresDB.WithContext(ctx).First(&userRecord, id).Error; err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Username != nil {
		updates["username"] = *req.Username
	}
	if req.Nationality != nil {
		updates["nationality"] = *req.Nationality
	}
	if req.BirthDate != nil {
		updates["birth_date"] = *req.BirthDate
	}

	if err := db.PostgresDB.WithContext(ctx).Model(&userRecord).Updates(updates).Error; err != nil {
		return nil, err
	}

	user := r.mapUser(userRecord)
	return &user, nil
}

func (r *Repository) CleanupExpiredTokens() {
	db.PostgresDB.Where("expires_at < ?", time.Now()).Delete(&entity.TokenEntity{})
}

// mappers
func (r *Repository) mapUsers(records []entity.UserEntity) []models.User {
	users := make([]models.User, 0, len(records))
	for _, rec := range records {
		users = append(users, r.mapUser(rec))
	}

	return users
}

func (r *Repository) mapUser(record entity.UserEntity) models.User {
	return models.User{
		ID:          record.ID,
		Name:        record.Name,
		Username:    record.Username,
		Password:    record.Password,
		Nationality: record.Nationality,
		BirthDate:   record.BirthDate,
	}
}
