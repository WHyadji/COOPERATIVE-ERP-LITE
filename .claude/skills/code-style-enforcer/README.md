# Code Style Enforcer Skill

A comprehensive code style enforcement skill that automatically formats code, enforces naming conventions, organizes imports, and maintains consistency across multiple programming languages.

## Features

### 🎨 Auto-Formatting
- Consistent indentation (spaces/tabs)
- Line length enforcement
- Whitespace normalization
- Bracket and parenthesis alignment
- Quote style consistency
- Semicolon enforcement

### 📝 Naming Conventions
- **Python**: snake_case, PascalCase, UPPER_SNAKE_CASE
- **JavaScript**: camelCase, PascalCase, UPPER_SNAKE_CASE
- **Go**: camelCase, PascalCase (exported)
- **Java**: camelCase, PascalCase
- Automatic detection and fixing suggestions

### 📦 Import Organization
- Groups imports by category
- Sorts alphabetically within groups
- Maintains proper spacing
- Removes unused imports
- Standard → Third-party → Local ordering

### ✅ Multiple Style Guides
- Google Style Guide
- Airbnb JavaScript Style
- PEP 8 (Python)
- Standard JS
- PSR-12 (PHP)
- Custom configurations

## Quick Start

### Check Style
```bash
# Check single file
python scripts/enforce_style.py --file app.py --check

# Check entire project
python scripts/enforce_style.py --project ./src --check

# Check with specific style guide
python scripts/enforce_style.py --project . --style google --check
```

### Auto-Fix Issues
```bash
# Fix single file
python scripts/enforce_style.py --file index.js --fix

# Fix entire project
python scripts/enforce_style.py --project . --fix

# Fix staged files before commit
python scripts/enforce_style.py --staged --fix
```

### Generate Reports
```bash
# HTML report
python scripts/enforce_style.py --project . --check --format html --output style-report.html

# Markdown report for PR comments
python scripts/enforce_style.py --git-diff main..feature --format markdown
```

## Supported Languages

| Language | Extensions | Formatter | Linter | Style Guides |
|----------|-----------|-----------|--------|--------------|
| Python | .py | Black, autopep8 | pylint, flake8 | PEP 8, Google |
| JavaScript | .js, .jsx | Prettier, ESLint | ESLint | Airbnb, Standard |
| TypeScript | .ts, .tsx | Prettier | TSLint | Airbnb, Google |
| Go | .go | gofmt | golint | Effective Go |
| Java | .java | google-java-format | Checkstyle | Google, Oracle |
| C/C++ | .c, .cpp, .h | clang-format | cpplint | Google, LLVM |
| Ruby | .rb | RuboCop | RuboCop | Ruby Style Guide |
| PHP | .php | PHP-CS-Fixer | PHPCS | PSR-12, Symfony |
| Rust | .rs | rustfmt | clippy | Rust Style Guide |
| SQL | .sql | sql-formatter | SQLFluff | SQL Style Guide |

## Configuration

### Basic Configuration (.styleguide.yml)
```yaml
style_guide: google
line_length: 100
indent_size: 4
indent_style: space

languages:
  python:
    formatter: black
    line_length: 100
    quote_style: double
    
  javascript:
    formatter: prettier
    indent_size: 2
    semicolons: true
    quote_style: single
```

### Detailed Configuration
```yaml
version: 1.0
extends: airbnb

global:
  line_length: 100
  indent_size: 2
  indent_style: space
  trailing_whitespace: remove
  insert_final_newline: true

python:
  formatter: black
  linter: pylint
  options:
    max_line_length: 100
    docstring_style: google
  naming:
    variables: snake_case
    functions: snake_case
    classes: PascalCase
    constants: UPPER_SNAKE_CASE

javascript:
  formatter: prettier
  linter: eslint
  options:
    max_line_length: 80
    semicolons: true
    trailing_comma: es5
  naming:
    variables: camelCase
    functions: camelCase
    classes: PascalCase
    constants: UPPER_SNAKE_CASE

ignored:
  - node_modules/
  - dist/
  - "*.min.js"
```

## Style Rules Examples

### Python Formatting
```python
# Before
class  my_class:
    def __init__(   self,name:str,    age:int):
        self.name=name;self.age=age
def   CalculateTotal(items):total=0
    for item in items:total+=item
    return total
MAX_value=100

# After (auto-fixed)
class MyClass:
    def __init__(self, name: str, age: int):
        self.name = name
        self.age = age

def calculate_total(items):
    total = 0
    for item in items:
        total += item
    return total

MAX_VALUE = 100
```

### JavaScript Formatting
```javascript
// Before
const   userData={name:"John",age:30}
function   processData(  data  ){
const result=data.map(item=>{return item*2})
return result}
class  user_account{
constructor(name){this.name=name}}

// After (auto-fixed)
const userData = { name: 'John', age: 30 };

function processData(data) {
  const result = data.map((item) => {
    return item * 2;
  });
  return result;
}

class UserAccount {
  constructor(name) {
    this.name = name;
  }
}
```

### Import Organization
```python
# Before
from mymodule import something
import os
from django.db import models
import sys
import numpy as np

# After (auto-fixed)
import os
import sys

import numpy as np
from django.db import models

from mymodule import something
```

## CI/CD Integration

### GitHub Actions
```yaml
name: Code Style
on: [push, pull_request]

jobs:
  style:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Check Code Style
        run: |
          python scripts/enforce_style.py --project . --check
      
      - name: Auto-fix (on main branch)
        if: github.ref == 'refs/heads/main'
        run: |
          python scripts/enforce_style.py --project . --fix
          git config user.name "Style Bot"
          git config user.email "bot@example.com"
          git add -A
          git diff --staged --quiet || git commit -m "Auto-fix code style"
          git push
```

### Pre-commit Hook
```bash
#!/bin/sh
# .git/hooks/pre-commit

python scripts/enforce_style.py --staged --fix
git add -A

python scripts/enforce_style.py --staged --check
if [ $? -ne 0 ]; then
    echo "Style check failed. Please fix issues before committing."
    exit 1
fi
```

### Pre-commit Configuration
```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: code-style
        name: Code Style Check
        entry: python scripts/enforce_style.py
        language: system
        args: ['--staged', '--fix']
        pass_filenames: false
```

## Editor Integration

### VS Code
```json
{
  "editor.formatOnSave": true,
  "python.formatting.provider": "black",
  "python.linting.pylintEnabled": true,
  "eslint.enable": true,
  "eslint.autoFixOnSave": true,
  "[python]": {
    "editor.rulers": [100],
    "editor.tabSize": 4
  },
  "[javascript]": {
    "editor.rulers": [80],
    "editor.tabSize": 2
  }
}
```

### Sublime Text
```json
{
  "rulers": [80, 100],
  "tab_size": 4,
  "translate_tabs_to_spaces": true,
  "trim_trailing_white_space_on_save": true,
  "ensure_newline_at_eof_on_save": true
}
```

## Reports

### Text Report
```
CODE STYLE REPORT
=====================================
Files Processed: 45
Total Issues: 234
Auto-Fixed: 189
Manual Fix Required: 45

ISSUES BY CATEGORY:
- formatting: 120 (105 fixed)
- naming-convention: 67 (45 fixed)
- import-order: 23 (23 fixed)
- whitespace: 24 (16 fixed)
```

### HTML Report
Beautiful HTML report with:
- Summary statistics
- Issue breakdown by file
- Color-coded severity
- Fix suggestions
- Before/after comparisons

### Markdown Report
GitHub/GitLab friendly format:
- Emoji indicators
- Code blocks
- Actionable recommendations
- PR-ready formatting

## Command Options

```bash
python scripts/enforce_style.py [OPTIONS]

Input Options (one required):
  --file FILE          Format single file
  --project PATH       Format entire project
  --git-diff DIFF      Format git changes
  --staged             Format staged files

Actions:
  --check              Check only, don't fix
  --fix                Auto-fix issues

Configuration:
  --config FILE        Configuration file
  --style STYLE        Style guide: google, airbnb, standard, pep8

Output:
  --format FORMAT      Report format: text, json, html, markdown
  --output FILE        Save report to file
```

## Style Guide Comparison

| Feature | Google | Airbnb | Standard | PEP 8 |
|---------|--------|--------|----------|-------|
| Line Length | 100 | 100 | 80 | 79 |
| Indent | 2 spaces | 2 spaces | 2 spaces | 4 spaces |
| Semicolons (JS) | Yes | Yes | No | N/A |
| Quotes (JS) | Single | Single | Single | N/A |
| Trailing Comma | Yes | Yes | No | No |

## Best Practices

1. **Start with auto-fixable rules** - Build confidence with automatic fixes
2. **Agree on style guide as team** - Consistency matters more than personal preference
3. **Use pre-commit hooks** - Catch issues before they enter version control
4. **Format on save** - Configure editors for automatic formatting
5. **Progressive adoption** - Start with new files, gradually fix existing
6. **Document exceptions** - Use comments for legitimate style exceptions
7. **Regular updates** - Keep formatters and linters updated
8. **Performance** - Use incremental formatting for large codebases
9. **CI/CD integration** - Enforce style in continuous integration
10. **Team training** - Educate team on style guidelines

## Common Issues & Solutions

### Issue: Line too long
```python
# Problem
result = very_long_function_name(parameter1, parameter2, parameter3, parameter4, parameter5)

# Solution (auto-fixed)
result = very_long_function_name(
    parameter1, parameter2, parameter3, 
    parameter4, parameter5
)
```

### Issue: Wrong naming convention
```python
# Problem
def GetUserData():  # Should be snake_case
    user_NAME = "John"  # Should be consistent
    
# Solution (suggested)
def get_user_data():
    user_name = "John"
```

### Issue: Unorganized imports
```javascript
// Problem
import Component from './Component';
import React from 'react';
import axios from 'axios';

// Solution (auto-fixed)
import React from 'react';
import axios from 'axios';

import Component from './Component';
```

## Exit Codes

- `0` - Success, no issues or all fixed
- `1` - Style issues found (check mode)
- `2` - Error during execution

## Requirements

- Python 3.7+
- Optional: black, prettier, eslint (for enhanced formatting)
- Git (for git-diff functionality)

## License

This skill is provided for use with Claude AI to enhance code style consistency.