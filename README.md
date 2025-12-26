sudo systemctl stop postgresql
make down
docker volume rm empty-room-game_auth-postgres-data empty-room-game_game-postgres-data 2>/dev/null || true
make up
go run ./cmd/game_service
go run ./cmd/auth_service