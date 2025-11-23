package services

import (
	"context"

	"github.com/xealgo/muddy/api"
	"github.com/xealgo/muddy/internal/config"
	"github.com/xealgo/muddy/internal/game"
	"google.golang.org/grpc"
)

// LoginService implements the login service.
type LoginService struct {
	api.LoginServiceServer
	cfg *config.Config
	sm  *game.SessionManager
}

// LoginService registers the LoginService with the gRPC server.
func RegisterLoginService(cfg *config.Config, server *grpc.Server, sm *game.SessionManager) {
	service := &LoginService{
		cfg: cfg,
		sm:  sm,
	}

	api.RegisterLoginServiceServer(server, service)
}

// Login handles user login requests.
func (s *LoginService) Login(ctx context.Context, req *api.LoginRequest) (*api.LoginResponse, error) {
	player := game.NewPlayer(req.Username, req.Username)

	err := s.sm.Register(player)
	if err != nil {
		return nil, err
	}

	return &api.LoginResponse{
		SessionUuid: player.GetUUID(),
	}, nil
}
