---
name: Code Review Assistant
description: Perform systematic code reviews with security, performance, and best practices checks. Analyzes code for vulnerabilities, performance issues, code smells, and provides actionable recommendations with severity levels.
---

# Code Review Assistant

## Overview

This skill provides comprehensive code review capabilities with automated checks for security vulnerabilities, performance issues, code quality, and best practices. It generates detailed reports with actionable recommendations and severity ratings.

## Quick Start

### Review Single File
```bash
python scripts/review_code.py --file app.py --checks all

python scripts/review_code.py --file index.js --checks security,performance
```

### Review Entire Project
```bash
python scripts/review_code.py --project ./src --output review-report.html

python scripts/review_code.py --project . --exclude node_modules,dist --format markdown
```

### Review Git Changes
```bash
python scripts/review_code.py --git-diff HEAD~1 --checks all

python scripts/review_code.py --git-diff main..feature-branch
```

## Review Categories

### 🔒 Security Checks
- SQL Injection vulnerabilities
- XSS (Cross-Site Scripting) risks
- Command injection points
- Path traversal vulnerabilities
- Hardcoded secrets and credentials
- Insecure cryptography usage
- Authentication/authorization issues
- CORS misconfigurations
- Sensitive data exposure
- Dependency vulnerabilities

### ⚡ Performance Analysis
- Algorithm complexity (O(n²), O(n³))
- Database query optimization
- N+1 query problems
- Memory leaks
- Inefficient loops
- Unnecessary re-renders (React/Vue)
- Bundle size issues
- Caching opportunities
- Async/await optimization
- Resource loading

### 🏗️ Code Quality
- Code duplication (DRY violations)
- Cyclomatic complexity
- Method/function length
- Class cohesion
- Dead code detection
- Magic numbers/strings
- Naming conventions
- Comment quality
- Test coverage gaps
- Error handling

### 📋 Best Practices
- SOLID principles violations
- Design pattern opportunities
- API design issues
- Documentation gaps
- Accessibility problems
- Type safety issues
- Dependency management
- Version control practices
- Configuration management
- Logging standards

## Supported Languages

- **JavaScript/TypeScript** - ES6+, Node.js, React, Vue, Angular
- **Python** - 3.6+, Django, Flask, FastAPI
- **Java** - 8+, Spring, Spring Boot
- **C#** - .NET Core, ASP.NET
- **Go** - 1.15+
- **Ruby** - Rails, Sinatra
- **PHP** - 7.4+, Laravel, Symfony
- **Rust** - Security and memory safety
- **SQL** - PostgreSQL, MySQL, SQLite
- **Shell** - Bash, PowerShell

## Security Vulnerability Detection

### SQL Injection
```python
# ❌ VULNERABLE CODE DETECTED
def get_user(user_id):
    query = f"SELECT * FROM users WHERE id = {user_id}"  # SQL Injection risk!
    cursor.execute(query)
    
# ✅ RECOMMENDATION: Use parameterized queries
def get_user(user_id):
    query = "SELECT * FROM users WHERE id = %s"
    cursor.execute(query, (user_id,))
```

### XSS Prevention
```javascript
// ❌ VULNERABLE CODE DETECTED
function displayUserInput(input) {
    document.getElementById('output').innerHTML = input;  // XSS vulnerability!
}

// ✅ RECOMMENDATION: Sanitize user input
function displayUserInput(input) {
    const sanitized = DOMPurify.sanitize(input);
    document.getElementById('output').textContent = sanitized;
}
```

### Authentication Issues
```python
# ❌ SECURITY ISSUE: Weak password hashing
def hash_password(password):
    return hashlib.md5(password.encode()).hexdigest()  # MD5 is broken!

# ✅ RECOMMENDATION: Use bcrypt or argon2
def hash_password(password):
    return bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
```

### Hardcoded Secrets
```javascript
// ❌ CRITICAL: Hardcoded API key
const API_KEY = "sk-1234567890abcdef";  // Never commit secrets!

// ✅ RECOMMENDATION: Use environment variables
const API_KEY = process.env.API_KEY;
```

## Performance Analysis

### Algorithm Complexity
```python
# ❌ PERFORMANCE ISSUE: O(n²) complexity
def find_duplicates(items):
    duplicates = []
    for i in range(len(items)):
        for j in range(i + 1, len(items)):
            if items[i] == items[j]:
                duplicates.append(items[i])
    return duplicates

# ✅ RECOMMENDATION: Use set for O(n) complexity
def find_duplicates(items):
    seen = set()
    duplicates = set()
    for item in items:
        if item in seen:
            duplicates.add(item)
        seen.add(item)
    return list(duplicates)
```

### Database Optimization
```python
# ❌ N+1 QUERY PROBLEM
def get_posts_with_authors():
    posts = Post.objects.all()
    for post in posts:
        print(post.author.name)  # Each iteration queries DB!

# ✅ RECOMMENDATION: Use select_related
def get_posts_with_authors():
    posts = Post.objects.select_related('author').all()
    for post in posts:
        print(post.author.name)  # No additional queries
```

### Memory Leaks
```javascript
// ❌ MEMORY LEAK: Event listener not removed
class Component {
    constructor() {
        window.addEventListener('resize', this.handleResize);
    }
    // Missing cleanup!
}

// ✅ RECOMMENDATION: Clean up listeners
class Component {
    constructor() {
        this.handleResize = this.handleResize.bind(this);
        window.addEventListener('resize', this.handleResize);
    }
    
    destroy() {
        window.removeEventListener('resize', this.handleResize);
    }
}
```

## Code Quality Checks

### DRY Violations
```python
# ❌ CODE DUPLICATION DETECTED
def process_user_data(data):
    if not data.get('email'):
        raise ValueError("Email is required")
    if not '@' in data.get('email', ''):
        raise ValueError("Invalid email format")
    # ... processing

def validate_user_email(email):
    if not email:
        raise ValueError("Email is required")
    if not '@' in email:
        raise ValueError("Invalid email format")
    # ... validation

# ✅ RECOMMENDATION: Extract common validation
def validate_email(email):
    if not email:
        raise ValueError("Email is required")
    if not '@' in email:
        raise ValueError("Invalid email format")
    return True

def process_user_data(data):
    validate_email(data.get('email'))
    # ... processing
```

### Cyclomatic Complexity
```javascript
// ❌ HIGH COMPLEXITY: Score 15 (should be < 10)
function processOrder(order) {
    if (order.type === 'standard') {
        if (order.amount > 100) {
            if (order.customer.isVIP) {
                // ... nested logic continues
            } else {
                if (order.customer.history.length > 5) {
                    // ... more nesting
                }
            }
        }
    }
}

// ✅ RECOMMENDATION: Extract to smaller functions
function processOrder(order) {
    const discount = calculateDiscount(order);
    const shipping = calculateShipping(order);
    return applyOrderRules(order, discount, shipping);
}

function calculateDiscount(order) {
    if (order.customer.isVIP) return 0.2;
    if (order.amount > 100) return 0.1;
    return 0;
}
```

### Error Handling
```python
# ❌ POOR ERROR HANDLING
def divide(a, b):
    return a / b  # Unhandled ZeroDivisionError

def read_file(path):
    return open(path).read()  # File not closed, no error handling

# ✅ RECOMMENDATION: Proper error handling
def divide(a, b):
    if b == 0:
        raise ValueError("Cannot divide by zero")
    return a / b

def read_file(path):
    try:
        with open(path, 'r') as f:
            return f.read()
    except FileNotFoundError:
        logger.error(f"File not found: {path}")
        raise
    except Exception as e:
        logger.error(f"Error reading file: {e}")
        raise
```

## Best Practices Analysis

### SOLID Principles
```python
# ❌ SINGLE RESPONSIBILITY VIOLATION
class UserManager:
    def create_user(self, data): pass
    def send_email(self, user): pass  # Should be in EmailService
    def generate_report(self): pass   # Should be in ReportService
    def validate_password(self): pass

# ✅ RECOMMENDATION: Separate concerns
class UserService:
    def create_user(self, data): pass
    def validate_user(self, user): pass

class EmailService:
    def send_welcome_email(self, user): pass

class ReportService:
    def generate_user_report(self): pass
```

### Dependency Injection
```python
# ❌ TIGHT COUPLING
class OrderService:
    def __init__(self):
        self.db = PostgresDB()  # Hard dependency
        self.email = EmailService()

# ✅ RECOMMENDATION: Dependency injection
class OrderService:
    def __init__(self, db: DatabaseInterface, email: EmailInterface):
        self.db = db
        self.email = email
```

### API Design
```python
# ❌ INCONSISTENT API
@app.route('/getUser/<id>')      # GET with verb in URL
@app.route('/user/delete/<id>')  # DELETE with verb in URL
@app.route('/api/v1/Users')      # Inconsistent casing

# ✅ RECOMMENDATION: RESTful design
@app.route('/users/<id>', methods=['GET'])
@app.route('/users/<id>', methods=['DELETE'])
@app.route('/api/v1/users', methods=['GET'])
```

## Review Report Format

### HTML Report
```html
<!DOCTYPE html>
<html>
<head>
    <title>Code Review Report</title>
    <style>
        .critical { color: #d32f2f; }
        .high { color: #f57c00; }
        .medium { color: #fbc02d; }
        .low { color: #388e3c; }
        .info { color: #1976d2; }
    </style>
</head>
<body>
    <h1>Code Review Report</h1>
    
    <div class="summary">
        <h2>Summary</h2>
        <p>Total Issues: 47</p>
        <ul>
            <li class="critical">Critical: 3</li>
            <li class="high">High: 8</li>
            <li class="medium">Medium: 15</li>
            <li class="low">Low: 21</li>
        </ul>
    </div>
    
    <div class="issues">
        <h2>Security Issues</h2>
        <div class="issue critical">
            <h3>🔴 SQL Injection Vulnerability</h3>
            <p>File: database.py, Line: 45</p>
            <pre><code>query = f"SELECT * FROM users WHERE id = {user_id}"</code></pre>
            <p><strong>Recommendation:</strong> Use parameterized queries</p>
            <pre><code>query = "SELECT * FROM users WHERE id = %s"
cursor.execute(query, (user_id,))</code></pre>
        </div>
    </div>
</body>
</html>
```

### Markdown Report
```markdown
# Code Review Report

## Summary Statistics
- **Files Reviewed**: 23
- **Lines of Code**: 4,567
- **Issues Found**: 47
  - 🔴 Critical: 3
  - 🟠 High: 8
  - 🟡 Medium: 15
  - 🟢 Low: 21

## Critical Issues

### 1. SQL Injection Vulnerability
**File**: `database.py:45`
**Severity**: Critical
**Category**: Security

```python
# Vulnerable Code
query = f"SELECT * FROM users WHERE id = {user_id}"
```

**Recommendation**: Use parameterized queries
```python
query = "SELECT * FROM users WHERE id = %s"
cursor.execute(query, (user_id,))
```

## Performance Issues

### 1. N+1 Query Problem
**File**: `views.py:89-95`
**Severity**: High
**Category**: Performance
```

## Automated Review Rules

### Security Rules
```yaml
security:
  sql_injection:
    patterns:
      - 'execute\(.*f".*{.*}.*"\)'
      - 'execute\(.*\.format\('
    severity: critical
    
  hardcoded_secrets:
    patterns:
      - 'api_key\s*=\s*["\']\w+'
      - 'password\s*=\s*["\']\w+'
      - 'secret\s*=\s*["\']\w+'
    severity: critical
    
  weak_crypto:
    patterns:
      - 'md5\('
      - 'sha1\('
      - 'Random\(\)'
    severity: high
```

### Performance Rules
```yaml
performance:
  nested_loops:
    max_depth: 3
    complexity: O(n^3)
    severity: high
    
  large_methods:
    max_lines: 50
    severity: medium
    
  database_queries:
    patterns:
      - 'select.*from.*where.*in\s*\(select'
      - 'for.*\.objects\.get\('
    severity: high
```

## Custom Rules Configuration

```yaml
# .codereview.yml
review:
  enabled_checks:
    - security
    - performance
    - quality
    - best_practices
    
  security:
    check_dependencies: true
    scan_secrets: true
    owasp_top_10: true
    
  performance:
    max_complexity: 10
    max_method_lines: 50
    warn_on_sync_io: true
    
  quality:
    min_test_coverage: 80
    max_file_lines: 300
    enforce_typing: true
    
  ignore_paths:
    - node_modules/
    - vendor/
    - .git/
    - "*.min.js"
    
  custom_rules:
    - pattern: "console.log"
      message: "Remove console.log before production"
      severity: low
    - pattern: "TODO|FIXME|HACK"
      message: "Unresolved TODO comment"
      severity: info
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
      
      - name: Run Code Review
        run: |
          python scripts/review_code.py \
            --git-diff origin/main..HEAD \
            --output review.md \
            --format markdown
      
      - name: Comment PR
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

# Run code review on staged files
files=$(git diff --cached --name-only --diff-filter=ACM | grep -E '\.(py|js|java)$')

if [ -n "$files" ]; then
    python scripts/review_code.py --files $files --checks security,critical
    
    if [ $? -ne 0 ]; then
        echo "Code review failed. Please fix critical issues before committing."
        exit 1
    fi
fi
```

## Review Metrics

### Code Quality Metrics
```python
class CodeMetrics:
    def calculate_complexity(self, code):
        """McCabe Cyclomatic Complexity"""
        # Count decision points
        complexity = 1
        complexity += code.count('if ')
        complexity += code.count('elif ')
        complexity += code.count('for ')
        complexity += code.count('while ')
        complexity += code.count('except ')
        complexity += code.count('case ')
        return complexity
    
    def calculate_maintainability(self, code):
        """Maintainability Index"""
        loc = len(code.splitlines())
        complexity = self.calculate_complexity(code)
        
        # Simplified maintainability index
        mi = 171 - 5.2 * math.log(loc) - 0.23 * complexity
        return max(0, min(100, mi))
    
    def detect_code_smells(self, code):
        smells = []
        
        # Long method
        if len(code.splitlines()) > 50:
            smells.append("Long method")
        
        # Too many parameters
        params = re.findall(r'def \w+\((.*?)\)', code)
        for param_list in params:
            if len(param_list.split(',')) > 5:
                smells.append("Too many parameters")
        
        # Duplicate code
        lines = code.splitlines()
        for i in range(len(lines) - 10):
            block = lines[i:i+5]
            for j in range(i+5, len(lines) - 5):
                if lines[j:j+5] == block:
                    smells.append("Duplicate code block")
        
        return smells
```

## Language-Specific Checks

### JavaScript/TypeScript
```javascript
// React-specific checks
const reactChecks = {
    // Hook rules violations
    checkHookRules: (code) => {
        const issues = [];
        
        // Hooks in conditions
        if (/if.*use[A-Z]/.test(code)) {
            issues.push({
                type: 'error',
                message: 'React Hook called conditionally'
            });
        }
        
        // Missing dependencies
        const deps = code.match(/useEffect\((.*?)\]/s);
        if (deps && !deps[1].includes('[')) {
            issues.push({
                type: 'warning',
                message: 'useEffect missing dependency array'
            });
        }
        
        return issues;
    },
    
    // Performance issues
    checkPerformance: (code) => {
        const issues = [];
        
        // Inline functions in render
        if (/onClick=\{.*=>/g.test(code)) {
            issues.push({
                type: 'warning',
                message: 'Inline arrow function causes re-renders'
            });
        }
        
        return issues;
    }
};
```

### Python
```python
# Python-specific checks
class PythonAnalyzer:
    def check_type_hints(self, code):
        """Check for missing type hints"""
        issues = []
        
        # Function without type hints
        functions = re.findall(r'def (\w+)\((.*?)\):', code)
        for func_name, params in functions:
            if params and '->' not in code:
                issues.append(f"Function '{func_name}' missing return type hint")
            if params and ':' not in params:
                issues.append(f"Function '{func_name}' missing parameter type hints")
        
        return issues
    
    def check_pythonic_code(self, code):
        """Check for non-Pythonic patterns"""
        issues = []
        
        # Using range(len()) antipattern
        if 'range(len(' in code:
            issues.append("Use enumerate() instead of range(len())")
        
        # Manual file closing
        if 'file.close()' in code:
            issues.append("Use context manager (with statement) for files")
        
        return issues
```

### SQL
```sql
-- SQL-specific checks
CREATE FUNCTION check_sql_issues(query TEXT)
RETURNS TABLE(issue_type TEXT, description TEXT) AS $$
BEGIN
    -- Check for SELECT *
    IF query LIKE '%SELECT * %' THEN
        RETURN QUERY SELECT 'performance', 'Avoid SELECT *, specify columns';
    END IF;
    
    -- Check for missing WHERE in DELETE/UPDATE
    IF (query LIKE '%DELETE FROM%' OR query LIKE '%UPDATE%') 
       AND query NOT LIKE '%WHERE%' THEN
        RETURN QUERY SELECT 'critical', 'DELETE/UPDATE without WHERE clause';
    END IF;
    
    -- Check for LIKE with leading wildcard
    IF query LIKE '%LIKE ''!%%' ESCAPE '!' THEN
        RETURN QUERY SELECT 'performance', 'Leading wildcard prevents index usage';
    END IF;
    
    RETURN;
END;
$$ LANGUAGE plpgsql;
```

## Integration with IDEs

### VS Code Extension Config
```json
{
    "codeReview.enabled": true,
    "codeReview.onSave": true,
    "codeReview.severity": "warning",
    "codeReview.checks": [
        "security",
        "performance",
        "quality"
    ],
    "codeReview.exclude": [
        "**/.git/**",
        "**/node_modules/**",
        "**/*.min.js"
    ],
    "codeReview.customRules": [
        {
            "pattern": "debugger",
            "message": "Remove debugger statement",
            "severity": "error"
        }
    ]
}
```

## Review Checklist Template

```markdown
## Code Review Checklist

### Security
- [ ] No hardcoded credentials
- [ ] Input validation present
- [ ] SQL injection prevention
- [ ] XSS protection
- [ ] CSRF tokens used
- [ ] Authentication properly implemented
- [ ] Authorization checks in place
- [ ] Sensitive data encrypted
- [ ] Security headers configured
- [ ] Dependencies up to date

### Performance
- [ ] Database queries optimized
- [ ] Caching implemented where needed
- [ ] No N+1 query problems
- [ ] Pagination for large datasets
- [ ] Async operations where appropriate
- [ ] Resource cleanup handled
- [ ] Connection pooling used
- [ ] CDN for static assets

### Code Quality
- [ ] Code is DRY
- [ ] Functions are single-purpose
- [ ] Clear naming conventions
- [ ] Adequate comments
- [ ] Error handling complete
- [ ] Unit tests present
- [ ] Integration tests where needed
- [ ] Documentation updated

### Best Practices
- [ ] SOLID principles followed
- [ ] Design patterns used appropriately
- [ ] RESTful conventions
- [ ] Logging implemented
- [ ] Configuration externalized
- [ ] Version control best practices
- [ ] Code formatted consistently
```