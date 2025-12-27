package pgstorage

import (
    "context"
    "fmt"

    "github.com/Masterminds/squirrel"
    "github.com/pkg/errors"
)

func (s *PGStorage) CreateGameSession(ctx context.Context, partnershipID, user1ID, user2ID string) error {
    schema := s.getSchema(partnershipID)
    
    query := squirrel.Insert(fmt.Sprintf("%s.game_sessions", schema)).
        Columns("partnership_id", "user1_id", "user2_id", "current_room").
        Values(partnershipID, user1ID, user2ID, 1).
        PlaceholderFormat(squirrel.Dollar)

    queryText, args, err := query.ToSql()
    if err != nil {
        return errors.Wrap(err, "generate query error")
    }

    _, err = s.db.Exec(ctx, queryText, args...)
    if err != nil {
        return errors.Wrap(err, "execute query error")
    }

    return nil
}

func (s *PGStorage) DeleteGameSession(ctx context.Context, partnershipID string) error {
    schema := s.getSchema(partnershipID)
    
    query := squirrel.Delete(fmt.Sprintf("%s.game_sessions", schema)).
        Where(squirrel.Eq{"partnership_id": partnershipID}).
        PlaceholderFormat(squirrel.Dollar)

    queryText, args, err := query.ToSql()
    if err != nil {
        return errors.Wrap(err, "generate query error")
    }

    result, err := s.db.Exec(ctx, queryText, args...)
    if err != nil {
        return errors.Wrap(err, "execute query error")
    }

    if result.RowsAffected() == 0 {
        return errors.New("game session not found")
    }

    return nil
}


func (s *PGStorage) UpdateCurrentRoom(ctx context.Context, partnershipID string, newRoom int32) error {
    schema := s.getSchema(partnershipID)
    
    query := squirrel.Update(fmt.Sprintf("%s.game_sessions", schema)).
        Set("current_room", newRoom).
        Where(squirrel.Eq{"partnership_id": partnershipID}).
        PlaceholderFormat(squirrel.Dollar)

    queryText, args, err := query.ToSql()
    if err != nil {
        return errors.Wrap(err, "generate query error")
    }

    result, err := s.db.Exec(ctx, queryText, args...)
    if err != nil {
        return errors.Wrap(err, "execute query error")
    }

    if result.RowsAffected() == 0 {
        return errors.New("game session not found")
    }

    return nil
}

func (s *PGStorage) GetMaxRoomNumber(ctx context.Context) (int32, error) {
    query := "SELECT MAX(room_number) FROM public.rooms"
    
    var maxRoom int32
    err := s.db.QueryRow(ctx, query).Scan(&maxRoom)
    if err != nil {
        return 0, errors.Wrap(err, "query error")
    }

    return maxRoom, nil
}