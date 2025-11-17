# Phase 1: MVP (Month 1-3)

**The 12-Week Sprint to Launch**

**Duration:** 12 weeks
**Team Size:** 5 people
**Goal:** Launch with 10 pilot cooperatives

---

## 📋 What's in This Folder

This folder contains **three main documents** for Phase 1 MVP:

| Document | Purpose | Update Frequency | Owner |
|----------|---------|------------------|-------|
| **[📚 Implementation Guide](./implementation-guide.md)** | Technical reference with code examples, architecture, and patterns | As needed | Tech Lead |
| **[📈 Progress Tracking](./progress-tracking.md)** | Current status, completed work, metrics, and blockers | Weekly (Monday) | Tech Lead / PM |
| **[✅ Todolist](./todolist.md)** | All tasks by week with checkboxes and priorities | Daily | Product Manager |

### Quick Guide to Using These Documents

**🆕 Starting Development?**
- Read → [Implementation Guide](./implementation-guide.md)
- Reference code examples and technical patterns

**📊 Need Current Status?**
- Check → [Progress Tracking](./progress-tracking.md)
- See what's done, what's blocked, metrics

**✅ Planning Your Week?**
- Review → [Todolist](./todolist.md)
- See all tasks, mark completed items

---

## 🎯 Phase 1 Overview

### What is MVP?

The **Minimum Viable Product** - the smallest version that:
- ✅ Solves the core problem (manual bookkeeping)
- ✅ Delivers value to 10 cooperatives
- ✅ Proves the concept works
- ✅ Generates initial revenue (IDR 10M/month)
- ✅ Provides foundation for Phase 2

### 8 Core Features

1. **User Authentication & Roles** - Secure login with role-based access
2. **Member Management** - Complete CRUD for cooperative members
3. **Share Capital Tracking** - Simpanan Pokok, Wajib, Sukarela
4. **Basic POS (Cash Only)** - Simple point-of-sale for retail
5. **Simple Accounting** - Manual journal entries, chart of accounts
6. **4 Essential Reports** - Financial position, income, cash flow, member balances
7. **Member Portal (Web)** - Members can view balances online
8. **Data Import** - Import members and products from Excel

---

## 📊 MVP Targets (End of Week 12)

| Metric | Target |
|--------|--------|
| **Cooperatives** | 10 pilot cooperatives |
| **Active Users** | 40+ users logging in weekly |
| **Transactions** | 2,000+ per month |
| **Revenue (MRR)** | IDR 10M |
| **Uptime** | 95%+ |
| **NPS Score** | 45+ |
| **Critical Bugs** | 0 |

---

## 📅 12-Week Timeline

### **Month 1: Foundation**
- **Week 1:** Backend foundation (database, auth)
- **Week 2:** Frontend + Member management
- **Week 3:** Share capital tracking
- **Week 4:** Simple accounting system

### **Month 2: Core Features**
- **Week 5:** Product management + POS backend
- **Week 6:** POS frontend
- **Week 7:** 4 Essential reports
- **Week 8:** Member portal

### **Month 3: Launch**
- **Week 9:** Testing & bug fixing
- **Week 10:** Deployment to production
- **Week 11:** Onboard first 6 cooperatives
- **Week 12:** Onboard final 4 cooperatives (total 10) ✅

---

## ✅ Success Criteria

**MVP is DONE when:**

### Technical ✅
- All 8 features working in production
- 10 cooperatives migrated from Excel
- 95%+ uptime maintained
- Zero critical bugs
- Tests passing (70%+ coverage)

### Business ✅
- 10 paying cooperatives
- MRR: IDR 10M
- NPS score 45+
- 0% churn (all cooperatives stay)

### User ✅
- Reports generate in < 5 minutes
- POS checkout < 30 seconds
- Members can check balances anytime
- Staff trained in < 1 day

---

## 🚀 Quick Start

**Week 0 (Now):**
1. Read [implementation-guide.md](implementation-guide.md) completely
2. Install Go 1.25.4, Node 20.18.1, PostgreSQL 17.2
3. Visit 2-3 local cooperatives
4. Collect sample Chart of Accounts
5. Team kickoff meeting

**Week 1 (MVP Starts):**
1. Create project folder structure
2. Initialize backend (Go + Gin)
3. Setup database (PostgreSQL)
4. Build authentication system
5. Initialize frontend (Next.js)

---

## 📁 Key Files

**Backend (Go):**
```
backend/
├── cmd/api/main.go              # Server entry point
├── internal/
│   ├── models/                  # 10 database models
│   ├── handlers/                # API handlers
│   ├── services/                # Business logic
│   └── middleware/              # Auth, CORS
```

**Frontend (Next.js):**
```
frontend/
├── src/app/
│   ├── (auth)/login/            # Login page
│   ├── (dashboard)/             # Main app
│   │   ├── members/             # Member management
│   │   ├── share-capital/       # Capital tracking
│   │   ├── pos/                 # Point of sale
│   │   ├── accounting/          # Journal entries
│   │   └── reports/             # 4 essential reports
│   └── portal/                  # Member portal
```

---

## 💰 Budget

**Total:** IDR 148.5M for 12 weeks

**Breakdown:**
- Personnel (5 people): IDR 90M (60%)
- Infrastructure: IDR 22M (15%)
- Tools & Services: IDR 15M (10%)
- Marketing: IDR 15M (10%)
- Contingency: IDR 6.5M (5%)

---

## 👥 Team

| Role | Responsibilities |
|------|------------------|
| Tech Lead / Backend Dev | Go backend, database, architecture |
| Frontend Developer | Next.js, UI/UX, responsive design |
| Full-Stack Developer | Backend + Frontend support |
| Product Manager | Requirements, testing, deployment |
| Designer / QA | UI design, testing, documentation |

**Daily standup:** 9:00 AM (15 minutes)
**Weekly planning:** Friday 2:00 PM (1 hour)

---

## 🎉 What MVP Achieves

### For Cooperatives

**Before:**
- Manual bookkeeping (paper/Excel)
- Hours to generate reports
- Calculation errors
- No member transparency

**After:**
- Digital bookkeeping (cloud)
- Reports in 5 minutes
- Accurate calculations
- Member portal access

### For Business

**Validation:**
- ✅ Technology works
- ✅ Users adopt it
- ✅ Revenue model proven
- ✅ Ready to scale (Phase 2)

**Foundation:**
- Stable codebase
- Scalable architecture
- 10 reference customers
- Product-market fit validated

---

## 🎯 What's NOT in MVP

Deferred to Phase 2:
- ❌ SHU Calculation (Phase 2, Week 4-5)
- ❌ QRIS Payments (Phase 2, Week 6-7)
- ❌ WhatsApp Integration (Phase 3)
- ❌ Native Mobile App (Phase 2, Week 8-9)
- ❌ Inventory Automation (Phase 3)
- ❌ Barcode Scanning (Phase 2, Week 10)
- ❌ Receipt Printing (Phase 2, Week 10)
- ❌ Loan Management (Phase 2, Week 2-3)

**Why?** Focus on core value. Ship fast. Learn. Iterate.

---

## 💡 Success Tips

### Do's ✅
- Ship features incrementally
- Test with real cooperatives early
- Focus on core value
- Daily communication
- Celebrate small wins

### Don'ts ❌
- Don't add features mid-sprint
- Don't skip testing
- Don't work in isolation
- Don't ignore user feedback
- Don't burn out the team

---

## 📈 Current Progress (Quick View)

**Last Updated:** November 17, 2025

```
Week 1:  [✅] Backend foundation (100%)
Week 2:  [🔄] Frontend + Members (80% - backend done)
Week 3:  [✅] Share capital (100%)
Week 4:  [✅] Accounting (100%)
Week 5:  [⏳] Products & POS backend (0%)
Week 6:  [⏳] POS frontend (0%)
Week 7:  [⏳] Reports (0%)
Week 8:  [⏳] Member portal (0%)
Week 9:  [⏳] Testing complete (0%)
Week 10: [⏳] Deployed to production (0%)
Week 11: [⏳] 6 cooperatives live (0%)
Week 12: [⏳] 10 cooperatives live (0%)

Overall Progress: 35% Complete
```

**📊 For detailed progress and metrics:**
- See [Progress Tracking](./progress-tracking.md)
- See [Todolist](./todolist.md) for task breakdown

---

## 🚀 After MVP

**Week 13:** Team rest (3-5 days off)
**Week 14:** Phase 2 planning
**Week 15:** Phase 2 begins

👉 **Next:** [Phase 2 Implementation Guide](../phase-2/implementation-guide.md)

---

**Ready to start?**

**Open [implementation-guide.md](implementation-guide.md) and begin Week 1! 🚀**

---

**Remember:**

> "Perfect is the enemy of shipped."
>
> Focus on 8 core features. Defer everything else.
>
> Your competition is paper books, not enterprise software.
>
> **Ship in 12 weeks. No excuses.** 💪
