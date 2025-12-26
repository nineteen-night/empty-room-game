package bootstrap

import (
	"context"

	"github.com/nineteen-night/empty-room-game/config"
	gameservice "github.com/nineteen-night/empty-room-game/internal/game/services/gameService"
	"github.com/nineteen-night/empty-room-game/internal/game/storage/pgstorage"
)

func InitGameService(storage *pgstorage.PGStorage, cfg *config.Config) *gameservice.GameService {
	return gameservice.NewGameService(context.Background(), storage)
}