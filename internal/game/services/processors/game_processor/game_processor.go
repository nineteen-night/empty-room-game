package game_processor

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/game/models"
)

type gameService interface {
    UpsertGameSessions(ctx context.Context, sessions []*models.GameSession) error
    UpsertGameStates(ctx context.Context, states []*models.GameState) error
}

type GameProcessor struct {
    gameService gameService
}

func NewGameProcessor(gameService gameService) *GameProcessor {
    return &GameProcessor{
        gameService: gameService,
    }
}