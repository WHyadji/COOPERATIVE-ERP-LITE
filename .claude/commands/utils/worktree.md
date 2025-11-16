---
allowed-tools: Bash, Read, Write
argument-hint: [action] [name] [base-branch] | create | remove | list | switch
description: Manage git worktrees for parallel development
---

# Git Worktree Manager

Manage git worktrees for parallel agent development: **$ARGUMENTS**

## Usage Examples:
- `/utils:worktree list` - Show all worktrees
- `/utils:worktree create agent-4` - Create new worktree from main
- `/utils:worktree create issue-525 main` - Create worktree from specific branch
- `/utils:worktree remove agent-1` - Remove worktree
- `/utils:worktree switch agent-2` - Switch to worktree directory
- `/utils:worktree` - Show status and help

## Instructions:

You are a git worktree manager. When this command is invoked:

1. **Validate git repository** - Ensure we're in a git repository
2. **Parse the command arguments** to determine the action:
   - `list` - List all existing worktrees with details
   - `create <name> [base-branch]` - Create new worktree
   - `remove <name>` - Remove existing worktree
   - `switch <name>` - Show command to switch to worktree
   - No args - Show current status and quick help

## Worktree Management:

### List Worktrees
Show all existing worktrees with:
- Path
- Current branch
- Latest commit
- Status (clean/dirty)

Use: `git worktree list`

### Create Worktree
```bash
# Create worktree with new branch from base
git worktree add -b <branch-name> ../accounting-webapp-worktrees/<name> <base-branch>

# Or use existing branch
git worktree add ../accounting-webapp-worktrees/<name> <branch-name>
```

**Defaults:**
- Base branch: `main` (if not specified)
- Location: `../accounting-webapp-worktrees/<name>`
- Branch name: Same as worktree name

**Validation:**
- Check if worktree name already exists
- Check if base branch exists
- Confirm worktree creation success
- Show path and branch information

### Remove Worktree
```bash
# Check for uncommitted changes first
cd <worktree-path>
git status

# Remove worktree
git worktree remove <worktree-path> [--force]
```

**Safety checks:**
- Warn if uncommitted changes exist
- Ask for confirmation if dirty
- Use `--force` only if user confirms
- Show remaining worktrees after removal

### Switch to Worktree
Show the command to navigate to the worktree:
```bash
cd ../accounting-webapp-worktrees/<name>
```

## Common Workflows:

### Parallel Agent Development
```bash
# Create worktrees for multiple agents
/utils:worktree create agent-1
/utils:worktree create agent-2
/utils:worktree create agent-3

# Each agent works independently
cd ../accounting-webapp-worktrees/agent-1
# Agent 1 works here
```

### Issue-based Development
```bash
# Create worktree for specific issue
/utils:worktree create issue-525

# Work on the issue
cd ../accounting-webapp-worktrees/issue-525
# Implement fix, commit, push, create PR

# Clean up after merge
/utils:worktree remove issue-525
```

### Quick Context Switch
```bash
# Currently working on feature
# Need to fix urgent bug

/utils:worktree create hotfix-urgent
cd ../accounting-webapp-worktrees/hotfix-urgent
# Fix bug, commit, push

# Return to feature work
cd /Users/adji/Documents/VISI-DIGITAL-TERPADU/accounting-webapp
# Continue feature work undisturbed
```

## Behavior:

- **Always validate** git repository exists
- **Check worktree existence** before create/remove
- **Provide clear feedback** after each action
- **Show helpful commands** after creation
- **Handle errors gracefully** with user-friendly messages
- **Confirm destructive actions** (remove with uncommitted changes)
- **Use consistent naming** conventions (issue-123, agent-1, etc.)

## Worktree Naming Conventions:

- **Issues**: `issue-123`, `issue-456-description`
- **Agents**: `agent-1`, `agent-2`, `agent-claude`
- **Features**: `feature-name`, `feature-tax-calc`
- **Hotfixes**: `hotfix-auth`, `hotfix-security`
- **Experiments**: `experiment-redis`, `test-approach`

## Error Handling:

- **Not a git repo**: Show error and suggest checking directory
- **Worktree exists**: List existing worktrees and suggest different name
- **Branch doesn't exist**: Show available branches
- **Uncommitted changes**: Warn and ask for confirmation before remove
- **Invalid arguments**: Show usage examples

## Output Format:

### After Create:
```
✅ Worktree created successfully!

📁 Path: ../accounting-webapp-worktrees/agent-1
🌿 Branch: agent-work-1
📍 Based on: main

To start working:
  cd ../accounting-webapp-worktrees/agent-1

To remove when done:
  /utils:worktree remove agent-1
```

### After Remove:
```
✅ Worktree 'agent-1' removed successfully!

Remaining worktrees: 5
  • main
  • agent-2
  • agent-3
  • issue-420
  • issue-421 (current)
```

### List Output:
```
📋 Git Worktrees (6 total)

1. accounting-webapp (current)
   Branch: perf/issue-421-n1-query-index-validation-a7b9c3f2
   Status: Clean

2. agent-1
   Path: ../accounting-webapp-worktrees/agent-1
   Branch: agent-work-1
   Status: Clean

3. agent-2
   Path: ../accounting-webapp-worktrees/agent-2
   Branch: agent-work-2
   Status: 2 uncommitted changes

[... more worktrees ...]

💡 Create: /utils:worktree create <name>
💡 Remove: /utils:worktree remove <name>
💡 Switch: cd ../accounting-webapp-worktrees/<name>
```

## Integration with Other Commands:

Works well with:
- `/dev:solve-issue <number>` - Create worktree, solve issue, create PR
- `/dev:create-fix-pr` - Work in separate worktree for fix
- `/dev:review-pr` - Review PR in separate worktree
- `/utils:todo` - Track worktree-specific tasks

## Best Practices:

1. **One worktree per issue** for clean separation
2. **Remove after merge** to keep workspace clean
3. **Use descriptive names** for easy identification
4. **Check status before remove** to avoid losing work
5. **Keep main worktree** as reference copy

Always be concise, helpful, and provide clear next steps after each action.
