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
    gameSessionsSQL := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        %s SERIAL PRIMARY KEY,
        %s BIGINT NOT NULL UNIQUE,
        %s VARCHAR(100) NOT NULL DEFAULT 'start_room',
        %s VARCHAR(20) DEFAULT 'active' CHECK (%s IN ('active', 'paused', 'completed')),
        %s BIGINT NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`, gameSessionsTable, gameSessionIDColumn, 
       partnershipIDColumn, currentRoomColumn,
       statusColumn, statusColumn, currentPlayerIDColumn)
    
    _, err := s.db.Exec(context.Background(), gameSessionsSQL)
    if err != nil {
        return errors.Wrap(err, "создание таблицы game_sessions")
    }

    gameStatesSQL := fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS %s (
        %s SERIAL PRIMARY KEY,
        %s BIGINT NOT NULL UNIQUE REFERENCES %s(%s),
        %s TEXT DEFAULT '{}',
        %s TEXT DEFAULT '[]',
        %s BIGINT DEFAULT 1,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`, gameStatesTable, gameStateIDColumn,
       gameSessionIDColumn2, gameSessionsTable, gameSessionIDColumn,
       inventoryColumn, puzzlesSolvedColumn, currentRoomIDColumn)
    
    _, err = s.db.Exec(context.Background(), gameStatesSQL)
    if err != nil {
        return errors.Wrap(err, "создание таблицы game_states")
    }

    indexesSQL := fmt.Sprintf(`
    CREATE INDEX IF NOT EXISTS idx_game_sessions_partnership ON %s(%s);
    CREATE INDEX IF NOT EXISTS idx_game_sessions_status ON %s(%s);
    CREATE INDEX IF NOT EXISTS idx_game_states_session ON %s(%s);
    `, gameSessionsTable, partnershipIDColumn,
       gameSessionsTable, statusColumn,
       gameStatesTable, gameSessionIDColumn2)
    
    _, err = s.db.Exec(context.Background(), indexesSQL)
    if err != nil {
        return errors.Wrap(err, "создание индексов")
    }

    return nil
}

func (s *PGstorage) Close() {
    s.db.Close()
}