package server

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/xealgo/muddy/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
)

// GrpcServer represents a gRPC server instance.
type GrpcServer struct {
	Server *grpc.Server

	cfg *config.Config
}

// NewGrpcServer creates a new GrpcServer instance.
func NewGrpcServer(cfg *config.Config) *GrpcServer {
	creds := credentials.NewTLS(cfg.TLSConfig)

	return &GrpcServer{
		Server: grpc.NewServer(grpc.Creds(creds)),
		cfg:    cfg,
	}
}

// StartServer starts the gRPC server and listens for incoming connections.
func (gs *GrpcServer) StartServer(ctx context.Context, wg *sync.WaitGroup) error {
	defer wg.Done()

	addr := fmt.Sprintf(":%d", gs.cfg.GrpcPort)

	slog.Info("Starting GRPC server", "address", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	// Enable reflection for debugging
	reflection.Register(gs.Server)

	// Handle errors created in the go routine
	errChan := make(chan error, 1)

	go func() {
		if err := gs.Server.Serve(lis); err != nil {
			errChan <- err
			return
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("grpc server failed to start: %w", err)
	case <-ctx.Done():
		slog.Info("Gracefully shutting down gRPC server")
		gs.Server.GracefulStop()
	}

	return nil
}
