package game_consumer

import (
    "context"
    "encoding/json"
    "log/slog"
    "time"

    "github.com/nineteen-night/empty-room-game/internal/game/models"
    "github.com/segmentio/kafka-go"
)

func (c *GameConsumer) Consume(ctx context.Context) {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:           c.kafkaBroker,
        GroupID:           "GameService_group",
        Topic:             c.topicName,
        HeartbeatInterval: 3 * time.Second,
        SessionTimeout:    30 * time.Second,
    })
    defer r.Close()

    for {
        msg, err := r.ReadMessage(ctx)
        if err != nil {
            slog.Error("GameConsumer.consume error", "error", err.Error())
            continue
        }

        if c.topicName == "game_sessions_upsert" {
            var session models.GameSession
            err = json.Unmarshal(msg.Value, &session)
            if err != nil {
                slog.Error("parse game sessions error", "error", err)
                continue
            }
            err = c.gameProcessor.HandleGameSessions(ctx, &session)
        } else if c.topicName == "game_states_upsert" {
            var state models.GameState
            err = json.Unmarshal(msg.Value, &state)
            if err != nil {
                slog.Error("parse game states error", "error", err)
                continue
            }
            err = c.gameProcessor.HandleGameStates(ctx, &state)
        }

        if err != nil {
            slog.Error("Handle error", "error", err)
        }
    }
}