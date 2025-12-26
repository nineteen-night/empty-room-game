package authService

import (
    "context"
    
    "github.com/nineteen-night/empty-room-game/internal/auth/models"
)

type AuthStorage interface {
    CreateUser(ctx context.Context, user *models.User) (string, error)
    GetUserByID(ctx context.Context, userID string) (*models.User, error)
    GetUserByUsername(ctx context.Context, username string) (*models.User, error)
    UpdateUserMaxRoom(ctx context.Context, userID string, roomNumber int32) (bool, error)
    
    CreatePartnership(ctx context.Context, user1ID, user2ID string) (string, error)
    GetPartnershipByID(ctx context.Context, partnershipID string) (*models.Partnership, error)
    GetPartnershipBetweenUsers(ctx context.Context, user1ID, user2ID string) (*models.Partnership, error) // ← ДОБАВИТЬ
    TerminatePartnership(ctx context.Context, partnershipID string) error
    
    HashPassword(password string) (string, error)
    CheckPasswordHash(password, hash string) bool
}

type EventSender interface {
    SendPartnershipCreated(ctx context.Context, partnershipID, user1ID, user2ID string) error
    SendPartnershipTerminated(ctx context.Context, partnershipID string) error
}

type AuthService struct {
    authStorage AuthStorage
    eventSender EventSender
    minNameLen  int
    maxNameLen  int
}

func NewAuthService(ctx context.Context, authStorage AuthStorage, minNameLen, maxNameLen int) *AuthService {
    return &AuthService{
        authStorage: authStorage,
        eventSender: nil,
        minNameLen:  minNameLen,
        maxNameLen:  maxNameLen,
    }
}

func (s *AuthService) SetEventSender(eventSender EventSender) {
    s.eventSender = eventSender
}