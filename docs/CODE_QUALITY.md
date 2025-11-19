# Code Quality & Linting Guide

This document outlines the code quality standards, linting configuration, and development workflow for the Cooperative ERP Lite project.

## Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Linting Tools](#linting-tools)
- [Running Linters](#running-linters)
- [Pre-commit Hooks](#pre-commit-hooks)
- [CI/CD Integration](#cicd-integration)
- [IDE Setup](#ide-setup)
- [Linting Rules](#linting-rules)
- [Troubleshooting](#troubleshooting)

## Overview

We use a comprehensive linting setup to ensure code quality, consistency, and security across both backend (Go) and frontend (TypeScript/React) codebases.

**Quality Goals:**
- ✅ Consistent code style across the team
- ✅ Early detection of bugs and security issues
- ✅ Automated code formatting
- ✅ Type safety enforcement
- ✅ Best practices compliance

## Quick Start

### Initial Setup

```bash
# Install Go linter (if not already installed)
brew install golangci-lint

# Install frontend dependencies (includes ESLint & Prettier)
cd frontend
npm install

# Install pre-commit hooks (optional but recommended)
pip install pre-commit
pre-commit install
```

### Run All Quality Checks

```bash
# From project root
./scripts/lint.sh

# Auto-fix issues where possible
./scripts/lint.sh --fix

# Run only backend checks
./scripts/lint.sh --backend

# Run only frontend checks
./scripts/lint.sh --frontend
```

## Linting Tools

### Backend (Go)

| Tool | Purpose | Configuration |
|------|---------|---------------|
| **golangci-lint** | Meta-linter aggregating 50+ linters | `.golangci.yml` |
| **gofmt** | Code formatting | Built-in |
| **gofumpt** | Stricter gofmt | Via golangci-lint |
| **go vet** | Static analysis | Built-in |
| **gosec** | Security scanner | Via golangci-lint |

**Enabled Linters (40+):**
- errcheck, staticcheck, unused, gosimple (core quality)
- gosec (security)
- gocyclo, gocognit (complexity)
- revive, stylecheck (style)
- And 30+ more specialized linters

### Frontend (TypeScript/React/Next.js)

| Tool | Purpose | Configuration |
|------|---------|---------------|
| **ESLint** | Code quality & patterns | `frontend/eslint.config.mjs` |
| **Prettier** | Code formatting | `frontend/.prettierrc` |
| **TypeScript** | Type checking | `frontend/tsconfig.json` |

**ESLint Plugins:**
- `@typescript-eslint` - TypeScript-specific rules
- `eslint-config-next` - Next.js best practices
- `eslint-plugin-jsx-a11y` - Accessibility checks
- `eslint-plugin-react-hooks` - React hooks rules

## Running Linters

### Backend

```bash
cd backend

# Run all linters
golangci-lint run

# Auto-fix issues
golangci-lint run --fix

# Run specific linter
golangci-lint run --disable-all -E errcheck

# Run with increased timeout
golangci-lint run --timeout=10m

# Format code
gofmt -w .

# Run go vet
go vet ./...

# Check go.mod tidiness
go mod tidy
```

### Frontend

```bash
cd frontend

# Run ESLint
npm run lint

# Auto-fix ESLint issues
npm run lint:fix

# Check code formatting
npm run format:check

# Auto-format code
npm run format

# Type check
npm run type-check

# Run all quality checks
npm run quality

# Auto-fix all issues
npm run quality:fix
```

## Pre-commit Hooks

Pre-commit hooks automatically run linters before each commit, preventing bad code from entering the repository.

### Installation

```bash
# Install pre-commit framework
pip install pre-commit

# Install hooks
pre-commit install

# Run manually on all files
pre-commit run --all-files
```

### What Gets Checked

**On Every Commit:**
- ✅ Trailing whitespace removal
- ✅ End-of-file fixes
- ✅ YAML/JSON validation
- ✅ Large file detection
- ✅ Merge conflict markers
- ✅ Private key detection
- ✅ Go: formatting, linting, imports
- ✅ Frontend: ESLint, Prettier
- ✅ Markdown linting
- ✅ Spell checking

**Configuration:** `.pre-commit-config.yaml`

### Bypassing Hooks

```bash
# Skip hooks for emergency commits (not recommended)
git commit --no-verify -m "Emergency fix"
```

## CI/CD Integration

GitHub Actions automatically runs all quality checks on:
- Every push to `main` or `develop`
- Every pull request

### Workflow Jobs

**1. Backend Linting** (`.github/workflows/lint.yml`)
- golangci-lint with full ruleset
- go vet static analysis
- go.mod tidiness check
- Test execution with coverage

**2. Frontend Linting**
- ESLint checks
- Prettier formatting validation
- TypeScript type checking
- Test execution with coverage

**3. Security Scanning**
- Trivy vulnerability scanner
- gosec security analysis
- Dependency review (PRs only)

**4. Quality Summary**
- Aggregates all check results
- Fails PR if any check fails

### Viewing Results

1. Go to the "Actions" tab in GitHub
2. Click on the latest workflow run
3. View detailed logs for each job
4. See annotations on PR files

## IDE Setup

### Visual Studio Code (Recommended)

The project includes VS Code configuration for automatic linting and formatting.

**1. Install Recommended Extensions**

When you open the project, VS Code will prompt you to install recommended extensions:
- Prettier - Code formatter
- ESLint
- Go (golang.go)
- GitLens
- Error Lens

**2. Configuration Applied**

Settings in `.vscode/settings.json`:
- ✅ Format on save enabled
- ✅ Auto-fix ESLint issues on save
- ✅ Organize imports on save
- ✅ Go: Use golangci-lint for linting
- ✅ Go: Use gofumpt for formatting
- ✅ TypeScript: Prefer absolute imports

**3. Debugging**

Pre-configured launch configurations:
- Debug Backend (Go)
- Debug Backend Tests
- Debug Frontend (Next.js)
- Debug Frontend Tests
- Full Stack Debug (both together)

**4. Tasks**

Run from Command Palette (Cmd+Shift+P → "Tasks: Run Task"):
- `go: lint` - Run golangci-lint
- `frontend: lint` - Run ESLint
- `quality: lint all` - Run all linters
- `quality: fix all` - Auto-fix all issues

### Other IDEs

**GoLand / WebStorm:**
1. Enable golangci-lint: Settings → Tools → golangci-lint
2. Enable ESLint: Settings → Languages → JavaScript → ESLint
3. Enable Prettier: Settings → Languages → JavaScript → Prettier

**Vim/Neovim:**
- Use `vim-go` plugin with `:GoMetaLinter`
- Use `coc-eslint` for ESLint
- Use `ale` or `null-ls` for general linting

## Linting Rules

### Go Rules Highlights

**Complexity Limits:**
- Cyclomatic complexity: max 15
- Cognitive complexity: max 20
- Function length: 100 lines / 50 statements
- Nesting depth: max 5

**Security Rules:**
- No `fmt.Print*` (use logger instead)
- No `panic` (use proper error handling)
- SQL injection detection
- Hardcoded credentials detection
- Crypto misuse detection

**Code Quality:**
- Unused code detection
- Error handling enforcement
- Proper context.Context usage
- Interface design validation
- Struct tag validation (GORM, JSON)

### TypeScript/React Rules Highlights

**TypeScript:**
- No `any` type (warn)
- Unused variables (with `_` prefix exception)
- Consistent type imports
- Prefer nullish coalescing (`??`)
- Prefer optional chaining (`?.`)

**React:**
- Hooks rules enforcement
- Exhaustive dependencies in useEffect
- Self-closing components
- No unnecessary JSX curly braces

**Next.js:**
- Proper Link component usage
- Image component for optimization
- No `<a>` tags for internal links

**Security:**
- No `eval()` or `new Function()`
- No inline scripts
- Target="_blank" security

**Accessibility:**
- Alt text for images
- ARIA attributes validation
- Role requirements
- Keyboard accessibility

### Exceptions and Overrides

**Go Test Files:**
Many linters are disabled for `_test.go` files:
- gocyclo, funlen (test functions can be longer)
- dupl (test code may have duplication)
- gosec (false positives in tests)

**Indonesian Terms:**
Spell checker ignores:
- koperasi, simpanan, pokok, wajib, sukarela
- anggota, keuangan

**Disable Specific Rules:**

```go
// In Go (use sparingly)
//nolint:errcheck // Justification required
result, _ := someFunction()
```

```typescript
// In TypeScript
// eslint-disable-next-line @typescript-eslint/no-explicit-any
const data: any = unknownData;
```

## Troubleshooting

### Common Issues

**1. "golangci-lint: command not found"**

```bash
# Install golangci-lint
brew install golangci-lint

# Or download binary
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

**2. "ESLint errors but can't auto-fix"**

```bash
cd frontend

# Some errors can't be auto-fixed, read the error message
npm run lint

# Try fixing what's possible
npm run lint:fix

# Manually fix remaining issues
```

**3. "Pre-commit hooks are slow"**

```bash
# Run only on changed files (default behavior)
git commit -m "message"

# To run on all files (slower)
pre-commit run --all-files
```

**4. "Linter false positives"**

```bash
# Disable specific linter for a line
//nolint:lintername // Explanation why

# Update .golangci.yml to exclude pattern
# See exclude-rules section
```

**5. "Go imports not organizing on save"**

In VS Code:
1. Check Go extension is installed
2. Verify settings.json has `"source.organizeImports": "explicit"`
3. Restart VS Code

### Performance Tips

**1. Speed up golangci-lint:**

```bash
# Use --fast mode (skips some slow linters)
golangci-lint run --fast

# Run only on changed files
golangci-lint run --new-from-rev=HEAD~1

# Increase parallel processing
golangci-lint run --concurrency=8
```

**2. Speed up ESLint:**

```bash
# Use cache
npm run lint -- --cache

# Run only on specific files
npx eslint app/components/*.tsx
```

## Best Practices

1. **Run linters before committing:**
   ```bash
   ./scripts/lint.sh --fix
   ```

2. **Never commit with `--no-verify` unless emergency**

3. **Fix linting issues immediately - don't accumulate debt**

4. **Use IDE auto-fix features while coding**

5. **Write code that passes linters from the start**

6. **If you need to disable a rule, document why**

7. **Keep configurations in sync across team**

8. **Update linters regularly but test thoroughly**

## Resources

- [golangci-lint documentation](https://golangci-lint.run/)
- [ESLint rules](https://eslint.org/docs/rules/)
- [Prettier options](https://prettier.io/docs/en/options.html)
- [TypeScript ESLint](https://typescript-eslint.io/)
- [Next.js ESLint](https://nextjs.org/docs/basic-features/eslint)

## Support

**Issues with linting setup?**
1. Check this documentation
2. Review error messages carefully
3. Ask in team chat
4. Create GitHub issue with `linting` label
