package auth_consumer

import (
    "context"
    "encoding/json"
    "log/slog"
    "time"

    "github.com/segmentio/kafka-go"
)

type RoomCompletedEvent struct {
    EventType  string `json:"event_type"`
    UserID     string `json:"user_id"`
    RoomNumber int32  `json:"room_number"`
}

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

        if c.topicName == "room_completed_events" {
            var event RoomCompletedEvent
            err = json.Unmarshal(msg.Value, &event)
            if err != nil {
                slog.Error("parse room_completed error", "error", err)
                continue
            }
            
            if event.EventType == "room_completed" {
                slog.Info("Processing room_completed", 
                    "user", event.UserID, 
                    "room", event.RoomNumber)
                
                err = c.authProcessor.HandleRoomCompleted(ctx, event.UserID, event.RoomNumber)
                if err != nil {
                    slog.Error("Failed to handle room_completed", 
                        "error", err, 
                        "user", event.UserID)
                } else {
                    slog.Info("Successfully processed room_completed", 
                        "user", event.UserID)
                }
            }
        }

        if err != nil {
            slog.Error("Handle error", "error", err)
        }
    }
}