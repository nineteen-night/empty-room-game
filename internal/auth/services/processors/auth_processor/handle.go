package auth_processor

import (
    "context"
    "encoding/json"
    "fmt"
    
    "github.com/segmentio/kafka-go"
)

func (p *AuthProcessor) HandleRoomCompleted(ctx context.Context, userID string, roomNumber int32) error {
    return p.authService.HandleRoomCompleted(ctx, userID, roomNumber)
}

func (p *AuthProcessor) SendPartnershipCreated(ctx context.Context, partnershipID, user1ID, user2ID string) error {
    event := map[string]interface{}{
        "event_type":     "partnership_created",
        "partnership_id": partnershipID,
        "user1_id":       user1ID,
        "user2_id":       user2ID,
    }
    
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }
    
    err = p.partnershipWriter.WriteMessages(ctx, kafka.Message{
        Value: data,
    })
    
    if err != nil {
        return fmt.Errorf("failed to write message to kafka: %w", err)
    }
    
    return nil
}

func (p *AuthProcessor) SendPartnershipTerminated(ctx context.Context, partnershipID string) error {
    event := map[string]interface{}{
        "event_type":     "partnership_terminated",
        "partnership_id": partnershipID,
    }
    
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }
    
    err = p.partnershipWriter.WriteMessages(ctx, kafka.Message{
        Value: data,
    })
    
    if err != nil {
        return fmt.Errorf("failed to write message to kafka: %w", err)
    }
    
    return nil
}