package pgstorage

import (
    "context"
    "fmt"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    "github.com/Masterminds/squirrel"
    "github.com/pkg/errors"
    "golang.org/x/crypto/bcrypt"
    "github.com/google/uuid"
)

func (storage *PGstorage) CreateUser(ctx context.Context, user *models.User) (string, error) {
    userID := uuid.New().String()
    
    query := squirrel.Insert(usersTable).
        Columns(userIDColumn, usernameColumn, emailColumn, passwordHashColumn, maxRoomReachedColumn).
        Values(userID, user.Username, user.Email, user.PasswordHash, 0).
        PlaceholderFormat(squirrel.Dollar)
    
    queryText, args, err := query.ToSql()
    if err != nil {
        return "", errors.Wrap(err, "generate query error")
    }
    
    _, err = storage.db.Exec(ctx, queryText, args...)
    if err != nil {
        return "", errors.Wrap(err, "execute query error")
    }
    
    return userID, nil
}

func (storage *PGstorage) CreatePartnership(ctx context.Context, user1ID, user2ID string) (string, error) {
    partnershipID := uuid.New().String()
    
    query := squirrel.Insert(partnershipsTable).
        Columns(partnershipIDColumn, user1IDColumn, user2IDColumn, statusColumn).
        Values(partnershipID, user1ID, user2ID, "active").
        PlaceholderFormat(squirrel.Dollar)
    
    queryText, args, err := query.ToSql()
    if err != nil {
        return "", errors.Wrap(err, "generate query error")
    }
    
    _, err = storage.db.Exec(ctx, queryText, args...)
    if err != nil {
        return "", errors.Wrap(err, "execute query error")
    }
    
    return partnershipID, nil
}

func (storage *PGstorage) TerminatePartnership(ctx context.Context, partnershipID string) error {
    query := squirrel.Delete(partnershipsTable).
        Where(squirrel.Eq{partnershipIDColumn: partnershipID}).
        PlaceholderFormat(squirrel.Dollar)
    
    queryText, args, err := query.ToSql()
    if err != nil {
        return errors.Wrap(err, "generate query error")
    }
    
    result, err := storage.db.Exec(ctx, queryText, args...)
    if err != nil {
        return errors.Wrap(err, "execute query error")
    }
    
    if result.RowsAffected() == 0 {
        return errors.New("partnership not found")
    }
    
    return nil
}

func (storage *PGstorage) UpdateUserMaxRoom(ctx context.Context, userID string, roomNumber int32) (bool, error) {
    var currentMax int32
    err := storage.db.QueryRow(ctx, 
        fmt.Sprintf("SELECT %s FROM %s WHERE %s = $1", 
            maxRoomReachedColumn, usersTable, userIDColumn),
        userID).Scan(&currentMax)
    
    if err != nil {
        return false, errors.Wrap(err, "get current max_room error")
    }
    
    if roomNumber > currentMax {
        query := squirrel.Update(usersTable).
            Set(maxRoomReachedColumn, roomNumber).
            Where(squirrel.Eq{userIDColumn: userID}).
            PlaceholderFormat(squirrel.Dollar)
        
        queryText, args, err := query.ToSql()
        if err != nil {
            return false, errors.Wrap(err, "generate query error")
        }
        
        _, err = storage.db.Exec(ctx, queryText, args...)
        if err != nil {
            return false, errors.Wrap(err, "execute query error")
        }
        
        return true, nil 
    }
    
    return false, nil 
}

func (storage *PGstorage) HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(hash), nil
}

func (storage *PGstorage) CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}