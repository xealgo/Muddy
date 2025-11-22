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
	"github.com/xealgo/muddy/internal/session"
)

const (
	DefaultMaxStreamBufferSize = 256
)

// Streaming represents a WebTransport server configuration
type Streaming struct {
	cfg  *config.Config
	addr string
	wt   *webtransport.Server
	sm   *session.SessionManager

	maxStreamBufferSize uint
}

// NewStreaming creates a new Streaming instance
func NewStreaming(cfg *config.Config, sm *session.SessionManager) (*Streaming, error) {
	s := &Streaming{
		cfg:  cfg,
		addr: fmt.Sprintf(":%d", cfg.WTPort),
		sm:   sm,
	}

	s.maxStreamBufferSize = DefaultMaxStreamBufferSize

	s.wt = &webtransport.Server{
		CheckOrigin: func(r *http.Request) bool {
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
func (s *Streaming) SetMaxStreamBufferSize(size uint) {
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

	session, err := s.wt.Upgrade(w, r)
	if err != nil {
		return fmt.Errorf("WebTransport upgrade failed: %w", err)
	}

	go s.handleSession(ctx, session)
	return nil
}

// handleSession manages a WebTransport session.
func (s *Streaming) handleSession(ctx context.Context, conn *webtransport.Session) {
	defer conn.CloseWithError(0, "connect closed")

	shutdownContext, cancel := context.WithCancel(conn.Context())

	go func() {
		<-ctx.Done()
		cancel()
	}()

	defer cancel()

	for {
		stream, err := conn.AcceptStream(conn.Context())
		if err != nil {
			if shutdownContext.Err() != nil {
				slog.Info("Shutting down session stream handler")
				return
			}

			slog.Error("Failed to accept stream", "error", err)
			return
		}

		go s.processStream(ctx, stream)
	}
}

// processStream handles an individual WebTransport stream.
func (s *Streaming) processStream(ctx context.Context, stream *webtransport.Stream) {
	defer stream.Close()

	buffer := make([]byte, s.maxStreamBufferSize)
	for {
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

		response := fmt.Sprintf("Server received message: %s", message)
		_, err = stream.Write(([]byte(response)))
		if err != nil {
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
