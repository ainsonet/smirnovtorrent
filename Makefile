.PHONY: build test clean help lint vet

# Build the application
build:
	go build -o smirnovtorrent.exe ./cmd/smirnovtorrent

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
	./smirnovtorrent.exe

# Show dependencies
deps:
	go mod graph

# Update dependencies
update-deps:
	go get -u ./...
	go mod tidy

# Help
help:
	@echo "Available targets:"
	@echo "  build       - Build the application"
	@echo "  test        - Run tests"
	@echo "  test-coverage - Run tests with coverage report"
	@echo "  clean       - Remove build artifacts"
	@echo "  lint        - Run linter (requires golangci-lint)"
	@echo "  vet         - Run go vet"
	@echo "  fmt         - Format code"
	@echo "  run         - Build and run"
	@echo "  deps        - Show dependencies"
	@echo "  update-deps - Update dependencies"
	@echo "  help        - Show this help"