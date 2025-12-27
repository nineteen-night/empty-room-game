package pgstorage

import (
    "context"
    "fmt"
    "hash/fnv"
    "sync"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/pkg/errors"
)

type PGStorage struct {
    db        *pgxpool.Pool 
    shardCount int        
    mu        sync.RWMutex
}


func NewPGStorage(connString string, shardCount int) (*PGStorage, error) {
    config, err := pgxpool.ParseConfig(connString)
    if err != nil {
        return nil, errors.Wrap(err, "ошибка парсинга конфига")
    }

    db, err := pgxpool.NewWithConfig(context.Background(), config)
    if err != nil {
        return nil, errors.Wrap(err, "ошибка подключения")
    }
    
    storage := &PGStorage{
        db:        db,
        shardCount: shardCount,
    }

    if err := storage.initSchemas(); err != nil {
        db.Close()
        return nil, fmt.Errorf("ошибка инициализации схем: %w", err)
    }

    return storage, nil
}

func (s *PGStorage) getSchema(partnershipID string) string {
    if partnershipID == "" || s.shardCount == 0 {
        return "shard_0"
    }
    
    h := fnv.New32a()
    h.Write([]byte(partnershipID))
    hash := h.Sum32()

    shardNum := int(hash % uint32(s.shardCount))
    return fmt.Sprintf("shard_%d", shardNum)
}

func (s *PGStorage) initSchemas() error {
    for i := 0; i < s.shardCount; i++ {
        schemaName := fmt.Sprintf("shard_%d", i)

        _, err := s.db.Exec(context.Background(), 
            fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schemaName))
        if err != nil {
            return errors.Wrapf(err, "ошибка создания схемы %s", schemaName)
        }
        
        createTableSQL := fmt.Sprintf(`
            CREATE TABLE IF NOT EXISTS %s.game_sessions (
                partnership_id UUID PRIMARY KEY,
                user1_id UUID NOT NULL,
                user2_id UUID NOT NULL,
                current_room INTEGER DEFAULT 1 NOT NULL
            )`, schemaName)
        
        _, err = s.db.Exec(context.Background(), createTableSQL)
        if err != nil {
            return errors.Wrapf(err, "ошибка создания таблицы в схеме %s", schemaName)
        }
        

        indexSQL := fmt.Sprintf(`
            CREATE INDEX IF NOT EXISTS idx_%s_game_sessions_user1 
            ON %s.game_sessions(user1_id);
            CREATE INDEX IF NOT EXISTS idx_%s_game_sessions_user2 
            ON %s.game_sessions(user2_id);
        `, schemaName, schemaName, schemaName, schemaName)
        
        _, err = s.db.Exec(context.Background(), indexSQL)
        if err != nil {
            return errors.Wrapf(err, "ошибка создания индексов в схеме %s", schemaName)
        }
    }
    

    roomsSQL := `
        CREATE TABLE IF NOT EXISTS public.rooms (
            room_number INTEGER PRIMARY KEY,
            name VARCHAR(50) NOT NULL,
            description TEXT NOT NULL
        )`
    
    _, err := s.db.Exec(context.Background(), roomsSQL)
    if err != nil {
        return errors.Wrap(err, "ошибка создания таблицы rooms")
    }

    return s.initRooms()
}

func (s *PGStorage) initRooms() error {

    var count int
    err := s.db.QueryRow(context.Background(), 
        "SELECT COUNT(*) FROM public.rooms").Scan(&count)
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
        _, err := s.db.Exec(context.Background(), 
            `INSERT INTO public.rooms (room_number, name, description) 
             VALUES ($1, $2, $3) 
             ON CONFLICT (room_number) DO NOTHING`,
            room.number, room.name, room.description)
        if err != nil {
            return errors.Wrapf(err, "ошибка добавления комнаты %d", room.number)
        }
    }

    return nil
}


func (s *PGStorage) getShardForPartnership(partnershipID string) string {
    return s.getSchema(partnershipID)
}

func (s *PGStorage) Close() {
    s.mu.Lock()
    defer s.mu.Unlock()
    
    if s.db != nil {
        s.db.Close()
        s.db = nil
    }
}

func (s *PGStorage) DebugSharding(partnershipID string) string {
    schema := s.getSchema(partnershipID)
    h := fnv.New32a()
    h.Write([]byte(partnershipID))
    hash := h.Sum32()
    
    return fmt.Sprintf("partnership_id: %s → hash: %d → schema: %s", 
        partnershipID, hash, schema)
}