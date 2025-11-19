# Linting Quick Reference Card

## 🚀 Quick Commands

```bash
# Run all quality checks
./scripts/lint.sh

# Auto-fix all issues
./scripts/lint.sh --fix

# Backend only
./scripts/lint.sh --backend

# Frontend only
./scripts/lint.sh --frontend
```

## 🔧 Backend (Go)

```bash
cd backend

# Lint
golangci-lint run

# Format
gofmt -w .

# Vet
go vet ./...

# Tidy modules
go mod tidy

# Auto-fix
golangci-lint run --fix
```

## 🎨 Frontend (TypeScript)

```bash
cd frontend

# Lint
npm run lint

# Format check
npm run format:check

# Type check
npm run type-check

# All checks
npm run quality

# Auto-fix everything
npm run quality:fix
```

## 🪝 Pre-commit Hooks

```bash
# Install
pip install pre-commit
pre-commit install

# Run manually
pre-commit run --all-files

# Skip (emergency only)
git commit --no-verify
```

## 📝 Disable Rules

```go
// Go - with explanation required
//nolint:errcheck // Explanation
code()
```

```typescript
// TypeScript
// eslint-disable-next-line rule-name
code()
```

## 🔍 IDE Quick Fixes

**VS Code:**
- Format: `Shift+Alt+F` (Mac: `Shift+Option+F`)
- Fix ESLint: `Cmd+Shift+P` → "ESLint: Fix all auto-fixable Problems"
- Organize Imports: Automatic on save

## ✅ CI/CD Status

Check `.github/workflows/lint.yml` for status on PRs

## 📚 Full Documentation

See `/docs/CODE_QUALITY.md` for complete guide
