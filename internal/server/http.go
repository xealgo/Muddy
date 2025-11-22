package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/xealgo/muddy/internal/config"
)

type HttpRouteHandler func(server *HttpServer) error
type HttpServerListener func() error

// HttpServerStaticFileConfig represents configuration for serving a static file
type HttpServerStaticFileConfig struct {
	Path     string
	FilePath string
}

// HttpServer represents an HTTP server configuration
type HttpServer struct {
	addr string
	cfg  *config.Config
}

// NewHttpServer creates a new HttpServer instance
func NewHttpServer(cfg *config.Config, handlers ...HttpRouteHandler) (*HttpServer, error) {
	hs := &HttpServer{
		addr: fmt.Sprintf("localhost:%d", cfg.HttpPort),
		cfg:  cfg,
	}

	for _, handler := range handlers {
		err := handler(hs)
		if err != nil {
			return nil, fmt.Errorf("failed to apply HTTP route handler - %v", err)
		}
	}

	return hs, nil
}

// Start starts the HTTP server
func (server HttpServer) Start(ctx context.Context, wg *sync.WaitGroup) error {
	h := &http.Server{
		Addr:      server.addr,
		TLSConfig: server.cfg.TLSConfig,
	}

	defer wg.Done()

	if h.Addr == "" {
		return fmt.Errorf("HTTP server address is empty, server will not start")
	}

	if h.TLSConfig == nil {
		return fmt.Errorf("HTTP server TLS configuration is nil, server will not start")
	}

	slog.Info("Starting HTTP server", "address", h.Addr)

	// Handle errors created in the go routine
	errChan := make(chan error, 1)

	go func() {
		if err := h.ListenAndServeTLS("", ""); err != nil {
			errChan <- err
			return
		}
	}()

	select {
	case err := <-errChan:
		return fmt.Errorf("HTTP server failed to start: %w", err)
	case <-ctx.Done():
		//
	}

	slog.Info("Shutting down HTTP server")

	shutdownCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := h.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("HTTP server shutdown failed: %w", err)
	}

	slog.Info("HTTP server shut down gracefully")

	return nil
}

// WithStaticFiles configures static file serving for the HTTP server
func WithStaticPageHandlers(entries ...HttpServerStaticFileConfig) HttpRouteHandler {
	return func(server *HttpServer) error {
		if len(entries) == 0 {
			return fmt.Errorf("no static files configured for http server")
		}

		for _, staticFile := range entries {
			finfo, err := os.Stat(staticFile.FilePath)
			if err != nil || finfo.IsDir() {
				slog.Error("Static file not found or is a directory", "FilePath", staticFile.FilePath)
				continue
			}

			http.HandleFunc(staticFile.Path, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/html")
				http.ServeFile(w, r, staticFile.FilePath)
			})
		}

		return nil
	}
}

// WithCORSHandler adds CORS headers to HTTP responses
func WithCORSHandler() HttpRouteHandler {
	return func(server *HttpServer) error {
		http.HandleFunc("/*", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
		})

		return nil
	}
}
