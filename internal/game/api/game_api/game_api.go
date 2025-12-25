package game_api

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/game/models"
    "github.com/nineteen-night/empty-room-game/internal/game/pb/game_api"
)

type gameService interface {
    GetGameSessionsByIDs(ctx context.Context, ids []uint64) ([]*models.GameSession, error)
    UpsertGameSessions(ctx context.Context, sessions []*models.GameSession) error
    GetGameStatesByIDs(ctx context.Context, ids []uint64) ([]*models.GameState, error)
    UpsertGameStates(ctx context.Context, states []*models.GameState) error
}

// GameServiceAPI реализует grpc GameServiceServer
type GameServiceAPI struct {
    game_api.UnimplementedGameServiceServer
    gameService gameService
}

func NewGameServiceAPI(gameService gameService) *GameServiceAPI {
    return &GameServiceAPI{
        gameService: gameService,
    }
}