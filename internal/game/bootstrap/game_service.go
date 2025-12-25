package bootstrap

import (
    "context"
    "github.com/nineteen-night/empty-room-game/config"
    gameService "github.com/nineteen-night/empty-room-game/internal/game/services/gameService"
    "github.com/nineteen-night/empty-room-game/internal/game/storage/pgstorage"
)

func InitGameService(storage *pgstorage.PGstorage, cfg *config.Config) *gameService.GameService {
    return gameService.NewGameService(context.Background(), storage, 
        cfg.GameServiceSettings.MinNameLen, cfg.GameServiceSettings.MaxNameLen)
}