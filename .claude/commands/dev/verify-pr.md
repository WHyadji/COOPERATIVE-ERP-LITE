---
allowed-tools: Read, Write, Edit, Bash, Grep, Glob
argument-hint: [pr-number] | --current | --fix-issues
description: Comprehensive pre-merge verification: tests, linting, builds, security, and quality checks
---

# Pre-Merge PR Verification Command

Performs comprehensive verification checks before merging a PR to ensure code quality, test coverage, build success, and security compliance.

## Context Gathering

Detecting PR: !`
# Determine PR reference from arguments
PR_REF="$ARGUMENTS"

if [ -z "$PR_REF" ] || [ "$PR_REF" = "--current" ]; then
    BRANCH=$(git branch --show-current)
    gh pr list --head "$BRANCH" --json number --jq '.[0].number' || echo "none"
elif [[ "$PR_REF" =~ ^[0-9]+$ ]]; then
    echo "$PR_REF"
else
    echo "$PR_REF" | grep -oP '[0-9]+' || echo "none"
fi
`

PR Details: !`gh pr view $PR_NUMBER --json number,title,author,state,isDraft,mergeable,additions,deletions,changedFiles,url 2>/dev/null || echo "{}"`

Changed Files: !`gh pr diff $PR_NUMBER --name-only 2>/dev/null || echo ""`

Current Branch: !`git branch --show-current`

Git Status: !`git status --short`

Has Uncommitted Changes: !`git status --porcelain | wc -l`

## Your Task

Perform comprehensive pre-merge verification for PR #$PR_NUMBER. This ensures the PR meets all quality, testing, build, and security requirements before merging.

### Phase 1: Pre-Flight Checks

1. **Validate PR exists and is mergeable**
   ```bash
   if PR_NUMBER is "none" or empty:
       ERROR: "Cannot find PR. Specify PR number or run from PR branch"
       EXIT

   if PR is merged or closed:
       ERROR: "PR #$PR_NUMBER is already merged/closed"
       EXIT

   if PR is draft:
       WARNING: "PR #$PR_NUMBER is still a draft"

   if PR has merge conflicts:
       ERROR: "PR #$PR_NUMBER has merge conflicts. Resolve before verification"
       EXIT
   ```

2. **Check for uncommitted changes**
   ```bash
   if uncommitted changes exist:
       WARNING: "You have uncommitted changes. Stash or commit them first."
       if not "--force" flag:
           EXIT
   ```

3. **Detect change type (Go backend / Frontend / Full-stack)**
   ```bash
   CHANGED_FILES=$(gh pr diff $PR_NUMBER --name-only)

   HAS_GO=$(echo "$CHANGED_FILES" | grep "\.go$" | wc -l)
   HAS_FRONTEND=$(echo "$CHANGED_FILES" | grep -E "\.(tsx?|jsx?|css)$" | wc -l)

   if HAS_GO > 0 && HAS_FRONTEND > 0:
       CHANGE_TYPE="full-stack"
   elif HAS_GO > 0:
       CHANGE_TYPE="backend"
   else:
       CHANGE_TYPE="frontend"
   ```

### Phase 2: Test Verification

Run comprehensive test suites based on change type.

#### Backend Tests (if Go files changed)

```bash
echo "🧪 Running Backend Tests..."

# 1. Unit Tests
echo "  → Running unit tests..."
cd services/accounting-ledger && go test ./... -v -cover -coverprofile=coverage.out
UNIT_STATUS=$?

# 2. Integration Tests
echo "  → Running integration tests..."
cd services/accounting-ledger && go test -tags=integration ./... -v
INTEGRATION_STATUS=$?

# 3. API Gateway Tests
echo "  → Running API Gateway tests..."
cd services/api-gateway && go test ./... -v
GATEWAY_STATUS=$?

# 4. Business Service Tests
echo "  → Running Business Service tests..."
cd services/business-service && go test ./... -v
BUSINESS_STATUS=$?

# Calculate coverage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

if [ "$COVERAGE" -lt 85 ]; then
    WARNING: "Test coverage is ${COVERAGE}% (target: 85%)"
fi
```

#### Frontend Tests (if TypeScript/React files changed)

```bash
echo "🧪 Running Frontend Tests..."

# 1. Main App Tests
echo "  → Running frontend unit tests..."
cd frontend && npm run test -- --coverage --watchAll=false
FRONTEND_UNIT_STATUS=$?

# 2. Admin Panel Tests
if changed files in admin-frontend:
    echo "  → Running admin panel tests..."
    cd admin-frontend && npm run test -- --coverage --watchAll=false
    ADMIN_TEST_STATUS=$?

# 3. E2E Tests (critical paths only)
echo "  → Running E2E tests..."
cd frontend && npm run test:e2e
E2E_STATUS=$?
```

**Test Results Summary:**
```
Backend Tests:
  ✅/❌ Unit Tests: $UNIT_STATUS
  ✅/❌ Integration Tests: $INTEGRATION_STATUS
  ✅/❌ API Gateway: $GATEWAY_STATUS
  ✅/❌ Business Service: $BUSINESS_STATUS
  📊 Coverage: ${COVERAGE}%

Frontend Tests:
  ✅/❌ Unit Tests: $FRONTEND_UNIT_STATUS
  ✅/❌ Admin Tests: $ADMIN_TEST_STATUS
  ✅/❌ E2E Tests: $E2E_STATUS
```

### Phase 3: Code Quality Checks

#### Backend Quality (Go)

```bash
echo "🔍 Running Go Code Quality Checks..."

# 1. Formatting Check
echo "  → Checking gofmt..."
cd services && gofmt -l $(find . -name "*.go") > /tmp/gofmt-issues.txt
GOFMT_ISSUES=$(cat /tmp/gofmt-issues.txt | wc -l)

if [ $GOFMT_ISSUES -gt 0 ]; then
    WARNING: "$GOFMT_ISSUES files need formatting"
    if "--fix-issues" flag:
        gofmt -w $(cat /tmp/gofmt-issues.txt)
fi

# 2. Go Vet
echo "  → Running go vet..."
cd services/accounting-ledger && go vet ./...
VET_STATUS=$?

# 3. Linting with golangci-lint
echo "  → Running golangci-lint..."
cd services/accounting-ledger && golangci-lint run --timeout=5m
LINT_STATUS=$?

if [ $LINT_STATUS -ne 0 ]; then
    if "--fix-issues" flag:
        golangci-lint run --fix
fi
```

#### Frontend Quality (TypeScript/React)

```bash
echo "🔍 Running Frontend Code Quality Checks..."

# 1. TypeScript Type Checking
echo "  → Running TypeScript compiler..."
cd frontend && npx tsc --noEmit
TSC_STATUS=$?

# 2. ESLint
echo "  → Running ESLint..."
cd frontend && npm run lint
ESLINT_STATUS=$?

if [ $ESLINT_STATUS -ne 0 ] && "--fix-issues" flag:
    npm run lint -- --fix

# 3. Format Check
echo "  → Checking code formatting..."
cd frontend && npx prettier --check "**/*.{ts,tsx,js,jsx,css}"
PRETTIER_STATUS=$?

if [ $PRETTIER_STATUS -ne 0 ] && "--fix-issues" flag:
    npx prettier --write "**/*.{ts,tsx,js,jsx,css}"
```

**Quality Check Results:**
```
Backend Quality:
  ✅/❌ Formatting: $GOFMT_ISSUES issues
  ✅/❌ Go Vet: $VET_STATUS
  ✅/❌ Linting: $LINT_STATUS

Frontend Quality:
  ✅/❌ TypeScript: $TSC_STATUS
  ✅/❌ ESLint: $ESLINT_STATUS
  ✅/❌ Prettier: $PRETTIER_STATUS
```

### Phase 4: Build Verification

Build all affected services to ensure compilation succeeds.

```bash
echo "🔨 Running Build Verification..."

if CHANGE_TYPE includes "backend":
    # Build Go Services
    echo "  → Building accounting-ledger..."
    cd services/accounting-ledger && go build -o bin/server ./cmd/server
    LEDGER_BUILD=$?

    echo "  → Building api-gateway..."
    cd services/api-gateway && go build -o bin/server ./cmd/server
    GATEWAY_BUILD=$?

    echo "  → Building business-service..."
    cd services/business-service && go build -o bin/server ./cmd/server
    BUSINESS_BUILD=$?

if CHANGE_TYPE includes "frontend":
    # Build Frontend
    echo "  → Building frontend..."
    cd frontend && npm run build
    FRONTEND_BUILD=$?

    if changed admin-frontend:
        echo "  → Building admin-frontend..."
        cd admin-frontend && npm run build
        ADMIN_BUILD=$?

# Alternative: Use build script
# ./scripts/build.sh all
```

**Build Results:**
```
Backend Builds:
  ✅/❌ Accounting Ledger: $LEDGER_BUILD
  ✅/❌ API Gateway: $GATEWAY_BUILD
  ✅/❌ Business Service: $BUSINESS_BUILD

Frontend Builds:
  ✅/❌ Frontend: $FRONTEND_BUILD
  ✅/❌ Admin Frontend: $ADMIN_BUILD
```

### Phase 5: Security Scanning

```bash
echo "🔒 Running Security Scans..."

# 1. Frontend Security (npm audit)
if CHANGE_TYPE includes "frontend":
    echo "  → Running npm audit..."
    cd frontend && npm audit --audit-level=moderate
    NPM_AUDIT=$?

    cd admin-frontend && npm audit --audit-level=moderate
    ADMIN_AUDIT=$?

# 2. Go Security (govulncheck)
if CHANGE_TYPE includes "backend":
    echo "  → Running Go vulnerability check..."
    if command -v govulncheck &> /dev/null; then
        cd services/accounting-ledger && govulncheck ./...
        GO_VULN=$?
    else:
        WARNING: "govulncheck not installed. Install: go install golang.org/x/vuln/cmd/govulncheck@latest"

# 3. Dependency Check
echo "  → Checking for outdated dependencies..."
if CHANGE_TYPE includes "backend":
    cd services/accounting-ledger && go list -u -m -json all | jq -r 'select(.Update != null) | "\(.Path): \(.Version) -> \(.Update.Version)"' > /tmp/go-updates.txt
    GO_UPDATES=$(cat /tmp/go-updates.txt | wc -l)

if CHANGE_TYPE includes "frontend":
    cd frontend && npm outdated > /tmp/npm-updates.txt 2>&1
    NPM_UPDATES=$(cat /tmp/npm-updates.txt | wc -l)
```

**Security Results:**
```
Security Scans:
  ✅/❌ NPM Audit: $NPM_AUDIT
  ✅/❌ Go Vulnerabilities: $GO_VULN
  📦 Outdated Go Deps: $GO_UPDATES
  📦 Outdated NPM Deps: $NPM_UPDATES
```

### Phase 6: CI Status Check

```bash
echo "🤖 Checking CI Status..."

CI_STATUS=$(gh pr checks $PR_NUMBER --json name,state,link)

FAILED_CHECKS=$(echo "$CI_STATUS" | jq '[.[] | select(.state != "SUCCESS" and .state != "SKIPPED")] | length')

if [ $FAILED_CHECKS -gt 0 ]; then
    ERROR: "$FAILED_CHECKS CI checks failed"
    echo "$CI_STATUS" | jq -r '.[] | select(.state != "SUCCESS") | "  ❌ \(.name): \(.state)"'
else:
    SUCCESS: "All CI checks passed"
```

### Phase 7: PR Metadata Validation

```bash
echo "📋 Validating PR Metadata..."

# 1. Check for issue references
PR_BODY=$(gh pr view $PR_NUMBER --json body --jq '.body')
CLOSES_ISSUES=$(echo "$PR_BODY" | grep -oP '(Closes|Fixes|Resolves) #\K[0-9]+' | wc -l)

if [ $CLOSES_ISSUES -eq 0 ]; then
    WARNING: "PR doesn't reference any issues. Add 'Closes #123' to PR description"

# 2. Check for breaking changes
HAS_BREAKING=$(echo "$PR_BODY" | grep -i "breaking change" | wc -l)
if [ $HAS_BREAKING -gt 0 ]; then
    WARNING: "PR contains BREAKING CHANGES. Ensure proper communication."

# 3. Check PR labels
PR_LABELS=$(gh pr view $PR_NUMBER --json labels --jq '.labels[].name' | tr '\n' ',' | sed 's/,$//')
if [ -z "$PR_LABELS" ]; then
    WARNING: "PR has no labels. Add appropriate labels."

# 4. Check for reviewers
REVIEWERS=$(gh pr view $PR_NUMBER --json reviewRequests --jq '.reviewRequests | length')
REVIEWS=$(gh pr view $PR_NUMBER --json reviews --jq '.reviews | length')

if [ $REVIEWERS -eq 0 ] && [ $REVIEWS -eq 0 ]; then
    WARNING: "No reviewers assigned. Consider requesting reviews."
```

### Phase 8: Generate Verification Report

Create comprehensive verification report:

```markdown
# 🔍 PR Verification Report

**PR**: #{PR_NUMBER} - {PR_TITLE}
**Author**: {AUTHOR}
**Branch**: {BRANCH}
**Changes**: +{ADDITIONS} / -{DELETIONS} lines across {FILES_CHANGED} files
**Type**: {CHANGE_TYPE}

---

## ✅ Verification Results

### Test Coverage
| Suite | Status | Details |
|-------|--------|---------|
| Backend Unit | {STATUS} | {COVERAGE}% coverage |
| Backend Integration | {STATUS} | - |
| API Gateway | {STATUS} | - |
| Business Service | {STATUS} | - |
| Frontend Unit | {STATUS} | - |
| E2E Tests | {STATUS} | - |

**Overall**: {PASS/FAIL}

### Code Quality
| Check | Status | Issues |
|-------|--------|--------|
| Go Formatting | {STATUS} | {COUNT} |
| Go Vet | {STATUS} | - |
| golangci-lint | {STATUS} | {COUNT} |
| TypeScript | {STATUS} | {COUNT} |
| ESLint | {STATUS} | {COUNT} |
| Prettier | {STATUS} | {COUNT} |

**Overall**: {PASS/FAIL}

### Build Status
| Component | Status |
|-----------|--------|
| Accounting Ledger | {STATUS} |
| API Gateway | {STATUS} |
| Business Service | {STATUS} |
| Frontend | {STATUS} |
| Admin Frontend | {STATUS} |

**Overall**: {PASS/FAIL}

### Security
| Check | Status | Findings |
|-------|--------|----------|
| NPM Audit | {STATUS} | {COUNT} vulnerabilities |
| Go Vulnerabilities | {STATUS} | {COUNT} issues |
| Outdated Dependencies | {STATUS} | Go: {COUNT}, NPM: {COUNT} |

**Overall**: {PASS/FAIL}

### CI/CD
| Check | Status |
|-------|--------|
{CI checks status}

**Overall**: {PASS/FAIL}

### PR Metadata
- **Linked Issues**: {COUNT} (Closes #{ISSUES})
- **Breaking Changes**: {YES/NO}
- **Labels**: {LABELS}
- **Reviewers**: {COUNT}

---

## 🎯 Verification Summary

**Status**: {✅ READY TO MERGE | ⚠️ WARNINGS | ❌ FAILED}

### Critical Issues ({COUNT})
{List of blocking issues that MUST be fixed}

### Warnings ({COUNT})
{List of non-blocking warnings}

### Recommendations
{List of recommended actions}

---

## 🚦 Merge Readiness

{If all critical checks pass:}
✅ **READY TO MERGE**
- All tests passing
- Code quality checks passed
- Builds successful
- No critical security issues
- CI checks passed

**Next Steps:**
1. Get final review approval
2. Run `/dev:merge-pr {PR_NUMBER}` to merge safely

{If critical issues exist:}
❌ **NOT READY TO MERGE**

**Required Actions:**
1. {Fix critical issue 1}
2. {Fix critical issue 2}
3. Re-run `/dev:verify-pr {PR_NUMBER}`

{If only warnings:}
⚠️ **MERGE WITH CAUTION**

**Recommended Actions:**
1. {Address warning 1}
2. {Address warning 2}
3. Consider fixing before merge

---

🤖 **Verification completed at {TIMESTAMP}**

Run with `--fix-issues` to auto-fix linting and formatting issues.
```

### Phase 9: Post Report to GitHub

```bash
# Save report to file
echo "$REPORT" > /tmp/pr-verification-report.md

# Post as PR comment
gh pr comment $PR_NUMBER -F /tmp/pr-verification-report.md

echo "✅ Verification report posted to PR #$PR_NUMBER"
echo "View at: {PR_URL}"
```

### Phase 10: Exit with Status

```bash
# Count critical failures
CRITICAL_FAILURES=0
if any test failed: CRITICAL_FAILURES++
if any build failed: CRITICAL_FAILURES++
if critical security issues: CRITICAL_FAILURES++
if CI checks failed: CRITICAL_FAILURES++

if [ $CRITICAL_FAILURES -eq 0 ]; then
    echo "✅ PR #$PR_NUMBER is ready to merge!"
    exit 0
else:
    echo "❌ PR #$PR_NUMBER has $CRITICAL_FAILURES critical failures"
    echo "Fix the issues and re-run verification"
    exit 1
fi
```

## Command Options

- **Default**: `/dev:verify-pr 123` - Full verification
- **Current Branch**: `/dev:verify-pr --current` - Auto-detect PR
- **Auto-Fix**: `/dev:verify-pr 123 --fix-issues` - Auto-fix linting/formatting
- **Quick**: `/dev:verify-pr 123 --quick` - Skip E2E and integration tests
- **Force**: `/dev:verify-pr 123 --force` - Skip pre-flight checks

## Example Output

```
🔍 Verifying PR #123: Fix tax calculation issues

✅ Pre-flight checks passed
🧪 Running tests... ████████████████████ 100%
  ✅ Backend unit tests: PASSED (87% coverage)
  ✅ Integration tests: PASSED
  ✅ Frontend tests: PASSED
🔍 Running code quality checks...
  ✅ Go formatting: OK
  ✅ golangci-lint: OK
  ✅ TypeScript: OK
  ✅ ESLint: OK
🔨 Building services...
  ✅ All builds successful
🔒 Running security scans...
  ✅ No vulnerabilities found
🤖 CI Status: All checks passed

✅ PR #123 is READY TO MERGE!

View full report: https://github.com/.../pull/123#issuecomment-...

Next steps:
  - Get approval from reviewers
  - Run /dev:merge-pr 123 to merge safely
```

## Notes

- This command is **read-only** - it never modifies code (unless `--fix-issues` specified)
- All checks run in parallel where possible for speed
- Failed checks don't stop execution - full report always generated
- Report posted as comment for transparency
- Use before requesting final review or merging
