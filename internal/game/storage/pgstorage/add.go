package pgstorage

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (s *PGStorage) CreateGameSession(ctx context.Context, partnershipID, user1ID, user2ID string) error {
	query := squirrel.Insert(gameSessionsTable).
		Columns(partnershipIDColumn, user1IDColumn, user2IDColumn, currentRoomColumn).
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
	query := squirrel.Delete(gameSessionsTable).
		Where(squirrel.Eq{partnershipIDColumn: partnershipID}).
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
	query := squirrel.Update(gameSessionsTable).
		Set(currentRoomColumn, newRoom).
		Where(squirrel.Eq{partnershipIDColumn: partnershipID}).
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
	query := fmt.Sprintf("SELECT MAX(%s) FROM %s", roomNumberColumn, roomsTable)
	
	var maxRoom int32
	err := s.db.QueryRow(ctx, query).Scan(&maxRoom)
	if err != nil {
		return 0, errors.Wrap(err, "query error")
	}

	return maxRoom, nil
}