package authService

import (
	"context"
	"errors"
	"testing"

	"github.com/nineteen-night/empty-room-game/internal/auth/models"
	"github.com/nineteen-night/empty-room-game/internal/auth/services/authService/mocks"
	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
)

type UserServiceSuite struct {
	suite.Suite
	ctx          context.Context
	authStorage  *mocks.AuthStorage
	authService  *AuthService
}

func (s *UserServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.authStorage = mocks.NewAuthStorage(s.T())
	s.authService = NewAuthService(s.ctx, s.authStorage, 3, 50)
}

// успешное получение пользователя
func (s *UserServiceSuite) TestGetUser_Success() {
	userID := "user-123"
	expectedUser := &models.User{
		ID:             userID,
		Username:       "testuser",
		Email:          "test@example.com",
		PasswordHash:   "hashed_password",
		MaxRoomReached: 5,
	}

	s.authStorage.EXPECT().GetUserByID(s.ctx, userID).Return(expectedUser, nil)

	user, err := s.authService.GetUser(s.ctx, userID)

	assert.NilError(s.T(), err)
	assert.Equal(s.T(), userID, user.ID)
	assert.Equal(s.T(), "testuser", user.Username)
	assert.Equal(s.T(), "test@example.com", user.Email)
	assert.Equal(s.T(), int32(5), user.MaxRoomReached)
}

// получение несуществующего пользователя
func (s *UserServiceSuite) TestGetUser_NotFound() {
	userID := "non-existent-id"
	
	s.authStorage.EXPECT().GetUserByID(s.ctx, userID).Return(nil, nil)

	user, err := s.authService.GetUser(s.ctx, userID)

	assert.ErrorContains(s.T(), err, "не найден")
	assert.Assert(s.T(), user == nil)
}

//ошибка хранилища
func (s *UserServiceSuite) TestGetUser_StorageError() {
	userID := "user-123"
	storageErr := errors.New("database error")

	s.authStorage.EXPECT().GetUserByID(s.ctx, userID).Return(nil, storageErr)

	user, err := s.authService.GetUser(s.ctx, userID)

	assert.ErrorContains(s.T(), err, "get user error")
	assert.ErrorContains(s.T(), err, storageErr.Error())
	assert.Assert(s.T(), user == nil)
}

// успешное обновление комнаты
func (s *UserServiceSuite) TestHandleRoomCompleted_Success() {
	userID := "user-123"
	roomNumber := int32(5)

	s.authStorage.EXPECT().UpdateUserMaxRoom(s.ctx, userID, roomNumber).Return(true, nil)

	err := s.authService.HandleRoomCompleted(s.ctx, userID, roomNumber)

	assert.NilError(s.T(), err)
}

//случай когда комната не новая
func (s *UserServiceSuite) TestHandleRoomCompleted_NoUpdate() {
	userID := "user-123"
	roomNumber := int32(3)

	s.authStorage.EXPECT().UpdateUserMaxRoom(s.ctx, userID, roomNumber).Return(false, nil)

	err := s.authService.HandleRoomCompleted(s.ctx, userID, roomNumber)

	assert.NilError(s.T(), err)
}

//ошибка обновления
func (s *UserServiceSuite) TestHandleRoomCompleted_StorageError() {
	// Arrange
	userID := "user-123"
	roomNumber := int32(5)
	storageErr := errors.New("update error")

	s.authStorage.EXPECT().UpdateUserMaxRoom(s.ctx, userID, roomNumber).Return(false, storageErr)

	err := s.authService.HandleRoomCompleted(s.ctx, userID, roomNumber)

	assert.ErrorContains(s.T(), err, "update max room error")
	assert.ErrorContains(s.T(), err, storageErr.Error())
}

// некорректные входные данные
func (s *UserServiceSuite) TestHandleRoomCompleted_InvalidInput() {
	tests := []struct {
		name       string
		userID     string
		roomNumber int32
		wantErr    string
	}{
		{"empty user id", "", 5, "ID пользователя не может быть пустым"},
		{"room zero", "user-123", 0, "номер комнаты должен быть положительным"},
		{"room negative", "user-123", -1, "номер комнаты должен быть положительным"},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.authService.HandleRoomCompleted(s.ctx, tt.userID, tt.roomNumber)

			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
			} else {
				assert.NilError(t, err)
			}
		})
	}
}

func TestUserServiceSuite(t *testing.T) {
	suite.Run(t, new(UserServiceSuite))
}