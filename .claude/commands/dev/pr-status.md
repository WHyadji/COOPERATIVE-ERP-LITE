---
allowed-tools: Bash
argument-hint: [pr-number] | --current
description: Quick PR health check showing CI status, reviews, conflicts, and linked issues
---

# PR Status Quick Check

Displays a quick health check dashboard for a PR: CI status, review status, merge conflicts, linked issues, and overall readiness.

## Context

PR Number: !`
if [ "$ARGUMENTS" = "--current" ] || [ -z "$ARGUMENTS" ]; then
    git branch --show-current | xargs -I {} gh pr list --head {} --json number --jq '.[0].number'
else
    echo "$ARGUMENTS" | grep -oP '\d+'
fi
`

## Your Task

```bash
PR_NUM=$PR_NUMBER

# Fetch PR data
PR_DATA=$(gh pr view $PR_NUM --json number,title,state,isDraft,mergeable,additions,deletions,reviewDecision,headRefName,url)

# Parse data
TITLE=$(echo "$PR_DATA" | jq -r '.title')
STATE=$(echo "$PR_DATA" | jq -r '.state')
IS_DRAFT=$(echo "$PR_DATA" | jq -r '.isDraft')
MERGEABLE=$(echo "$PR_DATA" | jq -r '.mergeable')
REVIEW_STATUS=$(echo "$PR_DATA" | jq -r '.reviewDecision')
BRANCH=$(echo "$PR_DATA" | jq -r '.headRefName')
URL=$(echo "$PR_DATA" | jq -r '.url')

# Get CI status
CI_DATA=$(gh pr checks $PR_NUM --json name,state)
CI_PASS=$(echo "$CI_DATA" | jq '[.[] | select(.state == "SUCCESS")] | length')
CI_FAIL=$(echo "$CI_DATA" | jq '[.[] | select(.state == "FAILURE")] | length')
CI_PENDING=$(echo "$CI_DATA" | jq '[.[] | select(.state == "PENDING" or .state == "IN_PROGRESS")] | length')

# Get linked issues
LINKED_ISSUES=$(gh pr view $PR_NUM --json body --jq '.body' | grep -oP '(Closes|Fixes|Resolves) #\K\d+' | wc -l)

# Display dashboard
echo ""
echo "╔════════════════════════════════════════════════════════╗"
echo "║  PR #$PR_NUM Health Check                            "
echo "╠════════════════════════════════════════════════════════╣"
echo "║"
echo "║  📋 $TITLE"
echo "║  🔗 $URL"
echo "║"
echo "╠════════════════════════════════════════════════════════╣"
echo "║  STATUS"
echo "╠════════════════════════════════════════════════════════╣"
echo "║"

# State
if [ "$STATE" = "OPEN" ]; then
    if [ "$IS_DRAFT" = "true" ]; then
        echo "║  🟡 Draft"
    else
        echo "║  🟢 Open"
    fi
else
    echo "║  🔴 $STATE"
fi

# Mergeable
if [ "$MERGEABLE" = "MERGEABLE" ]; then
    echo "║  ✅ No conflicts"
elif [ "$MERGEABLE" = "CONFLICTING" ]; then
    echo "║  ❌ Has merge conflicts"
else
    echo "║  ⚠️  Mergeable: $MERGEABLE"
fi

# Reviews
if [ "$REVIEW_STATUS" = "APPROVED" ]; then
    echo "║  ✅ Approved"
elif [ "$REVIEW_STATUS" = "CHANGES_REQUESTED" ]; then
    echo "║  ❌ Changes requested"
elif [ "$REVIEW_STATUS" = "REVIEW_REQUIRED" ]; then
    echo "║  ⏳ Review pending"
else
    echo "║  ⚪ No reviews"
fi

echo "║"
echo "╠════════════════════════════════════════════════════════╣"
echo "║  CI/CD CHECKS"
echo "╠════════════════════════════════════════════════════════╣"
echo "║"
echo "║  ✅ Passed:  $CI_PASS"
echo "║  ❌ Failed:  $CI_FAIL"
echo "║  ⏳ Pending: $CI_PENDING"
echo "║"
echo "╠════════════════════════════════════════════════════════╣"
echo "║  METADATA"
echo "╠════════════════════════════════════════════════════════╣"
echo "║"
echo "║  📦 Branch: $BRANCH"
echo "║  🔗 Linked Issues: $LINKED_ISSUES"
echo "║  📊 Changes: +$(echo "$PR_DATA" | jq -r '.additions') / -$(echo "$PR_DATA" | jq -r '.deletions')"
echo "║"
echo "╠════════════════════════════════════════════════════════╣"
echo "║  READY TO MERGE?"
echo "╠════════════════════════════════════════════════════════╣"
echo "║"

# Overall readiness
READY=true
BLOCKERS=""

if [ "$STATE" != "OPEN" ]; then
    READY=false
    BLOCKERS="$BLOCKERS PR not open."
fi

if [ "$IS_DRAFT" = "true" ]; then
    READY=false
    BLOCKERS="$BLOCKERS Still draft."
fi

if [ "$MERGEABLE" != "MERGEABLE" ]; then
    READY=false
    BLOCKERS="$BLOCKERS Has conflicts."
fi

if [ $CI_FAIL -gt 0 ]; then
    READY=false
    BLOCKERS="$BLOCKERS CI failing."
fi

if [ "$REVIEW_STATUS" != "APPROVED" ]; then
    READY=false
    BLOCKERS="$BLOCKERS Not approved."
fi

if [ "$READY" = true ]; then
    echo "║  ✅ YES - Ready to merge!"
    echo "║"
    echo "║  Next: /dev:merge-pr $PR_NUM"
else
    echo "║  ❌ NO - Blockers:"
    echo "║  $BLOCKERS"
    echo "║"
    echo "║  Actions:"
    if [ "$MERGEABLE" != "MERGEABLE" ]; then
        echo "║  - Resolve conflicts"
    fi
    if [ $CI_FAIL -gt 0 ]; then
        echo "║  - Fix failing CI checks"
    fi
    if [ "$REVIEW_STATUS" != "APPROVED" ]; then
        echo "║  - Get review approval"
    fi
    if [ "$IS_DRAFT" = "true" ]; then
        echo "║  - Mark as ready: gh pr ready $PR_NUM"
    fi
fi

echo "║"
echo "╚════════════════════════════════════════════════════════╝"
echo ""
```

## Usage

```bash
/dev:pr-status 123        # Check PR #123
/dev:pr-status --current  # Check current branch's PR
```
