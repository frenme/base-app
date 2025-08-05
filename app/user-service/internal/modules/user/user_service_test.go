package user

import (
	"context"
	"errors"
	"shared/pkg/models"
	"testing"
	"time"
	"user/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) GetUserByID(ctx context.Context, id int) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) UpdateUser(ctx context.Context, id int, req dto.UpdateUserDTO) (*models.User, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockRepository) GetUsers(ctx context.Context, req dto.UsersRequestDTO) ([]models.User, int64, error) {
	args := m.Called(ctx, req)
	return args.Get(0).([]models.User), args.Get(1).(int64), args.Error(2)
}

func TestGetCurrentUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := &Service{repo: mockRepo}

	testUser := &models.User{
		ID:       1,
		Username: "testUser",
		Name:     stringPtr("Test user"),
	}

	ctx := context.WithValue(context.Background(), UserID, 1)

	mockRepo.On("GetUserByID", ctx, 1).Return(testUser, nil)

	user, err := service.GetCurrentUser(ctx)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, testUser.ID, user.ID)
	assert.Equal(t, testUser.Username, user.Username)
	mockRepo.AssertExpectations(t)

	invalidCtx := context.Background()
	user, err = service.GetCurrentUser(invalidCtx)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "userID is not an integer", err.Error())
}

func TestUpdateUser(t *testing.T) {
	mockRepo := new(MockRepository)
	service := &Service{repo: mockRepo}

	userID := 1
	newName := "New name"
	birthday := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)

	updateRequest := dto.UpdateUserDTO{
		Name:      &newName,
		BirthDate: &birthday,
	}

	expectedUser := &models.User{
		ID:        userID,
		Username:  "existingUser",
		Name:      &newName,
		BirthDate: &birthday,
	}

	ctx := context.Background()

	mockRepo.On("UpdateUser", ctx, userID, updateRequest).Return(expectedUser, nil)

	user, err := service.UpdateUser(ctx, userID, updateRequest)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser.ID, user.ID)
	assert.Equal(t, *expectedUser.Name, *user.Name)
	assert.Equal(t, expectedUser.BirthDate.Format(time.RFC3339), user.BirthDate.Format(time.RFC3339))
	mockRepo.AssertExpectations(t)

	mockRepo.On("UpdateUser", ctx, 999, updateRequest).Return(nil, errors.New("user not found"))

	user, err = service.UpdateUser(ctx, 999, updateRequest)
	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "user not found", err.Error())
}

func stringPtr(s string) *string {
	return &s
}
