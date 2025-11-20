# Development Documentation

This directory contains technical documentation for developers working on the Cooperative ERP Lite system.

## Available Guides

### 📘 [Module Development Guide](./MODULE-DEVELOPMENT-GUIDE.md)
**Complete step-by-step guide for building custom modules**

Covers:
- Planning & Design
- Database Schema & Migrations
- Backend Development (Go + Gin + GORM)
- Frontend Development (Next.js + TypeScript + MUI)
- API Integration & Testing
- Deployment

**Use this guide when:** You need to build a new module from scratch

---

## Quick Start for Module Development

```bash
# 1. Read the planning section
docs/development/MODULE-DEVELOPMENT-GUIDE.md#step-1-planning--design

# 2. Design database schema
docs/development/MODULE-DEVELOPMENT-GUIDE.md#step-2-database-design

# 3. Follow step-by-step implementation
docs/development/MODULE-DEVELOPMENT-GUIDE.md#step-3-backend-development
```

---

## Development Workflow

```
Planning (1 day)
  ↓
Database Design (0.5 day)
  ↓
Backend Development (2-3 days)
  ↓
Frontend Development (2-3 days)
  ↓
Integration & Testing (1 day)
  ↓
Documentation (0.5 day)
  ↓
Deployment (0.5 day)

Total: 7-10 days per module
```

---

## Code Examples

The Module Development Guide includes complete code examples for:

✅ **Database Migrations** - SQL with proper indexing and constraints
✅ **Go Models** - GORM models with validation
✅ **Service Layer** - Business logic with multi-tenant filtering
✅ **HTTP Handlers** - Gin handlers with Swagger docs
✅ **TypeScript Types** - Complete type definitions
✅ **React Components** - MUI components with React Hook Form
✅ **API Integration** - Axios client setup
✅ **Unit Tests** - Go test examples

---

## Best Practices Checklist

When developing modules, always:

- ✅ Filter by `cooperative_id` (multi-tenancy)
- ✅ Use database transactions for multi-step operations
- ✅ Validate input on both frontend and backend
- ✅ Handle errors properly (specific error messages)
- ✅ Write unit tests for service layer
- ✅ Document API endpoints with Swagger
- ✅ Test multi-tenant isolation
- ✅ Follow naming conventions (camelCase JS, snake_case SQL)

---

## Common Module Types

### Essential Tier (IDR 99k/month)
- Simple, single-purpose modules
- 3-5 core features
- 5 days development time

**Examples:** Fuel Tracking, Customer Loyalty, Crop Planning

### Professional Tier (IDR 199k/month)
- Comprehensive modules
- 8-12 features
- 7 days development time

**Examples:** Vehicle Rental, Inventory Management, Harvest Tracking

### Enterprise Tier (IDR 399k/month)
- Complex, mission-critical modules
- 15+ features
- 10-15 days development time

**Examples:** Fleet Management (GPS), Micro-Lending, Patient Records

---

## Tech Stack Reference

### Backend
- **Language:** Go 1.24+
- **Framework:** Gin v1.10.0
- **ORM:** GORM v1.25.12
- **Database:** PostgreSQL 15
- **Auth:** JWT (golang-jwt/jwt/v5)
- **Validation:** go-playground/validator/v10

### Frontend
- **Framework:** Next.js 14 (App Router)
- **Language:** TypeScript 5
- **UI Library:** Material-UI (MUI) v5
- **Forms:** React Hook Form + Zod
- **API Client:** Axios
- **State:** React Context API

---

## Directory Structure

```
docs/development/
├── README.md                           # This file
├── MODULE-DEVELOPMENT-GUIDE.md         # Complete module dev guide
├── API-DOCUMENTATION.md                # API design standards (future)
├── TESTING-GUIDE.md                    # Testing strategies (future)
└── DEPLOYMENT-GUIDE.md                 # Deployment procedures (future)
```

---

## Related Documentation

**Business Documentation:**
- [Module Catalog](../business/module/MODULE-CATALOG-COMPLETE.md) - All 21 modules
- [Pricing Strategy](../business/module/OPTIMAL-PRICING-STRATEGY.md) - Pure SaaS pricing

**Architecture Documentation:**
- [System Architecture](../architecture.md) - Overall system design
- [Technical Stack](../technical-stack-versions.md) - Locked versions

**Project Planning:**
- [MVP Action Plan](../mvp-action-plan.md) - 12-week roadmap
- [Quick Start Guide](../quick-start-guide.md) - Environment setup

---

## Getting Help

**For technical questions:**
- Check the [Module Development Guide](./MODULE-DEVELOPMENT-GUIDE.md)
- Review code examples in the guide
- Check the [Troubleshooting section](./MODULE-DEVELOPMENT-GUIDE.md#troubleshooting)

**For architecture questions:**
- See [Module Architecture](./MODULE-DEVELOPMENT-GUIDE.md#module-architecture)
- Review modular monolith pattern
- Check multi-tenant design patterns

**For deployment questions:**
- See [Step 7: Deployment](./MODULE-DEVELOPMENT-GUIDE.md#step-7-deployment)
- Check Cloud Run deployment instructions

---

## Module Development Checklist

Use this checklist when building a new module:

### Planning Phase
- [ ] Requirements analysis complete
- [ ] Database schema designed
- [ ] API endpoints defined
- [ ] Frontend pages planned
- [ ] User roles identified

### Backend Development
- [ ] Migration file created
- [ ] Models defined with validation
- [ ] Service layer implemented
- [ ] HTTP handlers created
- [ ] Routes registered
- [ ] Unit tests written
- [ ] Multi-tenant filtering verified

### Frontend Development
- [ ] TypeScript types defined
- [ ] API client functions created
- [ ] React components built
- [ ] Forms with validation
- [ ] Pages created
- [ ] Responsive design tested

### Testing & Quality
- [ ] Unit tests pass
- [ ] Integration tests pass
- [ ] Manual testing complete
- [ ] Multi-tenant isolation verified
- [ ] Error handling tested
- [ ] Edge cases covered

### Documentation
- [ ] API docs (Swagger) complete
- [ ] User guide written
- [ ] Code comments added
- [ ] README updated

### Deployment
- [ ] Migrations run on staging
- [ ] Backend deployed to staging
- [ ] Frontend deployed to staging
- [ ] Staging testing complete
- [ ] Production deployment

---

**Last Updated:** 2025-01-19
**Maintained By:** Development Team
