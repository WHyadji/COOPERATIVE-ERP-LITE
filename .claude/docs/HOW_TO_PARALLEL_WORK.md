# Guide: Working with Multiple Claude Code Sessions in Parallel

**Version**: 1.0
**Last Updated**: 2025-11-06
**Applies To**: All developers using Claude Code for parallel development

---

## 📖 Table of Contents

1. [Introduction](#introduction)
2. [Why Parallel Work?](#why-parallel-work)
3. [Prerequisites](#prerequisites)
4. [Setup Guide](#setup-guide)
5. [Workflow: Step-by-Step](#workflow-step-by-step)
6. [Coordination Strategies](#coordination-strategies)
7. [Conflict Resolution](#conflict-resolution)
8. [Best Practices](#best-practices)
9. [Common Pitfalls](#common-pitfalls)
10. [Advanced Techniques](#advanced-techniques)
11. [Troubleshooting](#troubleshooting)
12. [Reference](#reference)

---

## Introduction

This guide explains how to run multiple Claude Code sessions simultaneously to accelerate development. Each terminal runs an independent Claude Code instance working on different issues in parallel.

### What You'll Learn:
- How to set up multiple Claude Code terminals
- Coordination strategies to prevent conflicts
- Best practices from Claude Code documentation
- Real-world workflow examples

---

## Why Parallel Work?

### Benefits:
✅ **Faster Development**: Work on 2-3 issues simultaneously
✅ **Efficient Time Use**: While one agent runs tests, another can implement features
✅ **Logical Separation**: Backend and frontend work independently
✅ **Better Focus**: Each agent handles one concern at a time

### When to Use:
- Independent features (e.g., different API endpoints)
- Separate layers (backend API + frontend UI)
- Different domains (authentication + reporting)
- Documentation + implementation

### When NOT to Use:
- Same file modifications
- Tightly coupled features
- Sequential dependencies (A must finish before B starts)
- Complex refactoring affecting many files

---

## Prerequisites

### Required:
- [ ] Claude Code CLI installed and authenticated
- [ ] Git repository with clean working directory
- [ ] Understanding of Issue-Driven Development (see `.github/issues/README.md`)
- [ ] Familiarity with basic Git operations

### Recommended:
- [ ] Terminal multiplexer (tmux, iTerm2 tabs, or separate terminal windows)
- [ ] Git GUI tool for visualizing parallel branches (optional)
- [ ] Coordination file (`.claude/PARALLEL_WORK.md`) set up

---

## Setup Guide

### Step 1: Prepare Your Environment

#### Option A: Multiple Terminal Windows (Easiest)
```bash
# macOS: Use Command+N to open new Terminal windows
# Linux: Use Ctrl+Shift+N
# Windows: Use Ctrl+Shift+N in Windows Terminal

# Each window will run one Claude Code session
```

#### Option B: Terminal Tabs
```bash
# macOS Terminal: Command+T for new tab
# iTerm2: Command+T for new tab
# Windows Terminal: Ctrl+Shift+T for new tab
```

#### Option C: tmux (Advanced)
```bash
# Install tmux
brew install tmux  # macOS
sudo apt install tmux  # Linux

# Create tmux session with 3 panes
tmux new-session \; \
  split-window -h \; \
  split-window -v \; \
  select-pane -t 0

# Navigate panes: Ctrl+b then arrow keys
```

### Step 2: Initialize Each Terminal

**In Terminal 1:**
```bash
cd /Users/adji/Documents/VISI-DIGITAL-TERPADU/accounting-webapp
git checkout feature/mvp-implementation
git pull origin feature/mvp-implementation
git status  # Ensure clean working directory
claude
```

**In Terminal 2:**
```bash
cd /Users/adji/Documents/VISI-DIGITAL-TERPADU/accounting-webapp
git checkout feature/mvp-implementation
git pull origin feature/mvp-implementation
git status
claude
```

**In Terminal 3:**
```bash
cd /Users/adji/Documents/VISI-DIGITAL-TERPADU/accounting-webapp
git checkout feature/mvp-implementation
git pull origin feature/mvp-implementation
git status
claude
```

### Step 3: Set Up Coordination File

**First time setup:**
```bash
# In any terminal (before starting Claude)
git pull origin feature/mvp-implementation
# Check .claude/PARALLEL_WORK.md exists
cat .claude/PARALLEL_WORK.md
```

---

## Workflow: Step-by-Step

### Before Starting Work in Any Terminal

#### 1. Check Coordination File
```bash
# In Terminal 1 (before starting Claude)
cat .claude/PARALLEL_WORK.md | grep "Session"
```

Look for:
- Which sessions are in use?
- Which files are locked?
- Are there any conflicts?

#### 2. Claim Your Session

**Update `.claude/PARALLEL_WORK.md`:**
```markdown
### Session 1: [IN USE - Alice]
- **Status**: 🔴 In Progress
- **Working on**: Implement Ledger Entry API
- **Issue**: #1
- **Files Being Modified**:
  - services/accounting-ledger/internal/handlers/ledger_handler.go
  - services/accounting-ledger/internal/service/ledger_service.go
  - services/accounting-ledger/internal/repository/ledger_repository_pg.go
- **Started**: 2025-11-06 09:00 AM
- **ETA**: 3 hours (12:00 PM)
- **Last Update**: 2025-11-06 09:00 AM
```

#### 3. Add File Locks

**Update the File Locks table:**
```markdown
| File Path | Session | Issue | Lock Time | Expected Release |
|-----------|---------|-------|-----------|------------------|
| services/accounting-ledger/internal/handlers/ledger_handler.go | Session 1 | #1 | 09:00 AM | 12:00 PM |
| services/accounting-ledger/internal/service/ledger_service.go | Session 1 | #1 | 09:00 AM | 12:00 PM |
```

#### 4. Commit Coordination File
```bash
git add .claude/PARALLEL_WORK.md
git commit -m "chore: session 1 starting work on Issue #1"
git push origin feature/mvp-implementation
```

#### 5. Start Claude Code
```bash
claude
# > "Work on Issue #1: Implement Ledger Entry API.
#    Follow the task specification in .github/issues/phase-1-task-1.1-ledger-entry-api.md"
```

---

### While Working

#### Every 30-60 Minutes:

**1. Pull Changes from Other Sessions:**
```bash
# In Claude: "Can you pull the latest changes from remote?"
# Or run in separate terminal:
cd /Users/adji/Documents/VISI-DIGITAL-TERPADU/accounting-webapp
git pull origin feature/mvp-implementation
```

**2. Commit Your Progress:**
```bash
# Claude automatically commits, but ensure regular commits:
# "Please commit the current progress with message:
#  'feat: implement ledger handler basic structure (#1)'"
```

**3. Update ETA if Needed:**
```markdown
- **ETA**: 4 hours (1:00 PM) - Extended due to additional validation logic
- **Last Update**: 2025-11-06 10:30 AM
```

**4. Watch for Conflicts:**
```bash
# Check if others modified your files
git status
git diff origin/feature/mvp-implementation
```

---

### After Completing Work

#### 1. Final Commit and Push
```bash
# In Claude: "Please commit all changes with message:
#  'feat: complete ledger entry API implementation (#1)
#
#  Implemented:
#  - GET /api/v1/accounts/:accountId/journal-entries endpoint
#  - Ledger service with balance calculation
#  - Repository layer with PostgreSQL implementation
#  - Unit and integration tests
#  - API documentation
#
#  Closes #1'"

# Ensure pushed:
git push origin feature/mvp-implementation
```

#### 2. Release Session and File Locks

**Update `.claude/PARALLEL_WORK.md`:**
```markdown
### Session 1: [AVAILABLE]
- **Status**: 🟢 Available
- **Working on**: N/A
- **Issue**: N/A
- **Files Being Modified**: N/A
- **Started**: N/A
- **ETA**: N/A
- **Last Update**: N/A
```

**Remove from File Locks table:**
```markdown
| File Path | Session | Issue | Lock Time | Expected Release |
|-----------|---------|-------|-----------|------------------|
| *(No files locked)* | - | - | - | - |
```

**Add to Completed Today:**
```markdown
| Time | Session | Issue | Description | Commit |
|------|---------|-------|-------------|--------|
| 12:30 PM | Session 1 | #1 | Implemented Ledger Entry API | abc123f |
```

#### 3. Commit Coordination File
```bash
git add .claude/PARALLEL_WORK.md
git commit -m "chore: session 1 completed Issue #1"
git push origin feature/mvp-implementation
```

#### 4. Exit Claude
```bash
# In Claude: "Thank you! I'm done with this issue. You can exit."
# Or press Ctrl+D to exit
```

---

## Coordination Strategies

### Strategy 1: Layer-Based Separation (Recommended for MVP)

**Separate by architectural layers:**

```
Terminal 1 (Backend API):
├── handlers/
├── services/
└── repositories/

Terminal 2 (Frontend Services):
├── types/
├── lib/api/
└── hooks/

Terminal 3 (UI Components):
├── components/
└── app/
```

**Example parallel work:**
```markdown
Session 1: Issue #1 - Ledger Entry API (Backend)
Session 2: Issue #4 - TypeScript Types (Frontend)
Session 3: Issue #19 - User Documentation (Docs)
```

### Strategy 2: Feature-Based Separation

**Separate by business features:**

```
Session 1: Account Management
├── Account API
├── Account Repository
└── Account Tests

Session 2: Journal Entries
├── Journal API
├── Journal Repository
└── Journal Tests

Session 3: Chart of Accounts
├── CoA API
├── CoA Repository
└── CoA Tests
```

### Strategy 3: Phase-Based Separation

**Work on same phase but different tasks:**

```
Phase 1 (Week 1-2):
├── Session 1: Task 1.1 (Ledger API)
├── Session 2: Task 1.2 (Journal API)
└── Session 3: Task 1.3 (CoA API)

Phase 2 (Week 2-3):
├── Session 1: Task 2.1 (TypeScript Types)
├── Session 2: Task 2.2 (API Service Layer)
└── Session 3: Task 2.3 (React Hooks)
```

### Strategy 4: Role-Based Separation

**Assign by developer role:**

```
Backend Developer (Session 1):
└── All Go services, repositories, handlers

Frontend Developer (Session 2):
└── All React components, hooks, API clients

DevOps/Testing (Session 3):
└── CI/CD, Docker, tests, documentation
```

---

## Conflict Resolution

### Types of Conflicts

#### 1. Git Merge Conflicts

**Symptoms:**
```bash
git pull
# CONFLICT (content): Merge conflict in services/accounting-ledger/internal/handlers/ledger_handler.go
```

**Resolution:**
```bash
# Option A: Resolve in Claude
# "There's a merge conflict in ledger_handler.go. Can you help resolve it?"

# Option B: Manual resolution
git status
# Open conflicted file in editor
# Look for conflict markers:
<<<<<<< HEAD
your changes
=======
their changes
>>>>>>> branch-name

# Edit to keep both changes or choose one
# Then:
git add services/accounting-ledger/internal/handlers/ledger_handler.go
git commit -m "fix: resolve merge conflict in ledger handler"
git push
```

#### 2. File Lock Conflicts

**Symptoms:**
- Another session is modifying a file you need

**Resolution:**

**Step 1: Check `.claude/PARALLEL_WORK.md`**
```markdown
| services/accounting-ledger/internal/domain/journal_types.go | Session 2 | #2 | 10:00 AM | 11:30 AM |
```

**Step 2: Add Conflict Note**
```markdown
## 🚨 Conflicts & Issues

- **Session 1 needs** `journal_types.go` (locked by Session 2)
- **Requested at**: 10:30 AM
- **Urgency**: Can wait until 11:30 AM
- **Reason**: Need to add new field to JournalEntry struct
```

**Step 3: Coordinate**
- If ETA is soon (< 30 min): Wait
- If urgent: Contact other developer to coordinate
- If possible: Work on different part of issue temporarily

#### 3. Dependency Conflicts

**Symptoms:**
- Issue #2 depends on Issue #1, but both running in parallel

**Resolution:**

**Step 1: Identify the dependency**
```markdown
Issue #5 (API Service Layer) depends on:
- Issue #4 (TypeScript Types) - MUST be completed first
```

**Step 2: Update Work Queue**
```markdown
### Blocked (Waiting for Dependencies)
- [ ] Issue #5: API Service Layer (Needs Issue #4)
```

**Step 3: Work on something else**
```bash
# In Session 2:
# "Let's work on Issue #6 (React Hooks) instead while waiting for Issue #4"
```

#### 4. Test Failures Due to Parallel Changes

**Symptoms:**
```bash
npm run test
# FAIL: API integration tests failed
# Expected endpoint: /api/v1/accounts/:id/entries
# But found: /api/v1/ledger/:accountId/entries
```

**Resolution:**

**Step 1: Pull latest changes**
```bash
git pull origin feature/mvp-implementation
```

**Step 2: Check what changed**
```bash
git log --oneline -5
git diff HEAD~1 HEAD
```

**Step 3: Update your code to match**
```bash
# In Claude: "The endpoint path was changed by another session.
#  Can you update our API client to use /api/v1/ledger/:accountId/entries?"
```

---

## Best Practices

### Communication Best Practices

#### 1. Update Coordination File Religiously
```bash
# Before starting: Update session status
# Every hour: Update ETA if changed
# After completing: Release locks and mark available
```

#### 2. Commit Frequently
```bash
# Good: Commit every 30-60 minutes
git commit -m "feat: add ledger handler basic structure (#1)"
git commit -m "feat: implement balance calculation logic (#1)"
git commit -m "test: add unit tests for ledger service (#1)"

# Bad: One huge commit at the end
git commit -m "feat: complete everything (#1)"
```

#### 3. Use Descriptive Commit Messages
```bash
# Good:
git commit -m "feat: implement GetLedgerEntriesByAccount handler (#1)

Adds new handler for retrieving ledger entries with:
- Pagination support (limit/offset)
- Date range filtering
- Running balance calculation
- Multi-tenant safety checks

Related to Issue #1"

# Bad:
git commit -m "updates"
```

#### 4. Pull Before Push
```bash
# Always pull first to catch conflicts early
git pull origin feature/mvp-implementation
# Then push
git push origin feature/mvp-implementation
```

### Code Organization Best Practices

#### 1. Minimize Shared Files

**Good: Separate files**
```
services/accounting-ledger/internal/handlers/
├── ledger_handler.go      # Session 1
├── journal_handler.go     # Session 2
└── account_handler.go     # Session 3
```

**Bad: Single large file**
```
services/accounting-ledger/internal/handlers/
└── handlers.go            # All sessions conflict here!
```

#### 2. Coordinate Interface Changes

**If multiple sessions need to modify shared interfaces:**

```go
// services/accounting-ledger/internal/interfaces/repositories.go

// Session 1 needs to add:
type LedgerRepository interface {
    GetLedgerEntries(ctx context.Context, accountID uuid.UUID) ([]*LedgerEntry, error)
}

// Session 2 needs to add:
type JournalRepository interface {
    GetJournalDetails(ctx context.Context, journalID uuid.UUID) (*JournalDetail, error)
}
```

**Solution:**
- One session adds both interfaces
- Other session waits or works on implementation
- Coordinate via `.claude/PARALLEL_WORK.md`

#### 3. Use Feature Flags for Incomplete Work

```go
// If Session 1 needs Session 2's work that's not ready:

func GetLedgerEntries(ctx context.Context, accountID uuid.UUID) ([]*LedgerEntry, error) {
    // TODO: Integration with journal detail API (Issue #2)
    // Placeholder for now
    if config.FeatureFlags.JournalDetailEnabled {
        // Will be implemented by Session 2
        details, err := journalDetailService.GetDetails(ctx, accountID)
        if err != nil {
            return nil, err
        }
        return enrichWithDetails(entries, details), nil
    }

    // Fallback: basic implementation
    return entries, nil
}
```

### Testing Best Practices

#### 1. Run Tests Before Pushing

```bash
# Backend tests
cd services/accounting-ledger
go test ./... -v

# Frontend tests
cd frontend
npm run test

# If all pass:
git push origin feature/mvp-implementation
```

#### 2. Test in Isolation

```bash
# Don't run services that others are testing
# Use different database/ports if needed

# Terminal 1: Testing backend
cd services/accounting-ledger
go test ./internal/handlers/... -v

# Terminal 2: Testing frontend
cd frontend
npm run test -- --testPathPattern=ledger

# Terminal 3: E2E tests (when services stable)
npm run test:e2e
```

#### 3. Coordinate Integration Tests

```markdown
## 🚨 Integration Test Schedule

- **10:00 AM**: Session 1 runs backend integration tests
- **11:00 AM**: Session 2 runs frontend integration tests
- **2:00 PM**: Session 3 runs full E2E test suite

Avoid running integration tests simultaneously to prevent port conflicts.
```

---

## Common Pitfalls

### Pitfall 1: Forgetting to Pull

**Problem:**
```bash
# Session 1 pushes changes
git push origin feature/mvp-implementation

# Session 2 tries to push (without pulling first)
git push origin feature/mvp-implementation
# ! [rejected] feature/mvp-implementation -> feature/mvp-implementation (non-fast-forward)
```

**Solution:**
```bash
# Always pull before pushing
git pull origin feature/mvp-implementation
# Resolve any conflicts
git push origin feature/mvp-implementation
```

**Prevention:**
```bash
# Set up Git alias for pull-then-push
git config alias.pushup '!git pull origin $(git branch --show-current) && git push origin $(git branch --show-current)'

# Now use:
git pushup
```

### Pitfall 2: Working on Dependent Issues Simultaneously

**Problem:**
```
Session 1: Issue #4 (TypeScript Types) - In Progress
Session 2: Issue #5 (API Service Layer) - In Progress, but DEPENDS on #4
```

**Result:**
- Session 2 can't compile because types don't exist yet
- Session 2 wastes time waiting or using placeholder types

**Solution:**
```markdown
# In .claude/PARALLEL_WORK.md

### Work Queue
#### Ready to Start (No Dependencies)
- [ ] Issue #1: Ledger Entry API
- [ ] Issue #2: Journal Detail API
- [x] Issue #4: TypeScript Types (Session 1 - In Progress)

#### Blocked (Waiting for Dependencies)
- [ ] Issue #5: API Service Layer (Needs Issue #4) - DO NOT START YET
```

**Prevention:**
- Always check issue dependencies before starting
- Refer to IMPLEMENTATION_PLAN.md for dependency graph

### Pitfall 3: Not Updating Coordination File

**Problem:**
```
Session 1: Working on ledger_handler.go (not documented)
Session 2: Also starts working on ledger_handler.go
Result: Massive merge conflict
```

**Solution:**
- ALWAYS update `.claude/PARALLEL_WORK.md` before starting
- Commit the coordination file immediately
- Check it before starting any work

**Prevention:**
```bash
# Create Git hook to remind you
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
if ! git diff --cached --name-only | grep -q ".claude/PARALLEL_WORK.md"; then
    echo "⚠️  Warning: .claude/PARALLEL_WORK.md not updated. Did you document your work?"
    echo "Press Enter to continue or Ctrl+C to cancel"
    read
fi
EOF

chmod +x .git/hooks/pre-commit
```

### Pitfall 4: Long-Running Locks

**Problem:**
```markdown
| services/accounting-ledger/internal/domain/types.go | Session 1 | #1 | 9:00 AM | 5:00 PM |
# 8 hours! Everyone else blocked
```

**Solution:**
- Break work into smaller chunks
- Release locks during lunch/breaks
- Commit partial work and push

**Prevention:**
```markdown
# Set maximum lock time: 2-3 hours
| File | Session | Issue | Lock Time | Expected Release | Max Lock |
|------|---------|-------|-----------|------------------|----------|
| types.go | Session 1 | #1 | 9:00 AM | 11:00 AM | 2h |
```

### Pitfall 5: Ignoring Conflicts

**Problem:**
```bash
git pull origin feature/mvp-implementation
# CONFLICT (content): Merge conflict in handler.go

# Developer marks resolved without actually fixing:
git add handler.go
git commit -m "merge"
git push

# Code is now broken with conflict markers:
<<<<<<< HEAD
func HandleRequest() {
=======
func HandleLedgerRequest() {
>>>>>>> feature/mvp-implementation
```

**Solution:**
- **NEVER** commit files with conflict markers
- Carefully review and test merged code
- Ask Claude for help: "Can you help resolve this merge conflict?"

**Prevention:**
```bash
# Add Git hook to prevent committing conflicts
cat > .git/hooks/pre-commit << 'EOF'
#!/bin/bash
if git diff --cached | grep -q "^<<<<<<< HEAD"; then
    echo "❌ Error: Conflict markers found in staged files!"
    echo "Please resolve conflicts before committing."
    exit 1
fi
EOF

chmod +x .git/hooks/pre-commit
```

---

## Advanced Techniques

### Technique 1: Stacked PRs (For Complex Features)

When a feature is too large for parallel work on single branch:

```bash
# Main branch
feature/mvp-implementation

# Stack 1: Backend API (based on mvp-implementation)
feature/issue-1-ledger-api

# Stack 2: Frontend Types (based on issue-1-ledger-api)
feature/issue-4-typescript-types

# Stack 3: UI Components (based on issue-4-typescript-types)
feature/issue-7-ledger-dialog
```

**Workflow:**
```bash
# Session 1: Create backend API branch
git checkout feature/mvp-implementation
git checkout -b feature/issue-1-ledger-api
# Work on Issue #1
git push origin feature/issue-1-ledger-api

# Session 2: Wait for Session 1, then create types branch
git checkout feature/issue-1-ledger-api  # Base on Session 1's work!
git checkout -b feature/issue-4-typescript-types
# Work on Issue #4
git push origin feature/issue-4-typescript-types

# Later: Merge in order
# 1. Merge issue-1 → mvp-implementation
# 2. Rebase issue-4 on mvp-implementation
# 3. Merge issue-4 → mvp-implementation
```

### Technique 2: Shared Stubs During Development

When two sessions need each other's work:

```go
// Session 1 creates stub that Session 2 can use:

// File: internal/interfaces/stubs/journal_service_stub.go
package stubs

// JournalServiceStub provides placeholder implementation for testing
// TODO: Replace with real implementation from Issue #2
type JournalServiceStub struct{}

func (s *JournalServiceStub) GetJournalDetails(ctx context.Context, id uuid.UUID) (*JournalDetail, error) {
    return &JournalDetail{
        ID:          id,
        Reference:   "STUB-001",
        Description: "Stub journal entry for testing",
    }, nil
}

// Session 2 can now use this stub while implementing their feature
// Later, Session 2 replaces stub with real implementation
```

### Technique 3: Real-Time Coordination with Watch

```bash
# Terminal 4: Dedicated coordination terminal
watch -n 5 'git pull && cat .claude/PARALLEL_WORK.md | grep -A 5 "Session"'

# Updates every 5 seconds showing:
# - Latest coordination file
# - Active sessions
# - File locks
```

### Technique 4: Parallel Testing Matrix

```yaml
# .github/workflows/parallel-test.yml
name: Parallel Test Matrix

on: [push]

jobs:
  test-matrix:
    strategy:
      matrix:
        test-suite:
          - backend-unit
          - backend-integration
          - frontend-unit
          - frontend-e2e
          - api-contract

    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Run ${{ matrix.test-suite }}
        run: npm run test:${{ matrix.test-suite }}
```

---

## Troubleshooting

### Issue: "Git push rejected (non-fast-forward)"

**Cause:** Someone else pushed changes while you were working

**Solution:**
```bash
# Pull and rebase
git pull --rebase origin feature/mvp-implementation

# If conflicts, resolve them:
git status
# Edit conflicted files
git add .
git rebase --continue

# Then push
git push origin feature/mvp-implementation
```

### Issue: "Port already in use"

**Cause:** Multiple sessions trying to run services on same port

**Solution:**
```bash
# Session 1: Use default ports
npm run dev  # Port 3000

# Session 2: Use different ports
PORT=3001 npm run dev

# Session 3: Use different ports
PORT=3002 npm run dev

# Backend services:
# Session 1: accounting-ledger on 8080
# Session 2: api-gateway on 8081
# Session 3: business-service on 8082
```

### Issue: "Database connection pool exhausted"

**Cause:** Multiple sessions using same database simultaneously

**Solution:**
```bash
# Use separate test databases
# Session 1:
DATABASE_URL=postgres://localhost:54321/accounting_test_1 go test ./...

# Session 2:
DATABASE_URL=postgres://localhost:54321/accounting_test_2 go test ./...

# Or use connection pooling limits:
# config.yaml
database:
  max_connections: 20  # Divide among sessions
```

### Issue: "Claude session lost context"

**Cause:** Session was idle too long or context limit reached

**Solution:**
```bash
# Restart Claude with context:
# 1. Commit current work
git add .
git commit -m "wip: save current progress on Issue #1"

# 2. Exit and restart Claude
exit
claude

# 3. Provide context:
# > "I was working on Issue #1 (Ledger Entry API).
#    I've completed the handler and service layer.
#    Next I need to implement the repository layer.
#    Please read .github/issues/phase-1-task-1.1-ledger-entry-api.md
#    and continue from step 1.1.3"
```

### Issue: "Coordination file merge conflicts"

**Cause:** Multiple sessions updated `.claude/PARALLEL_WORK.md` simultaneously

**Solution:**
```bash
git pull origin feature/mvp-implementation
# CONFLICT in .claude/PARALLEL_WORK.md

# Open file and resolve:
# Keep both session updates (they shouldn't conflict logically)

### Session 1: [IN USE - Alice]
# Status: 🔴 In Progress
...

### Session 2: [IN USE - Bob]
# Status: 🔴 In Progress
...

# Stage and commit:
git add .claude/PARALLEL_WORK.md
git commit -m "fix: resolve parallel work coordination merge"
git push
```

---

## Reference

### Quick Command Cheat Sheet

```bash
# Start parallel work
git checkout feature/mvp-implementation
git pull origin feature/mvp-implementation
# Update .claude/PARALLEL_WORK.md
git add .claude/PARALLEL_WORK.md && git commit -m "chore: starting Issue #X" && git push
claude

# During work
git pull origin feature/mvp-implementation  # Every 30-60 min
git add . && git commit -m "feat: progress on Issue #X (#X)" && git push  # Frequently

# Finish work
# Update .claude/PARALLEL_WORK.md (mark available, remove locks, add to completed)
git add .claude/PARALLEL_WORK.md && git commit -m "chore: completed Issue #X" && git push
exit  # Exit Claude

# Check coordination
cat .claude/PARALLEL_WORK.md | grep "Session"
cat .claude/PARALLEL_WORK.md | grep "File Locks" -A 10
```

### File Locations Reference

| File | Purpose | When to Update |
|------|---------|----------------|
| `.claude/PARALLEL_WORK.md` | Track active sessions and file locks | Before/after work, every hour |
| `.github/issues/*.md` | Issue templates | Read before starting issue |
| `IMPLEMENTATION_PLAN.md` | Overall MVP plan | Reference for dependencies |
| `.github/issues/README.md` | IDD methodology | Reference for issue workflow |
| `.github/issues/WORKING_WITH_ISSUES.md` | Daily workflow guide | Reference for commits/PRs |

### Recommended Parallel Groups by Phase

#### Phase 1 (Week 1-2): API Layer Foundation
```
✅ Safe to parallelize:
- Session 1: Issue #1 (Ledger Entry API)
- Session 2: Issue #2 (Journal Detail API)
- Session 3: Issue #3 (CoA API Enhancement)

⚠️ Coordination needed:
- internal/domain/journal_types.go (coordinate additions)
- internal/interfaces/repositories.go (coordinate interface additions)
```

#### Phase 2 (Week 2-3): Frontend Type System
```
✅ Safe to parallelize (after Phase 1):
- Session 1: Issue #4 (TypeScript Types)
- Session 2: Issue #5 (API Service Layer) - Wait for #4
- Session 3: Issue #6 (React Hooks) - Wait for #5

⚠️ Sequential dependency:
Must complete in order: #4 → #5 → #6
Alternative: Run #4, then parallelize #5 and #6 if possible
```

#### Phase 3 (Week 4-5): Core Accounting UI
```
✅ Safe to parallelize (after Phase 2):
- Session 1: Issue #7 (Ledger Dialog)
- Session 2: Issue #8 (Journal Drill-Down)
- Session 3: Issue #9 (CoA Integration)

⚠️ Coordination needed:
- Shared components in frontend/components/
- Shared hooks in frontend/hooks/
```

#### Phases 4-8: Can mix and match
```
Backend-focused:
- Session 1: Backend feature
- Session 2: Frontend feature
- Session 3: Testing/Documentation

Feature-focused:
- Session 1: Feature A (full stack)
- Session 2: Feature B (full stack)
- Session 3: Feature C (full stack)
```

### Status Emoji Legend

| Emoji | Meaning | When to Use |
|-------|---------|-------------|
| 🟢 | Available | Session is free |
| 🔴 | In Progress | Actively working |
| 🟡 | Testing | Running tests |
| 🔵 | In Review | Waiting for code review |
| ⚪ | Blocked | Waiting for dependency |
| ⛔ | Conflict | Merge conflict or file lock conflict |
| ✅ | Completed | Work finished |
| 🚧 | WIP | Work in progress (temporary state) |

---

## Related Documentation

- **Issue-Driven Development Guide**: `.github/issues/README.md`
- **Daily Workflow Guide**: `.github/issues/WORKING_WITH_ISSUES.md`
- **Parallel Work Tracker**: `.claude/PARALLEL_WORK.md`
- **Implementation Plan**: `IMPLEMENTATION_PLAN.md`
- **Project Conventions**: `CLAUDE.md`

---

## Need Help?

### Common Questions:

**Q: How many parallel sessions should I run?**
A: Start with 2, max 3. More than 3 increases coordination overhead.

**Q: Can I work on the same file in different sessions?**
A: Avoid it! Causes merge conflicts. Use file locks in `.claude/PARALLEL_WORK.md`.

**Q: What if I need a file that's locked?**
A: Check ETA in coordination file. If urgent, coordinate with other developer. Otherwise, work on different part of issue.

**Q: Should I use separate branches or same branch?**
A: For this MVP: Same branch (`feature/mvp-implementation`) with good coordination. For larger teams: Consider separate feature branches per issue.

**Q: How do I handle test failures from other sessions' changes?**
A: Pull latest changes, update your code to match, re-run tests. Communicate if you find issues.

**Q: Can I pause a session and resume later?**
A: Yes! Commit your work, update coordination file (mark available or "paused"), push. Resume by pulling and updating status.

---

**Remember**: The key to successful parallel work is **communication**. Update `.claude/PARALLEL_WORK.md` religiously, commit frequently, and pull often. When in doubt, over-communicate!

---

**Version History**:
- **v1.0** (2025-11-06): Initial guide based on Claude Code best practices

**Maintained By**: Development Team
**Last Reviewed**: 2025-11-06
