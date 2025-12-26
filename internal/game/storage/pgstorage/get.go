package pgstorage

import (
	"context"

	"github.com/nineteen-night/empty-room-game/internal/game/models"
	"github.com/Masterminds/squirrel"
	"github.com/pkg/errors"
)

func (s *PGStorage) GetGameSession(ctx context.Context, partnershipID string) (*models.GameSession, error) {
	query := squirrel.Select(
		partnershipIDColumn,
		user1IDColumn,
		user2IDColumn,
		currentRoomColumn,
	).
		From(gameSessionsTable).
		Where(squirrel.Eq{partnershipIDColumn: partnershipID}).
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
		roomNumberColumn,
		nameColumn,
		descriptionColumn,
	).
		From(roomsTable).
		Where(squirrel.Eq{roomNumberColumn: roomNumber}).
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
	query := squirrel.Select(
		partnershipIDColumn,
		user1IDColumn,
		user2IDColumn,
		currentRoomColumn,
	).
		From(gameSessionsTable).
		Where(
			squirrel.Or{
				squirrel.Eq{user1IDColumn: userID},
				squirrel.Eq{user2IDColumn: userID},
			},
		).
		PlaceholderFormat(squirrel.Dollar)

	queryText, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "generate query error")
	}

	rows, err := s.db.Query(ctx, queryText, args...)
	if err != nil {
		return nil, errors.Wrap(err, "query error")
	}
	defer rows.Close()

	var sessions []*models.GameSession
	for rows.Next() {
		var session models.GameSession
		err := rows.Scan(
			&session.PartnershipID,
			&session.User1ID,
			&session.User2ID,
			&session.CurrentRoom,
		)
		if err != nil {
			return nil, errors.Wrap(err, "scan error")
		}
		sessions = append(sessions, &session)
	}

	return sessions, nil
}