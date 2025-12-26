package authService

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
	"gotest.tools/v3/assert"
)

type ValidateServiceSuite struct {
	suite.Suite
	authService *AuthService
}

func (s *ValidateServiceSuite) SetupTest() {
	s.authService = &AuthService{
		minNameLen: 3,
		maxNameLen: 50,
	}
}

// валидация регистрации
func (s *ValidateServiceSuite) TestValidateRegistration() {
	tests := []struct {
		name     string
		username string
		email    string
		password string
		wantErr  string
	}{
		// Успешно
		{"valid", "validuser", "valid@example.com", "Password123", ""},
		{"valid min length", "usr", "test@example.com", "Pass12", ""},
		{"valid max length", strings.Repeat("u", 50), "test@example.com", "Pass12", ""},
		
		// Ошибки имени пользователя
		{"username too short", "ab", "test@example.com", "Password123", "имя пользователя должно быть от 3 до 50 символов"},
		{"username too long", strings.Repeat("u", 51), "test@example.com", "Password123", "имя пользователя должно быть от 3 до 50 символов"},
		{"empty username", "", "test@example.com", "Password123", "имя пользователя должно быть от 3 до 50 символов"},
		
		// Ошибки пароля
		{"password too short", "validuser", "test@example.com", "12345", "пароль должен быть не менее 6 символов"},
		{"empty password", "validuser", "test@example.com", "", "пароль должен быть не менее 6 символов"},
		
		// Ошибки email
		{"invalid email no @", "validuser", "invalid-email", "Password123", "некорректный email"},
		{"invalid email no domain", "validuser", "invalid@", "Password123", "некорректный email"},
		{"invalid email only domain", "validuser", "@domain.com", "Password123", "некорректный email"},
		{"invalid email spaces", "validuser", "test @example.com", "Password123", "некорректный email"},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.authService.ValidateRegistration(tt.username, tt.email, tt.password)

			if tt.wantErr == "" {
				assert.NilError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}

// валидация входа
func (s *ValidateServiceSuite) TestValidateLogin() {
	tests := []struct {
		name     string
		username string
		password string
		wantErr  string
	}{
		// Успешно
		{"valid", "validuser", "password", ""},
		{"valid min length", "usr", "password", ""},
		{"valid max length", strings.Repeat("u", 50), "password", ""},
		
		// Ошибки
		{"username too short", "ab", "password", "имя пользователя должно быть от 3 до 50 символов"},
		{"username too long", strings.Repeat("u", 51), "password", "имя пользователя должно быть от 3 до 50 символов"},
		{"empty username", "", "password", "имя пользователя должно быть от 3 до 50 символов"},
		{"empty password", "validuser", "", "пароль не может быть пустым"},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.authService.ValidateLogin(tt.username, tt.password)

			if tt.wantErr == "" {
				assert.NilError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}

//валидация создания партнёрства
func (s *ValidateServiceSuite) TestValidatePartnershipCreation() {
	tests := []struct {
		name    string
		user1ID string
		user2ID string
		wantErr string
	}{
		{"valid", "user1", "user2", ""},
		{"empty user1", "", "user2", "ID пользователей не могут быть пустыми"},
		{"empty user2", "user1", "", "ID пользователей не могут быть пустыми"},
		{"both empty", "", "", "ID пользователей не могут быть пустыми"},
		{"same user", "user1", "user1", "нельзя создать партнёрство с самим собой"},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.authService.ValidatePartnershipCreation(tt.user1ID, tt.user2ID)

			if tt.wantErr == "" {
				assert.NilError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}

//валидация расторжения партнёрства
func (s *ValidateServiceSuite) TestValidatePartnershipTermination() {
	tests := []struct {
		name           string
		partnershipID  string
		wantErr        string
	}{
		{"valid", "partnership-123", ""},
		{"empty", "", "ID партнёрства не может быть пустым"},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			err := s.authService.ValidatePartnershipTermination(tt.partnershipID)

			if tt.wantErr == "" {
				assert.NilError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}

//валидация завершения комнаты
func (s *ValidateServiceSuite) TestValidateRoomCompletion() {
	tests := []struct {
		name       string
		userID     string
		roomNumber int32
		wantErr    string
	}{
		{"valid", "user-123", 5, ""},
		{"empty user", "", 5, "ID пользователя не может быть пустым"},
		{"room zero", "user-123", 0, "номер комнаты должен быть положительным"},
		{"room negative", "user-123", -1, "номер комнаты должен быть положительным"},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Act
			err := s.authService.ValidateRoomCompletion(tt.userID, tt.roomNumber)

			// Assert
			if tt.wantErr == "" {
				assert.NilError(t, err)
			} else {
				assert.ErrorContains(t, err, tt.wantErr)
			}
		})
	}
}

//проверка email
func (s *ValidateServiceSuite) TestIsValidEmail() {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		// Валидные
		{"valid simple", "test@example.com", true},
		{"valid with dot", "test.name@example.com", true},
		{"valid with plus", "test+tag@example.com", true},
		{"valid with dash", "test-name@example.com", true},
		{"valid with numbers", "test123@example.com", true},
		{"valid short domain", "test@ex.co", true},
		{"valid long domain", "test@example-with-hyphen.com", true},
		{"valid very short", "a@b.c", true}, // Технически валиден по RFC
		
		// Невалидные
		{"no at symbol", "invalid-email", false},
		{"double at", "test@@example.com", false},
		{"no domain", "test@", false},
		{"no local part", "@example.com", false},
		{"space in email", "test @example.com", false},
		{"too long local part", strings.Repeat("a", 65) + "@example.com", false}, // > 64 символов
		{"empty local part", "@example.com", false},
		{"too long", strings.Repeat("a", 250) + "@example.com", false},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			result := s.authService.IsValidEmail(tt.email)

			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValidateServiceSuite(t *testing.T) {
	suite.Run(t, new(ValidateServiceSuite))
}