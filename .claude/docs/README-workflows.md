# 🚀 Quick Start: GitHub Issue Workflows

## The Fastest Way: Automated End-to-End

```bash
/dev:solve-issue <issue-number>
```

**Example:**
```bash
/dev:solve-issue 382
```

**What happens:**
1. ✅ Fetches issue #382 from GitHub
2. ✅ Creates unique branch: `fix/issue-382-description-a7b3c9f2`
3. ✅ AI implements solution (SOLID principles, tests, docs)
4. ✅ Runs quality checks (lint, test, build)
5. ✅ Commits with proper message
6. ✅ Pushes to GitHub
7. ✅ Creates PR linked to issue
8. ✅ **Issue auto-closes when PR merges**

**Time:** 5-15 minutes vs 2-4 hours manually ⚡

---

## All Available Commands

```bash
# 1. Find issues in codebase
/dev:identify-issues [path] --create-issues

# 2. Solve issue automatically (RECOMMENDED)
/dev:solve-issue <issue-number>

# 3. Create branch + PR, implement manually
/dev:create-fix-pr <issue-numbers>

# 4. Deep PR review with AI agents
/dev:review-pr <pr-number>

# 5. Pre-merge verification
/dev:verify-pr <pr-number>

# 6. Safe merge with cleanup
/dev:merge-pr <pr-number>

# 7. View all workflows
/dev:workflow-guide
```

---

## Quick Comparison

| Command | Speed | Control | Best For |
|---------|-------|---------|----------|
| `/dev:solve-issue` | ⚡⚡⚡⚡⚡ | 🎮 | Bug fixes, small features |
| `/dev:create-fix-pr` | ⚡⚡⚡ | 🎮🎮🎮🎮 | Complex features |
| Manual Git | ⚡ | 🎮🎮🎮🎮🎮 | Learning, exploration |

---

## Example Workflows

### Scenario 1: Fix a Bug
```bash
# One command, done!
/dev:solve-issue 425
```

### Scenario 2: Implement Complex Feature
```bash
# Create structure, implement yourself
/dev:create-fix-pr 401
# ... make changes ...
git add . && git commit -m "feat: implement feature"
git push
/dev:merge-pr --current
```

### Scenario 3: Clean Up Technical Debt
```bash
# Find all issues first
/dev:identify-issues --all --create-issues

# Solve each one
/dev:solve-issue 501
/dev:solve-issue 502
/dev:solve-issue 503
```

### Scenario 4: Review PR → Fix Issues
```bash
# Review creates sub-issues
/dev:review-pr 100

# Solve each sub-issue
/dev:solve-issue 101
/dev:solve-issue 102
```

---

## Branch Naming (Automatic)

Format: `{type}/issue-{number}-{description}-{unique}`

Examples:
- `fix/issue-382-sql-injection-a7b3c9f2`
- `feature/issue-401-export-report-3f9e2b1c`
- `security/issue-425-auth-bypass-e4d8a7b9`

**Benefits:**
- Easy to identify
- Never conflicts
- Auto-sorts by type

---

## Quality Guarantees

All automated commands ensure:

- ✅ SOLID principles
- ✅ Test coverage ≥85%
- ✅ Linting passes
- ✅ Type safety
- ✅ Build succeeds
- ✅ Security checks
- ✅ Indonesian compliance

---

## Common Questions

**Q: Does it auto-merge?**
A: No, unless you use `--auto-merge` flag. You control merging.

**Q: Can I review AI's work?**
A: Yes! Use `/dev:review-pr --current` after implementation.

**Q: What if it makes mistakes?**
A: You can edit the code and commit fixes before merging.

**Q: Can I use with existing issues?**
A: Yes! Works with any GitHub issue number or URL.

**Q: Does it work for features and bugs?**
A: Yes! Auto-detects type from issue labels.

---

## Pro Tips

1. **Start with issue identification:**
   ```bash
   /dev:identify-issues --create-issues
   ```

2. **Use draft mode for review:**
   ```bash
   /dev:solve-issue 382 --draft-only
   ```

3. **Combine with PR review:**
   ```bash
   /dev:solve-issue 382
   /dev:review-pr --current
   ```

4. **Batch related issues:**
   ```bash
   /dev:create-fix-pr 380 381 382
   ```

---

## Need More Details?

```bash
# Full workflow guide
/dev:workflow-guide

# Or read the command files
cat .claude/commands/dev/solve-issue.md
cat .claude/commands/dev/workflow-guide.md
```

---

🤖 **Powered by Claude Code**

Ready to 10x your development speed! 🚀
