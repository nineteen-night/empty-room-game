package bootstrap

import (
    gameProcessor "github.com/nineteen-night/empty-room-game/internal/game/services/processors/game_processor"
    gameService "github.com/nineteen-night/empty-room-game/internal/game/services/gameService"
)

func InitGameProcessor(gameService *gameService.GameService) *gameProcessor.GameProcessor {
    return gameProcessor.NewGameProcessor(gameService)
}