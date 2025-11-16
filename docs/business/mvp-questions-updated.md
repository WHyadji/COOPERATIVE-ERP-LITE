# MVP Questions - Updated Based on Answers

## Status Summary

✅ **ANSWERED** - Can proceed
⚠️ **PARTIALLY ANSWERED** - Need more details
❌ **NOT ANSWERED** - Blocking MVP development

---

## 1. Technology Stack ✅ FULLY ANSWERED

### Backend
**Answer**: ✅ **Go (Golang) + Gin/Echo framework**
- Excellent performance and scalability
- Low infrastructure costs
- Built-in concurrency for multi-tenant architecture
- Single binary deployment

### Frontend
**Answer**: ✅ **Next.js (React framework)**
- Built-in SSR for better SEO
- Modern, industry-standard
- Great developer experience

### Mobile
**Answer**: ✅ **React Native** (Deferred to Post-MVP)
- For MVP: Mobile-responsive web only
- Native app in Month 4-5 after pilot validation

### Database
**Answer**: ✅ **PostgreSQL**
- Recommended choice confirmed

### Cloud & Infrastructure
**Answer**: ✅ **Google Cloud Platform (Jakarta region)**
- Start with Cloud Run (managed containers)
- Plan Docker architecture from day 1
- Scale to Kubernetes at 100+ cooperatives
- On-premise support planned but not for MVP

**Status**: ✅ Can proceed with development environment setup

---

## 2. MVP Scope ✅ ANSWERED

### MUST HAVE (Month 1-2)
✅ Member Management (basic, no KTP validation yet)
✅ Share Capital Tracking (Simpanan Pokok/Wajib/Sukarela)
✅ Simple POS (cash only)
✅ Basic financial reports (for RAT meeting)
✅ Web app only (no mobile app yet)

### NICE TO HAVE (Month 3)
✅ Basic mobile view (responsive web, not native app)
✅ Simple accounting entries
✅ Transaction history

### DEFINITELY DEFER
✅ SHU Calculation (critical but complex - do it right post-MVP)
✅ Native mobile app
✅ QRIS payments
✅ WhatsApp integration
✅ Inventory management automation

**Status**: ✅ Clear MVP scope defined - can start development

---

## 3. Integration Priorities ✅ ANSWERED

### NIK Validation
**Answer**: ✅ **Skip for MVP**
- Manual KTP entry is fine
- Dukcapil API requires 6+ months government partnership process
- Add fake validation UI for demo purposes only
- Implement real validation in Phase 2

### QRIS Payment
**Answer**: ✅ **Skip for MVP**
- Cash-only POS for pilot
- When ready (Month 4-6): Use Xendit (best documentation, startup-friendly)

### WhatsApp Business API
**Answer**: ✅ **Skip for MVP**
- Not critical for pilot
- Use regular SMS for urgent notifications
- Add WhatsApp in Phase 2 (Month 6+)

### Bank Integration
**Answer**: ✅ **Skip for MVP**
- Manual reconciliation is normal for cooperatives
- Add BNI/Mandiri API in Phase 2

**Status**: ✅ Clear integration priorities - MVP is simpler than originally planned

---

## 4. Localization ✅ ANSWERED

### Language
**Answer**: ✅ **Bahasa Indonesia ONLY for MVP**
- English can wait (even government accepts Indonesian reports)
- Regional languages are Phase 3 nice-to-have

### Date/Number Formatting
**Answer**: ✅ **Indonesian Standard**
- Date: DD/MM/YYYY
- Number: 1.000.000,00 (Indonesian format with dot for thousands, comma for decimals)
- Currency: "Rp" for display, store as integer (satuan)
- Example display: Rp 1.500.000,00

### Fiscal Year
**Answer**: ✅ **Default Calendar Year (Jan-Dec)**
- Allow custom configuration (some cooperatives use July-June fiscal year)
- Make it configurable per cooperative in settings

**Status**: ✅ Clear localization requirements

---

## 5. Regulatory & Compliance ⚠️ PARTIALLY ANSWERED

### Ministry Reports
**Answer**: ⚠️ **Format known, need actual templates**
- Monthly: Simple transaction summary
- Quarterly: Financial position snapshot
- Annual: RAT package - comprehensive reports
- Format: Start with PDF/Excel, online submission coming in 2026

**ACTION NEEDED**: Get actual report templates from existing cooperatives in pilot area

### RAT Documents Required
**Answer**: ⚠️ **List known, need samples**
Required documents:
1. Neraca (Balance Sheet)
2. Laporan Laba Rugi (Income Statement)
3. Laporan Perubahan Modal (Changes in Equity)
4. Laporan SHU dan Pembagiannya (SHU Report & Distribution)
5. Daftar Hadir Anggota (Member Attendance List)

**ACTION NEEDED**: Get sample RAT document package from successful cooperative

### Accounting Standards
**Answer**: ✅ **Use SAK ETAP**
- SAK ETAP (simpler than full SAK)
- No mandatory COA but follow common practice

**ACTION NEEDED**: Get standard COA from successful cooperatives

### SHU Formula
**Answer**: ✅ **Typical formula provided**

```
Total SHU = Net Income after tax

Distribution (typical):
- Cadangan (Reserve): 25%
- Jasa Anggota (Member services): 25%
- Jasa Modal (Capital return): 20%
- Pengurus (Management): 10%
- Karyawan (Employees): 10%
- Sosial (Social): 5%
- Pembangunan Daerah (Regional): 5%
```

**Note**: Each cooperative can adjust percentages in their AD/ART (bylaws)

**ACTION NEEDED**:
- Confirm this is standard
- Make percentages configurable per cooperative
- Get 2-3 real cooperative AD/ART examples

### Data Privacy
**Answer**: ✅ **Follow UU ITE basics**
- Encrypt NIK/KTP data
- No specific cooperative data law yet
- Standard security best practices

**Status**: ⚠️ Partially answered - need real document samples to proceed

---

## 6. User Roles & Permissions ✅ ANSWERED

**Proposed Roles Confirmed**:
1. Super Admin (our team)
2. Cooperative Admin
3. Chairman/Board (read-only)
4. Manager (operational)
5. Finance Officer (accounting)
6. Cashier (POS)
7. Warehouse Staff (post-MVP)
8. Member (view only)

**Answers**:
- ✅ Roles sufficient? **Yes**
- ✅ One person multiple roles? **Yes**
- ✅ Roles customizable per cooperative? **Yes**
- ✅ Role-specific features needed? **Yes**

**Status**: ✅ Clear role structure

---

## 7. Business Rules & Workflows ⚠️ PARTIALLY ANSWERED

### Member Registration
**Answer**: ✅ **Admin registers members**
- No approval workflow initially (simplify MVP)
- Minimum age: 17 (has KTP)
- One cooperative per person for pilot (can expand later)
- No family membership for MVP

### Share Capital Rules
**Answer**: ⚠️ **Typical amounts provided, need to confirm**

**Simpanan Pokok (Principal Share)**:
- One-time payment at joining
- Typical: IDR 50K-100K
- ⚠️ **NEED TO CONFIRM**: Refundable when leaving? Fixed amount or variable?

**Simpanan Wajib (Mandatory Share)**:
- Monthly contribution
- Typical: IDR 10K-50K per month
- ⚠️ **NEED TO CONFIRM**: Penalty for late payment?

**Simpanan Sukarela (Voluntary Share)**:
- Anytime contribution
- Any amount accepted
- ⚠️ **NEED TO CONFIRM**: Withdrawable immediately or with notice?

**MVP Approach**: ✅ Track amounts only, no automation (no late penalties, no auto-deductions)

### Transaction Approval Workflows
**Answer**: ❌ **NOT ANSWERED**

**QUESTIONS**:
- Which transactions need approval for MVP?
- Approval thresholds (e.g., expenses > IDR 5M need manager approval)?
- Or skip approval workflow for MVP entirely?

**RECOMMENDATION FOR MVP**: Skip approval workflows, add in Phase 2

**Status**: ⚠️ Need to clarify share capital rules and confirm approach to approvals

---

## 8. POS System ⚠️ NEEDS CLARIFICATION

### POS Hardware
**Answer**: ✅ **Tablet (Android) and Smartphone**
- Primary: Android tablets
- Fallback: Android smartphones
- Design for touch interface

### POS Features
**Answer**: ❌ **"All of that" - Too broad for MVP**

User answered "all of that" for these features:
- Barcode scanning
- Receipt printing
- Cash drawer integration
- Weight scale integration
- Multiple payment methods
- Discount management
- Shift management
- Member pricing

**CRITICAL QUESTION**: All features means 2-3 months just for POS. Not realistic for MVP.

**RECOMMENDATION - MVP POS (Month 1-2)**:
- ✅ Product catalog (name, price, basic info)
- ✅ Add to cart, calculate total
- ✅ Cash payment only
- ✅ Digital receipt (on screen, no printing yet)
- ✅ Basic sales report
- ❌ No barcode scanning (manual entry)
- ❌ No receipt printing (add Month 3)
- ❌ No cash drawer integration (manual cash management)
- ❌ No weight scale
- ❌ No QRIS (defer to Phase 2)
- ❌ No discount management (fixed prices only)
- ❌ No shift management (simple daily report)
- ❌ No member pricing (same price for all)

**POST-MVP (Month 3-6)**:
- Add barcode scanning
- Add receipt printing
- Add member pricing
- Add shift management

**Phase 2 (Month 6+)**:
- Add QRIS payments
- Add discount/promo management
- Add cash drawer integration

### Offline Mode
**Answer**: ⚠️ **Conflicting - "No" but want sync**
- Offline important? **No**
- Should sync when back online? **Yes**
- Conflict handling? **Mark for later**

**CLARIFICATION NEEDED**: If not important, skip entirely for MVP. If some offline needed, must build from start.

**RECOMMENDATION**: Skip offline mode for MVP, require internet connection.

**Status**: ⚠️ **CRITICAL** - Need to narrow POS scope or MVP will take 6 months

---

## 9. Inventory Management ✅ ANSWERED (Defer to Phase 2)

### Stock Valuation
**Answer**: ✅ **Follow Indonesian accounting standard**
- Support FIFO (primary)
- Support Average Cost
- Allow cooperative to choose method

### Product Features
**Requested for Phase 2** (not MVP):
- Product variants (size, color) - YES
- Batch/lot tracking - YES
- Serial number tracking - YES
- Expiry date tracking - YES
- Supplier linking - NO (one product = one supplier for now)

### Stock Counting
**Answer**: ✅ **No mobile app for counting**
- Manual stock count entry
- No variance approval workflow for MVP

**Status**: ✅ Clear - Defer full inventory to Phase 2, basic product catalog only for MVP

---

## 10. Mobile App Features ✅ ANSWERED

**Answer**: ✅ **Responsive web only, no native app for MVP**

### For Members
**MVP**:
- View share capital balance
- View transaction history
- View announcements
- Contact cooperative

**Phase 2**:
- Digital member card with QR code
- Apply for loan
- Push notifications

### For Staff
**Answer**: ✅ **Same app with role-based views**
- Staff need on mobile:
  - ✅ Approve transactions
  - ✅ View reports
  - ❌ Stock counting (not for MVP)

**Status**: ✅ Clear scope for mobile web

---

## 11. Reporting Requirements ⚠️ UNREALISTIC FOR MVP

### Financial Reports Requested
User marked ALL as required for MVP:
- ✅ Balance Sheet (Neraca)
- ✅ Income Statement (Laba Rugi)
- ✅ Cash Flow Statement
- ✅ Changes in Equity
- ✅ Trial Balance
- ✅ General Ledger
- ✅ Journal Entry Report

### Operational Reports Requested
- ✅ Sales Report (daily, monthly)
- ✅ Purchase Report
- ✅ Inventory Report
- ✅ Member Transaction Report
- ✅ SHU Report

### Export Formats Requested
- ✅ PDF
- ✅ Excel
- ✅ Print directly

**CRITICAL ISSUE**: This is 4-6 weeks of work just for reports. Not realistic for 3-month MVP.

**RECOMMENDATION - MVP Reports (Essential Only)**:
**Month 1-2**:
1. Trial Balance (essential for accounting validation)
2. Member List with Share Capital
3. Daily Sales Report
4. Simple Transaction History

**Month 3**:
5. Balance Sheet (basic)
6. Income Statement (basic)
7. Member Transaction Report

**Phase 2 (Month 4-6)**:
8. Cash Flow Statement
9. Changes in Equity
10. General Ledger
11. SHU Report
12. Advanced exports (Excel, PDF)

**For MVP**: PDF export only. Excel in Phase 2.

**Status**: ⚠️ **CRITICAL** - Need to reduce reporting scope or extend timeline

---

## 12. Subscription/Billing ✅ ANSWERED

**Answer**: ✅ **Manual for MVP**
- Manual billing (send invoice)
- Bank transfer payment
- Manual activation
- No enforcement (trust-based for pilot)
- No automated limits

**Defer**: Automated billing system to Phase 2

**Status**: ✅ Simple approach for pilot

---

## 13. Data Migration ✅ ANSWERED

**Answer**: ✅ **Excel template import**

**Process**:
1. We provide Excel template
2. Cooperatives fill:
   - Member list (Name, NIK, Phone, Join Date)
   - Share capitals (Pokok, Wajib, Sukarela)
   - Opening cash balance
3. We import manually (one-time)
4. Start fresh for transactions (no historical transaction import)

**Timeline**: 2-3 days per cooperative = 20-30 days for 10 cooperatives

**Status**: ✅ Clear migration process

---

## 14. Support & Training ✅ ANSWERED

### Support Channels for Pilot
**Answer**: ✅
- Email
- WhatsApp (business account)
- Phone
- ❌ No live chat (too complex for pilot)
- ❌ No help desk system (manual for pilot)

### Training Approach
**Answer**: ✅
- On-site training at cooperative
- Virtual training (Zoom/Google Meet)
- Written documentation (Bahasa Indonesia only)
- In-app guided tours
- ❌ No video tutorials for MVP (Phase 2)
- ❌ No train-the-trainer yet

**Status**: ✅ Clear support plan

---

## 15. Security & Authentication ✅ ANSWERED

### Authentication Method
**Answer**: ✅
- Username + Password (primary)
- Phone + OTP (optional for member verification)
- ❌ No social login
- ❌ No email-based auth

### Password & Session
**Answer**: ✅ **Simple approach**
- Minimum 6 characters (no complexity requirements to reduce friction)
- 30-minute session timeout
- Single device login
- ❌ No 2FA for MVP (add for finance role in Phase 2)

**Status**: ✅ Balance security and usability

---

## 16. Multi-tenancy Architecture ✅ ANSWERED

**Answer**: ✅ **Shared Database, Shared Schema**
- Add `cooperative_id` to all tables
- Simple but sufficient for pilot
- Easy to migrate to separate schemas later if needed
- Cost-effective for early stage

**Customization**:
- COA customization: **Yes** (each cooperative can customize)
- Report formats: **Phase 2**
- Workflows: **Phase 2**
- Branding: **Phase 2**

**Status**: ✅ Clear architecture approach

---

## 17. Performance Targets ✅ ANSWERED

**MVP Targets**:
- 10 cooperatives
- 3,000 total members
- 1,000 transactions/day
- Page load < 3 seconds on 3G network

**Status**: ✅ Realistic targets for MVP

---

## 18. Testing Strategy ❌ NOT ANSWERED

**QUESTIONS**:
- Manual testing only for MVP?
- Or include automated tests (slower development, more reliable)?
- UAT process with pilot cooperatives?

**RECOMMENDATION**:
- Manual testing for MVP (faster)
- Add automated tests in Phase 2
- UAT with 2-3 pilot cooperatives before rolling out to all 10

**Status**: ❌ Need decision on testing approach

---

## 19. Legal & Compliance ✅ ANSWERED

**Answer**: ✅ **Simple for pilot**
- Simple Terms of Service
- Data owned by cooperative
- 30-day data retention after stop paying
- No SLA for pilot (best-effort basis)
- ❌ No insurance for pilot
- ❌ No lawyer review for pilot (add before commercial launch)

**Status**: ✅ Appropriate for pilot phase

---

## CRITICAL ISSUES TO RESOLVE

### 🔴 BLOCKING ISSUES (Must resolve before development):

1. **POS Feature Scope** ⚠️
   - User wants "all features" but that's 2-3 months of work
   - **ACTION**: Confirm simplified POS for MVP (basic only, no barcode/printing/etc)

2. **Reporting Scope** ⚠️
   - User wants ALL reports but that's 4-6 weeks of work
   - **ACTION**: Confirm reduced reporting for MVP (4 essential reports only)

3. **Government Report Templates** ⚠️
   - Need actual templates to build correctly
   - **ACTION**: Get sample reports from 2-3 existing cooperatives

4. **Chart of Accounts Template** ⚠️
   - Need standard COA to implement accounting properly
   - **ACTION**: Get COA from successful cooperative

5. **Share Capital Business Rules** ⚠️
   - Need to confirm refundability, withdrawal rules, penalties
   - **ACTION**: Interview 2-3 cooperative managers for exact rules

### 🟡 IMPORTANT BUT NOT BLOCKING:

6. **Offline POS Mode**
   - User said "no" but wants sync - conflicting
   - **DECISION**: Skip offline for MVP? Or build from start?

7. **Transaction Approval Workflows**
   - Not answered
   - **DECISION**: Skip for MVP, add Phase 2?

8. **Testing Strategy**
   - Not answered
   - **DECISION**: Manual only for MVP?

---

## REVISED MVP TIMELINE WITH CURRENT SCOPE

**If we build "all features" user requested**:
- POS with all features: 6-8 weeks
- All reports: 4-6 weeks
- Core features: 6-8 weeks
- **Total**: 16-22 weeks (4-5 months) ❌ Too long

**If we follow RECOMMENDED MVP scope**:
- Simple POS (basic only): 2-3 weeks
- Essential reports (4 reports): 2 weeks
- Core features: 4-6 weeks
- Testing & bug fixing: 2 weeks
- **Total**: 10-13 weeks (2.5-3 months) ✅ Achievable

---

## IMMEDIATE NEXT STEPS

### This Week:
1. **Schedule 2-hour scoping call** to resolve blocking issues
2. **Get real cooperative examples**:
   - Visit 2-3 cooperatives in pilot area
   - Get Chart of Accounts
   - Get RAT report samples
   - Get member registration forms
   - Interview about share capital rules

3. **Finalize MVP scope**:
   - Confirm simplified POS (basic only)
   - Confirm reduced reporting (4 essential reports)
   - Confirm skip offline mode
   - Document final MVP feature list

### Next Week:
4. **Create detailed PRD** (Product Requirements Document)
5. **Database schema design** based on confirmed COA and business rules
6. **Start development environment setup**:
   - Go backend skeleton
   - Next.js frontend skeleton
   - PostgreSQL setup
   - Docker configuration

---

## MVP FEATURE CHECKLIST (FINAL CONFIRMATION NEEDED)

### Core Features (CONFIRMED)
- [x] User authentication (username/password)
- [x] User roles (8 roles)
- [x] Cooperative settings
- [x] Member registration (admin only)
- [x] Member list and search
- [x] Share capital tracking (3 types)
- [x] Share capital transaction recording

### POS System (NEEDS CONFIRMATION)
- [ ] ⚠️ Basic POS (manual entry, cash only) - **OR**
- [ ] ⚠️ Full POS (barcode, printing, drawer, scale, etc.) - **WHICH ONE?**

### Accounting (NEEDS CONFIRMATION)
- [x] Chart of Accounts setup
- [ ] ⚠️ Manual journal entry only - **OR**
- [ ] ⚠️ Automated journal entry from POS/transactions - **WHICH ONE?**

### Reporting (NEEDS CONFIRMATION)
- [ ] ⚠️ 4 essential reports (Trial Balance, Member List, Sales, Transactions) - **OR**
- [ ] ⚠️ All 12 reports (Balance Sheet, P&L, Cash Flow, etc.) - **WHICH ONE?**

### Mobile (CONFIRMED)
- [x] Responsive web (works on mobile browser)
- [x] No native app for MVP

---

**Version**: 2.0.0 - Updated based on user answers
**Last Updated**: 2025-11-15
**Document Owner**: Product Management

**Status**:
✅ Technology stack: CONFIRMED
✅ MVP direction: MOSTLY CLEAR
⚠️ Feature scope: NEEDS REDUCTION
⚠️ Real examples: NEEDED
❌ Government templates: BLOCKING

**CRITICAL PATH**:
1. Get real cooperative examples (this week)
2. Reduce POS and reporting scope (this week)
3. Create PRD (next week)
4. Start development (week 3)
