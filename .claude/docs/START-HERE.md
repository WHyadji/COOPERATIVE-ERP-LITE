# 🚀 Start Here - Quick Guide to Workflow System

Welcome! This is your entry point to the automated GitHub issue resolution system.

---

## ⚡ Super Quick Start (30 seconds)

**To solve any GitHub issue:**
```bash
/dev:solve-issue <issue-number>
```

**That's it!** The command does everything:
- Creates branch
- Implements solution (AI)
- Runs tests
- Commits & pushes
- Creates PR
- Links to issue

**Time:** 5-15 minutes vs 2-4 hours manually

---

## 📖 What to Read First

### Just want to solve an issue?
→ [README-workflows.md](./README-workflows.md) (3 min read)

### Issue is from a PR review (e.g., Issue #382 from PR #188)?
→ [ANSWER-your-question.md](./ANSWER-your-question.md) (5 min read)

### Want to understand all workflow options?
→ [workflow-guide.md](./workflow-guide.md) (15 min read)

### Need visual examples?
→ [EXAMPLE-pr-review-workflow.md](./EXAMPLE-pr-review-workflow.md) (10 min read)

### Want complete documentation index?
→ [README.md](./README.md) (master index)

---

## 🎯 Common Scenarios

### Scenario 1: Simple Bug Fix
```bash
/dev:solve-issue 425
# Done! ✅
```

### Scenario 2: Issue from PR Review
```bash
# Issue #382 has title: "[PR #188] [CRITICAL] ..."
/dev:solve-issue 382
# Automatically detects it's from PR #188
# Creates separate fix PR
# Links everything together
# Done! ✅
```

### Scenario 3: Multiple Related Issues
```bash
/dev:solve-issue 501  # Fix first
/dev:solve-issue 502  # Fix second
/dev:solve-issue 503  # Fix third
# Each gets its own PR ✅
```

### Scenario 4: Quality Check Your Work
```bash
/dev:solve-issue 382        # AI implements
/dev:review-pr --current    # AI reviews
/dev:verify-pr --current    # Run all checks
/dev:merge-pr --current     # Merge when ready
```

---

## 📚 Documentation Overview

```
.claude/docs/
├── START-HERE.md ← YOU ARE HERE
├── README-workflows.md        Quick reference (3.7KB)
├── ANSWER-your-question.md    PR review issues (5.5KB)
├── workflow-guide.md          Complete guide (11KB)
├── EXAMPLE-pr-review-workflow.md  Visual examples (13KB)
└── README.md                  Full index (8.4KB)
```

**Total time to read everything:** ~1 hour
**Time to become productive:** ~20 minutes

---

## 🔑 Key Commands

```bash
# Solve any issue
/dev:solve-issue <issue-number>

# Find issues first
/dev:identify-issues --create-issues

# Review a PR (creates sub-issues)
/dev:review-pr <pr-number>

# Verify before merge
/dev:verify-pr <pr-number>

# Safe merge
/dev:merge-pr <pr-number>
```

---

## ❓ FAQ

**Q: Does it work for issues from PR reviews?**
A: Yes! Automatically detects parent PR and links everything.

**Q: Does it auto-merge?**
A: No, you control merging (unless you use `--auto-merge` flag).

**Q: Can I review AI's work?**
A: Yes! Use `/dev:review-pr --current` after implementation.

**Q: What if I want to implement myself?**
A: Use `/dev:create-fix-pr <issue>` for manual implementation.

**Q: How much time does it save?**
A: 90%+ time savings (5-15 min vs 2-4 hours).

---

## 🎓 Learning Path

1. **5 minutes:** Read [README-workflows.md](./README-workflows.md)
2. **5 minutes:** Try `/dev:solve-issue` on a simple issue
3. **5 minutes:** Review the created PR
4. **10 minutes:** Read [ANSWER-your-question.md](./ANSWER-your-question.md) for PR review scenario
5. **30 minutes:** Explore [workflow-guide.md](./workflow-guide.md) for advanced usage

**Total:** 1 hour to mastery

---

## 💡 Pro Tip

**For your exact scenario** (Issue #382 from PR #188 review):

```bash
/dev:solve-issue 382
```

The command will:
1. Detect it's from PR #188
2. Create separate fix PR #400
3. Reference PR #188 in PR #400
4. Comment on both #382 and #188
5. Link everything perfectly

**Result:** Clean, traceable, independent fix in 5-15 minutes! 🎉

---

## 🚀 Next Steps

1. **Try it now:**
   ```bash
   /dev:solve-issue <your-issue-number>
   ```

2. **Read for your case:**
   - Simple fix: [README-workflows.md](./README-workflows.md)
   - PR review issue: [ANSWER-your-question.md](./ANSWER-your-question.md)

3. **Master the system:**
   - [workflow-guide.md](./workflow-guide.md)
   - [README.md](./README.md)

---

🤖 **Powered by Claude Code**

Ready to 10x your development speed! 🚀
