package gameService

import (
	"context"
	"fmt"

	"github.com/nineteen-night/empty-room-game/internal/game/models"
)

type GameStorage interface {
	CreateGameSession(ctx context.Context, partnershipID, user1ID, user2ID string) error
	DeleteGameSession(ctx context.Context, partnershipID string) error
	GetGameSession(ctx context.Context, partnershipID string) (*models.GameSession, error)
	UpdateCurrentRoom(ctx context.Context, partnershipID string, newRoom int32) error
	GetRoomByNumber(ctx context.Context, roomNumber int32) (*models.Room, error)
	GetMaxRoomNumber(ctx context.Context) (int32, error)
	GetGameSessionsByUserID(ctx context.Context, userID string) ([]*models.GameSession, error)
}

type EventSender interface {
	SendRoomCompleted(ctx context.Context, userID string, roomNumber int32) error
}

type GameService struct {
	gameStorage GameStorage
	eventSender EventSender
}

func NewGameService(ctx context.Context, gameStorage GameStorage) *GameService {

	return &GameService{
		gameStorage: gameStorage,
		eventSender: nil,
	}
}

func (s *GameService) SetEventSender(eventSender EventSender) {
	s.eventSender = eventSender
}

func (s *GameService) HandlePartnershipCreated(ctx context.Context, partnershipID, user1ID, user2ID string) error {
    if err := s.ValidateGameSessionCreation(partnershipID, user1ID, user2ID); err != nil {
        return err
    }

    existingSession, err := s.gameStorage.GetGameSession(ctx, partnershipID)
    if err != nil {
        return fmt.Errorf("ошибка проверки существующей сессии: %w", err)
    }
    
    if existingSession != nil {
        return nil
    }
    
    err = s.gameStorage.CreateGameSession(ctx, partnershipID, user1ID, user2ID)
    if err != nil {
        return fmt.Errorf("ошибка создания игровой сессии: %w", err)
    }
    
    return nil
}

func (s *GameService) HandlePartnershipTerminated(ctx context.Context, partnershipID string) error {
    if err := s.ValidatePartnershipID(partnershipID); err != nil {
        return err
    }

    err := s.gameStorage.DeleteGameSession(ctx, partnershipID)
    if err != nil {
        return fmt.Errorf("ошибка удаления игровой сессии: %w", err)
    }
    
    return nil
}