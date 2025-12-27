package pgstorage

import (
    "context"
    "fmt"
    "hash/fnv"
    "sync"
    "log/slog"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/pkg/errors"
)

type PGStorage struct {
    shards []*pgxpool.Pool
    mu     sync.RWMutex
}

func NewPGStorage(shardConfigs []struct {
    Host, Username, Password, DBName, SSLMode string
    Port int
}) (*PGStorage, error) {
    storage := &PGStorage{
        shards: make([]*pgxpool.Pool, len(shardConfigs)),
    }

    for i, cfg := range shardConfigs {
        connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
            cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)
        
        config, err := pgxpool.ParseConfig(connStr)
        if err != nil {
            return nil, errors.Wrap(err, "ошибка парсинга конфига")
        }

        db, err := pgxpool.NewWithConfig(context.Background(), config)
        if err != nil {
            return nil, errors.Wrap(err, "ошибка подключения к шарду")
        }
        
        storage.shards[i] = db
        
        if err := storage.initShardTables(i); err != nil {
            return nil, fmt.Errorf("ошибка инициализации шарда %d: %w", i, err)
        }
        
        slog.Info("Шард инициализирован", "шард", i, "база", cfg.DBName)
    }

    return storage, nil
}

func (s *PGStorage) getShard(partnershipID string) *pgxpool.Pool {
    if partnershipID == "" {
        s.mu.RLock()
        defer s.mu.RUnlock()
        return s.shards[0]
    }
    
    h := fnv.New32a()
    h.Write([]byte(partnershipID))
    idx := int(h.Sum32() % uint32(len(s.shards)))
    
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if idx < len(s.shards) {
        return s.shards[idx]
    }
    return s.shards[0]
}

func (s *PGStorage) initShardTables(shardIndex int) error {
    s.mu.RLock()
    db := s.shards[shardIndex]
    s.mu.RUnlock()
    
    if db == nil {
        return fmt.Errorf("нет соединения с шардом %d", shardIndex)
    }

    _, err := db.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS rooms (
            room_number INTEGER PRIMARY KEY,
            name VARCHAR(50) NOT NULL,
            description TEXT NOT NULL
        )
    `)
    if err != nil {
        return errors.Wrap(err, "создание таблицы rooms")
    }

    _, err = db.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS game_sessions (
            partnership_id UUID PRIMARY KEY,
            user1_id UUID NOT NULL,
            user2_id UUID NOT NULL,
            current_room INTEGER DEFAULT 1 NOT NULL
        )
    `)
    if err != nil {
        return errors.Wrap(err, "создание таблицы game_sessions")
    }

    _, err = db.Exec(context.Background(), `
        CREATE INDEX IF NOT EXISTS idx_game_sessions_user1 ON game_sessions(user1_id);
        CREATE INDEX IF NOT EXISTS idx_game_sessions_user2 ON game_sessions(user2_id);
    `)
    if err != nil {
        return errors.Wrap(err, "создание индексов")
    }

    return s.initRoomsForShard(shardIndex)
}

func (s *PGStorage) initRoomsForShard(shardIndex int) error {
    s.mu.RLock()
    db := s.shards[shardIndex]
    s.mu.RUnlock()

    var count int
    err := db.QueryRow(context.Background(), 
        "SELECT COUNT(*) FROM rooms").Scan(&count)
    if err != nil {
        count = 0
    }

    if count > 0 {
        return nil
    }

    rooms := []struct {
        number      int32
        name        string
        description string
    }{
        {1, "Прихожая", "Вы в тёмной прихожей старого дома..."},
        {2, "Кабинет", "Старый кабинет с пыльными книгами..."},
        {3, "Библиотека", "Пыльная библиотека с древними фолиантами..."},
        {4, "Гостиная", "Простораная гостиная с камином..."},
        {5, "Спальня", "Старинная спальня с кованой кроватью..."},
    }

    for _, room := range rooms {
        _, err := db.Exec(context.Background(), 
            "INSERT INTO rooms (room_number, name, description) VALUES ($1, $2, $3) ON CONFLICT (room_number) DO NOTHING",
            room.number, room.name, room.description)
        if err != nil {
            return errors.Wrapf(err, "ошибка добавления комнаты %d", room.number)
        }
    }

    return nil
}

func (s *PGStorage) getFirstShard() *pgxpool.Pool {
    s.mu.RLock()
    defer s.mu.RUnlock()
    
    if len(s.shards) > 0 {
        return s.shards[0]
    }
    return nil
}

func (s *PGStorage) Close() {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    for i, db := range s.shards {
        if db != nil {
            db.Close()
            s.shards[i] = nil
        }
    }
}