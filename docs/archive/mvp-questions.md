# MVP-Specific Questions (3-Month Pilot Program)

## Purpose
These questions are CRITICAL to answer before starting MVP development for the 10 pilot cooperatives (Months 1-3).

**MVP Goal**: Prove the core value proposition with minimum features
**MVP Success**: 10 cooperatives using the system daily, satisfied enough to pay and refer others

---

## CRITICAL MVP QUESTIONS (Must Answer Before Coding)

### 1. MVP Technology Stack ✅ (ANSWERED)

**Backend**: Go (Golang) + Gin/Echo framework
**Frontend**: Next.js (React framework)
**Mobile**: React Native
**Database**: PostgreSQL
**Cloud**: Google Cloud Platform (Jakarta region)
**Infrastructure**: Start with Cloud Run (managed), plan for containers

**Why Go for Backend:**
- ✅ Excellent performance and low memory usage (cost-effective scaling)
- ✅ Built-in concurrency perfect for handling multiple cooperatives
- ✅ Single binary deployment (easier DevOps)
- ✅ Strong typing prevents common bugs
- ✅ Great for microservices if we need to split later
- ✅ Lower cloud infrastructure costs than Node.js/Python

**Status**: ✅ ANSWERED - Can proceed with development environment setup

---

### 2. MVP Scope Definition (NEEDS ANSWER)

**Question**: For the 10 pilot cooperatives to succeed, what is the MINIMUM they need?

#### Proposed MVP Features (Choose):

**Tier 1 - Absolutely Must Have:**
- [ ] User login and authentication
- [ ] Member registration (name, KTP/NIK, phone, address)
- [ ] Share capital recording (Simpanan Pokok, Wajib, Sukarela)
- [ ] Basic Chart of Accounts (cooperative template)
- [ ] Journal entry (manual posting)
- [ ] Cash receipts and payments
- [ ] Basic POS (cash only, 1 business unit)
- [ ] Simple reports (member list, transaction history, trial balance)

**Tier 2 - Important but Could Manual Workaround:**
- [ ] SHU calculation (or do manually in Excel for pilot?)
- [ ] Member mobile app (or just web browser on phone?)
- [ ] Inventory tracking (or just record sales, manage stock in Excel?)
- [ ] NIK validation via Dukcapil (or just record NIK without validation?)
- [ ] Digital member card with QR code (or print card later?)

**Tier 3 - Defer to Post-MVP:**
- [ ] Multiple business units (start with 1 unit only)
- [ ] QRIS payment (cash only for pilot)
- [ ] WhatsApp notifications (use SMS/email)
- [ ] Savings & Loans module (just basic member accounts)
- [ ] Advanced reporting (RAT automation)

**CRITICAL DECISION NEEDED:**
What goes in Tier 1 (MVP) vs Tier 2 (post-pilot) vs Tier 3 (later)?

**Recommendation**:
Focus on **Member Transparency** as core value:
- Members can see their share capital balance
- Members can see their transaction history
- Managers can record transactions easily
- Basic financial reports work

Everything else can wait or be done manually during pilot.

---

### 3. Share Capital Business Rules (NEEDS ANSWER)

**For MVP, we need to know the EXACT rules:**

**Simpanan Pokok (Principal Share):**
- [ ] Paid when? (At registration only, or can add later?)
- [ ] Fixed amount for all members? Or varies?
- [ ] Typical amount: IDR _______ per member
- [ ] Refundable when member leaves? Yes/No
- [ ] Can increase during membership? Yes/No

**Simpanan Wajib (Mandatory Share):**
- [ ] Frequency: Monthly / Quarterly / Annually?
- [ ] Fixed amount? IDR _______ per period
- [ ] Or percentage of transaction? ____%
- [ ] Late payment penalty? Yes/No, if yes: ____%
- [ ] Automatic deduction? Or manual payment?

**Simpanan Sukarela (Voluntary Share):**
- [ ] Minimum amount? IDR _______
- [ ] Maximum amount? IDR _______ or unlimited
- [ ] Can withdraw? Yes/No, if yes: How quickly?
- [ ] Interest paid? Yes/No, if yes: ___% per year

**CRITICAL**: Need actual examples from real cooperatives to configure MVP correctly.

---

### 4. Chart of Accounts for MVP (NEEDS TEMPLATE)

**Question**: Is there a standard Chart of Accounts for Indonesian cooperatives?

**Options:**

**Option A**: Use official Ministry template (if exists)
- **Need**: Official COA template from Ministry of Cooperatives
- **Pros**: Compliance guaranteed, familiar to cooperatives
- **Cons**: May be too complex for small cooperatives

**Option B**: Create simplified cooperative COA
- **Need**: Sample from 2-3 successful cooperatives
- **Pros**: Simpler, easier to understand
- **Cons**: May not be fully compliant

**Option C**: Customizable per cooperative
- **Need**: Base template + customization interface
- **Pros**: Flexible
- **Cons**: More complex to build and support

**CRITICAL DECISION**: Which option for MVP? Need actual COA template to proceed.

---

### 5. Member Registration Process for MVP (NEEDS ANSWER)

**Question**: What is the EXACT workflow for registering a new member?

**Process Steps (Confirm):**
1. [ ] Member fills application form (on paper or digital?)
2. [ ] Board/Manager reviews application
3. [ ] Board approves in meeting? Or manager can approve directly?
4. [ ] Member pays Simpanan Pokok
5. [ ] Member receives member ID number (auto-generated or manual?)
6. [ ] Member receives member card (digital or print?)

**For MVP:**
- Can manager register member directly? Or always needs approval workflow?
- Can member self-register online? Or only staff can register?
- Is approval workflow critical for MVP? Or can add later?

**RECOMMENDATION**: MVP = Simple direct registration by staff, add approval workflow post-pilot.

---

### 6. Transaction Recording for MVP (NEEDS ANSWER)

**Question**: What types of transactions must be recorded in MVP?

**Essential Transactions:**
- [ ] Share capital contribution (Simpanan Pokok/Wajib/Sukarela)
- [ ] Sales (POS transactions)
- [ ] Cash receipts (member payments)
- [ ] Cash disbursements (expenses)
- [ ] Loan disbursement (if doing loans in pilot)
- [ ] Loan repayment (if doing loans in pilot)

**For Each Transaction Type:**
- What information must be captured?
- What approvals needed?
- What documents generated (receipts, invoices)?

**CRITICAL**: Need sample transaction forms from real cooperatives

---

### 7. Member Mobile App for MVP (NEEDS DECISION)

**Question**: Is mobile app essential for MVP, or can members use web browser?

**Option A**: Build mobile app (React Native)
- **Timeline**: +3-4 weeks development
- **Pros**: Better user experience, works offline, push notifications
- **Cons**: Slower MVP launch, need app store approval

**Option B**: Mobile-responsive web only
- **Timeline**: Included in web development
- **Pros**: Faster to launch, no app store hassle
- **Cons**: Can't work offline, no push notifications

**Option C**: Web for MVP, mobile app post-pilot
- **Timeline**: MVP in 3 months, mobile in Month 4-5
- **Pros**: Fastest MVP, can validate before investing in mobile
- **Cons**: Members expect mobile app

**CRITICAL DECISION**: Which option? Consider pilot timeline and member expectations.

**RECOMMENDATION**: Option C - Web first (mobile-responsive), native app after pilot validation.

---

### 8. Reporting for MVP (NEEDS PRIORITIZATION)

**Question**: Which reports are ESSENTIAL for pilot cooperatives?

**Financial Reports (Priority?):**
- [ ] Trial Balance (Essential? Yes/No)
- [ ] Balance Sheet (Essential? Yes/No)
- [ ] Income Statement/P&L (Essential? Yes/No)
- [ ] Cash Flow Statement (Essential? Yes/No)
- [ ] General Ledger (Essential? Yes/No)

**Operational Reports (Priority?):**
- [ ] Member List (Essential? Yes/No)
- [ ] Member Share Capital Balance (Essential? Yes/No)
- [ ] Daily Sales Report (Essential? Yes/No)
- [ ] Member Transaction History (Essential? Yes/No)

**For MVP:**
What are the TOP 3 reports that cooperatives look at daily/weekly?

**RECOMMENDATION**:
1. Member share capital balance (transparency = core value)
2. Daily sales report (operational need)
3. Trial balance (financial control)

Everything else can be manual/Excel for pilot.

---

### 9. POS System for MVP (NEEDS SIMPLIFICATION)

**Question**: What is the MINIMUM POS functionality for pilot?

**Basic POS Features:**
- [ ] Product catalog (name, price, SKU)
- [ ] Add item to cart
- [ ] Calculate total
- [ ] Record cash payment
- [ ] Print/show receipt
- [ ] End-of-day sales report

**Advanced POS Features (Include in MVP?):**
- [ ] Barcode scanning (Need? Yes/No)
- [ ] Discount management (Need? Yes/No)
- [ ] Multiple payment methods (Need? Yes/No - or just cash?)
- [ ] Inventory deduction (Need? Yes/No - or track separately?)
- [ ] Member pricing (Need? Yes/No - or same price for all?)
- [ ] Shift management (Need? Yes/No)

**CRITICAL DECISION**: Cash-only simple POS for MVP? Or need more features?

**RECOMMENDATION**:
- MVP = Cash only, no barcode, no inventory deduction, same price for all
- Manual stock count daily
- Add advanced features post-pilot

---

### 10. Data Migration for Pilot (NEEDS PROCESS)

**Question**: What existing data do pilot cooperatives need to migrate?

**Data Migration Scope:**
- [ ] Member list (how many members per cooperative? Average: _____ )
- [ ] Opening balances (just totals, or full transaction history?)
- [ ] Share capital balances (current balances only? Or full history?)
- [ ] Product catalog (if doing POS)
- [ ] Opening cash balance
- [ ] Opening bank balance

**Migration Method for 10 Cooperatives:**
- [ ] **Option A**: Excel template, they fill, we import
- [ ] **Option B**: We visit and enter data manually
- [ ] **Option C**: They enter data themselves after training

**Timeline Estimate:**
- How long to migrate 1 cooperative's data? _____ hours/days
- Total for 10 cooperatives? _____ hours/days

**CRITICAL**: This affects pilot timeline. If migration takes 2 weeks per cooperative, might only do 5 cooperatives in 3 months.

**RECOMMENDATION**: Excel template + we help enter data. Budget 2-3 days per cooperative.

---

### 11. Training for Pilot (NEEDS PLAN)

**Question**: How to train pilot cooperative staff in limited time?

**Who Needs Training:**
- Managers (how many? ____ per cooperative)
- Finance officers (how many? ____ per cooperative)
- Cashiers (how many? ____ per cooperative)
- Members (train all? Or just show how to use app?)

**Training Approach:**
- [ ] **Option A**: On-site training at each cooperative (2 days)
- [ ] **Option B**: Bring all cooperatives to 1 location (2 days workshop)
- [ ] **Option C**: Virtual training via Zoom (4 sessions x 2 hours)
- [ ] **Option D**: Hybrid: 1 day workshop + on-site support

**Training Content:**
- Day 1: System overview, member registration, share capital, basic POS
- Day 2: Accounting basics, reports, troubleshooting

**CRITICAL DECISION**: Training budget and logistics for 10 cooperatives.

**RECOMMENDATION**: Option D - 1 day workshop with all 10 cooperatives + 1 day on-site per cooperative

---

### 12. Support During Pilot (NEEDS STRUCTURE)

**Question**: How to support 10 cooperatives during 3-month pilot?

**Support Channels for Pilot:**
- [ ] WhatsApp group (all 10 cooperatives + our team)
- [ ] Phone hotline (business hours only?)
- [ ] On-site visits (how often? Weekly? Monthly? On-demand?)
- [ ] Video calls (for troubleshooting)

**Response Time Expectations:**
- Critical issues (system down): _____ hours
- Urgent questions (can't process transaction): _____ hours
- Normal questions (how to do X): _____ hours

**Support Team Size:**
- _____ dedicated support person(s) for 10 cooperatives
- _____ developer on-call for bugs
- _____ product manager for feedback

**RECOMMENDATION**:
- 1 dedicated customer success person
- WhatsApp group for quick questions
- Weekly check-in calls
- On-site visit once per month per cooperative

---

### 13. Success Metrics for Pilot (NEEDS DEFINITION)

**Question**: How do we know if MVP/pilot is successful?

**Usage Metrics:**
- [ ] ___% of cooperatives use system daily
- [ ] ___% of members registered in system
- [ ] ___% of transactions recorded in system (vs manual)
- [ ] ___ transactions per day per cooperative (average)

**Satisfaction Metrics:**
- [ ] Customer Satisfaction Score: ___ / 10 (target)
- [ ] Net Promoter Score: ___ (target)
- [ ] Would recommend to other cooperatives: ___% (target)

**Business Metrics:**
- [ ] ___ cooperatives willing to pay (out of 10)
- [ ] ___ referrals generated
- [ ] ___ case studies completed (target: 3)

**CRITICAL**: Define success criteria BEFORE pilot, not after.

**RECOMMENDATION**:
- 8/10 cooperatives use daily
- CSAT > 8/10
- 7/10 willing to pay after free pilot
- 3 strong case studies

---

### 14. Pilot Cooperative Selection (NEEDS CRITERIA)

**Question**: How to select the 10 pilot cooperatives?

**Selection Criteria (Rank Importance):**
- [ ] Size: 100-500 members (manageable, not too small)
- [ ] Location: Same district (easier support)
- [ ] Internet: Reliable connectivity (reduce technical issues)
- [ ] Leadership: Open to technology (early adopters)
- [ ] Operations: Active business (real transactions to test)
- [ ] Diversity: Mix of types (retail, loans, agriculture)
- [ ] Influence: Well-connected (can refer others)

**Red Flags (Don't Select):**
- Dormant cooperatives (no active transactions)
- Resistant leadership (forced to try by others)
- Too small (<50 members) or too large (>1000 members)
- No smartphone/internet access

**CRITICAL**: Right pilot cooperatives = successful pilot. Wrong selection = wasted effort.

---

### 15. MVP Timeline Reality Check (NEEDS VALIDATION)

**Question**: Is 3 months realistic for MVP with 10 pilots?

**Development Time:**
- Environment setup: 1 week
- Core features development: 6-8 weeks
- Testing: 1-2 weeks
- Bug fixing: 1-2 weeks
- **Total**: 9-13 weeks (2-3 months)

**Pilot Onboarding Time:**
- Cooperative selection: 1 week
- Data migration per cooperative: 2-3 days
- Training per cooperative: 1-2 days
- Stabilization per cooperative: 1 week
- **Total per cooperative**: 2-3 weeks
- **For 10 cooperatives (parallel)**: 4-6 weeks

**Realistic Timeline:**
- Month 1: Development (Weeks 1-4)
- Month 2: Development + First 5 cooperatives onboarding (Weeks 5-8)
- Month 3: Final development + Last 5 cooperatives + Stabilization (Weeks 9-12)

**CRITICAL QUESTION**: Is this realistic? Or should we:
- Extend to 4 months?
- Reduce to 5 pilot cooperatives?
- Reduce MVP scope?

---

## MVP DECISIONS SUMMARY TABLE

| Decision Point | Option A | Option B | Option C | Recommended |
|---------------|----------|----------|----------|-------------|
| **Backend** | Node.js | PHP | Go | ✅ Go/Golang (Answered) |
| **Frontend** | Next.js | Vue | React | ✅ Next.js (Answered) |
| **Mobile** | Now | Post-pilot | Never | ⚠️ Post-pilot |
| **SHU Calc** | Automated | Manual | Skip | ⚠️ Manual (Excel for pilot) |
| **NIK Valid** | API | Manual entry | Skip | ⚠️ Manual entry |
| **POS** | Full features | Cash only | Skip | ⚠️ Cash only |
| **QRIS** | Integrated | Manual | Skip | ⚠️ Skip (post-pilot) |
| **Multi-Unit** | Yes | Single only | N/A | ⚠️ Single only |
| **Inventory** | Automated | Manual count | Skip | ⚠️ Manual count |
| **Training** | On-site | Workshop | Virtual | ⚠️ Hybrid |
| **Pilot Size** | 10 coops | 5 coops | 15 coops | ⚠️ 10 coops |
| **Timeline** | 3 months | 4 months | 2 months | ⚠️ 3 months |

⚠️ = Needs decision before proceeding

---

## IMMEDIATE NEXT STEPS

### Week 1: Critical Decisions Workshop
**Attendees**: Product team, Business team, 2-3 cooperative advisors

**Agenda**:
1. MVP scope finalization (Tier 1 vs Tier 2 vs Tier 3)
2. Share capital business rules definition
3. Chart of Accounts template selection
4. Pilot cooperative selection criteria
5. Timeline validation

**Deliverable**: MVP Requirements Document (PRD)

---

### Week 2: Requirements Gathering
1. **Get real examples**:
   - Visit 2-3 cooperatives
   - Observe their current processes
   - Collect sample forms and documents
   - Interview members, managers, finance officers

2. **Collect templates**:
   - Member registration form
   - Share capital receipt
   - Transaction forms
   - Monthly reports they currently use

3. **Government requirements**:
   - Contact Ministry of Cooperatives
   - Get official Chart of Accounts (if exists)
   - Get SHU calculation guidelines
   - Get reporting templates

---

### Week 3-4: Technical Foundation
1. **Backend Setup (Go)**:
   - Initialize Go project with Gin/Echo framework
   - PostgreSQL connection and ORM (GORM)
   - JWT authentication middleware
   - RESTful API structure
   - Docker containerization

2. **Frontend Setup (Next.js)**:
   - Initialize Next.js project with TypeScript
   - Authentication flow
   - Basic component library
   - API integration layer

3. **Database**:
   - Schema design v1
   - Migration scripts
   - Seed data for testing

**Parallel**: Pilot cooperative recruitment

---

## QUESTIONS FOR COOPERATIVE ADVISORS

**To validate our assumptions, please help answer:**

1. **What do cooperative managers spend most time on daily?**
   (This tells us what to automate first)

2. **What causes members to lose trust in cooperatives?**
   (This tells us what transparency features matter most)

3. **What reports do cooperatives look at every day/week?**
   (This tells us what reports to prioritize)

4. **When you visit a cooperative, what do you check to know if it's healthy?**
   (This tells us what metrics to display on dashboard)

5. **What makes RAT preparation painful?**
   (This tells us what to automate, even if not in MVP)

---

**Version**: 1.0.1
**Last Updated**: 2025-11-15
**Document Owner**: Product Management

**Status**:
- ✅ Tech stack decided (Go/Golang, Next.js, React Native, PostgreSQL, GCP)
- ⚠️ Need decisions on MVP scope, timeline, features
- ⚠️ Need real cooperative examples and templates
- ⚠️ Need government requirements and standards
