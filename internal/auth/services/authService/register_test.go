package authService

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/nineteen-night/empty-room-game/internal/auth/models"
	"github.com/nineteen-night/empty-room-game/internal/auth/services/authService/mocks"
	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
)

type RegisterServiceSuite struct {
	suite.Suite
	ctx          context.Context
	authStorage  *mocks.AuthStorage
	authService  *AuthService
}

func (s *RegisterServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.authStorage = mocks.NewAuthStorage(s.T())
	s.authService = NewAuthService(s.ctx, s.authStorage, 3, 50)
}

//успешная регистрация
func (s *RegisterServiceSuite) TestRegister_Success() {
	username := "newuser"
	email := "new@example.com"
	password := "Password123"
	hashedPassword := "hashed_password_123"
	userID := "new-user-id"

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, username).Return(nil, nil)
	s.authStorage.EXPECT().HashPassword(password).Return(hashedPassword, nil)
	s.authStorage.EXPECT().CreateUser(s.ctx, &models.User{
		Username:       username,
		Email:          email,
		PasswordHash:   hashedPassword,
		MaxRoomReached: 0,
	}).Return(userID, nil)

	user, err := s.authService.Register(s.ctx, username, email, password)

	assert.NilError(s.T(), err)
	assert.Equal(s.T(), userID, user.ID)
	assert.Equal(s.T(), username, user.Username)
	assert.Equal(s.T(), email, user.Email)
	assert.Equal(s.T(), int32(0), user.MaxRoomReached)
}

//слишком короткое имя
func (s *RegisterServiceSuite) TestRegister_UsernameTooShort() {
	username := "ab"
	email := "test@example.com"
	password := "Password123"

	user, err := s.authService.Register(s.ctx, username, email, password)

	assert.ErrorContains(s.T(), err, "имя пользователя должно быть от 3 до 50 символов")
	assert.Assert(s.T(), user == nil)
}

//слишком длинное имя
func (s *RegisterServiceSuite) TestRegister_UsernameTooLong() {
	username := strings.Repeat("a", 51)
	email := "test@example.com"
	password := "Password123"

	user, err := s.authService.Register(s.ctx, username, email, password)

	assert.ErrorContains(s.T(), err, "имя пользователя должно быть от 3 до 50 символов")
	assert.Assert(s.T(), user == nil)
}

//слишком короткий пароль
func (s *RegisterServiceSuite) TestRegister_PasswordTooShort() {
	username := "validuser"
	email := "test@example.com"
	password := "12345"

	user, err := s.authService.Register(s.ctx, username, email, password)

	assert.Error(s.T(), err, "пароль должен быть не менее 6 символов")
	assert.Assert(s.T(), user == nil)
}

//некорректный email
func (s *RegisterServiceSuite) TestRegister_InvalidEmail() {
	username := "validuser"
	email := "invalid-email"
	password := "Password123"

	user, err := s.authService.Register(s.ctx, username, email, password)

	assert.ErrorContains(s.T(), err, "некорректный email")
	assert.Assert(s.T(), user == nil)
}

//пользователь существует
func (s *RegisterServiceSuite) TestRegister_UserAlreadyExists() {
	username := "existinguser"
	email := "test@example.com"
	password := "Password123"
	existingUser := &models.User{
		ID:             "existing-id",
		Username:       username,
		Email:          email,
		PasswordHash:   "hash",
		MaxRoomReached: 5,
	}

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, username).Return(existingUser, nil)

	user, err := s.authService.Register(s.ctx, username, email, password)

	assert.Error(s.T(), err, "пользователь с таким именем уже существует")
	assert.Assert(s.T(), user == nil)
}

//ошибка хэширования пароля
func (s *RegisterServiceSuite) TestRegister_HashPasswordError() {
	username := "validuser"
	email := "test@example.com"
	password := "Password123"
	hashErr := errors.New("hash error")

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, username).Return(nil, nil)
	s.authStorage.EXPECT().HashPassword(password).Return("", hashErr)

	user, err := s.authService.Register(s.ctx, username, email, password)

	assert.ErrorContains(s.T(), err, "ошибка хэширования пароля")
	assert.ErrorContains(s.T(), err, hashErr.Error())
	assert.Assert(s.T(), user == nil)
}

//ошибка создания пользователя
func (s *RegisterServiceSuite) TestRegister_CreateUserError() {
	username := "validuser"
	email := "test@example.com"
	password := "Password123"
	hashedPassword := "hashed_password"
	createErr := errors.New("create user error")

	s.authStorage.EXPECT().GetUserByUsername(s.ctx, username).Return(nil, nil)
	s.authStorage.EXPECT().HashPassword(password).Return(hashedPassword, nil)
	s.authStorage.EXPECT().CreateUser(s.ctx, &models.User{
		Username:       username,
		Email:          email,
		PasswordHash:   hashedPassword,
		MaxRoomReached: 0,
	}).Return("", createErr)

	user, err := s.authService.Register(s.ctx, username, email, password)

	assert.ErrorContains(s.T(), err, "ошибка создания пользователя")
	assert.ErrorContains(s.T(), err, createErr.Error())
	assert.Assert(s.T(), user == nil)
}

func TestRegisterServiceSuite(t *testing.T) {
	suite.Run(t, new(RegisterServiceSuite))
}