package game_consumer

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/segmentio/kafka-go"
)

type PartnershipEvent struct {
	EventType     string `json:"event_type"`
	PartnershipID string `json:"partnership_id"`
	User1ID       string `json:"user1_id,omitempty"`
	User2ID       string `json:"user2_id,omitempty"`
}

func (c *GameConsumer) Consume(ctx context.Context) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:           c.kafkaBrokers,
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

		if c.topicName == "partnership_events" {
			var event PartnershipEvent
			err = json.Unmarshal(msg.Value, &event)
			if err != nil {
				slog.Error("parse partnership event error", "error", err)
				continue
			}

			switch event.EventType {
			case "partnership_created":
				slog.Info("Processing partnership_created",
					"partnership", event.PartnershipID,
					"user1", event.User1ID,
					"user2", event.User2ID)

				err = c.gameProcessor.HandlePartnershipCreated(ctx, event.PartnershipID, event.User1ID, event.User2ID)
				if err != nil {
					slog.Error("Failed to handle partnership_created",
						"error", err,
						"partnership", event.PartnershipID)
				} else {
					slog.Info("Successfully processed partnership_created",
						"partnership", event.PartnershipID)
				}

			case "partnership_terminated":
				slog.Info("Processing partnership_terminated",
					"partnership", event.PartnershipID)

				err = c.gameProcessor.HandlePartnershipTerminated(ctx, event.PartnershipID)
				if err != nil {
					slog.Error("Failed to handle partnership_terminated",
						"error", err,
						"partnership", event.PartnershipID)
				} else {
					slog.Info("Successfully processed partnership_terminated",
						"partnership", event.PartnershipID)
				}
			}
		}

		if err != nil {
			slog.Error("Handle error", "error", err)
		}
	}
}