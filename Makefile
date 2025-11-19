# Define the name of your main Go package (e.g., "main" if your main.go is in the root)
SERVER_PACKAGE := ./cmd/muddy

BUILD_PATH := ./builds

# Define the name of the executable binary
SERVER_BINARY := "$(BUILD_PATH)/muddy"

# Default target: builds the application
all: build-server

# Sets up UDP by increasing the UDP buffer size at the OS level
setup-udp:
	./scripts/setupudp.sh

# Build target: compiles the Go application
build-server:
	@go build -o "$(SERVER_BINARY)" $(SERVER_PACKAGE)

# Run target: executes the compiled binary
run-server: build-server
	@$(SERVER_BINARY)

# Test target: runs all Go tests
test:
	go test -v ./...

# Clean target: removes the compiled binary
clean:
	rm -f $(SERVER_BINARY)

# Builds proto files in api/proto and outputs the Go code to api/
proto:
	protoc --go_out=./api --go-grpc_out=./api ./api/proto/*.proto	

# Phony targets ensure that make always executes these rules, even if a file with the same name exists.
.PHONY: all build-server run-server test clean