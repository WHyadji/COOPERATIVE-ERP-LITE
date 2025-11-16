---
allowed-tools: Read, Write, Edit, Bash, Grep, Glob, Task
description: Proactive codebase scanning to identify issues before creating PRs
argument-hint: [path] | --all | --security-only | --create-issues
---

# Proactive Issue Identification Scanner

Scans codebase to proactively identify issues before creating PRs: code quality problems, security vulnerabilities, TODOs, missing tests, and SOLID violations.

## Context

Scan Target: !`echo "${ARGUMENTS:-services/}" | sed 's/--[a-z-]*//g' | xargs`

Project Root: !`pwd`

## Your Task

Scan the specified path for issues across multiple dimensions: code quality, security, performance, testing, and documentation.

### Phase 1: Determine Scan Scope

```bash
SCAN_PATH="${ARGUMENTS:-services/}"
SCAN_ALL=$(echo "$ARGUMENTS" | grep -c "\--all")
SECURITY_ONLY=$(echo "$ARGUMENTS" | grep -c "\--security-only")
CREATE_ISSUES=$(echo "$ARGUMENTS" | grep -c "\--create-issues")

if [ $SCAN_ALL -eq 1 ]; then
    SCAN_PATH="."
fi

echo "🔍 Scanning: $SCAN_PATH"
```

### Phase 2: Run Scans in Parallel

Launch specialized scanning agents in parallel:

#### Agent 1: Code Quality Scan
```
Task: Scan $SCAN_PATH for:
- SOLID principle violations
- Code smells and anti-patterns
- Complex functions (>20 lines)
- TODO/FIXME/HACK comments
- Placeholder code
- Unused imports/variables
- Naming convention violations
- Missing error handling

Return: List of issues with file:line, severity, description
```

#### Agent 2: Security Scan
```
Task: Scan $SCAN_PATH for:
- Hardcoded secrets/credentials
- SQL injection vulnerabilities
- Insecure dependencies
- Missing input validation
- Authentication bypasses
- Data exposure in logs
- Weak cryptography

Return: List of security issues with severity and remediation
```

#### Agent 3: Testing Gaps
```
Task: Analyze $SCAN_PATH for:
- Functions without tests
- Low test coverage areas (<85%)
- Missing integration tests
- Untested error paths
- Missing edge case tests

Return: List of testing gaps with priority
```

### Phase 3: Additional Checks

```bash
# 1. Find placeholder implementations
echo "📝 Checking for placeholders..."
grep -rn "TODO\|FIXME\|HACK\|XXX" $SCAN_PATH --include="*.go" --include="*.ts" --include="*.tsx"

# 2. Find mock data
echo "🎭 Checking for hardcoded mock data..."
grep -rn "mockData\|hardcoded\|placeholder" $SCAN_PATH --include="*.go" --include="*.ts" --include="*.tsx"

# 3. Run linters
if [ -d "services" ]; then
    echo "🔧 Running golangci-lint..."
    cd services/accounting-ledger && golangci-lint run --out-format=json > /tmp/lint-issues.json
fi

if [ -d "frontend" ]; then
    echo "🔧 Running ESLint..."
    cd frontend && npm run lint -- --format=json > /tmp/eslint-issues.json 2>/dev/null || true
fi

# 4. Security audit
if [ $SECURITY_ONLY -eq 1 ] || [ $SCAN_ALL -eq 1 ]; then
    echo "🔒 Running security audits..."
    cd frontend && npm audit --json > /tmp/npm-audit.json 2>/dev/null || true
fi
```

### Phase 4: Aggregate and Prioritize

Compile all findings and assign priorities:

```
CRITICAL Priority:
- Security vulnerabilities
- Data corruption risks
- Authentication bypasses
- SQL injection vectors

HIGH Priority:
- SOLID violations affecting maintainability
- Missing critical tests
- Performance bottlenecks
- Multi-tenancy isolation issues

MEDIUM Priority:
- Code quality issues
- Documentation gaps
- TODO/FIXME comments
- Test coverage < 85%

LOW Priority:
- Minor style violations
- Optimization opportunities
- Nice-to-have improvements
```

### Phase 5: Generate Report

```markdown
# 🔍 Issue Identification Report

**Scan Target**: {SCAN_PATH}
**Timestamp**: {TIMESTAMP}
**Total Issues Found**: {COUNT}

## Summary by Severity

| Severity | Count | Category Breakdown |
|----------|-------|-------------------|
| 🔴 CRITICAL | {COUNT} | Security: {N}, Data: {N} |
| 🟠 HIGH | {COUNT} | SOLID: {N}, Testing: {N}, Performance: {N} |
| 🟡 MEDIUM | {COUNT} | Quality: {N}, Docs: {N} |
| ⚪ LOW | {COUNT} | Style: {N}, Optimization: {N} |

## Critical Issues ({COUNT})

{For each critical issue:}
### [{FILE}:{LINE}] {TITLE}

**Severity**: CRITICAL
**Category**: {CATEGORY}
**Impact**: {IMPACT_DESCRIPTION}

**Description**:
{DETAILED_DESCRIPTION}

**Recommendation**:
{FIX_RECOMMENDATION}

**Example Fix**:
```{language}
{SUGGESTED_CODE}
```

---

## High Priority Issues ({COUNT})

{List high priority issues...}

## Suggested Fix Order

1. **Phase 1 - Critical** (Must fix immediately):
   - Issue #{N}: {TITLE}
   - Issue #{N}: {TITLE}

2. **Phase 2 - High** (Fix before next release):
   - Issue #{N}: {TITLE}
   - Issue #{N}: {TITLE}

3. **Phase 3 - Medium** (Address in sprint):
   - Issue #{N}: {TITLE}

4. **Phase 4 - Low** (Backlog):
   - Issue #{N}: {TITLE}

## Recommended Actions

1. Create fix PRs for critical and high issues: `/dev:create-fix-pr {ISSUE_NUMBERS}`
2. Address security issues first
3. Improve test coverage to ≥85%
4. Refactor SOLID violations
5. Document complex functions

---

🤖 Generated by `/dev:identify-issues` scanner
```

### Phase 6: Create GitHub Issues (if requested)

```bash
if [ $CREATE_ISSUES -eq 1 ]; then
    echo "📋 Creating GitHub issues..."

    for FINDING in ${CRITICAL_FINDINGS[@]}; do
        gh issue create \
            --title "[CRITICAL] ${FINDING_TITLE}" \
            --body "${FINDING_DESCRIPTION}" \
            --label "critical,${FINDING_CATEGORY}" \
            --assignee "@me"
    done

    for FINDING in ${HIGH_FINDINGS[@]}; do
        gh issue create \
            --title "[HIGH] ${FINDING_TITLE}" \
            --body "${FINDING_DESCRIPTION}" \
            --label "high-priority,${FINDING_CATEGORY}" \
            --assignee "@me"
    done

    echo "✅ Created $(expr ${#CRITICAL_FINDINGS[@]} + ${#HIGH_FINDINGS[@]}) issues"
fi
```

### Phase 7: Print Summary

```bash
echo ""
echo "════════════════════════════════════════"
echo "🔍 SCAN COMPLETE"
echo "════════════════════════════════════════"
echo "Total Issues: $TOTAL_ISSUES"
echo "  🔴 Critical: $CRITICAL_COUNT"
echo "  🟠 High:     $HIGH_COUNT"
echo "  🟡 Medium:   $MEDIUM_COUNT"
echo "  ⚪ Low:      $LOW_COUNT"
echo "════════════════════════════════════════"
echo ""
echo "📄 Full report saved to: ./issue-report.md"
echo ""
echo "Next steps:"
echo "  1. Review critical issues first"
echo "  2. Create fix PR: /dev:create-fix-pr {ISSUE_NUMBERS}"
echo "  3. Or create issues: /dev:identify-issues $SCAN_PATH --create-issues"
echo ""
```

## Usage Examples

```bash
# Scan specific service
/dev:identify-issues services/accounting-ledger/

# Scan entire codebase
/dev:identify-issues --all

# Security-focused scan
/dev:identify-issues --security-only

# Scan and auto-create GitHub issues
/dev:identify-issues services/api-gateway/ --create-issues

# Scan specific file
/dev:identify-issues frontend/components/Dashboard.tsx
```

## Notes

- Proactive issue discovery before PR creation
- Multi-dimensional analysis (quality, security, testing)
- Prioritized findings for efficient fixing
- Optional GitHub issue creation
- Integrates with existing lint tools
- Project-specific checks (SOLID, multi-tenancy, Indonesian compliance)
