package game_api

import (
    "context"
    "log"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
    proto_models "github.com/nineteen-night/empty-room-game/internal/game/pb/models"
    "github.com/nineteen-night/empty-room-game/internal/game/pb/game_api"
    "github.com/samber/lo"
)

func (s *GameServiceAPI) GetGameStatesByIDs(ctx context.Context, req *game_api.GetGameStatesByIDsRequest) (*game_api.GetGameStatesByIDsResponse, error) {
    log.Printf("Received GetGameStatesByIDs request with IDs: %v", req.Ids)

    states, err := s.gameService.GetGameStatesByIDs(ctx, req.Ids)
    if err != nil {
        return &game_api.GetGameStatesByIDsResponse{}, err
    }

    return &game_api.GetGameStatesByIDsResponse{
        GameStates: mapGameStatesToProto(states),
    }, nil
}

func mapGameStatesToProto(states []*models.GameState) []*proto_models.GameStateModel {
    return lo.Map(states, func(s *models.GameState, _ int) *proto_models.GameStateModel {
        return &proto_models.GameStateModel{
            Id:            s.ID,
            GameSessionId: s.GameSessionID,
            Inventory:     s.Inventory,
            PuzzlesSolved: s.PuzzlesSolved,
            CurrentRoomId: s.CurrentRoomID,
        }
    })
}