#!/bin/bash

cd "$(dirname "$0")/.." || exit

# Генерация для Auth Service
protoc -I ./api \
  -I ./api/google/api \
  --go_out=./internal/auth/pb --go_opt=paths=source_relative \
  --go-grpc_out=./internal/auth/pb --go-grpc_opt=paths=source_relative \
  ./api/auth_api/auth.proto ./api/models/auth_model.proto

# Генерация gRPC-Gateway для Auth Service
protoc -I ./api \
  -I ./api/google/api \
  --grpc-gateway_out=./internal/auth/pb \
  --grpc-gateway_opt paths=source_relative \
  --grpc-gateway_opt logtostderr=true \
  ./api/auth_api/auth.proto

# Генерация OpenAPI для Auth Service
protoc -I ./api \
  -I ./api/google/api \
  --openapiv2_out=./internal/auth/pb/swagger \
  --openapiv2_opt logtostderr=true \
  ./api/auth_api/auth.proto

# Генерация для Game Service
protoc -I ./api \
  -I ./api/google/api \
  --go_out=./internal/game/pb --go_opt=paths=source_relative \
  --go-grpc_out=./internal/game/pb --go-grpc_opt=paths=source_relative \
  ./api/game_api/game.proto ./api/models/game_model.proto

# Генерация gRPC-Gateway для Game Service
protoc -I ./api \
  -I ./api/google/api \
  --grpc-gateway_out=./internal/game/pb \
  --grpc-gateway_opt paths=source_relative \
  --grpc-gateway_opt logtostderr=true \
  ./api/game_api/game.proto

# Генерация OpenAPI для Game Service
protoc -I ./api \
  -I ./api/google/api \
  --openapiv2_out=./internal/game/pb/swagger \
  --openapiv2_opt logtostderr=true \
  ./api/game_api/game.proto
