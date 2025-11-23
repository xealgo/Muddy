package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gookit/color"

	"github.com/xealgo/muddy/internal/config"
	"github.com/xealgo/muddy/internal/game"
	"github.com/xealgo/muddy/internal/server"
	"github.com/xealgo/muddy/internal/services"
	"github.com/xealgo/muddy/internal/session"
	"github.com/xealgo/muddy/internal/world"
)

func main() {
	color.Green.Println("Starting Muddy!")

	cfg, err := config.NewConfig(
		config.WithEnvPath(".env"),
		config.WithDefaults(17000, 17001, 17002, "./certs/muddy.crt", "./certs/muddy-server.key"),
	)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Session manager instance used for managing player sessions
	sm := session.NewSessionManager(64)

	world := world.NewWorld()
	err = world.LoadRoomsFromYaml("./data/test-world.yml")
	if err != nil {
		slog.Error("Failed to load world data", "error", err)
		os.Exit(1)
	}

	game := game.NewGame(world)
	game.Sm = sm

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	wg := sync.WaitGroup{}

	// Handle shutdown signal
	go func() {
		<-sigChan
		cancel()
	}()

	// HTTP server setup
	httpServer, err := server.NewHttpServer(
		cfg,
		server.WithCORSHandler(),
		server.WithStaticPageHandlers(
			server.HttpServerStaticFileConfig{Path: "/", FilePath: "./public/index.html"},
			// server.HttpServerStaticFileConfig{Path: "/game-client", FilePath: "./public/client.html"},
		),
	)

	if err != nil {
		slog.Error("Failed to create HTTP server", "error", err)
		os.Exit(1)
	}

	fmt.Printf("https://localhost:%d/\n", cfg.HttpPort)

	wg.Add(1)
	go func() {
		if err := httpServer.Start(ctx, &wg); err != nil {
			slog.Error("HTTP server failed", "error", err)
		}
	}()

	// GRPC server setup
	grpcServer := server.NewGrpcServer(cfg)
	services.RegisterHealthService(cfg, grpcServer.Server, game.State(), sm)
	services.RegisterLoginService(cfg, grpcServer.Server, sm)

	wg.Add(1)
	go func() {
		if err := grpcServer.StartServer(ctx, &wg); err != nil {
			slog.Error("gRPC server failed", "error", err)
		}
	}()

	// WebTransport (streaming) server setup
	stream, err := server.NewStreaming(cfg, sm, game)
	if err != nil {
		slog.Error("Failed to create WebTransport server", "error", err)
		os.Exit(1)
	}

	wg.Add(1)
	go func() {
		if err = stream.StartServer(ctx, &wg); err != nil {
			slog.Error("WebTransport server failed", "error", err)
		}
	}()

	wg.Wait()
	slog.Info("Muddy server shutting down")
}
