---
allowed-tools: Read, Write, Edit, Bash, Task, Grep, Glob
argument-hint: [issue-number] | [issue-url] | --auto-merge | --draft-only
description: End-to-end workflow to solve GitHub issue - branch, implement, commit, push, PR
---

# Solve Issue - Complete Workflow

Automated end-to-end workflow to solve a GitHub issue from start to finish:
1. Fetch issue details from GitHub
2. Create uniquely named branch
3. Implement the fix/feature
4. Run tests and quality checks
5. Commit with proper message
6. Push and create PR
7. Link PR to issue for auto-close

## Context Gathering

**Arguments:** `$ARGUMENTS`

**Parse issue reference:**
```bash
ISSUE_INPUT=$(echo "$ARGUMENTS" | sed 's/--[a-z-]*//g' | xargs)
AUTO_MERGE=$(echo "$ARGUMENTS" | grep -c "\--auto-merge")
DRAFT_ONLY=$(echo "$ARGUMENTS" | grep -c "\--draft-only")

# Extract issue number from URL or direct number
if echo "$ISSUE_INPUT" | grep -q "github.com"; then
    ISSUE_NUMBER=$(echo "$ISSUE_INPUT" | grep -oP '/issues/\K\d+')
else
    ISSUE_NUMBER="$ISSUE_INPUT"
fi
```

**Current state:**
```bash
gh --version
git branch --show-current
git status --porcelain | wc -l
pwd
```

## Your Task

You will orchestrate a complete workflow to solve issue #$ISSUE_NUMBER from creation to PR, implementing clean code following SOLID principles.

### Phase 1: Validate and Gather Issue Details

1. **Check prerequisites**
   ```bash
   # Check GitHub CLI authenticated
   gh auth status || (echo "❌ GitHub CLI not authenticated. Run: gh auth login" && exit 1)

   # Check uncommitted changes
   DIRTY=$(git status --porcelain | wc -l)
   if [ $DIRTY -gt 0 ]; then
       echo "⚠️  WARNING: You have $DIRTY uncommitted changes"
       echo "Options:"
       echo "  1. Stash: git stash"
       echo "  2. Commit: git add . && git commit -m 'WIP'"
       echo "  3. Discard: git checkout ."
       read -p "Continue anyway? (y/N) " -n 1 -r
       if [[ ! $REPLY =~ ^[Yy]$ ]]; then
           exit 1
       fi
   fi

   # Check on main branch
   CURRENT_BRANCH=$(git branch --show-current)
   if [ "$CURRENT_BRANCH" != "main" ] && [ "$CURRENT_BRANCH" != "master" ]; then
       echo "⚠️  Not on main branch. Switching to main..."
       git checkout main || git checkout master
   fi

   # Pull latest changes
   echo "📥 Pulling latest changes..."
   git pull
   ```

2. **Fetch issue details from GitHub**
   ```bash
   echo "📋 Fetching issue #$ISSUE_NUMBER from GitHub..."

   ISSUE_DATA=$(gh issue view $ISSUE_NUMBER --json number,title,body,labels,state,assignees,milestone,url 2>&1)

   if echo "$ISSUE_DATA" | grep -q "Could not resolve"; then
       echo "❌ ERROR: Issue #$ISSUE_NUMBER not found"
       exit 1
   fi

   ISSUE_TITLE=$(echo "$ISSUE_DATA" | jq -r '.title')
   ISSUE_BODY=$(echo "$ISSUE_DATA" | jq -r '.body // ""')
   ISSUE_STATE=$(echo "$ISSUE_DATA" | jq -r '.state')
   ISSUE_LABELS=$(echo "$ISSUE_DATA" | jq -r '.labels[].name' | tr '\n' ',' | sed 's/,$//')
   ISSUE_URL=$(echo "$ISSUE_DATA" | jq -r '.url')

   echo "✅ Issue found:"
   echo "   Title: $ISSUE_TITLE"
   echo "   State: $ISSUE_STATE"
   echo "   Labels: $ISSUE_LABELS"

   if [ "$ISSUE_STATE" = "CLOSED" ]; then
       echo "⚠️  WARNING: Issue is already closed"
       read -p "Continue anyway? (y/N) " -n 1 -r
       if [[ ! $REPLY =~ ^[Yy]$]]; then
           exit 1
       fi
   fi

   # Detect if this is a sub-issue from PR review
   PARENT_PR=$(echo "$ISSUE_TITLE" | grep -oP '\[PR #\K\d+(?=\])' || echo "")
   if [ ! -z "$PARENT_PR" ]; then
       echo ""
       echo "🔗 This is a sub-issue from PR #$PARENT_PR review"
       echo "   Original PR: https://github.com/$(gh repo view --json nameWithOwner -q .nameWithOwner)/pull/$PARENT_PR"

       # Fetch parent PR context
       PARENT_PR_DATA=$(gh pr view $PARENT_PR --json title,body,url,headRefName 2>&1)
       if ! echo "$PARENT_PR_DATA" | grep -q "Could not resolve"; then
           PARENT_PR_TITLE=$(echo "$PARENT_PR_DATA" | jq -r '.title')
           PARENT_PR_URL=$(echo "$PARENT_PR_DATA" | jq -r '.url')
           PARENT_PR_BRANCH=$(echo "$PARENT_PR_DATA" | jq -r '.headRefName')

           echo "   Parent PR Title: $PARENT_PR_TITLE"
           echo "   Parent PR Branch: $PARENT_PR_BRANCH"
           echo ""
           echo "💡 This fix will be created in a separate PR that references the parent PR #$PARENT_PR"
       fi
   fi
   ```

### Phase 2: Create Smart Branch Name

Generate descriptive, unique branch name based on issue type and content:

```bash
echo "🌿 Generating branch name..."

# Determine branch prefix from labels
BRANCH_PREFIX="fix"

if echo "$ISSUE_LABELS" | grep -qi "feature"; then
    BRANCH_PREFIX="feature"
elif echo "$ISSUE_LABELS" | grep -qi "security\|critical"; then
    BRANCH_PREFIX="security"
elif echo "$ISSUE_LABELS" | grep -qi "enhancement"; then
    BRANCH_PREFIX="enhancement"
elif echo "$ISSUE_LABELS" | grep -qi "refactor\|technical-debt"; then
    BRANCH_PREFIX="refactor"
elif echo "$ISSUE_LABELS" | grep -qi "documentation\|docs"; then
    BRANCH_PREFIX="docs"
elif echo "$ISSUE_LABELS" | grep -qi "performance"; then
    BRANCH_PREFIX="perf"
elif echo "$ISSUE_LABELS" | grep -qi "test"; then
    BRANCH_PREFIX="test"
fi

# Create slug from title (max 40 chars, lowercase, alphanumeric + hyphens)
TITLE_SLUG=$(echo "$ISSUE_TITLE" | \
    sed 's/\[.*\]//g' | \
    tr '[:upper:]' '[:lower:]' | \
    sed 's/[^a-z0-9]/-/g' | \
    sed 's/--*/-/g' | \
    sed 's/^-\|-$//g' | \
    cut -c1-40 | \
    sed 's/-$//')

# Generate unique suffix (8-char hash from timestamp + issue number)
UNIQUE_SUFFIX=$(echo "$ISSUE_NUMBER-$(date +%s)" | md5sum | cut -c1-8)

# Construct final branch name
BRANCH_NAME="${BRANCH_PREFIX}/issue-${ISSUE_NUMBER}-${TITLE_SLUG}-${UNIQUE_SUFFIX}"

echo "✅ Branch name: $BRANCH_NAME"
echo "   Format: ${BRANCH_PREFIX}/issue-${ISSUE_NUMBER}-<description>-<unique-id>"

# Verify branch doesn't exist
if git show-ref --verify --quiet refs/heads/$BRANCH_NAME; then
    echo "❌ ERROR: Branch already exists locally: $BRANCH_NAME"
    echo "This should never happen due to unique suffix. Please report this!"
    exit 1
fi

if git ls-remote --heads origin $BRANCH_NAME 2>/dev/null | grep -q $BRANCH_NAME; then
    echo "❌ ERROR: Branch already exists on remote: $BRANCH_NAME"
    exit 1
fi
```

### Phase 3: Create Branch and Initial Commit

```bash
echo "🌿 Creating branch: $BRANCH_NAME"
git checkout -b $BRANCH_NAME

echo "✅ Branch created and checked out"

# Create initial tracking commit (empty - just to establish branch)
git commit --allow-empty -m "chore: initialize branch for issue #${ISSUE_NUMBER}

This branch tracks work on:
${ISSUE_TITLE}

Issue: ${ISSUE_URL}

🤖 Generated with Claude Code /dev:solve-issue"

echo "✅ Initial tracking commit created"
```

### Phase 4: Implement the Fix/Feature

This is the core implementation phase using AI-assisted development:

```bash
echo ""
echo "════════════════════════════════════════════════════════════"
echo "🔧 IMPLEMENTATION PHASE"
echo "════════════════════════════════════════════════════════════"
echo ""
```

**Now use the Task tool to launch implementation agent:**

```
Task: Implement the solution for GitHub issue #${ISSUE_NUMBER}

**Issue Title**: ${ISSUE_TITLE}

**Issue Description**:
${ISSUE_BODY}

**Labels**: ${ISSUE_LABELS}

**Project Context**:
- This is a multi-tenant accounting web application
- Backend: Go microservices (api-gateway, accounting-ledger, business-service)
- Frontend: Next.js with TypeScript
- Must follow SOLID principles (see /Users/adji/Documents/VISI-DIGITAL-TERPADU/accounting-webapp/CLAUDE.md)
- Must maintain test coverage ≥85%
- Must comply with Indonesian accounting standards (SAK EP)

**Your Implementation Requirements**:

1. **Analyze the issue** thoroughly:
   - Understand the problem or feature request
   - Identify affected files and services
   - Determine implementation approach

2. **Follow SOLID principles**:
   - Single Responsibility: Each component does one thing
   - Open/Closed: Extend behavior without modifying existing code
   - Liskov Substitution: Implementations are interchangeable
   - Interface Segregation: Cohesive, practical interfaces
   - Dependency Inversion: Depend on abstractions

3. **Implement with TDD approach**:
   - Write tests FIRST (use /implement:implement-tdd if helpful)
   - Make tests pass
   - Refactor while keeping tests green
   - Target ≥85% coverage

4. **Code Quality**:
   - Follow naming conventions from CLAUDE.md
   - Add proper error handling
   - Include documentation comments
   - Keep functions small and focused (<20 lines typically)

5. **Security Considerations**:
   - Validate all inputs
   - Check authentication/authorization
   - Maintain multi-tenant isolation (company_id filters)
   - Prevent SQL injection, XSS, etc.

6. **Make incremental commits** as you work:
   - Commit after each logical unit of work
   - Use conventional commit format:
     - feat: new feature
     - fix: bug fix
     - refactor: code restructuring
     - test: adding tests
     - docs: documentation
     - chore: maintenance
   - Example: "feat: add PPH21 tax calculation service"

7. **Final verification**:
   - All tests passing
   - No linting errors
   - Code builds successfully
   - Changes address the issue completely

**DO NOT use the TodoWrite tool** - focus entirely on implementation.

**Return a summary** of:
- Files created/modified
- Tests added
- Key implementation decisions
- Any concerns or trade-offs
- Verification results (tests, lint, build)
```

**Agent Type**: `clean-code-generator` (for new code) or `code-quality-guardian` (for fixes/refactoring)

### Phase 5: Run Quality Checks

After implementation completes, run comprehensive quality checks:

```bash
echo ""
echo "════════════════════════════════════════════════════════════"
echo "✅ QUALITY CHECKS"
echo "════════════════════════════════════════════════════════════"
echo ""

# Determine what changed
CHANGED_FILES=$(git diff --name-only main...HEAD)
HAS_GO=$(echo "$CHANGED_FILES" | grep -c "\.go$" || true)
HAS_TS=$(echo "$CHANGED_FILES" | grep -c "\.\(ts\|tsx\)$" || true)

CHECKS_PASSED=true

# Go backend checks
if [ $HAS_GO -gt 0 ]; then
    echo "🔍 Running Go checks..."

    # Find affected service
    SERVICE_DIR=$(echo "$CHANGED_FILES" | grep "\.go$" | head -1 | cut -d'/' -f1-2)

    if [ -d "$SERVICE_DIR" ]; then
        cd "$SERVICE_DIR"

        # Format check
        echo "  📝 Checking format..."
        UNFORMATTED=$(gofmt -l . | grep -v "^vendor" | wc -l)
        if [ $UNFORMATTED -gt 0 ]; then
            echo "  ⚠️  Running gofmt..."
            gofmt -w .
            git add .
            git commit -m "style: apply gofmt" || true
        fi

        # Linting
        echo "  🔍 Running linter..."
        if ! golangci-lint run ./... 2>&1; then
            echo "  ❌ Linting failed"
            CHECKS_PASSED=false
        fi

        # Tests
        echo "  🧪 Running tests..."
        if ! go test ./... -v -cover; then
            echo "  ❌ Tests failed"
            CHECKS_PASSED=false
        fi

        # Build
        echo "  🔨 Building..."
        if ! go build ./...; then
            echo "  ❌ Build failed"
            CHECKS_PASSED=false
        fi

        cd - > /dev/null
    fi
fi

# Frontend checks
if [ $HAS_TS -gt 0 ]; then
    echo "🔍 Running TypeScript checks..."

    # Determine frontend or admin
    FRONTEND_DIR="frontend"
    if echo "$CHANGED_FILES" | grep -q "^admin-frontend/"; then
        FRONTEND_DIR="admin-frontend"
    fi

    if [ -d "$FRONTEND_DIR" ]; then
        cd "$FRONTEND_DIR"

        # Linting
        echo "  🔍 Running ESLint..."
        if ! npm run lint 2>&1; then
            echo "  ⚠️  Attempting auto-fix..."
            npm run lint -- --fix || true
            git add .
            git commit -m "style: apply ESLint fixes" || true
        fi

        # Type check
        echo "  📝 Type checking..."
        if ! npx tsc --noEmit; then
            echo "  ❌ Type check failed"
            CHECKS_PASSED=false
        fi

        # Tests
        echo "  🧪 Running tests..."
        if ! npm run test 2>&1; then
            echo "  ❌ Tests failed"
            CHECKS_PASSED=false
        fi

        # Build
        echo "  🔨 Building..."
        if ! npm run build 2>&1; then
            echo "  ❌ Build failed"
            CHECKS_PASSED=false
        fi

        cd - > /dev/null
    fi
fi

if [ "$CHECKS_PASSED" = false ]; then
    echo ""
    echo "❌ Quality checks failed. Please fix the issues before proceeding."
    echo ""
    echo "You can:"
    echo "  1. Fix issues manually"
    echo "  2. Commit fixes: git add . && git commit -m 'fix: address quality issues'"
    echo "  3. Re-run checks"
    echo "  4. Push when ready: git push origin $BRANCH_NAME"
    exit 1
fi

echo "✅ All quality checks passed!"
```

### Phase 6: Create Final Commit and Push

```bash
echo ""
echo "════════════════════════════════════════════════════════════"
echo "📤 PREPARING FOR PUSH"
echo "════════════════════════════════════════════════════════════"
echo ""

# Check if there are uncommitted changes
if [ $(git status --porcelain | wc -l) -gt 0 ]; then
    echo "📝 Creating final commit for remaining changes..."

    # Generate commit message based on issue type
    COMMIT_TYPE="fix"
    if echo "$ISSUE_LABELS" | grep -qi "feature"; then
        COMMIT_TYPE="feat"
    elif echo "$ISSUE_LABELS" | grep -qi "refactor"; then
        COMMIT_TYPE="refactor"
    elif echo "$ISSUE_LABELS" | grep -qi "docs"; then
        COMMIT_TYPE="docs"
    elif echo "$ISSUE_LABELS" | grep -qi "test"; then
        COMMIT_TYPE="test"
    fi

    # Create meaningful commit message
    COMMIT_MSG="${COMMIT_TYPE}: ${ISSUE_TITLE}

${ISSUE_BODY}

Closes #${ISSUE_NUMBER}

🤖 Generated with Claude Code /dev:solve-issue"

    git add .
    echo "$COMMIT_MSG" | git commit -F -

    echo "✅ Final commit created"
fi

# Push to remote
echo "📤 Pushing to remote..."
git push -u origin $BRANCH_NAME

echo "✅ Pushed to origin/$BRANCH_NAME"
```

### Phase 7: Create Pull Request

```bash
echo ""
echo "════════════════════════════════════════════════════════════"
echo "📋 CREATING PULL REQUEST"
echo "════════════════════════════════════════════════════════════"
echo ""

# Generate PR title
PR_TITLE="${ISSUE_TITLE}"

# Generate comprehensive PR description
# Include parent PR reference if this is a sub-issue from review
PARENT_PR_SECTION=""
if [ ! -z "$PARENT_PR" ]; then
    PARENT_PR_SECTION="
## 🔗 Related to PR #${PARENT_PR}

This issue was identified during review of PR #${PARENT_PR}:
- **Parent PR**: ${PARENT_PR_URL}
- **Parent PR Title**: ${PARENT_PR_TITLE}
- **Parent PR Branch**: \`${PARENT_PR_BRANCH}\`

This fix addresses one of the issues found in that review and should be:
1. Reviewed independently
2. Merged into \`main\` separately
3. The parent PR #${PARENT_PR} can then be updated/rebased if needed

**Note**: Merging this PR will automatically close issue #${ISSUE_NUMBER}, which was created from the review of PR #${PARENT_PR}.
"
fi

PR_BODY="## 🎯 Purpose

This PR resolves issue #${ISSUE_NUMBER}:
**${ISSUE_TITLE}**

Closes #${ISSUE_NUMBER}
${PARENT_PR_SECTION}
## 📋 Issue Details

**Link**: ${ISSUE_URL}
**Labels**: ${ISSUE_LABELS}

### Description
${ISSUE_BODY}

## 🔧 Implementation

### Changes Made
$(git log --oneline main..HEAD | sed 's/^/- /')

### Files Changed
\`\`\`
$(git diff --stat main...HEAD)
\`\`\`

### Key Decisions
- Followed SOLID principles for maintainability
- Maintained test coverage ≥85%
- Ensured multi-tenant isolation where applicable
- Added comprehensive error handling

## 🧪 Testing

### Test Coverage
- [x] Unit tests added/updated
- [x] Integration tests added/updated (if applicable)
- [x] Test coverage ≥85%
- [x] All tests passing

### Manual Testing
- [x] Tested locally
- [x] Verified issue resolution
- [x] Checked for regressions

## ✅ Quality Checks

- [x] Code formatted (gofmt/prettier)
- [x] Linting passed (golangci-lint/ESLint)
- [x] Type checking passed (TypeScript)
- [x] Build successful
- [x] SOLID principles followed
- [x] Error handling comprehensive
- [x] Documentation updated

## 🔒 Security Checklist

- [x] Input validation implemented
- [x] Authentication/authorization checked
- [x] Multi-tenant isolation maintained
- [x] No sensitive data exposed
- [x] SQL injection prevented
- [x] XSS vulnerabilities addressed

## 🇮🇩 Indonesian Compliance

$(if echo "$ISSUE_LABELS" | grep -qi "accounting\|tax\|ledger"; then
    echo "- [x] SAK EP compliance verified"
    echo "- [x] Tax calculations accurate"
    echo "- [x] Audit trail maintained"
    echo "- [x] Currency handling correct (IDR)"
else
    echo "- [ ] N/A - No compliance implications"
fi)

## 📚 Documentation

- [x] Code comments added
- [x] API documentation updated (if applicable)
- [x] README updated (if applicable)

## 🚀 Deployment Notes

No special deployment steps required.

## 🔗 Related

- Issue: #${ISSUE_NUMBER}
- Branch: \`${BRANCH_NAME}\`

---

🤖 **Auto-generated PR** using Claude Code \`/dev:solve-issue\` workflow
"

# Save PR body to file
echo "$PR_BODY" > /tmp/pr-body-$ISSUE_NUMBER.md

# Determine if draft or ready
DRAFT_FLAG=""
if [ $DRAFT_ONLY -eq 1 ]; then
    DRAFT_FLAG="--draft"
    echo "📝 Creating DRAFT pull request..."
else
    echo "📝 Creating pull request..."
fi

# Create PR
gh pr create \
    $DRAFT_FLAG \
    --title "$PR_TITLE" \
    --body-file /tmp/pr-body-$ISSUE_NUMBER.md \
    --base main

# Get PR number
PR_NUMBER=$(gh pr list --head $BRANCH_NAME --json number --jq '.[0].number')
PR_URL=$(gh pr view $PR_NUMBER --json url --jq '.url')

echo "✅ Pull request created: #$PR_NUMBER"
echo "   $PR_URL"

# Add labels from issue to PR
if [ ! -z "$ISSUE_LABELS" ]; then
    echo "🏷️  Adding labels to PR..."
    gh pr edit $PR_NUMBER --add-label "$ISSUE_LABELS"
fi

# Add assignees from issue
ISSUE_ASSIGNEES=$(gh issue view $ISSUE_NUMBER --json assignees --jq '.assignees[].login' | tr '\n' ',' | sed 's/,$//')
if [ ! -z "$ISSUE_ASSIGNEES" ]; then
    echo "👤 Adding assignees to PR..."
    gh pr edit $PR_NUMBER --add-assignee "$ISSUE_ASSIGNEES"
else
    # Assign to current user
    CURRENT_USER=$(gh api user --jq '.login')
    gh pr edit $PR_NUMBER --add-assignee "$CURRENT_USER"
fi

# Link PR to issue via comment
PARENT_PR_REF=""
if [ ! -z "$PARENT_PR" ]; then
    PARENT_PR_REF="

**Context**: This issue was created from review of PR #${PARENT_PR}
**Parent PR**: ${PARENT_PR_URL}"
fi

gh issue comment $ISSUE_NUMBER --body "🔧 **PR Created**: #${PR_NUMBER}

This PR implements the solution for this issue.

**PR**: ${PR_URL}
**Branch**: \`${BRANCH_NAME}\`${PARENT_PR_REF}

The issue will automatically close when the PR is merged.

---
🤖 Auto-linked by Claude Code /dev:solve-issue"

echo "✅ PR linked to issue"

# Also comment on parent PR if this is a sub-issue fix
if [ ! -z "$PARENT_PR" ]; then
    echo "🔗 Commenting on parent PR #$PARENT_PR..."
    gh pr comment $PARENT_PR --body "✅ **Sub-Issue Fixed**: #${ISSUE_NUMBER}

A fix PR has been created for one of the issues found in this review:

**Issue**: #${ISSUE_NUMBER} - ${ISSUE_TITLE}
**Fix PR**: #${PR_NUMBER}
**Fix PR URL**: ${PR_URL}

This fix will be merged independently. Once merged, you may want to:
- Rebase this PR if needed
- Re-run the review to check if issue is resolved
- Continue with other review findings

---
🤖 Auto-linked by Claude Code /dev:solve-issue"
    echo "✅ Commented on parent PR #$PARENT_PR"
fi
```

### Phase 8: Optional Auto-Merge

```bash
if [ $AUTO_MERGE -eq 1 ]; then
    echo ""
    echo "════════════════════════════════════════════════════════════"
    echo "🚀 AUTO-MERGE ENABLED"
    echo "════════════════════════════════════════════════════════════"
    echo ""

    echo "⏳ Waiting for CI checks..."
    gh pr checks $PR_NUMBER --watch

    CI_STATE=$(gh pr checks $PR_NUMBER --json state --jq '.[0].state')

    if [ "$CI_STATE" = "SUCCESS" ]; then
        echo "✅ CI checks passed"

        # Request approval or auto-approve if owner
        echo "🔀 Merging PR..."
        gh pr merge $PR_NUMBER --squash --auto

        echo "✅ PR will auto-merge when approved"
    else
        echo "❌ CI checks failed. Not merging."
        echo "   Review PR: $PR_URL"
        exit 1
    fi
else
    echo ""
    echo "ℹ️  Auto-merge not enabled. PR is ready for review."
fi
```

### Phase 9: Print Summary and Next Steps

```bash
echo ""
echo "════════════════════════════════════════════════════════════"
echo "✅ WORKFLOW COMPLETE"
echo "════════════════════════════════════════════════════════════"
echo ""
echo "📋 Summary:"
echo "   Issue:   #${ISSUE_NUMBER} - ${ISSUE_TITLE}"
echo "   Branch:  ${BRANCH_NAME}"
echo "   PR:      #${PR_NUMBER} $([ $DRAFT_ONLY -eq 1 ] && echo '(DRAFT)' || echo '(READY)')"
echo "   URL:     ${PR_URL}"
echo "   Labels:  ${ISSUE_LABELS}"
echo ""
echo "✅ Completed Steps:"
echo "   [✓] Fetched issue from GitHub"
echo "   [✓] Created unique branch"
echo "   [✓] Implemented solution"
echo "   [✓] Ran quality checks (lint, test, build)"
echo "   [✓] Committed changes"
echo "   [✓] Pushed to remote"
echo "   [✓] Created pull request"
echo "   [✓] Linked PR to issue"
echo ""
echo "🎯 What's Next:"
echo ""
if [ $DRAFT_ONLY -eq 1 ]; then
    echo "1. 📝 MARK PR AS READY:"
    echo "   gh pr ready $PR_NUMBER"
    echo ""
fi
if [ $AUTO_MERGE -eq 0 ]; then
    echo "1. 🔍 REVIEW PR (optional automated review):"
    echo "   /dev:review-pr $PR_NUMBER"
    echo ""
    echo "2. ✅ VERIFY BEFORE MERGE:"
    echo "   /dev:verify-pr $PR_NUMBER"
    echo ""
    echo "3. 🚀 MERGE WHEN APPROVED:"
    echo "   /dev:merge-pr $PR_NUMBER"
    echo "   OR: gh pr merge $PR_NUMBER --squash"
    echo ""
fi
echo "════════════════════════════════════════════════════════════"
echo ""
echo "💡 Tips:"
echo "  - Issue #${ISSUE_NUMBER} will auto-close when PR merges"
echo "  - Branch follows naming: ${BRANCH_PREFIX}/issue-NNN-description-<unique>"
echo "  - Unique suffix prevents branch name conflicts"
echo "  - All commits follow conventional commit format"
echo "  - Quality checks already passed (lint, test, build)"
echo ""
```

## Usage Examples

### Basic Usage
```bash
# Solve issue by number (most common)
/dev:solve-issue 382

# Solve issue from URL
/dev:solve-issue https://github.com/owner/repo/issues/382

# Create draft PR only (manual merge later)
/dev:solve-issue 382 --draft-only

# Solve and auto-merge when CI passes
/dev:solve-issue 382 --auto-merge
```

### Complete Workflow Example
```bash
# 1. Identify issues first (optional)
/dev:identify-issues services/accounting-ledger/ --create-issues
# Creates issues: #380, #381, #382

# 2. Solve issue #382
/dev:solve-issue 382

# Output:
#   ✅ Branch: fix/issue-382-increase-test-coverage-a7b3c9f2
#   ✅ Implementation complete
#   ✅ Quality checks passed
#   ✅ PR created: #400
#   ✅ Linked to issue #382

# 3. Issue #382 auto-closes when PR #400 merges!
```

### Integration with Existing Commands
```bash
# Workflow 1: Review → Create Issues → Solve
/dev:review-pr 100                    # Reviews PR, creates issues
/dev:solve-issue 101                  # Solve first issue
/dev:solve-issue 102                  # Solve second issue

# Workflow 2: Identify → Create Fix PR → Manual Implementation
/dev:identify-issues --create-issues  # Scan and create issues
/dev:create-fix-pr 103 104 105        # Create PR for multiple issues
# Then manually implement...

# Workflow 3: Solve → Auto Review → Merge
/dev:solve-issue 106                  # Full auto implementation
/dev:review-pr --current              # Deep automated review
/dev:merge-pr --current               # Merge when ready
```

### Special Case: Solving Sub-Issues from PR Reviews

When `/dev:review-pr` creates sub-issues, they have titles like:
```
[PR #188] [CRITICAL] Increase test coverage to 85% minimum
```

The `/dev:solve-issue` command automatically detects these and:
1. ✅ Extracts the parent PR number (188)
2. ✅ Fetches parent PR context
3. ✅ Creates a separate fix PR (not modifying the parent PR)
4. ✅ References parent PR in the new PR description
5. ✅ Comments on both the issue AND the parent PR
6. ✅ Links everything together for traceability

**Example:**
```bash
# Step 1: Review PR #188
/dev:review-pr 188

# Output:
# Found 5 issues, created:
#   - Issue #382: [PR #188] [CRITICAL] Increase test coverage to 85%
#   - Issue #383: [PR #188] [HIGH] Fix SOLID violation in AccountService
#   - Issue #384: [PR #188] [MEDIUM] Add input validation

# Step 2: Solve issue #382 (it knows it's from PR #188)
/dev:solve-issue 382

# What happens:
# ✅ Detects parent PR: #188
# ✅ Creates branch: test/issue-382-increase-test-coverage-a7b3c9f2
# ✅ Implements test coverage improvements
# ✅ Creates NEW PR #400 (separate from PR #188)
# ✅ PR #400 description references PR #188
# ✅ Comments on PR #188: "Sub-issue #382 fixed in PR #400"
# ✅ Comments on issue #382: "Fixed by PR #400, related to PR #188"

# Step 3: When PR #400 merges
# ✅ Issue #382 auto-closes
# ✅ PR #188 now has one less issue to address
# ✅ Can continue solving other issues from PR #188 review
```

**Why separate PRs?**
- Each fix can be reviewed independently
- Fixes can be merged even if parent PR #188 isn't ready
- Parent PR #188 can be rebased/updated after fixes merge
- Better Git history and bisectability
- Easier to track which review finding was addressed

**Full workflow:**
```bash
# 1. Review a PR
/dev:review-pr 188
# Creates issues: #382, #383, #384, #385, #386

# 2. Solve CRITICAL issues first
/dev:solve-issue 382  # Creates PR #400
/dev:solve-issue 383  # Creates PR #401

# 3. Merge the fix PRs
/dev:merge-pr 400
/dev:merge-pr 401

# 4. Now PR #188 author can:
# - Rebase PR #188 on main (gets the fixes)
# - Or close PR #188 if fixes covered everything
# - Or continue with remaining issues

# 5. Solve remaining issues
/dev:solve-issue 384  # Creates PR #402
/dev:solve-issue 385  # Creates PR #403
```

## Branch Naming Format

The command creates branches following this pattern:
```
{type}/issue-{number}-{description}-{unique-id}
```

**Examples:**
- `fix/issue-382-increase-test-coverage-a7b3c9f2`
- `feature/issue-401-add-tax-report-export-3f9e2b1c`
- `security/issue-425-fix-sql-injection-e4d8a7b9`
- `refactor/issue-500-implement-solid-principles-9c2f1a6d`

**Benefits:**
- Easy to identify purpose at a glance
- Issue number clearly visible
- Description provides context
- Unique 8-char suffix prevents conflicts
- Auto-sorts by type in Git tools

## Commit Message Format

All commits follow [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <description>

<body>

Closes #<issue-number>

🤖 Generated with Claude Code /dev:solve-issue
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `refactor`: Code restructuring
- `test`: Adding tests
- `docs`: Documentation
- `style`: Formatting
- `perf`: Performance improvement
- `chore`: Maintenance

## Notes

- **Fully Automated**: From issue to PR in one command
- **Quality Guaranteed**: Runs linting, tests, and build before push
- **SOLID Principles**: AI agent enforces clean code practices
- **Auto-Linking**: PR automatically closes issue on merge
- **Unique Branches**: Never conflicts with existing branches
- **Test Coverage**: Maintains ≥85% coverage requirement
- **Comprehensive PR**: Includes all checklists and documentation
- **Flexible**: Can be draft-only or auto-merge
- **Safe**: Validates state before making changes

This command provides the most efficient workflow from issue to resolution! 🚀
