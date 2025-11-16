---
allowed-tools: Bash
argument-hint: [pr-number] [issue-numbers...]
description: Link existing PR to GitHub issues for auto-closing on merge
---

# Link Issues to Existing PR

Links GitHub issues to an existing PR by updating the PR description with "Closes #N" syntax, enabling automatic issue closing on merge.

## Context

PR Number: !`echo "$ARGUMENTS" | grep -oP '^\d+' | head -1`
Issue Numbers: !`echo "$ARGUMENTS" | grep -oP '\d+' | tail -n +2 | tr '\n' ' '`

## Your Task

```bash
PR_NUM=$(echo "$ARGUMENTS" | awk '{print $1}')
ISSUE_NUMS=$(echo "$ARGUMENTS" | awk '{for(i=2;i<=NF;i++)print $i}' | tr '\n' ' ')

if [ -z "$PR_NUM" ] || [ -z "$ISSUE_NUMS" ]; then
    echo "Usage: /dev:link-issues PR_NUMBER ISSUE_NUMBERS..."
    echo "Example: /dev:link-issues 123 45 67 89"
    exit 1
fi

# Get current PR body
CURRENT_BODY=$(gh pr view $PR_NUM --json body --jq '.body')

# Add issue links
NEW_BODY="$CURRENT_BODY

## Linked Issues

$(for ISSUE in $ISSUE_NUMS; do echo "Closes #$ISSUE"; done)"

# Update PR
echo "$NEW_BODY" | gh pr edit $PR_NUM --body-file -

echo "✅ Linked issues to PR #$PR_NUM: $ISSUE_NUMS"
echo "Issues will auto-close when PR is merged"
```
