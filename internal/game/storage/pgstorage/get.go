package pgstorage

import (
    "context"
    "fmt"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
    "github.com/Masterminds/squirrel"
    "github.com/pkg/errors"
)

func (s *PGStorage) GetGameSession(ctx context.Context, partnershipID string) (*models.GameSession, error) {
    schema := s.getSchema(partnershipID)
    
    query := squirrel.Select(
        "partnership_id", "user1_id", "user2_id", "current_room",
    ).
        From(fmt.Sprintf("%s.game_sessions", schema)).
        Where(squirrel.Eq{"partnership_id": partnershipID}).
        PlaceholderFormat(squirrel.Dollar)

    queryText, args, err := query.ToSql()
    if err != nil {
        return nil, errors.Wrap(err, "generate query error")
    }

    var session models.GameSession
    err = s.db.QueryRow(ctx, queryText, args...).Scan(
        &session.PartnershipID,
        &session.User1ID,
        &session.User2ID,
        &session.CurrentRoom,
    )

    if err != nil {
        if err.Error() == "no rows in result set" {
            return nil, nil
        }
        return nil, errors.Wrap(err, "query error")
    }

    return &session, nil
}

func (s *PGStorage) GetRoomByNumber(ctx context.Context, roomNumber int32) (*models.Room, error) {
    query := squirrel.Select(
        "room_number", "name", "description",
    ).
        From("public.rooms").
        Where(squirrel.Eq{"room_number": roomNumber}).
        PlaceholderFormat(squirrel.Dollar)

    queryText, args, err := query.ToSql()
    if err != nil {
        return nil, errors.Wrap(err, "generate query error")
    }

    var room models.Room
    err = s.db.QueryRow(ctx, queryText, args...).Scan(
        &room.RoomNumber,
        &room.Name,
        &room.Description,
    )

    if err != nil {
        if err.Error() == "no rows in result set" {
            return nil, nil
        }
        return nil, errors.Wrap(err, "query error")
    }

    return &room, nil
}

func (s *PGStorage) GetGameSessionsByUserID(ctx context.Context, userID string) ([]*models.GameSession, error) {
    var allSessions []*models.GameSession

    for i := 0; i < s.shardCount; i++ {
        schema := fmt.Sprintf("shard_%d", i)
        
        query := squirrel.Select(
            "partnership_id", "user1_id", "user2_id", "current_room",
        ).
            From(fmt.Sprintf("%s.game_sessions", schema)).
            Where(
                squirrel.Or{
                    squirrel.Eq{"user1_id": userID},
                    squirrel.Eq{"user2_id": userID},
                },
            ).
            PlaceholderFormat(squirrel.Dollar)

        queryText, args, err := query.ToSql()
        if err != nil {
            continue
        }

        rows, err := s.db.Query(ctx, queryText, args...)
        if err != nil {
            continue
        }
        
        for rows.Next() {
            var session models.GameSession
            err := rows.Scan(
                &session.PartnershipID,
                &session.User1ID,
                &session.User2ID,
                &session.CurrentRoom,
            )
            if err != nil {
                rows.Close()
                continue
            }
            allSessions = append(allSessions, &session)
        }
        rows.Close()
    }

    return allSessions, nil
}