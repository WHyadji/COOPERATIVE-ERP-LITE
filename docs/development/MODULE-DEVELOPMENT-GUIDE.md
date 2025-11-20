# Module Development Guide

**Document Version:** 1.0
**Purpose:** Comprehensive guide for developing custom modules
**Audience:** Backend & Frontend Developers
**Last Updated:** 2025-01-19

---

## Table of Contents

1. [Overview](#overview)
2. [Prerequisites](#prerequisites)
3. [Module Architecture](#module-architecture)
4. [Development Workflow](#development-workflow)
5. [Step 1: Planning & Design](#step-1-planning--design)
6. [Step 2: Database Design](#step-2-database-design)
7. [Step 3: Backend Development](#step-3-backend-development)
8. [Step 4: Frontend Development](#step-4-frontend-development)
9. [Step 5: Integration & Testing](#step-5-integration--testing)
10. [Step 6: Documentation](#step-6-documentation)
11. [Step 7: Deployment](#step-7-deployment)
12. [Best Practices](#best-practices)
13. [Common Pitfalls](#common-pitfalls)
14. [Troubleshooting](#troubleshooting)

---

## Overview

### What is a Module?

A **module** is a self-contained feature that extends the core ERP system with industry-specific functionality. Each module:

- ✅ Has its own database tables (with schema isolation)
- ✅ Provides REST API endpoints
- ✅ Has dedicated frontend pages/components
- ✅ Is independently activatable per cooperative
- ✅ Follows multi-tenant architecture (cooperative_id filtering)
- ✅ Integrates with core ERP (auth, accounting, members)

### Module Examples

**Simple Module (Essential Tier):**
- Fuel Tracking & Analytics
- Customer Loyalty Program
- Crop Planning & Calendar

**Medium Module (Professional Tier):**
- Vehicle Rental Management
- Advanced Inventory Management
- Harvest & Distribution Tracking

**Complex Module (Enterprise Tier):**
- Fleet Management System (GPS integration)
- Micro-Lending System (credit scoring, installments)
- Patient Records Management (BPJS integration)

---

## Prerequisites

### Required Knowledge

**Backend:**
- ✅ Go programming (1.24+)
- ✅ Gin web framework
- ✅ GORM ORM
- ✅ PostgreSQL
- ✅ RESTful API design
- ✅ JWT authentication

**Frontend:**
- ✅ TypeScript
- ✅ Next.js 14 (App Router)
- ✅ React 18
- ✅ Material-UI (MUI)
- ✅ React Hook Form
- ✅ Axios for API calls

**General:**
- ✅ Git version control
- ✅ Docker basics
- ✅ Multi-tenant architecture concepts

### Development Environment

**Required Software:**
```bash
# Go 1.24+
go version

# Node.js 20+ and npm
node -v
npm -v

# PostgreSQL 15+
psql --version

# Docker (for local database)
docker -v

# Git
git --version
```

**Project Structure:**
```
COOPERATIVE-ERP-LITE/
├── backend/
│   ├── cmd/
│   │   └── api/main.go
│   ├── internal/
│   │   ├── models/          # Database models
│   │   ├── handlers/        # HTTP handlers
│   │   ├── services/        # Business logic
│   │   ├── middleware/      # Auth, CORS, etc.
│   │   └── database/        # DB connection
│   ├── migrations/          # Database migrations
│   └── go.mod
├── frontend/
│   ├── app/                 # Next.js App Router
│   ├── components/          # React components
│   ├── lib/
│   │   └── api.ts          # Axios client
│   ├── types/               # TypeScript types
│   └── package.json
└── docs/
    └── development/         # This guide
```

---

## Module Architecture

### Modular Monolith Pattern

We use a **modular monolith** architecture:
- Single codebase (not microservices)
- Clear module boundaries
- Shared database with schema isolation
- Direct service injection (no message queue for MVP)

```
┌─────────────────────────────────────────────────┐
│              CORE ERP SYSTEM                    │
│  (Auth, Members, Accounting, Reports)           │
└──────────────────┬──────────────────────────────┘
                   │
                   │ Shared Services
                   │
    ┌──────────────┼──────────────┬───────────────┐
    │              │              │               │
┌───▼────┐   ┌────▼─────┐  ┌─────▼──────┐  ┌────▼─────┐
│Module 1│   │Module 2  │  │Module 3    │  │Module N  │
│Vehicle │   │Inventory │  │Lending     │  │Analytics │
│Rental  │   │Mgmt      │  │System      │  │& BI      │
└────────┘   └──────────┘  └────────────┘  └──────────┘
```

### Multi-Tenant Architecture

**CRITICAL:** All queries MUST filter by `cooperative_id`.

```go
// ❌ WRONG - No cooperative filter
db.Find(&vehicles)

// ✅ CORRECT - With cooperative filter
db.Where("cooperative_id = ?", cooperativeID).Find(&vehicles)
```

**Why?**
- Security: Prevent data leakage between cooperatives
- Isolation: Each cooperative sees only their data
- Compliance: Data privacy requirements

---

## Development Workflow

### Standard Development Process

```
1. Planning & Design (1 day)
   └─ Define requirements, database schema, API endpoints

2. Database Design (0.5 day)
   └─ Create migration files, define models

3. Backend Development (2-3 days)
   └─ Implement models, services, handlers, tests

4. Frontend Development (2-3 days)
   └─ Create pages, components, forms, API integration

5. Integration & Testing (1 day)
   └─ End-to-end testing, bug fixes

6. Documentation (0.5 day)
   └─ API docs, user guide

7. Deployment (0.5 day)
   └─ Deploy to staging, then production

Total: 7-10 days per module
```

### Git Workflow

```bash
# 1. Create feature branch
git checkout -b feature/module-vehicle-rental

# 2. Make changes and commit regularly
git add .
git commit -m "feat(rental): add vehicle model and migration"
git commit -m "feat(rental): add booking API endpoints"
git commit -m "feat(rental): add rental dashboard UI"

# 3. Push to remote
git push origin feature/module-vehicle-rental

# 4. Create Pull Request (PR)
gh pr create --title "Add Vehicle Rental Module" --body "..."

# 5. After review, merge to main
git checkout main
git pull origin main
```

---

## Step 1: Planning & Design

### 1.1 Requirements Analysis

Create a specification document:

**Example: Vehicle Rental Module**

```markdown
## Module: Vehicle Rental Management

### Overview
Enable transport cooperatives to manage vehicle rentals, bookings, and member payouts.

### Core Features
1. Vehicle Management (CRUD, photos, status)
2. Booking System (calendar, pricing, contracts)
3. Return Processing (inspection, late fees)
4. Maintenance Scheduling
5. Revenue Reports

### User Roles
- Admin: Full access
- Operations Manager: Create/manage bookings
- Vehicle Owner (Member): View their vehicle bookings & revenue
- Customer: N/A (external, not in system)

### Success Criteria
- Create 50 bookings in < 30 minutes
- Zero double-bookings
- Accurate revenue calculation (100% match manual calculation)
```

### 1.2 Database Schema Design

**Design Tables:**

```sql
-- Example: Vehicle Rental Module

1. vehicles
   - id, cooperative_id, owner_member_id
   - make, model, year, plate_number
   - status, base_rate_daily
   - created_at, updated_at

2. rental_contracts
   - id, cooperative_id, vehicle_id
   - customer_name, customer_phone
   - pickup_date, return_date
   - total_amount, status
   - created_at, updated_at

3. maintenance_records
   - id, vehicle_id, cooperative_id
   - maintenance_type, description
   - total_cost, status
   - created_at, updated_at
```

**Design Principles:**
- ✅ Always include `cooperative_id` (multi-tenancy)
- ✅ Use UUID for IDs (security, scalability)
- ✅ Add `created_at`, `updated_at` timestamps
- ✅ Use enums for status fields
- ✅ Add indexes on frequently queried columns
- ✅ Foreign keys for relationships

### 1.3 API Endpoint Design

**RESTful Convention:**

```
Vehicles:
GET    /api/v1/modules/vehicle-rental/vehicles           # List all
POST   /api/v1/modules/vehicle-rental/vehicles           # Create
GET    /api/v1/modules/vehicle-rental/vehicles/:id       # Get one
PUT    /api/v1/modules/vehicle-rental/vehicles/:id       # Update
DELETE /api/v1/modules/vehicle-rental/vehicles/:id       # Delete

Bookings:
GET    /api/v1/modules/vehicle-rental/bookings           # List all
POST   /api/v1/modules/vehicle-rental/bookings           # Create
GET    /api/v1/modules/vehicle-rental/bookings/:id       # Get one
PUT    /api/v1/modules/vehicle-rental/bookings/:id       # Update
DELETE /api/v1/modules/vehicle-rental/bookings/:id       # Delete
POST   /api/v1/modules/vehicle-rental/bookings/:id/confirm  # Confirm booking
POST   /api/v1/modules/vehicle-rental/bookings/:id/return   # Process return

Reports:
GET    /api/v1/modules/vehicle-rental/reports/revenue    # Revenue report
GET    /api/v1/modules/vehicle-rental/reports/utilization # Utilization report
```

### 1.4 Frontend Pages Design

**Page Structure:**

```
/dashboard/modules/vehicle-rental/
├── vehicles/                    # Vehicle management
│   ├── page.tsx                # List vehicles
│   ├── create/page.tsx         # Create vehicle form
│   └── [id]/edit/page.tsx      # Edit vehicle form
├── bookings/                    # Booking management
│   ├── page.tsx                # List bookings
│   ├── create/page.tsx         # Create booking form
│   └── [id]/page.tsx           # Booking details
├── maintenance/                 # Maintenance scheduling
│   └── page.tsx                # Maintenance list
└── reports/                     # Reports & analytics
    └── page.tsx                # Revenue dashboard
```

---

## Step 2: Database Design

### 2.1 Create Migration File

**Location:** `backend/migrations/`

**Naming Convention:** `YYYYMMDDHHMMSS_create_module_name_tables.sql`

**Example:** `backend/migrations/20250119100000_create_vehicle_rental_tables.sql`

```sql
-- Migration: Create Vehicle Rental Module Tables
-- Author: Your Name
-- Date: 2025-01-19

-- ============================================================================
-- TABLE: vehicles
-- Description: Store vehicle information for rental management
-- ============================================================================

CREATE TABLE IF NOT EXISTS vehicles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cooperative_id UUID NOT NULL REFERENCES cooperatives(id) ON DELETE CASCADE,
    owner_member_id UUID REFERENCES members(id) ON DELETE SET NULL,

    -- Basic Information
    vehicle_type VARCHAR(50) NOT NULL, -- car, motorcycle, truck
    make VARCHAR(100) NOT NULL,        -- Toyota, Honda, etc.
    model VARCHAR(100) NOT NULL,       -- Avanza, Beat, etc.
    year INTEGER NOT NULL,
    plate_number VARCHAR(20) NOT NULL,
    color VARCHAR(50),

    -- Specifications
    engine_capacity VARCHAR(20),
    transmission VARCHAR(20),          -- Manual, Automatic
    fuel_type VARCHAR(20),             -- Petrol, Diesel
    seating_capacity INTEGER,

    -- Status & Condition
    status VARCHAR(20) DEFAULT 'available' CHECK (status IN ('available', 'rented', 'maintenance', 'retired')),
    condition VARCHAR(20) DEFAULT 'good' CHECK (condition IN ('excellent', 'good', 'fair', 'poor')),
    current_mileage INTEGER DEFAULT 0,

    -- Pricing
    base_rate_hourly DECIMAL(15,2),
    base_rate_daily DECIMAL(15,2),
    base_rate_weekly DECIMAL(15,2),
    base_rate_monthly DECIMAL(15,2),

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT unique_plate_per_coop UNIQUE(cooperative_id, plate_number),
    CONSTRAINT positive_mileage CHECK (current_mileage >= 0),
    CONSTRAINT valid_year CHECK (year >= 1900 AND year <= EXTRACT(YEAR FROM CURRENT_DATE) + 1)
);

-- Indexes for performance
CREATE INDEX idx_vehicles_coop_status ON vehicles(cooperative_id, status);
CREATE INDEX idx_vehicles_owner ON vehicles(owner_member_id);
CREATE INDEX idx_vehicles_plate ON vehicles(plate_number);

-- Comments for documentation
COMMENT ON TABLE vehicles IS 'Vehicles available for rental in cooperative fleet';
COMMENT ON COLUMN vehicles.cooperative_id IS 'Multi-tenant: which cooperative owns this vehicle';
COMMENT ON COLUMN vehicles.owner_member_id IS 'Which member owns this vehicle (for revenue sharing)';
COMMENT ON COLUMN vehicles.status IS 'Current status: available, rented, maintenance, retired';

-- ============================================================================
-- TABLE: rental_contracts
-- Description: Rental booking and contract management
-- ============================================================================

CREATE TABLE IF NOT EXISTS rental_contracts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cooperative_id UUID NOT NULL REFERENCES cooperatives(id) ON DELETE CASCADE,
    contract_number VARCHAR(50) NOT NULL,
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE RESTRICT,

    -- Customer Information
    customer_name VARCHAR(200) NOT NULL,
    customer_phone VARCHAR(20) NOT NULL,
    customer_id_number VARCHAR(50),    -- KTP/SIM
    customer_address TEXT,

    -- Rental Period
    pickup_date TIMESTAMP NOT NULL,
    return_date TIMESTAMP NOT NULL,
    actual_return_date TIMESTAMP,
    duration_hours INTEGER,
    duration_days INTEGER,

    -- Pricing
    rate_type VARCHAR(20) CHECK (rate_type IN ('hourly', 'daily', 'weekly', 'monthly')),
    base_rate DECIMAL(15,2) NOT NULL,
    discount_amount DECIMAL(15,2) DEFAULT 0,
    additional_charges DECIMAL(15,2) DEFAULT 0,
    total_amount DECIMAL(15,2) NOT NULL,

    -- Payment
    deposit_amount DECIMAL(15,2) DEFAULT 0,
    payment_status VARCHAR(20) DEFAULT 'pending' CHECK (payment_status IN ('pending', 'partial', 'paid', 'refunded')),
    paid_amount DECIMAL(15,2) DEFAULT 0,
    outstanding_amount DECIMAL(15,2) DEFAULT 0,

    -- Contract Status
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'ongoing', 'completed', 'cancelled')),

    -- Return Details
    return_fuel_level VARCHAR(20),     -- Empty, 1/4, 1/2, 3/4, Full
    return_mileage INTEGER,
    return_condition VARCHAR(20),      -- excellent, good, damaged
    damage_notes TEXT,
    late_fee DECIMAL(15,2) DEFAULT 0,
    cleaning_fee DECIMAL(15,2) DEFAULT 0,
    damage_fee DECIMAL(15,2) DEFAULT 0,

    -- Staff Tracking
    created_by UUID REFERENCES users(id),
    completed_by UUID REFERENCES users(id),

    -- Additional Info
    notes TEXT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT unique_contract_number UNIQUE(cooperative_id, contract_number),
    CONSTRAINT valid_dates CHECK (return_date > pickup_date),
    CONSTRAINT positive_amounts CHECK (total_amount >= 0 AND deposit_amount >= 0)
);

-- Indexes
CREATE INDEX idx_contracts_coop ON rental_contracts(cooperative_id);
CREATE INDEX idx_contracts_vehicle ON rental_contracts(vehicle_id);
CREATE INDEX idx_contracts_status ON rental_contracts(status);
CREATE INDEX idx_contracts_dates ON rental_contracts(pickup_date, return_date);
CREATE INDEX idx_contracts_number ON rental_contracts(contract_number);

-- Comments
COMMENT ON TABLE rental_contracts IS 'Rental bookings and contracts';
COMMENT ON COLUMN rental_contracts.contract_number IS 'Unique contract identifier (e.g., RNT-20250119-0001)';

-- ============================================================================
-- TABLE: maintenance_records
-- Description: Vehicle maintenance history and scheduling
-- ============================================================================

CREATE TABLE IF NOT EXISTS maintenance_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vehicle_id UUID NOT NULL REFERENCES vehicles(id) ON DELETE CASCADE,
    cooperative_id UUID NOT NULL REFERENCES cooperatives(id) ON DELETE CASCADE,

    -- Maintenance Details
    maintenance_type VARCHAR(50) NOT NULL, -- service, repair, inspection, tire_change
    description TEXT NOT NULL,

    -- Scheduling
    scheduled_date DATE,
    completed_date DATE,
    status VARCHAR(20) DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'in_progress', 'completed', 'cancelled')),

    -- Costs
    parts_cost DECIMAL(15,2) DEFAULT 0,
    labor_cost DECIMAL(15,2) DEFAULT 0,
    total_cost DECIMAL(15,2) NOT NULL,

    -- Service Provider
    workshop_name VARCHAR(200),
    mechanic_name VARCHAR(200),

    -- Vehicle State
    mileage_at_service INTEGER,

    -- Staff & Notes
    performed_by UUID REFERENCES users(id),
    notes TEXT,

    -- Timestamps
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    -- Constraints
    CONSTRAINT positive_costs CHECK (total_cost >= 0)
);

-- Indexes
CREATE INDEX idx_maintenance_vehicle ON maintenance_records(vehicle_id);
CREATE INDEX idx_maintenance_coop ON maintenance_records(cooperative_id);
CREATE INDEX idx_maintenance_dates ON maintenance_records(scheduled_date, completed_date);
CREATE INDEX idx_maintenance_status ON maintenance_records(status);

-- Comments
COMMENT ON TABLE maintenance_records IS 'Vehicle maintenance history and scheduling';

-- ============================================================================
-- FUNCTION: Auto-update updated_at timestamp
-- ============================================================================

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to all tables
CREATE TRIGGER update_vehicles_updated_at BEFORE UPDATE ON vehicles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_rental_contracts_updated_at BEFORE UPDATE ON rental_contracts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_maintenance_records_updated_at BEFORE UPDATE ON maintenance_records
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================================================
-- SEED DATA (Optional - for development/testing)
-- ============================================================================

-- Note: Replace with actual cooperative_id from your database
-- INSERT INTO vehicles (cooperative_id, owner_member_id, ...) VALUES (...);

-- ============================================================================
-- END OF MIGRATION
-- ============================================================================
```

### 2.2 Run Migration

```bash
# Connect to database
psql -U postgres -d koperasi_erp

# Run migration
\i backend/migrations/20250119100000_create_vehicle_rental_tables.sql

# Verify tables created
\dt

# Check table structure
\d vehicles
\d rental_contracts
\d maintenance_records
```

### 2.3 Create Rollback Migration

**File:** `backend/migrations/20250119100000_create_vehicle_rental_tables_down.sql`

```sql
-- Rollback Migration: Drop Vehicle Rental Module Tables
-- Author: Your Name
-- Date: 2025-01-19

-- Drop triggers
DROP TRIGGER IF EXISTS update_vehicles_updated_at ON vehicles;
DROP TRIGGER IF EXISTS update_rental_contracts_updated_at ON rental_contracts;
DROP TRIGGER IF EXISTS update_maintenance_records_updated_at ON maintenance_records;

-- Drop tables (reverse order due to foreign keys)
DROP TABLE IF EXISTS maintenance_records CASCADE;
DROP TABLE IF EXISTS rental_contracts CASCADE;
DROP TABLE IF EXISTS vehicles CASCADE;

-- Drop function if no other tables use it
-- DROP FUNCTION IF EXISTS update_updated_at_column();
```

---

## Step 3: Backend Development

### 3.1 Create Models

**Location:** `backend/internal/models/vehicle.go`

```go
package models

import (
	"time"
	"github.com/google/uuid"
)

// Vehicle represents a vehicle in the rental fleet
type Vehicle struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CooperativeID  uuid.UUID  `gorm:"type:uuid;not null;index:idx_vehicles_coop_status" json:"cooperative_id"`
	OwnerMemberID  *uuid.UUID `gorm:"type:uuid;index:idx_vehicles_owner" json:"owner_member_id"`

	// Basic Information
	VehicleType    string     `gorm:"type:varchar(50);not null" json:"vehicle_type" validate:"required,oneof=car motorcycle truck"`
	Make           string     `gorm:"type:varchar(100);not null" json:"make" validate:"required"`
	Model          string     `gorm:"type:varchar(100);not null" json:"model" validate:"required"`
	Year           int        `gorm:"not null" json:"year" validate:"required,min=1900,max=2026"`
	PlateNumber    string     `gorm:"type:varchar(20);not null;uniqueIndex:idx_plate_coop" json:"plate_number" validate:"required"`
	Color          string     `gorm:"type:varchar(50)" json:"color"`

	// Specifications
	EngineCapacity   string `gorm:"type:varchar(20)" json:"engine_capacity"`
	Transmission     string `gorm:"type:varchar(20)" json:"transmission" validate:"omitempty,oneof=Manual Automatic"`
	FuelType         string `gorm:"type:varchar(20)" json:"fuel_type" validate:"omitempty,oneof=Petrol Diesel Electric"`
	SeatingCapacity  int    `gorm:"" json:"seating_capacity" validate:"omitempty,min=1,max=100"`

	// Status & Condition
	Status          string `gorm:"type:varchar(20);default:'available'" json:"status" validate:"oneof=available rented maintenance retired"`
	Condition       string `gorm:"type:varchar(20);default:'good'" json:"condition" validate:"oneof=excellent good fair poor"`
	CurrentMileage  int    `gorm:"default:0" json:"current_mileage" validate:"min=0"`

	// Pricing
	BaseRateHourly  float64 `gorm:"type:decimal(15,2)" json:"base_rate_hourly" validate:"omitempty,min=0"`
	BaseRateDaily   float64 `gorm:"type:decimal(15,2)" json:"base_rate_daily" validate:"omitempty,min=0"`
	BaseRateWeekly  float64 `gorm:"type:decimal(15,2)" json:"base_rate_weekly" validate:"omitempty,min=0"`
	BaseRateMonthly float64 `gorm:"type:decimal(15,2)" json:"base_rate_monthly" validate:"omitempty,min=0"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Cooperative  Cooperative        `gorm:"foreignKey:CooperativeID" json:"-"`
	OwnerMember  *Member            `gorm:"foreignKey:OwnerMemberID" json:"owner_member,omitempty"`
	Contracts    []RentalContract   `gorm:"foreignKey:VehicleID" json:"contracts,omitempty"`
	Maintenance  []MaintenanceRecord `gorm:"foreignKey:VehicleID" json:"maintenance,omitempty"`
}

// TableName specifies the table name for GORM
func (Vehicle) TableName() string {
	return "vehicles"
}

// RentalContract represents a rental booking
type RentalContract struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CooperativeID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"cooperative_id"`
	ContractNumber  string     `gorm:"type:varchar(50);not null;uniqueIndex:idx_contract_coop" json:"contract_number"`
	VehicleID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"vehicle_id"`

	// Customer Information
	CustomerName      string  `gorm:"type:varchar(200);not null" json:"customer_name" validate:"required"`
	CustomerPhone     string  `gorm:"type:varchar(20);not null" json:"customer_phone" validate:"required"`
	CustomerIDNumber  string  `gorm:"type:varchar(50)" json:"customer_id_number"`
	CustomerAddress   string  `gorm:"type:text" json:"customer_address"`

	// Rental Period
	PickupDate       time.Time  `gorm:"not null" json:"pickup_date" validate:"required"`
	ReturnDate       time.Time  `gorm:"not null" json:"return_date" validate:"required,gtfield=PickupDate"`
	ActualReturnDate *time.Time `json:"actual_return_date"`
	DurationHours    int        `json:"duration_hours"`
	DurationDays     int        `json:"duration_days"`

	// Pricing
	RateType           string  `gorm:"type:varchar(20)" json:"rate_type" validate:"oneof=hourly daily weekly monthly"`
	BaseRate           float64 `gorm:"type:decimal(15,2);not null" json:"base_rate" validate:"required,min=0"`
	DiscountAmount     float64 `gorm:"type:decimal(15,2);default:0" json:"discount_amount" validate:"min=0"`
	AdditionalCharges  float64 `gorm:"type:decimal(15,2);default:0" json:"additional_charges" validate:"min=0"`
	TotalAmount        float64 `gorm:"type:decimal(15,2);not null" json:"total_amount" validate:"required,min=0"`

	// Payment
	DepositAmount      float64 `gorm:"type:decimal(15,2);default:0" json:"deposit_amount" validate:"min=0"`
	PaymentStatus      string  `gorm:"type:varchar(20);default:'pending'" json:"payment_status" validate:"oneof=pending partial paid refunded"`
	PaidAmount         float64 `gorm:"type:decimal(15,2);default:0" json:"paid_amount" validate:"min=0"`
	OutstandingAmount  float64 `gorm:"type:decimal(15,2);default:0" json:"outstanding_amount" validate:"min=0"`

	// Contract Status
	Status string `gorm:"type:varchar(20);default:'pending'" json:"status" validate:"oneof=pending confirmed ongoing completed cancelled"`

	// Return Details
	ReturnFuelLevel  string  `gorm:"type:varchar(20)" json:"return_fuel_level"`
	ReturnMileage    int     `json:"return_mileage"`
	ReturnCondition  string  `gorm:"type:varchar(20)" json:"return_condition"`
	DamageNotes      string  `gorm:"type:text" json:"damage_notes"`
	LateFee          float64 `gorm:"type:decimal(15,2);default:0" json:"late_fee"`
	CleaningFee      float64 `gorm:"type:decimal(15,2);default:0" json:"cleaning_fee"`
	DamageFee        float64 `gorm:"type:decimal(15,2);default:0" json:"damage_fee"`

	// Staff Tracking
	CreatedBy   *uuid.UUID `gorm:"type:uuid" json:"created_by"`
	CompletedBy *uuid.UUID `gorm:"type:uuid" json:"completed_by"`

	// Additional Info
	Notes string `gorm:"type:text" json:"notes"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Cooperative Cooperative `gorm:"foreignKey:CooperativeID" json:"-"`
	Vehicle     Vehicle     `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Creator     *User       `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Completer   *User       `gorm:"foreignKey:CompletedBy" json:"completer,omitempty"`
}

// TableName specifies the table name for GORM
func (RentalContract) TableName() string {
	return "rental_contracts"
}

// BeforeCreate hook to generate contract number
func (r *RentalContract) BeforeCreate(tx *gorm.DB) error {
	if r.ContractNumber == "" {
		// Generate contract number: RNT-20250119-0001
		r.ContractNumber = fmt.Sprintf("RNT-%s-%04d",
			time.Now().Format("20060102"),
			// Query count of today's contracts + 1
			1, // Simplified - implement proper counter
		)
	}

	// Calculate duration
	duration := r.ReturnDate.Sub(r.PickupDate)
	r.DurationHours = int(duration.Hours())
	r.DurationDays = int(duration.Hours() / 24)
	if duration.Hours() > 0 && int(duration.Hours())%24 != 0 {
		r.DurationDays++ // Round up partial days
	}

	// Calculate outstanding amount
	r.OutstandingAmount = r.TotalAmount - r.PaidAmount

	return nil
}

// MaintenanceRecord represents vehicle maintenance history
type MaintenanceRecord struct {
	ID                uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	VehicleID         uuid.UUID  `gorm:"type:uuid;not null;index" json:"vehicle_id"`
	CooperativeID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"cooperative_id"`

	// Maintenance Details
	MaintenanceType string `gorm:"type:varchar(50);not null" json:"maintenance_type" validate:"required"`
	Description     string `gorm:"type:text;not null" json:"description" validate:"required"`

	// Scheduling
	ScheduledDate  *time.Time `gorm:"type:date" json:"scheduled_date"`
	CompletedDate  *time.Time `gorm:"type:date" json:"completed_date"`
	Status         string     `gorm:"type:varchar(20);default:'scheduled'" json:"status" validate:"oneof=scheduled in_progress completed cancelled"`

	// Costs
	PartsCost  float64 `gorm:"type:decimal(15,2);default:0" json:"parts_cost" validate:"min=0"`
	LaborCost  float64 `gorm:"type:decimal(15,2);default:0" json:"labor_cost" validate:"min=0"`
	TotalCost  float64 `gorm:"type:decimal(15,2);not null" json:"total_cost" validate:"required,min=0"`

	// Service Provider
	WorkshopName  string `gorm:"type:varchar(200)" json:"workshop_name"`
	MechanicName  string `gorm:"type:varchar(200)" json:"mechanic_name"`

	// Vehicle State
	MileageAtService int `json:"mileage_at_service"`

	// Staff & Notes
	PerformedBy *uuid.UUID `gorm:"type:uuid" json:"performed_by"`
	Notes       string     `gorm:"type:text" json:"notes"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	Vehicle     Vehicle     `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"`
	Cooperative Cooperative `gorm:"foreignKey:CooperativeID" json:"-"`
	Performer   *User       `gorm:"foreignKey:PerformedBy" json:"performer,omitempty"`
}

// TableName specifies the table name for GORM
func (MaintenanceRecord) TableName() string {
	return "maintenance_records"
}

// BeforeCreate hook to calculate total cost
func (m *MaintenanceRecord) BeforeCreate(tx *gorm.DB) error {
	if m.TotalCost == 0 {
		m.TotalCost = m.PartsCost + m.LaborCost
	}
	return nil
}
```

### 3.2 Create Service Layer

**Location:** `backend/internal/services/vehicle_service.go`

```go
package services

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"your-project/internal/models"
)

type VehicleService struct {
	db *gorm.DB
}

func NewVehicleService(db *gorm.DB) *VehicleService {
	return &VehicleService{db: db}
}

// GetAll retrieves all vehicles for a cooperative
func (s *VehicleService) GetAll(cooperativeID uuid.UUID, filters map[string]interface{}) ([]models.Vehicle, error) {
	var vehicles []models.Vehicle

	query := s.db.Where("cooperative_id = ?", cooperativeID)

	// Apply filters
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if vehicleType, ok := filters["vehicle_type"].(string); ok && vehicleType != "" {
		query = query.Where("vehicle_type = ?", vehicleType)
	}
	if search, ok := filters["search"].(string); ok && search != "" {
		query = query.Where(
			"make ILIKE ? OR model ILIKE ? OR plate_number ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%",
		)
	}

	// Preload relationships
	if err := query.Preload("OwnerMember").Order("created_at DESC").Find(&vehicles).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch vehicles: %w", err)
	}

	return vehicles, nil
}

// GetByID retrieves a single vehicle by ID
func (s *VehicleService) GetByID(cooperativeID, vehicleID uuid.UUID) (*models.Vehicle, error) {
	var vehicle models.Vehicle

	if err := s.db.Where("id = ? AND cooperative_id = ?", vehicleID, cooperativeID).
		Preload("OwnerMember").
		Preload("Contracts").
		Preload("Maintenance").
		First(&vehicle).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("vehicle not found")
		}
		return nil, fmt.Errorf("failed to fetch vehicle: %w", err)
	}

	return &vehicle, nil
}

// Create creates a new vehicle
func (s *VehicleService) Create(cooperativeID uuid.UUID, vehicle *models.Vehicle) error {
	// Set cooperative ID (security: prevent user from setting different cooperative)
	vehicle.CooperativeID = cooperativeID

	// Validate
	if err := validate.Struct(vehicle); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	// Check duplicate plate number
	var count int64
	s.db.Model(&models.Vehicle{}).
		Where("cooperative_id = ? AND plate_number = ?", cooperativeID, vehicle.PlateNumber).
		Count(&count)
	if count > 0 {
		return fmt.Errorf("vehicle with plate number %s already exists", vehicle.PlateNumber)
	}

	// Create vehicle
	if err := s.db.Create(vehicle).Error; err != nil {
		return fmt.Errorf("failed to create vehicle: %w", err)
	}

	return nil
}

// Update updates a vehicle
func (s *VehicleService) Update(cooperativeID, vehicleID uuid.UUID, updates map[string]interface{}) error {
	// Ensure cooperative_id is not changed
	delete(updates, "cooperative_id")
	delete(updates, "id")

	result := s.db.Model(&models.Vehicle{}).
		Where("id = ? AND cooperative_id = ?", vehicleID, cooperativeID).
		Updates(updates)

	if result.Error != nil {
		return fmt.Errorf("failed to update vehicle: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("vehicle not found")
	}

	return nil
}

// Delete soft deletes a vehicle
func (s *VehicleService) Delete(cooperativeID, vehicleID uuid.UUID) error {
	// Check if vehicle has active rentals
	var activeCount int64
	s.db.Model(&models.RentalContract{}).
		Where("vehicle_id = ? AND status IN (?)", vehicleID, []string{"pending", "confirmed", "ongoing"}).
		Count(&activeCount)

	if activeCount > 0 {
		return fmt.Errorf("cannot delete vehicle with active rentals")
	}

	result := s.db.Where("id = ? AND cooperative_id = ?", vehicleID, cooperativeID).
		Delete(&models.Vehicle{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete vehicle: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("vehicle not found")
	}

	return nil
}

// CheckAvailability checks if vehicle is available for given dates
func (s *VehicleService) CheckAvailability(vehicleID uuid.UUID, startDate, endDate time.Time) (bool, error) {
	var count int64

	s.db.Model(&models.RentalContract{}).
		Where("vehicle_id = ?", vehicleID).
		Where("status IN (?)", []string{"confirmed", "ongoing"}).
		Where("pickup_date < ? AND return_date > ?", endDate, startDate).
		Count(&count)

	return count == 0, nil
}
```

**(File continues with RentalService, MaintenanceService...)**

### 3.3 Create HTTP Handlers

**Location:** `backend/internal/handlers/vehicle_handler.go`

```go
package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"your-project/internal/models"
	"your-project/internal/services"
)

type VehicleHandler struct {
	service *services.VehicleService
}

func NewVehicleHandler(service *services.VehicleService) *VehicleHandler {
	return &VehicleHandler{service: service}
}

// GetAll godoc
// @Summary List all vehicles
// @Description Get all vehicles for the authenticated cooperative
// @Tags vehicles
// @Accept json
// @Produce json
// @Param status query string false "Filter by status"
// @Param vehicle_type query string false "Filter by vehicle type"
// @Param search query string false "Search by make, model, or plate"
// @Success 200 {array} models.Vehicle
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/modules/vehicle-rental/vehicles [get]
func (h *VehicleHandler) GetAll(c *gin.Context) {
	// Get cooperative ID from JWT token (set by auth middleware)
	cooperativeID, err := uuid.Parse(c.GetString("cooperative_id"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid cooperative ID"})
		return
	}

	// Parse query filters
	filters := map[string]interface{}{
		"status":       c.Query("status"),
		"vehicle_type": c.Query("vehicle_type"),
		"search":       c.Query("search"),
	}

	vehicles, err := h.service.GetAll(cooperativeID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vehicles,
		"count":   len(vehicles),
	})
}

// GetByID godoc
// @Summary Get vehicle by ID
// @Description Get a single vehicle with all details
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path string true "Vehicle ID"
// @Success 200 {object} models.Vehicle
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/modules/vehicle-rental/vehicles/{id} [get]
func (h *VehicleHandler) GetByID(c *gin.Context) {
	cooperativeID, _ := uuid.Parse(c.GetString("cooperative_id"))
	vehicleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	vehicle, err := h.service.GetByID(cooperativeID, vehicleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    vehicle,
	})
}

// Create godoc
// @Summary Create new vehicle
// @Description Add a new vehicle to the fleet
// @Tags vehicles
// @Accept json
// @Produce json
// @Param vehicle body models.Vehicle true "Vehicle object"
// @Success 201 {object} models.Vehicle
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/modules/vehicle-rental/vehicles [post]
func (h *VehicleHandler) Create(c *gin.Context) {
	cooperativeID, _ := uuid.Parse(c.GetString("cooperative_id"))

	var vehicle models.Vehicle
	if err := c.ShouldBindJSON(&vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Create(cooperativeID, &vehicle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Vehicle created successfully",
		"data":    vehicle,
	})
}

// Update godoc
// @Summary Update vehicle
// @Description Update vehicle details
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path string true "Vehicle ID"
// @Param updates body map[string]interface{} true "Fields to update"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/modules/vehicle-rental/vehicles/{id} [put]
func (h *VehicleHandler) Update(c *gin.Context) {
	cooperativeID, _ := uuid.Parse(c.GetString("cooperative_id"))
	vehicleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.Update(cooperativeID, vehicleID, updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Vehicle updated successfully",
	})
}

// Delete godoc
// @Summary Delete vehicle
// @Description Soft delete a vehicle
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path string true "Vehicle ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/modules/vehicle-rental/vehicles/{id} [delete]
func (h *VehicleHandler) Delete(c *gin.Context) {
	cooperativeID, _ := uuid.Parse(c.GetString("cooperative_id"))
	vehicleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	if err := h.service.Delete(cooperativeID, vehicleID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Vehicle deleted successfully",
	})
}

// CheckAvailability godoc
// @Summary Check vehicle availability
// @Description Check if vehicle is available for given dates
// @Tags vehicles
// @Accept json
// @Produce json
// @Param id path string true "Vehicle ID"
// @Param start_date query string true "Start date (ISO 8601)"
// @Param end_date query string true "End date (ISO 8601)"
// @Success 200 {object} AvailabilityResponse
// @Router /api/v1/modules/vehicle-rental/vehicles/{id}/availability [get]
func (h *VehicleHandler) CheckAvailability(c *gin.Context) {
	vehicleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid vehicle ID"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, c.Query("start_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start_date format"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, c.Query("end_date"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end_date format"})
		return
	}

	available, err := h.service.CheckAvailability(vehicleID, startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":   true,
		"available": available,
	})
}

// Response types for Swagger
type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type AvailabilityResponse struct {
	Success   bool `json:"success"`
	Available bool `json:"available"`
}
```

### 3.4 Register Routes

**Location:** `backend/cmd/api/main.go`

```go
package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"your-project/internal/database"
	"your-project/internal/handlers"
	"your-project/internal/middleware"
	"your-project/internal/services"
)

func main() {
	// Initialize database
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Initialize Gin router
	r := gin.Default()

	// Middleware
	r.Use(middleware.CORS())

	// Public routes (no auth required)
	public := r.Group("/api/v1")
	{
		public.POST("/auth/login", authHandler.Login)
		public.POST("/auth/register", authHandler.Register)
	}

	// Protected routes (auth required)
	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		// Core ERP routes
		protected.GET("/members", memberHandler.GetAll)
		// ... other core routes

		// Module routes
		RegisterModuleRoutes(protected, db)
	}

	// Start server
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

// RegisterModuleRoutes registers all module routes
func RegisterModuleRoutes(r *gin.RouterGroup, db *gorm.DB) {
	// Vehicle Rental Module
	vehicleRental := r.Group("/modules/vehicle-rental")
	{
		// Services
		vehicleService := services.NewVehicleService(db)
		rentalService := services.NewRentalService(db)
		maintenanceService := services.NewMaintenanceService(db)

		// Handlers
		vehicleHandler := handlers.NewVehicleHandler(vehicleService)
		rentalHandler := handlers.NewRentalHandler(rentalService)
		maintenanceHandler := handlers.NewMaintenanceHandler(maintenanceService)

		// Vehicle routes
		vehicles := vehicleRental.Group("/vehicles")
		{
			vehicles.GET("", vehicleHandler.GetAll)
			vehicles.POST("", vehicleHandler.Create)
			vehicles.GET("/:id", vehicleHandler.GetByID)
			vehicles.PUT("/:id", vehicleHandler.Update)
			vehicles.DELETE("/:id", vehicleHandler.Delete)
			vehicles.GET("/:id/availability", vehicleHandler.CheckAvailability)
		}

		// Rental/Booking routes
		bookings := vehicleRental.Group("/bookings")
		{
			bookings.GET("", rentalHandler.GetAll)
			bookings.POST("", rentalHandler.Create)
			bookings.GET("/:id", rentalHandler.GetByID)
			bookings.PUT("/:id", rentalHandler.Update)
			bookings.DELETE("/:id", rentalHandler.Cancel)
			bookings.POST("/:id/confirm", rentalHandler.Confirm)
			bookings.POST("/:id/return", rentalHandler.ProcessReturn)
		}

		// Maintenance routes
		maintenance := vehicleRental.Group("/maintenance")
		{
			maintenance.GET("", maintenanceHandler.GetAll)
			maintenance.POST("", maintenanceHandler.Create)
			maintenance.GET("/:id", maintenanceHandler.GetByID)
			maintenance.PUT("/:id", maintenanceHandler.Update)
			maintenance.DELETE("/:id", maintenanceHandler.Delete)
			maintenance.GET("/upcoming", maintenanceHandler.GetUpcoming)
		}

		// Reports routes
		reports := vehicleRental.Group("/reports")
		{
			reports.GET("/revenue", rentalHandler.RevenueReport)
			reports.GET("/utilization", vehicleHandler.UtilizationReport)
			reports.GET("/member-payouts", rentalHandler.MemberPayoutsReport)
		}
	}

	// Add more modules here...
}
```

### 3.5 Write Unit Tests

**Location:** `backend/internal/services/vehicle_service_test.go`

```go
package services

import (
	"testing"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"your-project/internal/models"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal("Failed to create test database:", err)
	}

	// Migrate schema
	db.AutoMigrate(&models.Cooperative{}, &models.Member{}, &models.Vehicle{})

	return db
}

func TestVehicleService_Create(t *testing.T) {
	db := setupTestDB(t)
	service := NewVehicleService(db)

	// Create test cooperative
	coop := models.Cooperative{Name: "Test Coop"}
	db.Create(&coop)

	t.Run("Success - Create valid vehicle", func(t *testing.T) {
		vehicle := &models.Vehicle{
			VehicleType:    "car",
			Make:           "Toyota",
			Model:          "Avanza",
			Year:           2020,
			PlateNumber:    "B 1234 ABC",
			BaseRateDaily:  300000,
			Status:         "available",
			Condition:      "good",
		}

		err := service.Create(coop.ID, vehicle)
		assert.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, vehicle.ID)
		assert.Equal(t, coop.ID, vehicle.CooperativeID)
	})

	t.Run("Fail - Duplicate plate number", func(t *testing.T) {
		vehicle := &models.Vehicle{
			VehicleType:    "car",
			Make:           "Honda",
			Model:          "Beat",
			Year:           2021,
			PlateNumber:    "B 1234 ABC", // Same as above
			BaseRateDaily:  150000,
		}

		err := service.Create(coop.ID, vehicle)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("Fail - Invalid year", func(t *testing.T) {
		vehicle := &models.Vehicle{
			VehicleType:    "car",
			Make:           "Toyota",
			Model:          "Avanza",
			Year:           1800, // Invalid
			PlateNumber:    "B 5678 DEF",
			BaseRateDaily:  300000,
		}

		err := service.Create(coop.ID, vehicle)
		assert.Error(t, err)
	})
}

func TestVehicleService_GetAll(t *testing.T) {
	db := setupTestDB(t)
	service := NewVehicleService(db)

	// Create test data
	coop1 := models.Cooperative{Name: "Coop 1"}
	coop2 := models.Cooperative{Name: "Coop 2"}
	db.Create(&coop1)
	db.Create(&coop2)

	// Create vehicles for coop1
	db.Create(&models.Vehicle{CooperativeID: coop1.ID, Make: "Toyota", Model: "Avanza", Year: 2020, PlateNumber: "B 1234"})
	db.Create(&models.Vehicle{CooperativeID: coop1.ID, Make: "Honda", Model: "Beat", Year: 2021, PlateNumber: "B 5678"})

	// Create vehicle for coop2
	db.Create(&models.Vehicle{CooperativeID: coop2.ID, Make: "Suzuki", Model: "Ertiga", Year: 2019, PlateNumber: "B 9012"})

	t.Run("Get all vehicles for coop1", func(t *testing.T) {
		vehicles, err := service.GetAll(coop1.ID, map[string]interface{}{})
		assert.NoError(t, err)
		assert.Len(t, vehicles, 2) // Only coop1's vehicles
	})

	t.Run("Get all vehicles for coop2", func(t *testing.T) {
		vehicles, err := service.GetAll(coop2.ID, map[string]interface{}{})
		assert.NoError(t, err)
		assert.Len(t, vehicles, 1) // Only coop2's vehicles
	})

	t.Run("Filter by make", func(t *testing.T) {
		vehicles, err := service.GetAll(coop1.ID, map[string]interface{}{"search": "Toyota"})
		assert.NoError(t, err)
		assert.Len(t, vehicles, 1)
		assert.Equal(t, "Toyota", vehicles[0].Make)
	})
}

func TestVehicleService_CheckAvailability(t *testing.T) {
	db := setupTestDB(t)
	service := NewVehicleService(db)

	// Create test data
	coop := models.Cooperative{Name: "Test Coop"}
	db.Create(&coop)

	vehicle := models.Vehicle{
		CooperativeID: coop.ID,
		Make:          "Toyota",
		Model:         "Avanza",
		Year:          2020,
		PlateNumber:   "B 1234",
	}
	db.Create(&vehicle)

	// Create a booking: Jan 1-5
	booking := models.RentalContract{
		CooperativeID:  coop.ID,
		VehicleID:      vehicle.ID,
		ContractNumber: "TEST-001",
		CustomerName:   "John Doe",
		CustomerPhone:  "08123456789",
		PickupDate:     time.Date(2025, 1, 1, 9, 0, 0, 0, time.UTC),
		ReturnDate:     time.Date(2025, 1, 5, 17, 0, 0, 0, time.UTC),
		BaseRate:       300000,
		TotalAmount:    1200000,
		Status:         "confirmed",
	}
	db.Create(&booking)

	t.Run("Available - No overlap", func(t *testing.T) {
		// Check Jan 6-10 (after existing booking)
		available, err := service.CheckAvailability(
			vehicle.ID,
			time.Date(2025, 1, 6, 9, 0, 0, 0, time.UTC),
			time.Date(2025, 1, 10, 17, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)
		assert.True(t, available)
	})

	t.Run("Not available - Full overlap", func(t *testing.T) {
		// Check Jan 2-4 (inside existing booking)
		available, err := service.CheckAvailability(
			vehicle.ID,
			time.Date(2025, 1, 2, 9, 0, 0, 0, time.UTC),
			time.Date(2025, 1, 4, 17, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)
		assert.False(t, available)
	})

	t.Run("Not available - Partial overlap", func(t *testing.T) {
		// Check Jan 4-8 (overlaps with existing booking)
		available, err := service.CheckAvailability(
			vehicle.ID,
			time.Date(2025, 1, 4, 9, 0, 0, 0, time.UTC),
			time.Date(2025, 1, 8, 17, 0, 0, 0, time.UTC),
		)
		assert.NoError(t, err)
		assert.False(t, available)
	})
}

// Run tests:
// go test ./internal/services -v
```

---

## Step 4: Frontend Development

### 4.1 Create TypeScript Types

**Location:** `frontend/types/vehicle-rental.ts`

```typescript
// Vehicle types
export interface Vehicle {
  id: string;
  cooperative_id: string;
  owner_member_id?: string;

  vehicle_type: 'car' | 'motorcycle' | 'truck';
  make: string;
  model: string;
  year: number;
  plate_number: string;
  color?: string;

  engine_capacity?: string;
  transmission?: 'Manual' | 'Automatic';
  fuel_type?: 'Petrol' | 'Diesel' | 'Electric';
  seating_capacity?: number;

  status: 'available' | 'rented' | 'maintenance' | 'retired';
  condition: 'excellent' | 'good' | 'fair' | 'poor';
  current_mileage: number;

  base_rate_hourly?: number;
  base_rate_daily?: number;
  base_rate_weekly?: number;
  base_rate_monthly?: number;

  created_at: string;
  updated_at: string;

  owner_member?: {
    id: string;
    name: string;
  };
}

// Rental Contract types
export interface RentalContract {
  id: string;
  cooperative_id: string;
  contract_number: string;
  vehicle_id: string;

  customer_name: string;
  customer_phone: string;
  customer_id_number?: string;
  customer_address?: string;

  pickup_date: string;
  return_date: string;
  actual_return_date?: string;
  duration_hours: number;
  duration_days: number;

  rate_type: 'hourly' | 'daily' | 'weekly' | 'monthly';
  base_rate: number;
  discount_amount: number;
  additional_charges: number;
  total_amount: number;

  deposit_amount: number;
  payment_status: 'pending' | 'partial' | 'paid' | 'refunded';
  paid_amount: number;
  outstanding_amount: number;

  status: 'pending' | 'confirmed' | 'ongoing' | 'completed' | 'cancelled';

  return_fuel_level?: string;
  return_mileage?: number;
  return_condition?: string;
  damage_notes?: string;
  late_fee?: number;
  cleaning_fee?: number;
  damage_fee?: number;

  notes?: string;
  created_at: string;
  updated_at: string;

  vehicle?: Vehicle;
}

// Form types
export interface VehicleFormData {
  vehicle_type: string;
  make: string;
  model: string;
  year: number;
  plate_number: string;
  color?: string;

  engine_capacity?: string;
  transmission?: string;
  fuel_type?: string;
  seating_capacity?: number;

  status: string;
  condition: string;
  current_mileage: number;

  base_rate_daily: number;
  owner_member_id?: string;
}

export interface BookingFormData {
  vehicle_id: string;
  customer_name: string;
  customer_phone: string;
  customer_id_number?: string;
  pickup_date: Date;
  return_date: Date;
  rate_type: string;
  deposit_amount: number;
  notes?: string;
}

// API Response types
export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
  count?: number;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    total_pages: number;
  };
}
```

### 4.2 Create API Client

**Location:** `frontend/lib/api/vehicle-rental.ts`

```typescript
import api from '@/lib/api'; // Base axios instance
import type { Vehicle, RentalContract, ApiResponse, PaginatedResponse, VehicleFormData, BookingFormData } from '@/types/vehicle-rental';

// Vehicle API
export const vehicleRentalAPI = {
  // Vehicles
  vehicles: {
    getAll: async (filters?: {
      status?: string;
      vehicle_type?: string;
      search?: string;
    }): Promise<Vehicle[]> => {
      const response = await api.get<ApiResponse<Vehicle[]>>('/modules/vehicle-rental/vehicles', { params: filters });
      return response.data.data;
    },

    getById: async (id: string): Promise<Vehicle> => {
      const response = await api.get<ApiResponse<Vehicle>>(`/modules/vehicle-rental/vehicles/${id}`);
      return response.data.data;
    },

    create: async (data: VehicleFormData): Promise<Vehicle> => {
      const response = await api.post<ApiResponse<Vehicle>>('/modules/vehicle-rental/vehicles', data);
      return response.data.data;
    },

    update: async (id: string, data: Partial<VehicleFormData>): Promise<void> => {
      await api.put(`/modules/vehicle-rental/vehicles/${id}`, data);
    },

    delete: async (id: string): Promise<void> => {
      await api.delete(`/modules/vehicle-rental/vehicles/${id}`);
    },

    checkAvailability: async (id: string, startDate: Date, endDate: Date): Promise<boolean> => {
      const response = await api.get<{ success: boolean; available: boolean }>(
        `/modules/vehicle-rental/vehicles/${id}/availability`,
        {
          params: {
            start_date: startDate.toISOString(),
            end_date: endDate.toISOString(),
          },
        }
      );
      return response.data.available;
    },
  },

  // Bookings/Rentals
  bookings: {
    getAll: async (filters?: { status?: string }): Promise<RentalContract[]> => {
      const response = await api.get<ApiResponse<RentalContract[]>>('/modules/vehicle-rental/bookings', { params: filters });
      return response.data.data;
    },

    getById: async (id: string): Promise<RentalContract> => {
      const response = await api.get<ApiResponse<RentalContract>>(`/modules/vehicle-rental/bookings/${id}`);
      return response.data.data;
    },

    create: async (data: BookingFormData): Promise<RentalContract> => {
      const response = await api.post<ApiResponse<RentalContract>>('/modules/vehicle-rental/bookings', data);
      return response.data.data;
    },

    update: async (id: string, data: Partial<BookingFormData>): Promise<void> => {
      await api.put(`/modules/vehicle-rental/bookings/${id}`, data);
    },

    confirm: async (id: string): Promise<void> => {
      await api.post(`/modules/vehicle-rental/bookings/${id}/confirm`);
    },

    cancel: async (id: string): Promise<void> => {
      await api.delete(`/modules/vehicle-rental/bookings/${id}`);
    },

    processReturn: async (id: string, returnData: {
      return_fuel_level: string;
      return_mileage: number;
      return_condition: string;
      damage_notes?: string;
      late_fee?: number;
      cleaning_fee?: number;
      damage_fee?: number;
    }): Promise<void> => {
      await api.post(`/modules/vehicle-rental/bookings/${id}/return`, returnData);
    },
  },

  // Reports
  reports: {
    revenue: async (startDate: Date, endDate: Date) => {
      const response = await api.get('/modules/vehicle-rental/reports/revenue', {
        params: {
          start_date: startDate.toISOString(),
          end_date: endDate.toISOString(),
        },
      });
      return response.data.data;
    },

    utilization: async (month: number, year: number) => {
      const response = await api.get('/modules/vehicle-rental/reports/utilization', {
        params: { month, year },
      });
      return response.data.data;
    },
  },
};
```

### 4.3 Create React Components

**Location:** `frontend/components/modules/vehicle-rental/VehicleList.tsx`

```typescript
'use client';

import { useState, useEffect } from 'react';
import {
  Box,
  Card,
  CardContent,
  Typography,
  Button,
  Chip,
  IconButton,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  MenuItem,
} from '@mui/material';
import { Add as AddIcon, Edit as EditIcon, Delete as DeleteIcon } from '@mui/icons-material';
import { useRouter } from 'next/navigation';
import { vehicleRentalAPI } from '@/lib/api/vehicle-rental';
import type { Vehicle } from '@/types/vehicle-rental';

export default function VehicleList() {
  const router = useRouter();
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [loading, setLoading] = useState(true);
  const [filters, setFilters] = useState({
    status: '',
    vehicle_type: '',
    search: '',
  });

  useEffect(() => {
    fetchVehicles();
  }, [filters]);

  const fetchVehicles = async () => {
    try {
      setLoading(true);
      const data = await vehicleRentalAPI.vehicles.getAll(filters);
      setVehicles(data);
    } catch (error) {
      console.error('Failed to fetch vehicles:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id: string) => {
    if (!confirm('Are you sure you want to delete this vehicle?')) return;

    try {
      await vehicleRentalAPI.vehicles.delete(id);
      fetchVehicles(); // Refresh list
    } catch (error) {
      console.error('Failed to delete vehicle:', error);
      alert('Failed to delete vehicle. It may have active bookings.');
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'available': return 'success';
      case 'rented': return 'warning';
      case 'maintenance': return 'info';
      case 'retired': return 'default';
      default: return 'default';
    }
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">Vehicles</Typography>
        <Button
          variant="contained"
          startIcon={<AddIcon />}
          onClick={() => router.push('/dashboard/modules/vehicle-rental/vehicles/create')}
        >
          Add Vehicle
        </Button>
      </Box>

      {/* Filters */}
      <Card sx={{ mb: 3 }}>
        <CardContent>
          <Box display="flex" gap={2}>
            <TextField
              label="Search"
              placeholder="Make, model, or plate..."
              value={filters.search}
              onChange={(e) => setFilters({ ...filters, search: e.target.value })}
              sx={{ flex: 1 }}
            />
            <TextField
              select
              label="Status"
              value={filters.status}
              onChange={(e) => setFilters({ ...filters, status: e.target.value })}
              sx={{ minWidth: 150 }}
            >
              <MenuItem value="">All</MenuItem>
              <MenuItem value="available">Available</MenuItem>
              <MenuItem value="rented">Rented</MenuItem>
              <MenuItem value="maintenance">Maintenance</MenuItem>
              <MenuItem value="retired">Retired</MenuItem>
            </TextField>
            <TextField
              select
              label="Type"
              value={filters.vehicle_type}
              onChange={(e) => setFilters({ ...filters, vehicle_type: e.target.value })}
              sx={{ minWidth: 150 }}
            >
              <MenuItem value="">All</MenuItem>
              <MenuItem value="car">Car</MenuItem>
              <MenuItem value="motorcycle">Motorcycle</MenuItem>
              <MenuItem value="truck">Truck</MenuItem>
            </TextField>
          </Box>
        </CardContent>
      </Card>

      {/* Vehicle Table */}
      <TableContainer component={Card}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Vehicle</TableCell>
              <TableCell>Plate Number</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Owner</TableCell>
              <TableCell>Status</TableCell>
              <TableCell>Daily Rate</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {loading ? (
              <TableRow>
                <TableCell colSpan={7} align="center">Loading...</TableCell>
              </TableRow>
            ) : vehicles.length === 0 ? (
              <TableRow>
                <TableCell colSpan={7} align="center">No vehicles found</TableCell>
              </TableRow>
            ) : (
              vehicles.map((vehicle) => (
                <TableRow key={vehicle.id} hover>
                  <TableCell>
                    <Typography variant="body1" fontWeight="medium">
                      {vehicle.make} {vehicle.model}
                    </Typography>
                    <Typography variant="caption" color="text.secondary">
                      {vehicle.year}
                    </Typography>
                  </TableCell>
                  <TableCell>{vehicle.plate_number}</TableCell>
                  <TableCell>
                    <Chip label={vehicle.vehicle_type} size="small" />
                  </TableCell>
                  <TableCell>
                    {vehicle.owner_member?.name || 'N/A'}
                  </TableCell>
                  <TableCell>
                    <Chip
                      label={vehicle.status}
                      size="small"
                      color={getStatusColor(vehicle.status)}
                    />
                  </TableCell>
                  <TableCell>
                    IDR {vehicle.base_rate_daily?.toLocaleString() || 'N/A'}
                  </TableCell>
                  <TableCell align="right">
                    <IconButton
                      size="small"
                      onClick={() => router.push(`/dashboard/modules/vehicle-rental/vehicles/${vehicle.id}/edit`)}
                    >
                      <EditIcon />
                    </IconButton>
                    <IconButton
                      size="small"
                      color="error"
                      onClick={() => handleDelete(vehicle.id)}
                    >
                      <DeleteIcon />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </TableContainer>
    </Box>
  );
}
```

### 4.4 Create Form Components

**Location:** `frontend/components/modules/vehicle-rental/VehicleForm.tsx`

```typescript
'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm, Controller } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import {
  Box,
  Button,
  Card,
  CardContent,
  Grid,
  TextField,
  MenuItem,
  Typography,
} from '@mui/material';
import { vehicleRentalAPI } from '@/lib/api/vehicle-rental';
import type { VehicleFormData } from '@/types/vehicle-rental';

// Validation schema
const vehicleSchema = z.object({
  vehicle_type: z.enum(['car', 'motorcycle', 'truck']),
  make: z.string().min(1, 'Make is required'),
  model: z.string().min(1, 'Model is required'),
  year: z.number().min(1900).max(new Date().getFullYear() + 1),
  plate_number: z.string().min(1, 'Plate number is required'),
  color: z.string().optional(),

  engine_capacity: z.string().optional(),
  transmission: z.enum(['Manual', 'Automatic']).optional(),
  fuel_type: z.enum(['Petrol', 'Diesel', 'Electric']).optional(),
  seating_capacity: z.number().min(1).max(100).optional(),

  status: z.enum(['available', 'rented', 'maintenance', 'retired']),
  condition: z.enum(['excellent', 'good', 'fair', 'poor']),
  current_mileage: z.number().min(0),

  base_rate_daily: z.number().min(0),
  owner_member_id: z.string().optional(),
});

interface Props {
  initialData?: Partial<VehicleFormData>;
  vehicleId?: string;
}

export default function VehicleForm({ initialData, vehicleId }: Props) {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const { control, handleSubmit, formState: { errors } } = useForm<VehicleFormData>({
    resolver: zodResolver(vehicleSchema),
    defaultValues: initialData || {
      vehicle_type: 'car',
      status: 'available',
      condition: 'good',
      current_mileage: 0,
      base_rate_daily: 0,
    },
  });

  const onSubmit = async (data: VehicleFormData) => {
    try {
      setLoading(true);

      if (vehicleId) {
        await vehicleRentalAPI.vehicles.update(vehicleId, data);
      } else {
        await vehicleRentalAPI.vehicles.create(data);
      }

      router.push('/dashboard/modules/vehicle-rental/vehicles');
      router.refresh();
    } catch (error) {
      console.error('Failed to save vehicle:', error);
      alert('Failed to save vehicle. Please try again.');
    } finally {
      setLoading(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)}>
      <Box display="flex" justifyContent="space-between" alignItems="center" mb={3}>
        <Typography variant="h4">
          {vehicleId ? 'Edit Vehicle' : 'Add New Vehicle'}
        </Typography>
        <Box display="flex" gap={2}>
          <Button
            variant="outlined"
            onClick={() => router.back()}
            disabled={loading}
          >
            Cancel
          </Button>
          <Button
            type="submit"
            variant="contained"
            disabled={loading}
          >
            {loading ? 'Saving...' : 'Save'}
          </Button>
        </Box>
      </Box>

      <Card>
        <CardContent>
          <Grid container spacing={3}>
            {/* Basic Information */}
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>Basic Information</Typography>
            </Grid>

            <Grid item xs={12} md={6}>
              <Controller
                name="vehicle_type"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    select
                    fullWidth
                    label="Vehicle Type"
                    error={!!errors.vehicle_type}
                    helperText={errors.vehicle_type?.message}
                  >
                    <MenuItem value="car">Car</MenuItem>
                    <MenuItem value="motorcycle">Motorcycle</MenuItem>
                    <MenuItem value="truck">Truck</MenuItem>
                  </TextField>
                )}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <Controller
                name="plate_number"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    label="Plate Number"
                    error={!!errors.plate_number}
                    helperText={errors.plate_number?.message}
                  />
                )}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <Controller
                name="make"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    label="Make"
                    placeholder="e.g., Toyota"
                    error={!!errors.make}
                    helperText={errors.make?.message}
                  />
                )}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <Controller
                name="model"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    fullWidth
                    label="Model"
                    placeholder="e.g., Avanza"
                    error={!!errors.model}
                    helperText={errors.model?.message}
                  />
                )}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <Controller
                name="year"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    type="number"
                    fullWidth
                    label="Year"
                    error={!!errors.year}
                    helperText={errors.year?.message}
                    onChange={(e) => field.onChange(parseInt(e.target.value))}
                  />
                )}
              />
            </Grid>

            {/* Specifications */}
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>Specifications</Typography>
            </Grid>

            <Grid item xs={12} md={6}>
              <Controller
                name="transmission"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    select
                    fullWidth
                    label="Transmission"
                    error={!!errors.transmission}
                    helperText={errors.transmission?.message}
                  >
                    <MenuItem value="">N/A</MenuItem>
                    <MenuItem value="Manual">Manual</MenuItem>
                    <MenuItem value="Automatic">Automatic</MenuItem>
                  </TextField>
                )}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <Controller
                name="fuel_type"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    select
                    fullWidth
                    label="Fuel Type"
                    error={!!errors.fuel_type}
                    helperText={errors.fuel_type?.message}
                  >
                    <MenuItem value="">N/A</MenuItem>
                    <MenuItem value="Petrol">Petrol</MenuItem>
                    <MenuItem value="Diesel">Diesel</MenuItem>
                    <MenuItem value="Electric">Electric</MenuItem>
                  </TextField>
                )}
              />
            </Grid>

            {/* Status & Pricing */}
            <Grid item xs={12}>
              <Typography variant="h6" gutterBottom>Status & Pricing</Typography>
            </Grid>

            <Grid item xs={12} md={4}>
              <Controller
                name="status"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    select
                    fullWidth
                    label="Status"
                    error={!!errors.status}
                    helperText={errors.status?.message}
                  >
                    <MenuItem value="available">Available</MenuItem>
                    <MenuItem value="rented">Rented</MenuItem>
                    <MenuItem value="maintenance">Maintenance</MenuItem>
                    <MenuItem value="retired">Retired</MenuItem>
                  </TextField>
                )}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <Controller
                name="condition"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    select
                    fullWidth
                    label="Condition"
                    error={!!errors.condition}
                    helperText={errors.condition?.message}
                  >
                    <MenuItem value="excellent">Excellent</MenuItem>
                    <MenuItem value="good">Good</MenuItem>
                    <MenuItem value="fair">Fair</MenuItem>
                    <MenuItem value="poor">Poor</MenuItem>
                  </TextField>
                )}
              />
            </Grid>

            <Grid item xs={12} md={4}>
              <Controller
                name="current_mileage"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    type="number"
                    fullWidth
                    label="Current Mileage (km)"
                    error={!!errors.current_mileage}
                    helperText={errors.current_mileage?.message}
                    onChange={(e) => field.onChange(parseInt(e.target.value))}
                  />
                )}
              />
            </Grid>

            <Grid item xs={12} md={6}>
              <Controller
                name="base_rate_daily"
                control={control}
                render={({ field }) => (
                  <TextField
                    {...field}
                    type="number"
                    fullWidth
                    label="Daily Rate (IDR)"
                    error={!!errors.base_rate_daily}
                    helperText={errors.base_rate_daily?.message}
                    onChange={(e) => field.onChange(parseFloat(e.target.value))}
                  />
                )}
              />
            </Grid>
          </Grid>
        </CardContent>
      </Card>
    </form>
  );
}
```

### 4.5 Create Pages

**Location:** `frontend/app/dashboard/modules/vehicle-rental/vehicles/page.tsx`

```typescript
import VehicleList from '@/components/modules/vehicle-rental/VehicleList';

export default function VehiclesPage() {
  return <VehicleList />;
}
```

**Location:** `frontend/app/dashboard/modules/vehicle-rental/vehicles/create/page.tsx`

```typescript
import VehicleForm from '@/components/modules/vehicle-rental/VehicleForm';

export default function CreateVehiclePage() {
  return <VehicleForm />;
}
```

**Location:** `frontend/app/dashboard/modules/vehicle-rental/vehicles/[id]/edit/page.tsx`

```typescript
import { vehicleRentalAPI } from '@/lib/api/vehicle-rental';
import VehicleForm from '@/components/modules/vehicle-rental/VehicleForm';

export default async function EditVehiclePage({ params }: { params: { id: string } }) {
  const vehicle = await vehicleRentalAPI.vehicles.getById(params.id);

  return <VehicleForm initialData={vehicle} vehicleId={params.id} />;
}
```

---

## Step 5: Integration & Testing

### 5.1 Manual Testing Checklist

```
□ Backend API Testing (using Postman/curl)
  □ Create vehicle (POST /vehicles)
  □ List vehicles (GET /vehicles)
  □ Get single vehicle (GET /vehicles/:id)
  □ Update vehicle (PUT /vehicles/:id)
  □ Delete vehicle (DELETE /vehicles/:id)
  □ Check multi-tenant isolation (different cooperative_id)

□ Frontend Testing
  □ Vehicle list page loads
  □ Create vehicle form works
  □ Edit vehicle form works
  □ Delete vehicle works
  □ Filters work correctly
  □ Validation messages display

□ End-to-End Flows
  □ Complete rental booking flow
  □ Return process flow
  □ Maintenance scheduling flow
  □ Revenue report generation

□ Edge Cases
  □ Duplicate plate number handling
  □ Double booking prevention
  □ Invalid date handling
  □ Permission checks (different roles)
```

### 5.2 Automated Testing

```bash
# Backend tests
cd backend
go test ./... -v -cover

# Frontend tests (if using Jest/Vitest)
cd frontend
npm test

# E2E tests (if using Playwright/Cypress)
npm run test:e2e
```

---

## Step 6: Documentation

### 6.1 API Documentation

Use Swagger/OpenAPI annotations in handlers:

```go
// @Summary Create new vehicle
// @Description Add a new vehicle to the fleet
// @Tags vehicles
// @Accept json
// @Produce json
// @Param vehicle body models.Vehicle true "Vehicle object"
// @Success 201 {object} models.Vehicle
// @Failure 400 {object} ErrorResponse
// @Router /api/v1/modules/vehicle-rental/vehicles [post]
func (h *VehicleHandler) Create(c *gin.Context) {
    // ... implementation
}
```

Generate Swagger docs:
```bash
swag init -g cmd/api/main.go
```

Access at: `http://localhost:8080/swagger/index.html`

### 6.2 User Guide

Create user documentation in `docs/user-guides/vehicle-rental-module.md`

---

## Step 7: Deployment

### 7.1 Staging Deployment

```bash
# 1. Build backend
cd backend
go build -o bin/api cmd/api/main.go

# 2. Run migrations on staging database
psql -h staging-db.example.com -U postgres -d koperasi_erp < migrations/xxx.sql

# 3. Deploy backend to Cloud Run (or your platform)
gcloud run deploy koperasi-api \
  --source . \
  --region asia-southeast2 \
  --allow-unauthenticated

# 4. Build frontend
cd frontend
npm run build

# 5. Deploy frontend to Vercel (or your platform)
vercel --prod
```

### 7.2 Production Deployment

Same as staging, but with production environment variables and database.

---

## Best Practices

### 1. Multi-Tenancy

✅ **ALWAYS filter by cooperative_id:**
```go
// Every query
db.Where("cooperative_id = ?", cooperativeID).Find(&records)
```

❌ **NEVER expose data from other cooperatives:**
```go
// Missing cooperative filter - SECURITY VULNERABILITY!
db.Find(&records) // BAD!
```

### 2. Transaction Management

✅ **Use transactions for multi-step operations:**
```go
err := db.Transaction(func(tx *gorm.DB) error {
    // Step 1: Create rental
    if err := tx.Create(&rental).Error; err != nil {
        return err
    }

    // Step 2: Update vehicle status
    if err := tx.Model(&vehicle).Update("status", "rented").Error; err != nil {
        return err
    }

    // Step 3: Post to accounting
    if err := postToAccounting(tx, rental); err != nil {
        return err
    }

    return nil // Commit
})
```

### 3. Error Handling

✅ **Proper error handling:**
```go
if err != nil {
    if errors.Is(err, gorm.ErrRecordNotFound) {
        return nil, fmt.Errorf("vehicle not found")
    }
    return nil, fmt.Errorf("database error: %w", err)
}
```

### 4. Validation

✅ **Validate on both frontend and backend:**

Frontend (Zod):
```typescript
const schema = z.object({
  year: z.number().min(1900).max(2026),
  plate_number: z.string().min(1),
});
```

Backend (validator):
```go
type Vehicle struct {
    Year int `validate:"required,min=1900,max=2026"`
}
```

### 5. Code Organization

```
✅ Good structure:
backend/
├── models/       # Data structures only
├── services/     # Business logic only
├── handlers/     # HTTP handling only
└── middleware/   # Cross-cutting concerns

❌ Bad structure:
backend/
└── everything_in_one_file.go  # Messy!
```

---

## Common Pitfalls

### 1. Forgetting Cooperative ID Filter

❌ **Problem:**
```go
db.Find(&vehicles) // Gets ALL vehicles from ALL cooperatives!
```

✅ **Solution:**
```go
db.Where("cooperative_id = ?", cooperativeID).Find(&vehicles)
```

### 2. Not Using Transactions

❌ **Problem:**
```go
db.Create(&rental)
db.Update(&vehicle) // If this fails, rental still created!
```

✅ **Solution:**
```go
db.Transaction(func(tx *gorm.DB) error {
    tx.Create(&rental)
    tx.Update(&vehicle)
    return nil
})
```

### 3. Hardcoding IDs

❌ **Problem:**
```go
cooperativeID := "abc-123" // Hardcoded!
```

✅ **Solution:**
```go
cooperativeID := c.GetString("cooperative_id") // From JWT
```

### 4. No Input Validation

❌ **Problem:**
```go
var vehicle Vehicle
c.ShouldBindJSON(&vehicle) // No validation!
db.Create(&vehicle) // SQL error or data corruption
```

✅ **Solution:**
```go
var vehicle Vehicle
if err := c.ShouldBindJSON(&vehicle); err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}
if err := validate.Struct(vehicle); err != nil {
    return c.JSON(400, gin.H{"error": err.Error()})
}
```

---

## Troubleshooting

### Backend Issues

**Problem:** Database connection fails
```
Solution:
1. Check .env file: DATABASE_URL correct?
2. Check PostgreSQL running: docker ps
3. Check credentials: psql -U postgres
```

**Problem:** CORS errors
```
Solution:
Add CORS middleware in main.go:
r.Use(cors.New(cors.Config{
    AllowOrigins: []string{"http://localhost:3000"},
    AllowCredentials: true,
}))
```

**Problem:** JWT token not working
```
Solution:
1. Check JWT_SECRET in .env
2. Verify middleware order (auth before routes)
3. Check token format: "Bearer <token>"
```

### Frontend Issues

**Problem:** API calls fail with 401
```
Solution:
1. Check if logged in
2. Verify token in cookies
3. Check API base URL in lib/api.ts
```

**Problem:** Form validation not working
```
Solution:
1. Check Zod schema matches backend
2. Verify React Hook Form integration
3. Check error messages display
```

---

## Next Steps

After completing this module:

1. ✅ **Create more modules** using this guide as template
2. ✅ **Add module activation system** (toggle modules per cooperative)
3. ✅ **Implement module billing** (track usage, charge fees)
4. ✅ **Build module marketplace** (browse, install modules)
5. ✅ **Create module SDK** (easier module development)

---

**Document Version:** 1.0
**Last Updated:** 2025-01-19
**Author:** Development Team
**Status:** Complete

For questions or improvements, contact the development team.
