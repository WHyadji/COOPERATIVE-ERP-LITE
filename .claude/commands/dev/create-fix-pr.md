---
allowed-tools: Read, Write, Edit, Bash
argument-hint: [issue-numbers...] | --security [issue] | --from-review [pr-number]
description: Creates fix branch and PR that properly references and closes GitHub issues
---

# Create Fix PR from Issues

Automatically creates a feature/fix branch and pull request that properly references and closes GitHub issues. Streamlines the transition from identified issues to fix implementation.

## Context Gathering

Issue Numbers: !`echo "$ARGUMENTS" | grep -oP '\d+' | tr '\n' ' '`

Current Branch: !`git branch --show-current`

Git Status: !`git status --porcelain | wc -l`

Remote: !`git remote get-url origin`

Repository Info: !`gh repo view --json name,owner --jq '{owner: .owner.login, repo: .name}'`

## Your Task

Create a fix branch and pull request for the specified GitHub issues. This command automates the issue-to-PR workflow by:
1. Fetching issue details from GitHub
2. Creating an appropriately named branch
3. Generating a comprehensive PR description with issue links
4. Setting up proper labels, assignees, and metadata

### Phase 1: Parse Arguments and Validate

1. **Extract issue numbers from arguments**
   ```bash
   ISSUE_NUMBERS=$(echo "$ARGUMENTS" | grep -oP '\d+' | tr '\n' ' ')

   if [ -z "$ISSUE_NUMBERS" ]; then
       ERROR: "No issue numbers provided"
       echo "Usage:"
       echo "  /dev:create-fix-pr 45 67 89"
       echo "  /dev:create-fix-pr --security 123"
       echo "  /dev:create-fix-pr --from-review 456"
       EXIT
   fi
   ```

2. **Check for special flags**
   ```bash
   IS_SECURITY=$(echo "$ARGUMENTS" | grep -o "\--security" | wc -l)
   FROM_REVIEW=$(echo "$ARGUMENTS" | grep -oP '\--from-review\s+\K\d+')

   if [ ! -z "$FROM_REVIEW" ]; then
       # Extract issues from PR review comments
       REVIEW_ISSUES=$(gh pr view $FROM_REVIEW --comments | grep -oP 'Issue: #\K\d+' | sort -u | tr '\n' ' ')
       ISSUE_NUMBERS="$REVIEW_ISSUES"
   fi
   ```

3. **Validate current state**
   ```bash
   # Check for uncommitted changes
   if [ $(git status --porcelain | wc -l) -gt 0 ]; then
       WARNING: "You have uncommitted changes"
       echo "Options:"
       echo "  1. Commit them: git add . && git commit -m 'msg'"
       echo "  2. Stash them: git stash"
       echo "  3. Discard them: git checkout ."
       EXIT
   fi

   # Ensure on main/master branch
   CURRENT_BRANCH=$(git branch --show-current)
   if [ "$CURRENT_BRANCH" != "main" ] && [ "$CURRENT_BRANCH" != "master" ] && [ "$CURRENT_BRANCH" != "develop" ]; then
       WARNING: "Not on main branch. Switch to main first"
       echo "Run: git checkout main && git pull"
       EXIT
   fi

   # Pull latest changes
   echo "📥 Pulling latest changes..."
   git pull origin $CURRENT_BRANCH
   ```

### Phase 2: Fetch Issue Details from GitHub

For each issue number, fetch complete details:

```bash
echo "📋 Fetching issue details from GitHub..."

for ISSUE_NUM in $ISSUE_NUMBERS; do
    ISSUE_DATA=$(gh issue view $ISSUE_NUM --json number,title,body,labels,assignees,state)

    if [ -z "$ISSUE_DATA" ]; then
        ERROR: "Issue #$ISSUE_NUM not found"
        continue
    fi

    # Extract issue details
    ISSUE_TITLE=$(echo "$ISSUE_DATA" | jq -r '.title')
    ISSUE_BODY=$(echo "$ISSUE_DATA" | jq -r '.body')
    ISSUE_LABELS=$(echo "$ISSUE_DATA" | jq -r '.labels[].name' | tr '\n' ',' | sed 's/,$//')
    ISSUE_STATE=$(echo "$ISSUE_DATA" | jq -r '.state')

    if [ "$ISSUE_STATE" = "CLOSED" ]; then
        WARNING: "Issue #$ISSUE_NUM is already closed"
    fi

    echo "  ✅ Issue #$ISSUE_NUM: $ISSUE_TITLE"
done
```

### Phase 3: Determine Branch Name and Type

Generate descriptive branch name based on issues:

```bash
# Determine branch prefix based on issue labels/content
BRANCH_PREFIX="fix"

if [ $IS_SECURITY -eq 1 ]; then
    BRANCH_PREFIX="security"
elif echo "$ISSUE_LABELS" | grep -qi "feature"; then
    BRANCH_PREFIX="feature"
elif echo "$ISSUE_LABELS" | grep -qi "enhancement"; then
    BRANCH_PREFIX="enhancement"
elif echo "$ISSUE_LABELS" | grep -qi "refactor"; then
    BRANCH_PREFIX="refactor"
elif echo "$ISSUE_LABELS" | grep -qi "documentation"; then
    BRANCH_PREFIX="docs"
fi

# Create branch name from issues
# Example: fix/issue-45-67-89-tax-calculation-bugs
ISSUE_LIST=$(echo "$ISSUE_NUMBERS" | tr ' ' '-')
BRANCH_SUFFIX=$(echo "$ISSUE_TITLE" | tr '[:upper:]' '[:lower:]' | sed 's/[^a-z0-9]/-/g' | sed 's/--*/-/g' | cut -c1-50 | sed 's/-$//')

BRANCH_NAME="${BRANCH_PREFIX}/issue-${ISSUE_LIST}-${BRANCH_SUFFIX}"

echo "🌿 Branch name: $BRANCH_NAME"
```

### Phase 4: Create Branch

```bash
echo "🌿 Creating branch: $BRANCH_NAME"

# Check if branch already exists
if git show-ref --verify --quiet refs/heads/$BRANCH_NAME; then
    ERROR: "Branch $BRANCH_NAME already exists locally"
    echo "Options:"
    echo "  1. Switch to it: git checkout $BRANCH_NAME"
    echo "  2. Delete it: git branch -D $BRANCH_NAME"
    echo "  3. Use different name"
    EXIT
fi

if git ls-remote --heads origin $BRANCH_NAME | grep -q $BRANCH_NAME; then
    ERROR: "Branch $BRANCH_NAME already exists on remote"
    EXIT
fi

# Create and checkout new branch
git checkout -b $BRANCH_NAME

echo "✅ Branch created and checked out: $BRANCH_NAME"
```

### Phase 5: Generate Comprehensive PR Description

Create detailed PR description that references all issues:

```markdown
# PR Description Template

## 🎯 Purpose

This PR addresses the following issues:

{For each issue:}
- Closes #{ISSUE_NUM}: {ISSUE_TITLE}

## 📋 Issue Details

### Issue #{ISSUE_NUM}: {ISSUE_TITLE}

**Status**: {OPEN/CLOSED}
**Labels**: {LABELS}

**Description**:
{ISSUE_BODY}

**Acceptance Criteria**:
{Extract checklist from issue body if present, or create based on description}

---

{Repeat for each issue}

## 🔧 Implementation Plan

### Changes Required

{Analyze issue descriptions and suggest implementation areas:}

**Backend Changes** (if applicable):
- [ ] Update service layer in `services/{service}/internal/service/`
- [ ] Modify repository layer if database changes needed
- [ ] Add/update domain models in `internal/domain/`
- [ ] Update API handlers if endpoints affected

**Frontend Changes** (if applicable):
- [ ] Update components in `frontend/components/`
- [ ] Modify API client calls in `frontend/lib/api/`
- [ ] Update types in `frontend/types/`
- [ ] Add/update hooks if needed

**Database Changes** (if applicable):
- [ ] Create migration file
- [ ] Update RLS policies
- [ ] Add indexes if needed

**Testing Requirements**:
- [ ] Add unit tests for new/modified functions
- [ ] Add integration tests if API changes
- [ ] Update E2E tests if user-facing changes
- [ ] Verify test coverage ≥ 85%

**Documentation Updates**:
- [ ] Update code comments
- [ ] Update API documentation if applicable
- [ ] Update README if behavior changes
- [ ] Add ADR if architectural decision made

## 🔒 Security Considerations

{If security-related issues:}
- [ ] Validate all inputs
- [ ] Check authentication/authorization
- [ ] Verify multi-tenancy isolation
- [ ] Review for data exposure
- [ ] Check for SQL injection vectors
- [ ] Validate JWT token handling

{Otherwise:}
- [ ] No security implications identified
- [ ] Standard security practices followed

## 🇮🇩 Indonesian Compliance

{If accounting/tax related:}
- [ ] SAK EP standards compliance verified
- [ ] Tax calculation accuracy checked (PPH/PPN)
- [ ] Currency formatting correct (IDR)
- [ ] Audit trail maintained
- [ ] Multi-tenant isolation preserved

{Otherwise:}
- [ ] No compliance implications

## 🧪 Testing Strategy

### Unit Tests
```bash
# Backend
cd services/{affected-service} && go test ./... -v

# Frontend
cd frontend && npm run test
```

### Integration Tests
```bash
cd services/{affected-service} && go test -tags=integration ./... -v
```

### Manual Testing Steps
1. {Step 1 based on issue}
2. {Step 2 based on issue}
3. {Verify expected behavior}

## 📊 SOLID Principles Checklist

- [ ] **Single Responsibility**: Each modified component has single reason to change
- [ ] **Open/Closed**: Changes extend behavior without modifying existing code where possible
- [ ] **Liskov Substitution**: Interface implementations are interchangeable
- [ ] **Interface Segregation**: Interfaces remain cohesive and practical
- [ ] **Dependency Inversion**: Dependencies on abstractions, not concretions

## ✅ Pre-Merge Checklist

- [ ] All tests passing (`/dev:verify-pr --current`)
- [ ] Code quality checks passed (linting, formatting)
- [ ] All builds successful
- [ ] No security vulnerabilities introduced
- [ ] Test coverage ≥ 85%
- [ ] Documentation updated
- [ ] SOLID principles followed
- [ ] Multi-tenancy isolation verified (if applicable)
- [ ] Indonesian compliance checked (if applicable)

## 🔗 Related Links

{For each issue:}
- Issue #{ISSUE_NUM}: https://github.com/{owner}/{repo}/issues/{ISSUE_NUM}

## 📝 Additional Notes

{If security issues:}
⚠️ **SECURITY**: This PR addresses security vulnerabilities. Review carefully before merge.

{If breaking changes detected:}
⚠️ **BREAKING CHANGE**: This PR may introduce breaking changes. Verify compatibility.

---

## 🚀 Next Steps After PR Creation

1. **Implement the fixes**:
   - Make code changes according to implementation plan
   - Follow SOLID principles and project conventions
   - Write tests as you go (TDD approach)

2. **Commit your changes**:
   ```bash
   git add .
   git commit -m "fix: {concise description of fixes}"
   git push origin $BRANCH_NAME
   ```

3. **Request review**:
   - Add reviewers to PR
   - Run `/dev:review-pr --current` for automated deep review
   - Address review feedback

4. **Verify before merge**:
   ```bash
   /dev:verify-pr --current
   ```

5. **Merge when ready**:
   ```bash
   /dev:merge-pr --current
   ```

---

🤖 **PR created by Claude Code `/dev:create-fix-pr` command**

This PR will automatically close the linked issues when merged.
```

### Phase 6: Create Draft PR

```bash
echo "📝 Creating draft PR..."

# Save PR description to file
echo "$PR_DESCRIPTION" > /tmp/fix-pr-description.md

# Determine PR title
if [ $(echo "$ISSUE_NUMBERS" | wc -w) -eq 1 ]; then
    PR_TITLE="Fix: $ISSUE_TITLE"
else
    PR_TITLE="Fix: Multiple issues - $(echo $ISSUE_NUMBERS | sed 's/ /, #/g' | sed 's/^/#/')"
fi

# Create draft PR
gh pr create \
    --draft \
    --title "$PR_TITLE" \
    --body-file /tmp/fix-pr-description.md \
    --base main

PR_NUMBER=$(gh pr list --head $BRANCH_NAME --json number --jq '.[0].number')
PR_URL=$(gh pr view $PR_NUMBER --json url --jq '.url')

echo "✅ Draft PR created: #$PR_NUMBER"
echo "   $PR_URL"
```

### Phase 7: Add Labels to PR

```bash
echo "🏷️  Adding labels to PR..."

# Collect unique labels from all issues
ALL_LABELS=$(for ISSUE_NUM in $ISSUE_NUMBERS; do
    gh issue view $ISSUE_NUM --json labels --jq '.labels[].name'
done | sort -u | tr '\n' ',' | sed 's/,$//')

# Add issue-derived labels
if [ ! -z "$ALL_LABELS" ]; then
    gh pr edit $PR_NUMBER --add-label "$ALL_LABELS"
fi

# Add fix-specific labels
if [ $IS_SECURITY -eq 1 ]; then
    gh pr edit $PR_NUMBER --add-label "security"
fi

gh pr edit $PR_NUMBER --add-label "ready-for-implementation"

echo "✅ Labels added"
```

### Phase 8: Link Issues to PR

```bash
echo "🔗 Linking issues to PR..."

# GitHub automatically links when PR description contains "Closes #123"
# But we'll also add comments to issues for visibility

for ISSUE_NUM in $ISSUE_NUMBERS; do
    gh issue comment $ISSUE_NUM --body "🔧 Fix PR created: #$PR_NUMBER

This PR will close this issue when merged.

**Next Steps**:
1. Implement the fix in branch \`$BRANCH_NAME\`
2. Push changes and mark PR as ready for review
3. Run \`/dev:verify-pr $PR_NUMBER\` before requesting merge

View PR: $PR_URL"

    echo "  ✅ Linked issue #$ISSUE_NUM to PR #$PR_NUMBER"
done
```

### Phase 9: Set Up Assignees

```bash
echo "👤 Setting up assignees..."

# Get current GitHub user
CURRENT_USER=$(gh api user --jq '.login')

# Assign PR to current user
gh pr edit $PR_NUMBER --add-assignee "$CURRENT_USER"

# Also collect assignees from issues
ISSUE_ASSIGNEES=$(for ISSUE_NUM in $ISSUE_NUMBERS; do
    gh issue view $ISSUE_NUM --json assignees --jq '.assignees[].login'
done | sort -u | tr '\n' ',' | sed 's/,$//')

if [ ! -z "$ISSUE_ASSIGNEES" ]; then
    gh pr edit $PR_NUMBER --add-assignee "$ISSUE_ASSIGNEES"
fi

echo "✅ Assignees set"
```

### Phase 10: Print Summary and Next Steps

```bash
echo ""
echo "════════════════════════════════════════════════════════════"
echo "✅ FIX PR CREATED SUCCESSFULLY"
echo "════════════════════════════════════════════════════════════"
echo ""
echo "📋 Summary:"
echo "  Branch:  $BRANCH_NAME"
echo "  PR:      #$PR_NUMBER (Draft)"
echo "  URL:     $PR_URL"
echo "  Issues:  $(echo $ISSUE_NUMBERS | sed 's/ /, #/g' | sed 's/^/#/')"
echo "  Labels:  $ALL_LABELS"
echo ""
echo "🎯 What's Next:"
echo ""
echo "1. 🔧 IMPLEMENT THE FIX:"
echo "   - You're now on branch: $BRANCH_NAME"
echo "   - Make your code changes"
echo "   - Follow the implementation plan in the PR description"
echo "   - Write tests as you go (TDD approach)"
echo ""
echo "2. 💾 COMMIT YOUR CHANGES:"
echo "   git add ."
echo "   git commit -m \"fix: <description>\""
echo "   git push origin $BRANCH_NAME"
echo ""
echo "3. 🔍 GET AUTOMATED REVIEW (Optional):"
echo "   /dev:review-pr $PR_NUMBER"
echo ""
echo "4. ✅ VERIFY BEFORE MERGE:"
echo "   /dev:verify-pr $PR_NUMBER"
echo ""
echo "5. 🎉 MARK AS READY FOR REVIEW:"
echo "   gh pr ready $PR_NUMBER"
echo ""
echo "6. 🚀 MERGE WHEN APPROVED:"
echo "   /dev:merge-pr $PR_NUMBER"
echo ""
echo "════════════════════════════════════════════════════════════"
echo ""
echo "💡 Tips:"
echo "  - PR is in DRAFT mode - mark ready when implementation done"
echo "  - Issues will auto-close when PR is merged"
echo "  - Use /dev:verify-pr before requesting final review"
echo "  - Follow SOLID principles and project conventions"
echo ""
```

## Command Variations

### Basic Usage
```bash
# Fix single issue
/dev:create-fix-pr 123

# Fix multiple related issues
/dev:create-fix-pr 45 67 89

# Fix security issue
/dev:create-fix-pr --security 234

# Create PR from review findings
/dev:create-fix-pr --from-review 456
```

### Advanced Usage
```bash
# With custom branch prefix
/dev:create-fix-pr 123 --type feature

# Skip issue validation (for external issues)
/dev:create-fix-pr 123 --no-validate

# Create non-draft PR immediately
/dev:create-fix-pr 123 --ready
```

## Example Workflow

```bash
# 1. Issues identified from review
/dev:review-pr 100  # Creates issues #101, #102, #103

# 2. Create fix PR for those issues
/dev:create-fix-pr 101 102 103

# Output:
#   ✅ Branch created: fix/issue-101-102-103-tax-calculation-errors
#   ✅ Draft PR created: #104
#   ✅ Issues linked to PR

# 3. Implement fixes
# (make your code changes...)

git add .
git commit -m "fix: resolve tax calculation errors"
git push origin fix/issue-101-102-103-tax-calculation-errors

# 4. Mark ready and verify
gh pr ready 104
/dev:verify-pr 104

# 5. Merge when approved
/dev:merge-pr 104

# Result: Issues #101, #102, #103 automatically closed ✅
```

## Notes

- **Draft by default**: PRs are created as drafts so you can implement first
- **Auto-linking**: Uses "Closes #123" syntax to auto-close issues on merge
- **Smart naming**: Branch names are descriptive and include issue numbers
- **Comprehensive description**: PR includes all context from linked issues
- **Label inheritance**: PR inherits labels from issues for consistency
- **Assignee tracking**: Assigns to current user and issue assignees
- **Implementation guidance**: PR description includes implementation plan and checklists

This command bridges the gap between issue identification and fix implementation! 🚀
