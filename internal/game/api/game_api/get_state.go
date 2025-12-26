package game_api

import (
	"context"

	"github.com/nineteen-night/empty-room-game/internal/game/pb/game_api"
	proto_models "github.com/nineteen-night/empty-room-game/internal/game/pb/models"
)

func (s *GameServiceAPI) GetGameState(ctx context.Context, req *game_api.GetGameStateRequest) (*game_api.GetGameStateResponse, error) {
	gameState, err := s.gameService.GetGameState(ctx, req.PartnershipId)
	if err != nil {
		return nil, err
	}

	return &game_api.GetGameStateResponse{
		GameState: &proto_models.GameState{
			CurrentRoom: gameState.CurrentRoom,
			RoomInfo: &proto_models.Room{
				RoomNumber:  gameState.RoomInfo.RoomNumber,
				Name:        gameState.RoomInfo.Name,
				Description: gameState.RoomInfo.Description,
			},
		},
	}, nil
}