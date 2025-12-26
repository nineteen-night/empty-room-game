package bootstrap

import (
    "fmt"
    
    "github.com/nineteen-night/empty-room-game/config"
    authProcessor "github.com/nineteen-night/empty-room-game/internal/auth/services/processors/auth_processor"
    authService "github.com/nineteen-night/empty-room-game/internal/auth/services/authService"
)

func InitAuthProcessor(authService *authService.AuthService, cfg *config.Config) *authProcessor.AuthProcessor {
    kafkaBrokers := []string{fmt.Sprintf("%v:%v", cfg.Kafka.Host, cfg.Kafka.Port)}
    return authProcessor.NewAuthProcessor(authService, kafkaBrokers)
}