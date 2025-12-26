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

type LoginServiceSuite struct {
	suite.Suite
	ctx          context.Context
	authStorage  *mocks.AuthStorage
	authService  *AuthService
}

func (s *LoginServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.authStorage = mocks.NewAuthStorage(s.T())
	s.authService = NewAuthService(s.ctx, s.authStorage, 3, 50)
}

//успешный вход
func (s *LoginServiceSuite) TestLogin_Success() {
	username := "testuser"
	password := "password123"
	expectedUser := &models.User{
		ID:             "user-id",
		Username:       username,
		Email:          "test@example.com",
		PasswordHash:   "hashed_password",
		MaxRoomReached: 3,
	}

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, username).Return(expectedUser, nil)
	s.authStorage.EXPECT().CheckPasswordHash(password, "hashed_password").Return(true)

	user, err := s.authService.Login(s.ctx, username, password)

	assert.NilError(s.T(), err)
	assert.Equal(s.T(), username, user.Username)
	assert.Equal(s.T(), "test@example.com", user.Email)
	assert.Equal(s.T(), int32(3), user.MaxRoomReached)
}

//вход несуществующего пользователя
func (s *LoginServiceSuite) TestLogin_UserNotFound() {
	username := "nonexistent"
	password := "password123"
	
	s.authStorage.EXPECT().GetUserByUsername(s.ctx, username).Return(nil, nil)

	user, err := s.authService.Login(s.ctx, username, password)

	assert.Error(s.T(), err, "user not found")
	assert.Assert(s.T(), user == nil)
}

//неверный пароль
func (s *LoginServiceSuite) TestLogin_InvalidPassword() {
	username := "testuser"
	password := "wrongpassword"
	expectedUser := &models.User{
		ID:             "user-id",
		Username:       username,
		Email:          "test@example.com",
		PasswordHash:   "hashed_password",
		MaxRoomReached: 3,
	}

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, username).Return(expectedUser, nil)
	s.authStorage.EXPECT().CheckPasswordHash(password, "hashed_password").Return(false)

	user, err := s.authService.Login(s.ctx, username, password)

	assert.Error(s.T(), err, "invalid password")
	assert.Assert(s.T(), user == nil)
}

//ошибка хранилища при входе
func (s *LoginServiceSuite) TestLogin_StorageError() {
	username := "testuser"
	password := "password123"
	storageErr := errors.New("database connection failed")

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, username).Return(nil, storageErr)

	user, err := s.authService.Login(s.ctx, username, password)

	assert.ErrorContains(s.T(), err, "get user error")
	assert.ErrorContains(s.T(), err, storageErr.Error())
	assert.Assert(s.T(), user == nil)
}

func TestLoginServiceSuite(t *testing.T) {
	suite.Run(t, new(LoginServiceSuite))
}