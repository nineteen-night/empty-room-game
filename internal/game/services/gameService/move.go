package gameService

import (
	"context"
	"fmt"

	"github.com/nineteen-night/empty-room-game/internal/game/models"
)

func (s *GameService) MoveToNextRoom(ctx context.Context, partnershipID string) (*models.GameState, error) {
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

	maxRoom, err := s.gameStorage.GetMaxRoomNumber(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения максимального номера комнаты: %w", err)
	}

	if session.CurrentRoom >= maxRoom {
		return nil, fmt.Errorf("вы уже достигли последней комнаты")
	}

	newRoom := session.CurrentRoom + 1

	err = s.gameStorage.UpdateCurrentRoom(ctx, partnershipID, newRoom)
	if err != nil {
		return nil, fmt.Errorf("ошибка обновления комнаты: %w", err)
	}

	if s.eventSender != nil {
		if err := s.eventSender.SendRoomCompleted(ctx, session.User1ID, newRoom); err != nil {
			fmt.Printf("Предупреждение: не удалось отправить событие для user1: %v\n", err)
		}

		if err := s.eventSender.SendRoomCompleted(ctx, session.User2ID, newRoom); err != nil {
			fmt.Printf("Предупреждение: не удалось отправить событие для user2: %v\n", err)
		}
	}

	room, err := s.gameStorage.GetRoomByNumber(ctx, newRoom)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения информации о комнате: %w", err)
	}

	return &models.GameState{
		CurrentRoom: newRoom,
		RoomInfo:    room,
	}, nil
}