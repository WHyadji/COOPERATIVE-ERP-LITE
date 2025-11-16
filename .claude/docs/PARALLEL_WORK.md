# Parallel Work Tracking

**Last Updated**: 2025-11-06
**Branch**: `feature/mvp-implementation`
**Active Sessions**: 0 / 3 Max Recommended

---

## 🎯 Purpose

This file tracks parallel Claude Code sessions to prevent conflicts and coordinate work across multiple terminals. Update this file before starting work and after completing tasks.

---

## 📋 Current Active Sessions

### Session 1: [AVAILABLE]
- **Status**: 🟢 Available
- **Working on**: N/A
- **Issue**: N/A
- **Files Being Modified**: N/A
- **Started**: N/A
- **ETA**: N/A
- **Last Update**: N/A

### Session 2: [AVAILABLE]
- **Status**: 🟢 Available
- **Working on**: N/A
- **Issue**: N/A
- **Files Being Modified**: N/A
- **Started**: N/A
- **ETA**: N/A
- **Last Update**: N/A

### Session 3: [AVAILABLE]
- **Status**: 🟢 Available
- **Working on**: N/A
- **Issue**: N/A
- **Files Being Modified**: N/A
- **Started**: N/A
- **ETA**: N/A
- **Last Update**: N/A

---

## 🔒 File Locks

Track files currently being modified to prevent conflicts:

| File Path | Session | Issue | Lock Time | Expected Release |
|-----------|---------|-------|-----------|------------------|
| *(No files locked)* | - | - | - | - |

---

## ✅ Completed Today

| Time | Session | Issue | Description | Commit |
|------|---------|-------|-------------|--------|
| *(No completions yet)* | - | - | - | - |

---

## 📅 Work Queue

### Ready to Start (No Dependencies)
- [ ] Issue #1: Ledger Entry API (Backend)
- [ ] Issue #2: Journal Detail API (Backend)
- [ ] Issue #3: CoA API Enhancement (Backend)

### Blocked (Waiting for Dependencies)
- [ ] Issue #4: TypeScript Types (Needs Issue #1, #2, #3)
- [ ] Issue #5: API Service Layer (Needs Issue #4)
- [ ] Issue #6: React Hooks (Needs Issue #5)

### In Review
*(No issues in review)*

---

## 🚨 Conflicts & Issues

*(No conflicts reported)*

---

## 📊 Recommended Parallel Groups

### Phase 1: API Layer Foundation (Week 1-2)
**Can run in parallel:**
```
Session 1: Issue #1 - Ledger Entry API
Session 2: Issue #2 - Journal Detail API
Session 3: Issue #3 - CoA API Enhancement
```
**Why safe**: Different handlers, minimal file overlap

**Potential conflicts**:
- `internal/domain/journal_types.go` (coordinate changes)
- `internal/interfaces/repositories.go` (coordinate interface additions)

---

### Phase 2: Frontend Type System (Week 2-3)
**Can run in parallel (after Phase 1):**
```
Session 1: Issue #4 - TypeScript Types
Session 2: Issue #5 - API Service Layer
Session 3: Issue #6 - React Hooks
```
**Why safe**: Different layers of frontend stack

**Potential conflicts**:
- Import chains (coordinate type definitions first)

---

### Phase 3: Core Accounting UI (Week 4-5)
**Can run in parallel (after Phase 2):**
```
Session 1: Issue #7 - Ledger Dialog
Session 2: Issue #8 - Journal Drill-Down
Session 3: Issue #9 - CoA Integration
```
**Why safe**: Separate UI components

**Potential conflicts**:
- Shared components in `frontend/components/`
- Shared hooks (use file locks)

---

## 🔄 How to Use This File

### Before Starting Work:

1. **Pull latest changes:**
   ```bash
   git checkout feature/mvp-implementation
   git pull origin feature/mvp-implementation
   ```

2. **Read this file** to see what others are working on

3. **Check file locks** to avoid conflicts

4. **Update your session status:**
   ```markdown
   ### Session 1: [IN USE - Your Name]
   - **Status**: 🔴 In Progress
   - **Working on**: Implement Ledger Entry API
   - **Issue**: #1
   - **Files Being Modified**:
     - services/accounting-ledger/internal/handlers/ledger_handler.go
     - services/accounting-ledger/internal/service/ledger_service.go
     - services/accounting-ledger/internal/repository/ledger_repository_pg.go
   - **Started**: 2025-11-06 10:00 AM
   - **ETA**: 3 hours
   - **Last Update**: 2025-11-06 10:00 AM
   ```

5. **Add file locks:**
   ```markdown
   | services/accounting-ledger/internal/handlers/ledger_handler.go | Session 1 | #1 | 10:00 AM | 1:00 PM |
   ```

6. **Commit this file:**
   ```bash
   git add .claude/PARALLEL_WORK.md
   git commit -m "chore: update parallel work tracking - starting Issue #1"
   git push
   ```

### While Working:

1. **Update ETA** if timeline changes
2. **Commit frequently** to avoid conflicts
3. **Add notes** if you discover new dependencies

### After Completing Work:

1. **Mark session as available:**
   ```markdown
   ### Session 1: [AVAILABLE]
   - **Status**: 🟢 Available
   ```

2. **Release file locks** (remove from table)

3. **Add to completed list:**
   ```markdown
   | 1:00 PM | Session 1 | #1 | Implemented Ledger Entry API | abc123f |
   ```

4. **Update work queue** (check off completed items)

5. **Commit and push:**
   ```bash
   git add .claude/PARALLEL_WORK.md
   git commit -m "chore: update parallel work tracking - completed Issue #1"
   git push
   ```

---

## 🎯 Quick Status Emojis

- 🟢 Available
- 🔴 In Progress
- 🟡 Testing
- 🔵 In Review
- ⚪ Blocked
- ⛔ Conflict Detected

---

## 📞 Communication Protocol

### If You Need a File That's Locked:

1. **Check ETA** in the session status
2. **Add a note** in "Conflicts & Issues" section:
   ```markdown
   ## 🚨 Conflicts & Issues
   - **Session 2 needs** `ledger_handler.go` (locked by Session 1)
   - **Requested at**: 11:30 AM
   - **Urgency**: Can wait / Need ASAP
   ```
3. **Coordinate** via your preferred method (Slack, email, etc.)

### If You Discover a Dependency:

1. **Update the Work Queue** to mark items as blocked
2. **Add note** explaining the dependency
3. **Notify** other sessions that may be affected

---

## 💡 Best Practices

### ✅ Do:
- Update this file before starting work
- Commit and push updates immediately
- Use descriptive commit messages with issue numbers
- Pull frequently to get others' changes
- Communicate early about potential conflicts

### ❌ Don't:
- Modify locked files without coordination
- Leave sessions marked "In Progress" if you're not working
- Skip updating this file (defeats the purpose!)
- Work on dependent issues in parallel
- Ignore conflicts hoping they'll resolve themselves

---

## 📚 Related Documentation

- **Issue-Driven Development**: `.github/issues/README.md`
- **Working with Issues**: `.github/issues/WORKING_WITH_ISSUES.md`
- **Parallel Work Guide**: `.claude/HOW_TO_PARALLEL_WORK.md`
- **Implementation Plan**: `IMPLEMENTATION_PLAN.md`

---

## 🔮 Future Enhancements

When the team grows, consider:
- Using GitHub Projects for visual tracking
- Setting up Slack notifications for file locks
- Automated conflict detection tools
- Session time limits to prevent long locks

---

**Remember**: Communication is key! This file is only useful if everyone keeps it updated. When in doubt, over-communicate rather than under-communicate.
