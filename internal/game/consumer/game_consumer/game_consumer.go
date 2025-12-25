package game_consumer

import (
    "context"
    "github.com/nineteen-night/empty-room-game/internal/game/models"
)

type gameProcessor interface {
    HandleGameSessions(ctx context.Context, session *models.GameSession) error
    HandleGameStates(ctx context.Context, state *models.GameState) error
}

type GameConsumer struct {
    gameProcessor gameProcessor
    kafkaBroker   []string
    topicName     string
}

func NewGameConsumer(gameProcessor gameProcessor, kafkaBroker []string, topicName string) *GameConsumer {
    return &GameConsumer{
        gameProcessor: gameProcessor,
        kafkaBroker:   kafkaBroker,
        topicName:     topicName,
    }
}