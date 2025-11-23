package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gookit/color"
	"github.com/manifoldco/promptui"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/webtransport-go"
	"github.com/xealgo/muddy/api"
	"github.com/xealgo/muddy/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {
	color.Green.Println("Welcome to Muddy!")

	cfg, err := config.NewClientConfig(
		config.WithEnvPath(".env"),
	)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Create gRPC client
	client, conn, err := createGRPCClient(cfg)
	if err != nil {
		log.Fatalf("Failed to create gRPC client: %v", err)
	}
	defer conn.Close()

	uuid, err := login(client)
	if err != nil {
		log.Fatalf("Login failed: %v", err)
		os.Exit(1)
	}

	// fmt.Printf("Your session ID is %s\n", uuid)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Handle shutdown signal
	go func() {
		<-sigChan
		cancel()
	}()

	if err = connectWebTransport(ctx, cfg, uuid); err != nil {
		log.Fatalf("Failed to connect to WebTransport server: %v", err)
	}
}

// createGRPCClient establishes a connection to the gRPC server
func createGRPCClient(cfg *config.Config) (api.LoginServiceClient, *grpc.ClientConn, error) {
	creds := credentials.NewTLS(&tls.Config{
		InsecureSkipVerify: true, // Only for development with self-signed certs
	})

	// Connect to gRPC server
	conn, err := grpc.NewClient(fmt.Sprintf("localhost:%d", cfg.GrpcPort), grpc.WithTransportCredentials(creds))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	client := api.NewLoginServiceClient(conn)

	return client, conn, nil
}

// login prompts the user for a username with validation.
func login(client api.LoginServiceClient) (string, error) {
	validate := func(input string) error {
		l := len(input)
		if l < 3 || l > 12 {
			return fmt.Errorf("username must be between 3 and 12 characters")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    "Please enter a username",
		Validate: validate,
	}

	username, err := prompt.Run()
	if err != nil {
		return "", err
	}

	// Call login service
	ctx := context.Background()
	req := &api.LoginRequest{
		Username: username,
	}

	resp, err := client.Login(ctx, req)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}

	return resp.SessionUuid, nil
}

// connectWebTransport connects to the WebTransport server and starts the game session
func connectWebTransport(ctx context.Context, cfg *config.Config, uuid string) error {
	addr := fmt.Sprintf("https://localhost:%d/wt?uuid=%s", cfg.WTPort, uuid)

	color.Blue.Println("Connecting to Game Server")

	dialer := webtransport.Dialer{
		// Configure TLS, QUIC options, ALPN, etc., here if needed.
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		QUICConfig: &quic.Config{
			EnableDatagrams: true,
			KeepAlivePeriod: 30 * time.Minute,
		},
	}

	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, sess, err := dialer.Dial(timeoutCtx, addr, nil)
	if err != nil {
		return fmt.Errorf("failed to start session: %w", err)
	}

	cancel()

	// Open a bidi stream
	stream, err := sess.OpenStream()
	if err != nil {
		log.Fatalf("open stream: %v", err)
	}

	defer stream.Close()

	// Start message handling
	go handleIncomingMessages(stream)

	// Handle user input
	return handleUserInput(ctx, stream)
}

// handleIncomingMessages listens for messages from the server
//
// NOTE: This will eventually be entirely re-written.. just hacking away
// for the time being.
func handleIncomingMessages(stream *webtransport.Stream) {
	stream.Write([]byte{'\n'})

	buffer := make([]byte, 1024)
	for {
		n, err := stream.Read(buffer)
		if err != nil {
			var appError *quic.ApplicationError

			if errors.As(err, &appError) {
				color.Red.Println("Disconnected from server.")
				os.Exit(0)
				return
			}

			if err.Error() != "EOF" {
				color.Red.Printf("Error reading from server: %v\n", err)
			}
			return
		}

		message := string(buffer[:n])

		if strings.Contains(message, "event:") {
			type Event struct {
				Type      string
				Timestamp time.Time
				Data      interface{}
			}

			e := Event{}

			m := strings.ReplaceAll(message, "event:", "")
			err = json.Unmarshal([]byte(m), &e)
			if err != nil {
				slog.Error("failed to unmarshal event", "error", err)
				return
			}

			if e.Type == "RoomChat" {
				color.Yellow.Println(e.Data)
			}
		} else {
			color.Cyan.Print(message)
		}
	}
}

// handleUserInput processes user commands and sends them to server
func handleUserInput(ctx context.Context, stream *webtransport.Stream) error {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "quit" || input == "exit" {
			color.Yellow.Println("Until next time!")
			return nil
		}

		// Send command to server
		_, err := stream.Write([]byte(input))
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("input error: %w", err)
	}

	return nil
}
