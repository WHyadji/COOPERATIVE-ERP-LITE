---
description: Implement code changes following Test-Driven Development (TDD) with automated test generation and agent guidance for your accounting web application
argument-hint: <file_path> [--issue "description"] [--auto|--continue|--refactor]
allowed-tools: read_file, write_file, str_replace_editor, run_bash, list_files, view_file
---

Implement code changes following Test-Driven Development (TDD) methodology with intelligent test generation and multi-agent guidance.

## Command Usage

```
/implement-tdd <file_path> [options]
```

### Options
- `--issue "description"`: Describe the issue to fix or feature to implement
- `--auto`: Automatic implementation mode (implements all tests automatically) 
- `--continue`: Continue from the last TDD session
- `--refactor`: Start refactoring phase after all tests pass

### Examples
```
/implement-tdd services/accounting-ledger/internal/service/tax_calculation_service.go --issue "Add PPH21 tax validation"
/implement-tdd internal/handlers/account_handler.go --auto
/implement-tdd user_service.py --continue
/implement-tdd payment_processor.go --refactor
```

## TDD Implementation Process

### Phase 1: Red - Write Failing Tests

The command analyzes your code and generates comprehensive test cases following these strategies:

#### Agent-Guided Test Generation
Consults multiple specialist agents for test recommendations:
- **Bugfix Agents**: Analyze potential issues and generate defensive tests
- **TDD Integration Agent**: Provide comprehensive test coverage strategies
- **Security Agent**: Generate tests for input validation and security

#### Code Analysis Test Generation
Automatically generates tests for:
- **Public Functions**: Create basic functionality tests with valid/invalid inputs
- **Validation Functions**: Test with valid data and error conditions
- **Creation Functions**: Verify object creation and initialization
- **Processing Functions**: Test data transformation and business logic

#### Accounting Domain-Specific Tests
Generates specialized tests for:
- **Tax Calculations**: PPH21, PPH22, PPH23, PPN with edge cases
- **Currency Operations**: Decimal precision, rounding, multi-currency
- **Financial Validation**: Account balances, transaction integrity
- **Compliance Rules**: Indonesian accounting standards (SAK EP)

### Phase 2: Green - Minimal Implementation

Implements minimal code to make tests pass:

#### Smart Code Generation
- **Function Stubs**: Generate minimal function implementations based on test requirements
- **Validation Logic**: Create input validation with appropriate error handling
- **Business Logic**: Implement core functionality to satisfy test assertions
- **Error Handling**: Add proper exception handling for edge cases

#### Pattern-Based Implementation
- **Validators**: Return boolean results with proper error messages
- **Creators**: Return objects with required fields (id, status, etc.)
- **Processors**: Transform input data and return processed results
- **Calculators**: Perform computations with proper type checking

### Phase 3: Refactor - Improve Design

After all tests pass, provides guided refactoring:

#### SOLID Principle Application
- **Single Responsibility**: Extract focused, single-purpose functions
- **Open/Closed**: Design for extension without modification
- **Dependency Inversion**: Use interfaces and dependency injection

#### Safe Refactoring Suggestions
1. **Extract Constants**: Replace magic numbers with named constants
2. **Rename Variables**: Use descriptive, domain-specific names  
3. **Extract Methods**: Break complex functions into smaller, focused ones
4. **Remove Duplication**: Consolidate repeated code patterns

## Interactive Implementation Mode

### Step-by-Step Guidance
1. **Current Test Display**: Shows the failing test to implement
2. **Implementation Hint**: Provides guidance on what to implement
3. **Automatic Implementation**: Generates minimal code to pass the test
4. **Progress Tracking**: Shows completion status and next steps
5. **Diff Visualization**: Displays exactly what code was added/changed

### Session Management
- **Session Persistence**: Save progress and resume later with `--continue`
- **Step Navigation**: Move through implementation steps systematically
- **Rollback Capability**: Revert changes if tests fail unexpectedly

## For Accounting Web Application Context

### Domain-Specific Implementation Patterns

#### Tax Calculation Services
```go
// Generated test
func TestCalculatePPH21_ValidSalary(t *testing.T) {
    salary := decimal.NewFromFloat(10000000) // 10M IDR
    tax, err := CalculatePPH21(salary, decimal.Zero)
    
    assert.NoError(t, err)
    assert.True(t, tax.GreaterThan(decimal.Zero))
    assert.Equal(t, decimal.NewFromFloat(500000), tax)
}

// Generated implementation
func CalculatePPH21(grossSalary, taxableAllowances decimal.Decimal) (decimal.Decimal, error) {
    if grossSalary.IsNegative() {
        return decimal.Zero, errors.New("gross salary cannot be negative")
    }
    // Minimal implementation to pass test
    return decimal.NewFromFloat(500000), nil
}
```

#### Account Service Implementation
```go
// Generated test
func TestCreateAccount_Success(t *testing.T) {
    account := &Account{
        Code: "1001",
        Name: "Cash in Bank", 
        Type: "ASSET",
        CompanyID: testCompanyID,
    }
    
    err := accountService.CreateAccount(ctx, account)
    assert.NoError(t, err)
    assert.NotEmpty(t, account.ID)
}

// Generated implementation
func (s *AccountService) CreateAccount(ctx context.Context, account *Account) error {
    if account == nil {
        return errors.New("account cannot be nil")
    }
    account.ID = uuid.New()
    return s.repository.Save(ctx, account)
}
```

#### Financial Validation Implementation
```go
// Generated test
func TestValidateJournalEntry_Balanced(t *testing.T) {
    entry := &JournalEntry{
        Entries: []Entry{
            {AccountID: "1001", DebitAmount: decimal.NewFromFloat(1000)},
            {AccountID: "2001", CreditAmount: decimal.NewFromFloat(1000)},
        },
    }
    
    err := ValidateJournalEntry(entry)
    assert.NoError(t, err)
}

// Generated implementation
func ValidateJournalEntry(entry *JournalEntry) error {
    if entry == nil {
        return errors.New("journal entry cannot be nil")
    }
    
    totalDebits := decimal.Zero
    totalCredits := decimal.Zero
    
    for _, e := range entry.Entries {
        totalDebits = totalDebits.Add(e.DebitAmount)
        totalCredits = totalCredits.Add(e.CreditAmount)
    }
    
    if !totalDebits.Equal(totalCredits) {
        return errors.New("journal entry must be balanced")
    }
    
    return nil
}
```

## Integration Benefits

### Multi-Agent Coordination
- **Comprehensive Analysis**: Multiple agents analyze different aspects of the code
- **Smart Test Generation**: Tests cover security, performance, data integrity
- **Domain Expertise**: Accounting-specific test patterns and validations
- **Quality Assurance**: Automated code quality checks throughout the process

### Development Workflow Integration
- **Git Integration**: Track implementation progress with meaningful commits
- **Test Framework Support**: Works with existing Go testing, pytest, Jest
- **CI/CD Ready**: Generated tests integrate with existing pipeline
- **Code Review Ready**: Clean, documented code with comprehensive test coverage

### Accounting Domain Benefits
- **Compliance Testing**: Ensures Indonesian accounting standards compliance
- **Financial Accuracy**: Decimal precision and rounding tests
- **Multi-Tenancy**: Company isolation and data integrity tests
- **Audit Trail**: Event sourcing and audit logging verification

The command transforms feature requests or bug reports into fully tested, production-ready code following TDD best practices and accounting domain expertise.

## Arguments Processing

The command processes `$ARGUMENTS` as follows:
1. **File Path** (required): The source file to implement TDD on
2. **Options**: Various flags for different modes and behaviors
3. **Issue Description**: Natural language description for agent analysis
4. **Mode Selection**: Interactive, automatic, continue, or refactor modes

Use this command to implement robust, well-tested code that follows TDD principles and accounting domain best practices.