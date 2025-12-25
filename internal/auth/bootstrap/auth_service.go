package bootstrap

import (
    "context"

    "github.com/nineteen-night/empty-room-game/config"
    authService "github.com/nineteen-night/empty-room-game/internal/auth/services/authService"
    "github.com/nineteen-night/empty-room-game/internal/auth/storage/pgstorage"
)

func InitAuthService(storage *pgstorage.PGstorage, cfg *config.Config) *authService.AuthService {
    return authService.NewAuthService(context.Background(), storage, 
        cfg.AuthServiceSettings.MinNameLen, cfg.AuthServiceSettings.MaxNameLen)
}