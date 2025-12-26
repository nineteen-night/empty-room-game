package game_api

import (
	"context"

	"github.com/nineteen-night/empty-room-game/internal/game/pb/game_api"
	proto_models "github.com/nineteen-night/empty-room-game/internal/game/pb/models"
)

func (s *GameServiceAPI) MoveToNextRoom(ctx context.Context, req *game_api.MoveToNextRoomRequest) (*game_api.MoveToNextRoomResponse, error) {
	newGameState, err := s.gameService.MoveToNextRoom(ctx, req.PartnershipId)
	if err != nil {
		return nil, err
	}

	return &game_api.MoveToNextRoomResponse{
		NewGameState: &proto_models.GameState{
			CurrentRoom: newGameState.CurrentRoom,
			RoomInfo: &proto_models.Room{
				RoomNumber:  newGameState.RoomInfo.RoomNumber,
				Name:        newGameState.RoomInfo.Name,
				Description: newGameState.RoomInfo.Description,
			},
		},
	}, nil
}