.PHONY: help build run test clean migrate-up migrate-down docker-up docker-down

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

build: ## Build the application
	go build -o bin/server cmd/server/main.go

run: ## Run the application
	go run cmd/server/main.go

test: ## Run tests
	go test -v -cover ./...

clean: ## Clean build artifacts
	rm -rf bin/

deps: ## Download dependencies
	go mod download
	go mod tidy

migrate-up: ## Run database migrations up
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/project_management?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/project_management?sslmode=disable" down

migrate-create: ## Create a new migration (use NAME=migration_name)
	migrate create -ext sql -dir migrations -seq $(NAME)

docker-up: ## Start Docker containers
	docker-compose up -d

docker-down: ## Stop Docker containers
	docker-compose down

docker-logs: ## View Docker logs
	docker-compose logs -f

lint: ## Run linter
	golangci-lint run

fmt: ## Format code
	go fmt ./...

watch: ## Run with auto-reload (requires air)
	air
