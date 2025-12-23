# Makefile untuk Golang Clean Architecture

.PHONY: help build run test clean migrate-up migrate-down migrate-create swagger deps install-tools

# Load .env file
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Variables
APP_NAME=go-rest-scaffold
MAIN_PATH=cmd/web/main.go
BUILD_DIR=build
BINARY_NAME=app
MIGRATION_DIR=db/migrations

# Application Configuration
APP_PORT ?= 3000

# Database Configuration
# Default values if not set in .env
DB_USER ?= postgres
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= golang_clean_architecture

# Map DB_PASSWORD from .env to DB_PASS, default to postgres
DB_PASS ?= $(if $(DB_PASSWORD),$(DB_PASSWORD),postgres)

DB_URL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Default target
help:
	@echo "Available commands:"
	@echo "  make build          - Build aplikasi"
	@echo "  make run            - Run aplikasi"
	@echo "  make test           - Run unit tests"
	@echo "  make swagger        - Generate swagger documentation"
	@echo "  make migrate-up     - Run database migrations up"
	@echo "  make migrate-down   - Run database migrations down"
	@echo "  make migrate-create - Create new migration (name=create_table_xxx)"
	@echo "  make deps           - Download dependencies"
	@echo "  make install-tools  - Install required tools (swag, migrate)"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make dev            - Run in development mode with hot reload"
	@echo "  make lint           - Run linter"
	@echo "  make fmt            - Format code"
	@echo ""
	@echo "Docker commands:"
	@echo "  make docker-build   - Build Docker image"
	@echo "  make docker-up      - Start all services (app + postgres)"
	@echo "  make docker-up-build- Rebuild and start all services"
	@echo "  make docker-down    - Stop all Docker services"
	@echo "  make docker-logs    - View Docker logs"
	@echo "  make docker-migrate - Run migrations in Docker"
	@echo "  make docker-clean   - Remove Docker volumes and images"

# Build aplikasi
build:
	@echo "Building application..."
	@go build -o $(BUILD_DIR)/$(BINARY_NAME).exe $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME).exe"

# Run aplikasi
run:
	@echo "Running application..."
	@go run $(MAIN_PATH)

# Run development mode dengan air (hot reload)
dev:
	@echo "Starting development mode..."
	@go install github.com/air-verse/air@latest
	@air

# Run unit tests
test:
	@echo "Running tests..."
	@go test -v ./test/

# Run tests dengan coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -cover -coverprofile=coverage.out ./test/
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Generate swagger documentation
swagger:
	@echo "Generating swagger documentation..."
	@swag init -g $(MAIN_PATH)
	@echo "Swagger docs generated in docs/"

# Download dependencies
deps:
	@echo "Downloading dependencies..."
	@go mod download
	@go mod tidy
	@echo "Dependencies downloaded"

# Install required tools
install-tools:
	@echo "Installing required tools..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Tools installed successfully"

# Database migrations
migrate-up:
	@echo "Running migrations up..."
	@migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) up
	@echo "Migrations completed"

migrate-down:
	@echo "Running migrations down..."
	@migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) down
	@echo "Migrations rolled back"

migrate-force:
	@echo "Force migration to version: $(version)"
	@migrate -database "$(DB_URL)" -path $(MIGRATION_DIR) force $(version)

migrate-create:
ifndef name
	$(error name is undefined. Usage: make migrate-create name=create_table_xxx)
endif
	@echo "Creating migration: $(name)"
	@migrate create -ext sql -dir $(MIGRATION_DIR) $(name)
	@echo "Migration created"

# Linting
lint:
	@echo "Running linter..."
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint run ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...
	@echo "Code formatted"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@rm -f app.exe
	@echo "Clean complete"

# Setup project (untuk first time setup)
setup: install-tools deps migrate-up swagger
	@echo "Project setup complete!"
	@echo "Run 'make run' to start the application"

# Build untuk production
build-prod:
	@echo "Building for production..."
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "Production builds complete"

# Docker commands
docker-build:
	@echo "Building Docker image..."
	@docker build -t $(APP_NAME):latest .

docker-run:
	@echo "Running Docker container on port $(APP_PORT)..."
	@docker run -p $(APP_PORT):$(APP_PORT) \
		-e APP_PORT=$(APP_PORT) \
		-e APP_NAME=$(APP_NAME) \
		--env-file .env \
		$(APP_NAME):latest

docker-up:
	@echo "Starting services with Docker Compose..."
	@docker compose up -d

docker-up-build:
	@echo "Building and starting services with Docker Compose..."
	@docker compose up -d --build

docker-down:
	@echo "Stopping Docker Compose services..."
	@docker compose down

docker-logs:
	@echo "Showing Docker Compose logs..."
	@docker compose logs -f

docker-migrate:
	@echo "Running database migrations in Docker..."
	@docker compose --profile migrate up migrate

docker-clean:
	@echo "Cleaning Docker resources..."
	@docker compose down -v --rmi local
	@echo "Docker resources cleaned"

docker-shell:
	@echo "Opening shell in app container..."
	@docker compose exec app sh

# Database commands
db-create:
	@echo "Creating database..."
	@createdb -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) $(DB_NAME)
	@echo "Database created: $(DB_NAME)"

db-drop:
	@echo "Dropping database..."
	@dropdb -U $(DB_USER) -h $(DB_HOST) -p $(DB_PORT) $(DB_NAME)
	@echo "Database dropped: $(DB_NAME)"

db-reset: db-drop db-create migrate-up
	@echo "Database reset complete"

