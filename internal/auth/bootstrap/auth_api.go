package bootstrap

import (
    server "github.com/nineteen-night/empty-room-game/internal/auth/api/auth_api"
    "github.com/nineteen-night/empty-room-game/internal/auth/services/authService"
)

func InitAuthServiceAPI(authService *authService.AuthService) *server.AuthServiceAPI {
    return server.NewAuthServiceAPI(authService)
}