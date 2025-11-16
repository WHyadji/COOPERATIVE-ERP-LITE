# COOPERATIVE ERP LITE - System Architecture

## Table of Contents
1. [System Overview](#system-overview)
2. [Architecture Patterns](#architecture-patterns)
3. [Technology Stack](#technology-stack)
4. [Module Architecture](#module-architecture)
5. [Database Design](#database-design)
6. [API Architecture](#api-architecture)
7. [Security Architecture](#security-architecture)
8. [Integration Architecture](#integration-architecture)

---

## 1. System Overview

### High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         PRESENTATION LAYER                       │
├─────────────────────────────────────────────────────────────────┤
│  Web App (React)  │  Mobile App (React Native)  │  Admin Portal │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         API GATEWAY LAYER                        │
├─────────────────────────────────────────────────────────────────┤
│  Authentication  │  Rate Limiting  │  Request Routing  │  Logging│
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                      APPLICATION LAYER                           │
├─────────────────────────────────────────────────────────────────┤
│ Members │ Finance │ Business Units │ Inventory │ Operations │ Reports │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         DATA ACCESS LAYER                        │
├─────────────────────────────────────────────────────────────────┤
│    ORM/Query Builder    │    Caching    │    Data Validation   │
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                          DATA LAYER                              │
├─────────────────────────────────────────────────────────────────┤
│  PostgreSQL  │  Redis Cache  │  File Storage  │  Message Queue  │
└─────────────────────────────────────────────────────────────────┘
```

---

## 2. Architecture Patterns

### 2.1 Microservices Architecture (Modular Monolith Initially)

**Phase 1: Modular Monolith**
- Start with well-separated modules in a single codebase
- Clear module boundaries and interfaces
- Shared database with schema separation

**Phase 2: Transition to Microservices** (if needed)
- Extract high-traffic modules to separate services
- Independent deployment and scaling
- Event-driven communication

### 2.2 Design Patterns

**Domain-Driven Design (DDD)**
- Bounded contexts per business domain
- Aggregate roots for consistency
- Domain events for cross-module communication

**CQRS (Command Query Responsibility Segregation)**
- Separate read and write operations
- Optimized read models for reporting
- Event sourcing for audit trail

**Repository Pattern**
- Abstract data access layer
- Testable business logic
- Database agnostic code

---

## 3. Technology Stack

### 3.1 Backend Stack

**Option A: Node.js Stack (Recommended)**
```
Runtime:        Node.js 20 LTS
Framework:      Express.js / NestJS
ORM:            Prisma / TypeORM
Language:       TypeScript
Validation:     Zod / Joi
Testing:        Jest / Supertest
```

**Option B: PHP Stack**
```
Framework:      Laravel 10+
ORM:            Eloquent
Language:       PHP 8.2+
Testing:        PHPUnit
```

**Option C: Python Stack**
```
Framework:      Django / FastAPI
ORM:            Django ORM / SQLAlchemy
Language:       Python 3.11+
Testing:        Pytest
```

### 3.2 Frontend Stack

**Web Application**
```
Framework:      React 18+ with TypeScript
State Mgmt:     Redux Toolkit / Zustand
UI Library:     Material-UI / Ant Design
Forms:          React Hook Form + Zod
Charts:         Recharts / ApexCharts
Build Tool:     Vite
```

**Mobile Application**
```
Framework:      React Native / Flutter
Navigation:     React Navigation
State:          Redux Toolkit
UI:             React Native Paper
```

### 3.3 Database & Storage

```
Primary DB:     PostgreSQL 15+
Cache:          Redis 7+
File Storage:   MinIO / AWS S3
Search:         ElasticSearch (optional)
Message Queue:  RabbitMQ / Redis Pub/Sub
```

### 3.4 DevOps & Infrastructure

```
Containerization: Docker
Orchestration:    Docker Compose / Kubernetes
CI/CD:           GitHub Actions / GitLab CI
Monitoring:      Prometheus + Grafana
Logging:         ELK Stack (ElasticSearch, Logstash, Kibana)
Backup:          Automated PostgreSQL backups
```

---

## 4. Module Architecture

### 4.1 Member Management Module

```
members/
├── domain/
│   ├── entities/
│   │   ├── Member.ts
│   │   ├── ShareCapital.ts
│   │   └── MemberCard.ts
│   ├── repositories/
│   │   └── IMemberRepository.ts
│   └── services/
│       ├── MemberRegistrationService.ts
│       ├── NIKValidationService.ts
│       └── ShareCapitalService.ts
├── application/
│   ├── commands/
│   │   ├── RegisterMember.ts
│   │   ├── UpdateMember.ts
│   │   └── AddShareCapital.ts
│   ├── queries/
│   │   ├── GetMember.ts
│   │   └── ListMembers.ts
│   └── handlers/
├── infrastructure/
│   ├── persistence/
│   │   ├── MemberRepository.ts
│   │   └── migrations/
│   └── external/
│       └── DukcapilNIKValidator.ts
└── presentation/
    ├── controllers/
    ├── routes/
    └── validators/
```

### 4.2 Financial Management Module

```
finance/
├── domain/
│   ├── entities/
│   │   ├── Account.ts
│   │   ├── JournalEntry.ts
│   │   ├── Transaction.ts
│   │   ├── Budget.ts
│   │   └── SHU.ts
│   ├── value-objects/
│   │   ├── Money.ts
│   │   ├── AccountCode.ts
│   │   └── FiscalYear.ts
│   └── services/
│       ├── AccountingService.ts
│       ├── SHUCalculationService.ts
│       └── BudgetService.ts
├── application/
│   ├── commands/
│   │   ├── CreateJournalEntry.ts
│   │   ├── PostTransaction.ts
│   │   └── CalculateSHU.ts
│   └── queries/
│       ├── GetTrialBalance.ts
│       ├── GetIncomeStatement.ts
│       └── GetBalanceSheet.ts
└── infrastructure/
    └── persistence/
        └── FinanceRepository.ts
```

### 4.3 Business Unit Module

```
business-units/
├── domain/
│   ├── entities/
│   │   ├── BusinessUnit.ts
│   │   ├── RetailUnit.ts
│   │   ├── LoanUnit.ts
│   │   └── AgriTradingUnit.ts
│   └── services/
│       ├── UnitPnLService.ts
│       └── InterUnitTransferService.ts
├── application/
└── infrastructure/
```

### 4.4 Inventory Module

```
inventory/
├── domain/
│   ├── entities/
│   │   ├── Product.ts
│   │   ├── Warehouse.ts
│   │   ├── Stock.ts
│   │   ├── PurchaseOrder.ts
│   │   └── Supplier.ts
│   └── services/
│       ├── StockMovementService.ts
│       └── InventoryValuationService.ts
├── application/
└── infrastructure/
```

---

## 5. Database Design

### 5.1 Core Tables Schema

**Members Schema**
```sql
-- members
CREATE TABLE members (
    id UUID PRIMARY KEY,
    nik VARCHAR(16) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    category ENUM('regular', 'honorary', 'founder'),
    status ENUM('active', 'inactive', 'suspended'),
    join_date DATE NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- share_capital
CREATE TABLE share_capital (
    id UUID PRIMARY KEY,
    member_id UUID REFERENCES members(id),
    capital_type ENUM('pokok', 'wajib', 'sukarela'),
    amount DECIMAL(15,2) NOT NULL,
    transaction_date DATE NOT NULL,
    created_at TIMESTAMP
);
```

**Finance Schema**
```sql
-- chart_of_accounts
CREATE TABLE chart_of_accounts (
    id UUID PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    account_type ENUM('asset', 'liability', 'equity', 'revenue', 'expense'),
    parent_id UUID REFERENCES chart_of_accounts(id),
    is_active BOOLEAN DEFAULT true
);

-- journal_entries
CREATE TABLE journal_entries (
    id UUID PRIMARY KEY,
    entry_number VARCHAR(50) UNIQUE NOT NULL,
    entry_date DATE NOT NULL,
    description TEXT,
    total_debit DECIMAL(15,2) NOT NULL,
    total_credit DECIMAL(15,2) NOT NULL,
    status ENUM('draft', 'posted', 'void'),
    posted_by UUID REFERENCES users(id),
    posted_at TIMESTAMP,
    created_at TIMESTAMP
);

-- journal_entry_lines
CREATE TABLE journal_entry_lines (
    id UUID PRIMARY KEY,
    journal_entry_id UUID REFERENCES journal_entries(id),
    account_id UUID REFERENCES chart_of_accounts(id),
    debit DECIMAL(15,2) DEFAULT 0,
    credit DECIMAL(15,2) DEFAULT 0,
    description TEXT,
    line_number INT
);
```

### 5.2 Database Indexing Strategy

```sql
-- Performance indexes
CREATE INDEX idx_members_nik ON members(nik);
CREATE INDEX idx_members_status ON members(status);
CREATE INDEX idx_journal_entries_date ON journal_entries(entry_date);
CREATE INDEX idx_journal_entry_lines_account ON journal_entry_lines(account_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
```

### 5.3 Database Partitioning (for scale)

```sql
-- Partition large tables by year
CREATE TABLE journal_entries_2025 PARTITION OF journal_entries
    FOR VALUES FROM ('2025-01-01') TO ('2026-01-01');
```

---

## 6. API Architecture

### 6.1 RESTful API Design

**Base URL**: `https://api.cooperative-erp.com/v1`

**Authentication**: JWT Bearer Token

**Standard Response Format**:
```json
{
  "success": true,
  "data": { },
  "message": "Operation successful",
  "meta": {
    "timestamp": "2025-11-15T10:00:00Z",
    "request_id": "uuid"
  }
}
```

**Error Response Format**:
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid input data",
    "details": [
      {
        "field": "nik",
        "message": "NIK must be 16 digits"
      }
    ]
  },
  "meta": {
    "timestamp": "2025-11-15T10:00:00Z",
    "request_id": "uuid"
  }
}
```

### 6.2 API Endpoints Structure

**Members API**
```
POST   /members                    # Register new member
GET    /members                    # List all members
GET    /members/:id                # Get member details
PUT    /members/:id                # Update member
DELETE /members/:id                # Deactivate member
POST   /members/:id/share-capital  # Add share capital
GET    /members/:id/transactions   # Get member transactions
```

**Finance API**
```
POST   /finance/journal-entries    # Create journal entry
GET    /finance/journal-entries    # List journal entries
POST   /finance/journal-entries/:id/post  # Post entry
GET    /finance/trial-balance      # Get trial balance
GET    /finance/income-statement   # Get P&L
GET    /finance/balance-sheet      # Get balance sheet
POST   /finance/shu/calculate      # Calculate SHU
```

### 6.3 API Versioning

- URL versioning: `/v1/`, `/v2/`
- Backward compatibility for at least 2 versions
- Deprecation warnings in response headers

---

## 7. Security Architecture

### 7.1 Authentication & Authorization

**Authentication Methods**:
- JWT (JSON Web Tokens) for API access
- Session-based for web application
- OAuth2 for third-party integrations

**Authorization Model**: Role-Based Access Control (RBAC)
```
Roles:
├── Super Admin (system management)
├── Board Director (full cooperative access)
├── Manager (operational management)
├── Finance Officer (finance module access)
├── Cashier (POS and transactions)
├── Warehouse Staff (inventory access)
└── Member (read-only portal access)
```

### 7.2 Security Measures

**Data Protection**:
- Encryption at rest (AES-256)
- Encryption in transit (TLS 1.3)
- Sensitive data masking (NIK, bank accounts)
- PII data anonymization for reports

**Application Security**:
- Input validation and sanitization
- SQL injection prevention (parameterized queries)
- XSS protection
- CSRF tokens
- Rate limiting (100 requests/minute per user)
- API key rotation

**Audit Trail**:
- All financial transactions logged
- User action tracking
- Data modification history
- Failed login attempts monitoring

---

## 8. Integration Architecture

### 8.1 External Integrations

**NIK Validation (Dukcapil)**
```
Integration Type: REST API
Purpose: Validate member NIK against government database
Data Flow: Outbound
```

**QRIS Payment Gateway**
```
Integration Type: REST API + Webhook
Purpose: Process payments via QRIS
Providers: GoPay, OVO, DANA, ShopeePay
Data Flow: Bidirectional
```

**WhatsApp Business API**
```
Integration Type: REST API
Purpose: Send reports and notifications
Provider: Twilio / Meta
Data Flow: Outbound
```

**Bank Integration**
```
Integration Type: API / File Upload
Purpose: Bank reconciliation
Data Flow: Inbound
```

### 8.2 Integration Patterns

**Event-Driven Architecture**
```
Event Bus: RabbitMQ / Redis Pub/Sub

Events:
- MemberRegistered
- ShareCapitalAdded
- TransactionPosted
- SHUCalculated
- ReportGenerated
```

**Webhook Handling**
```
Endpoint: /webhooks/{provider}
Security: HMAC signature validation
Retry: Exponential backoff (3 attempts)
```

---

## 9. Scalability & Performance

### 9.1 Horizontal Scaling

- Load balancer (Nginx / AWS ALB)
- Stateless application servers
- Database read replicas
- Distributed caching (Redis Cluster)

### 9.2 Performance Optimization

**Caching Strategy**:
- Redis for session storage
- Application-level caching for reports
- Database query result caching
- CDN for static assets

**Database Optimization**:
- Connection pooling
- Query optimization and indexing
- Materialized views for reports
- Database partitioning for large tables

---

## 10. Disaster Recovery & Backup

**Backup Strategy**:
- Automated daily database backups
- Transaction log backups every hour
- 30-day backup retention
- Offsite backup storage

**Recovery Objectives**:
- RPO (Recovery Point Objective): 1 hour
- RTO (Recovery Time Objective): 4 hours

---

## 11. Monitoring & Observability

**Application Monitoring**:
- Uptime monitoring
- API response time tracking
- Error rate monitoring
- Resource utilization

**Business Metrics**:
- Daily transaction volume
- Active user count
- Module usage statistics
- Report generation time

**Alerting**:
- System downtime alerts
- Database connection failures
- Failed payment notifications
- Unusual transaction patterns

---

**Version**: 1.0.0
**Last Updated**: 2025-11-15
