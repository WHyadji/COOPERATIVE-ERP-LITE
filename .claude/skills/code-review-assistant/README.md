# Code Review Assistant Skill

A comprehensive code review skill that performs systematic analysis of code for security vulnerabilities, performance issues, code quality, and best practices violations.

## Features

### 🔒 Security Analysis
- SQL injection detection
- XSS vulnerability scanning
- Hardcoded secrets detection
- Weak cryptography identification
- Authentication/authorization issues
- Command injection risks
- Path traversal vulnerabilities
- Insecure dependencies

### ⚡ Performance Optimization
- Algorithm complexity analysis (O(n²), O(n³))
- N+1 query detection
- Memory leak identification
- Inefficient loops
- Database query optimization
- Caching opportunities
- Bundle size issues
- Async/await optimization

### 🏗️ Code Quality
- Code duplication (DRY violations)
- Cyclomatic complexity measurement
- Method/function length analysis
- Dead code detection
- Magic numbers/strings
- Naming convention violations
- Error handling gaps
- Test coverage analysis

### 📋 Best Practices
- SOLID principles compliance
- Design pattern recommendations
- API design issues
- Documentation completeness
- Type safety checks
- Dependency management
- Configuration issues
- Logging standards

## Quick Start

### Review a Single File
```bash
# Review with all checks
python scripts/review_code.py --file app.py

# Review with specific checks
python scripts/review_code.py --file index.js --checks security,performance

# Output to HTML
python scripts/review_code.py --file main.go --format html --output review.html
```

### Review Entire Project
```bash
# Review all files in project
python scripts/review_code.py --project ./src

# Exclude certain paths
python scripts/review_code.py --project . --exclude node_modules,dist,build

# Generate markdown report
python scripts/review_code.py --project . --format markdown --output REVIEW.md
```

### Review Git Changes
```bash
# Review changes in last commit
python scripts/review_code.py --git-diff HEAD~1

# Review changes between branches
python scripts/review_code.py --git-diff main..feature-branch

# Review staged changes
python scripts/review_code.py --git-diff --cached
```

## Supported Languages

- **JavaScript/TypeScript** - Node.js, React, Vue, Angular
- **Python** - Django, Flask, FastAPI
- **Java** - Spring, Spring Boot
- **C#** - .NET Core, ASP.NET
- **Go** - Standard library, popular frameworks
- **Ruby** - Rails, Sinatra
- **PHP** - Laravel, Symfony
- **Rust** - Memory safety, ownership
- **SQL** - PostgreSQL, MySQL, SQLite
- **Shell** - Bash, PowerShell

## Report Formats

### Text Report
```
CODE REVIEW REPORT
=====================================
Generated: 2024-01-10 14:30:00
Files reviewed: 23
Lines of code: 4,567

SUMMARY
-------
Total issues: 47
  CRITICAL: 3
  HIGH: 8
  MEDIUM: 15
  LOW: 21

CRITICAL ISSUES
---------------
1. SQL Injection vulnerability
   File: database.py:45
   Code: query = f"SELECT * FROM users WHERE id = {user_id}"
   Recommendation: Use parameterized queries
```

### HTML Report
Beautiful, interactive HTML report with:
- Color-coded severity levels
- Code snippets with syntax highlighting
- Expandable issue details
- Summary statistics
- Filter and search capabilities

### Markdown Report
GitHub/GitLab friendly markdown with:
- Issue grouping by severity
- Code blocks with language hints
- Emoji indicators for severity
- Links to file locations
- Actionable recommendations

### JSON Report
Structured data for CI/CD integration:
```json
{
  "metadata": {
    "generated": "2024-01-10T14:30:00",
    "files_reviewed": 23,
    "lines_of_code": 4567
  },
  "summary": {
    "critical": 3,
    "high": 8,
    "medium": 15,
    "low": 21
  },
  "issues": [...]
}
```

## Configuration

### Using Configuration File
Create `.codereview.yml`:
```yaml
enabled_checks:
  - security
  - performance
  - quality
  - best_practices

severity_threshold: low

max_file_lines: 500
max_method_lines: 50
max_complexity: 10

ignore_paths:
  - node_modules
  - vendor
  - .git
  - dist
  - build

file_extensions:
  - .py
  - .js
  - .ts
  - .java
  - .cs
  - .go

custom_rules:
  - pattern: "console.log"
    message: "Remove console.log statements"
    severity: low
  - pattern: "debugger"
    message: "Remove debugger statements"
    severity: high
```

Use with:
```bash
python scripts/review_code.py --project . --config .codereview.yml
```

## CI/CD Integration

### GitHub Actions
```yaml
name: Code Review
on: [pull_request]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Python
        uses: actions/setup-python@v2
        with:
          python-version: '3.9'
      
      - name: Run Code Review
        run: |
          python scripts/review_code.py \
            --git-diff origin/main..HEAD \
            --format markdown \
            --output review.md
      
      - name: Comment PR
        if: always()
        uses: actions/github-script@v6
        with:
          script: |
            const fs = require('fs');
            const review = fs.readFileSync('review.md', 'utf8');
            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: review
            });
```

### Pre-commit Hook
```bash
#!/bin/sh
# .git/hooks/pre-commit

python scripts/review_code.py --git-diff --cached --checks security,critical

if [ $? -ne 0 ]; then
    echo "Critical issues found. Fix before committing."
    exit 1
fi
```

## Issue Severity Levels

### 🔴 Critical
- SQL injection vulnerabilities
- Hardcoded passwords/API keys
- Command injection risks
- Unencrypted sensitive data
- Missing authentication

### 🟠 High
- XSS vulnerabilities
- Weak cryptography
- N+1 query problems
- Memory leaks
- Missing authorization

### 🟡 Medium
- Large methods/files
- High cyclomatic complexity
- Missing error handling
- Performance bottlenecks
- Code duplication

### 🟢 Low
- Magic numbers
- Naming convention violations
- Missing type hints
- Unused imports
- Console.log statements

### 🔵 Info
- TODO/FIXME comments
- Documentation gaps
- Style inconsistencies
- Optimization opportunities
- Best practice suggestions

## Command Options

```bash
python scripts/review_code.py [OPTIONS]

Input Options (one required):
  --file FILE          Review single file
  --project PATH       Review entire project
  --git-diff DIFF      Review git changes

Check Options:
  --checks CHECKS      Comma-separated: security,performance,quality,best_practices,all
                      Default: all

Output Options:
  --output FILE        Save report to file
  --format FORMAT      Output format: text,json,html,markdown
                      Default: text

Configuration:
  --config FILE        Configuration file path
  --exclude PATHS      Comma-separated paths to exclude
```

## Examples

### Security-Only Review
```bash
python scripts/review_code.py --project . --checks security --format html --output security-report.html
```

### Performance Analysis
```bash
python scripts/review_code.py --file api.py --checks performance --format markdown
```

### Pre-Release Check
```bash
python scripts/review_code.py --project ./src --checks security,critical --format json --output release-check.json
```

### Pull Request Review
```bash
python scripts/review_code.py --git-diff main..$(git branch --show-current) --format markdown
```

## Best Practices

1. **Run on every PR** - Integrate with CI/CD
2. **Fix critical issues immediately** - Block merge if critical
3. **Track metrics over time** - Monitor code quality trends
4. **Customize rules** - Add project-specific checks
5. **Regular full scans** - Weekly/monthly comprehensive reviews
6. **Team training** - Share reports in code reviews
7. **Incremental improvement** - Fix high-priority issues first

## Exit Codes

- `0` - Success, no critical issues
- `1` - Critical issues found
- `2` - Error during execution

## Requirements

- Python 3.7+
- Git (for git-diff functionality)
- No additional dependencies for basic functionality

## License

This skill is provided for use with Claude AI to enhance code review capabilities.