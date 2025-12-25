package game_api

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
    proto_models "github.com/nineteen-night/empty-room-game/internal/game/pb/models"
    "github.com/nineteen-night/empty-room-game/internal/game/pb/game_api"
    "github.com/samber/lo"
)

func (s *GameServiceAPI) UpsertGameStates(ctx context.Context, req *game_api.UpsertGameStatesRequest) (*game_api.UpsertGameStatesResponse, error) {
    states := mapGameStatesFromProto(req.GameStates)
    err := s.gameService.UpsertGameStates(ctx, states)
    if err != nil {
        return &game_api.UpsertGameStatesResponse{}, err
    }
    return &game_api.UpsertGameStatesResponse{}, nil
}

func mapGameStatesFromProto(states []*proto_models.GameStateUpsertModel) []*models.GameState {
    return lo.Map(states, func(s *proto_models.GameStateUpsertModel, _ int) *models.GameState {
        return &models.GameState{
            GameSessionID: s.GameSessionId,
            Inventory:     s.Inventory,
            PuzzlesSolved: s.PuzzlesSolved,
            CurrentRoomID: s.CurrentRoomId,
        }
    })
}