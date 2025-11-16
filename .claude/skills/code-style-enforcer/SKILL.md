---
name: Code Style Enforcer
description: Apply consistent formatting, linting, and naming conventions across languages. Automatically format code, enforce style guides, check naming patterns, and maintain consistency across Python, JavaScript, TypeScript, Java, Go, and more.
---

# Code Style Enforcer

## Overview

This skill enforces consistent code style across multiple programming languages. It automatically formats code, applies linting rules, enforces naming conventions, and maintains consistency throughout your codebase with support for popular style guides like PEP8, Airbnb, Google, and Standard.

## Quick Start

### Format Single File
```bash
python scripts/enforce_style.py --file app.py --fix

python scripts/enforce_style.py --file index.js --style airbnb --fix
```

### Format Entire Project
```bash
python scripts/enforce_style.py --project ./src --fix

python scripts/enforce_style.py --project . --exclude node_modules,dist --fix
```

### Check Without Fixing
```bash
python scripts/enforce_style.py --project . --check

python scripts/enforce_style.py --file main.go --check --report style-report.json
```

## Supported Languages & Style Guides

### Python
- **PEP 8** - Standard Python style guide
- **Black** - Opinionated formatter
- **Google Python Style**
- **NumPy Style** - For scientific computing

### JavaScript/TypeScript
- **Airbnb** - Most popular style guide
- **Standard** - No configuration needed
- **Google JavaScript Style**
- **Prettier** - Opinionated formatter

### Java
- **Google Java Style**
- **Sun/Oracle Conventions**
- **Spring Framework Style**

### Go
- **gofmt** - Standard Go formatting
- **Effective Go** - Best practices

### C/C++
- **Google C++ Style**
- **LLVM Coding Standards**
- **Linux Kernel Style**

### Additional Languages
- **Ruby** - RuboCop standards
- **PHP** - PSR-12 standards
- **Rust** - rustfmt conventions
- **C#** - Microsoft conventions
- **Swift** - Swift.org API guidelines

## Formatting Rules

### Indentation & Spacing
```python
# ❌ Before
def calculate(x,y,z):
 if x>0:
   result=x*y+z
   return    result

# ✅ After
def calculate(x, y, z):
    if x > 0:
        result = x * y + z
        return result
```

### Line Length
```javascript
// ❌ Before (line too long)
const result = someVeryLongFunctionName(argumentOne, argumentTwo, argumentThree, argumentFour, argumentFive);

// ✅ After
const result = someVeryLongFunctionName(
    argumentOne,
    argumentTwo,
    argumentThree,
    argumentFour,
    argumentFive
);
```

### Import Organization
```python
# ❌ Before
import os
from myapp import models
import sys
from django.db import models as django_models
import json
from .utils import helper

# ✅ After
# Standard library imports
import json
import os
import sys

# Third-party imports
from django.db import models as django_models

# Local imports
from myapp import models
from .utils import helper
```

## Naming Conventions

### Python Naming
```python
# ❌ Incorrect naming
class myClass:
    def MyMethod(self):
        MyVariable = 10
        CONSTANT_value = 20

# ✅ Correct naming (PEP 8)
class MyClass:
    def my_method(self):
        my_variable = 10
        CONSTANT_VALUE = 20
```

### JavaScript/TypeScript Naming
```javascript
// ❌ Incorrect naming
const user_name = "John";
function GetUserData() {}
class user_manager {}
const APIKey = "secret";

// ✅ Correct naming
const userName = "John";
function getUserData() {}
class UserManager {}
const API_KEY = "secret";
```

### Variable Naming Patterns
```yaml
naming_conventions:
  python:
    classes: PascalCase
    functions: snake_case
    variables: snake_case
    constants: UPPER_SNAKE_CASE
    private: _leading_underscore
    
  javascript:
    classes: PascalCase
    functions: camelCase
    variables: camelCase
    constants: UPPER_SNAKE_CASE
    components: PascalCase  # React/Vue
    
  java:
    classes: PascalCase
    methods: camelCase
    variables: camelCase
    constants: UPPER_SNAKE_CASE
    packages: lowercase
```

## Code Organization

### File Structure
```python
"""
Module docstring describing the purpose of this module.
"""

# Imports (organized by type)
import standard_library
import third_party
import local_modules

# Constants
CONSTANT_VALUE = 100

# Module-level variables
module_variable = None

# Classes
class MyClass:
    """Class docstring."""
    pass

# Functions
def my_function():
    """Function docstring."""
    pass

# Main execution
if __name__ == "__main__":
    main()
```

### Class Organization
```python
class WellOrganizedClass:
    """Class with proper organization."""
    
    # Class variables
    class_variable = None
    
    def __init__(self):
        """Constructor."""
        # Instance variables
        self.public_variable = None
        self._protected_variable = None
        self.__private_variable = None
    
    # Special methods
    def __str__(self):
        pass
    
    def __repr__(self):
        pass
    
    # Properties
    @property
    def my_property(self):
        pass
    
    # Public methods
    def public_method(self):
        pass
    
    # Protected methods
    def _protected_method(self):
        pass
    
    # Private methods
    def __private_method(self):
        pass
    
    # Static methods
    @staticmethod
    def static_method():
        pass
    
    # Class methods
    @classmethod
    def class_method(cls):
        pass
```

## Documentation Standards

### Python Docstrings (Google Style)
```python
def complex_function(param1: str, param2: int, param3: List[str] = None) -> Dict[str, Any]:
    """
    Brief description of function.
    
    Longer description explaining the function's behavior,
    assumptions, and any important details.
    
    Args:
        param1: Description of param1.
        param2: Description of param2.
        param3: Optional list of strings. Defaults to None.
    
    Returns:
        Dictionary containing:
            - 'result': The processed result
            - 'status': Status code
            - 'data': Additional data
    
    Raises:
        ValueError: If param2 is negative.
        TypeError: If param1 is not a string.
    
    Example:
        >>> result = complex_function("test", 42)
        >>> print(result['status'])
        200
    """
    pass
```

### JavaScript JSDoc
```javascript
/**
 * Calculate the sum of two numbers.
 * 
 * @param {number} a - The first number.
 * @param {number} b - The second number.
 * @returns {number} The sum of a and b.
 * @throws {TypeError} If parameters are not numbers.
 * 
 * @example
 * const result = add(5, 3);
 * console.log(result); // 8
 */
function add(a, b) {
    if (typeof a !== 'number' || typeof b !== 'number') {
        throw new TypeError('Parameters must be numbers');
    }
    return a + b;
}
```

## Configuration Files

### .editorconfig (Universal)
```ini
# .editorconfig
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.py]
indent_style = space
indent_size = 4
max_line_length = 88

[*.{js,jsx,ts,tsx}]
indent_style = space
indent_size = 2
max_line_length = 80

[*.{java,c,cpp,cs}]
indent_style = space
indent_size = 4

[*.go]
indent_style = tab
indent_size = 4

[*.md]
trim_trailing_whitespace = false
```

### Python Configuration (pyproject.toml)
```toml
[tool.black]
line-length = 88
target-version = ['py38']
include = '\.pyi?$'

[tool.isort]
profile = "black"
line_length = 88
multi_line_output = 3
include_trailing_comma = true

[tool.pylint]
max-line-length = 88
disable = ["C0111", "R0903"]
good-names = ["i", "j", "k", "df", "ex", "_"]

[tool.mypy]
python_version = "3.8"
warn_return_any = true
warn_unused_configs = true
disallow_untyped_defs = true
```

### JavaScript/TypeScript (.eslintrc.json)
```json
{
  "extends": ["airbnb-base"],
  "env": {
    "es2021": true,
    "node": true
  },
  "parserOptions": {
    "ecmaVersion": 12,
    "sourceType": "module"
  },
  "rules": {
    "indent": ["error", 2],
    "linebreak-style": ["error", "unix"],
    "quotes": ["error", "single"],
    "semi": ["error", "always"],
    "no-unused-vars": "error",
    "no-console": "warn",
    "prefer-const": "error",
    "arrow-spacing": "error",
    "object-curly-spacing": ["error", "always"],
    "comma-dangle": ["error", "always-multiline"],
    "max-len": ["error", { "code": 80 }],
    "naming-convention": [
      "error",
      {
        "selector": "variable",
        "format": ["camelCase", "UPPER_CASE"]
      },
      {
        "selector": "function",
        "format": ["camelCase"]
      },
      {
        "selector": "class",
        "format": ["PascalCase"]
      }
    ]
  }
}
```

### Prettier Configuration (.prettierrc)
```json
{
  "printWidth": 80,
  "tabWidth": 2,
  "useTabs": false,
  "semi": true,
  "singleQuote": true,
  "quoteProps": "as-needed",
  "jsxSingleQuote": false,
  "trailingComma": "es5",
  "bracketSpacing": true,
  "jsxBracketSameLine": false,
  "arrowParens": "always",
  "endOfLine": "lf"
}
```

## Style Enforcement Rules

### Python Rules
```python
class PythonStyleEnforcer:
    """Enforce Python style rules."""
    
    def check_indentation(self, code: str) -> List[Issue]:
        """Check for consistent 4-space indentation."""
        issues = []
        for i, line in enumerate(code.split('\n'), 1):
            if line and line[0] == ' ':
                indent = len(line) - len(line.lstrip())
                if indent % 4 != 0:
                    issues.append(Issue(
                        line=i,
                        message=f"Indentation not a multiple of 4 (found {indent} spaces)"
                    ))
            elif line and line[0] == '\t':
                issues.append(Issue(
                    line=i,
                    message="Tabs used for indentation (use 4 spaces)"
                ))
        return issues
    
    def check_naming(self, code: str) -> List[Issue]:
        """Check PEP 8 naming conventions."""
        issues = []
        
        # Class names should be PascalCase
        class_pattern = r'class\s+([a-z_][a-zA-Z0-9_]*)'
        for match in re.finditer(class_pattern, code):
            class_name = match.group(1)
            if not class_name[0].isupper():
                issues.append(Issue(
                    line=code[:match.start()].count('\n') + 1,
                    message=f"Class '{class_name}' should be PascalCase"
                ))
        
        # Function names should be snake_case
        func_pattern = r'def\s+([A-Z][a-zA-Z0-9_]*)'
        for match in re.finditer(func_pattern, code):
            func_name = match.group(1)
            issues.append(Issue(
                line=code[:match.start()].count('\n') + 1,
                message=f"Function '{func_name}' should be snake_case"
            ))
        
        return issues
    
    def check_imports(self, code: str) -> List[Issue]:
        """Check import organization."""
        import_lines = []
        lines = code.split('\n')
        
        for i, line in enumerate(lines):
            if line.startswith('import ') or line.startswith('from '):
                import_lines.append((i, line))
        
        # Check if imports are grouped properly
        if import_lines:
            issues = []
            last_type = None
            for i, line in import_lines:
                current_type = self._get_import_type(line)
                if last_type and current_type < last_type:
                    issues.append(Issue(
                        line=i + 1,
                        message="Imports not properly organized (standard → third-party → local)"
                    ))
                last_type = current_type
        
        return issues
```

### JavaScript/TypeScript Rules
```javascript
class JavaScriptStyleEnforcer {
    checkIndentation(code) {
        const issues = [];
        const lines = code.split('\n');
        
        lines.forEach((line, index) => {
            if (line.match(/^\t/)) {
                issues.push({
                    line: index + 1,
                    message: 'Use spaces instead of tabs'
                });
            }
            
            const indent = line.match(/^( *)/)[1].length;
            if (indent % 2 !== 0) {
                issues.push({
                    line: index + 1,
                    message: `Indentation should be 2 spaces (found ${indent})`
                });
            }
        });
        
        return issues;
    }
    
    checkSemicolons(code) {
        const issues = [];
        const lines = code.split('\n');
        
        lines.forEach((line, index) => {
            const trimmed = line.trim();
            if (trimmed && 
                !trimmed.endsWith(';') && 
                !trimmed.endsWith('{') && 
                !trimmed.endsWith('}') &&
                !trimmed.startsWith('//') &&
                !trimmed.includes('if') &&
                !trimmed.includes('for') &&
                !trimmed.includes('while')) {
                issues.push({
                    line: index + 1,
                    message: 'Missing semicolon'
                });
            }
        });
        
        return issues;
    }
    
    checkNaming(code) {
        const issues = [];
        
        // Check for snake_case variables (should be camelCase)
        const snakeCaseVar = /(?:const|let|var)\s+([a-z]+_[a-z_]+)/g;
        let match;
        while ((match = snakeCaseVar.exec(code)) !== null) {
            const line = code.substring(0, match.index).split('\n').length;
            issues.push({
                line,
                message: `Variable '${match[1]}' should be camelCase`
            });
        }
        
        // Check for lowercase class names
        const classPattern = /class\s+([a-z][a-zA-Z0-9]*)/g;
        while ((match = classPattern.exec(code)) !== null) {
            const line = code.substring(0, match.index).split('\n').length;
            issues.push({
                line,
                message: `Class '${match[1]}' should be PascalCase`
            });
        }
        
        return issues;
    }
}
```

## Auto-Formatting

### Python Auto-Formatter
```python
def format_python_code(code: str, style: str = "black") -> str:
    """Auto-format Python code."""
    
    if style == "black":
        import black
        return black.format_str(code, mode=black.Mode())
    
    elif style == "autopep8":
        import autopep8
        return autopep8.fix_code(code)
    
    elif style == "yapf":
        from yapf.yapflib.yapf_api import FormatCode
        return FormatCode(code, style_config='google')[0]
    
    return code

def organize_imports(code: str) -> str:
    """Organize Python imports."""
    import isort
    return isort.code(code)
```

### JavaScript Auto-Formatter
```javascript
function formatJavaScriptCode(code, style = 'prettier') {
    if (style === 'prettier') {
        const prettier = require('prettier');
        return prettier.format(code, {
            parser: 'babel',
            singleQuote: true,
            semi: true,
            tabWidth: 2,
        });
    }
    
    if (style === 'standard') {
        const standard = require('standard');
        const { results } = standard.lintTextSync(code, { fix: true });
        return results[0].output || code;
    }
    
    return code;
}
```

## Whitespace & Formatting

### Trailing Whitespace
```python
# ❌ Before (trailing spaces shown as •)
def function():••
    x = 10•••
    return x••

# ✅ After
def function():
    x = 10
    return x
```

### Blank Lines
```python
# ❌ Before (no blank lines)
import os
import sys
class MyClass:
    def method1(self):
        pass
    def method2(self):
        pass
def function():
    pass

# ✅ After (PEP 8 spacing)
import os
import sys


class MyClass:
    
    def method1(self):
        pass
    
    def method2(self):
        pass


def function():
    pass
```

### Line Breaks
```javascript
// ❌ Before
const obj = {a: 1, b: 2, c: 3, d: 4, e: 5};
const arr = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10];

// ✅ After
const obj = {
  a: 1,
  b: 2,
  c: 3,
  d: 4,
  e: 5,
};

const arr = [
  1, 2, 3, 4, 5,
  6, 7, 8, 9, 10,
];
```

## Comments & Documentation

### Comment Style
```python
# ❌ Bad comments
x = x + 1  # Increment x by 1
# Set name to John
name = "John"

#TODO: fix this
broken_function()

# ✅ Good comments
# Calculate compound interest using the formula:
# A = P(1 + r/n)^(nt)
amount = principal * (1 + rate/n) ** (n * time)

# TODO(john): Implement caching by 2024-02-01
# See issue #123 for details

# HACK: Workaround for third-party library bug
# Remove when library version > 2.0 is released
```

### Inline Documentation
```python
def process_data(
    data: List[Dict[str, Any]],
    validate: bool = True,
    transform: Optional[Callable] = None,
) -> ProcessResult:
    """
    Process raw data with optional validation and transformation.
    
    This function handles the complete data processing pipeline including
    validation, transformation, and aggregation. It's designed to be
    flexible and extensible through the transform parameter.
    
    Args:
        data: List of dictionaries containing raw data. Each dictionary
            must have 'id' and 'value' keys at minimum.
        validate: Whether to validate data before processing. If True,
            invalid records are logged and skipped.
        transform: Optional callable to transform each record. Should
            accept a dictionary and return a dictionary.
    
    Returns:
        ProcessResult containing:
            - processed_count: Number of successfully processed records
            - error_count: Number of failed records
            - results: List of processed data
            - errors: List of error messages
    
    Raises:
        ValueError: If data is empty or None
        TypeError: If data is not a list of dictionaries
    
    Example:
        >>> data = [{'id': 1, 'value': 100}, {'id': 2, 'value': 200}]
        >>> result = process_data(data, validate=True)
        >>> print(f"Processed {result.processed_count} records")
        Processed 2 records
    
    Note:
        This function is thread-safe and can be used in parallel processing.
        For large datasets (>10,000 records), consider using process_data_batch()
        instead for better performance.
    """
    pass
```

## CI/CD Integration

### Pre-commit Configuration
```yaml
# .pre-commit-config.yaml
repos:
  # Python
  - repo: https://github.com/psf/black
    rev: 22.10.0
    hooks:
      - id: black
        language_version: python3.9

  - repo: https://github.com/PyCQA/isort
    rev: 5.10.1
    hooks:
      - id: isort

  - repo: https://github.com/PyCQA/flake8
    rev: 5.0.4
    hooks:
      - id: flake8
        args: ['--max-line-length=88']

  # JavaScript/TypeScript
  - repo: https://github.com/pre-commit/mirrors-prettier
    rev: v2.7.1
    hooks:
      - id: prettier
        files: \.(js|jsx|ts|tsx|json|css|scss|md)$

  - repo: https://github.com/pre-commit/mirrors-eslint
    rev: v8.24.0
    hooks:
      - id: eslint
        files: \.(js|jsx|ts|tsx)$

  # General
  - repo: https://github.com/pre-commit/pre-commit-hooks
    rev: v4.3.0
    hooks:
      - id: trailing-whitespace
      - id: end-of-file-fixer
      - id: check-yaml
      - id: check-json
      - id: check-merge-conflict
      - id: check-added-large-files
        args: ['--maxkb=500']
```

### GitHub Actions Workflow
```yaml
name: Code Style Check

on: [push, pull_request]

jobs:
  style-check:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Python
      uses: actions/setup-python@v4
      with:
        python-version: '3.9'
    
    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '16'
    
    - name: Install dependencies
      run: |
        pip install black isort flake8 pylint
        npm install -g prettier eslint
    
    - name: Check Python style
      run: |
        black --check .
        isort --check-only .
        flake8 .
    
    - name: Check JavaScript style
      run: |
        prettier --check "**/*.{js,jsx,ts,tsx}"
        eslint "**/*.{js,jsx,ts,tsx}"
    
    - name: Auto-fix (on push to main)
      if: github.ref == 'refs/heads/main'
      run: |
        black .
        isort .
        prettier --write "**/*.{js,jsx,ts,tsx}"
        git config --local user.email "action@github.com"
        git config --local user.name "GitHub Action"
        git add -A
        git diff --staged --quiet || git commit -m "Auto-format code"
        git push
```

## Style Report Format

### JSON Report
```json
{
  "summary": {
    "files_checked": 45,
    "files_with_issues": 12,
    "total_issues": 89,
    "auto_fixed": 67,
    "manual_required": 22
  },
  "issues": [
    {
      "file": "src/main.py",
      "line": 45,
      "column": 12,
      "rule": "indentation",
      "severity": "error",
      "message": "Expected 4 spaces, found 3",
      "fixed": true
    },
    {
      "file": "src/utils.js",
      "line": 23,
      "column": 8,
      "rule": "naming-convention",
      "severity": "warning",
      "message": "Variable 'user_name' should be camelCase",
      "fixed": false,
      "suggestion": "userName"
    }
  ],
  "statistics": {
    "by_rule": {
      "indentation": 23,
      "line-length": 15,
      "naming": 18,
      "imports": 12,
      "whitespace": 21
    },
    "by_language": {
      "python": 34,
      "javascript": 28,
      "typescript": 27
    }
  }
}
```

## Custom Style Rules

```yaml
# .style-enforcer.yml
rules:
  # Global rules
  global:
    max_line_length: 100
    indent_size: 4
    indent_style: space
    end_of_line: lf
    charset: utf-8
    trim_trailing_whitespace: true
    insert_final_newline: true

  # Python specific
  python:
    style_guide: pep8
    formatter: black
    max_line_length: 88
    docstring_style: google
    import_order:
      - standard
      - third_party
      - first_party
      - local
    naming:
      classes: PascalCase
      functions: snake_case
      constants: UPPER_SNAKE_CASE

  # JavaScript/TypeScript
  javascript:
    style_guide: airbnb
    formatter: prettier
    indent_size: 2
    quotes: single
    semicolons: true
    trailing_comma: es5
    naming:
      classes: PascalCase
      functions: camelCase
      constants: UPPER_SNAKE_CASE
      react_components: PascalCase

  # Custom patterns
  custom_rules:
    - name: no-console-log
      pattern: 'console\.log'
      message: "Remove console.log statements"
      severity: warning
      languages: [javascript, typescript]
      
    - name: no-print
      pattern: '^print\('
      message: "Use logging instead of print"
      severity: warning
      languages: [python]
      
    - name: todo-format
      pattern: 'TODO(?!\\([a-z]+\\):)'
      message: "TODO should include author: TODO(name):"
      severity: info
      languages: all
```

## Best Practices

1. **Consistent Configuration** - Use same rules across team
2. **Automated Formatting** - Format on save in IDE
3. **Pre-commit Hooks** - Catch issues before commit
4. **CI/CD Integration** - Block PRs with style violations
5. **Progressive Adoption** - Start with warnings, move to errors
6. **Team Agreement** - Document and agree on style guide
7. **Regular Updates** - Keep tools and rules up to date
8. **Language-Specific** - Use appropriate tools per language
9. **Editor Integration** - Configure IDEs for consistency
10. **Documentation** - Document custom rules and exceptions