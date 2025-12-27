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