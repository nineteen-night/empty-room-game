package game_api

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
    proto_models "github.com/nineteen-night/empty-room-game/internal/game/pb/models"
    "github.com/nineteen-night/empty-room-game/internal/game/pb/game_api"
    "github.com/samber/lo"
)

func (s *GameServiceAPI) UpsertGameSessions(ctx context.Context, req *game_api.UpsertGameSessionsRequest) (*game_api.UpsertGameSessionsResponse, error) {
    sessions := mapGameSessionsFromProto(req.GameSessions)
    err := s.gameService.UpsertGameSessions(ctx, sessions)
    if err != nil {
        return &game_api.UpsertGameSessionsResponse{}, err
    }
    return &game_api.UpsertGameSessionsResponse{}, nil
}

func mapGameSessionsFromProto(sessions []*proto_models.GameSessionUpsertModel) []*models.GameSession {
    return lo.Map(sessions, func(s *proto_models.GameSessionUpsertModel, _ int) *models.GameSession {
        return &models.GameSession{
            PartnershipID:   s.PartnershipId,
            CurrentRoom:     s.CurrentRoom,
            Status:          s.Status,
            CurrentPlayerID: s.CurrentPlayerId,
        }
    })
}