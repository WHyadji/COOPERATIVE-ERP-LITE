# Backend - Cooperative ERP Lite

Backend API server for Indonesia's first cooperative operating system - a digital platform designed specifically for the Indonesian cooperative model.

## Overview

This is a RESTful API built with Go and the Gin framework, providing backend services for the Cooperative ERP Lite platform. The backend handles authentication, member management, financial transactions, accounting, POS operations, and reporting for Indonesian cooperatives.

**Philosophy**: "Better than paper books = WIN" - We're competing with manual record-keeping, so even 20% improvement is a major win.

## Prerequisites

- **Go**: 1.25.4 or later ([Download](https://golang.org/dl/))
- **PostgreSQL**: 15 or later
- **Docker**: Optional, for running PostgreSQL in a container
- **Air**: Optional, for hot-reload during development
  ```bash
  go install github.com/cosmtrek/air@latest
  ```

## Technology Stack

### Core Framework
- **Gin v1.11.0** - HTTP web framework
- **GORM v1.31.1** - ORM for database operations
- **PostgreSQL Driver v1.6.0** - Database driver

### Authentication & Security
- **JWT v5.3.0** - JSON Web Tokens for authentication
- **bcrypt** - Password hashing (via golang.org/x/crypto v0.44.0)

### Utilities
- **UUID v1.6.0** - Unique identifier generation
- **godotenv v1.5.1** - Environment variable management
- **CORS v1.7.6** - Cross-Origin Resource Sharing middleware
- **Validator v10.27.0** - Request validation

### Testing
- **testify v1.11.1** - Testing framework with assertions

## Project Structure

```
backend/
├── cmd/
│   ├── api/              # API server entry point
│   │   └── main.go
│   └── hashpass/         # Password hashing utility
│       └── main.go
├── internal/
│   ├── config/           # Configuration management
│   │   ├── config.go
│   │   └── database.go
│   ├── constants/        # Application constants
│   │   └── messages.go
│   ├── errors/           # Custom error types
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware (auth, CORS, logging)
│   │   ├── auth.go
│   │   ├── cors.go
│   │   └── logger.go
│   ├── models/           # Database models (GORM)
│   ├── services/         # Business logic layer
│   └── utils/            # Utility functions
│       ├── errors.go
│       ├── jwt.go
│       ├── logger.go
│       ├── pagination.go
│       ├── password.go
│       └── response.go
├── pkg/
│   └── validasi/         # Custom validation logic
├── docs/                 # API documentation & guides
│   ├── swagger.yaml      # OpenAPI specification
│   ├── swagger.json
│   ├── docs.go
│   ├── HANDLER_IMPLEMENTATION_GUIDE.md
│   └── MULTI_TENANT_SECURITY_REVIEW.md
├── bin/                  # Compiled binaries (gitignored)
├── .env.example          # Environment variables template
├── .gitignore
├── go.mod                # Go module dependencies
├── go.sum                # Dependency checksums
└── README.md             # This file
```

## Quick Start

### 1. Clone and Navigate

```bash
cd backend
```

### 2. Install Dependencies

```bash
go mod download
```

### 3. Set Up Database

**Option A: Using Docker (Recommended)**

```bash
# Start PostgreSQL container
docker run --name koperasi-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=koperasi_erp \
  -p 5432:5432 \
  -d postgres:15

# Verify it's running
docker ps
```

**Option B: Local PostgreSQL Installation**

Install PostgreSQL 15 and create a database named `koperasi_erp`.

### 4. Configure Environment

```bash
# Copy example environment file
cp .env.example .env

# Edit .env with your database credentials and JWT secret
nano .env
```

**Required environment variables:**
```env
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=koperasi_erp
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-super-secret-key-change-this-in-production

# Server
PORT=8080
GIN_MODE=debug
```

### 5. Run Database Migrations

The application automatically runs migrations on startup. Models are defined in `internal/models/`.

### 6. Start the Server

**Development mode (with hot-reload using Air):**

```bash
air
```

**Or run directly:**

```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## Development

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/services/...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Formatting

```bash
# Format all Go files
go fmt ./...

# Check for common mistakes
go vet ./...
```

### Building

```bash
# Build for current platform
go build -o bin/accounting-ledger cmd/api/main.go

# Run the binary
./bin/accounting-ledger

# Build for production (with optimizations)
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/accounting-ledger cmd/api/main.go
```

### Database Management

```bash
# Connect to database (Docker)
docker exec -it koperasi-postgres psql -U postgres -d koperasi_erp

# Backup database
docker exec koperasi-postgres pg_dump -U postgres koperasi_erp > backup.sql

# Restore database
docker exec -i koperasi-postgres psql -U postgres koperasi_erp < backup.sql

# View database logs
docker logs koperasi-postgres
```

### Useful Commands

```bash
# Hash a password for testing
go run cmd/hashpass/main.go "your-password"

# List all dependencies
go list -m all

# Update dependencies
go get -u ./...
go mod tidy

# Check for dependency vulnerabilities
go list -json -deps | grep "CVE"
```

## API Documentation

### Available Endpoints

The API follows RESTful conventions with the following main routes:

#### Authentication
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `POST /api/auth/logout` - User logout
- `GET /api/auth/me` - Get current user info

#### Members
- `GET /api/members` - List all members
- `GET /api/members/:id` - Get member details
- `POST /api/members` - Create new member
- `PUT /api/members/:id` - Update member
- `DELETE /api/members/:id` - Delete member

#### Share Capital
- `GET /api/share-capital` - List share capital records
- `POST /api/share-capital` - Record share capital transaction

#### Transactions
- `GET /api/transactions` - List transactions
- `POST /api/transactions` - Create transaction

#### POS
- `GET /api/products` - List products
- `POST /api/sales` - Record sale
- `GET /api/sales/:id` - Get sale details

#### Reports
- `GET /api/reports/balance-sheet` - Generate balance sheet
- `GET /api/reports/profit-loss` - Generate P&L statement
- `GET /api/reports/member-balances` - Get member balances

### Swagger Documentation

API documentation is available via Swagger UI at:
- **Development**: `http://localhost:8080/swagger/index.html`

To regenerate Swagger docs:
```bash
# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Generate docs
swag init -g cmd/api/main.go -o docs/
```

### Authentication

All protected endpoints require a JWT token. Include it in the Authorization header:

```bash
Authorization: Bearer <your-jwt-token>
```

Tokens are issued upon successful login and expire after 24 hours.

## Multi-Tenancy

This application is **multi-tenant** by design. All data queries are automatically filtered by `cooperative_id` to ensure data isolation between cooperatives.

**Important**: Always include the cooperative context in middleware to prevent data leakage.

See `docs/MULTI_TENANT_SECURITY_REVIEW.md` for security guidelines.

## Testing Strategy

### Unit Tests
- Service layer functions (`internal/services/*_test.go`)
- Utility functions (`internal/utils/*_test.go`)
- Data validation logic

### Integration Tests
- API endpoints with database
- Authentication flows
- Multi-tenant data isolation

### Test Coverage Goals
- **Critical paths**: 80%+ coverage
- **Business logic**: 70%+ coverage
- **Overall**: 60%+ coverage

## Security Considerations

### Password Security
- All passwords are hashed using bcrypt
- Default cost factor: 10
- Never log or expose passwords

### JWT Security
- Tokens stored in HTTP-only cookies
- 24-hour expiration
- Include cooperative_id in claims

### Multi-Tenant Security
- All queries filtered by cooperative_id
- Middleware validates tenant access
- No cross-tenant data access

### Input Validation
- Request validation using validator/v10
- SQL injection prevention via GORM parameterized queries
- XSS prevention in JSON responses

## Performance Optimization

### Database
- Indexes on frequently queried fields
- Connection pooling (configured in `config/database.go`)
- Prepared statements via GORM

### Caching
- Future: Redis for session management
- Future: Query result caching for reports

## Deployment

### Environment Variables (Production)
```env
GIN_MODE=release
DB_SSLMODE=require
JWT_SECRET=<strong-random-secret>
ALLOWED_ORIGINS=https://yourdomain.com
```

### Cloud Run Deployment (GCP)

```bash
# Build container
docker build -t gcr.io/your-project/cooperative-erp-backend .

# Push to registry
docker push gcr.io/your-project/cooperative-erp-backend

# Deploy to Cloud Run
gcloud run deploy cooperative-erp-backend \
  --image gcr.io/your-project/cooperative-erp-backend \
  --platform managed \
  --region asia-southeast2 \
  --allow-unauthenticated \
  --set-env-vars="DB_HOST=...,DB_NAME=...,JWT_SECRET=..."
```

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps

# Check PostgreSQL logs
docker logs koperasi-postgres

# Test connection
psql -h localhost -U postgres -d koperasi_erp
```

### Port Already in Use

```bash
# Find process using port 8080
lsof -i :8080

# Kill the process
kill -9 <PID>
```

### Module Dependencies

```bash
# Clean module cache
go clean -modcache

# Re-download dependencies
go mod download
```

## Indonesian Cooperative Context

This system is designed for Indonesian cooperatives following:
- **Legal Framework**: UU No. 25 Tahun 1992
- **Accounting Standards**: SAK ETAP (Indonesian GAAP)
- **Share Capital Types**:
  - Simpanan Pokok (Principal deposit)
  - Simpanan Wajib (Mandatory deposit)
  - Simpanan Sukarela (Voluntary deposit)

## Contributing

### Code Style
- Follow Go best practices and idioms
- Use `gofmt` for formatting
- Write meaningful commit messages
- Add tests for new features

### Pull Request Process
1. Create a feature branch
2. Write tests for your changes
3. Ensure all tests pass
4. Update documentation
5. Submit PR with clear description

## Support

For issues, questions, or contributions:
- Check existing documentation in `docs/`
- Review implementation guides
- Follow the 12-week MVP action plan

## License

Copyright (c) 2025 Cooperative ERP Lite

## Project Status

**Current Phase**: Pre-Development (Week 0)
**Target**: MVP Launch in 12 weeks
**Goal**: 10 pilot cooperatives

**MVP Features (Locked Scope)**:
1. Authentication & Authorization
2. Member Management
3. Share Capital Tracking
4. Basic POS (cash-only)
5. Simple Accounting (manual entries)
6. Essential Reports
7. Member Portal
8. Data Import (Excel templates)

Features explicitly deferred to Phase 2:
- SHU calculation
- QRIS/digital payments
- Barcode scanning
- Receipt printing
- Inventory automation
- WhatsApp integration
- Native mobile apps
- Offline mode
