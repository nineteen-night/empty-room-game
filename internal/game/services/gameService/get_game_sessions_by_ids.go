package gameService

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/game/models"
)

func (s *GameService) GetGameSessionsByIDs(ctx context.Context, ids []uint64) ([]*models.GameSession, error) {
    return s.gameStorage.GetGameSessionsByIDs(ctx, ids)
}