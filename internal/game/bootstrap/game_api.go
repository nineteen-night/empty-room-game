package bootstrap

import (
	server "github.com/nineteen-night/empty-room-game/internal/game/api/game_api"
	gameservice "github.com/nineteen-night/empty-room-game/internal/game/services/gameService"
)

func InitGameServiceAPI(gameService *gameservice.GameService) *server.GameServiceAPI {
	return server.NewGameServiceAPI(gameService)
}