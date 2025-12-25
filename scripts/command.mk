.PHONY: generate-api
generate-api:
	@./scripts/generate.sh

.PHONY: up
up:
	docker-compose up -d

.PHONY: down
down:
	docker-compose down

.PHONY: cov
cov:
	go test -cover ./...

.PHONY: mock
mock:
	mockery

.PHONY: build-auth
build-auth:
	go build -o bin/auth-service ./cmd/auth_service

.PHONY: build-game
build-game:
	go build -o bin/game-service ./cmd/game_service

.PHONY: run-auth
run-auth: build-auth
	./bin/auth-service

.PHONY: run-game
run-game: build-game
	./bin/game-service

.PHONY: test-auth
test-auth:
	go test ./internal/auth/...

.PHONY: test-game
test-game:
	go test ./internal/game/...

.PHONY: logs
logs:
	docker-compose logs -f

.PHONY: mock-auth
mock-auth:
	mockery --dir internal/auth/services/authService --name AuthStorage --output internal/auth/services/authService/mocks --with-expecter