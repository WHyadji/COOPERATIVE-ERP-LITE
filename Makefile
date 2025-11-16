.PHONY: help build up down restart logs clean seed swagger test

# Default target
.DEFAULT_GOAL := help

##@ General

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Docker Operations

build: ## Build Docker images
	@echo "🔨 Building Docker images..."
	docker-compose build

up: ## Start all services (detached mode)
	@echo "🚀 Starting services..."
	docker-compose up -d
	@echo "✅ Services started!"
	@echo "   - Backend API: http://localhost:8080"
	@echo "   - Swagger UI: http://localhost:8080/swagger/index.html"
	@echo "   - Adminer: http://localhost:8081"
	@echo "   - PostgreSQL: localhost:5432"

down: ## Stop all services
	@echo "🛑 Stopping services..."
	docker-compose down
	@echo "✅ Services stopped!"

restart: down up ## Restart all services

logs: ## View logs (follow mode)
	docker-compose logs -f

logs-backend: ## View backend logs only
	docker-compose logs -f backend

logs-postgres: ## View postgres logs only
	docker-compose logs -f postgres

ps: ## Show running containers
	docker-compose ps

##@ Development

dev: ## Run backend in development mode (with auto-reload)
	@echo "🔥 Starting development server..."
	cd backend && air

run: ## Run backend directly (without Docker)
	@echo "▶️  Running backend server..."
	cd backend && go run cmd/api/main.go

swagger: ## Generate Swagger documentation
	@echo "📚 Generating Swagger documentation..."
	cd backend && swag init -g cmd/api/main.go -o docs
	@echo "✅ Swagger documentation generated!"
	@echo "   View at: http://localhost:8080/swagger/index.html"

seed: ## Run database seed
	@echo "🌱 Seeding database..."
	cd backend && go run cmd/seed/main.go
	@echo "✅ Database seeded successfully!"

##@ Database Operations

db-create: ## Create database (if not exists)
	docker-compose exec postgres createdb -U postgres koperasi_erp

db-drop: ## Drop database (WARNING: destructive!)
	@echo "⚠️  WARNING: This will delete all data!"
	@read -p "Are you sure? [y/N]: " confirm && \
	if [ "$$confirm" = "y" ] || [ "$$confirm" = "Y" ]; then \
		docker-compose exec postgres dropdb -U postgres koperasi_erp; \
		echo "✅ Database dropped!"; \
	else \
		echo "❌ Cancelled!"; \
	fi

db-connect: ## Connect to PostgreSQL shell
	docker-compose exec postgres psql -U postgres -d koperasi_erp

db-backup: ## Backup database to file
	@echo "💾 Creating database backup..."
	@mkdir -p backups
	docker-compose exec -T postgres pg_dump -U postgres koperasi_erp > backups/backup_$(shell date +%Y%m%d_%H%M%S).sql
	@echo "✅ Backup created in backups/ directory"

db-restore: ## Restore database from latest backup
	@echo "📥 Restoring database from latest backup..."
	@LATEST=$$(ls -t backups/*.sql | head -1); \
	if [ -z "$$LATEST" ]; then \
		echo "❌ No backup files found!"; \
	else \
		echo "Restoring from: $$LATEST"; \
		docker-compose exec -T postgres psql -U postgres koperasi_erp < "$$LATEST"; \
		echo "✅ Database restored!"; \
	fi

##@ Testing

test: ## Run all tests
	@echo "🧪 Running tests..."
	cd backend && go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "🧪 Running tests with coverage..."
	cd backend && go test -v -coverprofile=coverage.out ./...
	cd backend && go tool cover -html=coverage.out -o coverage.html
	@echo "✅ Coverage report: backend/coverage.html"

bench: ## Run benchmarks
	@echo "⚡ Running benchmarks..."
	cd backend && go test -bench=. -benchmem ./...

##@ Code Quality

lint: ## Run linter
	@echo "🔍 Running linter..."
	cd backend && golangci-lint run ./...

fmt: ## Format code
	@echo "✨ Formatting code..."
	cd backend && go fmt ./...

vet: ## Run go vet
	@echo "🔎 Running go vet..."
	cd backend && go vet ./...

tidy: ## Tidy go modules
	@echo "🧹 Tidying go modules..."
	cd backend && go mod tidy

##@ Cleanup

clean: ## Clean build artifacts and Docker volumes
	@echo "🧹 Cleaning up..."
	docker-compose down -v
	rm -rf backend/bin
	rm -rf backend/coverage.out backend/coverage.html
	@echo "✅ Cleanup complete!"

clean-all: clean ## Clean everything including Docker images
	@echo "🧹 Deep cleaning..."
	docker-compose down -v --rmi all
	@echo "✅ Deep cleanup complete!"

##@ Setup

setup: ## Initial project setup
	@echo "⚙️  Setting up project..."
	@echo "1. Copying environment file..."
	@if [ ! -f backend/.env ]; then \
		cp backend/.env.example backend/.env; \
		echo "   ✅ .env file created"; \
	else \
		echo "   ℹ️  .env already exists"; \
	fi
	@echo "2. Installing Go dependencies..."
	cd backend && go mod download
	@echo "3. Installing swag CLI..."
	go install github.com/swaggo/swag/cmd/swag@latest
	@echo "4. Generating Swagger docs..."
	cd backend && swag init -g cmd/api/main.go -o docs
	@echo "5. Building Docker images..."
	docker-compose build
	@echo ""
	@echo "✅ Setup complete!"
	@echo ""
	@echo "Next steps:"
	@echo "  1. Update backend/.env with your configuration"
	@echo "  2. Run 'make up' to start services"
	@echo "  3. Run 'make seed' to populate with sample data"
	@echo "  4. Visit http://localhost:8080/swagger/index.html"

##@ Quick Start

quick-start: setup up seed ## Complete setup and start (for first-time users)
	@echo ""
	@echo "🎉 Everything is ready!"
	@echo ""
	@echo "📡 Services:"
	@echo "   - API: http://localhost:8080/api/v1"
	@echo "   - Swagger UI: http://localhost:8080/swagger/index.html"
	@echo "   - Adminer: http://localhost:8081"
	@echo ""
	@echo "🔑 Login Credentials:"
	@echo "   - Admin: username=admin, password=admin123"
	@echo "   - Bendahara: username=bendahara, password=bendahara123"
	@echo "   - Kasir: username=kasir, password=kasir123"
	@echo ""
	@echo "📚 Try the API:"
	@echo "   curl http://localhost:8080/health"
	@echo "   curl -X POST http://localhost:8080/api/v1/auth/login \\"
	@echo "        -H 'Content-Type: application/json' \\"
	@echo "        -d '{\"namaPengguna\":\"admin\",\"kataSandi\":\"admin123\"}'"

##@ Info

version: ## Show version info
	@echo "📦 Cooperative ERP Lite"
	@echo "   Go: $$(go version)"
	@echo "   Docker: $$(docker --version)"
	@echo "   Docker Compose: $$(docker-compose --version)"

status: ps ## Alias for ps (show running containers)
