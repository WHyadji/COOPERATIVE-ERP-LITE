# MVP Action Plan - 12 Week Sprint to Pilot

## Executive Summary

**Goal**: Launch working system with 10 pilot cooperatives in 12 weeks
**Strategy**: Start simple, ship fast, iterate based on real usage
**Philosophy**: Better than paper books = WIN

---

## FINAL MVP SCOPE (LOCKED) ✅

### What We're Building (8 Core Features)

1. **User Authentication**
   - Username + password
   - 4 roles: Admin, Treasurer, Cashier, Member
   - Simple session management

2. **Member Management**
   - Register members (admin only)
   - Track 3 types of share capital
   - Member list and search

3. **Basic POS**
   - Select product from catalog (50-100 items)
   - Cash payment only
   - Digital receipt (on screen)
   - Daily sales summary

4. **Simple Accounting**
   - Pre-configured COA (standard Indonesian cooperative)
   - Manual transaction entry
   - Trial balance

5. **4 Essential Reports**
   - Member & Capital List
   - Daily Transaction Report
   - Simple Balance Sheet
   - Simple P&L

6. **Member Portal (Web)**
   - View share capital balance
   - View transaction history
   - View announcements

7. **Cooperative Settings**
   - Basic info (name, address, logo)
   - Fiscal year configuration
   - User management

8. **Data Import**
   - Excel template for initial data
   - One-time manual import per cooperative

### What We're NOT Building (Phase 2)

❌ SHU calculation (complex, do it right later)
❌ QRIS payments (cash only for now)
❌ Barcode scanning (manual entry is fine)
❌ Receipt printing (WhatsApp receipt instead)
❌ Inventory automation (manual stock count)
❌ WhatsApp integration (SMS for urgent only)
❌ Offline mode (require internet)
❌ Native mobile app (responsive web only)
❌ Approval workflows (trust-based for pilot)
❌ NIK validation API (manual entry)
❌ Multiple business units (single unit only)
❌ Advanced reports (Cash Flow, GL, etc.)

---

## 12-WEEK DETAILED TIMELINE

### Week 1-2: Foundation (Mar 1-14)

**Monday-Tuesday (Days 1-2)**:
- [ ] Development environment setup
  - Install Go 1.21+
  - Install Node.js 20+
  - Install PostgreSQL 15
  - Install Docker
  - Setup VS Code / GoLand

**Wednesday-Thursday (Days 3-4)**:
- [ ] Backend skeleton
  ```bash
  # Initialize Go project
  go mod init github.com/yourorg/koperasi-erp

  # Install dependencies
  go get github.com/gin-gonic/gin
  go get gorm.io/gorm
  go get gorm.io/driver/postgres
  go get github.com/golang-jwt/jwt/v5
  go get golang.org/x/crypto/bcrypt

  # Project structure
  mkdir -p cmd/api
  mkdir -p internal/{models,handlers,middleware,services,database}
  mkdir -p pkg/{utils,config}
  ```

**Friday (Day 5)**:
- [ ] Frontend skeleton
  ```bash
  # Initialize Next.js
  npx create-next-app@latest koperasi-web --typescript

  # Install dependencies
  npm install axios react-hook-form @hookform/resolvers zod
  npm install @mui/material @emotion/react @emotion/styled
  ```

**Week 2 (Days 6-10)**:
- [ ] Database schema implementation
- [ ] JWT authentication endpoints
- [ ] User CRUD operations
- [ ] Login/logout flow (backend + frontend)
- [ ] Protected route middleware
- [ ] Basic dashboard layout

**Deliverable Week 2**: Login working, can create users

---

### Week 3-4: Core Features (Mar 15-28)

**Week 3 (Days 11-15)**:
- [ ] Cooperative settings (CRUD)
- [ ] Member registration form
- [ ] Member list with pagination
- [ ] Member search functionality
- [ ] Share capital recording
- [ ] Share capital balance display

**Week 4 (Days 16-20)**:
- [ ] Chart of Accounts seeding
- [ ] Manual transaction entry form
- [ ] Transaction validation (debit = credit)
- [ ] Transaction list view
- [ ] Basic trial balance calculation
- [ ] Account balance calculation

**Deliverable Week 4**: Can register members, record transactions, view trial balance

---

### Week 5-6: POS & Reports (Mar 29 - Apr 11)

**Week 5 (Days 21-25)**:
- [ ] Product catalog (CRUD)
- [ ] POS interface (select products, calculate total)
- [ ] Cash payment recording
- [ ] Sale receipt display (on screen)
- [ ] Daily sales summary
- [ ] Link POS sales to accounting (auto journal entry)

**Week 6 (Days 26-30)**:
- [ ] Report 1: Member & Capital List (PDF)
- [ ] Report 2: Daily Transaction Report (PDF)
- [ ] Report 3: Simple Balance Sheet (PDF)
- [ ] Report 4: Simple P&L (PDF)
- [ ] Report export functionality
- [ ] Print preview

**Deliverable Week 6**: POS works, 4 reports generate correctly

---

### Week 7-8: Mobile Web & Polish (Apr 12-25)

**Week 7 (Days 31-35)**:
- [ ] Responsive design (mobile breakpoints)
- [ ] Member portal pages:
  - Share capital balance
  - Transaction history
  - Announcements view
- [ ] Member authentication (separate from staff)
- [ ] Mobile-optimized navigation

**Week 8 (Days 36-40)**:
- [ ] UI polish and consistency
- [ ] Error handling improvements
- [ ] Loading states
- [ ] Success/error messages
- [ ] Basic user guide (in-app)
- [ ] Bug fixes from testing

**Deliverable Week 8**: Fully responsive, member portal works on phones

---

### Week 9-10: Pilot Deployment (Apr 26 - May 9)

**Week 9 (Days 41-45)**:
- [ ] Deploy to Google Cloud Run
- [ ] Setup production database
- [ ] Configure domain and SSL
- [ ] Data migration for first 3 cooperatives:
  - [ ] Cooperative A
  - [ ] Cooperative B
  - [ ] Cooperative C
- [ ] On-site training (Day 1 each)
- [ ] Go live with 3 cooperatives

**Week 10 (Days 46-50)**:
- [ ] Daily monitoring and bug fixes
- [ ] User feedback collection
- [ ] Quick feature adjustments
- [ ] Data migration for next 3:
  - [ ] Cooperative D
  - [ ] Cooperative E
  - [ ] Cooperative F
- [ ] Training and go-live

**Deliverable Week 10**: 6 cooperatives live and using daily

---

### Week 11-12: Complete Rollout & Stabilize (May 10-23)

**Week 11 (Days 51-55)**:
- [ ] Final 4 cooperatives:
  - [ ] Cooperative G
  - [ ] Cooperative H
  - [ ] Cooperative I
  - [ ] Cooperative J
- [ ] All 10 cooperatives trained and live
- [ ] Critical bug fixes
- [ ] Performance optimization

**Week 12 (Days 56-60)**:
- [ ] User satisfaction survey
- [ ] Collect testimonials
- [ ] Create 3 case studies
- [ ] Document lessons learned
- [ ] Plan Phase 2 features
- [ ] Celebrate success! 🎉

**Deliverable Week 12**: All 10 cooperatives using daily, ready to scale

---

## THIS WEEK ACTION ITEMS (Week 0 - Prep Week)

### Monday (Today):
1. ✅ Review and approve this action plan
2. [ ] Schedule site visits to 2-3 local cooperatives
3. [ ] Contact Dinas Koperasi for report templates

### Tuesday-Wednesday:
4. [ ] Visit cooperatives, gather:
   - [ ] Current COA (Excel/paper)
   - [ ] RAT report samples (last year)
   - [ ] Member registration forms
   - [ ] Daily transaction forms
   - [ ] Interview treasurer (30 min)

### Thursday:
5. [ ] Create Product Requirements Document (PRD)
6. [ ] Design simple wireframes (paper or Figma)
7. [ ] Write user stories for Week 1-2 features

### Friday:
8. [ ] Setup development environment
9. [ ] Create GitHub repository
10. [ ] Initialize Go and Next.js projects
11. [ ] Setup local PostgreSQL
12. [ ] First commit!

---

## TECHNOLOGY STACK (CONFIRMED)

### Backend
```
Language:     Go 1.21+
Framework:    Gin (web framework)
ORM:          GORM
Database:     PostgreSQL 15
Auth:         JWT
Password:     bcrypt
Validation:   go-playground/validator
Config:       viper
```

### Frontend
```
Framework:    Next.js 14 (App Router)
Language:     TypeScript
UI Library:   Material-UI (MUI)
Forms:        React Hook Form + Zod
State:        React Context (simple for MVP)
API Client:   Axios
PDF:          jsPDF or react-pdf
Charts:       Recharts (for dashboard)
```

### Infrastructure
```
Cloud:        Google Cloud Platform
Hosting:      Cloud Run (managed containers)
Database:     Cloud SQL (PostgreSQL)
Storage:      Cloud Storage (for uploads)
Domain:       .id domain
SSL:          Let's Encrypt (automatic)
CDN:          Cloud CDN
```

### DevOps
```
Version:      Git + GitHub
CI/CD:        GitHub Actions
Container:    Docker
Monitoring:   Google Cloud Monitoring
Logs:         Cloud Logging
Backup:       Automated daily backups
```

---

## DATABASE SCHEMA (START SIMPLE)

```sql
-- Core 8 tables only

CREATE TABLE cooperatives (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name VARCHAR(255) NOT NULL,
  address TEXT,
  phone VARCHAR(20),
  logo_url VARCHAR(500),
  fiscal_year_start INTEGER DEFAULT 1, -- 1 = January
  settings JSONB,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id UUID NOT NULL REFERENCES cooperatives(id),
  username VARCHAR(100) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  full_name VARCHAR(255) NOT NULL,
  role VARCHAR(50) NOT NULL, -- admin, treasurer, cashier, member
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE members (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id UUID NOT NULL REFERENCES cooperatives(id),
  member_no VARCHAR(50) NOT NULL,
  nik VARCHAR(16),
  full_name VARCHAR(255) NOT NULL,
  phone VARCHAR(20),
  address TEXT,
  join_date DATE NOT NULL,
  status VARCHAR(20) DEFAULT 'active', -- active, inactive, suspended
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(cooperative_id, member_no)
);

CREATE TABLE share_capitals (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id UUID NOT NULL REFERENCES cooperatives(id),
  member_id UUID NOT NULL REFERENCES members(id),
  capital_type VARCHAR(20) NOT NULL, -- pokok, wajib, sukarela
  amount BIGINT NOT NULL, -- in satuan (Rp 1 = 1)
  transaction_date DATE NOT NULL,
  notes TEXT,
  created_by UUID REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id UUID NOT NULL REFERENCES cooperatives(id),
  code VARCHAR(20) NOT NULL,
  name VARCHAR(255) NOT NULL,
  account_type VARCHAR(20) NOT NULL, -- asset, liability, equity, revenue, expense
  parent_id UUID REFERENCES accounts(id),
  is_active BOOLEAN DEFAULT true,
  UNIQUE(cooperative_id, code)
);

CREATE TABLE transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id UUID NOT NULL REFERENCES cooperatives(id),
  transaction_no VARCHAR(50) NOT NULL,
  transaction_date DATE NOT NULL,
  description TEXT,
  total_amount BIGINT NOT NULL,
  transaction_type VARCHAR(50), -- manual, pos_sale, share_capital
  created_by UUID REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(cooperative_id, transaction_no)
);

CREATE TABLE transaction_lines (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
  account_id UUID NOT NULL REFERENCES accounts(id),
  debit BIGINT DEFAULT 0,
  credit BIGINT DEFAULT 0,
  notes TEXT,
  CONSTRAINT check_debit_or_credit CHECK (
    (debit > 0 AND credit = 0) OR (credit > 0 AND debit = 0)
  )
);

CREATE TABLE products (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id UUID NOT NULL REFERENCES cooperatives(id),
  sku VARCHAR(50),
  name VARCHAR(255) NOT NULL,
  description TEXT,
  unit VARCHAR(20) DEFAULT 'pcs',
  price BIGINT NOT NULL,
  stock INTEGER DEFAULT 0,
  is_active BOOLEAN DEFAULT true,
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(cooperative_id, sku)
);

CREATE TABLE sales (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  cooperative_id UUID NOT NULL REFERENCES cooperatives(id),
  sale_no VARCHAR(50) NOT NULL,
  member_id UUID REFERENCES members(id),
  total BIGINT NOT NULL,
  payment_method VARCHAR(20) DEFAULT 'cash',
  transaction_id UUID REFERENCES transactions(id),
  cashier_id UUID REFERENCES users(id),
  created_at TIMESTAMP DEFAULT NOW(),
  UNIQUE(cooperative_id, sale_no)
);

CREATE TABLE sale_items (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  sale_id UUID NOT NULL REFERENCES sales(id) ON DELETE CASCADE,
  product_id UUID NOT NULL REFERENCES products(id),
  quantity INTEGER NOT NULL,
  unit_price BIGINT NOT NULL,
  subtotal BIGINT NOT NULL
);

-- Indexes for performance
CREATE INDEX idx_users_cooperative ON users(cooperative_id);
CREATE INDEX idx_members_cooperative ON members(cooperative_id);
CREATE INDEX idx_members_member_no ON members(member_no);
CREATE INDEX idx_share_capitals_member ON share_capitals(member_id);
CREATE INDEX idx_accounts_cooperative ON accounts(cooperative_id);
CREATE INDEX idx_transactions_cooperative ON transactions(cooperative_id);
CREATE INDEX idx_transactions_date ON transactions(transaction_date);
CREATE INDEX idx_products_cooperative ON products(cooperative_id);
CREATE INDEX idx_sales_cooperative ON sales(cooperative_id);
CREATE INDEX idx_sales_date ON sales(created_at);
```

---

## STANDARD CHART OF ACCOUNTS (SEED DATA)

```sql
-- Seed COA for all cooperatives
INSERT INTO accounts (cooperative_id, code, name, account_type) VALUES
-- ASSETS (1xxx)
('...', '1000', 'ASET', 'asset'),
('...', '1100', 'Kas', 'asset'),
('...', '1200', 'Bank', 'asset'),
('...', '1300', 'Piutang Anggota', 'asset'),
('...', '1400', 'Persediaan Barang', 'asset'),
('...', '1500', 'Aset Tetap', 'asset'),
('...', '1510', 'Peralatan', 'asset'),
('...', '1520', 'Akumulasi Penyusutan', 'asset'),

-- LIABILITIES (2xxx)
('...', '2000', 'KEWAJIBAN', 'liability'),
('...', '2100', 'Hutang Usaha', 'liability'),
('...', '2200', 'Hutang Bank', 'liability'),
('...', '2300', 'Simpanan Sukarela Anggota', 'liability'),

-- EQUITY (3xxx)
('...', '3000', 'MODAL', 'equity'),
('...', '3100', 'Simpanan Pokok', 'equity'),
('...', '3200', 'Simpanan Wajib', 'equity'),
('...', '3300', 'Cadangan', 'equity'),
('...', '3400', 'SHU Tahun Berjalan', 'equity'),

-- REVENUE (4xxx)
('...', '4000', 'PENDAPATAN', 'revenue'),
('...', '4100', 'Penjualan', 'revenue'),
('...', '4200', 'Pendapatan Jasa Simpan Pinjam', 'revenue'),
('...', '4300', 'Pendapatan Lain-lain', 'revenue'),

-- EXPENSES (5xxx)
('...', '5000', 'BIAYA', 'expense'),
('...', '5100', 'Harga Pokok Penjualan', 'expense'),
('...', '5200', 'Biaya Gaji', 'expense'),
('...', '5300', 'Biaya Operasional', 'expense'),
('...', '5400', 'Biaya Listrik & Air', 'expense'),
('...', '5500', 'Biaya Administrasi', 'expense');
```

---

## API ENDPOINTS (MVP ONLY)

```
# Authentication
POST   /api/auth/login
POST   /api/auth/logout
GET    /api/auth/me

# Users
GET    /api/users
POST   /api/users
PUT    /api/users/:id
DELETE /api/users/:id

# Members
GET    /api/members
POST   /api/members
GET    /api/members/:id
PUT    /api/members/:id
DELETE /api/members/:id

# Share Capital
GET    /api/members/:id/share-capital
POST   /api/members/:id/share-capital

# Accounts (COA)
GET    /api/accounts
POST   /api/accounts
PUT    /api/accounts/:id

# Transactions
GET    /api/transactions
POST   /api/transactions
GET    /api/transactions/:id

# POS
GET    /api/products
POST   /api/products
POST   /api/sales
GET    /api/sales/daily-summary

# Reports
GET    /api/reports/members
GET    /api/reports/transactions
GET    /api/reports/balance-sheet
GET    /api/reports/income-statement
GET    /api/reports/trial-balance

# Cooperative Settings
GET    /api/cooperative
PUT    /api/cooperative
```

---

## UI PAGES (MVP ONLY)

```
Public:
/login
/forgot-password

Admin/Staff Dashboard:
/dashboard
/members
/members/new
/members/:id
/share-capital
/transactions
/transactions/new
/pos
/products
/reports
/reports/members
/reports/transactions
/reports/balance-sheet
/reports/income-statement
/settings
/users

Member Portal:
/portal
/portal/share-capital
/portal/transactions
/portal/announcements
```

---

## SUCCESS METRICS

### Week 4 (End of Core Development)
- [ ] Can register 100 members in < 30 minutes
- [ ] Can record 50 transactions correctly
- [ ] Trial balance balances (debit = credit)

### Week 8 (Before Pilot)
- [ ] All 8 core features working
- [ ] 4 reports generate correctly
- [ ] Mobile-responsive (works on phone)
- [ ] < 10 critical bugs

### Week 12 (End of Pilot)
- [ ] 8+ cooperatives using daily
- [ ] 1,000+ transactions recorded
- [ ] Zero data loss
- [ ] 3 strong testimonials
- [ ] CSAT > 7/10
- [ ] Ready to scale to 100 cooperatives

---

## RISK MITIGATION

### Technical Risks

**Risk**: Database performance with multiple cooperatives
- **Mitigation**: Add `cooperative_id` index on all tables
- **Mitigation**: Connection pooling from day 1

**Risk**: Report generation slow
- **Mitigation**: Pre-calculate balances nightly
- **Mitigation**: Queue heavy reports

**Risk**: Cloud costs higher than expected
- **Mitigation**: Start with smallest Cloud Run instance
- **Mitigation**: Monitor daily, set budget alerts

### Operational Risks

**Risk**: Cooperatives don't have accurate data to migrate
- **Mitigation**: Help them reconcile before migration
- **Mitigation**: Start with opening balances only, backfill later

**Risk**: Internet connectivity issues
- **Mitigation**: Choose cooperatives with stable connection for pilot
- **Mitigation**: Mobile data as backup

**Risk**: User resistance to change
- **Mitigation**: Hands-on training
- **Mitigation**: Daily check-ins first week
- **Mitigation**: WhatsApp support group

---

## BUDGET (12-Week MVP)

### Development Costs
- 2 Developers × 3 months × IDR 15M = IDR 90M
- 1 Designer (part-time) × 3 months × IDR 5M = IDR 15M
- **Subtotal**: IDR 105M

### Infrastructure (3 months)
- Google Cloud (Cloud Run + Cloud SQL): IDR 2M/month × 3 = IDR 6M
- Domain + SSL: IDR 500K
- **Subtotal**: IDR 6.5M

### Pilot Support
- Travel (visit 10 cooperatives): IDR 10M
- Training materials: IDR 2M
- **Subtotal**: IDR 12M

### Contingency (20%): IDR 25M

**TOTAL MVP BUDGET**: IDR 148.5M (~USD 10,000)

---

## DAILY STANDUP FORMAT

**Time**: 9:00 AM daily (15 minutes)

**Questions**:
1. What did I complete yesterday?
2. What will I work on today?
3. Any blockers?

**Tools**: Google Meet or Discord

---

## COMMUNICATION PLAN

### Internal Team
- Daily standup (15 min)
- Weekly sprint planning (1 hour)
- Bi-weekly sprint review (30 min)
- Slack/Discord for async communication

### Pilot Cooperatives
- Weekly progress update (email)
- WhatsApp group for quick questions
- On-site visits for critical issues
- Monthly satisfaction survey

### Stakeholders
- Bi-weekly progress report
- Demo at Week 4, 8, 12
- Budget tracking monthly

---

## WHAT CAN GO WRONG & BACKUP PLANS

### Scenario 1: Development takes longer than expected
**Backup**: Reduce pilot from 10 to 5 cooperatives, use extra time for quality

### Scenario 2: Cooperatives drop out of pilot
**Backup**: Have 15 cooperatives on waitlist, replace dropouts immediately

### Scenario 3: Critical bug in production
**Backup**: Rollback capability, daily backups, 4-hour response SLA

### Scenario 4: Cloud costs exceed budget
**Backup**: Migrate to cheaper VPS (DigitalOcean), optimize queries

### Scenario 5: Team member unavailable
**Backup**: Cross-train team, document everything, use freelancers

---

## POST-MVP (Phase 2 Planning)

### Month 4-6 Features (Based on Pilot Feedback)
- SHU calculation and distribution
- QRIS payment integration
- Receipt printing support
- Advanced reports (Cash Flow, GL)
- Mobile native app (React Native)
- Inventory automation
- Multi-business unit support

### Month 7-9 Features (Scale Preparation)
- WhatsApp Business API
- NIK validation (Dukcapil)
- Approval workflows
- Automated billing
- Advanced analytics
- API for third-party integration

### Month 10-12 (Scale to 100 Cooperatives)
- Performance optimization
- Multi-region deployment
- Advanced security
- White-label options
- Marketplace beta

---

## FINAL CHECKLIST BEFORE STARTING

- [ ] ✅ MVP scope locked (no changes!)
- [ ] ✅ Technology stack confirmed
- [ ] ✅ Database schema designed
- [ ] ✅ Team assembled
- [ ] ✅ Budget approved
- [ ] ✅ Development environment ready
- [ ] ✅ GitHub repository created
- [ ] ✅ First cooperative visit scheduled
- [ ] ✅ COA template obtained
- [ ] ✅ All stakeholders aligned

---

## MOTIVATIONAL REMINDER

> "Perfect is the enemy of good. Ship something that works, then make it better."

> "Your competition is paper books and Excel. If your software is 20% better, you win."

> "The best time to start was yesterday. The second best time is NOW."

**Let's build this! 🚀**

---

**Version**: 1.0.0
**Created**: 2025-11-15
**Owner**: Product & Engineering Team
**Status**: READY TO START

**Next Action**: Setup development environment and make first commit!
