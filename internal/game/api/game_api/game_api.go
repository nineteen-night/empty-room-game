package game_api

import (
	"context"

	"github.com/nineteen-night/empty-room-game/internal/game/models"
	"github.com/nineteen-night/empty-room-game/internal/game/pb/game_api"
)

type gameService interface {
	GetGameState(ctx context.Context, partnershipID string) (*models.GameState, error)
	MoveToNextRoom(ctx context.Context, partnershipID string) (*models.GameState, error)
}

type GameServiceAPI struct {
	game_api.UnimplementedGameServiceServer
	gameService gameService
}

func NewGameServiceAPI(gameService gameService) *GameServiceAPI {
	return &GameServiceAPI{
		gameService: gameService,
	}
}