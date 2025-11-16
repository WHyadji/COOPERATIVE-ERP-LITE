---
description: Analyze code for SOLID principle violations and generate comprehensive TDD refactoring plan with step-by-step implementation guidance for your accounting web application
argument-hint: <file_path> [output_file]
allowed-tools: read_file, write_file, str_replace_editor, run_bash, list_files, view_file
---

Analyze code for SOLID principle violations and generate a comprehensive Test-Driven Development (TDD) refactoring plan with detailed implementation steps.

## Command Usage

```
/solid-refactor <file_path> [output_file]
```

### Arguments
- `file_path` (required): Path to the source file to analyze
- `output_file` (optional): Custom output filename for the refactoring plan

### Examples
```
/solid-refactor services/accounting-ledger/internal/service/tax_calculation_service.go
/solid-refactor internal/handlers/user_handler.go custom_refactor_plan.md
/solid-refactor src/payment_processor.py
```

## SOLID Principles Analysis

The command performs comprehensive analysis for all five SOLID principles:

### Single Responsibility Principle (SRP)
**Detects:**
- Classes with too many methods (>7 methods suggests multiple responsibilities)
- Mixed concerns (business logic + persistence in same class)
- Long methods (>20 lines, too complex for single responsibility)
- God classes (>200 lines, handling too many concerns)

**Refactoring Patterns:**
- Extract Class for separating responsibilities
- Extract Method for breaking down complex functions
- Separate concerns into focused components

### Open/Closed Principle (OCP)
**Detects:**
- Complex conditional statements (long if-elif chains >3 branches)
- Switch statements that require modification for new types
- Hardcoded type checking in method names
- Direct type inspection instead of polymorphism

**Refactoring Patterns:**
- Strategy Pattern for replacing conditionals
- Template Method Pattern for shared algorithms
- Polymorphism instead of type checking

### Liskov Substitution Principle (LSP)
**Detects:**
- Methods throwing NotImplementedError in subclasses
- Refused bequest (subclasses rejecting parent functionality)
- Strengthened preconditions in derived classes
- Weakened postconditions in derived classes

**Refactoring Patterns:**
- Rethink inheritance hierarchy
- Use composition over inheritance
- Proper interface design

### Interface Segregation Principle (ISP)
**Detects:**
- Fat interfaces (classes with >5 public methods)
- Clients forced to depend on methods they don't use
- Multi-concern interfaces mixing different responsibilities

**Refactoring Patterns:**
- Split into multiple focused interfaces
- Create role-based interfaces
- Use composition for complex behavior

### Dependency Inversion Principle (DIP)
**Detects:**
- Concrete dependencies in constructors
- Direct instantiation of concrete classes
- High-level modules depending on low-level modules
- Missing abstraction layers

**Refactoring Patterns:**
- Dependency Injection with interfaces
- Abstract Factory Pattern
- Inversion of Control container

## TDD Refactoring Methodology

### Phase 1: Characterization Tests (Red)
**Purpose:** Capture existing behavior before refactoring
- Document current system behavior with comprehensive tests
- Establish performance baselines
- Create safety net for refactoring
- Test edge cases and error conditions

**Generated Tests:**
```python
class CharacterizationTests(unittest.TestCase):
    """Tests that capture current behavior before refactoring"""
    
    def test_current_behavior_with_valid_input(self):
        """Document how the system currently handles valid input"""
        # Comprehensive testing of existing functionality
        
    def test_performance_baseline(self):
        """Establish performance baseline"""
        # Measure current performance for comparison
```

### Phase 2: Design Tests (Red)
**Purpose:** Define desired behavior after refactoring
- Write tests for new class structure
- Test separated concerns individually
- Verify interface contracts
- Test integration between refactored components

### Phase 3: Implementation (Green)
**Purpose:** Implement refactored code to pass all tests
- Extract classes following SRP
- Implement dependency injection
- Create strategy patterns for conditionals
- Segregate fat interfaces

### Phase 4: Integration Verification (Green)
**Purpose:** Ensure refactored system works end-to-end
- Run all characterization tests (should still pass)
- Run new design tests (should pass)
- Verify performance improvements
- Test system integration

## For Accounting Web Application Context

### Domain-Specific SOLID Violations

#### Financial Service Classes
Common violations in accounting services:
- **Mixed Concerns**: Calculation + persistence + validation in single class
- **God Classes**: Tax services handling all tax types (PPH21, PPH22, PPH23, PPN)
- **Concrete Dependencies**: Direct database connections instead of repository interfaces

#### Example Refactoring: Tax Calculation Service
```go
// Before: Violates SRP, DIP
type TaxCalculationService struct {
    db *sql.DB // Concrete dependency
}

func (s *TaxCalculationService) CalculateAllTaxes(salary decimal.Decimal) (*TaxResult, error) {
    // 200+ lines handling PPH21, PPH22, PPH23, PPN, validation, persistence
}

// After: Follows SOLID principles
type TaxCalculationOrchestrator struct {
    pph21Calculator TaxCalculator
    pph22Calculator TaxCalculator  
    pph23Calculator TaxCalculator
    ppnCalculator   TaxCalculator
    repository      TaxRepository
    validator       TaxValidator
}

type TaxCalculator interface {
    Calculate(income decimal.Decimal, allowances decimal.Decimal) (decimal.Decimal, error)
}

type PPH21Calculator struct{}
func (c *PPH21Calculator) Calculate(income, allowances decimal.Decimal) (decimal.Decimal, error) {
    // Focused PPH21 calculation logic only
}
```

#### Microservices Architecture Considerations
- **Service Boundaries**: Ensure each service has single responsibility
- **Interface Contracts**: Use interfaces for service communication
- **Dependency Injection**: Inject repository and external service dependencies
- **Strategy Patterns**: Handle different accounting standards (SAK EP variations)

## Generated Refactoring Plan Structure

### Executive Summary
- Total violations found with severity breakdown
- Estimated improvement in code quality metrics
- Time investment required for refactoring

### SOLID Compliance Scorecard
```
| Principle | Current Score | Target Score | Status |
|-----------|---------------|--------------|---------|
| Single Responsibility (S) | 45% | 85% | =4 |
| Open/Closed (O) | 60% | 90% | =á |
| Liskov Substitution (L) | 80% | 95% | =â |
| Interface Segregation (I) | 40% | 90% | =4 |
| Dependency Inversion (D) | 30% | 85% | =4 |
```

### Step-by-Step Implementation Guide

Each refactoring step includes:

#### Step Information
- **Principle Focus**: Which SOLID principle is being addressed
- **Test-First Approach**: Tests written before implementation
- **Dependencies**: Which previous steps must be completed first
- **Time Estimate**: Realistic implementation time

#### Test Code
Complete test implementations for:
- New class interfaces
- Behavior verification
- Edge cases and error conditions
- Integration testing

#### Implementation Code
Production-ready code with:
- Proper separation of concerns
- Dependency injection setup
- Interface implementations
- Documentation and comments

### Complexity Analysis
- **Before**: Current cyclomatic complexity
- **After**: Projected complexity reduction
- **Improvement**: Percentage improvement in maintainability

## Integration with Development Workflow

### Generated Artifacts
1. **Main Refactoring Plan** (`solid_refactoring_plan_YYYYMMDD_HHMMSS.md`)
2. **Characterization Tests** (`test_filename_characterization.py`)
3. **Refactoring Tests** (`test_filename_refactored.py`)

### Continuous Integration Support
- Tests integrate with existing CI/CD pipelines
- Performance benchmarks for regression detection
- Code quality metrics tracking

### Accounting Domain Best Practices
The refactoring plan includes specific guidance for:

#### Financial Calculation Services
- **Separate calculation logic from persistence**
- **Use strategy pattern for different tax types**
- **Implement decimal precision handling consistently**
- **Create audit trail interfaces**

#### Multi-Tenant Architecture
- **Inject company context through interfaces**
- **Separate tenant isolation concerns**
- **Use repository pattern for data access**

#### Compliance and Reporting
- **Segregate reporting interfaces by audience**
- **Use decorator pattern for different compliance standards**
- **Implement observer pattern for audit logging**

## Usage Examples

### Basic Analysis
```bash
/solid-refactor services/accounting-ledger/internal/service/account_service.go
```
**Output:**
- Comprehensive refactoring plan
- Characterization tests
- Step-by-step TDD implementation guide

### Custom Output File
```bash
/solid-refactor internal/handlers/tax_handler.go tax_refactoring_plan.md
```
**Output:**
- Custom-named refactoring plan
- Focused on tax calculation SOLID violations
- Indonesian tax compliance considerations

### Large Class Refactoring
```bash
/solid-refactor src/financial_calculator.py
```
**Output:**
- God class decomposition plan
- Multiple extracted class designs
- Comprehensive test coverage strategy

The command transforms monolithic, tightly-coupled code into maintainable, testable, and extensible solutions following SOLID principles and accounting domain best practices.

## Arguments Processing

The command processes `$ARGUMENTS` as follows:
1. **File Path** (required): Source file to analyze for SOLID violations
2. **Output File** (optional): Custom filename for the generated refactoring plan
3. **Analysis Mode**: Comprehensive SOLID analysis with TDD implementation guidance

Use this command to systematically improve code architecture, reduce technical debt, and enhance maintainability while preserving existing functionality through comprehensive test coverage.