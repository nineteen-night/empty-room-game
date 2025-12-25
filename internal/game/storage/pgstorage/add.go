package pgstorage

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
    "github.com/Masterminds/squirrel"
    "github.com/pkg/errors"
    "github.com/samber/lo"
)

func (storage *PGstorage) UpsertGameSessions(ctx context.Context, domainSessions []*models.GameSession) error {
    if len(domainSessions) == 0 {
        return nil
    }

    query := storage.upsertGameSessionsQuery(domainSessions)
    queryText, args, err := query.ToSql()
    if err != nil {
        return errors.Wrap(err, "generate query error")
    }
    
    _, err = storage.db.Exec(ctx, queryText, args...)
    return errors.Wrap(err, "execute query error")
}

func (storage *PGstorage) upsertGameSessionsQuery(domainSessions []*models.GameSession) squirrel.Sqlizer {
    sessions := lo.Map(domainSessions, func(session *models.GameSession, _ int) *GameSession {
        return &GameSession{
            PartnershipID:   session.PartnershipID,
            CurrentRoom:     session.CurrentRoom,
            Status:          session.Status,
            CurrentPlayerID: session.CurrentPlayerID,
        }
    })

    q := squirrel.Insert(gameSessionsTable).
        Columns(partnershipIDColumn, currentRoomColumn, statusColumn, currentPlayerIDColumn).
        PlaceholderFormat(squirrel.Dollar).
        Suffix("ON CONFLICT (partnership_id) DO UPDATE SET current_room = EXCLUDED.current_room, status = EXCLUDED.status, current_player_id = EXCLUDED.current_player_id RETURNING id")
    
    for _, session := range sessions {
        q = q.Values(session.PartnershipID, session.CurrentRoom, session.Status, session.CurrentPlayerID)
    }
    return q
}

func (storage *PGstorage) UpsertGameStates(ctx context.Context, domainStates []*models.GameState) error {
    if len(domainStates) == 0 {
        return nil
    }

    query := storage.upsertGameStatesQuery(domainStates)
    queryText, args, err := query.ToSql()
    if err != nil {
        return errors.Wrap(err, "generate query error")
    }
    
    _, err = storage.db.Exec(ctx, queryText, args...)
    return errors.Wrap(err, "execute query error")
}

func (storage *PGstorage) upsertGameStatesQuery(domainStates []*models.GameState) squirrel.Sqlizer {
    states := lo.Map(domainStates, func(state *models.GameState, _ int) *GameState {
        return &GameState{
            GameSessionID: state.GameSessionID,
            Inventory:     state.Inventory,
            PuzzlesSolved: state.PuzzlesSolved,
            CurrentRoomID: state.CurrentRoomID,
        }
    })

    q := squirrel.Insert(gameStatesTable).
        Columns(gameSessionIDColumn2, inventoryColumn, puzzlesSolvedColumn, currentRoomIDColumn).
        PlaceholderFormat(squirrel.Dollar).
        Suffix("ON CONFLICT (game_session_id) DO UPDATE SET inventory = EXCLUDED.inventory, puzzles_solved = EXCLUDED.puzzles_solved, current_room_id = EXCLUDED.current_room_id RETURNING id")
    
    for _, state := range states {
        q = q.Values(state.GameSessionID, state.Inventory, state.PuzzlesSolved, state.CurrentRoomID)
    }
    return q
}