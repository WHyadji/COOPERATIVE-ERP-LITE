# Phase 2 Implementation Guide (Month 4-6)

**Enhanced Features Phase**
**Duration:** 12 weeks (Month 4-6)
**Start:** After MVP is stable with 10 pilot cooperatives
**Goal:** Add advanced features that make the platform indispensable

---

## 📋 Table of Contents

1. [Phase 2 Overview](#phase-2-overview)
2. [Prerequisites](#prerequisites)
3. [Success Criteria](#success-criteria)
4. [Folder Structure](#folder-structure)
5. [Week-by-Week Plan](#week-by-week-plan)
6. [Feature Implementation Details](#feature-implementation-details)
7. [Testing & QA](#testing--qa)
8. [Deployment Strategy](#deployment-strategy)
9. [Progress Tracking](#progress-tracking)

---

## 🎯 Phase 2 Overview

### What We're Building

Phase 2 transforms the basic MVP into a **comprehensive cooperative management platform** by adding:

1. **SHU Calculation Engine** - Automated profit distribution (Sisa Hasil Usaha)
2. **Savings & Loan Module** - Full lending operations with interest calculations
3. **QRIS Payment Integration** - Digital payments (Xendit/Midtrans)
4. **Bank API Integration** - Automated bank reconciliation
5. **Native Mobile App** - React Native app for Android/iOS
6. **Advanced POS Features** - Barcode scanning, receipt printing, cash drawer

---

### Why These Features?

**From MVP Feedback (Month 3):**
- ✅ Cooperatives love the basic system BUT...
- ❌ "We need SHU calculation for RAT (Annual Meeting)"
- ❌ "Members want to pay with QRIS, not just cash"
- ❌ "We need loan tracking for savings cooperatives"
- ❌ "Staff want mobile app to record transactions anywhere"
- ❌ "We need barcode scanner for faster checkout"

**Business Impact:**
- **SHU Engine** → Required for cooperative compliance (RAT reporting)
- **Loans** → Opens market to **savings cooperatives** (50% of market)
- **QRIS** → Increases transaction volume by 30-40%
- **Mobile App** → Enables field operations (mobile cooperatives)
- **Advanced POS** → 3x faster checkout, reduces errors

---

### Target Numbers

| Metric | Month 3 (MVP End) | Month 6 (Phase 2 End) | Growth |
|--------|-------------------|----------------------|--------|
| **Cooperatives** | 10 | 50 | 5x |
| **Active Users** | 40 | 250 | 6.25x |
| **Monthly Transactions** | 2,000 | 15,000 | 7.5x |
| **Revenue (MRR)** | IDR 10M | IDR 125M | 12.5x |
| **Team Size** | 5 | 12 | 2.4x |
| **NPS Score** | 45 | 65+ | Target |

---

## ✅ Prerequisites

### Before Starting Phase 2

**Technical Prerequisites:**
- [ ] MVP is deployed and stable (no critical bugs)
- [ ] 10 pilot cooperatives using system daily
- [ ] Database backups automated
- [ ] Monitoring and alerting in place
- [ ] CI/CD pipeline working
- [ ] Technical debt documented

**Business Prerequisites:**
- [ ] MVP feedback collected from all 10 cooperatives
- [ ] Phase 2 features prioritized based on feedback
- [ ] 20+ cooperatives in sales pipeline
- [ ] Budget approved (IDR 450M for Phase 2)
- [ ] Team expansion completed (5 → 12 people)

**Team Readiness:**
- [ ] 2 additional backend developers onboarded
- [ ] 1 mobile developer hired
- [ ] 1 QA engineer hired
- [ ] 1 DevOps engineer hired
- [ ] Team trained on Phase 2 architecture

---

## 🎯 Success Criteria

### Phase 2 is "DONE" When:

**Technical Success:**
- ✅ All 6 features deployed to production
- ✅ 95%+ uptime maintained
- ✅ Mobile app published to Play Store & App Store
- ✅ QRIS payments processing successfully
- ✅ SHU calculation validated by accountants
- ✅ Loan module handles all interest types
- ✅ Zero critical bugs in production

**Business Success:**
- ✅ 50 cooperatives using the platform
- ✅ 30+ cooperatives using QRIS payments
- ✅ 20+ cooperatives using loan module
- ✅ 15+ cooperatives completed SHU calculation
- ✅ MRR reaches IDR 125M
- ✅ NPS score 65+
- ✅ Churn rate < 5%

**User Success:**
- ✅ SHU calculation saves 20+ hours per cooperative
- ✅ QRIS adoption rate > 60%
- ✅ Mobile app used by 100+ field staff
- ✅ Loan processing time reduced by 50%
- ✅ POS checkout time reduced by 60%

---

## 📁 Folder Structure

### Complete Phase 2 Project Structure

```
COOPERATIVE-ERP-LITE/
│
├── backend/                           # Go backend
│   ├── cmd/
│   │   ├── api/
│   │   │   └── main.go               # API server entry point
│   │   ├── worker/                    # NEW: Background job processor
│   │   │   └── main.go               # Job worker for SHU, loans, etc.
│   │   └── migrate/                   # NEW: Database migration tool
│   │       └── main.go
│   │
│   ├── internal/
│   │   ├── models/                    # Database models
│   │   │   ├── cooperative.go
│   │   │   ├── user.go
│   │   │   ├── member.go
│   │   │   ├── share_capital.go
│   │   │   ├── loan.go               # NEW: Loan model
│   │   │   ├── loan_payment.go       # NEW: Loan payment tracking
│   │   │   ├── savings.go            # NEW: Savings account
│   │   │   ├── shu_calculation.go    # NEW: SHU calculation
│   │   │   ├── shu_distribution.go   # NEW: SHU distribution
│   │   │   └── payment.go            # NEW: Payment transactions
│   │   │
│   │   ├── handlers/                  # HTTP request handlers
│   │   │   ├── auth_handler.go
│   │   │   ├── member_handler.go
│   │   │   ├── loan_handler.go       # NEW: Loan endpoints
│   │   │   ├── savings_handler.go    # NEW: Savings endpoints
│   │   │   ├── shu_handler.go        # NEW: SHU endpoints
│   │   │   ├── payment_handler.go    # NEW: Payment endpoints
│   │   │   └── qris_handler.go       # NEW: QRIS payment endpoints
│   │   │
│   │   ├── services/                  # Business logic
│   │   │   ├── auth_service.go
│   │   │   ├── member_service.go
│   │   │   ├── loan_service.go       # NEW: Loan calculations
│   │   │   ├── interest_calculator.go # NEW: Interest calculation engine
│   │   │   ├── shu_service.go        # NEW: SHU calculation engine
│   │   │   ├── payment_service.go    # NEW: Payment processing
│   │   │   ├── qris_service.go       # NEW: QRIS integration
│   │   │   └── bank_service.go       # NEW: Bank API integration
│   │   │
│   │   ├── middleware/
│   │   │   ├── auth.go
│   │   │   ├── cors.go
│   │   │   └── rate_limit.go         # NEW: API rate limiting
│   │   │
│   │   ├── repository/                # Data access layer
│   │   │   ├── member_repo.go
│   │   │   ├── loan_repo.go          # NEW: Loan data access
│   │   │   ├── savings_repo.go       # NEW: Savings data access
│   │   │   ├── shu_repo.go           # NEW: SHU data access
│   │   │   └── payment_repo.go       # NEW: Payment data access
│   │   │
│   │   ├── config/
│   │   │   ├── config.go
│   │   │   ├── database.go
│   │   │   └── payment_config.go     # NEW: QRIS & bank config
│   │   │
│   │   └── utils/
│   │       ├── jwt.go
│   │       ├── validation.go
│   │       ├── decimal.go            # NEW: Decimal calculations
│   │       └── date.go               # NEW: Indonesian date handling
│   │
│   ├── pkg/                           # Shared packages
│   │   ├── qris/                      # NEW: QRIS SDK wrapper
│   │   │   ├── xendit.go             # Xendit integration
│   │   │   └── midtrans.go           # Midtrans integration
│   │   │
│   │   └── bank/                      # NEW: Bank API wrapper
│   │       ├── bca.go
│   │       ├── bri.go
│   │       └── mandiri.go
│   │
│   ├── migrations/                    # NEW: Database migrations
│   │   ├── 001_initial_schema.sql
│   │   ├── 002_add_loans.sql         # NEW: Loan tables
│   │   ├── 003_add_shu.sql           # NEW: SHU tables
│   │   └── 004_add_payments.sql      # NEW: Payment tables
│   │
│   ├── tests/
│   │   ├── unit/
│   │   │   ├── loan_service_test.go  # NEW: Loan tests
│   │   │   ├── shu_service_test.go   # NEW: SHU tests
│   │   │   └── interest_test.go      # NEW: Interest calculation tests
│   │   │
│   │   └── integration/
│   │       ├── qris_test.go          # NEW: QRIS integration tests
│   │       └── bank_test.go          # NEW: Bank integration tests
│   │
│   ├── docs/
│   │   ├── api/                       # NEW: API documentation
│   │   │   ├── loan_api.md
│   │   │   ├── shu_api.md
│   │   │   └── payment_api.md
│   │   │
│   │   └── examples/                  # NEW: Code examples
│   │       ├── shu_calculation_example.md
│   │       └── loan_calculation_example.md
│   │
│   ├── go.mod
│   ├── go.sum
│   ├── Dockerfile
│   ├── .env.example
│   └── README.md
│
├── frontend/                          # Next.js frontend
│   ├── src/
│   │   ├── app/                       # Next.js 15 app directory
│   │   │   ├── (auth)/
│   │   │   │   ├── login/
│   │   │   │   └── layout.tsx
│   │   │   │
│   │   │   ├── (dashboard)/
│   │   │   │   ├── dashboard/
│   │   │   │   ├── members/
│   │   │   │   ├── loans/             # NEW: Loan management pages
│   │   │   │   │   ├── page.tsx      # Loan list
│   │   │   │   │   ├── new/          # Create loan
│   │   │   │   │   ├── [id]/         # Loan detail
│   │   │   │   │   └── components/   # Loan components
│   │   │   │   │
│   │   │   │   ├── savings/           # NEW: Savings pages
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   └── components/
│   │   │   │   │
│   │   │   │   ├── shu/               # NEW: SHU calculation pages
│   │   │   │   │   ├── page.tsx      # SHU dashboard
│   │   │   │   │   ├── calculate/    # SHU calculation wizard
│   │   │   │   │   ├── history/      # Past calculations
│   │   │   │   │   └── components/
│   │   │   │   │
│   │   │   │   ├── payments/          # NEW: Payment management
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── qris/         # QRIS payment pages
│   │   │   │   │   └── components/
│   │   │   │   │
│   │   │   │   ├── pos/               # Enhanced POS
│   │   │   │   │   ├── page.tsx
│   │   │   │   │   ├── components/
│   │   │   │   │   │   ├── BarcodeScanner.tsx  # NEW
│   │   │   │   │   │   ├── ReceiptPrinter.tsx  # NEW
│   │   │   │   │   │   └── QRISPayment.tsx     # NEW
│   │   │   │   │
│   │   │   │   └── layout.tsx
│   │   │   │
│   │   │   ├── api/                   # API routes (if needed)
│   │   │   │   └── webhook/
│   │   │   │       └── qris/
│   │   │   │           └── route.ts   # NEW: QRIS webhook
│   │   │   │
│   │   │   └── layout.tsx
│   │   │
│   │   ├── components/
│   │   │   ├── ui/                    # Shared UI components
│   │   │   │   ├── Button.tsx
│   │   │   │   ├── Input.tsx
│   │   │   │   ├── Modal.tsx
│   │   │   │   └── QRCode.tsx        # NEW: QR code display
│   │   │   │
│   │   │   ├── loans/                 # NEW: Loan components
│   │   │   │   ├── LoanForm.tsx
│   │   │   │   ├── LoanCalculator.tsx
│   │   │   │   ├── PaymentSchedule.tsx
│   │   │   │   └── LoanSummary.tsx
│   │   │   │
│   │   │   ├── shu/                   # NEW: SHU components
│   │   │   │   ├── SHUCalculationWizard.tsx
│   │   │   │   ├── AllocationEditor.tsx
│   │   │   │   ├── MemberDistribution.tsx
│   │   │   │   └── SHUReport.tsx
│   │   │   │
│   │   │   └── payments/              # NEW: Payment components
│   │   │       ├── QRISButton.tsx
│   │   │       ├── PaymentStatus.tsx
│   │   │       └── BankSelector.tsx
│   │   │
│   │   ├── lib/
│   │   │   ├── api/
│   │   │   │   ├── client.ts
│   │   │   │   ├── loanApi.ts        # NEW: Loan API calls
│   │   │   │   ├── shuApi.ts         # NEW: SHU API calls
│   │   │   │   └── paymentApi.ts     # NEW: Payment API calls
│   │   │   │
│   │   │   ├── utils/
│   │   │   │   ├── currency.ts
│   │   │   │   ├── date.ts
│   │   │   │   ├── interest.ts       # NEW: Interest calculations
│   │   │   │   └── validation.ts
│   │   │   │
│   │   │   └── hooks/
│   │   │       ├── useLoans.ts       # NEW: Loan hooks
│   │   │       ├── useSHU.ts         # NEW: SHU hooks
│   │   │       └── usePayments.ts    # NEW: Payment hooks
│   │   │
│   │   └── types/
│   │       ├── loan.ts               # NEW: Loan types
│   │       ├── shu.ts                # NEW: SHU types
│   │       └── payment.ts            # NEW: Payment types
│   │
│   ├── public/
│   │   ├── icons/
│   │   └── images/
│   │
│   ├── package.json
│   ├── tsconfig.json
│   ├── next.config.js
│   └── README.md
│
├── mobile/                            # NEW: React Native mobile app
│   ├── src/
│   │   ├── screens/
│   │   │   ├── auth/
│   │   │   │   ├── LoginScreen.tsx
│   │   │   │   └── SplashScreen.tsx
│   │   │   │
│   │   │   ├── dashboard/
│   │   │   │   └── DashboardScreen.tsx
│   │   │   │
│   │   │   ├── pos/
│   │   │   │   ├── POSScreen.tsx
│   │   │   │   └── BarcodeScanScreen.tsx
│   │   │   │
│   │   │   ├── members/
│   │   │   │   ├── MemberListScreen.tsx
│   │   │   │   └── MemberDetailScreen.tsx
│   │   │   │
│   │   │   ├── loans/
│   │   │   │   ├── LoanListScreen.tsx
│   │   │   │   ├── LoanDetailScreen.tsx
│   │   │   │   └── PaymentScreen.tsx
│   │   │   │
│   │   │   └── payments/
│   │   │       └── QRISPaymentScreen.tsx
│   │   │
│   │   ├── components/
│   │   │   ├── common/
│   │   │   │   ├── Button.tsx
│   │   │   │   ├── Input.tsx
│   │   │   │   └── Card.tsx
│   │   │   │
│   │   │   ├── pos/
│   │   │   │   ├── ProductGrid.tsx
│   │   │   │   ├── Cart.tsx
│   │   │   │   └── BarcodeScanner.tsx
│   │   │   │
│   │   │   └── payments/
│   │   │       └── QRISDisplay.tsx
│   │   │
│   │   ├── navigation/
│   │   │   ├── AppNavigator.tsx
│   │   │   ├── AuthNavigator.tsx
│   │   │   └── types.ts
│   │   │
│   │   ├── services/
│   │   │   ├── api.ts
│   │   │   ├── auth.ts
│   │   │   ├── storage.ts
│   │   │   └── printer.ts           # Receipt printer integration
│   │   │
│   │   ├── hooks/
│   │   │   ├── useAuth.ts
│   │   │   ├── useBarcode.ts
│   │   │   └── usePrinter.ts
│   │   │
│   │   └── utils/
│   │       ├── currency.ts
│   │       └── validation.ts
│   │
│   ├── android/                      # Android-specific code
│   ├── ios/                          # iOS-specific code
│   ├── package.json
│   ├── app.json
│   ├── babel.config.js
│   └── README.md
│
├── docs/                              # Project documentation
│   ├── phase-2-implementation-guide.md  # THIS FILE
│   ├── technical-stack-versions.md
│   ├── mvp-action-plan.md
│   │
│   ├── phase-2/                       # NEW: Phase 2 specific docs
│   │   ├── loan-module-spec.md       # Loan feature specification
│   │   ├── shu-calculation-spec.md   # SHU calculation specification
│   │   ├── qris-integration-guide.md # QRIS integration guide
│   │   ├── bank-integration-guide.md # Bank API integration
│   │   ├── mobile-app-guide.md       # Mobile app development guide
│   │   └── testing-plan.md           # Phase 2 testing plan
│   │
│   └── api/                           # API documentation
│       ├── loan-api.yaml             # OpenAPI spec for loans
│       ├── shu-api.yaml              # OpenAPI spec for SHU
│       └── payment-api.yaml          # OpenAPI spec for payments
│
├── scripts/                           # Utility scripts
│   ├── seed-phase2-data.sql          # NEW: Phase 2 seed data
│   ├── migrate-to-phase2.sh          # NEW: Migration script
│   └── test-qris.sh                  # NEW: QRIS testing script
│
├── .github/
│   └── workflows/
│       ├── backend-ci.yml
│       ├── frontend-ci.yml
│       └── mobile-ci.yml             # NEW: Mobile CI/CD
│
├── docker-compose.yml
├── docker-compose.phase2.yml          # NEW: Phase 2 services
└── README.md
```

---

## 📅 Week-by-Week Plan

### Month 4: Foundation & Core Features

#### **Week 1: Phase 2 Setup & Planning**

**Goals:**
- Setup Phase 2 development environment
- Expand team and onboard new members
- Finalize technical architecture

**Deliverables:**
- [ ] New team members onboarded (2 backend, 1 mobile, 1 QA, 1 DevOps)
- [ ] Phase 2 folder structure created
- [ ] Database migrations for new tables written
- [ ] QRIS provider selected (Xendit vs Midtrans)
- [ ] Bank API partner selected (BCA, BRI, or Mandiri)
- [ ] Mobile app project initialized (React Native)
- [ ] Development environment ready for all team members

**Tasks:**
```bash
# Backend setup
cd backend
mkdir -p internal/services/{loan,shu,payment,qris,bank}
mkdir -p pkg/{qris,bank}
mkdir -p migrations

# Frontend setup
cd frontend
mkdir -p src/app/(dashboard)/{loans,savings,shu,payments}
mkdir -p src/components/{loans,shu,payments}
mkdir -p src/lib/api

# Mobile app setup
npx react-native init KoperasiMobile
cd mobile
mkdir -p src/{screens,components,services,navigation}
```

**Success Metrics:**
- All 5 new team members productive by Friday
- All development environments working
- Phase 2 kickoff meeting completed

---

#### **Week 2: Savings & Loan Module - Backend**

**Goals:**
- Build loan calculation engine
- Create loan management APIs
- Implement savings account tracking

**Deliverables:**
- [ ] Loan database schema implemented
- [ ] Interest calculation service (Flat, Declining, Effective)
- [ ] Loan CRUD APIs
- [ ] Loan payment tracking
- [ ] Savings account APIs
- [ ] Unit tests for loan calculations

**Key Files to Create:**
```
backend/internal/models/loan.go
backend/internal/models/loan_payment.go
backend/internal/models/savings.go
backend/internal/services/loan_service.go
backend/internal/services/interest_calculator.go
backend/internal/handlers/loan_handler.go
backend/tests/unit/loan_service_test.go
```

**API Endpoints to Build:**
```
POST   /api/loans                    # Create loan
GET    /api/loans                    # List loans
GET    /api/loans/:id                # Get loan detail
PUT    /api/loans/:id                # Update loan
DELETE /api/loans/:id                # Delete loan
POST   /api/loans/:id/payments       # Record payment
GET    /api/loans/:id/schedule       # Get payment schedule
GET    /api/loans/:id/balance        # Get outstanding balance

POST   /api/savings                  # Create savings account
GET    /api/savings/:member_id       # Get member savings
POST   /api/savings/deposit          # Deposit to savings
POST   /api/savings/withdraw         # Withdraw from savings
```

**Success Metrics:**
- All loan calculation tests passing
- API endpoints responding < 200ms
- Interest calculations accurate to 2 decimal places

---

#### **Week 3: Savings & Loan Module - Frontend**

**Goals:**
- Build loan management UI
- Create loan calculator
- Build payment schedule viewer

**Deliverables:**
- [ ] Loan list page with filters
- [ ] Create/edit loan form
- [ ] Loan detail page with payment history
- [ ] Loan calculator component
- [ ] Payment schedule table
- [ ] Savings account UI
- [ ] Deposit/withdrawal forms

**Key Files to Create:**
```
frontend/src/app/(dashboard)/loans/page.tsx
frontend/src/app/(dashboard)/loans/new/page.tsx
frontend/src/app/(dashboard)/loans/[id]/page.tsx
frontend/src/components/loans/LoanForm.tsx
frontend/src/components/loans/LoanCalculator.tsx
frontend/src/components/loans/PaymentSchedule.tsx
frontend/src/lib/api/loanApi.ts
frontend/src/lib/utils/interest.ts
frontend/src/types/loan.ts
```

**UI Components:**
- Loan application form
- Interest rate calculator
- Amortization schedule table
- Payment recording modal
- Outstanding balance dashboard

**Success Metrics:**
- Loan creation flow < 3 minutes
- Calculator shows real-time updates
- All forms validated properly

---

#### **Week 4: SHU Calculation Engine - Backend**

**Goals:**
- Build SHU calculation service
- Implement allocation rules
- Create distribution tracking

**Deliverables:**
- [ ] SHU database schema
- [ ] SHU calculation algorithm
- [ ] Allocation rule engine (reserve, member share, management, etc.)
- [ ] Member contribution tracking
- [ ] Distribution calculation
- [ ] SHU APIs
- [ ] Comprehensive calculation tests

**Key Files to Create:**
```
backend/internal/models/shu_calculation.go
backend/internal/models/shu_distribution.go
backend/internal/services/shu_service.go
backend/internal/handlers/shu_handler.go
backend/tests/unit/shu_service_test.go
```

**SHU Calculation Logic:**
```go
// Pseudocode for SHU calculation
type SHUCalculation struct {
    Period         string
    NetIncome      decimal.Decimal
    Reserve        decimal.Decimal  // 25% (mandatory)
    MemberShare    decimal.Decimal  // 50% (split by contribution)
    ManagementFee  decimal.Decimal  // 10%
    EmployeeFee    decimal.Decimal  // 10%
    Education      decimal.Decimal  // 5%
}

func CalculateSHU(period, netIncome, rules) {
    // 1. Apply allocation rules
    // 2. Calculate member contributions (transactions + capital)
    // 3. Distribute member share proportionally
    // 4. Generate distribution report
}
```

**API Endpoints:**
```
POST   /api/shu/calculate            # Run SHU calculation
GET    /api/shu/calculations         # List calculations
GET    /api/shu/calculations/:id     # Get calculation detail
GET    /api/shu/distributions/:id    # Get member distributions
POST   /api/shu/distributions/:id/approve  # Approve distribution
GET    /api/shu/reports/:id          # Generate SHU report
```

**Success Metrics:**
- Calculation completes in < 5 seconds for 1000 members
- All allocation percentages total 100%
- Member distributions accurate to Rupiah

---

### Month 5: Payments & Mobile

#### **Week 5: SHU Calculation - Frontend**

**Goals:**
- Build SHU calculation wizard
- Create distribution UI
- Build SHU reports

**Deliverables:**
- [ ] SHU calculation wizard (multi-step form)
- [ ] Allocation rule editor
- [ ] Member contribution preview
- [ ] Distribution table
- [ ] SHU report generator (PDF)
- [ ] Historical calculation viewer

**Key Files to Create:**
```
frontend/src/app/(dashboard)/shu/page.tsx
frontend/src/app/(dashboard)/shu/calculate/page.tsx
frontend/src/app/(dashboard)/shu/history/page.tsx
frontend/src/components/shu/SHUCalculationWizard.tsx
frontend/src/components/shu/AllocationEditor.tsx
frontend/src/components/shu/MemberDistribution.tsx
frontend/src/components/shu/SHUReport.tsx
frontend/src/lib/api/shuApi.ts
frontend/src/types/shu.ts
```

**Wizard Steps:**
1. Select period & net income
2. Configure allocation rules
3. Review member contributions
4. Preview distributions
5. Confirm & calculate
6. Generate report

**Success Metrics:**
- Wizard completion < 10 minutes
- Report generates in < 30 seconds
- Report printable and exportable (PDF)

---

#### **Week 6: QRIS Integration - Backend**

**Goals:**
- Integrate QRIS payment provider
- Build payment processing service
- Implement webhook handling

**Deliverables:**
- [ ] QRIS provider SDK integration (Xendit or Midtrans)
- [ ] Payment creation API
- [ ] Payment status tracking
- [ ] Webhook endpoint for payment notifications
- [ ] Payment reconciliation
- [ ] Refund handling

**Key Files to Create:**
```
backend/pkg/qris/xendit.go           # Or midtrans.go
backend/internal/models/payment.go
backend/internal/services/qris_service.go
backend/internal/services/payment_service.go
backend/internal/handlers/payment_handler.go
backend/internal/handlers/qris_handler.go
backend/tests/integration/qris_test.go
```

**QRIS Flow:**
```
1. Customer wants to pay → Generate QR code
2. Customer scans QR → Xendit/Midtrans processes
3. Payment successful → Webhook notification received
4. Update transaction status → Record in database
5. Send receipt → Email/WhatsApp notification
```

**API Endpoints:**
```
POST   /api/payments/qris/create     # Create QRIS payment
GET    /api/payments/:id/status      # Check payment status
POST   /api/payments/webhook/qris    # QRIS webhook (public)
POST   /api/payments/:id/refund      # Refund payment
GET    /api/payments/transactions    # List all payments
```

**Success Metrics:**
- QR code generates in < 2 seconds
- Webhook processes in < 1 second
- 99%+ payment success rate in testing

---

#### **Week 7: QRIS Integration - Frontend & POS**

**Goals:**
- Add QRIS to POS
- Build payment UI
- Create payment status tracking

**Deliverables:**
- [ ] QRIS button in POS
- [ ] QR code display component
- [ ] Payment status polling
- [ ] Payment success/failure handling
- [ ] Payment history page
- [ ] Receipt with QR code

**Key Files to Create:**
```
frontend/src/app/(dashboard)/pos/components/QRISPayment.tsx
frontend/src/app/(dashboard)/payments/page.tsx
frontend/src/components/payments/QRISButton.tsx
frontend/src/components/payments/PaymentStatus.tsx
frontend/src/components/ui/QRCode.tsx
frontend/src/lib/api/paymentApi.ts
frontend/src/types/payment.ts
```

**POS Payment Flow:**
1. Staff adds items to cart
2. Customer chooses QRIS payment
3. System generates QR code
4. Display QR on screen/receipt printer
5. Customer scans and pays
6. Real-time status update (polling every 2s)
7. Show success message
8. Print receipt

**Success Metrics:**
- QRIS flow completion < 60 seconds
- Payment status updates in real-time
- QR code readable on all devices

---

#### **Week 8: Mobile App - Core Features**

**Goals:**
- Build mobile app foundation
- Implement auth and navigation
- Build mobile POS

**Deliverables:**
- [ ] React Native project setup complete
- [ ] Login/logout flow
- [ ] Bottom tab navigation
- [ ] Mobile POS screen
- [ ] Barcode scanner integration
- [ ] Cart management
- [ ] Cash & QRIS payment
- [ ] Offline mode (basic)

**Key Files to Create:**
```
mobile/src/screens/auth/LoginScreen.tsx
mobile/src/screens/dashboard/DashboardScreen.tsx
mobile/src/screens/pos/POSScreen.tsx
mobile/src/screens/pos/BarcodeScanScreen.tsx
mobile/src/components/pos/ProductGrid.tsx
mobile/src/components/pos/Cart.tsx
mobile/src/components/pos/BarcodeScanner.tsx
mobile/src/navigation/AppNavigator.tsx
mobile/src/services/api.ts
mobile/src/services/storage.ts
```

**Mobile POS Features:**
- Product search with barcode
- Manual product selection
- Cart with quantity adjustment
- Multiple payment methods
- Receipt generation
- Sync with backend

**Success Metrics:**
- App runs on Android & iOS
- Barcode scanning works smoothly
- Offline cart persists
- Sync works when online

---

### Month 6: Advanced Features & Launch

#### **Week 9: Mobile App - Advanced Features**

**Goals:**
- Add member management to mobile
- Build loan module for mobile
- Add receipt printing

**Deliverables:**
- [ ] Member list and search
- [ ] Member detail view
- [ ] Loan list for mobile
- [ ] Loan payment recording
- [ ] Receipt printer integration (Bluetooth)
- [ ] Push notifications
- [ ] Background sync

**Key Files to Create:**
```
mobile/src/screens/members/MemberListScreen.tsx
mobile/src/screens/members/MemberDetailScreen.tsx
mobile/src/screens/loans/LoanListScreen.tsx
mobile/src/screens/loans/PaymentScreen.tsx
mobile/src/services/printer.ts
mobile/src/services/notifications.ts
mobile/src/hooks/usePrinter.ts
```

**Receipt Printing:**
- Support Bluetooth thermal printers
- Receipt templates (sale, loan payment)
- Printer configuration
- Test print functionality

**Success Metrics:**
- All core features work offline
- Receipt prints in < 5 seconds
- Push notifications delivered

---

#### **Week 10: Bank API Integration & Advanced POS**

**Goals:**
- Integrate bank APIs for reconciliation
- Add advanced POS hardware
- Build cash drawer integration

**Deliverables:**
- [ ] Bank API integration (BCA/BRI/Mandiri)
- [ ] Automatic transaction matching
- [ ] Bank reconciliation report
- [ ] Barcode scanner (USB) integration
- [ ] Receipt printer (USB/LAN) integration
- [ ] Cash drawer control
- [ ] Weight scale integration (optional)

**Key Files to Create:**
```
backend/pkg/bank/bca.go
backend/pkg/bank/bri.go
backend/internal/services/bank_service.go
backend/internal/handlers/bank_handler.go
frontend/src/app/(dashboard)/pos/components/BarcodeScanner.tsx
frontend/src/app/(dashboard)/pos/components/ReceiptPrinter.tsx
```

**Bank Integration Flow:**
1. Fetch bank statements via API
2. Match with system transactions
3. Flag unmatched items
4. Generate reconciliation report
5. Export for accounting

**Advanced POS Hardware:**
- USB barcode scanner (plug & play)
- Thermal receipt printer (ESC/POS protocol)
- Cash drawer (triggered by printer)
- Weight scale (for groceries, optional)

**Success Metrics:**
- Bank reconciliation 95%+ accurate
- Barcode scan < 1 second
- Receipt prints in < 3 seconds

---

#### **Week 11: Testing & Bug Fixing**

**Goals:**
- Comprehensive testing all Phase 2 features
- Fix critical bugs
- Performance optimization

**Deliverables:**
- [ ] All unit tests passing
- [ ] Integration tests for QRIS & bank APIs
- [ ] End-to-end testing (Playwright/Cypress)
- [ ] Mobile app testing (iOS & Android)
- [ ] Load testing (1000+ concurrent users)
- [ ] Security audit
- [ ] Bug fixes completed

**Testing Checklist:**

**Loan Module:**
- [ ] Interest calculations verified by accountant
- [ ] Payment schedules accurate
- [ ] Outstanding balances correct
- [ ] Late payment penalties work

**SHU Module:**
- [ ] Allocation rules total 100%
- [ ] Member distributions accurate
- [ ] Reports match manual calculations
- [ ] Edge cases handled (negative income, zero members)

**QRIS Payments:**
- [ ] QR codes generate correctly
- [ ] Webhook handling reliable
- [ ] Payment status updates real-time
- [ ] Refunds process correctly

**Mobile App:**
- [ ] Works on Android 8+
- [ ] Works on iOS 13+
- [ ] Offline mode functional
- [ ] Barcode scanning accurate
- [ ] Receipt printing works

**Performance:**
- [ ] API response time < 200ms (p95)
- [ ] Page load time < 2s
- [ ] Mobile app startup < 3s
- [ ] SHU calculation < 5s for 1000 members

**Success Metrics:**
- Zero critical bugs
- All tests passing
- Performance targets met

---

#### **Week 12: Deployment & Rollout**

**Goals:**
- Deploy Phase 2 to production
- Migrate MVP cooperatives
- Onboard new cooperatives

**Deliverables:**
- [ ] Production deployment completed
- [ ] Database migrations run successfully
- [ ] 10 MVP cooperatives upgraded to Phase 2
- [ ] 40 new cooperatives onboarded
- [ ] Mobile app published to Play Store & App Store
- [ ] Team trained on new features
- [ ] Documentation updated
- [ ] Marketing materials ready

**Deployment Checklist:**

**Backend:**
- [ ] Database backup taken
- [ ] Migrations tested in staging
- [ ] Environment variables configured
- [ ] QRIS credentials configured
- [ ] Bank API credentials configured
- [ ] Background workers running
- [ ] Monitoring alerts configured

**Frontend:**
- [ ] Build optimized for production
- [ ] CDN configured
- [ ] Analytics tracking enabled
- [ ] Error tracking (Sentry) configured

**Mobile:**
- [ ] Android APK built and signed
- [ ] iOS IPA built and signed
- [ ] Play Store listing created
- [ ] App Store listing created
- [ ] Beta testing completed
- [ ] Apps published

**Rollout Plan:**
1. **Week 12, Monday-Tuesday:** Deploy to production
2. **Week 12, Wednesday:** Upgrade 10 MVP cooperatives
3. **Week 12, Thursday-Friday:** Onboard 5 new cooperatives
4. **Week 13-15:** Onboard remaining 35 cooperatives
5. **Week 16:** Retrospective and Phase 3 planning

**Success Metrics:**
- Zero downtime during deployment
- All 10 MVP cooperatives successfully upgraded
- 50 total cooperatives by end of Month 6
- NPS score 65+

---

## 🔧 Feature Implementation Details

### Feature 1: SHU Calculation Engine

**What It Does:**
Automates the calculation and distribution of **Sisa Hasil Usaha** (cooperative profit sharing) according to Indonesian cooperative law.

**Business Rules:**
```
Net Income = Revenue - Expenses

Mandatory Allocations:
- Reserve Fund: 25% (UU No. 25/1992)
- Member Share: 50% minimum
- Other allocations: 25% (flexible)

Member Share Distribution:
- Based on member contributions:
  - Transaction volume (50%)
  - Share capital (50%)

Example:
Net Income: IDR 100,000,000
- Reserve: IDR 25,000,000 (25%)
- Member Share: IDR 50,000,000 (50%)
- Management: IDR 10,000,000 (10%)
- Employee: IDR 10,000,000 (10%)
- Education: IDR 5,000,000 (5%)

Member A:
- Transactions: IDR 10M (10% of total)
- Share capital: IDR 5M (5% of total)
- Contribution score: (10% + 5%) / 2 = 7.5%
- Member share: IDR 50M × 7.5% = IDR 3,750,000
```

**Database Schema:**
```sql
CREATE TABLE shu_calculations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cooperative_id UUID NOT NULL REFERENCES cooperatives(id),
    period VARCHAR(10) NOT NULL,  -- e.g., "2025"
    calculation_date DATE NOT NULL,
    net_income DECIMAL(15,2) NOT NULL,

    -- Allocations
    reserve_amount DECIMAL(15,2) NOT NULL,
    reserve_percentage DECIMAL(5,2) NOT NULL DEFAULT 25.00,

    member_share_amount DECIMAL(15,2) NOT NULL,
    member_share_percentage DECIMAL(5,2) NOT NULL DEFAULT 50.00,

    management_amount DECIMAL(15,2),
    management_percentage DECIMAL(5,2),

    employee_amount DECIMAL(15,2),
    employee_percentage DECIMAL(5,2),

    education_amount DECIMAL(15,2),
    education_percentage DECIMAL(5,2),

    other_amount DECIMAL(15,2),
    other_percentage DECIMAL(5,2),

    -- Metadata
    status VARCHAR(20) DEFAULT 'draft',  -- draft, approved, distributed
    calculated_by UUID REFERENCES users(id),
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,
    notes TEXT,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE shu_distributions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    calculation_id UUID NOT NULL REFERENCES shu_calculations(id),
    member_id UUID NOT NULL REFERENCES members(id),

    -- Member contributions
    total_transactions DECIMAL(15,2) NOT NULL,
    transaction_percentage DECIMAL(5,4) NOT NULL,

    total_capital DECIMAL(15,2) NOT NULL,
    capital_percentage DECIMAL(5,4) NOT NULL,

    -- Distribution
    contribution_score DECIMAL(5,4) NOT NULL,
    distribution_amount DECIMAL(15,2) NOT NULL,

    -- Payment
    payment_status VARCHAR(20) DEFAULT 'pending',  -- pending, paid
    payment_date DATE,
    payment_method VARCHAR(50),

    created_at TIMESTAMP DEFAULT NOW()
);
```

**Go Implementation:**
```go
// backend/internal/services/shu_service.go
package services

import (
    "github.com/shopspring/decimal"
    "koperasi-erp/internal/models"
)

type SHUService struct {
    repo SHURepository
}

type SHUAllocationRules struct {
    Reserve        decimal.Decimal  // 0.25 (25%)
    MemberShare    decimal.Decimal  // 0.50 (50%)
    Management     decimal.Decimal  // 0.10 (10%)
    Employee       decimal.Decimal  // 0.10 (10%)
    Education      decimal.Decimal  // 0.05 (5%)
}

func (s *SHUService) CalculateSHU(
    cooperativeID string,
    period string,
    netIncome decimal.Decimal,
    rules SHUAllocationRules,
) (*models.SHUCalculation, error) {

    // 1. Validate allocation rules total 100%
    total := rules.Reserve.Add(rules.MemberShare).
            Add(rules.Management).Add(rules.Employee).
            Add(rules.Education)

    if !total.Equal(decimal.NewFromInt(1)) {
        return nil, errors.New("allocations must total 100%")
    }

    // 2. Calculate allocations
    calculation := &models.SHUCalculation{
        CooperativeID: cooperativeID,
        Period:        period,
        NetIncome:     netIncome,

        ReserveAmount:     netIncome.Mul(rules.Reserve),
        ReservePercentage: rules.Reserve.Mul(decimal.NewFromInt(100)),

        MemberShareAmount:     netIncome.Mul(rules.MemberShare),
        MemberSharePercentage: rules.MemberShare.Mul(decimal.NewFromInt(100)),

        ManagementAmount:     netIncome.Mul(rules.Management),
        ManagementPercentage: rules.Management.Mul(decimal.NewFromInt(100)),

        EmployeeAmount:     netIncome.Mul(rules.Employee),
        EmployeePercentage: rules.Employee.Mul(decimal.NewFromInt(100)),

        EducationAmount:     netIncome.Mul(rules.Education),
        EducationPercentage: rules.Education.Mul(decimal.NewFromInt(100)),

        Status: "draft",
    }

    // 3. Save calculation
    if err := s.repo.CreateCalculation(calculation); err != nil {
        return nil, err
    }

    // 4. Calculate member distributions
    distributions, err := s.calculateMemberDistributions(
        calculation.ID,
        cooperativeID,
        period,
        calculation.MemberShareAmount,
    )
    if err != nil {
        return nil, err
    }

    // 5. Save distributions
    for _, dist := range distributions {
        if err := s.repo.CreateDistribution(dist); err != nil {
            return nil, err
        }
    }

    return calculation, nil
}

func (s *SHUService) calculateMemberDistributions(
    calculationID string,
    cooperativeID string,
    period string,
    memberShareAmount decimal.Decimal,
) ([]*models.SHUDistribution, error) {

    // 1. Get all members
    members, err := s.repo.GetActiveMembers(cooperativeID)
    if err != nil {
        return nil, err
    }

    // 2. Calculate total transactions and capital
    var totalTransactions decimal.Decimal
    var totalCapital decimal.Decimal

    memberData := make(map[string]struct {
        Transactions decimal.Decimal
        Capital      decimal.Decimal
    })

    for _, member := range members {
        // Get member transactions in period
        txns, err := s.repo.GetMemberTransactions(member.ID, period)
        if err != nil {
            return nil, err
        }

        // Get member capital
        capital, err := s.repo.GetMemberCapital(member.ID, period)
        if err != nil {
            return nil, err
        }

        memberData[member.ID] = struct {
            Transactions decimal.Decimal
            Capital      decimal.Decimal
        }{
            Transactions: txns,
            Capital:      capital,
        }

        totalTransactions = totalTransactions.Add(txns)
        totalCapital = totalCapital.Add(capital)
    }

    // 3. Calculate distributions
    distributions := make([]*models.SHUDistribution, 0)

    for _, member := range members {
        data := memberData[member.ID]

        // Transaction percentage
        var txnPct decimal.Decimal
        if totalTransactions.IsPositive() {
            txnPct = data.Transactions.Div(totalTransactions)
        }

        // Capital percentage
        var capPct decimal.Decimal
        if totalCapital.IsPositive() {
            capPct = data.Capital.Div(totalCapital)
        }

        // Contribution score (50% transactions + 50% capital)
        contributionScore := txnPct.Mul(decimal.NewFromFloat(0.5)).
                           Add(capPct.Mul(decimal.NewFromFloat(0.5)))

        // Distribution amount
        distAmount := memberShareAmount.Mul(contributionScore)

        dist := &models.SHUDistribution{
            CalculationID:        calculationID,
            MemberID:             member.ID,
            TotalTransactions:    data.Transactions,
            TransactionPercentage: txnPct.Mul(decimal.NewFromInt(100)),
            TotalCapital:         data.Capital,
            CapitalPercentage:    capPct.Mul(decimal.NewFromInt(100)),
            ContributionScore:    contributionScore.Mul(decimal.NewFromInt(100)),
            DistributionAmount:   distAmount,
            PaymentStatus:        "pending",
        }

        distributions = append(distributions, dist)
    }

    return distributions, nil
}
```

**Frontend Implementation:**
```tsx
// frontend/src/components/shu/SHUCalculationWizard.tsx
'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

const shuSchema = z.object({
  period: z.string().min(4),
  netIncome: z.number().positive(),
  reserve: z.number().min(25).max(100),
  memberShare: z.number().min(50).max(100),
  management: z.number().min(0).max(100),
  employee: z.number().min(0).max(100),
  education: z.number().min(0).max(100),
}).refine((data) => {
  const total = data.reserve + data.memberShare + data.management +
                data.employee + data.education;
  return total === 100;
}, {
  message: "Allocations must total 100%",
});

export function SHUCalculationWizard() {
  const [step, setStep] = useState(1);
  const [calculation, setCalculation] = useState(null);

  const { register, handleSubmit, watch, formState: { errors } } = useForm({
    resolver: zodResolver(shuSchema),
    defaultValues: {
      period: new Date().getFullYear().toString(),
      netIncome: 0,
      reserve: 25,
      memberShare: 50,
      management: 10,
      employee: 10,
      education: 5,
    },
  });

  const watchAllFields = watch();
  const allocationTotal = watchAllFields.reserve + watchAllFields.memberShare +
                          watchAllFields.management + watchAllFields.employee +
                          watchAllFields.education;

  const onSubmit = async (data) => {
    try {
      const response = await fetch('/api/shu/calculate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(data),
      });

      const result = await response.json();
      setCalculation(result);
      setStep(3); // Move to review step
    } catch (error) {
      console.error('SHU calculation failed:', error);
    }
  };

  return (
    <div className="max-w-4xl mx-auto p-6">
      <h1 className="text-2xl font-bold mb-6">Perhitungan SHU</h1>

      {/* Progress indicator */}
      <div className="mb-8">
        <div className="flex justify-between items-center">
          <Step number={1} active={step === 1} completed={step > 1}
                label="Data Dasar" />
          <Step number={2} active={step === 2} completed={step > 2}
                label="Alokasi" />
          <Step number={3} active={step === 3} completed={step > 3}
                label="Review" />
        </div>
      </div>

      <form onSubmit={handleSubmit(onSubmit)}>
        {step === 1 && (
          <div>
            <h2 className="text-xl font-semibold mb-4">Step 1: Data Dasar</h2>

            <div className="mb-4">
              <label className="block mb-2">Periode</label>
              <input
                type="text"
                {...register('period')}
                className="w-full p-2 border rounded"
                placeholder="2025"
              />
              {errors.period && (
                <p className="text-red-500 text-sm">{errors.period.message}</p>
              )}
            </div>

            <div className="mb-4">
              <label className="block mb-2">Sisa Hasil Usaha (SHU)</label>
              <input
                type="number"
                {...register('netIncome', { valueAsNumber: true })}
                className="w-full p-2 border rounded"
                placeholder="100000000"
              />
              {errors.netIncome && (
                <p className="text-red-500 text-sm">{errors.netIncome.message}</p>
              )}
              <p className="text-sm text-gray-600 mt-1">
                Format: {formatCurrency(watchAllFields.netIncome || 0)}
              </p>
            </div>

            <button
              type="button"
              onClick={() => setStep(2)}
              className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700"
            >
              Lanjut ke Alokasi
            </button>
          </div>
        )}

        {step === 2 && (
          <div>
            <h2 className="text-xl font-semibold mb-4">Step 2: Alokasi SHU</h2>

            <div className="bg-gray-100 p-4 rounded mb-6">
              <p className="font-semibold">Total SHU: {formatCurrency(watchAllFields.netIncome)}</p>
              <p className={allocationTotal === 100 ? 'text-green-600' : 'text-red-600'}>
                Total Alokasi: {allocationTotal}%
              </p>
            </div>

            <AllocationInput
              label="Dana Cadangan (Minimal 25%)"
              name="reserve"
              register={register}
              value={watchAllFields.reserve}
              amount={watchAllFields.netIncome * (watchAllFields.reserve / 100)}
              error={errors.reserve}
            />

            <AllocationInput
              label="Jasa Anggota (Minimal 50%)"
              name="memberShare"
              register={register}
              value={watchAllFields.memberShare}
              amount={watchAllFields.netIncome * (watchAllFields.memberShare / 100)}
              error={errors.memberShare}
            />

            <AllocationInput
              label="Jasa Pengurus"
              name="management"
              register={register}
              value={watchAllFields.management}
              amount={watchAllFields.netIncome * (watchAllFields.management / 100)}
              error={errors.management}
            />

            <AllocationInput
              label="Jasa Karyawan"
              name="employee"
              register={register}
              value={watchAllFields.employee}
              amount={watchAllFields.netIncome * (watchAllFields.employee / 100)}
              error={errors.employee}
            />

            <AllocationInput
              label="Dana Pendidikan"
              name="education"
              register={register}
              value={watchAllFields.education}
              amount={watchAllFields.netIncome * (watchAllFields.education / 100)}
              error={errors.education}
            />

            <div className="flex gap-4 mt-6">
              <button
                type="button"
                onClick={() => setStep(1)}
                className="bg-gray-300 text-gray-700 px-6 py-2 rounded hover:bg-gray-400"
              >
                Kembali
              </button>
              <button
                type="submit"
                disabled={allocationTotal !== 100}
                className="bg-blue-600 text-white px-6 py-2 rounded hover:bg-blue-700 disabled:bg-gray-400"
              >
                Hitung SHU
              </button>
            </div>
          </div>
        )}

        {step === 3 && calculation && (
          <div>
            <h2 className="text-xl font-semibold mb-4">Step 3: Review Hasil</h2>

            <div className="bg-green-50 p-6 rounded mb-6">
              <h3 className="font-semibold text-lg mb-4">Perhitungan Berhasil!</h3>

              <div className="grid grid-cols-2 gap-4">
                <div>
                  <p className="text-sm text-gray-600">Total SHU</p>
                  <p className="text-xl font-bold">{formatCurrency(calculation.netIncome)}</p>
                </div>
                <div>
                  <p className="text-sm text-gray-600">Jasa Anggota</p>
                  <p className="text-xl font-bold">{formatCurrency(calculation.memberShareAmount)}</p>
                </div>
              </div>
            </div>

            <MemberDistributionTable
              calculationId={calculation.id}
            />

            <div className="flex gap-4 mt-6">
              <button
                type="button"
                onClick={() => window.location.href = `/shu/${calculation.id}/report`}
                className="bg-green-600 text-white px-6 py-2 rounded hover:bg-green-700"
              >
                Lihat Laporan Lengkap
              </button>
              <button
                type="button"
                onClick={() => window.location.href = '/shu'}
                className="bg-gray-300 text-gray-700 px-6 py-2 rounded hover:bg-gray-400"
              >
                Kembali ke Dashboard
              </button>
            </div>
          </div>
        )}
      </form>
    </div>
  );
}

function AllocationInput({ label, name, register, value, amount, error }) {
  return (
    <div className="mb-4 p-4 border rounded">
      <label className="block mb-2 font-semibold">{label}</label>
      <div className="flex gap-4 items-center">
        <div className="flex-1">
          <input
            type="number"
            step="0.01"
            {...register(name, { valueAsNumber: true })}
            className="w-full p-2 border rounded"
          />
        </div>
        <div className="flex-shrink-0 w-32 text-right">
          <span className="text-2xl font-bold">{value}%</span>
        </div>
        <div className="flex-shrink-0 w-48 text-right">
          <span className="text-lg">{formatCurrency(amount)}</span>
        </div>
      </div>
      {error && (
        <p className="text-red-500 text-sm mt-1">{error.message}</p>
      )}
    </div>
  );
}

function formatCurrency(amount) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
}
```

**What This Achieves:**
- ✅ Saves 20+ hours of manual SHU calculation per cooperative
- ✅ Eliminates calculation errors
- ✅ Ensures compliance with UU No. 25/1992
- ✅ Generates professional reports for RAT (Annual Meeting)
- ✅ Transparent member distribution tracking
- ✅ Audit trail for all calculations

---

### Feature 2: Savings & Loan Module

**What It Does:**
Complete loan management system for savings cooperatives (Koperasi Simpan Pinjam) with multiple interest calculation methods.

**Interest Calculation Methods:**

**1. Flat Interest (Bunga Flat)**
```
Simple interest on principal
Monthly payment = (Principal + Total Interest) / Term

Example:
Principal: IDR 10,000,000
Interest: 12% per year
Term: 12 months

Total Interest = 10,000,000 × 12% = 1,200,000
Monthly Payment = (10,000,000 + 1,200,000) / 12 = 933,333
```

**2. Declining Balance (Bunga Efektif/Anuitas)**
```
Interest on remaining balance
Common for most consumer loans

Example:
Principal: IDR 10,000,000
Interest: 12% per year (1% per month)
Term: 12 months

Month 1:
  Interest = 10,000,000 × 1% = 100,000
  Principal = Monthly Payment - Interest

Month 2:
  Remaining = Principal - Principal Paid Month 1
  Interest = Remaining × 1%
  ...and so on
```

**3. Effective Interest**
```
True interest rate calculation
Used for compliance and comparison
```

**Database Schema:**
```sql
CREATE TABLE loans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cooperative_id UUID NOT NULL REFERENCES cooperatives(id),
    member_id UUID NOT NULL REFERENCES members(id),

    -- Loan details
    loan_number VARCHAR(50) UNIQUE NOT NULL,
    loan_date DATE NOT NULL,
    principal_amount DECIMAL(15,2) NOT NULL,
    interest_rate DECIMAL(5,2) NOT NULL,  -- Annual percentage
    interest_method VARCHAR(20) NOT NULL,  -- flat, declining, effective
    term_months INTEGER NOT NULL,

    -- Payment details
    monthly_payment DECIMAL(15,2) NOT NULL,
    total_interest DECIMAL(15,2) NOT NULL,
    total_repayment DECIMAL(15,2) NOT NULL,

    -- Status
    status VARCHAR(20) DEFAULT 'active',  -- active, paid_off, defaulted
    disbursement_date DATE,
    first_payment_date DATE,
    last_payment_date DATE,

    -- Tracking
    outstanding_balance DECIMAL(15,2) NOT NULL,
    total_paid DECIMAL(15,2) DEFAULT 0,
    payments_made INTEGER DEFAULT 0,

    -- Late payment
    late_payment_penalty_rate DECIMAL(5,2),  -- Percentage per day

    -- Collateral
    collateral_type VARCHAR(100),
    collateral_value DECIMAL(15,2),
    collateral_description TEXT,

    -- Metadata
    purpose TEXT,
    notes TEXT,
    approved_by UUID REFERENCES users(id),
    approved_at TIMESTAMP,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE loan_payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    loan_id UUID NOT NULL REFERENCES loans(id),

    -- Payment details
    payment_number INTEGER NOT NULL,
    payment_date DATE NOT NULL,
    due_date DATE NOT NULL,

    -- Amounts
    principal_amount DECIMAL(15,2) NOT NULL,
    interest_amount DECIMAL(15,2) NOT NULL,
    late_penalty DECIMAL(15,2) DEFAULT 0,
    total_amount DECIMAL(15,2) NOT NULL,

    -- Balance after payment
    balance_before DECIMAL(15,2) NOT NULL,
    balance_after DECIMAL(15,2) NOT NULL,

    -- Payment info
    payment_method VARCHAR(50),  -- cash, transfer, qris
    payment_reference VARCHAR(100),
    received_by UUID REFERENCES users(id),

    -- Status
    status VARCHAR(20) DEFAULT 'scheduled',  -- scheduled, paid, late, missed
    paid_date DATE,
    days_late INTEGER DEFAULT 0,

    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE savings_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cooperative_id UUID NOT NULL REFERENCES cooperatives(id),
    member_id UUID NOT NULL REFERENCES members(id),

    -- Account details
    account_number VARCHAR(50) UNIQUE NOT NULL,
    account_type VARCHAR(20) NOT NULL,  -- regular, time_deposit

    -- Balance
    balance DECIMAL(15,2) DEFAULT 0,

    -- Interest (for time deposits)
    interest_rate DECIMAL(5,2),
    interest_calculation_method VARCHAR(20),

    -- Status
    status VARCHAR(20) DEFAULT 'active',  -- active, frozen, closed
    opened_date DATE NOT NULL,
    closed_date DATE,

    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE savings_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID NOT NULL REFERENCES savings_accounts(id),

    -- Transaction details
    transaction_date DATE NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,  -- deposit, withdrawal, interest
    amount DECIMAL(15,2) NOT NULL,

    -- Balance
    balance_before DECIMAL(15,2) NOT NULL,
    balance_after DECIMAL(15,2) NOT NULL,

    -- Payment info
    payment_method VARCHAR(50),
    reference VARCHAR(100),
    processed_by UUID REFERENCES users(id),

    notes TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

**Go Implementation:**
```go
// backend/internal/services/interest_calculator.go
package services

import (
    "math"
    "github.com/shopspring/decimal"
)

type InterestCalculator struct{}

type LoanScheduleItem struct {
    PaymentNumber   int
    DueDate         time.Time
    PrincipalAmount decimal.Decimal
    InterestAmount  decimal.Decimal
    TotalAmount     decimal.Decimal
    BalanceBefore   decimal.Decimal
    BalanceAfter    decimal.Decimal
}

// Flat interest calculation
func (c *InterestCalculator) CalculateFlat(
    principal decimal.Decimal,
    annualRate decimal.Decimal,
    termMonths int,
) ([]LoanScheduleItem, error) {

    // Total interest = principal × annual rate
    totalInterest := principal.Mul(annualRate.Div(decimal.NewFromInt(100)))

    // Total repayment
    totalRepayment := principal.Add(totalInterest)

    // Monthly payment (fixed)
    monthlyPayment := totalRepayment.Div(decimal.NewFromInt(int64(termMonths)))

    // Monthly principal
    monthlyPrincipal := principal.Div(decimal.NewFromInt(int64(termMonths)))

    // Monthly interest
    monthlyInterest := totalInterest.Div(decimal.NewFromInt(int64(termMonths)))

    schedule := make([]LoanScheduleItem, termMonths)
    balance := principal

    for i := 0; i < termMonths; i++ {
        balanceBefore := balance
        balance = balance.Sub(monthlyPrincipal)

        schedule[i] = LoanScheduleItem{
            PaymentNumber:   i + 1,
            DueDate:         time.Now().AddDate(0, i+1, 0),
            PrincipalAmount: monthlyPrincipal,
            InterestAmount:  monthlyInterest,
            TotalAmount:     monthlyPayment,
            BalanceBefore:   balanceBefore,
            BalanceAfter:    balance,
        }
    }

    return schedule, nil
}

// Declining balance (effective/annuity) calculation
func (c *InterestCalculator) CalculateDeclining(
    principal decimal.Decimal,
    annualRate decimal.Decimal,
    termMonths int,
) ([]LoanScheduleItem, error) {

    // Monthly interest rate
    monthlyRate := annualRate.Div(decimal.NewFromInt(100)).
                   Div(decimal.NewFromInt(12))

    // Calculate monthly payment using annuity formula
    // PMT = P × (r × (1 + r)^n) / ((1 + r)^n - 1)

    // Convert to float64 for math.Pow
    r := monthlyRate.InexactFloat64()
    n := float64(termMonths)
    p := principal.InexactFloat64()

    // (1 + r)^n
    onePlusRPowN := math.Pow(1+r, n)

    // Monthly payment
    monthlyPaymentFloat := p * (r * onePlusRPowN) / (onePlusRPowN - 1)
    monthlyPayment := decimal.NewFromFloat(monthlyPaymentFloat)

    schedule := make([]LoanScheduleItem, termMonths)
    balance := principal

    for i := 0; i < termMonths; i++ {
        // Interest = balance × monthly rate
        interestAmount := balance.Mul(monthlyRate)

        // Principal = payment - interest
        principalAmount := monthlyPayment.Sub(interestAmount)

        balanceBefore := balance
        balance = balance.Sub(principalAmount)

        // Last payment adjustment (handle rounding)
        if i == termMonths-1 {
            principalAmount = principalAmount.Add(balance)
            balance = decimal.Zero
        }

        schedule[i] = LoanScheduleItem{
            PaymentNumber:   i + 1,
            DueDate:         time.Now().AddDate(0, i+1, 0),
            PrincipalAmount: principalAmount,
            InterestAmount:  interestAmount,
            TotalAmount:     monthlyPayment,
            BalanceBefore:   balanceBefore,
            BalanceAfter:    balance,
        }
    }

    return schedule, nil
}

// Calculate late payment penalty
func (c *InterestCalculator) CalculateLatePenalty(
    overdueAmount decimal.Decimal,
    penaltyRatePerDay decimal.Decimal,
    daysLate int,
) decimal.Decimal {

    // Penalty = amount × rate × days
    penalty := overdueAmount.
               Mul(penaltyRatePerDay.Div(decimal.NewFromInt(100))).
               Mul(decimal.NewFromInt(int64(daysLate)))

    return penalty
}
```

**What This Achieves:**
- ✅ Opens market to **savings cooperatives** (50% of total market)
- ✅ Automated loan calculations (no Excel errors)
- ✅ Multiple interest methods (flat, declining, effective)
- ✅ Automated payment schedules
- ✅ Late payment tracking and penalties
- ✅ Outstanding balance tracking
- ✅ Savings account management
- ✅ Complete audit trail

---

### Feature 3: QRIS Payment Integration

**What It Does:**
Enables digital payments through QRIS (Quick Response Code Indonesian Standard) using Xendit or Midtrans payment gateway.

**Payment Flow:**
```
1. Customer completes purchase → Staff selects QRIS payment
2. System calls Xendit/Midtrans API → Generate QR code
3. Display QR on screen/receipt → Customer scans with mobile banking
4. Customer confirms payment → Payment gateway processes
5. Webhook notification received → Update transaction status
6. Print receipt → Payment complete
```

**Provider Selection:**

| Feature | Xendit | Midtrans |
|---------|--------|----------|
| **QRIS Fee** | 0.7% | 0.7% |
| **Setup Fee** | Free | Free |
| **Settlement** | T+1 | T+1 |
| **API Quality** | Excellent | Good |
| **Documentation** | Excellent | Good |
| **Support** | 24/7 | Business hours |
| **Min Transaction** | IDR 1,500 | IDR 1,000 |
| **Max Transaction** | IDR 10,000,000 | IDR 10,000,000 |

**Recommendation:** **Xendit** (better API, documentation, and support)

**Go Implementation:**
```go
// backend/pkg/qris/xendit.go
package qris

import (
    "bytes"
    "encoding/json"
    "net/http"
    "github.com/shopspring/decimal"
)

type XenditClient struct {
    apiKey     string
    baseURL    string
    httpClient *http.Client
}

type CreateQRISRequest struct {
    ExternalID  string          `json:"external_id"`
    Amount      decimal.Decimal `json:"amount"`
    CallbackURL string          `json:"callback_virtual_account_id"`
}

type CreateQRISResponse struct {
    ID          string          `json:"id"`
    ExternalID  string          `json:"external_id"`
    Amount      decimal.Decimal `json:"amount"`
    QRString    string          `json:"qr_string"`
    Status      string          `json:"status"`
    ExpiryDate  string          `json:"expiry_date"`
}

func NewXenditClient(apiKey string) *XenditClient {
    return &XenditClient{
        apiKey:     apiKey,
        baseURL:    "https://api.xendit.co",
        httpClient: &http.Client{Timeout: 30 * time.Second},
    }
}

func (c *XenditClient) CreateQRIS(
    externalID string,
    amount decimal.Decimal,
    callbackURL string,
) (*CreateQRISResponse, error) {

    reqBody := CreateQRISRequest{
        ExternalID:  externalID,
        Amount:      amount,
        CallbackURL: callbackURL,
    }

    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest(
        "POST",
        c.baseURL+"/qr_codes",
        bytes.NewBuffer(jsonData),
    )
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.SetBasicAuth(c.apiKey, "")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var qrisResp CreateQRISResponse
    if err := json.NewDecoder(resp.Body).Decode(&qrisResp); err != nil {
        return nil, err
    }

    return &qrisResp, nil
}

func (c *XenditClient) GetQRISStatus(qrCodeID string) (*CreateQRISResponse, error) {
    req, err := http.NewRequest(
        "GET",
        c.baseURL+"/qr_codes/"+qrCodeID,
        nil,
    )
    if err != nil {
        return nil, err
    }

    req.SetBasicAuth(c.apiKey, "")

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var qrisResp CreateQRISResponse
    if err := json.NewDecoder(resp.Body).Decode(&qrisResp); err != nil {
        return nil, err
    }

    return &qrisResp, nil
}

// Webhook payload from Xendit
type XenditWebhook struct {
    ID             string          `json:"id"`
    ExternalID     string          `json:"external_id"`
    Amount         decimal.Decimal `json:"amount"`
    Status         string          `json:"status"`
    QRString       string          `json:"qr_string"`
    CallbackURL    string          `json:"callback_virtual_account_id"`
    TransactionID  string          `json:"transaction_id"`
    TransactionDate string         `json:"transaction_date"`
}

func (c *XenditClient) VerifyWebhook(signature, payload string) bool {
    // Implement webhook signature verification
    // Using Xendit's webhook verification method
    return true  // Simplified
}
```

**Frontend Implementation:**
```tsx
// frontend/src/app/(dashboard)/pos/components/QRISPayment.tsx
'use client';

import { useState, useEffect } from 'react';
import { QRCodeSVG } from 'qrcode.react';

interface QRISPaymentProps {
  amount: number;
  transactionId: string;
  onSuccess: () => void;
  onCancel: () => void;
}

export function QRISPayment({ amount, transactionId, onSuccess, onCancel }: QRISPaymentProps) {
  const [qrCode, setQrCode] = useState<string | null>(null);
  const [status, setStatus] = useState<'generating' | 'waiting' | 'success' | 'failed'>('generating');
  const [timeLeft, setTimeLeft] = useState(300); // 5 minutes

  useEffect(() => {
    generateQRCode();
  }, []);

  useEffect(() => {
    if (status === 'waiting') {
      // Poll payment status every 2 seconds
      const pollInterval = setInterval(checkPaymentStatus, 2000);

      // Countdown timer
      const timerInterval = setInterval(() => {
        setTimeLeft((prev) => {
          if (prev <= 1) {
            setStatus('failed');
            clearInterval(pollInterval);
            return 0;
          }
          return prev - 1;
        });
      }, 1000);

      return () => {
        clearInterval(pollInterval);
        clearInterval(timerInterval);
      };
    }
  }, [status]);

  const generateQRCode = async () => {
    try {
      const response = await fetch('/api/payments/qris/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          transaction_id: transactionId,
          amount: amount,
        }),
      });

      const data = await response.json();
      setQrCode(data.qr_string);
      setStatus('waiting');
    } catch (error) {
      console.error('Failed to generate QR code:', error);
      setStatus('failed');
    }
  };

  const checkPaymentStatus = async () => {
    try {
      const response = await fetch(`/api/payments/${transactionId}/status`);
      const data = await response.json();

      if (data.status === 'paid') {
        setStatus('success');
        setTimeout(onSuccess, 1500); // Show success for 1.5s then close
      }
    } catch (error) {
      console.error('Failed to check payment status:', error);
    }
  };

  const formatTime = (seconds: number) => {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg p-8 max-w-md w-full">
        <h2 className="text-2xl font-bold mb-4 text-center">Pembayaran QRIS</h2>

        {status === 'generating' && (
          <div className="text-center py-12">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
            <p>Membuat QR Code...</p>
          </div>
        )}

        {status === 'waiting' && qrCode && (
          <div>
            <div className="bg-gray-100 p-6 rounded-lg mb-4">
              <p className="text-center text-3xl font-bold text-blue-600 mb-2">
                {formatCurrency(amount)}
              </p>
              <p className="text-center text-sm text-gray-600">
                Waktu tersisa: <span className="font-semibold text-red-600">{formatTime(timeLeft)}</span>
              </p>
            </div>

            <div className="flex justify-center mb-4">
              <div className="bg-white p-4 rounded-lg shadow-lg">
                <QRCodeSVG
                  value={qrCode}
                  size={256}
                  level="H"
                  includeMargin={true}
                />
              </div>
            </div>

            <div className="bg-blue-50 p-4 rounded mb-4">
              <p className="text-sm text-center">
                📱 Scan QR code dengan aplikasi mobile banking atau e-wallet Anda
              </p>
            </div>

            <div className="flex gap-2 text-center text-xs text-gray-500 mb-4">
              <img src="/icons/gopay.png" alt="GoPay" className="h-8" />
              <img src="/icons/ovo.png" alt="OVO" className="h-8" />
              <img src="/icons/dana.png" alt="DANA" className="h-8" />
              <img src="/icons/shopeepay.png" alt="ShopeePay" className="h-8" />
              <img src="/icons/linkaja.png" alt="LinkAja" className="h-8" />
            </div>

            <div className="animate-pulse text-center text-sm text-gray-600 mb-4">
              ⏳ Menunggu pembayaran...
            </div>

            <button
              onClick={onCancel}
              className="w-full bg-gray-300 text-gray-700 py-2 rounded hover:bg-gray-400"
            >
              Batal
            </button>
          </div>
        )}

        {status === 'success' && (
          <div className="text-center py-12">
            <div className="text-green-600 text-6xl mb-4">✓</div>
            <p className="text-2xl font-bold text-green-600 mb-2">Pembayaran Berhasil!</p>
            <p className="text-gray-600">Transaksi telah berhasil diproses</p>
          </div>
        )}

        {status === 'failed' && (
          <div className="text-center py-12">
            <div className="text-red-600 text-6xl mb-4">✗</div>
            <p className="text-2xl font-bold text-red-600 mb-2">Pembayaran Gagal</p>
            <p className="text-gray-600 mb-6">QR Code telah kedaluwarsa atau pembayaran dibatalkan</p>
            <div className="flex gap-4">
              <button
                onClick={generateQRCode}
                className="flex-1 bg-blue-600 text-white py-2 rounded hover:bg-blue-700"
              >
                Coba Lagi
              </button>
              <button
                onClick={onCancel}
                className="flex-1 bg-gray-300 text-gray-700 py-2 rounded hover:bg-gray-400"
              >
                Batal
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}

function formatCurrency(amount: number) {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
  }).format(amount);
}
```

**What This Achieves:**
- ✅ **30-40% increase in transactions** (customers prefer digital payment)
- ✅ Reduced cash handling (safer, less errors)
- ✅ Automatic reconciliation (webhook updates)
- ✅ Real-time payment confirmation (< 2 seconds)
- ✅ Multiple e-wallet support (GoPay, OVO, DANA, ShopeePay, LinkAja)
- ✅ Lower transaction fees (0.7% vs 2-3% for credit cards)
- ✅ Modern customer experience

---

### Feature 4: Native Mobile App

**What It Does:**
React Native mobile app for Android and iOS enabling field staff to:
- Record POS transactions anywhere
- Scan barcodes
- Process loan payments
- View member information
- Work offline with sync

**Tech Stack:**
```yaml
Framework: React Native 0.76.5
Navigation: React Navigation 7.x
UI Library: React Native Paper 5.x
State Management: Zustand
Storage: AsyncStorage + SQLite (offline)
Camera: react-native-camera
Barcode: react-native-barcode-scanner
Printer: react-native-thermal-printer
```

**Key Features:**
1. **Offline-First Architecture**
   - Local SQLite database
   - Queue sync when online
   - Conflict resolution

2. **Barcode Scanning**
   - Camera-based scanning
   - Support EAN-13, Code-128
   - Product lookup

3. **Receipt Printing**
   - Bluetooth thermal printers
   - ESC/POS protocol
   - Custom templates

4. **Mobile POS**
   - Touch-optimized UI
   - Quick product selection
   - Multiple payment methods

**React Native Implementation:**
```tsx
// mobile/src/screens/pos/POSScreen.tsx
import React, { useState } from 'react';
import { View, FlatList, TouchableOpacity } from 'react-native';
import { Text, Button, Card, FAB } from 'react-native-paper';
import { useCart } from '../../hooks/useCart';
import { ProductGrid } from '../../components/pos/ProductGrid';
import { Cart } from '../../components/pos/Cart';
import { QRISPayment } from '../../components/payments/QRISPayment';

export function POSScreen({ navigation }) {
  const { items, total, addItem, removeItem, clear } = useCart();
  const [showQRIS, setShowQRIS] = useState(false);

  const handleBarcodeScan = () => {
    navigation.navigate('BarcodeScanner', {
      onScan: (product) => {
        addItem(product);
      },
    });
  };

  const handleCheckout = async (paymentMethod) => {
    if (paymentMethod === 'qris') {
      setShowQRIS(true);
    } else {
      // Handle cash payment
      await processCashPayment();
    }
  };

  const processCashPayment = async () => {
    try {
      const transaction = await api.createTransaction({
        items: items,
        total: total,
        payment_method: 'cash',
      });

      // Print receipt
      await printReceipt(transaction);

      // Clear cart
      clear();

      // Show success
      navigation.navigate('Success', { transaction });
    } catch (error) {
      console.error('Transaction failed:', error);
    }
  };

  return (
    <View style={styles.container}>
      <View style={styles.products}>
        <ProductGrid onSelectProduct={addItem} />
      </View>

      <View style={styles.cart}>
        <Cart
          items={items}
          total={total}
          onRemove={removeItem}
          onCheckout={handleCheckout}
        />
      </View>

      <FAB
        style={styles.fab}
        icon="barcode-scan"
        onPress={handleBarcodeScan}
        label="Scan Barcode"
      />

      {showQRIS && (
        <QRISPayment
          amount={total}
          onSuccess={() => {
            setShowQRIS(false);
            clear();
            navigation.navigate('Success');
          }}
          onCancel={() => setShowQRIS(false)}
        />
      )}
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    flexDirection: 'row',
  },
  products: {
    flex: 2,
    backgroundColor: '#f5f5f5',
  },
  cart: {
    flex: 1,
    backgroundColor: 'white',
    borderLeftWidth: 1,
    borderLeftColor: '#ddd',
  },
  fab: {
    position: 'absolute',
    bottom: 16,
    right: 16,
  },
});
```

```tsx
// mobile/src/screens/pos/BarcodeScanScreen.tsx
import React, { useState } from 'react';
import { View, Text, StyleSheet } from 'react-native';
import { RNCamera } from 'react-native-camera';
import { useProducts } from '../../hooks/useProducts';

export function BarcodeScanScreen({ route, navigation }) {
  const { onScan } = route.params;
  const { getProductByBarcode } = useProducts();
  const [scanning, setScanning] = useState(true);

  const handleBarCodeRead = async ({ data }) => {
    if (!scanning) return;

    setScanning(false);

    try {
      const product = await getProductByBarcode(data);

      if (product) {
        onScan(product);
        navigation.goBack();
      } else {
        alert('Produk tidak ditemukan');
        setScanning(true);
      }
    } catch (error) {
      console.error('Barcode scan error:', error);
      setScanning(true);
    }
  };

  return (
    <View style={styles.container}>
      <RNCamera
        style={styles.camera}
        type={RNCamera.Constants.Type.back}
        onBarCodeRead={handleBarCodeRead}
        captureAudio={false}
      >
        <View style={styles.overlay}>
          <View style={styles.scanArea}>
            <Text style={styles.instructions}>
              Arahkan kamera ke barcode
            </Text>
          </View>
        </View>
      </RNCamera>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  camera: {
    flex: 1,
  },
  overlay: {
    flex: 1,
    backgroundColor: 'rgba(0,0,0,0.5)',
    justifyContent: 'center',
    alignItems: 'center',
  },
  scanArea: {
    width: 300,
    height: 200,
    borderWidth: 2,
    borderColor: 'white',
    backgroundColor: 'transparent',
    justifyContent: 'center',
    alignItems: 'center',
  },
  instructions: {
    color: 'white',
    fontSize: 16,
    textAlign: 'center',
  },
});
```

**What This Achieves:**
- ✅ **Field operations enabled** (mobile cooperatives, door-to-door)
- ✅ **3x faster checkout** with barcode scanning
- ✅ **Works offline** (sync when connected)
- ✅ **Bluetooth receipt printing** (thermal printers)
- ✅ **100+ field staff** can use the system
- ✅ **Modern UX** optimized for mobile
- ✅ **Reduced hardware costs** (use smartphones vs POS terminals)

---

## 📊 Testing & QA

### Testing Strategy

**Unit Testing:**
```bash
# Backend (Go)
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out

# Target: 80%+ code coverage
```

**Integration Testing:**
```bash
# QRIS payment flow
go test ./tests/integration/qris_test.go -v

# Bank API integration
go test ./tests/integration/bank_test.go -v
```

**End-to-End Testing:**
```bash
# Frontend (Playwright)
cd frontend
npm run test:e2e

# Test scenarios:
# - Complete loan application
# - SHU calculation flow
# - QRIS payment in POS
# - Mobile app checkout
```

**Mobile Testing:**
```bash
# Android
cd mobile/android
./gradlew test

# iOS
cd mobile/ios
xcodebuild test

# Devices to test:
# - Android 8, 9, 10, 11, 12, 13
# - iOS 13, 14, 15, 16
```

**Load Testing:**
```bash
# Use k6 or Apache Bench
k6 run loadtest.js

# Scenarios:
# - 100 concurrent users
# - 1000 transactions/hour
# - 50 SHU calculations simultaneously
```

### Quality Gates

**Before deployment, ALL must pass:**
- [ ] Unit tests: 80%+ coverage
- [ ] Integration tests: All passing
- [ ] E2E tests: All critical paths passing
- [ ] Load tests: < 200ms p95 response time
- [ ] Mobile tests: Both platforms passing
- [ ] Security scan: No critical vulnerabilities
- [ ] Code review: Approved by 2+ developers
- [ ] Documentation: Updated
- [ ] Changelog: Written

---

## 🚀 Deployment Strategy

### Week 12: Production Deployment

**Step 1: Pre-Deployment (Monday)**
```bash
# 1. Backup production database
pg_dump production_db > backup_phase1_$(date +%Y%m%d).sql

# 2. Test migrations in staging
cd backend/migrations
./migrate-up.sh staging

# 3. Smoke test staging
./run-smoke-tests.sh staging

# 4. Freeze code (no more commits to main)
git tag phase-2-release-candidate
```

**Step 2: Deploy Backend (Tuesday Morning)**
```bash
# 1. Deploy to Cloud Run
gcloud run deploy koperasi-backend \
  --image gcr.io/project-id/backend:phase2 \
  --platform managed \
  --region asia-southeast2 \
  --allow-unauthenticated

# 2. Run migrations
./migrate-up.sh production

# 3. Verify health check
curl https://api.koperasi.id/health

# 4. Monitor logs
gcloud logging tail --project=project-id
```

**Step 3: Deploy Frontend (Tuesday Afternoon)**
```bash
# 1. Build production
cd frontend
npm run build

# 2. Deploy to Vercel/Cloud Run
vercel --prod

# 3. Verify deployment
curl https://app.koperasi.id

# 4. Test critical paths
npm run test:smoke-prod
```

**Step 4: Publish Mobile Apps (Wednesday)**
```bash
# Android
cd mobile/android
./gradlew bundleRelease
# Upload to Play Store Console

# iOS
cd mobile/ios
xcodebuild archive
# Upload to App Store Connect

# Wait for review (1-2 days for Android, 1-3 days for iOS)
```

**Step 5: Rollout to Users (Thursday-Friday)**
```
Thursday:
- 09:00: Announce Phase 2 to 10 MVP cooperatives
- 10:00: Upgrade first 3 cooperatives
- 14:00: Upgrade remaining 7 cooperatives
- 16:00: Training session for new features
- 17:00: Monitor usage and feedback

Friday:
- 09:00: Onboard 5 new cooperatives
- 13:00: Onboard 5 more new cooperatives
- Total: 20 cooperatives by end of day
```

**Rollback Plan:**
```bash
# If critical issues found:

# 1. Rollback backend
gcloud run deploy koperasi-backend \
  --image gcr.io/project-id/backend:phase1

# 2. Rollback database
psql production_db < backup_phase1_YYYYMMDD.sql

# 3. Rollback frontend
vercel rollback

# 4. Notify users
# Send email/WhatsApp notification
```

---

## 📈 Progress Tracking

### Weekly Metrics Dashboard

Track these every Friday:

**Development Progress:**
```
Week 1:  [ ] Team onboarded       [ ] Environment ready
Week 2:  [ ] Loan backend done    [ ] APIs tested
Week 3:  [ ] Loan frontend done   [ ] UI tested
Week 4:  [ ] SHU backend done     [ ] Calculation verified
Week 5:  [ ] SHU frontend done    [ ] Report generated
Week 6:  [ ] QRIS backend done    [ ] Payments working
Week 7:  [ ] QRIS POS done        [ ] QR codes scanning
Week 8:  [ ] Mobile app core      [ ] Barcode scanning
Week 9:  [ ] Mobile advanced      [ ] Printing working
Week 10: [ ] Bank API + POS HW    [ ] All integrated
Week 11: [ ] Testing complete     [ ] Bugs fixed
Week 12: [ ] Deployed to prod     [ ] Users migrated
```

**Feature Completion:**
```
SHU Calculation:        [================        ] 75%
Savings & Loans:        [========================] 100%
QRIS Payments:          [============            ] 50%
Bank Integration:       [======                  ] 25%
Mobile App:             [========                ] 33%
Advanced POS:           [====                    ] 20%
```

**Business Metrics:**
```
Cooperatives:           10 → 50  [====          ] 40/50
Monthly Users:          40 → 250 [=======       ] 150/250
Transactions/Month:     2K → 15K [=====         ] 8K/15K
MRR:                    10M → 125M [====        ] 50M/125M
NPS Score:              45 → 65+ [======        ] 58/65
```

**Team Velocity:**
```
Story Points/Week:      [Current: 45, Target: 60]
Code Reviews/Week:      [Current: 12, Target: 15]
Bugs Created:           [Current: 8, Target: < 5]
Bugs Fixed:             [Current: 10, Target: > 10]
```

### Daily Standup Template

```
What I did yesterday:
  -

What I'm doing today:
  -

Blockers:
  -

Help needed:
  -
```

### Weekly Review Template

```
Week X Review (Month Y)

Completed:
  -
  -

In Progress:
  -
  -

Blocked:
  -
  -

Metrics:
  - Feature completion: X%
  - Tests written: X
  - Bugs fixed: X
  - Code coverage: X%

Next Week Goals:
  -
  -

Risks:
  -
  -
```

---

## ✅ Phase 2 Success Checklist

### Technical Success

**Backend:**
- [ ] All new database tables created
- [ ] 6 new feature modules implemented
- [ ] 50+ new API endpoints
- [ ] Unit tests: 80%+ coverage
- [ ] Integration tests passing
- [ ] API response time < 200ms (p95)
- [ ] Zero memory leaks
- [ ] Background workers stable

**Frontend:**
- [ ] 20+ new pages created
- [ ] 50+ new components
- [ ] Mobile-responsive design
- [ ] Accessibility compliant
- [ ] Bundle size optimized
- [ ] Page load < 2s
- [ ] No console errors

**Mobile:**
- [ ] Android app published
- [ ] iOS app published
- [ ] Offline mode working
- [ ] Barcode scanning accurate
- [ ] Receipt printing working
- [ ] App size < 50MB
- [ ] Startup time < 3s

**Infrastructure:**
- [ ] CI/CD pipeline for mobile
- [ ] Monitoring for new services
- [ ] Alerts configured
- [ ] Backups automated
- [ ] Disaster recovery tested
- [ ] Security audit passed

### Business Success

**User Adoption:**
- [ ] 50 cooperatives onboarded
- [ ] 250+ active users
- [ ] 15,000+ monthly transactions
- [ ] 60%+ QRIS adoption
- [ ] 40%+ using loan module
- [ ] 30%+ completed SHU calculation

**Financial:**
- [ ] MRR: IDR 125M
- [ ] Churn rate < 5%
- [ ] CAC < IDR 2.5M
- [ ] LTV > IDR 250M
- [ ] Gross margin > 75%

**Quality:**
- [ ] NPS score 65+
- [ ] Support tickets < 50/week
- [ ] Response time < 2 hours
- [ ] Resolution time < 24 hours
- [ ] Uptime 95%+
- [ ] Zero data loss incidents

### Team Success

**People:**
- [ ] 12 team members productive
- [ ] Team satisfaction 8/10+
- [ ] Zero regretted attrition
- [ ] Skills matrix updated
- [ ] Career paths defined

**Process:**
- [ ] Agile ceremonies running
- [ ] Code review < 24 hours
- [ ] Documentation updated
- [ ] Knowledge sharing weekly
- [ ] Retrospectives actionable

---

## 🎉 What Phase 2 Achieves

### For Cooperatives

**Before Phase 2:**
- ❌ Manual SHU calculation (20+ hours)
- ❌ Excel for loan tracking (errors common)
- ❌ Cash-only transactions
- ❌ No mobile access
- ❌ Manual receipt writing

**After Phase 2:**
- ✅ Automated SHU (5 minutes)
- ✅ Digital loan management (zero errors)
- ✅ QRIS payments (modern)
- ✅ Mobile app (anywhere access)
- ✅ Printed receipts (professional)

### For Members

**Before:**
- ❌ Wait for manual calculations
- ❌ Unclear loan balances
- ❌ Must pay with cash
- ❌ Paper receipts lost
- ❌ No transparency

**After:**
- ✅ Real-time SHU distribution
- ✅ View loan details anytime
- ✅ Pay with e-wallet
- ✅ Digital receipts
- ✅ Full transparency

### For Your Business

**Scale:**
- 5x cooperative growth (10 → 50)
- 6x user growth (40 → 250)
- 7.5x transaction growth
- 12.5x revenue growth

**Market Position:**
- ✅ Only platform with SHU automation
- ✅ Only one with full loan module
- ✅ Only one with QRIS for cooperatives
- ✅ Only one with mobile app
- ✅ **Market leader position secured**

---

## 📅 Next Steps After Phase 2

**Month 7: Phase 3 Kickoff**
- Inventory management
- Multi-unit support
- WhatsApp integration
- Business intelligence

**Review this document weekly during Phase 2. Track progress ruthlessly. Ship features incrementally. Celebrate wins. Fix issues fast. Stay focused.**

**Let's build Indonesia's dominant cooperative platform! 🚀**

---

**Document Version:** 1.0
**Last Updated:** November 16, 2025
**Owner:** Product Team
**Status:** Active Development Guide
