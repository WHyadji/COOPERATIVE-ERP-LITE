# Contributing to Cooperative ERP Lite

Thank you for your interest in contributing to Cooperative ERP Lite! This document provides guidelines and instructions for contributing to the project.

## 🌟 Ways to Contribute

- **Code**: Implement new features, fix bugs, improve performance
- **Documentation**: Improve docs, add examples, fix typos
- **Testing**: Write tests, report bugs, improve test coverage
- **Design**: UI/UX improvements, wireframes, mockups
- **Translation**: Help translate to Indonesian and regional languages
- **Feedback**: User testing, feature suggestions, bug reports

## 🚀 Getting Started

### 1. Fork and Clone

```bash
# Fork the repository on GitHub
# Then clone your fork
git clone https://github.com/YOUR_USERNAME/COOPERATIVE-ERP-LITE.git
cd COOPERATIVE-ERP-LITE
```

### 2. Setup Development Environment

Follow the [Quick Start Guide](docs/quick-start-guide.md) to setup your development environment.

### 3. Create a Branch

```bash
# Create a feature branch
git checkout -b feature/your-feature-name

# Or a bugfix branch
git checkout -b bugfix/issue-description
```

## 📋 Development Workflow

### Branch Naming Convention

- `feature/` - New features (e.g., `feature/member-import`)
- `bugfix/` - Bug fixes (e.g., `bugfix/login-validation`)
- `hotfix/` - Critical production fixes
- `docs/` - Documentation updates
- `refactor/` - Code refactoring
- `test/` - Test additions/improvements

### Commit Message Format

We follow the [Conventional Commits](https://www.conventionalcommits.org/) specification:

```
<type>(<scope>): <subject>

<body>

<footer>
```

**Types:**
- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Code style changes (formatting, etc.)
- `refactor`: Code refactoring
- `test`: Adding or updating tests
- `chore`: Maintenance tasks

**Examples:**
```bash
feat(auth): add JWT refresh token mechanism

Implement refresh token to improve security and user experience.
Users can now stay logged in longer without re-authentication.

Closes #123

---

fix(pos): correct total calculation for multiple items

Fixed bug where discount was applied incorrectly when cart
had more than 5 items.

Fixes #456

---

docs(readme): update installation instructions

Added Docker setup instructions and troubleshooting section.
```

## 💻 Code Style Guide

### Go (Backend)

- Follow [Effective Go](https://go.dev/doc/effective_go)
- Use `gofmt` for formatting
- Use `golangci-lint` for linting
- Write descriptive variable names in **Bahasa Indonesia**
- Add comments for exported functions
- Keep functions small and focused

**Example:**
```go
// BuatAnggota membuat anggota baru dengan validasi
func (s *AnggotaService) BuatAnggota(idKoperasi uuid.UUID, req *BuatAnggotaRequest) (*models.AnggotaResponse, error) {
    // Generate nomor anggota otomatis
    nomorAnggota, err := s.GenerateNomorAnggota(idKoperasi)
    if err != nil {
        return nil, err
    }

    // ... rest of implementation
}
```

### TypeScript/React (Frontend)

- Follow [Airbnb JavaScript Style Guide](https://github.com/airbnb/javascript)
- Use TypeScript for type safety
- Use functional components with hooks
- Follow React best practices
- Use ESLint and Prettier

**Example:**
```typescript
// useAnggota.ts - Custom hook untuk manajemen anggota
export function useAnggota() {
  const [anggotaList, setAnggotaList] = useState<Anggota[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchAnggota = async () => {
    setLoading(true);
    try {
      const response = await api.get('/api/members');
      setAnggotaList(response.data);
    } catch (error) {
      console.error('Gagal mengambil data anggota:', error);
    } finally {
      setLoading(false);
    }
  };

  return { anggotaList, loading, fetchAnggota };
}
```

### Naming Conventions

**Backend (Go):**
- Variables: `camelCase` in Bahasa Indonesia (e.g., `namaAnggota`, `jumlahSetoran`)
- Functions: `PascalCase` with descriptive names (e.g., `BuatAnggota`, `HitungSaldo`)
- Constants: `UPPER_SNAKE_CASE` (e.g., `STATUS_AKTIF`, `PERAN_ADMIN`)

**Frontend (TypeScript):**
- Components: `PascalCase` (e.g., `MemberList`, `LoginForm`)
- Files: `kebab-case` (e.g., `member-list.tsx`, `login-form.tsx`)
- Functions: `camelCase` (e.g., `fetchMembers`, `handleSubmit`)

## 🧪 Testing

### Running Tests

**Backend:**
```bash
cd backend

# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./internal/services/...
```

**Frontend:**
```bash
cd frontend

# Run tests
npm test

# Run with coverage
npm test -- --coverage
```

### Writing Tests

- Write unit tests for all business logic
- Write integration tests for API endpoints
- Aim for >70% code coverage
- Test edge cases and error scenarios

## 📝 Pull Request Process

### 1. Before Submitting

- [ ] Code follows style guidelines
- [ ] All tests pass
- [ ] New features have tests
- [ ] Documentation is updated
- [ ] Commit messages follow convention
- [ ] Branch is up to date with `develop`

### 2. Submit Pull Request

```bash
# Push your branch
git push origin feature/your-feature-name

# Create PR on GitHub
```

**PR Title Format:**
```
[Type] Brief description

Examples:
[Feature] Add member import from Excel
[Fix] Correct balance calculation in reports
[Docs] Update API documentation
```

**PR Description Template:**
```markdown
## Description
Brief description of changes

## Related Issue
Closes #123

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## How Has This Been Tested?
- [ ] Unit tests
- [ ] Integration tests
- [ ] Manual testing

## Screenshots (if applicable)

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-reviewed code
- [ ] Commented hard-to-understand areas
- [ ] Updated documentation
- [ ] No new warnings
- [ ] Added tests
- [ ] All tests pass
```

### 3. Code Review

- Address reviewer feedback promptly
- Make requested changes in new commits
- Don't force push after review started
- Be open to suggestions and discussion

### 4. Merge

- PRs are merged by maintainers after approval
- Squash and merge for clean history
- Delete branch after merge

## 🐛 Bug Reports

### Before Reporting

1. Check if bug already reported in [Issues](https://github.com/YOUR_USERNAME/COOPERATIVE-ERP-LITE/issues)
2. Try to reproduce in latest version
3. Collect relevant information

### Bug Report Template

Use the bug report template when creating an issue. Include:

- Clear description
- Steps to reproduce
- Expected vs actual behavior
- Screenshots/logs
- Environment details (OS, browser, version)

## 💡 Feature Requests

### Before Requesting

1. Check if feature already requested
2. Ensure it aligns with project goals
3. Consider MVP scope constraints

### Feature Request Template

Use the feature request template when creating an issue. Include:

- Problem description
- Proposed solution
- Alternatives considered
- Additional context

## 📖 Documentation

### Types of Documentation

- **Code comments**: For complex logic
- **README**: Project overview and quick start
- **API docs**: Endpoint specifications
- **User guides**: How-to for end users
- **Architecture docs**: System design

### Documentation Style

- Clear and concise
- Use examples
- Keep up to date with code
- Bilingual (English + Indonesian) preferred

## 🎯 Priority Labels

Issues and PRs are labeled by priority:

- `P0-critical`: Must fix immediately
- `P1-high`: Fix in current sprint
- `P2-medium`: Fix in next sprint
- `P3-low`: Fix when possible

## 🏷️ Other Labels

- `good-first-issue`: Good for newcomers
- `help-wanted`: Need community help
- `enhancement`: New feature
- `bug`: Something broken
- `documentation`: Documentation improvements
- `wontfix`: Will not be fixed

## 👥 Community

### Code of Conduct

- Be respectful and inclusive
- Welcome newcomers
- Provide constructive feedback
- Focus on what's best for the community
- Show empathy

### Communication Channels

- **GitHub Issues**: Bug reports, feature requests
- **GitHub Discussions**: Q&A, ideas, general discussion
- **Email**: [support@cooperative-erp.com](mailto:support@cooperative-erp.com) (example)

## 📜 License

By contributing, you agree that your contributions will be licensed under the MIT License.

## 🙏 Recognition

Contributors will be recognized in:
- README contributors section
- Release notes
- Project website (when available)

## ❓ Questions?

If you have questions about contributing, please:

1. Check the documentation
2. Search existing issues
3. Ask in GitHub Discussions
4. Contact maintainers

---

**Thank you for contributing to Cooperative ERP Lite! Together we're building the future of Indonesian cooperatives. 🚀**
