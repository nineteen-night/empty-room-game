package bootstrap

import (
    server "github.com/nineteen-night/empty-room-game/internal/game/api/game_api"
    "github.com/nineteen-night/empty-room-game/internal/game/services/gameService"
)

func InitGameServiceAPI(gameService *gameService.GameService) *server.GameServiceAPI {
    return server.NewGameServiceAPI(gameService)
}