package bootstrap

import (
    "context"
    "fmt"
    "log/slog"
    "net"
    "net/http"
    "os"

    server "github.com/nineteen-night/empty-room-game/internal/auth/api/auth_api"
    authconsumer "github.com/nineteen-night/empty-room-game/internal/auth/consumer/auth_consumer"
    
    "github.com/nineteen-night/empty-room-game/internal/auth/pb/auth_api"
    "github.com/go-chi/chi/v5"
    "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
    httpSwagger "github.com/swaggo/http-swagger"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func AppRun(api server.AuthServiceAPI, authConsumer *authconsumer.AuthConsumer) {
    go authConsumer.Consume(context.Background())
    
    go func() {
        if err := runGRPCServer(api); err != nil {
            panic(fmt.Errorf("failed to run gRPC server: %v", err))
        }
    }()

    if err := runGatewayServer(); err != nil {
        panic(fmt.Errorf("failed to run gateway server: %v", err))
    }
}

func runGRPCServer(api server.AuthServiceAPI) error {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        return err
    }

    s := grpc.NewServer()
    auth_api.RegisterAuthServiceServer(s, &api)

    slog.Info("gRPC-server server listening on :50051")
    return s.Serve(lis)
}

func runGatewayServer() error {
    ctx := context.Background()
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    swaggerPath := os.Getenv("swaggerPath")
    if swaggerPath == "" {
        swaggerPath = "./internal/auth/pb/swagger/auth_api/auth.swagger.json"
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

    err := auth_api.RegisterAuthServiceHandlerFromEndpoint(ctx, mux, ":50051", opts)
    if err != nil {
        panic(err)
    }

    r.Mount("/", mux)

    slog.Info("gRPC-Gateway server listening on :8080")
    return http.ListenAndServe(":8080", r)
}