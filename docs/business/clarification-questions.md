# Clarification Questions for App Development

## Purpose
Before translating the business requirements into an actual application, we need to clarify specific details, make technology decisions, and resolve ambiguities. This document organizes all questions that need answers before development begins.

---

## CRITICAL DECISIONS NEEDED

### 1. Technology Stack Selection

**Backend Framework:**
- [ ] **Option A**: Node.js + Express/NestJS + TypeScript
  - Pros: Fast development, large ecosystem, JavaScript everywhere
  - Cons: Less mature for enterprise

- [ ] **Option B**: PHP + Laravel
  - Pros: Mature, great documentation, popular in Indonesia
  - Cons: Slower than Node.js for real-time features

- [ ] **Option C**: Python + Django/FastAPI
  - Pros: Great for data/analytics, clean code
  - Cons: Smaller talent pool in Indonesia

- [x] **Option D**: Go (Golang) + Gin/Echo/Fiber
  - Pros: Excellent performance, built-in concurrency, low memory footprint, single binary deployment
  - Cons: Smaller web framework ecosystem than Node.js, more verbose than Python/Node

**Question**: Which backend technology do you prefer? Consider your team's expertise and local talent availability.

**Answer**: ✅ **Go (Golang)**
- Excellent choice for high-performance, scalable systems
- Perfect for handling 4,000+ cooperatives with concurrent transactions
- Low infrastructure costs (efficient resource usage)
- Strong typing prevents bugs
- Great for microservices architecture later

---

**Frontend Framework:**
- [ ] **Option A**: React + TypeScript
  - Pros: Industry standard, huge ecosystem, good for complex UIs
  - Cons: Steeper learning curve

- [ ] **Option B**: Vue.js
  - Pros: Easier to learn, good documentation in Indonesian
  - Cons: Smaller ecosystem than React

- [ ] **Option C**: Next.js (React framework)
  - Pros: Built-in SSR, better SEO, modern
  - Cons: More complex setup

**Question**: Which frontend framework? Do you have existing team expertise in any of these?

answer: Next.js (React framework)
  - Pros: Built-in SSR, better SEO, modern
  - Cons: More complex setup

---

**Mobile App Framework:**
- [x] **Option A**: React Native
  - Pros: Share code with web (if using React), single codebase for iOS/Android
  - Cons: Performance limitations for complex UIs

- [ ] **Option B**: Flutter
  - Pros: Better performance, beautiful UI, growing in Indonesia
  - Cons: Different language (Dart), can't share with web

**Question**: React Native (JavaScript) or Flutter (Dart)? Consider team skills and performance needs.

---

**Database:**
- [x] **PostgreSQL** (Recommended in docs)
- [ ] **MySQL** (More common in Indonesia)
- [ ] **MongoDB** (NoSQL, flexible schema)

**Question**: PostgreSQL as recommended, or do you prefer MySQL for better local support?

---

### 2. Deployment & Hosting Strategy

**Cloud Provider:**
- [ ] **AWS** (Amazon Web Services) - Most mature, expensive
- [x] **Google Cloud Platform** - Good for Indonesia, competitive pricing
- [ ] **Microsoft Azure** - Good government relationships
- [ ] **Local Indonesian Cloud** (Telkom Cloud, Biznet Gio) - Data sovereignty
- [ ] **Hybrid** - On-premise option for enterprise customers

**Questions:**
1. Which cloud provider? Consider cost, local presence, and government requirements
2. Do we need on-premise deployment option for government/large cooperatives?
3. What about data sovereignty requirements? (Indonesian data must stay in Indonesia?)

answer: Cloud Provider: Start with Google Cloud Platform
Has Jakarta region (data sovereignty ✓)
Cheaper than AWS
Good Indonesian documentation
Easy migration path later

---

**Infrastructure Approach:**
- [ ] **Fully Managed** (Heroku, Vercel, Railway) - Easiest, more expensive
- [ ] **Container-based** (Docker + Kubernetes) - Flexible, complex
- [ ] **Serverless** (AWS Lambda, Cloud Functions) - Scalable, pay-per-use
- [ ] **Traditional VPS** (DigitalOcean, Linode) - Cheapest, manual scaling

**Question**: Start simple (managed) then migrate to containers, or build for scale from day one?
answer: Start Managed, Plan for Containers
Begin with Cloud Run (managed containers)
No Kubernetes complexity initially
Plan Docker architecture from day 1
Scale to K8s when you hit 100+ cooperatives
On-premise: Not for MVP, but prepare architecture for it (some government cooperatives will require it)

---

### 3. MVP Scope Definition

**CRITICAL QUESTION**: What is the Minimum Viable Product (MVP) for pilot program?

**Proposed MVP Scope (3 months):**
- [ ] Member Management (registration, KTP validation, basic profile)
- [ ] Share Capital Tracking (Simpanan Pokok/Wajib/Sukarela)
- [ ] Basic Accounting (COA, Journal Entry, Trial Balance)
- [ ] Simple POS (1 business unit, cash only)
- [ ] Basic Reporting (Member list, Transaction history, Simple financial reports)
- [ ] Mobile App (Member view balances, transaction history)
- [ ] User Management (Login, roles, permissions)

**Defer to Post-MVP:**
- [ ] SHU Calculation (can be manual for pilot)
- [ ] Multiple Business Units
- [ ] Inventory Management
- [ ] Savings & Loans Module
- [ ] QRIS Payment Integration
- [ ] WhatsApp Integration
- [ ] Advanced Reporting

**Question**: Does this MVP scope make sense for the 10 pilot cooperatives? What MUST be in MVP vs can wait?

answer: MUST HAVE (Month 1-2):

Member Management (basic, no KTP validation yet)
Share Capital Tracking
Simple POS (cash only)
Basic financial reports (for RAT meeting)
Web app only (no mobile yet)

NICE TO HAVE (Month 3):

Basic mobile view (responsive web, not app)
Simple accounting entries
Transaction history

DEFINITELY DEFER:

SHU Calculation (critical but complex - do it right post-MVP)
Mobile app
QRIS payments
WhatsApp
Inventory management

---

### 4. Integration Priorities

**NIK Validation (Dukcapil API):**
- **Question**: Do we have access to Dukcapil API? Is there a cost?
- **Alternative**: Manual KTP entry for MVP, add validation later?
- **Question**: Is NIK validation mandatory or nice-to-have for pilot?

---

**QRIS Payment Integration:**
- **Question**: Which payment gateway to start with?
  - [ ] Midtrans (popular, good docs)
  - [ ] Xendit (growing fast)
  - [ ] OVO/GoPay/Dana (direct integration)
  - [ ] Multiple (more complex)

- **Question**: Is QRIS needed for MVP or can we start with cash-only POS?

---

**WhatsApp Business API:**
- **Question**: Do we have WhatsApp Business API access? (Requires Meta approval)
- **Cost**: ~IDR 100-500K/month depending on volume
- **Alternative**: Start with SMS or email, add WhatsApp later?
- **Question**: Mandatory for MVP or can defer?

---

**Bank Integration:**
- **Question**: Which banks to integrate first for reconciliation?
- **Question**: Is this critical for pilot or can cooperatives manually reconcile initially?

answer: NIK Validation:

Skip for MVP - manual entry is fine
Dukcapil API requires government partnership (6+ months process)
Add fake validation UI for demo purposes

QRIS Payment:

Skip for MVP - cash only
When ready: Use Xendit (best documentation, startup-friendly)
Add in month 4-6

WhatsApp Business API:

Skip for MVP - not critical
Use regular SMS for urgent notifications
WhatsApp in Phase 2 (month 6+)

Bank Integration:

Skip for MVP
Manual reconciliation is normal for cooperatives
Add BNI/Mandiri API in Phase 2

---

### 5. Indonesian Localization Requirements

**Language:**
- [ ] **Bahasa Indonesia only** (simpler, faster development)
- [ ] **Bilingual** (Indonesian + English for government/reports)
- [ ] **Multi-language** (support regional languages: Javanese, Sundanese, etc.)

**Question**: Indonesian only for MVP, or is English needed for government reports?

---

**Date/Number Formatting:**
- Date format: DD/MM/YYYY or DD-MM-YYYY?
- Number format: 1.000.000,00 (Indonesian) or 1,000,000.00 (International)?
- Currency: Always "Rp" or "IDR"?

**Question**: Confirm Indonesian formatting preferences

---

**Fiscal Year:**
- Does Indonesian cooperative fiscal year = calendar year (Jan-Dec)?
- Or can cooperatives set custom fiscal year?

**Question**: Clarify fiscal year handling

answer: Language:

Bahasa Indonesia ONLY for MVP
English can wait (even government accepts Indonesian reports)
Regional languages are nice-to-have for Phase 3

Formatting:

Date: DD/MM/YYYY (Indonesian standard)
Number: 1.000.000,00 (Indonesian format)
Currency: "Rp" for display, store as integer (satuan)
Example: Rp 1.500.000,00

Fiscal Year:

Default: Calendar year (Jan-Dec)
Allow custom configuration (some cooperatives use July-June)

---

### 6. Regulatory & Compliance Questions

**Ministry of Cooperatives Reporting:**
- **Question**: What EXACTLY are the required government reports?
  - List of specific reports needed
  - Format (PDF, Excel, online submission?)
  - Frequency (monthly, quarterly, annually?)
  - Where to submit (online portal, email, physical?)

**Critical**: We need actual report templates/samples from Ministry

---

**RAT (Rapat Anggota Tahunan) Requirements:**
- **Question**: What documents must be prepared for RAT?
  - Financial statements (which ones exactly?)
  - SHU calculation report format?
  - Member list format?
  - Other required documents?

**Critical**: Need sample RAT document package

---

**Accounting Standards:**
- **Question**: Must follow SAK ETAP (Indonesian accounting standard)?
- Is there specific Chart of Accounts template for cooperatives?
- Are there mandatory account codes?

**Critical**: Need official cooperative COA template if exists

---

**SHU (Sisa Hasil Usaha) Calculation:**
- **Question**: What is the EXACT formula for SHU calculation?
- How to split between:
  - Member services (jasanggota)
  - Member capital (modalnggota)
  - Reserves
  - Other allocations
- Is there a government regulation defining this?

**Critical**: Need official SHU calculation methodology

---

**Data Privacy & Security:**
- **Question**: Are there Indonesian data privacy regulations we must comply with?
- KTP/NIK data - encryption requirements?
- Member financial data - storage requirements?
- Data breach notification requirements?

answer: Ministry Reports:

Monthly: Simple transaction summary
Quarterly: Financial position snapshot
Annual: RAT package - comprehensive reports
Format: Start with PDF/Excel, online submission coming in 2026

RAT Documents:

Neraca (Balance Sheet)
Laporan Laba Rugi (Income Statement)
Laporan Perubahan Modal (Changes in Equity)
Laporan SHU dan Pembagiannya
Daftar Hadir Anggota
Action: Get samples from existing cooperatives in pilot area

Accounting Standards:

Use SAK ETAP (simpler than full SAK)
No mandatory COA but follow common practice
Action: Use standard COA from successful cooperatives

SHU Formula (typical):
Total SHU = Net Income after tax
- Cadangan (Reserve): 25%
- Jasa Anggota (Member services): 25%
- Jasa Modal (Capital return): 20%
- Pengurus (Management): 10%
- Karyawan (Employees): 10%
- Sosial (Social): 5%
- Pembangunan Daerah (Regional): 5%
Note: Each cooperative can adjust percentages in their AD/ART
Data Privacy:

Follow UU ITE basics
Encrypt NIK/KTP data
No specific cooperative data law yet

---

### 7. User Roles & Permissions

**Proposed Roles:**
1. **Super Admin** (Our team - system management)
2. **Cooperative Admin** (Cooperative IT person)
3. **Chairman/Board** (Read-only, dashboards)
4. **Manager** (Operational management, approvals)
5. **Finance Officer** (Accounting, reporting)
6. **Cashier** (POS, transactions)
7. **Warehouse Staff** (Inventory - post-MVP)
8. **Member** (Mobile app - view only)

**Questions:**
1. Are these roles sufficient or missing any? yes
2. Can one person have multiple roles? yes
3. Should roles be customizable per cooperative? yes
4. Any role-specific features needed? yes

---

### 8. Business Rules & Workflows

**Member Registration:**
- Can members self-register or must admin register?
- Approval workflow needed?
- Minimum age requirement?
- Can one person be member of multiple cooperatives?
- Family membership allowed?

---

**Share Capital Rules:**
- **Simpanan Pokok** (Principal Share):
  - Paid once when joining?
  - Refundable when leaving?
  - Can amount vary per member or fixed?

- **Simpanan Wajib** (Mandatory Share):
  - Frequency (monthly, quarterly, annually)?
  - Amount fixed or variable?
  - Penalty for late payment?

- **Simpanan Sukarela** (Voluntary Share):
  - Any time contribution?
  - Minimum amount?
  - Withdrawable?

**Question**: Confirm the exact rules for each share capital type

---

**Transaction Approval Workflows:**
- Which transactions need approval?
- Approval thresholds (e.g., > IDR 10M needs manager approval)?
- Multi-level approval needed?
- Can approver edit or only approve/reject?

---

**Business Unit Management:**
- Can one product/service be in multiple business units?
- How to handle inter-unit transfers?
- Shared inventory or separate per unit?
- Consolidated P&L required?

answer 
Member Registration:
Admin registers
No approval workflow initially
Minimum age: 17 (has KTP)
One cooperative per person for pilot

Share Capital Rules:
Simpanan Pokok: One-time, IDR 50K-100K typical
Simpanan Wajib: Monthly, IDR 10K-ultimateK typical
Simpanan Sukarela: Anytime, any amount
MVP: Track amounts only, no automation
---

### 9. POS System Specifics

**POS Hardware:**
- Will cooperatives use:
  - [ ] Desktop/laptop with browser
  - [x] Tablet (Android/iPad)
  - [ ] Dedicated POS terminal
  - [x] Smartphone only

**Question**: What devices will be used? Affects UI design. andorid

---

**POS Features Priority:**
- [ ] Barcode scanning (need camera or scanner?)
- [ ] Receipt printing (thermal printer support?)
- [ ] Cash drawer integration
- [ ] Weight scale integration (for bulk items)
- [ ] Multiple payment methods (cash + QRIS + credit)
- [ ] Discount management
- [ ] Shift management
- [ ] Member pricing (different price for members?)

**Question**: Which POS features are essential for pilot? all of that

---

**Offline Mode:**
- **Question**: How important is offline POS? no
- Should it sync when back online? yes
- How to handle conflicts (e.g., sold same item twice)? give option to user to choose. mark this for leter

---

### 10. Inventory Management Questions

**Stock Valuation Method:**
- [x] FIFO (First In First Out)
- [ ] LIFO (Last In First Out)
- [x] Average Cost
- [x] Allow cooperative to choose?

**Question**: Which method to implement? Indonesian accounting standard preference? follow Indonesian accounting standard

---

**Product Management:**
- Product variants (size, color, etc.) needed? yrs
- Batch/lot tracking needed? yes
- Serial number tracking needed? yes
- Expiry date tracking needed? yes
- Supplier linking (one product from multiple suppliers)? no

---

**Stock Opname (Cycle Counting):**
- Mobile app for stock counting? no
- Variance approval workflow? no
- Adjustment posting rules?

---

### 11. Mobile App Features

**Member Mobile App - Priority Features:**

**Must Have:**
- [ ] Login / authentication
- [ ] View share capital balance
- [ ] View transaction history
- [ ] View cooperative announcements
- [ ] Contact cooperative

**Should Have:**
- [ ] Digital member card (QR code)
- [ ] Apply for loan (if loans module exists)
- [ ] View savings balance
- [ ] Push notifications

**Nice to Have:**
- [ ] Buy products (e-commerce)
- [ ] Pay bills
- [ ] Transfer to other members
- [ ] Vote in meetings

**Question**: What should be in MVP mobile app vs later versions? Responsive web only, no native app

---

**Staff Mobile App:**
- Separate app or same app with role-based views? same app with role based views
- What functions do staff need on mobile?
  - Stock counting? n
  - Approve transactions? y
  - View reports? y

---

### 12. Reporting Requirements

**Financial Reports - Which are mandatory?**
- [x] Balance Sheet (Neraca)
- [x] Income Statement (Laba Rugi)
- [x] Cash Flow Statement
- [x] Changes in Equity
- [x] Trial Balance
- [x] General Ledger
- [x] Journal Entry Report

**Question**: Which reports MUST be in MVP? all of that

---

**Operational Reports:**
- [x] Sales Report (daily, monthly)
- [x] Purchase Report
- [x] Inventory Report
- [x] Member Transaction Report
- [x] SHU Report
- [ ] Attendance Report (if HR module)

---

**Government Reports:**
**Question**: Please provide list of actual government report requirements with templates

---

**Report Export Formats:**
- [x] PDF
- [x] Excel
- [ ] CSV
- [x] Print directly

**Question**: Which formats are essential?

---

### 13. Pricing & Subscription Management

**Billing Questions:**
- **Question**: Who handles payment processing for subscriptions?
  - Direct bank transfer?
  - Payment gateway (credit card, e-wallet)?
  - Invoice + manual verification?

---

**Free Tier Limitations:**
- How to enforce 50 member limit?
- What happens when cooperative exceeds limit?
  - Block new members?
  - Force upgrade?
  - Grace period?

---

**Trial Period:**
- Free trial duration (7, 14, 30 days)?
- Credit card required upfront?
- Auto-convert to paid or manual?

---

**Subscription Changes:**
- Can cooperative downgrade mid-contract?
- Prorated refunds?
- What happens to data when downgrading?

---
answer: MVP Approach:

  Manual billing (send invoice)
  Bank transfer payment
  Manual activation
  No enforcement (trust-based for pilot)

Defer: Automated billing to Phase 2

### 14. Data Migration Strategy

**From Manual/Excel:**
- **Question**: What data needs migration for pilot cooperatives?
  - Member list (what fields?)
  - Opening balances
  - Historical transactions (how far back?)
  - Inventory (products and current stock)

---

**Migration Approach:**
- [ ] **Template-based**: We provide Excel template, they fill, we import
- [ ] **Manual entry**: Our team enters the data
- [ ] **API import**: For larger cooperatives
- [ ] **Gradual**: Start fresh, migrate historical as needed

**Question**: What's feasible for 10 pilot cooperatives in 3 months?

---
answer: 
Excel template we provide
They fill: Members, Opening Balances
We import manually (one-time)
Start fresh for transactions

Data needed:
Member list (Name, NIK, Phone, Join Date)
Share capitals (Pokok, Wajib, Sukarela)
Opening cash balance

### 15. Support & Training

**Support Channels:**
- [x] Email
- [x] WhatsApp (business account)
- [x] Phone
- [ ] Live chat
- [ ] In-app chat
- [ ] Help desk system
- [ ] whatsapp

**Question**: Which support channels for pilot? Affects tooling decisions.

---

**Training Approach:**
- [x] On-site training at cooperative
- [x] Virtual training (Zoom/Google Meet)
- [ ] Video tutorials (pre-recorded)
- [x] Written documentation
- [x] In-app guided tours
- [ ] Train-the-trainer

**Question**: What training format for pilot program?

---

**Documentation:**
- **Question**: What languages for documentation?
  - Bahasa Indonesia only? indonesia only
  - English for technical docs?

- **Question**: What format?
  - Video tutorials
  - PDF manuals
  - Online help center
  - In-app contextual help

---

### 16. Security & Authentication

**Authentication Method:**
- [ ] Email + Password
- [x] Phone + OTP
- [ ] Email + OTP
- [ ] Social login (Google, Facebook)
- [x] Username + Password

**Question**: What's most accessible for rural cooperative users?

---

**Password Requirements:**
- Minimum length?
- Complexity requirements (uppercase, numbers, symbols)?
- Expiry (force reset every X days)?
- Remember previous passwords?

**Question**: Balance security vs usability for non-tech users

---

**Two-Factor Authentication:**
- Mandatory for which roles (all, admin only, finance only)?
- OTP via SMS or email or authenticator app?

---

**Session Management:**
- Session timeout duration (15 min, 30 min, 2 hours)?
- Remember me option?
- Single device login or multiple?

---
answer: 
No password hassle
30-minute session timeout
Single device login

### 17. Multi-tenancy Architecture

**Data Isolation:**
- [ ] **Separate Database per Cooperative** (most isolated, expensive)
- [ ] **Shared Database, Separate Schema per Cooperative** (balanced)
- [ ] **Shared Database, Shared Schema** (cheapest, least isolated)

**Question**: How important is data isolation? Consider security vs cost.

---

**Customization Level:**
- Can each cooperative customize:
  - Chart of Accounts?
  - Report formats?
  - Workflows?
  - Branding (logo, colors)?
  - Email templates?

**Question**: How much customization to allow?

---
answer:
MVP: Shared Database, Shared Schema
  Add cooperative_id to all tables
  Simple but sufficient for pilot
  Easy to migrate later if needed


### 18. Performance & Scale Targets

**Pilot Scale (10 cooperatives):**
- Average 300 members each = 3,000 total members
- Estimate 100 transactions/day per cooperative = 1,000 transactions/day
- Estimate 10 concurrent users per cooperative = 100 concurrent users

---

**Year 1 Scale (4,000 cooperatives):**
- 1.2 million members
- 400,000 transactions/day
- 40,000 concurrent users

**Questions:**
1. Do these estimates seem reasonable?
2. Should we build for Year 1 scale from day 1, or start simple?
3. Performance targets: Page load < 2 seconds acceptable?

---
answer:
  10 cooperatives
  3,000 total members
  1,000 transactions/day
  Page load < 3 seconds on 3G

### 19. Testing Strategy

**Testing Approach:**
- [ ] Manual testing only (faster, less rigorous)
- [ ] Automated unit tests (slower development, more reliable)
- [ ] Integration tests
- [ ] End-to-end tests
- [ ] User acceptance testing (UAT)

**Question**: Testing rigor for MVP? Balance speed vs quality.

---

**Test Data:**
- Need realistic test data for demos
- How many test cooperatives?
- How much test data (members, transactions)?

---

### 20. Legal & Compliance

**Terms of Service:**
- Who owns the cooperative data?
- What happens if cooperative stops paying?
- Data deletion policy?
- Service level agreement (SLA)?

**Question**: Need lawyer review before pilot?

---

**Data Retention:**
- How long to keep data after cooperative leaves?
- Backup retention policy?
- Audit log retention?

---

**Insurance:**
- Cyber insurance needed?
- Professional liability insurance?
- Data breach insurance?

---
answer:
  Simple Terms of Service
  Data owned by cooperative
  30-day data retention after stop paying
  No SLA for pilot


## PRIORITIZATION MATRIX

### MUST ANSWER BEFORE DEVELOPMENT STARTS (Blocking)

1. **Technology stack** (Backend, Frontend, Mobile, Database)
2. **MVP scope** (What features in pilot vs later)
3. **SHU calculation formula** (Core business logic)
4. **Government reporting requirements** (Compliance)
5. **Chart of Accounts structure** (Accounting foundation)
6. **Share capital rules** (Business rules)
7. **Deployment strategy** (Infrastructure)

### SHOULD ANSWER DURING PILOT (Important but not blocking)

8. Integration priorities (QRIS, WhatsApp, NIK validation)
9. Offline mode requirements
10. Mobile app feature priorities
11. Data migration approach
12. Support channels
13. User roles refinement

### NICE TO ANSWER (Can decide later)

14. White-label customization
15. Advanced reporting
16. Inter-cooperative marketplace
17. API for third-party integrations
18. Multi-language support

---

## NEXT STEPS

### For Business Team:
1. **Review and answer** all questions in "MUST ANSWER" section
2. **Gather** actual government report templates and requirements
3. **Interview** 2-3 cooperative managers to validate assumptions
4. **Provide** sample data for realistic testing

### For Technical Team:
5. **Technology decision meeting** - finalize stack
6. **Architecture design** based on answers
7. **Database schema** design v1.0
8. **Development environment** setup

### For Both:
9. **Workshop** - walk through all questions together
10. **Document decisions** in separate file
11. **Create** detailed product requirements document (PRD)
12. **Begin** sprint planning for MVP

---

## QUESTIONS TRACKER

| Question # | Category | Priority | Status | Answered By | Answer |
|-----------|----------|----------|--------|-------------|--------|
| 1.1 | Tech Stack - Backend | MUST | ❌ Pending | | |
| 1.2 | Tech Stack - Frontend | MUST | ❌ Pending | | |
| 1.3 | Tech Stack - Mobile | MUST | ❌ Pending | | |
| 2.1 | Cloud Provider | MUST | ❌ Pending | | |
| 3.1 | MVP Scope | MUST | ❌ Pending | | |
| 6.1 | Government Reports | MUST | ❌ Pending | | |
| 6.3 | SHU Calculation | MUST | ❌ Pending | | |
| ... | ... | ... | ... | | |

---

**Version**: 1.0.0
**Last Updated**: 2025-11-15
**Document Owner**: Product Management

**Instructions**:
1. Review all questions
2. Schedule clarification meeting
3. Document all answers
4. Update tracker as questions are resolved
5. Don't start development until "MUST ANSWER" questions resolved
