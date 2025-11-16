.PHONY: help install dev-backend dev-frontend dev stop clean

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

install: ## Install all dependencies
	@echo "Installing backend dependencies..."
	cd apps/backend && go mod download
	@echo "Installing frontend dependencies..."
	cd apps/frontend && npm install
	@echo "Done!"

dev-db: ## Start PostgreSQL database
	docker-compose up -d postgres
	@echo "Waiting for database to be ready..."
	@sleep 3
	@echo "Database is ready!"

dev-backend: ## Start backend development server
	cd apps/backend && air

dev-frontend: ## Start frontend development server
	cd apps/frontend && npm run dev

dev: dev-db ## Start all development servers (in separate terminals)
	@echo "Database started!"
	@echo ""
	@echo "To start the full stack, run these commands in separate terminals:"
	@echo "  make dev-backend"
	@echo "  make dev-frontend"

stop: ## Stop all services
	docker-compose down

clean: ## Clean up generated files and dependencies
	rm -rf apps/backend/tmp
	rm -rf apps/frontend/node_modules
	rm -rf apps/frontend/dist
	docker-compose down -v

migrate-up: ## Run database migrations
	@cd apps/backend && \
		export $$(grep -v '^#' .env | xargs) && \
		migrate -path migrations -database "$$DATABASE_URL" up

migrate-down: ## Rollback last database migration
	@cd apps/backend && \
		export $$(grep -v '^#' .env | xargs) && \
		migrate -path migrations -database "$$DATABASE_URL" down 1

migrate-create: ## Create a new migration (usage: make migrate-create NAME=migration_name)
	@if [ -z "$(NAME)" ]; then \
		echo "Error: NAME is required. Usage: make migrate-create NAME=migration_name"; \
		exit 1; \
	fi
	cd apps/backend && migrate create -ext sql -dir migrations -seq $(NAME)

test-db-setup: ## Set up test database with migrations
	@echo "Setting up test database..."
	@cd apps/backend && \
		export $$(grep -v '^#' .env.test | xargs) && \
		migrate -path migrations -database "$$TEST_DATABASE_URL" up
	@echo "Test database ready!"

test-db-reset: ## Reset test database (drop and recreate all tables)
	@echo "Resetting test database..."
	@cd apps/backend && \
		export $$(grep -v '^#' .env.test | xargs) && \
		echo "y" | migrate -path migrations -database "$$TEST_DATABASE_URL" down || true
	@cd apps/backend && \
		export $$(grep -v '^#' .env.test | xargs) && \
		migrate -path migrations -database "$$TEST_DATABASE_URL" up
	@echo "Test database reset complete!"

test: test-db-setup ## Run all tests (requires database)
	@echo "Running all tests..."
	@cd apps/backend && \
		export $$(grep -v '^#' .env.test | xargs) && \
		go test -v ./...

test-unit: ## Run only unit tests (no database required)
	@echo "Running unit tests..."
	cd apps/backend && go test -v -short ./...

test-integration: test-db-setup ## Run only integration tests (requires database)
	@echo "Running integration tests..."
	@cd apps/backend && \
		export $$(grep -v '^#' .env.test | xargs) && \
		go test -v -run TestSuite ./...

test-coverage: test-db-setup ## Run tests with coverage report
	@echo "Running tests with coverage..."
	@cd apps/backend && \
		export $$(grep -v '^#' .env.test | xargs) && \
		go test -cover ./...
	@echo ""
	@echo "For detailed coverage report, run: make test-coverage-html"

test-coverage-html: test-db-setup ## Generate HTML coverage report
	@echo "Generating HTML coverage report..."
	@cd apps/backend && \
		export $$(grep -v '^#' .env.test | xargs) && \
		go test -coverprofile=coverage.out ./...
	@cd apps/backend && go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: apps/backend/coverage.html"
