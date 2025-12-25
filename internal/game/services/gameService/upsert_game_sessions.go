package gameService

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/game/models"
)

func (s *GameService) UpsertGameSessions(ctx context.Context, sessions []*models.GameSession) error {
    if err := s.validateGameSessions(sessions); err != nil {
        return err
    }

    return s.gameStorage.UpsertGameSessions(ctx, sessions)
}