package bootstrap

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"

	server "github.com/nineteen-night/empty-room-game/internal/game/api/game_api"
	gameconsumer "github.com/nineteen-night/empty-room-game/internal/game/consumer/game_consumer"

	"github.com/nineteen-night/empty-room-game/internal/game/pb/game_api"
	"github.com/go-chi/chi/v5"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func AppRun(api server.GameServiceAPI, gameConsumer *gameconsumer.GameConsumer) {
	go gameConsumer.Consume(context.Background())

	go func() {
		if err := runGRPCServer(api); err != nil {
			panic(fmt.Errorf("failed to run gRPC server: %v", err))
		}
	}()

	if err := runGatewayServer(); err != nil {
		panic(fmt.Errorf("failed to run gateway server: %v", err))
	}
}

func runGRPCServer(api server.GameServiceAPI) error {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	game_api.RegisterGameServiceServer(s, &api)

	slog.Info("Game gRPC-server listening on :50052")
	return s.Serve(lis)
}

func runGatewayServer() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	swaggerPath := os.Getenv("GAME_SWAGGER_PATH")
	if swaggerPath == "" {
		swaggerPath = "./internal/game/pb/swagger/game_api/game.swagger.json"
	}

	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		panic(fmt.Errorf("swagger file not found: %s", swaggerPath))
	}

	r := chi.NewRouter()
	r.Get("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, swaggerPath)
	})

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger.json"),
	))

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	err := game_api.RegisterGameServiceHandlerFromEndpoint(ctx, mux, ":50052", opts)
	if err != nil {
		panic(err)
	}

	r.Mount("/", mux)

	slog.Info("Game gRPC-Gateway server listening on :8081")
	return http.ListenAndServe(":8081", r)
}