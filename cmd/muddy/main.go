package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/gookit/color"
	_ "github.com/manifoldco/promptui"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"

	"github.com/xealgo/muddy/internal/config"
	"github.com/xealgo/muddy/internal/server"
)

func main() {
	color.Green.Println("Starting Muddy!")

	cfg, err := config.NewConfig(
		config.WithEnvPath(".env"),
		config.WithDefaults(17000, "./certs/muddy.crt", "./certs/muddy-server.key"),
	)

	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	addr := fmt.Sprintf(":%d", cfg.Port)

	httpServer, err := server.NewHttpServer(
		cfg,
		server.WithCORSHandler(),
		server.WithStaticPageHandlers(
			server.HttpServerStaticFileConfig{Path: "/game-client", FilePath: "./test-webclient.html"},
		),
	)

	if err != nil {
		slog.Error("Failed to create HTTP server", "error", err)
		os.Exit(1)
	}

	wt := &webtransport.Server{
		CheckOrigin: func(r *http.Request) bool {
			return strings.Contains(r.Host, "localhost:1700")
		},
		H3: http3.Server{
			Addr:      addr,
			TLSConfig: cfg.TLSConfig,
		},
	}

	http.HandleFunc("/wt", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Received request method: %s, URL: %s\n", r.Method, r.URL.Path)

		// Check if this is a WebTransport upgrade request
		if r.Method == "CONNECT" && r.Proto == "webtransport" {
			session, err := wt.Upgrade(w, r)
			if err != nil {
				slog.Error("Failed to upgrade to WebTransport session", "error", err)
				http.Error(w, "WebTransport upgrade failed", http.StatusBadRequest)
				return
			}

			go handleWebTransportSession(session)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Muddy WebTransport Server Ready")
	})

	go func() {
		if err := httpServer.Setup(); err != nil {
			slog.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	defer wt.Close()

	slog.Info("Starting WebTransport server", "address", addr)
	if err := wt.ListenAndServe(); err != nil {
		slog.Error("Failed to start WebTransport server", "error", err)
		os.Exit(1)
	}
}

// TODO: Move this to the appropriate package
func handleWebTransportSession(conn *webtransport.Session) {
	defer conn.CloseWithError(0, "connect closed")

	for {
		stream, err := conn.AcceptStream(context.Background())
		if err != nil {
			slog.Error("Failed to accept stream", "error", err)
			return
		}

		go handleWebTransportStream(stream)
	}
}

// TODO: Move this to the appropriate package
func handleWebTransportStream(stream *webtransport.Stream) {
	defer stream.Close()

	buffer := make([]byte, 256)
	for {
		n, err := stream.Read(buffer)
		if err != nil {
			if err != io.EOF {
				slog.Error("Failed to read from stream", "error", err)
			}
			return
		}

		message := string(buffer[:n])
		slog.Info("Received message over WebTransport", "message", message)

		response := fmt.Sprintf("Server received message: %s", message)
		_, err = stream.Write(([]byte(response)))
		if err != nil {
			slog.Error("Failed to write to stream", "error", err)
			return
		}
	}
}
