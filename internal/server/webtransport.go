// ///////////////////////////////////////////////////////////////////////////////
// webtransport.go
// Work in progress.
//
// TODO: Refactor...the code is a mess right now.
// TODO: Properly integrate with the command system once it's ready
// TODO: Improve error handling, context, etc.
// ////////////////////////////////////////////////////////////////////////////////
package server

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"

	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
	"github.com/xealgo/muddy/internal/config"
	"github.com/xealgo/muddy/internal/game"
	"github.com/xealgo/muddy/internal/session"
)

const (
	MaxBufferSize           = 1024
	MinBufferSize           = 256
	DefaultStreamBufferSize = 256
)

// Streaming represents a WebTransport server configuration
type Streaming struct {
	cfg  *config.Config
	addr string
	wt   *webtransport.Server
	sm   *session.SessionManager
	game *game.Game

	maxStreamBufferSize uint
}

// NewStreaming creates a new Streaming instance
func NewStreaming(cfg *config.Config, sm *session.SessionManager, game *game.Game) (*Streaming, error) {
	s := &Streaming{
		cfg:  cfg,
		addr: fmt.Sprintf(":%d", cfg.WTPort),
		sm:   sm,
		game: game,
	}

	s.maxStreamBufferSize = DefaultStreamBufferSize

	s.wt = &webtransport.Server{
		CheckOrigin: func(r *http.Request) bool {
			// TODO: Proper origin validation
			// strings.Contains(r.Host, ":1700")
			return true
		},
		H3: http3.Server{
			Addr:      s.addr,
			TLSConfig: cfg.TLSConfig,
		},
	}

	return s, nil
}

// SetMaxStreamBufferSize sets the maximum buffer size for each stream.
// size will be clamped between MinBufferSize and MaxBufferSize.
func (s *Streaming) SetMaxStreamBufferSize(size uint) {
	// Clamp between min and max buffer size
	if size > MaxBufferSize {
		size = MaxBufferSize
	} else if size < MinBufferSize {
		size = MinBufferSize
	}

	s.maxStreamBufferSize = size
}

// StartServer starts the WebTransport server.
func (s *Streaming) StartServer(ctx context.Context, wg *sync.WaitGroup) error {
	slog.Info("Starting WebTransport server", "address", s.addr)

	defer wg.Done()

	http.HandleFunc("/wt", func(w http.ResponseWriter, r *http.Request) {
		// Check if this is a WebTransport upgrade request
		if isWtConnectRequest(r) {
			err := s.handleUpgrade(ctx, w, r)
			if err != nil {
				slog.Error("Failed to upgrade to WebTransport session", "error", err)
				http.Error(w, "WebTransport upgrade failed", http.StatusBadRequest)
			}
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
	})

	// Handle errors created within the go routine
	errChan := make(chan error, 1)

	go func() {
		if err := s.wt.ListenAndServe(); err != nil {
			errChan <- err
			return
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("WebTransport server failed to start: %w", err)
	case <-ctx.Done():
		slog.Info("Shutting down WebTransport server")
		s.wt.Close()
		return nil
	}
}

// handleUpgrade handles the WebTransport upgrade request.
func (s *Streaming) handleUpgrade(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if !isWtConnectRequest(r) {
		return nil
	}

	uuid := r.URL.Query().Get("uuid")
	if len(uuid) == 0 {
		return fmt.Errorf("player session uuid required")
	}

	session, err := s.wt.Upgrade(w, r)
	if err != nil {
		return fmt.Errorf("WebTransport upgrade failed: %w", err)
	}

	// Handle the WebTransport session
	go s.handleSession(ctx, session, uuid)

	return nil
}

// handleSession manages a WebTransport session.
func (s *Streaming) handleSession(ctx context.Context, conn *webtransport.Session, sessionUUID string) {
	defer conn.CloseWithError(0, "connect closed")

	if len(sessionUUID) == 0 {
		slog.Error("Player session is nil")
		return
	}

	shutdownContext, cancel := context.WithCancel(conn.Context())

	go func() {
		<-ctx.Done()
		cancel()
	}()

	defer cancel()

	if _, ok := s.sm.GetSession(sessionUUID); ok {
		slog.Info("Player session already connected", "uuid", sessionUUID)
		return
	}

	// We'll need another go routine that periodically clears out pending sessions
	// after a certain time has passed.
	//
	// if ok := s.sm.RemovePending(sessionUUID); !ok {
	// 	slog.Error("Failed to remove pending player session", "uuid", sessionUUID)
	// }

	for {
		stream, err := conn.AcceptStream(conn.Context())
		if err != nil {
			player, exists := s.sm.GetSession(sessionUUID)
			if exists {
				fmt.Printf("%s has left the game\n", player.GetData().DisplayName)
				s.sm.RemovePlayerBySession(conn)
			}

			if shutdownContext.Err() != nil {
				slog.Info("Shutting down session stream handler")
				return
			}

			slog.Error("Failed to accept stream", "error", err)
			return
		}

		player, err := s.sm.Connect(sessionUUID, conn, stream)
		if err != nil {
			slog.Error("Failed to connect player session", "uuid", sessionUUID)
			stream.Write([]byte("Error creating player session. Disconnecting..."))
			return
		}

		// Eventually broadcast this..
		fmt.Printf("%s has joined the game\n", player.GetData().DisplayName)

		player.WriteString(fmt.Sprintf("Greetings %s!\n", player.GetData().DisplayName))

		go func(player *session.PlayerSession) {
			s.processStream(ctx, player)

			if player != nil {
				fmt.Printf("%s has left the game\n", player.GetData().DisplayName)
				s.sm.RemovePlayer(player.GetData().GetUUID())
			}
		}(player)
	}
}

// processStream handles an individual WebTransport stream.
func (s *Streaming) processStream(ctx context.Context, player *session.PlayerSession) {
	stream := player.GetStream()

	defer stream.Close()

	buffer := make([]byte, s.maxStreamBufferSize)
	for {
		select {
		case <-ctx.Done():
			slog.Info("Shutting down stream processor due to context cancellation")
			return
		default:
		}

		n, err := stream.Read(buffer)
		if err != nil {
			if ctx.Err() != nil {
				slog.Info("Shutting down stream processor", "error", ctx.Err())
				return
			}

			if err != io.EOF {
				slog.Error("Failed to read from stream", "error", err)
			}
			return
		}

		message := string(buffer[:n])
		if message == "PING" {
			if err := player.WriteString("PONG"); err != nil {
				slog.Error("Failed to write to stream", "error", err)
				return
			}
			continue
		}

		response := s.game.ProcessPlayerCommand(player, message)

		if err = player.WriteString(response); err != nil {
			slog.Error("Failed to write to stream", "error", err)
			return
		}

		if ctx.Err() != nil {
			slog.Info("Shutting down stream processor due to session closure", "error", ctx.Err())
			return
		}
	}
}

// isWtConnectRequest checks if the incoming HTTP request is a WebTransport CONNECT request.
func isWtConnectRequest(req *http.Request) bool {
	return req.Method == "CONNECT" && req.Proto == "webtransport"
}
