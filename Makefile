SERVER_PACKAGE := ./cmd/muddy
CLIENT_PACKAGE := ./cmd/client
BUILD_PATH := ./builds
SERVER_BINARY := "$(BUILD_PATH)/muddy"
CLIENT_BINARY := "$(BUILD_PATH)/client"

# Default target: builds the application
all: build-server

# Sets up UDP by increasing the UDP buffer size at the OS level
setup-udp:
	./scripts/setupudp.sh

# Build target: compiles the server
build-server:
	@go build -o "$(SERVER_BINARY)" $(SERVER_PACKAGE)

# Build target: compiles the client
build-client:
	@go build -o "$(CLIENT_BINARY)" $(CLIENT_PACKAGE)

# Run target: builds and runs the server
run-server: build-server
	@$(SERVER_BINARY)

# Run target: builds and runs the client
run-client: build-client
	@$(CLIENT_BINARY)

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