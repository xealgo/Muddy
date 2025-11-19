package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/xealgo/muddy/internal/config"
)

type HttpRouteHandler func(server *HttpServer) error
type HttpServerListener func() error

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
		addr: fmt.Sprintf("localhost:%d", cfg.Port+1),
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

// Setup initializes and returns the HTTP server listener
func (server HttpServer) Setup() error {
	h := &http.Server{
		Addr:      server.addr,
		TLSConfig: server.cfg.TLSConfig,
	}

	if h.Addr == "" {
		return fmt.Errorf("HTTP server address is empty, server will not start")
	}

	if h.TLSConfig == nil {
		return fmt.Errorf("HTTP server TLS configuration is nil, server will not start")
	}

	slog.Info("Starting HTTP server", "address", h.Addr)

	if err := h.ListenAndServeTLS("", ""); err != nil {
		return fmt.Errorf("HTTP server failed - %v", err)
	}

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
