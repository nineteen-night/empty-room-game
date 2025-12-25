package bootstrap

import (
    authProcessor "github.com/nineteen-night/empty-room-game/internal/auth/services/processors/auth_processor"
    authService "github.com/nineteen-night/empty-room-game/internal/auth/services/authService"
)

func InitAuthProcessor(authService *authService.AuthService) *authProcessor.AuthProcessor {
    return authProcessor.NewAuthProcessor(authService)
}