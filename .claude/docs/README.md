# 📚 Claude Code Documentation

Comprehensive documentation for Claude Code workflows and commands.

## 📖 Quick Start

**New to the workflow system?** Start here:
- [README-workflows.md](./README-workflows.md) - Quick reference card for GitHub workflows

**Want to solve an issue?** Read this:
- [ANSWER-your-question.md](./ANSWER-your-question.md) - How to solve issues from PR reviews

---

## 🔍 Documentation Index

### GitHub Workflow Documentation

#### Main Workflows
- **[workflow-guide.md](./workflow-guide.md)** (11KB)
  - Complete guide to all GitHub issue resolution workflows
  - Comparison of automated vs manual approaches
  - Best practices and pro tips
  - Time-saving strategies

- **[README-workflows.md](./README-workflows.md)** (3.7KB)
  - Quick reference card
  - Most common commands
  - Fast lookup for workflows
  - Command cheat sheet

#### Specific Use Cases
- **[ANSWER-your-question.md](./ANSWER-your-question.md)** (5.5KB)
  - How to solve issues created from PR reviews
  - Example: Issue #382 from PR #188 review
  - Why separate PRs are better
  - Complete workflow explanation

- **[EXAMPLE-pr-review-workflow.md](./EXAMPLE-pr-review-workflow.md)** (13KB)
  - Visual workflow diagram
  - Step-by-step example with real numbers
  - What gets linked where
  - Benefits vs alternative approaches

### Development Workflows

#### Test-Driven Development
- **[tdd-implement-docs.md](./tdd-implement-docs.md)** (7.6KB)
  - TDD workflow with `/implement:implement-tdd`
  - Red-Green-Refactor cycle
  - Integration with testing agents
  - Best practices for TDD

#### Code Refactoring
- **[solid-refactor-docs.md](./solid-refactor-docs.md)** (8.2KB)
  - SOLID principles refactoring guide
  - Using `/plan:solid-refactor`
  - Step-by-step refactoring process
  - Clean code principles

#### Bug Fixing
- **[bugfix-todo-docs.md](./bugfix-todo-docs.md)** (6.4KB)
  - Bug fixing workflow with `/plan:bugfix-plan`
  - Systematic debugging approach
  - Todo list generation for fixes
  - Severity assessment

### Parallel Work
- **[PARALLEL_WORK.md](./PARALLEL_WORK.md)** (6.6KB)
  - Working on multiple tasks simultaneously
  - Branch management strategies
  - Context switching best practices

- **[HOW_TO_PARALLEL_WORK.md](./HOW_TO_PARALLEL_WORK.md)** (29KB)
  - Detailed guide to parallel development
  - Advanced Git workflows
  - Team collaboration patterns
  - Conflict resolution

---

## 🎯 Common Use Cases

### I want to...

#### Fix a GitHub Issue
```bash
# Read first:
cat .claude/docs/README-workflows.md

# Then run:
/dev:solve-issue <issue-number>
```

#### Solve an Issue from PR Review
```bash
# Read first:
cat .claude/docs/ANSWER-your-question.md

# Example: Issue #382 from PR #188 review
/dev:solve-issue 382
# Automatically detects parent PR and links everything
```

#### Understand All Workflow Options
```bash
# Comprehensive guide:
cat .claude/docs/workflow-guide.md
```

#### Implement with TDD
```bash
# Read guide:
cat .claude/docs/tdd-implement-docs.md

# Then use:
/implement:implement-tdd <file_path>
```

#### Refactor Code Following SOLID
```bash
# Read guide:
cat .claude/docs/solid-refactor-docs.md

# Then use:
/plan:solid-refactor <file_path>
```

#### Fix a Bug Systematically
```bash
# Read guide:
cat .claude/docs/bugfix-todo-docs.md

# Then use:
/plan:bugfix-plan <bug_description> <file_path>
```

---

## 📊 Documentation by Topic

### Workflows & Automation
1. [workflow-guide.md](./workflow-guide.md) - All GitHub workflows
2. [README-workflows.md](./README-workflows.md) - Quick reference
3. [ANSWER-your-question.md](./ANSWER-your-question.md) - PR review sub-issues
4. [EXAMPLE-pr-review-workflow.md](./EXAMPLE-pr-review-workflow.md) - Visual examples

### Code Quality
1. [solid-refactor-docs.md](./solid-refactor-docs.md) - SOLID refactoring
2. [bugfix-todo-docs.md](./bugfix-todo-docs.md) - Bug fixing
3. [tdd-implement-docs.md](./tdd-implement-docs.md) - Test-driven development

### Team Collaboration
1. [PARALLEL_WORK.md](./PARALLEL_WORK.md) - Parallel workflows
2. [HOW_TO_PARALLEL_WORK.md](./HOW_TO_PARALLEL_WORK.md) - Advanced parallel work

---

## 🚀 Quick Command Reference

### GitHub Issue Workflows
```bash
/dev:identify-issues [path] --create-issues   # Find and create issues
/dev:solve-issue <issue-number>               # Solve issue end-to-end
/dev:create-fix-pr <issue-numbers>            # Create PR for manual fix
/dev:review-pr <pr-number>                    # Deep PR review
/dev:verify-pr <pr-number>                    # Pre-merge checks
/dev:merge-pr <pr-number>                     # Safe merge
```

### Development Workflows
```bash
/implement:implement-tdd <file_path>          # TDD implementation
/plan:solid-refactor <file_path>              # Plan SOLID refactor
/plan:bugfix-plan <bug> <file_path>           # Plan bug fix
```

### Other Commands
```bash
/dev:api                                      # Create API endpoint
/dev:component                                # Create React component
/security:security-audit                      # Security audit
/refactor:optimize-code                       # Optimize code
```

---

## 📈 Workflow Efficiency Comparison

| Task | Manual Time | With Commands | Time Saved |
|------|-------------|---------------|------------|
| Fix GitHub Issue | 2-4 hours | 5-15 min | 90%+ |
| PR Review + Fix | 3-6 hours | 15-30 min | 85%+ |
| TDD Implementation | 1-2 hours | 10-20 min | 80%+ |
| SOLID Refactoring | 4-8 hours | 30-60 min | 85%+ |
| Bug Fix Planning | 30-60 min | 5-10 min | 80%+ |

---

## 🎓 Learning Path

### Beginner
1. Start with [README-workflows.md](./README-workflows.md) for quick overview
2. Read [ANSWER-your-question.md](./ANSWER-your-question.md) for common scenario
3. Try `/dev:solve-issue` on a simple issue

### Intermediate
1. Read [workflow-guide.md](./workflow-guide.md) for all options
2. Explore [EXAMPLE-pr-review-workflow.md](./EXAMPLE-pr-review-workflow.md)
3. Use `/dev:review-pr` and solve sub-issues

### Advanced
1. Read [HOW_TO_PARALLEL_WORK.md](./HOW_TO_PARALLEL_WORK.md)
2. Study [solid-refactor-docs.md](./solid-refactor-docs.md)
3. Combine multiple workflows for complex tasks

---

## 💡 Pro Tips

### Fastest Workflow
```bash
# For most issues, just run:
/dev:solve-issue <issue-number>
# Done! 90% faster than manual.
```

### Quality Assurance
```bash
# After AI implementation, review:
/dev:review-pr --current
/dev:verify-pr --current
```

### Batch Processing
```bash
# Identify all issues first:
/dev:identify-issues --all --create-issues

# Then solve one by one:
/dev:solve-issue 501
/dev:solve-issue 502
/dev:solve-issue 503
```

### PR Review Workflow
```bash
# Review creates issues:
/dev:review-pr 188

# Solve each issue (auto-detects parent PR):
/dev:solve-issue 382
/dev:solve-issue 383
/dev:solve-issue 384
```

---

## 🔗 Related Resources

### Project Documentation
- [CLAUDE.md](../../CLAUDE.md) - Project overview and guidelines
- [.claude/commands/](../commands/) - All available commands

### External Links
- [Conventional Commits](https://www.conventionalcommits.org/)
- [SOLID Principles](https://en.wikipedia.org/wiki/SOLID)
- [GitHub CLI](https://cli.github.com/)

---

## 📝 Documentation Standards

All documentation in this directory follows:
- **Clear structure** - Headings, sections, examples
- **Code examples** - Real, runnable commands
- **Visual aids** - Diagrams, tables, emojis
- **Cross-references** - Links to related docs
- **Quick start** - TL;DR at the top
- **Complete coverage** - Deep dives for advanced users

---

## 🤝 Contributing

To add new documentation:

1. Create `.md` file in `.claude/docs/`
2. Follow existing format and structure
3. Add entry to this README index
4. Include practical examples
5. Cross-reference related docs

---

## 📞 Getting Help

- Read the relevant documentation above
- Try the examples in the docs
- Check command help: `cat .claude/commands/dev/<command>.md`
- Run `/help` in Claude Code

---

## 🎯 Summary

**Total Documentation:** 9 files, ~100KB

**Coverage:**
- ✅ GitHub workflows (4 docs)
- ✅ Development practices (3 docs)
- ✅ Team collaboration (2 docs)

**Time Savings:** 80-90% faster development with automated workflows

**Start Here:** [README-workflows.md](./README-workflows.md)

---

🤖 **Powered by Claude Code**

Your complete guide to efficient development workflows! 🚀
