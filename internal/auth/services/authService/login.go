package authService

import (
    "context"
    "errors"
    "fmt"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

func (s *AuthService) Login(ctx context.Context, username, password string) (*models.User, error) {
    if err := s.ValidateLogin(username, password); err != nil {
        return nil, err
    }

    user, err := s.authStorage.GetUserByUsername(ctx, username)
    if err != nil {
        return nil, fmt.Errorf("get user error: %w", err)
    }
    if user == nil {
        return nil, errors.New("user not found")
    }
    
    if !s.authStorage.CheckPasswordHash(password, user.PasswordHash) {
        return nil, errors.New("invalid password")
    }
    
    return user, nil
}