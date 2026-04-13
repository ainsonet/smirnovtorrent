.PHONY: build test clean help lint vet run build-linux build-windows build-macos

VERSION=$(shell grep "const version" cmd/smirnovtorrent/main.go | cut -d'"' -f2)
BINARY_NAME=smirnovtorrent
BUILD_DIR=bin

# Build the application
build:
	@echo "Building $(BINARY_NAME) v$(VERSION)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/smirnovtorrent
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Cross-platform builds
build-linux:
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux ./cmd/smirnovtorrent

build-windows:
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME).exe ./cmd/smirnovtorrent

build-macos:
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-macos ./cmd/smirnovtorrent

build-all: build-linux build-windows build-macos
	@echo "All platforms built successfully!"

# Run tests
test:
	go test -v ./...

# Run tests with coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Clean build artifacts
clean:
	rm -f smirnovtorrent.exe
	rm -f coverage.out
	rm -f coverage.html

# Run linter (requires golangci-lint)
lint:
	golangci-lint run ./...

# Run go vet
vet:
	go vet ./...

# Format code
fmt:
	go fmt ./...

# Build and run
run: build
	@echo "Running $(BINARY_NAME) v$(VERSION)..."
	./$(BUILD_DIR)/$(BINARY_NAME)

# Show dependencies
deps:
	go mod graph

# Update dependencies
update-deps:
	go get -u ./...
	go mod tidy

# Help
help:
	@echo "SmirnovTorrent v$(VERSION) - BitTorrent Client in Go"
	@echo ""
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-linux    - Build for Linux"
	@echo "  build-windows  - Build for Windows"
	@echo "  build-macos    - Build for macOS"
	@echo "  build-all      - Build for all platforms"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  clean          - Remove build artifacts"
	@echo "  lint           - Run linter (requires golangci-lint)"
	@echo "  vet            - Run go vet"
	@echo "  fmt            - Format code"
	@echo "  run            - Build and run"
	@echo "  deps           - Show dependencies"
	@echo "  update-deps    - Update dependencies"
	@echo "  help           - Show this help"
	@echo ""
	@echo "Examples:"
	@echo "  make build          # Build for current platform"
	@echo "  make build-all      # Build for all platforms"
	@echo "  make test           # Run all tests"
	@echo "  make run            # Build and run"