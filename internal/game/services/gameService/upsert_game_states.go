package gameService

import (
    "context"
    "fmt"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
)

func (s *GameService) UpsertGameStates(ctx context.Context, states []*models.GameState) error {
    if err := s.validateGameStates(states); err != nil {
        return err
    }

    for _, state := range states {
        sessions, err := s.GetGameSessionsByIDs(ctx, []uint64{state.GameSessionID})
        if err != nil {
            return fmt.Errorf("failed to check game session: %w", err)
        }
        
        if len(sessions) == 0 {
            return fmt.Errorf("game session not found")
        }
    }

    return s.gameStorage.UpsertGameStates(ctx, states)
}