.PHONY: help setup install-swag swagger build up down restart logs clean test quick-start dev

# Default target
help:
	@echo "🏗️  Cooperative ERP Lite - Makefile Commands"
	@echo ""
	@echo "Quick Start:"
	@echo "  make quick-start    - Complete setup and start services"
	@echo ""
	@echo "Setup Commands:"
	@echo "  make setup          - Initial project setup"
	@echo "  make install-swag   - Install Swagger CLI"
	@echo "  make swagger        - Generate Swagger documentation"
	@echo ""
	@echo "Docker Commands:"
	@echo "  make build          - Build Docker images"
	@echo "  make up             - Start all services"
	@echo "  make down           - Stop all services"
	@echo "  make restart        - Restart all services"
	@echo "  make logs           - View service logs"
	@echo ""
	@echo "Development Commands:"
	@echo "  make dev            - Run backend in development mode"
	@echo "  make test           - Run tests"
	@echo "  make clean          - Clean build artifacts"
	@echo ""
	@echo "Database Commands:"
	@echo "  make db-shell       - Connect to PostgreSQL shell"
	@echo "  make db-migrate     - Run database migrations"
	@echo ""

# Quick start - setup and run everything
quick-start: setup build up
	@echo "✅ Cooperative ERP Lite is running!"
	@echo "📍 Backend API: http://localhost:8080"
	@echo "📍 Swagger Docs: http://localhost:8080/swagger/index.html"
	@echo ""
	@echo "To view logs: make logs"
	@echo "To stop: make down"

# Setup project
setup:
	@echo "⚙️  Setting up project..."
	@echo "1. Copying environment file..."
	@if [ ! -f backend/.env ]; then \
		if [ -f backend/.env.example ]; then \
			cp backend/.env.example backend/.env; \
			echo "   ✅ Created .env from .env.example"; \
		else \
			echo "   ℹ️  No .env.example found, skipping"; \
		fi \
	else \
		echo "   ℹ️  .env already exists"; \
	fi
	@echo "2. Installing Go dependencies..."
	@cd backend && go mod download
	@echo "3. Installing swag CLI..."
	@$(MAKE) install-swag
	@echo "4. Generating Swagger docs..."
	@$(MAKE) swagger
	@echo "✅ Setup complete!"

# Install Swagger CLI
install-swag:
	@if ! command -v swag >/dev/null 2>&1; then \
		echo "   📦 Installing swag..."; \
		go install github.com/swaggo/swag/cmd/swag@latest; \
	else \
		echo "   ✅ swag already installed"; \
	fi

# Generate Swagger documentation
swagger:
	@echo "📝 Generating Swagger documentation..."
	@if [ -d backend/cmd/api ]; then \
		cd backend && swag init -g cmd/api/main.go -o docs; \
	else \
		echo "   ⚠️  backend/cmd/api/main.go not found yet"; \
		echo "   Skipping Swagger generation"; \
	fi

# Build Docker images
build:
	@echo "🐳 Building Docker images..."
	@docker compose build

# Start services
up:
	@echo "🚀 Starting services..."
	@docker compose up -d
	@echo "⏳ Waiting for services to be ready..."
	@sleep 5
	@docker compose ps

# Stop services
down:
	@echo "🛑 Stopping services..."
	@docker compose down

# Restart services
restart: down up

# View logs
logs:
	@docker compose logs -f

# Clean build artifacts
clean:
	@echo "🧹 Cleaning build artifacts..."
	@rm -rf backend/bin
	@rm -rf backend/docs
	@docker compose down -v
	@docker system prune -f
	@echo "✅ Clean complete!"

# Run tests
test:
	@echo "🧪 Running tests..."
	@cd backend && go test -v ./...

# Development mode (run backend without Docker)
dev:
	@echo "🔧 Starting backend in development mode..."
	@if [ ! -f backend/.env ]; then \
		echo "❌ Error: backend/.env not found. Run 'make setup' first."; \
		exit 1; \
	fi
	@cd backend && go run cmd/api/main.go

# Database shell
db-shell:
	@docker compose exec postgres psql -U postgres -d koperasi_erp

# Database migrations (placeholder)
db-migrate:
	@echo "📊 Running database migrations..."
	@echo "⚠️  Migration system not implemented yet"
	@echo "   Migrations will be handled by GORM Auto Migrate in the application"

# Check Go version
check-go:
	@echo "🔍 Checking Go version..."
	@go version | grep -q "go1.25" || (echo "❌ Go 1.25.x required" && exit 1)
	@echo "✅ Go version OK"

# Install development tools
install-tools: install-swag
	@echo "🔧 Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@echo "✅ Tools installed!"

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@cd backend && go fmt ./...
	@echo "✅ Code formatted!"

# Lint code (requires golangci-lint)
lint:
	@echo "🔍 Linting code..."
	@if command -v golangci-lint >/dev/null 2>&1; then \
		cd backend && golangci-lint run; \
	else \
		echo "⚠️  golangci-lint not installed"; \
		echo "   Install with: brew install golangci-lint"; \
	fi

# Show service status
status:
	@echo "📊 Service Status:"
	@docker compose ps
