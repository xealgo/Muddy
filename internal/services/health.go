package services

import (
	"context"

	"github.com/xealgo/muddy/api"
	"github.com/xealgo/muddy/internal/config"
	"github.com/xealgo/muddy/internal/game"
	"github.com/xealgo/muddy/internal/session"
	"google.golang.org/grpc"
)

// HealthService implements the health checking service.
type HealthService struct {
	api.HealthServiceServer
	cfg   *config.Config
	state *game.GameState
	sm    *session.SessionManager
}

// NewHealthService creates a new HealthService instance.
func RegisterHealthService(cfg *config.Config, server *grpc.Server, state *game.GameState, sm *session.SessionManager) {
	hs := &HealthService{
		cfg:   cfg,
		state: state,
		sm:    sm,
	}

	api.RegisterHealthServiceServer(server, hs)
}

// GetStatus returns the health status of the service.
func (s *HealthService) GetStatus(ctx context.Context, in *api.StatusRequest) (*api.StatusResponse, error) {
	return &api.StatusResponse{
		// Doesn't really make sense to have IsActive if it's just always going to be true,
		// figure out how to meaningfully define this.
		IsActive:      true,
		UptimeSeconds: int32(s.state.Uptime().Seconds()),
		ActiveUsers:   int32(s.sm.GetActiveSessionCount()),
	}, nil
}
