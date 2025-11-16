---
description: Generate comprehensive bug-fixing todo lists with intelligent code analysis, severity assessment, and step-by-step resolution guidance for your accounting web application
argument-hint: "<bug_description>" <file_path> [line_number] [--output file]
allowed-tools: read_file, write_file, str_replace_editor, run_bash, list_files, view_file
---

Generate comprehensive, intelligent bug-fixing todo lists by analyzing code context, categorizing issues, and providing systematic resolution steps.

## Command Usage

```
/bugfix-plan "<bug_description>" <file_path> [line_number] [--output file]
```

### Arguments
- `bug_description` (required): Detailed description of the bug or issue
- `file_path` (required): Path to the affected source file
- `line_number` (optional): Specific line number where the bug occurs
- `--output file` (optional): Custom output filename for the todo list

### Examples
```
/bugfix-plan "PPH21 tax calculation returns incorrect amount for high salaries" services/accounting-ledger/internal/service/tax_calculation_service.go 245
/bugfix-plan "User authentication fails with 500 error" internal/handlers/auth_handler.go --output auth_fix_plan.md
/bugfix-plan "Database connection pool exhausted during report generation" pkg/database/pool.go 128
```

## Intelligent Bug Analysis System

### Automated Code Context Analysis
The command performs deep analysis of the code surrounding the bug:

#### Function and Class Detection
- Identifies the containing function and class
- Extracts method signatures and parameter types
- Analyzes code structure and dependencies

#### Import and Dependency Mapping
- Maps all file imports and external dependencies
- Identifies related modules and services
- Finds associated test files automatically

#### Code Pattern Recognition
Detects specific code patterns that influence bug fixing approach:
- **Error Handling**: `try/catch`, exception handling patterns
- **Logging**: Logging statements and debug output
- **Database Operations**: Query patterns, transactions, connections
- **API Endpoints**: Request/response handling, routing
- **Validation Logic**: Input validation, assertion patterns
- **State Management**: State updates, data flow patterns

### Bug Categorization System

#### Logic Errors
**Detected by:** Incorrect conditions, wrong operators, unexpected behavior
**Common in Accounting:** Tax calculation errors, rounding mistakes, business rule violations
**Example Fixes:**
- Review conditional logic in tax bracket calculations
- Verify operator precedence in financial formulas
- Add edge case handling for zero amounts

#### Runtime Errors
**Detected by:** Exceptions, crashes, null references
**Common in Accounting:** Data validation failures, type mismatches
**Example Fixes:**
- Add null checks for decimal calculations
- Implement proper error handling for API calls
- Validate input data before processing

#### Performance Issues
**Detected by:** Slow execution, timeouts, memory problems
**Common in Accounting:** Large report generation, complex queries
**Example Fixes:**
- Optimize database queries with proper indexing
- Implement pagination for large datasets
- Add caching for frequently accessed data

#### Security Vulnerabilities
**Detected by:** Authentication issues, data exposure, injection risks
**Common in Accounting:** Unauthorized access to financial data
**Example Fixes:**
- Implement proper input sanitization
- Add authorization checks for sensitive operations
- Encrypt sensitive financial data

#### Data Integrity Issues
**Detected by:** Inconsistent data, validation failures
**Common in Accounting:** Journal entry imbalances, audit trail gaps
**Example Fixes:**
- Add database constraints for referential integrity
- Implement double-entry bookkeeping validation
- Create audit logging for all data changes

### Severity Assessment

#### High Severity
**Triggers:** Data loss, security breaches, production crashes, critical business impact
**Accounting Context:** Financial calculation errors, compliance violations, multi-tenant data leaks
**Response:** Immediate attention, hotfix deployment, stakeholder notification

#### Medium Severity
**Triggers:** Feature malfunctions, performance degradation, user workflow disruption
**Accounting Context:** Report generation failures, integration issues, user interface problems
**Response:** Schedule for next release, thorough testing, documentation updates

#### Low Severity
**Triggers:** Minor UI issues, non-critical functionality, cosmetic problems
**Accounting Context:** Display formatting, minor usability improvements
**Response:** Include in regular maintenance cycle, document for future improvement

## Systematic Todo Generation

### Phase 1: Investigation (Priority 1-3)

#### Task 1: Bug Reproduction
**Objective:** Establish consistent reproduction of the issue
**Accounting Context:** Set up test environment with relevant financial data
**Steps:**
- Configure test database with sample accounting data
- Create test company with appropriate tax settings
- Document exact steps to reproduce the issue
- Capture detailed error logs and system state

#### Task 2: Code Analysis
**Objective:** Understand the root cause through code examination
**Accounting Context:** Review financial calculation logic and data flow
**Steps:**
- Trace execution path through accounting service layers
- Examine business rule implementations
- Check data validation and transformation logic
- Review integration points with external tax services

#### Task 3: Related Component Investigation
**Objective:** Identify all components affected by the bug
**Accounting Context:** Map dependencies across microservices architecture
**Steps:**
- Check related repository implementations
- Review API contract compliance
- Examine event sourcing and audit trail impact
- Verify multi-tenant data isolation

### Phase 2: Testing Foundation (Priority 4)

#### Comprehensive Test Case Creation
**Objective:** Create failing tests that capture the bug
**Accounting Context:** Financial accuracy and compliance testing
**Test Categories:**
- **Unit Tests**: Individual calculation functions
- **Integration Tests**: Service-to-service communication
- **Business Rule Tests**: Compliance with Indonesian tax law
- **Edge Case Tests**: Boundary conditions and error scenarios

**Example Test Structure:**
```go
func TestPPH21Calculation_HighSalaryEdgeCase(t *testing.T) {
    // Arrange
    salary := decimal.NewFromFloat(100000000) // 100M IDR
    allowances := decimal.NewFromFloat(5000000)
    
    // Act
    tax, err := CalculatePPH21(salary, allowances)
    
    // Assert
    assert.NoError(t, err)
    assert.True(t, tax.GreaterThan(decimal.Zero))
    // Verify against expected tax bracket calculation
    expected := calculateExpectedPPH21(salary, allowances)
    assert.Equal(t, expected, tax)
}
```

### Phase 3: Implementation (Priority 5)

#### Intelligent Fix Strategy Generation
Based on bug analysis, generates specific fix approaches:

**For Logic Errors:**
- Review and correct conditional statements
- Add missing edge case handling  
- Verify algorithm correctness against business requirements
- Update calculation formulas with proper decimal precision

**For Runtime Errors:**
- Add comprehensive null/type checks
- Implement graceful error handling with user-friendly messages
- Validate all inputs before processing
- Add retry logic for transient failures

**For Performance Issues:**
- Optimize database queries with proper indexing
- Implement caching strategies for frequently accessed data
- Add pagination for large result sets
- Use async processing for long-running operations

### Phase 4: Validation (Priority 6)

#### Multi-Level Testing Strategy
**Unit Testing:** Individual function validation
**Integration Testing:** Service interaction verification
**System Testing:** End-to-end workflow validation
**Performance Testing:** Load and stress test execution
**Security Testing:** Vulnerability assessment
**Compliance Testing:** Regulatory requirement verification

### Phase 5: Quality Assurance (Priority 7-8)

#### Documentation and Review Preparation
- Add comprehensive inline comments
- Update API documentation
- Create detailed commit messages
- Prepare pull request with context and testing details

#### Specialized Verification (When Applicable)
**Performance Verification:**
- Run profiling tools to measure improvements
- Compare execution times with baseline measurements
- Monitor memory usage patterns
- Validate caching effectiveness

**Security Verification:**
- Execute security scanning tools
- Test input validation boundaries
- Verify authentication and authorization flows
- Check data sanitization and encryption

## For Accounting Web Application Context

### Domain-Specific Bug Patterns

#### Tax Calculation Issues
**Common Problems:**
- Incorrect tax bracket calculations
- Rounding errors in decimal arithmetic
- Missing validation for tax rate changes
- Integration failures with Bank Indonesia API

**Generated Fix Approaches:**
- Verify tax bracket thresholds against current regulations
- Implement banker's rounding for financial calculations
- Add validation for historical tax rate data
- Create fallback mechanisms for external API failures

#### Multi-Tenant Data Issues
**Common Problems:**
- Data leakage between companies
- Incorrect company context in calculations
- Missing tenant isolation in queries

**Generated Fix Approaches:**
- Add company_id filters to all database queries
- Implement row-level security policies
- Create comprehensive tenant isolation tests
- Add logging for cross-tenant access attempts

#### Financial Reporting Errors
**Common Problems:**
- Incorrect balance calculations
- Missing transaction validation
- Report generation timeouts
- Currency conversion errors

**Generated Fix Approaches:**
- Implement double-entry bookkeeping validation
- Add comprehensive transaction integrity checks
- Optimize report queries with proper indexing
- Create currency conversion rate validation

#### Compliance and Audit Issues
**Common Problems:**
- Missing audit trail entries
- Incorrect SAK EP compliance
- Data retention policy violations

**Generated Fix Approaches:**
- Implement comprehensive audit logging
- Add compliance validation rules
- Create data retention cleanup procedures
- Generate compliance reports for verification

## Generated Todo List Structure

### Executive Summary
- Bug description and affected components
- Severity assessment and impact analysis
- Total estimated time and task breakdown
- Priority focus areas

### Categorized Task List
Tasks organized by category with clear priorities:

#### Investigation Tasks
- Bug reproduction with specific steps
- Code analysis with focus areas
- Related component examination

#### Testing Tasks
- Test case creation with examples
- Validation criteria definition
- Edge case identification

#### Implementation Tasks  
- Fix strategy with multiple approaches
- Code modification guidelines
- Integration considerations

#### Quality Assurance Tasks
- Testing validation procedures
- Documentation requirements
- Review preparation steps

### Time Estimation
- Individual task time estimates
- Total project duration
- Dependency mapping
- Resource allocation guidance

### Verification Checklists
Each task includes specific verification steps:
- [ ] Concrete action items
- [ ] Success criteria
- [ ] Testing requirements
- [ ] Documentation needs

## Output Formats

### Console Display
Immediate feedback with formatted todo list showing:
- Priority-ordered tasks
- Time estimates and dependencies
- Category organization
- Summary statistics

### Markdown File Generation
Comprehensive documentation saved as:
- `bugfix_todo_YYYYMMDD_HHMMSS.md` (default)
- Custom filename if specified
- Complete with checkboxes for task tracking
- Professional formatting for team sharing

### Integration Benefits
- **Git Integration**: Ready for commit messages and PR descriptions
- **Project Management**: Compatible with issue tracking systems
- **Team Collaboration**: Sharable format for code reviews
- **Knowledge Base**: Searchable documentation for similar issues

## Arguments Processing

The command processes `$ARGUMENTS` as follows:
1. **Bug Description** (required): Natural language description for intelligent categorization
2. **File Path** (required): Source file location for context analysis
3. **Line Number** (optional): Specific location for focused analysis
4. **Output File** (optional): Custom filename for generated todo list

Use this command to transform bug reports into actionable, systematic fix plans that leverage intelligent code analysis and accounting domain expertise.