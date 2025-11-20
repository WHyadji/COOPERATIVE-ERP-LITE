# Codebase Organization Guide

**Last Updated**: November 18, 2025
**Version**: 1.0.0

This document provides a comprehensive overview of the Cooperative ERP Lite codebase structure, organization, and navigation guidelines.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Directory Structure](#directory-structure)
3. [Backend Architecture](#backend-architecture)
4. [Frontend Architecture](#frontend-architecture)
5. [Documentation](#documentation)
6. [Scripts & Utilities](#scripts--utilities)
7. [Configuration Files](#configuration-files)
8. [Development Workflow](#development-workflow)
9. [Navigation Guide](#navigation-guide)
10. [Best Practices](#best-practices)

---

## Project Overview

**Cooperative ERP Lite** is Indonesia's first cooperative operating system - a comprehensive digital platform for Indonesian cooperatives.

### Project Stats
- **Backend**: 19,305 lines of Go code
- **Frontend**: 4,606 lines of TypeScript/React
- **Documentation**: 64,100+ lines of markdown
- **Tests**: 22 test files with comprehensive coverage
- **Scripts**: 12 utility scripts organized in 3 categories

### Technology Stack
- **Backend**: Go 1.25.4 + Gin + GORM + PostgreSQL
- **Frontend**: Next.js 15.5.0 + React 19 + TypeScript + MUI
- **Infrastructure**: Docker + PostgreSQL + Nginx

---

## Directory Structure

```
COOPERATIVE-ERP-LITE/
├── backend/                 # Go backend application
│   ├── bin/                # Compiled binaries
│   ├── cmd/                # Application entry points
│   ├── internal/           # Private application code
│   ├── pkg/                # Public/shared packages
│   └── docs/               # Swagger/OpenAPI documentation
│
├── frontend/               # Next.js frontend application
│   ├── app/                # Next.js App Router pages
│   ├── components/         # React components
│   ├── lib/                # Utility libraries
│   ├── types/              # TypeScript type definitions
│   └── __tests__/          # Frontend tests
│
├── docs/                   # Comprehensive project documentation
│   ├── business/           # Business documentation
│   ├── phases/             # Phase-specific guides
│   └── *.md                # Core documentation files
│
├── scripts/                # Organized utility scripts
│   ├── seed/               # Database seeding scripts
│   ├── test/               # Testing scripts
│   ├── debug/              # Debugging utilities
│   └── README.md           # Scripts documentation
│
├── nginx/                  # Nginx configuration
│   ├── conf.d/             # Server configurations
│   ├── logs/               # Nginx logs
│   └── ssl/                # SSL certificates
│
├── .claude/                # Claude Code configuration
│   ├── agents/             # Specialized agents
│   ├── commands/           # Custom slash commands
│   ├── skills/             # Reusable skills
│   └── docs/               # Claude-specific docs
│
├── .github/                # GitHub configuration
│   ├── workflows/          # CI/CD pipelines
│   └── ISSUE_TEMPLATE/     # Issue templates
│
└── [Root configuration files]
```

---

## Backend Architecture

### Structure Overview

```
backend/
├── cmd/
│   ├── api/main.go         # API server entry point (HTTP server)
│   └── hashpass/main.go    # Password hashing utility
│
├── internal/               # Private application code (not importable)
│   ├── config/
│   │   ├── config.go       # Configuration management
│   │   └── database.go     # Database setup & connection
│   │
│   ├── constants/
│   │   ├── constants.go    # Application constants
│   │   └── messages.go     # Response messages (Indonesian)
│   │
│   ├── errors/
│   │   └── errors.go       # Custom error definitions
│   │
│   ├── handlers/           # HTTP request handlers (10 handlers)
│   │   ├── anggota_handler.go      # Member management
│   │   ├── auth_handler.go         # Authentication
│   │   ├── koperasi_handler.go     # Cooperative management
│   │   ├── akun_handler.go         # Chart of accounts
│   │   ├── transaksi_handler.go    # Journal entries
│   │   ├── pengguna_handler.go     # User management
│   │   ├── simpanan_handler.go     # Share capital
│   │   ├── produk_handler.go       # Product catalog
│   │   ├── penjualan_handler.go    # POS transactions
│   │   └── laporan_handler.go      # Reports
│   │
│   ├── middleware/
│   │   ├── auth.go         # JWT authentication middleware
│   │   ├── cors.go         # CORS middleware
│   │   └── logger.go       # Request/response logging
│   │
│   ├── models/             # Data models (8 models)
│   │   ├── anggota.go      # Cooperative member
│   │   ├── koperasi.go     # Cooperative organization
│   │   ├── pengguna.go     # System user with roles
│   │   ├── akun.go         # Chart of accounts
│   │   ├── transaksi.go    # Journal entry + line items
│   │   ├── simpanan.go     # Share capital deposit
│   │   ├── produk.go       # Product for POS
│   │   └── penjualan.go    # POS sale + line items
│   │
│   ├── services/           # Business logic layer (10 services)
│   │   ├── auth_service.go         # Authentication & JWT
│   │   ├── anggota_service.go      # Member CRUD
│   │   ├── koperasi_service.go     # Cooperative CRUD
│   │   ├── akun_service.go         # Chart of accounts
│   │   ├── transaksi_service.go    # Double-entry accounting
│   │   ├── pengguna_service.go     # User management
│   │   ├── simpanan_service.go     # Deposit tracking
│   │   ├── produk_service.go       # Product management
│   │   ├── penjualan_service.go    # Sales transactions
│   │   └── laporan_service.go      # Report generation
│   │
│   └── utils/
│       ├── jwt.go          # JWT token utilities
│       ├── password.go     # Password hashing (bcrypt)
│       ├── logger.go       # Logging utilities
│       └── response.go     # HTTP response formatting
│
├── pkg/                    # Public packages (importable)
│   └── validasi/
│       ├── validator.go    # Custom validation rules
│       └── validation-implementation-guide.md
│
└── docs/                   # Generated API documentation
    ├── swagger.yaml        # OpenAPI specification
    └── docs.go             # Swagger Go code
```

### Key Backend Components

#### 1. **Models** (Data Layer)
- Use GORM for ORM
- All models include `cooperative_id` for multi-tenancy
- Soft delete support with `deleted_at`
- UUID primary keys
- Audit timestamps (`created_at`, `updated_at`)

#### 2. **Services** (Business Logic)
- Handle all business logic
- Database transactions
- Data validation
- Error handling
- Multi-tenant filtering

#### 3. **Handlers** (HTTP Layer)
- HTTP request/response handling
- Request validation
- JWT token extraction
- Role-based access control
- Error response formatting

#### 4. **Middleware**
- **AuthMiddleware**: JWT validation, token parsing
- **CORSMiddleware**: Cross-origin request handling
- **LoggerMiddleware**: Request/response logging
- **RequireRole**: Role-based authorization

### API Endpoints

40+ REST endpoints organized under `/api/v1`:

```
POST   /api/v1/auth/login          # Login
POST   /api/v1/auth/logout         # Logout
GET    /api/v1/auth/profile        # Get profile
PUT    /api/v1/auth/password       # Change password

GET    /api/v1/anggota             # List members
POST   /api/v1/anggota             # Create member
GET    /api/v1/anggota/:id         # Get member
PUT    /api/v1/anggota/:id         # Update member
DELETE /api/v1/anggota/:id         # Delete member

GET    /api/v1/akun                # List accounts
POST   /api/v1/akun                # Create account
...

GET    /api/v1/laporan/neraca      # Balance sheet
GET    /api/v1/laporan/laba-rugi   # Income statement
...
```

Full API documentation: `/backend/docs/swagger.yaml`

---

## Frontend Architecture

### Structure Overview

```
frontend/
├── app/                    # Next.js App Router
│   ├── (auth)/
│   │   └── login/
│   │       └── page.tsx    # Login page
│   │
│   ├── (dashboard)/        # Protected routes
│   │   ├── layout.tsx      # Dashboard layout
│   │   ├── page.tsx        # Dashboard home
│   │   │
│   │   ├── members/        # Member management
│   │   │   ├── page.tsx    # Member list
│   │   │   ├── new/page.tsx # Create member
│   │   │   └── [id]/page.tsx # Member detail
│   │   │
│   │   └── simpanan/       # Deposit management
│   │       ├── page.tsx    # Deposit list
│   │       ├── new/page.tsx # Record deposit
│   │       └── saldo/page.tsx # Balance summary
│   │
│   ├── layout.tsx          # Root layout with AuthProvider
│   └── page.tsx            # Root redirect page
│
├── components/             # React components
│   └── layout/
│       ├── Header.tsx      # Navigation header
│       └── Sidebar.tsx     # Navigation sidebar
│
├── lib/                    # Utility libraries
│   ├── api/
│   │   ├── client.ts       # Axios API client
│   │   ├── memberApi.ts    # Member API endpoints
│   │   └── simpananApi.ts  # Deposit API endpoints
│   │
│   └── context/
│       └── AuthContext.tsx # Authentication state
│
├── types/                  # TypeScript definitions
│   ├── member.ts
│   ├── user.ts
│   └── api.ts
│
└── __tests__/              # Frontend tests
    ├── components/
    └── utils/
```

### Key Frontend Components

#### 1. **App Router** (Next.js 15)
- Server Components by default
- Client Components with 'use client'
- Route groups for organization: `(auth)`, `(dashboard)`
- Dynamic routes: `[id]`

#### 2. **Authentication**
- Context-based auth state: `AuthContext`
- JWT stored in HTTP-only cookies
- Automatic token refresh
- Protected routes with middleware

#### 3. **API Integration**
- Axios client with interceptors
- Request/response transformation
- Error handling
- Type-safe API calls

#### 4. **UI Components**
- Material-UI (MUI) components
- Custom layout components (Header, Sidebar)
- Form components with React Hook Form + Zod
- Responsive design (mobile-first)

### Pages Implemented

| Route | Purpose | Status |
|-------|---------|--------|
| `/` | Root redirect | ✅ Implemented |
| `/login` | Login page | ✅ Implemented |
| `/dashboard` | Dashboard home | ✅ Implemented |
| `/members` | Member list | ✅ Implemented |
| `/members/new` | Create member | ✅ Implemented |
| `/members/[id]` | Member detail | ✅ Implemented |
| `/simpanan` | Deposit list | ✅ Implemented |
| `/simpanan/new` | Record deposit | ✅ Implemented |
| `/simpanan/saldo` | Balance summary | ✅ Implemented |

---

## Documentation

### Documentation Structure

```
docs/
├── INDEX.md                        # Documentation hub
├── README.md                       # Docs overview
│
├── Core Documentation
├── architecture.md                 # System architecture (597 lines)
├── mvp-action-plan.md             # 12-week sprint plan (790 lines)
├── quick-start-guide.md           # Setup guide (752 lines)
├── technical-stack-versions.md    # Locked versions (590 lines)
├── post-mvp-roadmap.md            # Phase 2-6 features (2,491 lines)
│
├── business/                       # Business documentation
│   ├── business-model.md          # Business model (741 lines)
│   ├── market-analysis.md         # Market research (633 lines)
│   ├── go-to-market-strategy.md   # GTM strategy (1,270 lines)
│   ├── product-positioning.md     # Positioning (647 lines)
│   ├── clarifications.md          # Q&A (918 lines)
│   └── README.md                  # Business docs index
│
└── phases/                         # Phase-specific guides
    ├── phase-1/
    │   ├── implementation-guide.md # Week-by-week plan (3,011 lines)
    │   ├── progress-tracking.md    # Progress metrics
    │   └── README.md
    │
    └── phase-2/
        └── README.md
```

### Documentation Categories

#### 1. **Core Documentation**
Essential documents for understanding and working with the project:
- **CLAUDE.md**: Instructions for Claude Code
- **README.md**: Project overview
- **DOCKER-SETUP.md**: Docker guide
- **CONTRIBUTING.md**: Contribution guidelines

#### 2. **Technical Documentation**
- **architecture.md**: System design, patterns, technology choices
- **technical-stack-versions.md**: Exact dependency versions (LOCKED)
- **quick-start-guide.md**: Development environment setup

#### 3. **Project Planning**
- **mvp-action-plan.md**: 12-week sprint to MVP launch
- **post-mvp-roadmap.md**: Long-term roadmap (Phase 2-6)
- **implementation-guide.md**: Detailed week-by-week tasks

#### 4. **Business Documentation**
- Market analysis and competitive landscape
- Business model and revenue streams
- Go-to-market strategy
- Product positioning

#### 5. **API Documentation**
- **backend/docs/swagger.yaml**: OpenAPI specification
- Auto-generated from code comments
- Includes request/response examples

---

## Scripts & Utilities

All scripts are organized in the `scripts/` directory with clear categorization.

### Directory Structure

```
scripts/
├── seed/           # Database seeding
├── test/           # Testing utilities
├── debug/          # Debugging tools
└── README.md       # Comprehensive guide
```

### Seed Scripts

Located in `scripts/seed/`:

| Script | Purpose |
|--------|---------|
| `seed-initial.sh` | Seed initial cooperative, users, and basic data |
| `seed-coa.sh` | Seed Chart of Accounts (SAK ETAP compliant) |
| `test-data.sh` | Generate test data for development |
| `create-test-data-sequential.sh` | Create test data sequentially |
| `cleanup-test-data.sh` | Remove all test data from database |

**Usage**:
```bash
./scripts/seed/seed-initial.sh
./scripts/seed/seed-coa.sh
./scripts/seed/test-data.sh
```

### Test Scripts

Located in `scripts/test/`:

| Script | Purpose |
|--------|---------|
| `test-member.sh` | Test member management endpoints |
| `test-deposit.sh` | Test simpanan (deposit) endpoints |
| `test-concurrent.sh` | Test concurrent transactions and race conditions |
| `test-api-summary.sh` | Run comprehensive API test suite |
| `test-edge-cases.sh` | Test edge cases and error handling |

**Usage**:
```bash
./scripts/test/test-api-summary.sh    # Full test suite
./scripts/test/test-member.sh         # Specific tests
```

### Debug Scripts

Located in `scripts/debug/`:

| Script | Purpose |
|--------|---------|
| `debug-member.sh` | Debug member-related issues |
| `check-balance.sh` | Check account balances and trial balance |

**Usage**:
```bash
./scripts/debug/check-balance.sh
./scripts/debug/debug-member.sh
```

### Script Documentation

Full documentation for all scripts: `scripts/README.md`

---

## Configuration Files

### Root Level Configuration

```
.
├── .gitignore              # Git ignore rules
├── .env                    # Backend environment variables (not in git)
├── docker-compose.yml      # Docker orchestration
├── Dockerfile              # Backend Docker image
├── Makefile                # Build & development commands
├── go.mod                  # Go module dependencies
└── go.sum                  # Go dependency checksums
```

### Backend Configuration

```
backend/
├── .env                    # Backend environment variables
└── cmd/api/main.go         # Configuration loading
```

**Environment Variables**:
```bash
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=koperasi_erp
DB_SSLMODE=disable

# JWT
JWT_SECRET=your-secret-key-here

# Server
PORT=8080
GIN_MODE=debug
```

### Frontend Configuration

```
frontend/
├── .env.local              # Frontend environment variables (not in git)
├── next.config.ts          # Next.js configuration
├── tsconfig.json           # TypeScript configuration
├── package.json            # NPM dependencies
└── vitest.config.ts        # Test configuration
```

**Environment Variables**:
```bash
NEXT_PUBLIC_API_URL=http://localhost:8080/api/v1
```

### Docker Configuration

```
nginx/
├── nginx.conf              # Main Nginx config
├── conf.d/
│   └── default.conf        # Server block configuration
└── ssl/                    # SSL certificates
```

### Claude Code Configuration

```
.claude/
├── config.json             # Claude CLI settings
├── agents/                 # 16+ specialized agents
├── commands/               # 15+ custom slash commands
├── skills/                 # 8 reusable skills
└── docs/                   # Claude-specific docs
```

---

## Development Workflow

### Quick Start

```bash
# 1. Clone repository
git clone <repository-url>
cd COOPERATIVE-ERP-LITE

# 2. Start services with Docker
make quick-start

# 3. Seed database
./scripts/seed/seed-initial.sh
./scripts/seed/seed-coa.sh

# 4. Access application
# Frontend: http://localhost:3000
# Backend:  http://localhost:8080
# API Docs: http://localhost:8080/swagger/index.html
```

### Common Development Commands

```bash
# Development
make dev                    # Run backend in development mode
make dev-frontend          # Run frontend in development mode
make dev-all               # Run both backend and frontend

# Docker
make build                 # Build Docker images
make up                    # Start all services
make down                  # Stop all services
make logs                  # View logs
make restart               # Restart all services

# Database
make db-shell              # Connect to PostgreSQL shell
make db-migrate            # Run database migrations
make db-seed               # Seed database

# Testing
make test                  # Run backend tests
make test-frontend         # Run frontend tests
./scripts/test/test-api-summary.sh  # API integration tests

# Code Quality
make fmt                   # Format Go code
make lint                  # Run linters
make vet                   # Run go vet

# Build
make build-backend         # Build backend binary
make build-frontend        # Build frontend for production
```

### Git Workflow

```bash
# Check status
git status

# Create feature branch
git checkout -b feature/my-feature

# Commit changes
git add .
git commit -m "feat: add new feature"

# Push to remote
git push origin feature/my-feature

# Create pull request (use gh CLI or GitHub web)
gh pr create --title "Add new feature" --body "Description"
```

---

## Navigation Guide

### Finding Things Fast

#### By Feature/Module

| What You Need | Where to Look |
|---------------|---------------|
| **Member Management** | `backend/internal/models/anggota.go`<br>`backend/internal/services/anggota_service.go`<br>`backend/internal/handlers/anggota_handler.go`<br>`frontend/app/(dashboard)/members/` |
| **Authentication** | `backend/internal/services/auth_service.go`<br>`backend/internal/handlers/auth_handler.go`<br>`frontend/lib/context/AuthContext.tsx`<br>`frontend/app/(auth)/login/` |
| **Accounting** | `backend/internal/models/transaksi.go`<br>`backend/internal/services/transaksi_service.go`<br>`backend/internal/models/akun.go` |
| **Share Capital** | `backend/internal/models/simpanan.go`<br>`backend/internal/services/simpanan_service.go`<br>`frontend/app/(dashboard)/simpanan/` |
| **POS** | `backend/internal/models/penjualan.go`<br>`backend/internal/services/penjualan_service.go` |
| **Reports** | `backend/internal/services/laporan_service.go`<br>`backend/internal/handlers/laporan_handler.go` |

#### By Task

| Task | Where to Start |
|------|----------------|
| **Add new API endpoint** | 1. Create/update model in `internal/models/`<br>2. Add service in `internal/services/`<br>3. Add handler in `internal/handlers/`<br>4. Register route in `cmd/api/main.go` |
| **Add new frontend page** | 1. Create page in `app/`<br>2. Add API calls in `lib/api/`<br>3. Create components in `components/` |
| **Add database table** | 1. Create model in `internal/models/`<br>2. Run auto-migration in `internal/config/database.go` |
| **Modify business logic** | Look in `internal/services/` |
| **Change UI/UX** | Look in `frontend/app/` and `frontend/components/` |
| **Debug API issue** | 1. Check logs<br>2. Use `./scripts/debug/` scripts<br>3. Test with `./scripts/test/` scripts |
| **Setup environment** | Follow `docs/quick-start-guide.md` |
| **Understand architecture** | Read `docs/architecture.md` |

#### By File Type

| File Type | Locations |
|-----------|-----------|
| **Go Code** | `backend/cmd/`, `backend/internal/`, `backend/pkg/` |
| **TypeScript/React** | `frontend/app/`, `frontend/components/`, `frontend/lib/` |
| **Configuration** | Root level, `backend/.env`, `frontend/.env.local` |
| **Documentation** | `docs/`, `scripts/README.md`, root `*.md` files |
| **Tests** | `backend/*_test.go`, `frontend/__tests__/` |
| **Scripts** | `scripts/seed/`, `scripts/test/`, `scripts/debug/` |
| **Docker** | `Dockerfile`, `docker-compose.yml`, `nginx/` |

### Quick Reference: Important Files

| File | Purpose |
|------|---------|
| `CLAUDE.md` | Instructions for Claude Code |
| `README.md` | Project overview and setup |
| `docs/mvp-action-plan.md` | 12-week implementation plan |
| `docs/architecture.md` | System architecture |
| `docs/quick-start-guide.md` | Development setup guide |
| `backend/cmd/api/main.go` | Backend entry point |
| `frontend/app/layout.tsx` | Frontend root layout |
| `docker-compose.yml` | Service orchestration |
| `Makefile` | Development commands |
| `scripts/README.md` | Scripts documentation |

---

## Best Practices

### Code Organization

#### Backend (Go)

1. **Separation of Concerns**
   - Models: Data structures only
   - Services: Business logic
   - Handlers: HTTP request/response
   - Middleware: Cross-cutting concerns

2. **Package Structure**
   - `internal/`: Private code, not importable
   - `pkg/`: Public packages, can be imported
   - `cmd/`: Application entry points

3. **Naming Conventions**
   - Use Indonesian for domain concepts: `Anggota`, `Simpanan`, `Koperasi`
   - Use English for technical terms: `Service`, `Handler`, `Repository`
   - Files: lowercase with underscores: `anggota_service.go`

4. **Error Handling**
   - Return errors, don't panic
   - Use custom error types in `internal/errors/`
   - Log errors with context

#### Frontend (TypeScript/React)

1. **Component Organization**
   - One component per file
   - Co-locate related files
   - Use barrel exports (`index.ts`)

2. **Naming Conventions**
   - Components: PascalCase (`MemberList.tsx`)
   - Utilities: camelCase (`formatCurrency.ts`)
   - Types: PascalCase with Type suffix (`MemberType`)

3. **State Management**
   - Use Context for global state (auth)
   - Use local state for component-specific data
   - Server state with React Query (future)

4. **Type Safety**
   - Define types for all API responses
   - Use strict TypeScript config
   - Avoid `any` type

### Git Practices

1. **Branch Naming**
   - Feature: `feature/member-management`
   - Fix: `fix/login-validation`
   - Chore: `chore/update-dependencies`

2. **Commit Messages**
   - Use conventional commits format
   - Examples: `feat:`, `fix:`, `docs:`, `chore:`
   - Be descriptive but concise

3. **Pull Requests**
   - One feature/fix per PR
   - Include description and screenshots
   - Link related issues

### Documentation

1. **Code Documentation**
   - Document exported functions and types
   - Use godoc for Go
   - Use JSDoc for TypeScript

2. **README Files**
   - Each major directory should have a README
   - Include examples and usage

3. **Keep Docs Updated**
   - Update docs when changing features
   - Document breaking changes
   - Maintain changelog

### Testing

1. **Test Coverage**
   - Write tests for business logic
   - Test edge cases and error handling
   - Integration tests for critical flows

2. **Test Organization**
   - Tests alongside code: `*_test.go`, `*.test.ts`
   - Use descriptive test names
   - Arrange-Act-Assert pattern

3. **Test Data**
   - Use scripts in `scripts/seed/` for test data
   - Clean up after tests
   - Don't commit test databases

---

## Summary

This codebase is well-organized with:

✅ **Clear separation** between backend and frontend
✅ **Modular architecture** following best practices
✅ **Comprehensive documentation** (64,100+ lines)
✅ **Organized utility scripts** in categorized directories
✅ **Extensive testing** infrastructure
✅ **Production-ready** Docker setup
✅ **Multi-tenant architecture** from the ground up

### Quick Navigation

| Need | Go To |
|------|-------|
| Setup environment | `docs/quick-start-guide.md` |
| Understand architecture | `docs/architecture.md` |
| See implementation plan | `docs/mvp-action-plan.md` |
| Work with backend | `backend/` directory |
| Work with frontend | `frontend/` directory |
| Run scripts | `scripts/` directory |
| Read API docs | `backend/docs/swagger.yaml` |
| Configure Claude | `.claude/` directory |

---

**For detailed information on any specific component, refer to the respective section in this document or the linked documentation files.**

**Last Updated**: November 18, 2025
**Maintained By**: Development Team
**Questions?**: See `CONTRIBUTING.md` or create an issue
