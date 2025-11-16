# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Cooperative ERP Lite** is Indonesia's first cooperative operating system - a digital platform designed specifically for the Indonesian cooperative model. The project is in pre-development phase (Week 0), targeting an MVP launch in 12 weeks with 10 pilot cooperatives.

**Key Philosophy**: "Better than paper books = WIN" - The competition is manual record-keeping, so 20% improvement is a major win.

## Technology Stack

### Backend
- **Language**: Go 1.25.4 (latest stable, released November 5, 2025)
- **Framework**: Gin v1.10.0
- **ORM**: GORM v1.25.12 with PostgreSQL driver v1.5.9
- **Authentication**: JWT (golang-jwt/jwt/v5 v5.2.1) + bcrypt (golang.org/x/crypto v0.31.0)
- **Validation**: go-playground/validator/v10 v10.22.1

### Frontend
- **Framework**: Next.js 14 with TypeScript
- **UI Library**: Material-UI (MUI)
- **Forms**: React Hook Form + Zod validation
- **State Management**: React Context API
- **API Client**: Axios

### Infrastructure
- **Database**: PostgreSQL 15
- **Cloud**: Google Cloud Platform (Cloud Run, Cloud SQL, Cloud Storage)
- **Development**: Docker for local PostgreSQL

## Development Commands

### Backend (Go)

```bash
# Run server (from backend directory)
go run cmd/api/main.go

# Run with auto-reload (install air first: go install github.com/cosmtrek/air@latest)
air

# Run tests
go test ./...

# Format code
go fmt ./...

# Build binary
go build -o bin/api cmd/api/main.go

# Run seed data
go run cmd/seed/main.go

# Install dependencies
go mod download
go mod tidy
```

### Frontend (Next.js)

```bash
# Development server (from frontend directory)
npm run dev

# Build for production
npm run build

# Start production server
npm start

# Run linter
npm run lint

# Install dependencies
npm install
```

### Database

```bash
# Start PostgreSQL with Docker
docker run --name koperasi-postgres \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=koperasi_erp \
  -p 5432:5432 \
  -d postgres:15

# Start existing container
docker start koperasi-postgres

# Stop container
docker stop koperasi-postgres

# Connect to database
docker exec -it koperasi-postgres psql -U postgres -d koperasi_erp

# Backup database
docker exec koperasi-postgres pg_dump -U postgres koperasi_erp > backup.sql

# Restore database
docker exec -i koperasi-postgres psql -U postgres koperasi_erp < backup.sql
```

## Project Structure

### Current State (Pre-Development)
```
COOPERATIVE-ERP-LITE/
├── src/                    # Source code (currently empty placeholders)
│   ├── config/            # Configuration files
│   ├── modules/           # Feature modules
│   └── shared/            # Shared utilities
├── docs/                  # Comprehensive documentation
│   ├── business/          # Business documentation
│   ├── phases/            # Phase-specific implementation guides
│   │   ├── phase-1/      # MVP (Month 1-3, 12 weeks)
│   │   └── phase-2/      # Enhanced features (Month 4-6)
│   ├── architecture.md    # System architecture
│   ├── mvp-action-plan.md # 12-week sprint plan
│   ├── quick-start-guide.md # Development setup
│   └── technical-stack-versions.md # Locked dependency versions
├── tests/                 # Test files
├── .claude/              # Claude Code configuration
│   ├── commands/         # Custom slash commands
│   ├── agents/           # Specialized agent configurations
│   ├── skills/           # Reusable skills
│   └── docs/             # Claude-specific documentation
└── README.md             # Project overview

Backend structure (to be created):
backend/
├── cmd/
│   ├── api/main.go       # API server entry point
│   └── seed/main.go      # Database seeding
├── internal/
│   ├── models/           # Data models (User, Member, Transaction, etc.)
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # Authentication, CORS, logging
│   ├── services/         # Business logic
│   ├── database/         # Database connection and migrations
│   └── config/           # Configuration loader
├── pkg/
│   ├── utils/            # Utility functions
│   └── logger/           # Logging utilities
├── go.mod
├── go.sum
└── .env                  # Environment variables (not in git)

Frontend structure (to be created):
frontend/
├── app/                  # Next.js App Router
│   ├── login/page.tsx
│   ├── dashboard/page.tsx
│   └── layout.tsx
├── components/           # React components
├── lib/
│   └── api.ts           # Axios API client
├── package.json
└── .env.local           # Frontend environment variables (not in git)
```

## Architecture Principles

### Modular Monolith (Phase 1)
- Start with well-separated modules in a single codebase
- Clear module boundaries and interfaces
- Shared database with schema separation
- Can transition to microservices in Phase 2+ if needed

### Design Patterns
- **Domain-Driven Design (DDD)**: Bounded contexts per business domain
- **CQRS**: Separate read and write operations where beneficial
- **Repository Pattern**: Abstract data access layer for testability

### Core Modules (8 MVP Features)
1. **Authentication & Authorization**: JWT-based auth with 4 roles (Admin, Treasurer, Cashier, Member)
2. **Member Management**: CRUD operations for cooperative members
3. **Share Capital Tracking**: Simpanan Pokok, Wajib, Sukarela
4. **Basic POS**: Cash-only point of sale (no QRIS, no printing)
5. **Simple Accounting**: Manual journal entries with pre-configured Chart of Accounts
6. **Essential Reports**: Member balances, daily transactions, balance sheet, P&L
7. **Member Portal**: Web-based self-service (view balances, transactions, announcements)
8. **Data Import**: Excel template-based initial data migration

## Important Constraints and Decisions

### MVP Scope is LOCKED
The following features are explicitly OUT OF SCOPE for MVP (deferred to Phase 2):
- ❌ SHU (profit sharing) calculation
- ❌ QRIS/digital payments (cash only)
- ❌ Barcode scanning (manual entry)
- ❌ Receipt printing (digital receipts only)
- ❌ Inventory automation (manual stock count)
- ❌ WhatsApp integration
- ❌ Native mobile app (responsive web only)
- ❌ Offline mode (internet required)
- ❌ Approval workflows
- ❌ NIK validation API

### Indonesian Context
- **Accounting Standards**: SAK ETAP (Indonesian GAAP for cooperatives)
- **Legal Framework**: UU No. 25 Tahun 1992 (Indonesian Cooperative Law)
- **Chart of Accounts**: Pre-configured for Indonesian cooperative accounting
- **Language**: UI should support Indonesian (Bahasa Indonesia)
- **Share Capital Types**:
  - Simpanan Pokok (Principal deposit - paid once)
  - Simpanan Wajib (Mandatory deposit - monthly)
  - Simpanan Sukarela (Voluntary deposit - optional)

### Target Market
- **Primary**: 80,000 Koperasi Merah Putih (launching July 2025, IDR 3B funding each)
- **Secondary**: 127,000+ existing cooperatives
- **User Base**: Non-technical cooperative treasurers and staff
- **Deployment**: Multi-tenant SaaS platform

## Database Schema (Core Tables)

```sql
-- Key tables to implement in Phase 1
cooperatives       -- Cooperative organizations (multi-tenant)
users             -- System users with roles
members           -- Cooperative members
share_capital     -- Member share capital transactions
accounts          -- Chart of Accounts
transactions      -- Accounting journal entries
transaction_lines -- Journal entry line items
products          -- Product catalog for POS
sales             -- POS sales transactions
sale_items        -- Line items for sales
```

## Authentication & Security

- **JWT Tokens**: Stored in HTTP-only cookies (frontend) with 24-hour expiration
- **Password Hashing**: bcrypt with default cost
- **CORS**: Configured for frontend domain
- **Roles**: Admin (full access), Treasurer (finance), Cashier (POS only), Member (read-only portal)
- **Multi-tenancy**: All queries filtered by cooperative_id from JWT token

## Critical Documentation References

When working on this project, always reference:
1. **docs/mvp-action-plan.md** - The single source of truth for what to build
2. **docs/quick-start-guide.md** - Development environment setup
3. **docs/technical-stack-versions.md** - Exact dependency versions (LOCKED)
4. **docs/phases/phase-1/implementation-guide.md** - Week-by-week implementation plan
5. **docs/architecture.md** - System architecture and design patterns

## Testing Strategy

### Unit Tests
- All service layer functions
- Utility functions
- Data validation logic

### Integration Tests
- API endpoints with database
- Authentication flows
- Multi-tenant data isolation

### E2E Tests (Week 8+)
- Critical user journeys
- POS transaction flow
- Report generation

## Common Development Patterns

### Backend Handler Pattern
```go
// handlers/member_handler.go
type MemberHandler struct {
    service *services.MemberService
}

func (h *MemberHandler) Create(c *gin.Context) {
    var req CreateMemberRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    // Get cooperative_id from JWT token
    coopID := c.GetString("cooperative_id")

    member, err := h.service.Create(coopID, &req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(201, member)
}
```

### Frontend API Call Pattern
```typescript
// Use the configured API client from lib/api.ts
import api from '@/lib/api';

// API calls automatically include JWT token from cookies
const response = await api.post('/members', {
  name: 'John Doe',
  // ... other fields
});
```

### Multi-tenant Query Pattern
```go
// All queries must filter by cooperative_id
var members []models.Member
db.Where("cooperative_id = ?", cooperativeID).Find(&members)
```

## Success Metrics

### Week 4 Goals
- Register 100 members in < 30 minutes
- Record 50 transactions correctly
- Trial balance validates

### Week 8 Goals
- All 8 features working
- 4 reports generate correctly
- Mobile-responsive UI

### Week 12 Goals (MVP Launch)
- 8+ cooperatives using daily
- 1,000+ transactions recorded
- Zero data loss
- 3 strong testimonials
- Customer satisfaction > 7/10

## Development Workflow

1. **Environment Setup**: Follow docs/quick-start-guide.md exactly
2. **Version Lock**: Use exact versions from docs/technical-stack-versions.md
3. **Feature Development**: Check docs/mvp-action-plan.md for scope
4. **Code Structure**: Follow the backend/frontend structure defined above
5. **Testing**: Write tests alongside feature code
6. **Deployment**: Cloud Run (backend), Vercel/Cloud Run (frontend)

## Important Notes

- **Ship Fast, Iterate**: Perfect is the enemy of good. Get it working, then improve.
- **Multi-tenant First**: Every query must filter by cooperative_id
- **Indonesian Standards**: Follow SAK ETAP accounting, Indonesian UX patterns
- **Mobile-First**: Design for mobile even though it's responsive web
- **Data Integrity**: Accounting data must be accurate and auditable
- **Offline Not Supported**: Require internet connection for MVP (simplifies architecture)

## When in Doubt

1. Check if the feature is in MVP scope (docs/mvp-action-plan.md)
2. Reference the implementation guide (docs/phases/phase-1/implementation-guide.md)
3. Use locked versions (docs/technical-stack-versions.md)
4. Follow Indonesian cooperative accounting standards
5. Ask: "Is this better than paper books?" If yes, ship it.
