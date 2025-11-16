---
allowed-tools: Read
description: Quick reference guide for GitHub issue resolution workflows
argument-hint:
---

# 🚀 GitHub Issue Resolution Workflow Guide

This guide shows the different workflows available for solving GitHub issues efficiently.

## Quick Command Reference

```bash
# 🔍 DISCOVERY: Find issues
/dev:identify-issues [path] --create-issues

# 🛠️ SOLVE: Automated end-to-end solution
/dev:solve-issue <issue-number>

# 🔧 MANUAL: Create branch + PR, implement yourself
/dev:create-fix-pr <issue-numbers>

# 👀 REVIEW: Deep PR analysis
/dev:review-pr <pr-number>

# ✅ VERIFY: Pre-merge checks
/dev:verify-pr <pr-number>

# 🚀 MERGE: Safe merge with cleanup
/dev:merge-pr <pr-number>
```

---

## Workflow Options

### Option 1: 🤖 Fully Automated (Recommended)

**Use when:** You want AI to implement the entire solution

```bash
# One command does it all!
/dev:solve-issue 382

# What it does:
# 1. ✅ Fetches issue #382 from GitHub
# 2. ✅ Creates unique branch: fix/issue-382-description-a7b3c9f2
# 3. ✅ AI implements the solution (following SOLID principles)
# 4. ✅ Runs tests, linting, build
# 5. ✅ Commits with proper message
# 6. ✅ Pushes to remote
# 7. ✅ Creates PR linked to issue
# 8. ✅ Issue auto-closes when PR merges
```

**Variants:**
```bash
# Create draft PR (review before marking ready)
/dev:solve-issue 382 --draft-only

# Auto-merge when CI passes
/dev:solve-issue 382 --auto-merge

# Use GitHub URL instead
/dev:solve-issue https://github.com/owner/repo/issues/382
```

**Best for:**
- Bug fixes
- Small features
- Refactoring tasks
- Test coverage improvements
- Documentation updates

---

### Option 2: 📋 Manual Implementation

**Use when:** You want to implement the solution yourself

```bash
# Step 1: Create branch + draft PR
/dev:create-fix-pr 382

# Step 2: Implement manually
# (make your code changes)

# Step 3: Commit and push
git add .
git commit -m "fix: resolve issue"
git push

# Step 4: Mark ready and merge
gh pr ready <PR_NUMBER>
/dev:merge-pr <PR_NUMBER>
```

**Best for:**
- Complex features requiring human decisions
- Large architectural changes
- Cross-cutting concerns
- Learning/exploration

---

### Option 3: 🔍 Discovery → Solve

**Use when:** Starting from scratch, need to find issues first

```bash
# Step 1: Scan codebase for issues
/dev:identify-issues services/accounting-ledger/ --create-issues

# Output:
#   Created issues: #380, #381, #382, #383

# Step 2: Solve each issue
/dev:solve-issue 380
/dev:solve-issue 381
# ... etc
```

**Best for:**
- New to codebase
- Improving code quality
- Security audits
- Technical debt cleanup

---

### Option 4: 📊 PR Review → Fix Issues

**Use when:** Reviewing existing PRs and finding problems

```bash
# Step 1: Review a PR
/dev:review-pr 100

# Output:
#   Found 5 issues
#   Created sub-issues: #101, #102, #103, #104, #105

# Step 2: Fix issues individually
/dev:solve-issue 101  # Fix critical security issue
/dev:solve-issue 102  # Fix SOLID violation
# ... etc

# Step 3: Re-review original PR
/dev:review-pr 100
```

**Best for:**
- Code review processes
- Quality improvement
- Before merging large PRs

---

### Option 5: 🎯 Batch Multiple Issues

**Use when:** Multiple related issues should be fixed together

```bash
# Create single PR for multiple issues
/dev:create-fix-pr 380 381 382

# What it does:
# - Creates branch: fix/issue-380-381-382-description-unique
# - Creates PR that closes all 3 issues
# - You implement all fixes in one go

# Then implement and merge
# (all 3 issues close when PR merges)
```

**Best for:**
- Related bugs
- Coordinated refactoring
- Feature sets
- Compliance fixes

---

## Complete Workflow Example

### Scenario: Increase Test Coverage

```bash
# Issue #382: [CRITICAL] Increase test coverage to 85% minimum

# Solution: Use automated workflow
/dev:solve-issue 382

# Behind the scenes:
# ✅ Fetches issue details
# ✅ Creates branch: test/issue-382-increase-test-coverage-85-a7b3c9f2
# ✅ AI analyzes codebase
# ✅ AI writes missing tests
# ✅ AI runs tests to verify coverage
# ✅ Commits: "test: increase coverage to 85% minimum"
# ✅ Pushes branch
# ✅ Creates PR #400 linked to issue #382

# Review the PR
/dev:review-pr 400

# If good, merge
/dev:merge-pr 400

# Result: Issue #382 auto-closes! ✅
```

---

## Branch Naming Convention

All commands create branches following this pattern:

```
{type}/issue-{number}-{description}-{unique}
```

**Examples:**
- `fix/issue-382-sql-injection-a7b3c9f2`
- `feature/issue-401-export-tax-report-3f9e2b1c`
- `security/issue-425-fix-auth-bypass-e4d8a7b9`
- `refactor/issue-500-apply-solid-9c2f1a6d`

**Type prefixes:**
- `fix/` - Bug fixes
- `feature/` - New features
- `security/` - Security issues
- `refactor/` - Code improvements
- `docs/` - Documentation
- `test/` - Test additions
- `perf/` - Performance

---

## Commit Message Format

All workflows follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <short description>

<detailed description>

Closes #<issue>

🤖 Generated with Claude Code
```

**Example:**
```
fix: resolve SQL injection in account query

Sanitized user input in GetAccountsByCompany query.
Added parameterized queries and input validation.

Closes #425

🤖 Generated with Claude Code /dev:solve-issue
```

---

## Quality Guarantees

All automated workflows ensure:

- ✅ **SOLID Principles** followed
- ✅ **Test Coverage** ≥85%
- ✅ **Linting** passes (golangci-lint/ESLint)
- ✅ **Type Safety** verified (TypeScript)
- ✅ **Build** succeeds
- ✅ **Security** checks (no vulnerabilities)
- ✅ **Multi-tenant** isolation maintained
- ✅ **Indonesian Compliance** (SAK EP, tax calculations)

---

## Integration with Git Workflow

### Standard Git Flow
```bash
main → feature/fix branch → PR → main
```

### With Claude Code Commands
```bash
# Traditional (manual)
git checkout main
git pull
git checkout -b fix/issue-382-description
# ... implement ...
git add .
git commit -m "fix: ..."
git push
gh pr create
gh pr merge

# With /dev:solve-issue (automated)
/dev:solve-issue 382
# Done! 🎉
```

---

## When to Use Each Workflow

| Workflow | Use Case | Time Saved | Control |
|----------|----------|------------|---------|
| `/dev:solve-issue` | Bug fixes, small features | 90% | Low |
| `/dev:create-fix-pr` + manual | Complex features | 50% | High |
| `/dev:identify-issues` → solve | Code quality improvement | 80% | Medium |
| `/dev:review-pr` → solve issues | PR quality gates | 70% | Medium |
| Batch with `/dev:create-fix-pr` | Related issues | 60% | Medium |

---

## Pro Tips

### 1. **Always start with issue identification**
```bash
/dev:identify-issues --all --create-issues
# Creates prioritized issue list
```

### 2. **Use draft PRs for experimentation**
```bash
/dev:solve-issue 382 --draft-only
# Review implementation before marking ready
```

### 3. **Combine with automated review**
```bash
/dev:solve-issue 382
/dev:review-pr --current
# Double-check AI implementation
```

### 4. **Verify before merging**
```bash
/dev:verify-pr <PR_NUMBER>
# Runs all checks: tests, lint, build, security
```

### 5. **Use for learning**
```bash
# Create draft PR manually
/dev:create-fix-pr 382 --draft

# Implement yourself
# (learn the codebase)

# Then get AI review
/dev:review-pr --current
# Learn what you missed
```

---

## Troubleshooting

### "Branch already exists"
**Cause:** Unique suffix prevents this, but if it happens:
```bash
# Delete local branch
git branch -D <branch-name>

# Delete remote branch
git push origin --delete <branch-name>

# Try again
/dev:solve-issue 382
```

### "Quality checks failed"
**Cause:** Tests/linting failed during automated implementation
```bash
# Command will stop and show errors
# Fix manually:
git add .
git commit -m "fix: address quality issues"
git push

# Or start over:
git checkout main
git branch -D <branch-name>
/dev:solve-issue 382
```

### "Issue not found"
**Cause:** Issue number incorrect or doesn't exist
```bash
# Check issue exists:
gh issue view 382

# Or use URL:
/dev:solve-issue https://github.com/owner/repo/issues/382
```

---

## Advanced Workflows

### Workflow: Security Issue → Immediate Fix → Auto-Merge
```bash
# Critical security issue found
/dev:solve-issue 425 --auto-merge

# AI implements fix
# Tests pass
# Auto-merges when CI green
# Issue auto-closes
```

### Workflow: Feature Request → Implementation → Review → Merge
```bash
# 1. Implement feature
/dev:solve-issue 401 --draft-only

# 2. Get deep review
/dev:review-pr --current

# 3. Fix review issues (if any)
/dev:solve-issue 410  # Fix sub-issue from review

# 4. Verify everything
/dev:verify-pr <PR_NUMBER>

# 5. Merge
/dev:merge-pr <PR_NUMBER>
```

### Workflow: Technical Debt Cleanup Sprint
```bash
# 1. Identify all issues
/dev:identify-issues --all --create-issues

# 2. Solve SOLID violations
/dev:solve-issue 501  # SRP violation
/dev:solve-issue 502  # DIP violation
/dev:solve-issue 503  # ISP violation

# 3. Solve test coverage gaps
/dev:solve-issue 504
/dev:solve-issue 505

# 4. Solve documentation gaps
/dev:solve-issue 506
```

---

## Comparison with Manual Process

### Manual Process (Traditional)
```bash
# Time: ~2-4 hours per issue

1. Read issue (5 min)
2. Create branch (2 min)
3. Implement solution (60-120 min)
4. Write tests (30-60 min)
5. Fix linting errors (10 min)
6. Commit and push (5 min)
7. Create PR (10 min)
8. Write PR description (15 min)
9. Link to issue (2 min)
10. Request review (5 min)

Total: 2-4 hours
```

### With `/dev:solve-issue` (Automated)
```bash
# Time: ~5-15 minutes per issue

1. Run command (1 min)
2. Wait for AI implementation (3-10 min)
3. Review result (1-3 min)
4. Merge (1 min)

Total: 5-15 minutes
```

**Time saved: 85-95%** ⚡

---

## Best Practices

1. **Use descriptive issue titles**
   - Good: "Fix SQL injection in GetAccountsByCompany"
   - Bad: "Bug in accounts"

2. **Add labels to issues**
   - Helps determine branch prefix
   - Auto-applied to PRs

3. **Write detailed issue descriptions**
   - AI uses this for implementation
   - Include acceptance criteria

4. **Review AI implementations**
   - Use `/dev:review-pr --current`
   - Verify SOLID principles

5. **Keep issues focused**
   - One issue = one concern
   - Use batch for related issues

6. **Close stale branches**
   - Use `/dev:merge-pr` for cleanup
   - Auto-deletes merged branches

---

## Summary

**For maximum efficiency:**
1. Use `/dev:identify-issues` to find problems
2. Use `/dev:solve-issue` for automated fixes
3. Use `/dev:review-pr` for quality assurance
4. Use `/dev:merge-pr` for safe merging

**Result:**
- ⚡ 90% faster issue resolution
- ✅ Higher code quality (SOLID + tests)
- 🔒 Better security (automated checks)
- 📝 Better documentation (automated)
- 🎯 Fewer bugs (comprehensive testing)

---

🤖 **All workflows powered by Claude Code**

Ready to solve issues at lightning speed! 🚀
