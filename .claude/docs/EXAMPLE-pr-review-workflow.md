# 📊 Example: PR Review → Solve Sub-Issues Workflow

## Scenario: Your Issue #382 from PR #188

Your exact case: Issue #382 was created from a review of PR #188.

**Issue Title:** `[PR #188] [CRITICAL] Increase test coverage to 85% minimum`

---

## 🔄 Complete Workflow Diagram

```
┌─────────────────────────────────────────────────────────────┐
│ STEP 1: Review PR #188                                      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    /dev:review-pr 188
                              │
                              ▼
        ┌─────────────────────────────────────────┐
        │ Review finds 5 issues                   │
        │ Creates GitHub issues:                  │
        │  • #382 [CRITICAL] Test coverage        │
        │  • #383 [HIGH] SOLID violation          │
        │  • #384 [HIGH] Security issue           │
        │  • #385 [MEDIUM] Code quality           │
        │  • #386 [LOW] Documentation             │
        └─────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ STEP 2: Solve Issue #382                                    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    /dev:solve-issue 382
                              │
                              ▼
        ┌─────────────────────────────────────────┐
        │ Command detects:                        │
        │  • Issue #382 is from PR #188           │
        │  • Parent PR: #188                      │
        │  • Parent PR Branch: feature/xyz        │
        └─────────────────────────────────────────┘
                              │
                              ▼
        ┌─────────────────────────────────────────┐
        │ Creates NEW separate branch:            │
        │  test/issue-382-coverage-a7b3c9f2       │
        │                                         │
        │ NOT on PR #188's branch!                │
        └─────────────────────────────────────────┘
                              │
                              ▼
        ┌─────────────────────────────────────────┐
        │ AI Implementation:                      │
        │  • Analyzes codebase                    │
        │  • Identifies untested code             │
        │  • Writes missing tests                 │
        │  • Runs tests (verifies ≥85%)           │
        │  • Commits with proper message          │
        └─────────────────────────────────────────┘
                              │
                              ▼
        ┌─────────────────────────────────────────┐
        │ Creates NEW PR #400                     │
        │                                         │
        │ PR #400 Description includes:           │
        │  • Closes #382                          │
        │  • Related to PR #188                   │
        │  • Context: Found in review of #188     │
        │  • Suggestion: Merge independently      │
        └─────────────────────────────────────────┘
                              │
                              ▼
        ┌─────────────────────────────────────────┐
        │ Auto-comments created:                  │
        │                                         │
        │ On Issue #382:                          │
        │  "Fixed by PR #400, from review of #188"│
        │                                         │
        │ On PR #188:                             │
        │  "Sub-issue #382 fixed in PR #400"      │
        └─────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ STEP 3: Review & Merge PR #400                              │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
                    /dev:verify-pr 400
                              │
                              ▼
                    /dev:merge-pr 400
                              │
                              ▼
        ┌─────────────────────────────────────────┐
        │ Results:                                │
        │  ✅ PR #400 merged to main              │
        │  ✅ Issue #382 auto-closed              │
        │  ✅ Test coverage now ≥85%              │
        │  ✅ PR #188 can rebase on main          │
        └─────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│ STEP 4: Continue with other issues                          │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
            /dev:solve-issue 383  (Creates PR #401)
            /dev:solve-issue 384  (Creates PR #402)
            /dev:solve-issue 385  (Creates PR #403)
                              │
                              ▼
            Each creates separate PR, merges independently
```

---

## 📝 Actual Commands You Would Run

```bash
# Assuming PR #188 already exists and issue #382 was created from its review

# 1. Solve issue #382
/dev:solve-issue 382

# That's it! The command handles everything:
# ✅ Detects it's from PR #188
# ✅ Creates new branch (not on #188's branch)
# ✅ Implements fix
# ✅ Creates PR #400
# ✅ Links everything
# ✅ Comments on both #382 and #188

# 2. Optionally review the AI implementation
/dev:review-pr 400

# 3. Verify before merge
/dev:verify-pr 400

# 4. Merge when ready
/dev:merge-pr 400
```

---

## 🔗 What Gets Linked Where

### Issue #382
```
Title: [PR #188] [CRITICAL] Increase test coverage to 85% minimum

Body: (original issue description from review)

Comments:
  └─ 🔧 PR Created: #400
     This PR implements the solution for this issue.
     PR: https://github.com/.../pull/400
     Branch: test/issue-382-coverage-a7b3c9f2
     Context: This issue was created from review of PR #188
     Parent PR: https://github.com/.../pull/188

     The issue will automatically close when the PR is merged.

Status: Open → Closed (when PR #400 merges)
```

### PR #400 (New Fix PR)
```
Title: [PR #188] [CRITICAL] Increase test coverage to 85% minimum

Body:
  ## 🎯 Purpose
  This PR resolves issue #382

  Closes #382

  ## 🔗 Related to PR #188
  This issue was identified during review of PR #188:
  - Parent PR: https://github.com/.../pull/188
  - Parent PR Title: Add new tax calculation feature
  - Parent PR Branch: feature/tax-calculation

  This fix addresses one of the issues found in that review and should be:
  1. Reviewed independently
  2. Merged into main separately
  3. The parent PR #188 can then be updated/rebased if needed

  Note: Merging this PR will automatically close issue #382,
        which was created from the review of PR #188.

  ## 📋 Issue Details
  (issue description)

  ## 🔧 Implementation
  (what was implemented)

  (... rest of PR description)

Status: Open → Review → Merged
```

### PR #188 (Original PR Being Reviewed)
```
Title: Add new tax calculation feature

Body: (original PR description)

Comments:
  (... other comments ...)

  └─ ✅ Sub-Issue Fixed: #382
     A fix PR has been created for one of the issues found in this review:

     Issue: #382 - [PR #188] [CRITICAL] Increase test coverage
     Fix PR: #400
     Fix PR URL: https://github.com/.../pull/400

     This fix will be merged independently. Once merged, you may want to:
     - Rebase this PR if needed
     - Re-run the review to check if issue is resolved
     - Continue with other review findings

Status: Open (unchanged - you decide what to do with it)
```

---

## 🎯 Key Benefits of This Approach

### 1. **Independent Merge**
- PR #400 can merge even if PR #188 isn't ready
- Fixes go to main immediately
- No waiting for large PR to be approved

### 2. **Better Git History**
```
main
  ├─ PR #400: Increase test coverage (Issue #382)
  ├─ PR #401: Fix SOLID violation (Issue #383)
  ├─ PR #402: Fix security issue (Issue #384)
  └─ PR #188: Add tax calculation (rebased on above)
```

### 3. **Traceability**
- Easy to see which review finding was addressed
- Clear link from issue → fix PR → parent PR
- Audit trail for compliance

### 4. **Flexibility**
After fixes merge, PR #188 author can:
- **Option A**: Rebase PR #188 on main (includes fixes)
- **Option B**: Close PR #188 if fixes covered everything
- **Option C**: Continue with remaining changes

---

## 💡 Real-World Example

```bash
# Monday: Review finds issues
/dev:review-pr 188
# Creates: #382 (CRITICAL), #383 (HIGH), #384 (HIGH), #385 (MEDIUM), #386 (LOW)

# Monday afternoon: Fix CRITICAL first
/dev:solve-issue 382
# Creates PR #400, merges same day → Issue #382 closed ✅

# Tuesday: Fix HIGH priority issues
/dev:solve-issue 383
# Creates PR #401, merges → Issue #383 closed ✅

/dev:solve-issue 384
# Creates PR #402, merges → Issue #384 closed ✅

# Wednesday: Author of PR #188 rebases
git checkout feature/tax-calculation
git rebase main  # Now includes all fixes!

# Thursday: Fix remaining issues
/dev:solve-issue 385
/dev:solve-issue 386

# Friday: All issues fixed, PR #188 can merge or close
```

---

## 🆚 Alternative Approach (Not Recommended)

**What if we modified PR #188 directly?**

❌ **Problems:**
- Can't merge fixes until entire PR #188 is ready
- Mixing concerns (original feature + fixes)
- Harder to review (what's new vs fixes?)
- Complicated Git history
- If PR #188 is abandoned, fixes are lost

✅ **Our approach:**
- Fixes merge immediately
- Each change reviewed independently
- Clean Git history
- Fixes survive even if PR #188 closes
- PR #188 author gets fixes via rebase

---

## 📋 Summary

**Your exact workflow:**

```bash
# Your issue #382 is from PR #188 review
/dev:solve-issue 382

# Command automatically:
# 1. Detects parent PR: #188
# 2. Creates separate branch: test/issue-382-coverage-a7b3c9f2
# 3. Implements test coverage improvements
# 4. Creates NEW PR #400 (not modifying #188)
# 5. PR #400 references PR #188
# 6. Comments on both #382 and #188
# 7. Issue #382 closes when PR #400 merges
# 8. PR #188 is notified and can rebase
```

**Time:** 5-15 minutes per issue
**Result:** Clean, traceable, independent fixes

---

🤖 **Powered by Claude Code**

This workflow keeps your Git history clean and issues traceable! 🚀
