package gameService

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/game/models"
)

func (s *GameService) GetGameStatesByIDs(ctx context.Context, ids []uint64) ([]*models.GameState, error) {
    return s.gameStorage.GetGameStatesByIDs(ctx, ids)
}