package bootstrap

import (
    "fmt"
    "log"

    "github.com/nineteen-night/empty-room-game/config"
    "github.com/nineteen-night/empty-room-game/internal/game/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGStorage {
    connectionString := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=%s",
        cfg.Database.GameDB.Username,  
        cfg.Database.GameDB.Password,
        cfg.Database.GameDB.Host,
        cfg.Database.GameDB.Port,
        cfg.Database.GameDB.DBName,
        cfg.Database.GameDB.SSLMode,
    )
    
    log.Printf("Connecting to database: %s@%s:%d/%s", 
        cfg.Database.GameDB.Username,
        cfg.Database.GameDB.Host,
        cfg.Database.GameDB.Port,
        cfg.Database.GameDB.DBName)
    

    shardCount := 2  
    
    storage, err := pgstorage.NewPGStorage(connectionString, shardCount)
    if err != nil {
        log.Panic(fmt.Sprintf("ошибка инициализации БД, %v", err))
        panic(err)
    }
    
    log.Println("Database initialized successfully with schema-based sharding")
    return storage
}