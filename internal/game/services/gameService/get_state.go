package gameService

import (
	"context"
	"fmt"

	"github.com/nineteen-night/empty-room-game/internal/game/models"
)

func (s *GameService) GetGameState(ctx context.Context, partnershipID string) (*models.GameState, error) {
	if err := s.ValidatePartnershipID(partnershipID); err != nil {
		return nil, err
	}

	session, err := s.gameStorage.GetGameSession(ctx, partnershipID)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения игровой сессии: %w", err)
	}
	if session == nil {
		return nil, fmt.Errorf("игровая сессия не найдена")
	}

	room, err := s.gameStorage.GetRoomByNumber(ctx, session.CurrentRoom)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения информации о комнате: %w", err)
	}
	if room == nil {
		return nil, fmt.Errorf("комната не найдена")
	}

	return &models.GameState{
		CurrentRoom: session.CurrentRoom,
		RoomInfo:    room,
	}, nil
}