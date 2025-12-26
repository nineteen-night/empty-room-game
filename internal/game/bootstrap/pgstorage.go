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

	storage, err := pgstorage.NewPGStorage(connectionString)
	if err != nil {
		log.Panic(fmt.Sprintf("ошибка инициализации БД, %v", err))
		panic(err)
	}
	return storage
}