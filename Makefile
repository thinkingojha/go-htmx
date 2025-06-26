.PHONY: build run test tidy clean docker-build docker-run dev lint security-check deps

# Variables
BINARY_NAME=gohtmx
VERSION?=1.0.0
BUILD_DIR=./bin
DOCKER_IMAGE=go-htmx
DOCKER_TAG?=latest

# Build the application
build: clean
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
		-ldflags="-w -s -X main.version=$(VERSION)" \
		-o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Build for current platform
build-local: clean
	@echo "Building $(BINARY_NAME) for local platform..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

# Run the application in development mode
dev: build-local
	@echo "Starting development server..."
	@cp config.dev.yaml config.yaml || true
	@GOHTMX_APP_ENVIRONMENT=development $(BUILD_DIR)/$(BINARY_NAME)

# Run the application
run: build-local
	@echo "Starting server..."
	@$(BUILD_DIR)/$(BINARY_NAME)

# Run tests
test:
	@echo "Running tests..."
	@go test -v -race -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with short flag
test-short:
	@echo "Running short tests..."
	@go test -short -v ./...

# Benchmark tests
bench:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Format, lint and tidy
tidy:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Tidying modules..."
	@go mod tidy
	@echo "Vetting code..."
	@go vet ./...

# Lint the code (requires golangci-lint)
lint:
	@echo "Linting code..."
	@which golangci-lint > /dev/null || (echo "golangci-lint not installed. Install with: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest" && exit 1)
	@golangci-lint run ./...

# Security check (requires gosec)
security-check:
	@echo "Running security checks..."
	@which gosec > /dev/null || (echo "gosec not installed. Install with: go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest" && exit 1)
	@gosec ./...

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@npm ci

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html

# Build Tailwind CSS
build-css:
	@echo "Building CSS..."
	@npx tailwindcss build ./internal/static/css/main.css -o ./internal/static/css/tailwind.css --minify

# Docker build
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

# Docker run
docker-run: docker-build
	@echo "Running Docker container..."
	@docker run -p 8080:8080 --rm $(DOCKER_IMAGE):$(DOCKER_TAG)

# Docker run in development mode
docker-dev:
	@echo "Building and running Docker container in development mode..."
	@docker build -f Dockerfile.dev -t $(DOCKER_IMAGE):dev .
	@docker run -p 3000:3000 -v $(PWD):/app --rm $(DOCKER_IMAGE):dev

# Health check
health-check:
	@echo "Checking application health..."
	@curl -f http://localhost:3000/health || exit 1

# Full CI pipeline
ci: deps tidy lint security-check test build

# Production release build
release: ci docker-build
	@echo "Release $(VERSION) built successfully!"

# Help
help:
	@echo "Available targets:"
	@echo "  build           - Build the application binary"
	@echo "  build-local     - Build for current platform"
	@echo "  run             - Build and run the application"
	@echo "  dev             - Run in development mode"
	@echo "  test            - Run tests with coverage"
	@echo "  test-short      - Run short tests"
	@echo "  bench           - Run benchmarks"
	@echo "  tidy            - Format code and tidy modules"
	@echo "  lint            - Run linters"
	@echo "  security-check  - Run security checks"
	@echo "  deps            - Install dependencies"
	@echo "  clean           - Clean build artifacts"
	@echo "  build-css       - Build Tailwind CSS"
	@echo "  docker-build    - Build Docker image"
	@echo "  docker-run      - Build and run Docker container"
	@echo "  docker-dev      - Run in development mode with Docker"
	@echo "  health-check    - Check application health"
	@echo "  ci              - Run full CI pipeline"
	@echo "  release         - Build production release"
	@echo "  help            - Show this help message"