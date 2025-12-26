package authService

import (
    "context"
    "errors"
    "fmt"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

func (s *AuthService) Register(ctx context.Context, username, email, password string) (*models.User, error) {
    if err := s.ValidateRegistration(username, email, password); err != nil {
        return nil, err
    }

    existing, err := s.authStorage.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, fmt.Errorf("ошибка проверки пользователя: %w", err)
    }
    if existing != nil {
        return nil, errors.New("пользователь с таким именем уже существует")
    }

    hash, err := s.authStorage.HashPassword(password)
    if err != nil {
        return nil, fmt.Errorf("ошибка хэширования пароля: %w", err)
    }

    user := &models.User{
        Username:       username,
        Email:          email,
        PasswordHash:   hash,
        MaxRoomReached: 0,
    }
    
    userID, err := s.authStorage.CreateUser(ctx, user)
    if err != nil {
        return nil, fmt.Errorf("ошибка создания пользователя: %w", err)
    }
    
    user.ID = userID
    return user, nil
}