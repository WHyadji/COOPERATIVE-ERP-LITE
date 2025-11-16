# /solid-refactor

Analyze code for SOLID principle violations and generate a comprehensive refactoring plan following Test-Driven Development (TDD) methodology.

## Overview

The `/solid-refactor` command helps developers transform their code to follow SOLID principles while maintaining test coverage throughout the process. It analyzes existing code, identifies violations, and creates a step-by-step refactoring plan where tests are written before each change.

## SOLID Principles

- **S**ingle Responsibility Principle (SRP): A class should have one reason to change
- **O**pen/Closed Principle (OCP): Open for extension, closed for modification
- **L**iskov Substitution Principle (LSP): Subtypes must be substitutable for their base types
- **I**nterface Segregation Principle (ISP): Clients shouldn't depend on interfaces they don't use
- **D**ependency Inversion Principle (DIP): Depend on abstractions, not concretions

## Usage

```bash
/solid-refactor <file_path> [output_file]
```

### Parameters

- `<file_path>` (required): Path to the Python file to analyze
- `[output_file]` (optional): Custom output filename (default: `solid_refactoring_plan_YYYYMMDD_HHMMSS.md`)

## Examples

### Basic usage
```bash
/solid-refactor src/user_manager.py
```

### With custom output
```bash
/solid-refactor src/payment_processor.py payment_refactoring.md
```

### Complex class refactoring
```bash
/solid-refactor src/services/order_service.py order_service_solid.md
```

## Features

### 1. Comprehensive Code Analysis
- **AST-based analysis**: Deep code structure understanding
- **Pattern detection**: Identifies common anti-patterns
- **Complexity measurement**: Cyclomatic complexity calculation
- **Relationship mapping**: Understands class hierarchies and dependencies

### 2. Violation Detection

#### Single Responsibility Violations
- Classes with too many methods (>7)
- Mixed concerns (e.g., business logic + persistence)
- Long methods (>20 lines)
- God classes (>200 lines)

#### Open/Closed Violations
- Long if-elif chains (>3 branches)
- Type checking in method names
- Switch statements on types

#### Liskov Substitution Violations
- Methods throwing NotImplementedError
- Strengthened preconditions
- Weakened postconditions

#### Interface Segregation Violations
- Fat interfaces (>5 public methods)
- Multiple concerns in one interface
- Unused method dependencies

#### Dependency Inversion Violations
- Concrete class instantiation in constructors
- Direct dependencies on implementations
- Missing abstractions

### 3. TDD Refactoring Process

Each refactoring follows this cycle:
1. **Write characterization tests** - Capture current behavior
2. **Write target tests** - Define desired behavior
3. **Refactor code** - Make tests pass
4. **Verify all tests** - Ensure no regression

### 4. Generated Artifacts

The command generates:
- **Refactoring plan** (markdown): Complete step-by-step guide
- **Characterization tests**: `test_[filename]_characterization.py`
- **Refactoring tests**: `test_[filename]_refactored.py`

## Output Format

The generated plan includes:

### Executive Summary
- Total violations found
- Severity breakdown
- Key improvements
- Time estimates

### SOLID Compliance Score
```
| Principle | Current Score | Target Score | Status |
|-----------|---------------|--------------|---------|
| Single Responsibility (S) | 45% | 85% | 🔴 |
| Open/Closed (O) | 70% | 90% | 🟡 |
| Liskov Substitution (L) | 90% | 95% | 🟢 |
| Interface Segregation (I) | 60% | 90% | 🟡 |
| Dependency Inversion (D) | 40% | 85% | 🔴 |
```

### Refactoring Steps
Each step includes:
- Test code (TDD approach)
- Implementation code
- Time estimate
- Dependencies on other steps

## Refactoring Patterns

### Extract Class (SRP)
```python
# Before: Class with multiple responsibilities
class UserManager:
    def create_user(self, data): ...
    def save_to_db(self, user): ...
    def send_email(self, user): ...
    def validate_email(self, email): ...

# After: Separated responsibilities
class UserService:
    def __init__(self, repository, email_service):
        self.repository = repository
        self.email_service = email_service
    
    def create_user(self, data): ...

class UserRepository:
    def save(self, user): ...

class EmailService:
    def send_welcome_email(self, user): ...
```

### Strategy Pattern (OCP)
```python
# Before: Long if-elif chain
def calculate_discount(customer_type, amount):
    if customer_type == "gold":
        return amount * 0.2
    elif customer_type == "silver":
        return amount * 0.1
    elif customer_type == "bronze":
        return amount * 0.05
    else:
        return 0

# After: Strategy pattern
class DiscountStrategy(ABC):
    @abstractmethod
    def calculate(self, amount): pass

class GoldDiscount(DiscountStrategy):
    def calculate(self, amount):
        return amount * 0.2

# Easy to add new strategies without modifying existing code
```

### Dependency Injection (DIP)
```python
# Before: Concrete dependency
class OrderService:
    def __init__(self):
        self.db = MySQLDatabase()  # Concrete dependency

# After: Dependency injection
class OrderService:
    def __init__(self, repository: RepositoryInterface):
        self.repository = repository  # Abstraction
```

## Best Practices

### 1. Follow the Red-Green-Refactor Cycle
- **Red**: Write a failing test
- **Green**: Write minimal code to pass
- **Refactor**: Improve design while keeping tests green

### 2. Test Behavior, Not Implementation
```python
# Bad: Testing implementation details
def test_uses_mysql_query(self):
    self.assertIn("SELECT *", service._build_query())

# Good: Testing behavior
def test_retrieves_active_users(self):
    users = service.get_active_users()
    self.assertTrue(all(u.is_active for u in users))
```

### 3. Keep Tests Fast and Independent
- Mock external dependencies
- Use in-memory databases for tests
- Each test should run in milliseconds

### 4. Incremental Refactoring
- Make small, safe changes
- Run tests after each change
- Commit frequently

## Common Scenarios

### Scenario 1: God Class
**Problem**: A class that does everything
**Solution**: Extract multiple focused classes
**Time**: 4-6 hours

### Scenario 2: Rigid Design
**Problem**: Hard to add new features
**Solution**: Apply Open/Closed with strategies
**Time**: 2-3 hours

### Scenario 3: Fragile Base Class
**Problem**: Changes break subclasses
**Solution**: Fix inheritance hierarchy (LSP)
**Time**: 3-4 hours

### Scenario 4: Coupled Components
**Problem**: Can't test in isolation
**Solution**: Dependency injection
**Time**: 2-3 hours

## Integration with CI/CD

### Pre-commit Hook
```bash
#!/bin/bash
# Check SOLID compliance before commit
/solid-refactor src/main.py --check-only
if [ $? -ne 0 ]; then
    echo "SOLID violations found. Run /solid-refactor to fix."
    exit 1
fi
```

### GitHub Action
```yaml
name: SOLID Compliance Check
on: [push, pull_request]
jobs:
  solid-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Check SOLID Principles
        run: |
          /solid-refactor src/ --report-only
          if [ -f solid_violations.json ]; then
            echo "::warning::SOLID violations found"
          fi
```

## Troubleshooting

### "Too many violations to refactor"
Start with high-severity violations first:
```bash
/solid-refactor src/file.py --severity high
```

### "Tests are breaking after refactoring"
Run characterization tests first:
```bash
python -m pytest test_file_characterization.py -v
```

### "Unclear which pattern to apply"
The tool suggests patterns, but you can override:
```bash
/solid-refactor src/file.py --prefer-pattern strategy
```

## Metrics and Reporting

The command tracks:
- **Complexity reduction**: Before/after comparison
- **Test coverage**: Ensures coverage doesn't drop
- **SOLID scores**: Quantified compliance
- **Time tracking**: Actual vs estimated

## See Also

- [SOLID Principles Guide](https://docs.example.com/solid)
- [Test-Driven Development](https://docs.example.com/tdd)
- [Refactoring Patterns](https://docs.example.com/patterns)
- [Code Quality Metrics](https://docs.example.com/metrics)

## Changelog

### Version 1.0.0
- Initial release
- AST-based code analysis
- SOLID principle detection
- TDD refactoring plans
- Test generation
- Markdown report output