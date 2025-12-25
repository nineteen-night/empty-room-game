package pgstorage

import (
    "context"
    "fmt"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    "github.com/Masterminds/squirrel"
    "github.com/pkg/errors"
    "github.com/samber/lo"
    "golang.org/x/crypto/bcrypt"
)

func (storage *PGstorage) UpsertUsers(ctx context.Context, domainUsers []*models.User) error {
    query := storage.upsertUsersQuery(domainUsers)
    queryText, args, err := query.ToSql()
    if err != nil {
        return errors.Wrap(err, "generate query error")
    }
    _, err = storage.db.Exec(ctx, queryText, args...)
    return errors.Wrap(err, "execute query error")
}

func (storage *PGstorage) upsertUsersQuery(domainUsers []*models.User) squirrel.Sqlizer {
    users := lo.Map(domainUsers, func(user *models.User, _ int) *User {
        return &User{
            Username:     user.Username,
            Email:        user.Email,
            PasswordHash: user.PasswordHash,
        }
    })

    q := squirrel.Insert(usersTable).Columns(usernameColumn, emailColumn, passwordHashColumn).
        PlaceholderFormat(squirrel.Dollar).
        Suffix("ON CONFLICT (username) DO UPDATE SET email = EXCLUDED.email, password_hash = EXCLUDED.password_hash RETURNING id")
    
    for _, user := range users {
        q = q.Values(user.Username, user.Email, user.PasswordHash)
    }
    return q
}

func (storage *PGstorage) UpsertPartnerships(ctx context.Context, domainPartnerships []*models.Partnership) error {
    query := storage.upsertPartnershipsQuery(domainPartnerships)
    queryText, args, err := query.ToSql()
    if err != nil {
        return errors.Wrap(err, "generate query error")
    }
    _, err = storage.db.Exec(ctx, queryText, args...)
    return errors.Wrap(err, "execute query error")
}

func (storage *PGstorage) upsertPartnershipsQuery(domainPartnerships []*models.Partnership) squirrel.Sqlizer {
    partnerships := lo.Map(domainPartnerships, func(partnership *models.Partnership, _ int) *Partnership {
        return &Partnership{
            Player1ID: partnership.Player1ID,
            Player2ID: partnership.Player2ID,
            Status:    partnership.Status,
        }
    })

    q := squirrel.Insert(partnershipsTable).Columns(player1IDColumn, player2IDColumn, statusColumn).
        PlaceholderFormat(squirrel.Dollar).
        Suffix("ON CONFLICT (player1_id, player2_id) DO UPDATE SET status = EXCLUDED.status RETURNING id")
    
    for _, partnership := range partnerships {
        q = q.Values(partnership.Player1ID, partnership.Player2ID, partnership.Status)
    }
    return q
}

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("failed to hash password: %w", err)
    }
    return string(hash), nil
}

func CheckPasswordHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}