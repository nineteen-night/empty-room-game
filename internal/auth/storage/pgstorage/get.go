package pgstorage

import (
    "context"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    "github.com/Masterminds/squirrel"
    "github.com/pkg/errors"
)

func (storage *PGstorage) GetUserByID(ctx context.Context, userID string) (*models.User, error) {
    query := squirrel.Select(
        userIDColumn, 
        usernameColumn, 
        emailColumn, 
        passwordHashColumn,
        maxRoomReachedColumn,
    ).
    From(usersTable).
    Where(squirrel.Eq{userIDColumn: userID}).
    PlaceholderFormat(squirrel.Dollar)
    
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
        &user.MaxRoomReached, 
    )
    if err != nil {
        return nil, errors.Wrap(err, "query error")
    }
    
    return &user, nil
}

func (storage *PGstorage) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
    query := squirrel.Select(
        userIDColumn, 
        usernameColumn, 
        emailColumn, 
        passwordHashColumn,
        maxRoomReachedColumn,
    ).
    From(usersTable).
    Where(squirrel.Eq{usernameColumn: username}).
    PlaceholderFormat(squirrel.Dollar)
    
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
        &user.MaxRoomReached,
    )
    
    if err != nil {
        if err.Error() == "no rows in result set" {
            return nil, nil
        }
        return nil, errors.Wrap(err, "query error")
    }
    
    return &user, nil
}

func (storage *PGstorage) GetPartnershipByID(ctx context.Context, partnershipID string) (*models.Partnership, error) {
    query := squirrel.Select(
        partnershipIDColumn, 
        user1IDColumn,
        user2IDColumn,
    ).
    From(partnershipsTable).
    Where(squirrel.Eq{partnershipIDColumn: partnershipID}).
    PlaceholderFormat(squirrel.Dollar)
    
    queryText, args, err := query.ToSql()
    if err != nil {
        return nil, errors.Wrap(err, "generate query error")
    }
    
    var partnership models.Partnership
    err = storage.db.QueryRow(ctx, queryText, args...).Scan(
        &partnership.ID,
        &partnership.User1ID,
        &partnership.User2ID,
    )
    if err != nil {
        return nil, errors.Wrap(err, "query error")
    }
    
    return &partnership, nil
}

func (storage *PGstorage) GetPartnershipBetweenUsers(ctx context.Context, user1ID, user2ID string) (*models.Partnership, error) {
    query := squirrel.Select(
        partnershipIDColumn, 
        user1IDColumn,
        user2IDColumn,
    ).
    From(partnershipsTable).
    Where(
        squirrel.Or{
            squirrel.And{
                squirrel.Eq{user1IDColumn: user1ID},
                squirrel.Eq{user2IDColumn: user2ID},
            },
            squirrel.And{
                squirrel.Eq{user1IDColumn: user2ID},
                squirrel.Eq{user2IDColumn: user1ID},
            },
        },
    ).
    PlaceholderFormat(squirrel.Dollar)
    
    queryText, args, err := query.ToSql()
    if err != nil {
        return nil, errors.Wrap(err, "generate query error")
    }
    
    var partnership models.Partnership
    err = storage.db.QueryRow(ctx, queryText, args...).Scan(
        &partnership.ID,
        &partnership.User1ID,
        &partnership.User2ID,
    )
    
    if err != nil {
        if err.Error() == "no rows in result set" {
            return nil, nil
        }
        return nil, errors.Wrap(err, "query error")
    }
    
    return &partnership, nil
}