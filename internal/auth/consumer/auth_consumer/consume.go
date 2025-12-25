package auth_consumer

import (
    "context"
    "encoding/json"
    "log/slog"
    "time"

    "github.com/nineteen-night/empty-room-game/internal/auth/models"
    "github.com/segmentio/kafka-go"
)

func (c *AuthConsumer) Consume(ctx context.Context) {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:           c.kafkaBroker,
        GroupID:           "AuthService_group",
        Topic:             c.topicName,
        HeartbeatInterval: 3 * time.Second,
        SessionTimeout:    30 * time.Second,
    })
    defer r.Close()

    for {
        msg, err := r.ReadMessage(ctx)
        if err != nil {
            slog.Error("AuthConsumer.consume error", "error", err.Error())
            continue
        }

        if c.topicName == "users_upsert" {
            var user models.User
            err = json.Unmarshal(msg.Value, &user)
            if err != nil {
                slog.Error("parse users error", "error", err)
                continue
            }
            err = c.authProcessor.HandleUsers(ctx, &user)
        } else if c.topicName == "partnerships_upsert" {
            var partnership models.Partnership
            err = json.Unmarshal(msg.Value, &partnership)
            if err != nil {
                slog.Error("parse partnerships error", "error", err)
                continue
            }
            err = c.authProcessor.HandlePartnerships(ctx, &partnership)
        }

        if err != nil {
            slog.Error("Handle error", "error", err)
        }
    }
}