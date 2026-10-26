.PHONY: help build run dev test test-coverage lint fmt clean docker-up docker-down migrate-create migrate-up migrate-down migrate-status migrate-reset migrate-refresh migrate-up-to migrate-down-to swag install-tools seed-all seed

# Variables
BINARY_NAME=coursepack
BUILD_DIR=bin
MIGRATIONS_DIR=database/migrations

-include .env
export
DB_DSN := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

install-tools: ## Install development tools
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/cosmtrek/air@latest
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

build: ## Build the application
	mkdir -p $(BUILD_DIR)
	CGO_ENABLED=0 go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/api

run: build ## Run the application locally
	./$(BUILD_DIR)/$(BINARY_NAME)

dev: ## Run the application with hot-reload (requires air)
	swag init -q -g cmd/api/main.go -o docs/ --outputTypes go,json
	air

test: ## Run unit tests
	go test -v -race -coverprofile=coverage.out ./...

test-coverage: test ## Run tests with coverage report
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

lint: ## Run golangci-lint
	golangci-lint run --timeout=5m ./...

fmt: ## Format code
	go fmt ./...
	gofmt -s -w .

clean: ## Clean build artifacts
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html
	rm -rf ./uploads/*

docker-build: ## Build Docker image
	docker-compose build

docker-up: ## Start Docker containers
	docker-compose up -d
	@echo "Waiting for services to be healthy..."
	@sleep 5
	docker-compose ps

docker-down: ## Stop Docker containers
	docker-compose down

migrate-create: ## Create a new migration (usage: make migrate-create NAME=your_migration_name)
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is required"; exit 1; fi
	goose -dir $(MIGRATIONS_DIR) create $(NAME) sql

migrate-up: ## Run all pending migrations
	@echo "Running pending migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up

migrate-down: ## Rollback the last migration
	@echo "Rolling back last migration..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down

migrate-down-to: ## Rollback to a specific version (usage: make migrate-down-to VERSION=5)
	@if [ -z "$(VERSION)" ]; then echo "Error: VERSION is required"; exit 1; fi
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" down-to $(VERSION)

migrate-up-to: ## Migrate up to a specific version (usage: make migrate-up-to VERSION=10)
	@if [ -z "$(VERSION)" ]; then echo "Error: VERSION is required"; exit 1; fi
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" up-to $(VERSION)

migrate-status: ## Show migration status
	@echo "Migration status:"
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" status

migrate-reset: ## Roll back ALL migrations
	@echo "Rolling back all migrations..."
	goose -dir $(MIGRATIONS_DIR) postgres "$(DB_DSN)" reset
	@echo "All migrations have been rolled down"

migrate-refresh: migrate-reset migrate-up ## Reset and re-run all migrations
	@echo "Database migration refresh complete"

swag: ## Generate Swagger documentation
	swag init -g cmd/api/main.go -o docs/ --outputTypes go,json

seed-all: ## Run all seeders in registration order
	go run ./cmd/seed -all

seed: ## Run a specific seeder by name (usage: make seed NAME=roles)
	@if [ -z "$(NAME)" ]; then echo "Error: NAME is required. Usage: make seed NAME=<seeder_name>"; exit 1; fi
	go run ./cmd/seed -seeder=$(NAME)

health: ## Check server health
	curl -v http://localhost:8000/health

.DEFAULT_GOAL := help
