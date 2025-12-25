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
	//users
	usersSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s SERIAL PRIMARY KEY,
		%s VARCHAR(50) UNIQUE NOT NULL,
		%s VARCHAR(100) UNIQUE NOT NULL,
		%s VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`, usersTable, userIDColumn, usernameColumn, emailColumn, passwordHashColumn)
	
	_, err := s.db.Exec(context.Background(), usersSQL)
	if err != nil {
		return errors.Wrap(err, "создание таблицы users")
	}

	//partnerships
	partnershipsSQL := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS %s (
		%s SERIAL PRIMARY KEY,
		%s BIGINT NOT NULL REFERENCES %s(%s),
		%s BIGINT NOT NULL REFERENCES %s(%s),
		%s VARCHAR(20) DEFAULT 'pending' CHECK (%s IN ('pending', 'active', 'completed')),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(%s, %s)
	)`, partnershipsTable, partnershipIDColumn, 
	   player1IDColumn, usersTable, userIDColumn,
	   player2IDColumn, usersTable, userIDColumn,
	   statusColumn, statusColumn,
	   player1IDColumn, player2IDColumn)
	
	_, err = s.db.Exec(context.Background(), partnershipsSQL)
	if err != nil {
		return errors.Wrap(err, "создание таблицы partnerships")
	}

	//индексы
	indexesSQL := fmt.Sprintf(`
	CREATE INDEX IF NOT EXISTS idx_users_username ON %s(%s);
	CREATE INDEX IF NOT EXISTS idx_users_email ON %s(%s);
	CREATE INDEX IF NOT EXISTS idx_partnerships_player1 ON %s(%s);
	CREATE INDEX IF NOT EXISTS idx_partnerships_player2 ON %s(%s);
	`, usersTable, usernameColumn,
	   usersTable, emailColumn,
	   partnershipsTable, player1IDColumn,
	   partnershipsTable, player2IDColumn)
	
	_, err = s.db.Exec(context.Background(), indexesSQL)
	if err != nil {
		return errors.Wrap(err, "создание индексов")
	}

	return nil
}

func (s *PGstorage) Close() {
	s.db.Close()
}