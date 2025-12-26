package authService

import (
    "context"
    "errors"
    "fmt"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

func (s *AuthService) GetUser(ctx context.Context, userID string) (*models.User, error) {
    user, err := s.authStorage.GetUserByID(ctx, userID)
    if err != nil {
        return nil, fmt.Errorf("get user error: %w", err)
    }
    if user == nil {
        return nil, errors.New("пользователь не найден")
    }
    return user, nil
}

func (s *AuthService) HandleRoomCompleted(ctx context.Context, userID string, roomNumber int32) error {
    updated, err := s.authStorage.UpdateUserMaxRoom(ctx, userID, roomNumber)
    if err != nil {
        return fmt.Errorf("update max room error: %w", err)
    }
    
    if updated {
        fmt.Printf("User %s max room updated to %d\n", userID, roomNumber)
    }
    
    return nil
}