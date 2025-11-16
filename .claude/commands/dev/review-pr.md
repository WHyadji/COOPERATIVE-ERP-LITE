---
allowed-tools: Read, Write, Edit, Bash, Task, Grep, Glob
argument-hint: [pr-number] | [pr-url] | --current | --no-issues
description: Comprehensive PR review with multi-agent analysis, automated GitHub comments, and sub-issue creation
---

# Deep PR Review Command

Performs comprehensive pull request review using specialized agents for code quality, security, and performance analysis. Automatically posts review comments to GitHub and creates sub-issues for all findings.

## Context Gathering

**Arguments provided:** `$ARGUMENTS`

**Your first task:** Parse the arguments to determine the PR number. Follow this logic:
- If empty or `--current`: Use `gh pr list --head "$(git branch --show-current)" --json number --jq '.[0].number'` to get PR for current branch
- If it matches a URL pattern (contains `github.com/`): Extract the number after `/pull/`
- If it's just a number: Use it directly
- If `--no-issues` is present: Remember to skip issue creation later

**Once you have the PR_NUMBER**, use the Bash tool to gather the following information (replace `{PR_NUM}` with the actual number):

```bash
# Get PR details
gh pr view {PR_NUM} --json number,title,author,body,additions,deletions,changedFiles,state,isDraft,labels,reviewDecision,url

# Get changed files
gh pr diff {PR_NUM} --name-only

# Get commit history
gh pr view {PR_NUM} --json commits --jq '.commits[] | "\(.oid[:8]) \(.messageHeadline)"'

# Get CI status
gh pr checks {PR_NUM} --json name,state,link --jq '.[] | "\(.name): \(.state)"'
```

**Also read:**
Current Project Structure: @/Users/adji/Documents/VISI-DIGITAL-TERPADU/accounting-webapp/CLAUDE.md

## Your Task

You are performing a comprehensive deep review of Pull Request #$PR_NUMBER using multiple specialized agents in parallel. This review must be thorough, actionable, and automatically integrated with GitHub.

### Phase 1: Parse Arguments and Validate

1. Extract PR number from the context gathering above
2. Check if `--no-issues` flag is present in $ARGUMENTS
3. Validate that the PR exists and is in reviewable state
4. Determine if this is a Go backend, TypeScript frontend, or full-stack change

### Phase 2: Parallel Agent Analysis

Launch the following agents **IN PARALLEL** using the Task tool with multiple tool calls in a single message:

#### Agent 1: Code Quality Guardian
```
Task: Analyze PR #$PR_NUMBER for:
- SOLID principles violations
- Code quality issues and anti-patterns
- Naming conventions compliance
- Function complexity and size
- Error handling patterns
- Testing coverage and quality
- Documentation completeness
- Project-specific patterns from CLAUDE.md

Return findings with:
- File path and line number
- Severity: CRITICAL, HIGH, MEDIUM, LOW
- Issue description
- Recommended fix
- SOLID principle violated (if applicable)
```

#### Agent 2: Security Specialist
```
Task: Perform security analysis on PR #$PR_NUMBER:
- SQL injection vulnerabilities
- Authentication/authorization issues
- Sensitive data exposure in logs
- Input validation gaps
- CSRF/XSS vulnerabilities
- Insecure dependencies
- Multi-tenancy isolation violations
- JWT token handling issues
- Rate limiting bypass potential

Return findings with:
- File path and line number
- Severity: CRITICAL, HIGH, MEDIUM, LOW
- Vulnerability type and impact
- Exploit scenario
- Remediation steps
```

#### Agent 3: Performance Engineer
```
Task: Analyze PR #$PR_NUMBER for performance issues:
- Database query optimization (N+1, missing indexes)
- Memory leaks and resource management
- Inefficient algorithms
- Bundle size impact (frontend)
- API response time impact
- Caching opportunities
- Unnecessary re-renders (React)
- Connection pool management

Return findings with:
- File path and line number
- Severity: CRITICAL, HIGH, MEDIUM, LOW
- Performance impact (quantified if possible)
- Bottleneck description
- Optimization recommendation
```

### Phase 3: Indonesian Compliance Review

After agents complete, perform domain-specific compliance check:

1. **Accounting Standards (SAK EP)**
   - Double-entry bookkeeping compliance
   - Chart of accounts structure
   - Journal entry validation
   - Financial statement accuracy

2. **Indonesian Tax Compliance**
   - PPH21/22/23/4(2) calculations correct
   - PPN handling proper
   - Tax rate historical accuracy
   - Audit trail completeness

3. **Multi-Tenancy Isolation**
   - All queries include company_id filter
   - RLS policies properly implemented
   - No data leakage between companies
   - Proper authorization checks

4. **Currency Handling**
   - Indonesian Rupiah formatting correct
   - Decimal precision appropriate
   - Use of centralized currency utils

### Phase 4: Synthesize Findings

Aggregate all findings from agents and compliance review:

1. **Categorize by severity:**
   - CRITICAL: Must fix before merge (security holes, data corruption, compliance violations)
   - HIGH: Should fix before merge (SOLID violations, performance issues, bugs)
   - MEDIUM: Should address soon (code quality, documentation)
   - LOW: Nice to have (minor optimizations, style improvements)

2. **Organize by file and line number**

3. **Calculate metrics:**
   - Total issues by severity
   - SOLID principles compliance score (0-100%)
   - Security risk level (LOW/MEDIUM/HIGH/CRITICAL)
   - Performance impact assessment
   - Test coverage percentage
   - Documentation completeness score

### Phase 5: Post Line-by-Line Comments

For each finding with a specific line number, post an inline comment using:

```bash
gh pr review $PR_NUMBER \
  --comment \
  --body "**[{SEVERITY}] {AGENT_NAME}**

{ISSUE_DESCRIPTION}

**Impact:** {IMPACT_DESCRIPTION}

**Recommendation:**
{DETAILED_FIX}

**Code Example:**
\`\`\`{language}
{SUGGESTED_CODE}
\`\`\`

---
🤖 Generated by Claude Code PR Review"
```

Post inline comments for:
- All CRITICAL findings
- All HIGH findings
- Top 10 MEDIUM findings (if more than 10, include rest in summary)

### Phase 6: Post Summary Comment

Create a comprehensive summary comment with the following structure:

```markdown
# 🔍 Deep PR Review Summary

## 📊 Overview

- **PR**: #{PR_NUMBER} - {PR_TITLE}
- **Author**: {AUTHOR}
- **Changes**: +{ADDITIONS} / -{DELETIONS} lines across {FILES_CHANGED} files
- **Type**: {Go Backend | Frontend | Full-Stack}

## 🎯 Review Results

### Severity Breakdown
- 🔴 **CRITICAL**: {COUNT} issues - MUST FIX
- 🟠 **HIGH**: {COUNT} issues - SHOULD FIX
- 🟡 **MEDIUM**: {COUNT} issues
- ⚪ **LOW**: {COUNT} issues

### Quality Metrics
- **SOLID Compliance**: {SCORE}% ✅/❌
- **Security Risk**: {LEVEL} 🔒
- **Performance Impact**: {ASSESSMENT} ⚡
- **Test Coverage**: {PERCENTAGE}% 🧪
- **Documentation**: {SCORE}% 📝

## 🔴 Critical Issues ({COUNT})

{List each critical issue with:}
1. **[File:Line]** Issue description
   - **Impact**: Why this is critical
   - **Action**: What must be done
   - **Issue**: #{ISSUE_NUMBER} (auto-created)

## 🟠 High Priority Issues ({COUNT})

{List each high issue with:}
1. **[File:Line]** Issue description
   - **Impact**: Why this matters
   - **Recommendation**: How to fix
   - **Issue**: #{ISSUE_NUMBER} (auto-created)

## 🟡 Medium Priority Issues ({COUNT})

{Summarize medium issues by category}
- **SOLID Violations**: {COUNT}
- **Code Quality**: {COUNT}
- **Documentation**: {COUNT}

{Top 5 medium issues listed}

## 💡 SOLID Principles Analysis

| Principle | Status | Issues |
|-----------|--------|--------|
| Single Responsibility | ✅/❌ | {COUNT} |
| Open/Closed | ✅/❌ | {COUNT} |
| Liskov Substitution | ✅/❌ | {COUNT} |
| Interface Segregation | ✅/❌ | {COUNT} |
| Dependency Inversion | ✅/❌ | {COUNT} |

**Details**: {Summary of SOLID violations}

## 🔒 Security Assessment

**Risk Level**: {CRITICAL | HIGH | MEDIUM | LOW}

{List security findings by category:}
- **Authentication/Authorization**: {FINDINGS}
- **Data Protection**: {FINDINGS}
- **Input Validation**: {FINDINGS}
- **Multi-Tenancy Isolation**: {FINDINGS}

## ⚡ Performance Analysis

**Impact**: {POSITIVE | NEUTRAL | NEGATIVE}

{Key performance findings:}
- **Database Queries**: {ANALYSIS}
- **API Response Time**: {ANALYSIS}
- **Memory Usage**: {ANALYSIS}
- **Bundle Size** (if frontend): {ANALYSIS}

## 🇮🇩 Indonesian Compliance Check

- **SAK EP Standards**: ✅/❌ {NOTES}
- **Tax Calculations**: ✅/❌ {NOTES}
- **Currency Handling**: ✅/❌ {NOTES}
- **Audit Trails**: ✅/❌ {NOTES}

## 🧪 Testing Assessment

- **Unit Tests**: {COVERAGE}%
- **Integration Tests**: {STATUS}
- **Test Quality**: {ASSESSMENT}
- **Missing Tests**: {LIST}

## 📚 Documentation Review

- **Code Comments**: {SCORE}%
- **API Documentation**: ✅/❌
- **README Updates**: ✅/❌
- **Architecture Docs**: ✅/❌

## ✅ What's Good

{Highlight positive aspects:}
- Well-structured code
- Good test coverage
- Clear documentation
- Performance improvements
- Security best practices

## 🎯 Recommendations

### Before Merge
1. {Action item with priority}
2. {Action item with priority}

### Post-Merge
1. {Technical debt to address}
2. {Future improvements}

## 📋 Action Items

{LIST_OF_CREATED_ISSUES}
- #{ISSUE_NUMBER}: {ISSUE_TITLE}
- #{ISSUE_NUMBER}: {ISSUE_TITLE}

## 🚦 Review Decision

**Status**: {REQUEST CHANGES | APPROVE | COMMENT}

**Rationale**: {Explanation based on findings}

{If REQUEST CHANGES:}
**Blocking Issues**: {COUNT} CRITICAL + {COUNT} HIGH must be resolved

{If APPROVE:}
**All checks passed!** No critical issues found. Some minor improvements suggested in issues.

---
🤖 **Deep Review powered by Claude Code**
Generated by multi-agent analysis system with:
- Code Quality Guardian 🛡️
- Security Specialist 🔒
- Performance Engineer ⚡
- Compliance Checker 🇮🇩

**Note**: {COUNT} sub-issues have been automatically created for tracking. Review them individually for detailed implementation guidance.
```

Post this summary using:
```bash
gh pr review $PR_NUMBER --comment -F /tmp/pr-review-summary.md
```

### Phase 7: Auto-Create Sub-Issues

For **ALL findings** (as requested by user), create GitHub issues:

```bash
for each finding:
  gh issue create \
    --title "[PR #$PR_NUMBER] [{SEVERITY}] {SHORT_DESCRIPTION}" \
    --body "## Issue from PR Review

**PR**: #{PR_NUMBER} - {PR_TITLE}
**File**: \`{FILE_PATH}\`
**Line**: {LINE_NUMBER}
**Severity**: {SEVERITY}
**Category**: {CODE_QUALITY | SECURITY | PERFORMANCE | COMPLIANCE}
**Agent**: {AGENT_NAME}

### Description
{DETAILED_DESCRIPTION}

### Impact
{IMPACT_ANALYSIS}

### Current Code
\`\`\`{language}
{CODE_SNIPPET}
\`\`\`

### Recommended Fix
{DETAILED_FIX_DESCRIPTION}

### Suggested Code
\`\`\`{language}
{SUGGESTED_CODE}
\`\`\`

### Acceptance Criteria
- [ ] {CRITERION_1}
- [ ] {CRITERION_2}
- [ ] Tests added/updated
- [ ] Documentation updated

### References
- CLAUDE.md: {RELEVANT_SECTION}
- Related PR: #{PR_NUMBER}
- Related files: {OTHER_AFFECTED_FILES}

---
🤖 Auto-generated from PR deep review" \
    --label "pr-review,{severity-label},{category-label}" \
    --assignee "{PR_AUTHOR}"
```

**Labels to use:**
- Severity: `critical`, `high-priority`, `medium-priority`, `low-priority`
- Category: `code-quality`, `security`, `performance`, `compliance`, `testing`, `documentation`
- Type: `bug`, `enhancement`, `technical-debt`

### Phase 8: Final Review Decision

Based on the findings, post the final review decision:

```bash
if CRITICAL_COUNT > 0 || HIGH_COUNT > 3:
    gh pr review $PR_NUMBER --request-changes -b "See detailed review comments. {CRITICAL_COUNT} critical and {HIGH_COUNT} high-priority issues must be addressed."
elif HIGH_COUNT > 0:
    gh pr review $PR_NUMBER --comment -b "Good work! Please address {HIGH_COUNT} high-priority issues before merge. See detailed review and created issues."
else:
    gh pr review $PR_NUMBER --approve -b "Excellent work! All checks passed. Minor improvements tracked in issues."
fi
```

## Important Notes

1. **User Retains Merge Control**: This command NEVER auto-merges. You review the findings and decide when to merge.

2. **Issue Creation**: Unless `--no-issues` flag is present, create issues for ALL findings to ensure proper tracking.

3. **Project-Specific Checks**: Always validate against:
   - Clean architecture layers (Go)
   - Multi-tenancy isolation
   - Indonesian compliance (SAK EP, tax calculations)
   - Security best practices
   - SOLID principles

4. **Agent Coordination**: Launch all 3 agents in parallel using a single message with multiple Task tool calls. Wait for all to complete before synthesizing.

5. **Comment Format**: Use GitHub Markdown for all comments with proper code blocks, tables, and emojis for readability.

6. **Error Handling**: If gh CLI fails or PR doesn't exist, provide clear error message and exit gracefully.

7. **Respect CI Status**: Note CI failures in review but focus on code-level issues.

## Example Usage

```bash
# Review PR by number
/dev:review-pr 123

# Review current branch's PR
/dev:review-pr --current

# Review PR from URL
/dev:review-pr https://github.com/owner/repo/pull/456

# Review without creating issues (just comments)
/dev:review-pr 123 --no-issues

# Auto-detect current branch
/dev:review-pr
```

## Output

You will receive:
- ✅ Line-by-line inline comments on code
- ✅ Comprehensive summary comment with all findings
- ✅ Auto-created GitHub issues for all findings (linked to PR)
- ✅ Review decision (approve/request changes/comment)
- ⚠️ YOU decide when to merge (command does not auto-merge)

Ready to perform deep PR review! 🚀
