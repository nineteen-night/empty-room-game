package pgstorage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
)

type PGstorage struct {
	db *pgxpool.Pool
}

func NewPGStorage(connString string) (*PGstorage, error) {
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка парсинга конфига")
	}

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, errors.Wrap(err, "ошибка подключения")
	}
	
	storage := &PGstorage{
		db: db,
	}
	
	err = storage.initTables()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

func (s *PGstorage) initTables() error {
	usersSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		%s VARCHAR(50) UNIQUE NOT NULL,
		%s VARCHAR(100) UNIQUE NOT NULL,
		%s VARCHAR(255) NOT NULL,
		%s INTEGER DEFAULT 0 NOT NULL
	)`, usersTable, userIDColumn, usernameColumn, emailColumn, 
	   passwordHashColumn, maxRoomReachedColumn)
	
	_, err := s.db.Exec(context.Background(), usersSQL)
	if err != nil {
		return errors.Wrap(err, "создание таблицы users")
	}

    partnershipsSQL := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        %s UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        %s UUID NOT NULL REFERENCES %s(%s),
        %s UUID NOT NULL REFERENCES %s(%s),
        UNIQUE(%s, %s)
    )`, partnershipsTable, partnershipIDColumn, 
       user1IDColumn, usersTable, userIDColumn,
       user2IDColumn, usersTable, userIDColumn,
       user1IDColumn, user2IDColumn)
    
    _, err = s.db.Exec(context.Background(), partnershipsSQL)
    if err != nil {
        return errors.Wrap(err, "создание таблицы partnerships")
    }

	indexesSQL := fmt.Sprintf(`
	CREATE INDEX IF NOT EXISTS idx_users_username ON %s(%s);
	CREATE INDEX IF NOT EXISTS idx_users_email ON %s(%s);
	CREATE INDEX IF NOT EXISTS idx_partnerships_user1 ON %s(%s);
	CREATE INDEX IF NOT EXISTS idx_partnerships_user2 ON %s(%s);
	`, usersTable, usernameColumn,
	   usersTable, emailColumn,
	   partnershipsTable, user1IDColumn,
	   partnershipsTable, user2IDColumn)
	
	_, err = s.db.Exec(context.Background(), indexesSQL)
	if err != nil {
		return errors.Wrap(err, "создание индексов")
	}

	return nil
}

func (s *PGstorage) Close() {
	s.db.Close()
}