package pgstorage

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
    "github.com/Masterminds/squirrel"
    "github.com/pkg/errors"
)

func (storage *PGstorage) GetGameSessionsByIDs(ctx context.Context, IDs []uint64) ([]*models.GameSession, error) {
    if len(IDs) == 0 {
        return []*models.GameSession{}, nil
    }
    
    query := storage.getGameSessionsQuery(IDs)
    queryText, args, err := query.ToSql()
    if err != nil {
        return nil, errors.Wrap(err, "generate query error")
    }
    rows, err := storage.db.Query(ctx, queryText, args...)
    if err != nil {
        return nil, errors.Wrap(err, "query error")
    }
    defer rows.Close()
    
    var sessions []*models.GameSession
    for rows.Next() {
        var session models.GameSession
        if err := rows.Scan(&session.ID, &session.PartnershipID, &session.CurrentRoom, &session.Status, &session.CurrentPlayerID); err != nil {
            return nil, errors.Wrap(err, "failed to scan row")
        }
        sessions = append(sessions, &session)
    }
    return sessions, nil
}

func (storage *PGstorage) getGameSessionsQuery(IDs []uint64) squirrel.Sqlizer {
    q := squirrel.Select(
        gameSessionIDColumn,
        partnershipIDColumn,
        currentRoomColumn,
        statusColumn,
        currentPlayerIDColumn,
    ).
        From(gameSessionsTable).
        Where(squirrel.Eq{gameSessionIDColumn: IDs}).
        PlaceholderFormat(squirrel.Dollar)
    return q
}

func (storage *PGstorage) GetGameStatesByIDs(ctx context.Context, IDs []uint64) ([]*models.GameState, error) {
    if len(IDs) == 0 {
        return []*models.GameState{}, nil
    }
    
    query := storage.getGameStatesQuery(IDs)
    queryText, args, err := query.ToSql()
    if err != nil {
        return nil, errors.Wrap(err, "generate query error")
    }
    rows, err := storage.db.Query(ctx, queryText, args...)
    if err != nil {
        return nil, errors.Wrap(err, "query error")
    }
    defer rows.Close()
    
    var states []*models.GameState
    for rows.Next() {
        var state models.GameState
        if err := rows.Scan(&state.ID, &state.GameSessionID, &state.Inventory, &state.PuzzlesSolved, &state.CurrentRoomID); err != nil {
            return nil, errors.Wrap(err, "failed to scan row")
        }
        states = append(states, &state)
    }
    return states, nil
}

func (storage *PGstorage) getGameStatesQuery(IDs []uint64) squirrel.Sqlizer {
    q := squirrel.Select(
        gameStateIDColumn,
        gameSessionIDColumn,
        inventoryColumn,
        puzzlesSolvedColumn,
        currentRoomIDColumn,
    ).
        From(gameStatesTable).
        Where(squirrel.Eq{gameStateIDColumn: IDs}).
        PlaceholderFormat(squirrel.Dollar)
    return q
}