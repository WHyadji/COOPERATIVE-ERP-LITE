# Custom Claude Slash Commands

This directory contains custom slash commands organized by category to streamline development workflow for the accounting web application.

## Command Organization

Commands are organized in subdirectories by functional area:

- **admin/** - Admin panel features and management
- **analytics/** - Analytics and tracking implementation
- **code-quality/** - Code quality analysis and improvement
- **debug/** - Debugging and auto-fixing tools
- **dev/** - Development workflow commands
- **devops/** - DevOps, containerization, and infrastructure
- **docs/** - Documentation generation and management
- **implement/** - TDD implementation workflows
- **plan/** - Planning and architecture commands
- **refactor/** - Code refactoring and optimization
- **security/** - Security audits and fixes
- **session/** - Session management commands
- **supabase/** - Database and Supabase management
- **testing/** - Testing utilities and test generation
- **utils/** - General utility commands

## Available Commands

### Development Commands (`dev/`)

#### `/dev:setup`
Sets up the complete development environment including dependencies, environment variables, and services.
```bash
/dev:setup
```

#### `/dev:component`
Generates a new React component with TypeScript, styling, and best practices.
```bash
/dev:component Button with loading state
/dev:component UserProfile card with avatar
```

#### `/dev:api`
Creates a new API endpoint with validation, authentication, and error handling.
```bash
/dev:api user profile endpoint
/dev:api journal entry creation
```

#### `/dev:debug-issue`
Systematically debugs and fixes issues in the codebase.
```bash
/dev:debug-issue
```

#### `/dev:commit`
Generates descriptive git commit messages by analyzing staged changes.
```bash
/dev:commit
```

#### `/dev:review-pr [pr-number]`
Performs comprehensive PR review with multi-agent analysis (Code Quality, Security, Performance). Automatically posts line-by-line comments, creates a detailed summary, and generates sub-issues for all findings via GitHub CLI.
```bash
/dev:review-pr 123                    # Review by PR number
/dev:review-pr --current               # Review current branch's PR
/dev:review-pr https://github.com/... # Review by URL
/dev:review-pr 123 --no-issues        # Review without creating issues
```

**Features:**
- 🤖 Multi-agent parallel analysis (Code Quality Guardian, Security Specialist, Performance Engineer)
- 💬 Automated inline comments on code
- 📊 Comprehensive summary with metrics and scores
- 🎯 SOLID principles compliance check
- 🔒 Security vulnerability assessment
- ⚡ Performance impact analysis
- 🇮🇩 Indonesian compliance validation (SAK EP, tax calculations)
- 📋 Auto-creates GitHub issues for all findings
- ✅ Posts review decision (approve/request changes)
- 🚫 Never auto-merges (you control merge)

#### `/dev:verify-pr [pr-number]`
Comprehensive pre-merge verification: runs all tests, code quality checks, builds, and security scans. Posts detailed verification report to PR.
```bash
/dev:verify-pr 123           # Verify PR #123
/dev:verify-pr --current     # Verify current branch's PR
/dev:verify-pr 123 --fix-issues  # Auto-fix linting/formatting issues
```

#### `/dev:create-fix-pr [issue-numbers...]`
Creates fix branch and PR that properly references and closes GitHub issues. Automates issue-to-PR workflow with smart branch naming and comprehensive PR description.
```bash
/dev:create-fix-pr 45 67 89          # Create PR to fix multiple issues
/dev:create-fix-pr --security 123    # Create security fix PR
/dev:create-fix-pr --from-review 456 # Create PR from review findings
```

#### `/dev:merge-pr [pr-number]`
Safely merges PR with final validation, proper merge strategy, automatic branch cleanup, and issue closing verification.
```bash
/dev:merge-pr 123            # Standard merge
/dev:merge-pr 123 --squash   # Squash merge
/dev:merge-pr 123 --rebase   # Rebase merge
/dev:merge-pr --current      # Merge current branch's PR
```

#### `/dev:identify-issues [path]`
Proactive codebase scanning to identify issues before creating PRs: code quality, security, performance, testing gaps, and SOLID violations.
```bash
/dev:identify-issues services/accounting-ledger/  # Scan specific service
/dev:identify-issues --all                        # Scan entire codebase
/dev:identify-issues --security-only              # Security-focused scan
/dev:identify-issues --create-issues              # Auto-create GitHub issues
```

#### `/dev:link-issues [pr-number] [issue-numbers...]`
Links GitHub issues to existing PR for automatic closing on merge.
```bash
/dev:link-issues 123 45 67 89  # Link PR #123 to issues #45, #67, #89
```

#### `/dev:pr-status [pr-number]`
Quick PR health check dashboard: CI status, review status, conflicts, linked issues, and merge readiness.
```bash
/dev:pr-status 123       # Check PR #123 status
/dev:pr-status --current # Check current branch's PR
```

### Debug Commands (`debug/`)

#### `/debug:autofix`
Automatically identifies and fixes common code issues.
```bash
/debug:autofix
```

### DevOps Commands (`devops/`)

#### `/devops:containerize-application [application-type]`
Containerizes application with optimized Docker configuration and security.
```bash
/devops:containerize-application --node
/devops:containerize-application --multi-stage
```

#### `/devops:setup-linting [language]`
Configures comprehensive code linting and quality analysis tools.
```bash
/devops:setup-linting --typescript
/devops:setup-linting --multi-language
```

#### `/devops:design-database-schema [schema-type]`
Designs optimized database schemas with proper relationships and constraints.
```bash
/devops:design-database-schema --relational
/devops:design-database-schema fixed-assets --normalize
```

### Documentation Commands (`docs/`)

#### `/docs:create-prd [feature-name]`
Creates Product Requirements Document for new features.
```bash
/docs:create-prd tax-calculation-service
/docs:create-prd --template
/docs:create-prd --interactive
```

#### `/docs:add-changelog [version]`
Generates and maintains project changelog with Keep a Changelog format.
```bash
/docs:add-changelog 1.2.0
/docs:add-changelog patch "Fix transaction leak"
```

#### `/docs:update-docs [doc-type]`
Systematically updates project documentation.
```bash
/docs:update-docs --api
/docs:update-docs --architecture
/docs:update-docs --sync
```

#### `/docs:create-architecture-documentation [framework]`
Generates comprehensive architecture documentation with diagrams and ADRs.
```bash
/docs:create-architecture-documentation --c4-model
/docs:create-architecture-documentation --full-suite
```

### Supabase Commands (`supabase/`)

#### `/supabase:migrate`
Creates and manages database migrations.
```bash
/supabase:migrate create users table
/supabase:migrate add indexes to journal_entries
```

#### `/supabase:rls`
Creates Row Level Security policies for tables.
```bash
/supabase:rls accounts table
/supabase:rls journal_entries with company isolation
```

#### `/supabase:supabase-data-explorer [table-name]`
Explores and analyzes Supabase database data with intelligent querying.
```bash
/supabase:supabase-data-explorer accounts
/supabase:supabase-data-explorer --query "SELECT * FROM journal_entries WHERE company_id = '...'"
```

#### `/supabase:supabase-schema-sync [action]`
Synchronizes database schema with Supabase using MCP integration.
```bash
/supabase:supabase-schema-sync --pull
/supabase:supabase-schema-sync --diff
```

#### `/supabase:supabase-performance-optimizer [optimization-type]`
Optimizes Supabase database performance with intelligent analysis.
```bash
/supabase:supabase-performance-optimizer --queries
/supabase:supabase-performance-optimizer --indexes
/supabase:supabase-performance-optimizer --rls
```

### Planning Commands (`plan/`)

#### `/plan:solid-refactor <file_path>`
Analyzes code for SOLID principle violations and generates comprehensive TDD refactoring plan.
```bash
/plan:solid-refactor services/accounting-ledger/internal/service/account_service.go
/plan:solid-refactor frontend/components/Dashboard.tsx
```

#### `/plan:bugfix-plan <bug_description> <file_path>`
Generates comprehensive bug-fixing todo lists with intelligent code analysis and step-by-step resolution guidance.
```bash
/plan:bugfix-plan "transaction leak in database connections" services/accounting-ledger/internal/repository/dbstore.go
/plan:bugfix-plan "auth token expiry handling" frontend/lib/api/client.ts
```

### Implementation Commands (`implement/`)

#### `/implement:implement-tdd <file_path>`
Implements code changes following Test-Driven Development with automated test generation.
```bash
/implement:implement-tdd services/accounting-ledger/internal/service/tax_service.go --issue "PPH21 calculation"
/implement:implement-tdd frontend/components/JournalEntry.tsx --auto
```

#### `/implement:implement-feature`
Implements a new feature with proper structure and tests.
```bash
/implement:implement-feature
```

#### `/implement:implement-fix`
Implements bug fixes with proper testing and verification.
```bash
/implement:implement-fix
```

### Testing Commands (`testing/`)

#### `/testing:write-test`
Writes comprehensive tests for components or features.
```bash
/testing:write-test ContactForm component
/testing:write-test authentication flow
```

### Session Commands (`session/`)

#### `/session:session-start`
Starts a new development session with tracking.
```bash
/session:session-start
```

#### `/session:session-end`
Ends the current development session.
```bash
/session:session-end
```

#### `/session:session-list`
Lists all development sessions.
```bash
/session:session-list
```

#### `/session:session-current`
Shows the current active session.
```bash
/session:session-current
```

#### `/session:session-update`
Updates the current session details.
```bash
/session:session-update
```

#### `/session:session-help`
Shows help for session management commands.
```bash
/session:session-help
```

### Utility Commands (`utils/`)

#### `/utils:todo [action] [task-description]`
Manages project todos in todos.md file.
```bash
/utils:todo add "Implement PPH23 tax calculation"
/utils:todo complete "Fix auth middleware"
/utils:todo list
```

#### `/utils:ultra-think [problem or question]`
Deep analysis and problem solving with multi-dimensional thinking.
```bash
/utils:ultra-think "How should we architect the multi-tenant tax calculation system?"
/utils:ultra-think "What's the best approach for handling currency precision in accounting?"
```

#### `/utils:worktree [action] [name] [base-branch]`
Manages git worktrees for parallel agent development and isolated branch work.
```bash
/utils:worktree list                    # Show all worktrees
/utils:worktree create agent-4          # Create new worktree from main
/utils:worktree create issue-525 main   # Create worktree from specific branch
/utils:worktree remove agent-1          # Remove worktree safely
/utils:worktree switch agent-2          # Show command to switch to worktree
```

**Use cases:**
- Multiple Claude Code agents working simultaneously on different issues
- Parallel development without branch switching
- Quick context switching for hotfixes while preserving feature work
- Testing different implementation approaches in isolation

### Admin Commands (`admin/`)

#### `/admin:add-feature`
Adds new features to the admin panel with proper authentication and permissions.
```bash
/admin:add-feature user management dashboard
/admin:add-feature company settings configuration
```

### Analytics Commands (`analytics/`)

#### `/analytics:analytics`
Implements analytics tracking for user behavior and business metrics.
```bash
/analytics:analytics button clicks on dashboard
/analytics:analytics conversion funnel for journal entry creation
```

### Refactoring Commands (`refactor/`)

#### `/refactor:optimize-code`
Optimizes code for performance, bundle size, and best practices.
```bash
/refactor:optimize-code
```

#### `/refactor:refactor`
Systematic code refactoring workflow.
```bash
/refactor:refactor
```

#### `/refactor:fix-issue`
Fixes specific code issues with proper implementation.
```bash
/refactor:fix-issue
```

### Code Quality Commands (`code-quality/`)

#### `/code-quality:find-placeholder`
Comprehensive placeholder and incomplete code scanner.
```bash
/code-quality:find-placeholder services/accounting-ledger
/code-quality:find-placeholder frontend/components
```

#### `/code-quality:find-placeholder-api`
Comprehensive API placeholder and incomplete implementation scanner.
```bash
/code-quality:find-placeholder-api accounting-ledger
/code-quality:find-placeholder-api api-gateway
```

#### `/code-quality:find-mock-data`
Mock data and hardcoded response detector for APIs and services.
```bash
/code-quality:find-mock-data business-service
/code-quality:find-mock-data frontend/lib/api
```

#### `/code-quality:find-vulnerabilities`
Scans code for security vulnerabilities and issues.
```bash
/code-quality:find-vulnerabilities
```

#### `/code-quality:security_review`
Performs comprehensive security review of the codebase.
```bash
/code-quality:security_review
```

### Security Commands (`security/`)

#### `/security:security-audit`
Performs comprehensive security audit and fixes vulnerabilities.
```bash
/security:security-audit
```

## Command Features

All commands support:
- **Dynamic arguments** using `$ARGUMENTS` variable
- **Bash command execution** with `!` prefix for context gathering
- **File references** with `@` prefix for including file content
- **Context awareness** through automatic project state inspection
- **Tool restrictions** via `allowed-tools` in command frontmatter

## Command Structure

Each command file follows this structure:

```markdown
---
allowed-tools: ReadFile, WriteFile, Bash(npm:*), etc.
description: Brief description of what the command does
---

## Context
- Dynamic context using !`bash commands`
- File references with @path/to/file

## Your task
Clear instructions for Claude including $ARGUMENTS

Specific requirements and patterns to follow.
```

## Best Practices

1. **Use specific arguments** for better results
2. **Commands are context-aware** and analyze project state automatically
3. **Chain commands** for complex workflows when needed
4. **Follow project standards** - all commands respect CLAUDE.md conventions
5. **Provide clear descriptions** when arguments are required

## Creating New Commands

To create a new command:

1. **Choose the appropriate directory** based on functional area
2. **Create a `.md` file** with the command name (e.g., `new-feature.md`)
3. **Add YAML frontmatter** with allowed tools and description
4. **Include context gathering** using bash commands or file references
5. **Define clear task instructions** with specific requirements
6. **Update this README** with the new command documentation

Example:

```markdown
---
allowed-tools: Read, Write, Edit, Bash(go:*)
description: Creates a new Go service following clean architecture
---

## Context
Current services: !`ls -d services/*/`
Project structure: @CLAUDE.md

## Your task
Create a new Go microservice named "$ARGUMENTS" following the clean architecture pattern used in accounting-ledger.

Requirements:
- Follow repository pattern
- Include health check endpoints
- Add to docker-compose.yml
- Create basic tests
```

## Project-Specific Conventions

This accounting web application follows specific patterns:

- **Go services** use clean architecture with domain/service/repository layers
- **Frontend** uses Next.js 15 with TypeScript and Tailwind CSS
- **Testing** follows simplified approach with stubs instead of heavy mocking
- **Multi-tenancy** requires company-based data isolation
- **Indonesian compliance** with SAK EP standards and tax regulations
- **SOLID principles** enforced across all code

Refer to `/CLAUDE.md` for complete project guidelines and conventions.
