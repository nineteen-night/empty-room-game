package bootstrap

import (
    "fmt"
    "log"

    "github.com/nineteen-night/empty-room-game/config"
    "github.com/nineteen-night/empty-room-game/internal/auth/storage/pgstorage"
)

func InitPGStorage(cfg *config.Config) *pgstorage.PGstorage {
    connectionString := fmt.Sprintf(
        "postgres://%s:%s@%s:%d/%s?sslmode=%s",
        cfg.Database.AuthDB.Username,
        cfg.Database.AuthDB.Password,
        cfg.Database.AuthDB.Host,
        cfg.Database.AuthDB.Port,
        cfg.Database.AuthDB.DBName,
        cfg.Database.AuthDB.SSLMode,
    )
    
    storage, err := pgstorage.NewPGStorage(connectionString)
    if err != nil {
        log.Panic(fmt.Sprintf("ошибка инициализации БД, %v", err))
        panic(err)
    }
    return storage
}