package gameService

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/game/models"
)

type GameStorage interface {
    GetGameSessionsByIDs(ctx context.Context, ids []uint64) ([]*models.GameSession, error)
    UpsertGameSessions(ctx context.Context, sessions []*models.GameSession) error
    GetGameStatesByIDs(ctx context.Context, ids []uint64) ([]*models.GameState, error)
    UpsertGameStates(ctx context.Context, states []*models.GameState) error
}

type GameService struct {
    gameStorage GameStorage
    minNameLen  int
    maxNameLen  int
}

func NewGameService(ctx context.Context, gameStorage GameStorage, minNameLen, maxNameLen int) *GameService {
    return &GameService{
        gameStorage: gameStorage,
        minNameLen:  minNameLen,
        maxNameLen:  maxNameLen,
    }
}