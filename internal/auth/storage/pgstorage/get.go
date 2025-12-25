package pgstorage

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    "github.com/Masterminds/squirrel"
    "github.com/pkg/errors"
)

func (storage *PGstorage) GetUsersByIDs(ctx context.Context, IDs []uint64) ([]*models.User, error) {
    if len(IDs) == 0 {
        return []*models.User{}, nil  // Вернул проверку!
    }
    
    query := storage.getUsersQuery(IDs)
    queryText, args, err := query.ToSql()
    if err != nil {
        return nil, errors.Wrap(err, "generate query error")
    }
    rows, err := storage.db.Query(ctx, queryText, args...)
    if err != nil {
        return nil, errors.Wrap(err, "query error")
    }
    defer rows.Close()
    
    var users []*models.User
    for rows.Next() {
        var user models.User
        if err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash); err != nil {
            return nil, errors.Wrap(err, "failed to scan row")
        }
        users = append(users, &user)
    }
    return users, nil
}

func (storage *PGstorage) getUsersQuery(IDs []uint64) squirrel.Sqlizer {
    q := squirrel.Select(userIDColumn, usernameColumn, emailColumn, passwordHashColumn).
        From(usersTable).
        Where(squirrel.Eq{userIDColumn: IDs}).
        PlaceholderFormat(squirrel.Dollar)
    return q
}

func (storage *PGstorage) GetPartnershipsByIDs(ctx context.Context, IDs []uint64) ([]*models.Partnership, error) {
    if len(IDs) == 0 {
        return []*models.Partnership{}, nil  // Вернул проверку!
    }
    
    query := storage.getPartnershipsQuery(IDs)
    queryText, args, err := query.ToSql()
    if err != nil {
        return nil, errors.Wrap(err, "generate query error")
    }
    rows, err := storage.db.Query(ctx, queryText, args...)
    if err != nil {
        return nil, errors.Wrap(err, "query error")
    }
    defer rows.Close()
    
    var partnerships []*models.Partnership
    for rows.Next() {
        var partnership models.Partnership
        if err := rows.Scan(&partnership.ID, &partnership.Player1ID, &partnership.Player2ID, &partnership.Status); err != nil {
            return nil, errors.Wrap(err, "failed to scan row")
        }
        partnerships = append(partnerships, &partnership)
    }
    return partnerships, nil
}

func (storage *PGstorage) getPartnershipsQuery(IDs []uint64) squirrel.Sqlizer {
    q := squirrel.Select(partnershipIDColumn, player1IDColumn, player2IDColumn, statusColumn).
        From(partnershipsTable).
        Where(squirrel.Eq{partnershipIDColumn: IDs}).
        PlaceholderFormat(squirrel.Dollar)
    return q
}

func (storage *PGstorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
    query := storage.getUserByUsernameQuery(username)
    queryText, args, err := query.ToSql()
    if err != nil {
        return nil, errors.Wrap(err, "generate query error")
    }
    
    var user models.User
    err = storage.db.QueryRow(ctx, queryText, args...).Scan(
        &user.ID,
        &user.Username,
        &user.Email,
        &user.PasswordHash,
    )
    if err != nil {
        if err.Error() == "no rows in result set" {
            return nil, nil
        }
        return nil, errors.Wrap(err, "query error")
    }
    
    return &user, nil
}

func (storage *PGstorage) getUserByUsernameQuery(username string) squirrel.Sqlizer {
    q := squirrel.Select(userIDColumn, usernameColumn, emailColumn, passwordHashColumn).
        From(usersTable).
        Where(squirrel.Eq{usernameColumn: username}).
        PlaceholderFormat(squirrel.Dollar)
    return q
}