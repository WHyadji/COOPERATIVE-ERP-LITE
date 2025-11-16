---
allowed-tools: Read, Write, Bash
argument-hint: [pr-number] | --current | --squash | --rebase | --skip-verify
description: Safely merge PR with final validation, issue closing verification, and cleanup
---

# Smart PR Merge Command

Safely merges a PR with final validation checks, proper merge strategy, automatic branch cleanup, and issue closing verification.

## Context

PR Number: !`
if [ "$ARGUMENTS" = "--current" ] || [ -z "$ARGUMENTS" ]; then
    git branch --show-current | xargs -I {} gh pr list --head {} --json number --jq '.[0].number'
else
    echo "$ARGUMENTS" | grep -oP '\d+'
fi
`

PR Details: !`gh pr view $PR_NUMBER --json number,title,state,mergeable,isDraft,additions,deletions,reviews,reviewDecision,headRefName`

Linked Issues: !`gh pr view $PR_NUMBER --json body --jq '.body' | grep -oP '(Closes|Fixes|Resolves) #\K\d+'`

## Your Task

Safely merge PR #$PR_NUMBER with comprehensive pre-merge validation, proper merge strategy, and post-merge cleanup.

### Phase 1: Pre-Merge Validation

```bash
# Check merge strategy from arguments
MERGE_STRATEGY="merge"
if echo "$ARGUMENTS" | grep -q "\--squash"; then
    MERGE_STRATEGY="squash"
elif echo "$ARGUMENTS" | grep -q "\--rebase"; then
    MERGE_STRATEGY="rebase"
fi

SKIP_VERIFY=$(echo "$ARGUMENTS" | grep -c "\--skip-verify")

# 1. Validate PR state
if [ "$PR_STATE" != "OPEN" ]; then
    ERROR: "PR #$PR_NUMBER is not open (state: $PR_STATE)"
    EXIT
fi

if [ "$IS_DRAFT" = "true" ]; then
    ERROR: "PR #$PR_NUMBER is still a draft. Mark as ready first: gh pr ready $PR_NUMBER"
    EXIT
fi

if [ "$MERGEABLE" != "MERGEABLE" ]; then
    ERROR: "PR #$PR_NUMBER has merge conflicts or is not mergeable"
    echo "Resolve conflicts and try again"
    EXIT
fi

# 2. Check review status
if [ "$REVIEW_DECISION" != "APPROVED" ]; then
    WARNING: "PR #$PR_NUMBER is not approved (status: $REVIEW_DECISION)"
    echo "Continue anyway? (y/N)"
    # Prompt for confirmation
fi

# 3. Run verification unless skipped
if [ $SKIP_VERIFY -eq 0 ]; then
    echo "🔍 Running pre-merge verification..."
    /dev:verify-pr $PR_NUMBER

    if [ $? -ne 0 ]; then
        ERROR: "Verification failed. Fix issues or use --skip-verify to override"
        EXIT
    fi
else
    WARNING: "Skipping verification (--skip-verify flag used)"
fi

# 4. Check CI status
CI_FAILED=$(gh pr checks $PR_NUMBER --json state --jq '[.[] | select(.state != "SUCCESS" and .state != "SKIPPED")] | length')

if [ $CI_FAILED -gt 0 ]; then
    ERROR: "$CI_FAILED CI checks failed"
    gh pr checks $PR_NUMBER
    EXIT
fi
```

### Phase 2: Verify Issue Links

```bash
echo "🔗 Verifying linked issues..."

LINKED_ISSUES=$(gh pr view $PR_NUMBER --json body --jq '.body' | grep -oP '(Closes|Fixes|Resolves) #\K\d+' | sort -u)

if [ -z "$LINKED_ISSUES" ]; then
    WARNING: "No issues linked to this PR"
    echo "Issues that should be closed: (enter numbers separated by space, or press Enter to skip)"
    # Optionally prompt for issues
else
    echo "Issues that will be closed on merge:"
    for ISSUE in $LINKED_ISSUES; do
        ISSUE_TITLE=$(gh issue view $ISSUE --json title --jq '.title' 2>/dev/null)
        if [ $? -eq 0 ]; then
            echo "  ✅ #$ISSUE: $ISSUE_TITLE"
        else
            WARNING: "Issue #$ISSUE not found or inaccessible"
        fi
    done
fi
```

### Phase 3: Final Confirmation

```bash
echo ""
echo "═══════════════════════════════════════"
echo "🚀 READY TO MERGE"
echo "═══════════════════════════════════════"
echo "PR:       #$PR_NUMBER - $PR_TITLE"
echo "Strategy: $MERGE_STRATEGY"
echo "Branch:   $HEAD_REF_NAME"
echo "Changes:  +$ADDITIONS / -$DELETIONS"
echo "Issues:   $(echo $LINKED_ISSUES | wc -w) will be closed"
echo "═══════════════════════════════════════"
echo ""
echo "Proceed with merge? (Y/n)"
# Auto-proceed unless user configured to prompt
```

### Phase 4: Execute Merge

```bash
echo "🔀 Merging PR #$PR_NUMBER..."

case $MERGE_STRATEGY in
    squash)
        gh pr merge $PR_NUMBER --squash --delete-branch
        ;;
    rebase)
        gh pr merge $PR_NUMBER --rebase --delete-branch
        ;;
    *)
        gh pr merge $PR_NUMBER --merge --delete-branch
        ;;
esac

MERGE_STATUS=$?

if [ $MERGE_STATUS -ne 0 ]; then
    ERROR: "Merge failed"
    EXIT
fi

echo "✅ PR #$PR_NUMBER merged successfully"
```

### Phase 5: Post-Merge Verification

```bash
echo "🔍 Verifying post-merge state..."

# 1. Verify branch deleted
if git ls-remote --heads origin $HEAD_REF_NAME | grep -q $HEAD_REF_NAME; then
    WARNING: "Remote branch $HEAD_REF_NAME still exists"
    echo "Delete manually: git push origin --delete $HEAD_REF_NAME"
fi

# Cleanup local branch if exists
if git show-ref --verify --quiet refs/heads/$HEAD_REF_NAME; then
    git checkout main 2>/dev/null || git checkout master
    git branch -D $HEAD_REF_NAME
    echo "  ✅ Local branch deleted"
fi

# 2. Verify issues closed
sleep 2  # Give GitHub time to process
for ISSUE in $LINKED_ISSUES; do
    ISSUE_STATE=$(gh issue view $ISSUE --json state --jq '.state' 2>/dev/null)
    if [ "$ISSUE_STATE" = "CLOSED" ]; then
        echo "  ✅ Issue #$ISSUE closed"
    else
        WARNING: "Issue #$ISSUE not yet closed (state: $ISSUE_STATE)"
        echo "     Check: https://github.com/$(gh repo view --json nameWithOwner --jq '.nameWithOwner')/issues/$ISSUE"
    fi
done

# 3. Update local main branch
echo "📥 Updating local main branch..."
git checkout main 2>/dev/null || git checkout master
git pull origin main 2>/dev/null || git pull origin master
```

### Phase 6: Generate Merge Summary

```markdown
# ✅ Merge Complete

**PR**: #{PR_NUMBER} - {PR_TITLE}
**Merged**: {TIMESTAMP}
**Strategy**: {MERGE_STRATEGY}
**Changes**: +{ADDITIONS} / -{DELETIONS} lines

## Closed Issues

{For each linked issue:}
- ✅ #{ISSUE_NUMBER}: {ISSUE_TITLE}

## Post-Merge Actions

- ✅ PR merged to main
- ✅ Feature branch deleted
- ✅ Issues closed automatically
- ✅ Local main branch updated

## Next Steps

1. **Verify deployment**: Check that changes deploy correctly
2. **Monitor logs**: Watch for any runtime issues
3. **Update changelog**: Run `/docs:add-changelog` if needed
4. **Notify team**: Inform stakeholders of the merge

---

🎉 **Great work!** The fixes are now in main and issues are resolved.
```

Post summary to PR as final comment:

```bash
echo "$MERGE_SUMMARY" > /tmp/merge-summary.md
gh pr comment $PR_NUMBER -F /tmp/merge-summary.md
```

### Phase 7: Optional Changelog Update

```bash
echo ""
echo "📝 Update changelog? (y/N)"
# If user wants:
# /docs:add-changelog "Merged PR #$PR_NUMBER: $PR_TITLE"
```

## Usage Examples

```bash
# Standard merge (creates merge commit)
/dev:merge-pr 123

# Squash merge (single commit)
/dev:merge-pr 123 --squash

# Rebase merge (linear history)
/dev:merge-pr 123 --rebase

# Merge current branch's PR
/dev:merge-pr --current

# Skip verification (not recommended)
/dev:merge-pr 123 --skip-verify

# Squash merge current PR
/dev:merge-pr --current --squash
```

## Notes

- **Safety first**: Runs `/dev:verify-pr` before merging (unless `--skip-verify`)
- **Smart defaults**: Uses merge commit strategy by default
- **Auto-cleanup**: Deletes feature branch automatically
- **Issue tracking**: Verifies all linked issues close properly
- **Rollback friendly**: Merge commits make rollback easier than squash
- **Team notification**: Posts merge summary for visibility

**Recommended merge strategies:**
- `--merge`: Feature branches, collaborative work (preserves history)
- `--squash`: Single-developer features, cleanup commits (cleaner history)
- `--rebase`: Linear history preference, small changes (advanced users)
