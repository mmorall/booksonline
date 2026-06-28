.PHONY: help db-up db-down dev test lint build

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

db-up: ## Start the local database
	docker-compose up -d

db-down: ## Stop and remove the local database
	docker-compose down -v

dev: ## Run the application locally
	go run cmd/api/main.go

test: ## Run unit tests with race detection and coverage
	go test -v -race -cover ./...

lint: ## Run golangci-lint
	golangci-lint run

build: ## Build the binary locally
	go build -o bin/api cmd/api/main.go