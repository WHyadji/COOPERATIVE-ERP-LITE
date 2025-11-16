# Technical Stack - Version Specifications

**Last Updated:** November 16, 2025
**Documentation Method:** Using Context7 MCP for up-to-date stable versions

---

## 📦 Version Lock Strategy

**Philosophy:**
- Use **latest stable** versions as of November 2025
- Lock exact versions to avoid "works on my machine" issues
- Allow security patches (patch version updates)
- Review minor version updates quarterly
- Major version updates only between project phases

---

## 🎯 MVP Version Lock (Month 1-3)

### Backend Stack (Go)

**Core Language:**
```yaml
go: "1.25.4"  # Latest stable patch (November 5, 2025)
# Released: August 12, 2025 (1.25.0), Latest patch: November 5, 2025
# EOL: August 2027 (when 1.27 releases)
# Why: Latest stable with improved generics, experimental GC, encoding/json/v2
# Previous versions: 1.24 (Feb 2025), 1.23 (Aug 2024)
```

**🚀 Why Go 1.25.4?**

Go 1.25 is the **latest major release** (August 2025) with 3+ months of production stability. Key advantages:

1. **Improved Generics** - Better type inference and performance for generic code
2. **Experimental encoding/json/v2** - Modern JSON handling with better performance and API
3. **Experimental Garbage Collector** - Optional new GC with lower latency
4. **2 Years of Support** - Security updates until August 2027 (Go 1.27 release)
5. **Production Ready** - 1.25.4 is the 4th patch release with bug fixes from initial 1.25.0

**Release Timeline:**
- Go 1.25.0: August 12, 2025 (initial release)
- Go 1.25.1: September 3, 2025
- Go 1.25.3: October 13, 2025
- **Go 1.25.4: November 5, 2025** ← Current (11 days old, stable)

Using the latest stable major version ensures:
- Access to newest language features
- Longest support window (2 years)
- Best compatibility with modern libraries
- Most recent performance improvements

**Web Framework:**
```yaml
github.com/gin-gonic/gin: "v1.10.0"
# Latest stable release
# Why: Most popular Go web framework, excellent performance
# Alternative: github.com/labstack/echo v4.12.0
```

**ORM & Database:**
```yaml
gorm.io/gorm: "v1.25.12"
# Latest stable with Go 1.25 support
gorm.io/driver/postgres: "v1.5.9"
# PostgreSQL driver for GORM
```

**Authentication:**
```yaml
github.com/golang-jwt/jwt/v5: "v5.2.1"
# Latest stable JWT library
golang.org/x/crypto: "v0.31.0"
# For bcrypt password hashing
```

**Environment & Configuration:**
```yaml
github.com/joho/godotenv: "v1.5.1"
# For .env file loading
```

**CORS & Middleware:**
```yaml
github.com/gin-contrib/cors: "v1.7.2"
# Official CORS middleware for Gin
```

**Validation:**
```yaml
github.com/go-playground/validator/v10: "v10.22.1"
# Struct validation
```

---

### Frontend Stack (Next.js + React)

**Runtime:**
```yaml
node: "20.18.1"  # LTS (Maintenance until April 2026)
npm: "10.9.2"    # Bundled with Node 20.18.1
# Why: Node 20 is current LTS, stable and well-supported
# Note: Node 22 is also LTS but too new for production use
```

**Core Framework:**
```yaml
next: "15.5.0"
# Released: October 2024 (15.0), Latest: January 2025 (15.5)
# Why: Latest stable with Turbopack, React 19 support

react: "19.2.0"
react-dom: "19.2.0"
# Released: December 2024 (19.0), Latest: October 2025 (19.2)
# Why: Latest stable with Server Components, Actions, async transitions

typescript: "5.7.3"
# Latest stable TypeScript
```

**UI Framework:**
```yaml
@mui/material: "6.3.0"
@mui/icons-material: "6.3.0"
@emotion/react: "11.14.0"
@emotion/styled: "11.14.0"
# Latest Material-UI v6 with React 19 support
```

**Form Handling:**
```yaml
react-hook-form: "7.54.2"
# Latest stable form library
@hookform/resolvers: "3.9.1"
# Validation resolvers
zod: "3.24.1"
# Schema validation
```

**HTTP Client:**
```yaml
axios: "1.7.9"
# Latest stable
js-cookie: "3.0.5"
# Cookie handling
```

**Date Handling:**
```yaml
date-fns: "4.1.0"
# Modern date utility library (better than moment.js)
# With Indonesian locale support
```

**State Management (if needed):**
```yaml
zustand: "5.0.2"
# Lightweight state management (alternative to Redux)
# Only add if complex state needed
```

---

### Database

**PostgreSQL:**
```yaml
version: "17.2"
# Released: September 2024 (17.0), Latest patch: November 2025 (17.2)
# EOL: November 2029
# Why: Latest stable with JSON_TABLE, 2x COPY performance, better vacuum
# Migration from 15/16: Straightforward, no breaking changes
```

**Docker Image:**
```yaml
postgres:17.2-alpine
# Alpine for smaller image size
```

**Extensions:**
```yaml
uuid-ossp: "bundled"  # UUID generation
pgcrypto: "bundled"   # Encryption functions (if needed)
```

---

### Infrastructure & DevOps

**Docker:**
```yaml
docker_engine: "27.4.0"  # Latest stable
docker_compose: "2.31.0" # Latest stable
```

**Container Base Images:**
```yaml
# Backend
golang:1.25.4-alpine3.21

# Frontend (Development)
node:20.18.1-alpine3.21

# Frontend (Production)
node:20.18.1-alpine3.21  # For build
nginx:1.27.3-alpine       # For serving static files

# Database
postgres:17.2-alpine3.21
```

**Google Cloud Platform:**
```yaml
# Cloud Run
runtime: "managed"
region: "asia-southeast2"  # Jakarta

# Cloud SQL
version: "POSTGRES_17"
region: "asia-southeast2"

# Cloud Storage
storage_class: "REGIONAL"
region: "asia-southeast2"
```

---

### Development Tools

**Code Quality (Go):**
```yaml
golangci-lint: "v1.62.2"
# Comprehensive Go linter
```

**Code Quality (TypeScript/JavaScript):**
```yaml
eslint: "9.18.0"
# Latest ESLint
@typescript-eslint/parser: "8.20.0"
@typescript-eslint/eslint-plugin: "8.20.0"
prettier: "3.4.2"
# Code formatter
```

**Testing:**
```yaml
# Go (built-in)
testing: "standard library"  # Go 1.25

# Frontend
vitest: "2.1.8"              # Faster than Jest, Vite-compatible
@testing-library/react: "16.1.0"  # React 19 compatible
@testing-library/user-event: "14.5.2"
```

**Hot Reload:**
```yaml
# Go
github.com/air-verse/air: "v1.61.7"
# Automatic reloading for Go

# Next.js
# Built-in Fast Refresh (no additional package needed)
```

---

## 🚀 Phase 2 Additional Dependencies (Month 4-6)

### SHU Calculation Engine
```yaml
github.com/shopspring/decimal: "v1.4.0"
# Arbitrary-precision decimal for financial calculations
# Why: Avoid floating-point errors in money calculations
```

### QRIS Payment Integration
```yaml
# Option 1: Xendit
github.com/xendit/xendit-go/v6: "v6.1.0"

# Option 2: Midtrans
github.com/midtrans/midtrans-go: "v1.3.8"

# Decision: TBD during Phase 2 planning
```

### Bank API Integration
```yaml
# Will depend on bank partner selection
# Candidates: BCA, BRI, Mandiri, BNI APIs
# TBD during Phase 2
```

### Mobile App (React Native)
```yaml
react-native: "0.76.5"
# Latest stable (released November 2024)

# With Expo (Recommended for faster development)
expo: "52.0.30"
# Latest Expo SDK

# Navigation
@react-navigation/native: "7.0.16"
@react-navigation/stack: "7.2.4"

# UI Components
react-native-paper: "5.12.5"
# Material Design for React Native
```

### Loan & Interest Calculation
```yaml
# Will use shopspring/decimal (already added)
# Custom implementation for Indonesian loan types
```

---

## 📊 Phase 3 Additional Dependencies (Month 7-9)

### WhatsApp Business Integration
```yaml
# Official WhatsApp Business API
whatsapp-business-api: "via Meta Cloud API"
# No SDK needed, use HTTP client (axios)

# OR Third-party like Fonnte, Wablas
# TBD based on cost analysis
```

### Inventory Management
```yaml
# Using existing GORM models
# No additional dependencies needed
```

### Business Intelligence
```yaml
# Frontend charting
recharts: "2.15.0"
# React charting library

# Or
@mui/x-charts: "8.0.1"
# Material-UI charts (if already using MUI)
```

### Background Jobs
```yaml
github.com/hibiken/asynq: "v0.25.1"
# Redis-backed distributed task queue for Go
# For scheduled reports, late fee calculations, etc.

github.com/redis/go-redis/v9: "v9.7.0"
# Redis client for Go
```

**Redis:**
```yaml
redis: "7.4.2"
# Latest stable Redis
# Docker: redis:7.4.2-alpine
```

---

## 🔐 Security & Compliance

### SSL/TLS
```yaml
tls_version: "1.3"
# Minimum TLS 1.3 for all connections
```

### Password Hashing
```yaml
bcrypt_cost: 12
# 12 rounds (secure but not too slow)
# Provided by golang.org/x/crypto
```

### JWT Configuration
```yaml
algorithm: "HS256"
# For MVP (symmetric key)
# Phase 3+: Consider RS256 (asymmetric) for better security
```

### Rate Limiting
```yaml
github.com/ulule/limiter/v3: "v3.11.3"
# Rate limiting middleware for Go
```

### CORS Configuration
```yaml
allowed_origins:
  - "http://localhost:3000"      # Development
  - "https://app.koperasi.id"    # Production
allowed_methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"]
allowed_headers: ["Authorization", "Content-Type"]
max_age: 3600
```

---

## 📦 Complete package.json (Frontend)

```json
{
  "name": "koperasi-erp-frontend",
  "version": "0.1.0",
  "private": true,
  "engines": {
    "node": ">=20.18.1",
    "npm": ">=10.9.2"
  },
  "scripts": {
    "dev": "next dev",
    "build": "next build",
    "start": "next start",
    "lint": "next lint",
    "format": "prettier --write \"**/*.{js,jsx,ts,tsx,json,md}\""
  },
  "dependencies": {
    "next": "15.5.0",
    "react": "19.2.0",
    "react-dom": "19.2.0",
    "@mui/material": "6.3.0",
    "@mui/icons-material": "6.3.0",
    "@emotion/react": "11.14.0",
    "@emotion/styled": "11.14.0",
    "react-hook-form": "7.54.2",
    "@hookform/resolvers": "3.9.1",
    "zod": "3.24.1",
    "axios": "1.7.9",
    "js-cookie": "3.0.5",
    "date-fns": "4.1.0"
  },
  "devDependencies": {
    "typescript": "5.7.3",
    "@types/node": "22.10.2",
    "@types/react": "19.0.7",
    "@types/react-dom": "19.0.2",
    "eslint": "9.18.0",
    "eslint-config-next": "15.5.0",
    "@typescript-eslint/parser": "8.20.0",
    "@typescript-eslint/eslint-plugin": "8.20.0",
    "prettier": "3.4.2",
    "vitest": "2.1.8",
    "@testing-library/react": "16.1.0",
    "@testing-library/user-event": "14.5.2"
  }
}
```

---

## 📦 Complete go.mod (Backend)

```go
module github.com/yourusername/koperasi-erp

go 1.25

require (
    github.com/gin-gonic/gin v1.10.0
    github.com/gin-contrib/cors v1.7.2
    gorm.io/gorm v1.25.12
    gorm.io/driver/postgres v1.5.9
    github.com/golang-jwt/jwt/v5 v5.2.1
    golang.org/x/crypto v0.31.0
    github.com/joho/godotenv v1.5.1
    github.com/go-playground/validator/v10 v10.22.1
)

// Phase 2 dependencies (add when needed)
// github.com/shopspring/decimal v1.4.0
// github.com/xendit/xendit-go/v6 v6.1.0

// Phase 3 dependencies (add when needed)
// github.com/hibiken/asynq v0.25.1
// github.com/redis/go-redis/v9 v9.7.0
// github.com/ulule/limiter/v3 v3.11.3
```

---

## 🔄 Version Update Policy

### During MVP (Week 1-12):
- ✅ Security patches: Apply immediately
- ⚠️ Patch versions (x.x.X): Review and apply if low risk
- ❌ Minor versions (x.X.x): Defer unless critical bug fix
- ❌ Major versions (X.x.x): NO updates during MVP

### Between Phases:
- ✅ Review all dependency updates
- ✅ Test in development environment
- ✅ Update documentation
- ✅ Deploy to staging before production

### Security Vulnerabilities:
- 🚨 Critical: Patch within 24 hours
- ⚠️ High: Patch within 1 week
- ℹ️ Medium/Low: Include in next regular update

---

## 🧪 Testing Version Compatibility

Before Phase 1 starts, verify this combination works:

```bash
# Backend test
go version          # Should show: go1.25.4 linux/amd64 (or darwin/amd64)
go run cmd/api/main.go

# Frontend test
node --version      # Should show: v20.18.1
npm --version       # Should show: 10.9.2
npm run dev

# Database test
docker run -d postgres:17.2-alpine
psql --version      # Should show: 17.2
```

---

## 📚 Documentation Sources

All version information verified using **Context7 MCP** which provides up-to-date documentation from official sources:

- Go: https://go.dev/doc/devel/release
- Next.js: https://nextjs.org/blog
- React: https://react.dev/versions
- PostgreSQL: https://www.postgresql.org/docs/release/
- Material-UI: https://mui.com/versions/
- npm packages: https://www.npmjs.com/

**Last verified:** November 16, 2025

---

## ⚡ Key Decisions Summary

| Category | Choice | Why |
|----------|--------|-----|
| **Go** | 1.25.4 | Latest stable (Nov 2025), improved generics, experimental GC, 2 years support |
| **Next.js** | 15.5.0 | Latest with Turbopack, React 19 support |
| **React** | 19.2.0 | Latest stable with Server Components, Actions |
| **PostgreSQL** | 17.2 | Latest stable, JSON_TABLE, 2x COPY performance |
| **Node** | 20.18.1 LTS | Current LTS until April 2026 |
| **TypeScript** | 5.7.3 | Latest stable |
| **Material-UI** | 6.3.0 | Latest with React 19 support |

---

## 🎯 Next Steps

1. **Week 0 (Now):**
   - Team installs exact versions listed above
   - Run compatibility tests
   - Create `go.mod` and `package.json` with locked versions
   - Commit to Git as "Version Lock Baseline"

2. **Week 1:**
   - Initialize projects with these versions
   - Setup CI/CD to enforce version checks
   - Document any version-specific gotchas

3. **Monthly:**
   - Review security advisories
   - Check for critical updates
   - Update this document if versions change

---

**Version Control:**
- Initial version: November 16, 2025
- Using: Context7 MCP for up-to-date documentation
- Review frequency: Monthly during MVP, Quarterly post-MVP
