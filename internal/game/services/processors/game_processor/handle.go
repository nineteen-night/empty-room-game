package game_processor

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/game/models"
)

func (p *GameProcessor) HandleGameSessions(ctx context.Context, session *models.GameSession) error {
    return p.gameService.UpsertGameSessions(ctx, []*models.GameSession{session})
}

func (p *GameProcessor) HandleGameStates(ctx context.Context, state *models.GameState) error {
    return p.gameService.UpsertGameStates(ctx, []*models.GameState{state})
}