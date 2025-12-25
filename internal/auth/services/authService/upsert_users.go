package authService

import (
    "context"
    "fmt"
    
    "golang.org/x/crypto/bcrypt"
    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

func (s *AuthService) UpsertUsers(ctx context.Context, users []*models.User) error {
    if err := s.validateUsers(users); err != nil {
        return err
    }

    for _, user := range users {
        if user.ID == 0 { 
            existingUser, err := s.authStorage.GetUserByUsername(ctx, user.Username)
            if err != nil {
                return fmt.Errorf("failed to check existing user: %w", err)
            }
            if existingUser != nil {
                return fmt.Errorf("user with username %s already exists", user.Username)
            }
        }
    }

    for _, user := range users {
        if user.Password != "" && user.ID == 0 {
            hash, err := hashPassword(user.Password)
            if err != nil {
                return fmt.Errorf("failed to hash password: %w", err)
            }
            user.PasswordHash = hash
            user.Password = "" 
        }
    }

    return s.authStorage.UpsertUsers(ctx, users)
}

func hashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(hash), nil
}