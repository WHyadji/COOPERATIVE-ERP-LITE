# /tdd-implement

Automatically implements code changes following Test-Driven Development (TDD) principles with integrated agent guidance.

## Overview

The `/tdd-implement` command automates the TDD cycle by:
1. **Consulting agents** for test case suggestions
2. **Writing failing tests** (Red phase)
3. **Implementing minimal code** to pass tests (Green phase)
4. **Refactoring** while keeping tests green (Refactor phase)

This command integrates with all existing agents to provide intelligent implementation guidance.

## Usage

```bash
/tdd-implement <file_path> [options]
```

### Options

- `--issue "description"`: Describe the issue to implement/fix
- `--auto`: Fully automatic implementation mode
- `--continue`: Continue from last session
- `--refactor`: Start refactoring phase (after tests pass)
- `--help`: Show help information

## Examples

### Basic Implementation
```bash
# Implement with agent guidance
/tdd-implement user_service.py --issue "Add email validation"

# Auto-implement all tests
/tdd-implement api_handler.py --auto

# Continue previous session
/tdd-implement user_service.py --continue

# Start refactoring (after implementation)
/tdd-implement user_service.py --refactor
```

### Complex Implementation
```bash
# Fix security issue with full TDD
/tdd-implement auth.py --issue "Fix SQL injection vulnerability"

# Implement new feature
/tdd-implement payment.py --issue "Add refund functionality"
```

## How It Works

### 1. Agent Consultation Phase
When you provide an `--issue`, the command:
- Runs `/bugfix-agents` for bug analysis
- Runs `/solid-refactor` for design analysis  
- Runs `/tdd-integrate` for test strategies
- Combines all recommendations

### 2. Test Generation Phase
Creates test cases from:
- Agent recommendations
- Code analysis (finds untested functions)
- Issue description patterns

### 3. Implementation Phase

#### Interactive Mode (default)
```
Step 1: Implement code to pass test_validate_email
Test to implement:
```python
def test_validate_email_valid():
    """Test email validation with valid input"""
    result = validate_email("user@example.com")
    assert result == True
```
Hint: Function should return True for valid input

Current test status: 🔴 Failing
✅ Test now passes! Implementation added:

```diff
+def validate_email(email):
+    """Validate email format"""
+    if not email or '@' not in email:
+        return False
+    return True
```

Next: Step 2
Run `/tdd-implement user_service.py --continue` to proceed
```

#### Auto Mode
```bash
/tdd-implement user_service.py --auto

🤖 Implementing 5 test cases automatically...

Step 1: test_validate_email ✅
Step 2: test_create_user ✅
Step 3: test_process_payment ✅
Step 4: test_send_notification ❌ Failed - fix manually

Summary:
- Successfully implemented: 3/5
- Success rate: 60%
```

### 4. Refactoring Phase
After all tests pass:
```bash
/tdd-implement user_service.py --refactor

# Refactoring Phase (All Tests Green)

## SOLID Analysis
Single Responsibility: 65% 🟡
Open/Closed: 80% 🟢

## Safe Refactorings
1. Extract Constants ✅ Very Safe
2. Rename Variables ✅ Very Safe  
3. Extract Methods ⚠️ Safe with care
```

## Features

### 1. Intelligent Test Generation
- Analyzes function names to create appropriate tests
- Uses agent recommendations for complex scenarios
- Prioritizes based on code complexity

### 2. Minimal Implementation
The command generates just enough code to pass tests:
```python
# For validation functions → returns boolean
# For creation functions → returns object with id
# For processing functions → returns result dict
```

### 3. Session Management
- Saves progress automatically
- Resume with `--continue`
- Tracks implementation history

### 4. Agent Integration
Seamlessly works with:
- **SecurityAgent**: Security-focused tests
- **PerformanceAgent**: Performance benchmarks
- **DataIntegrityAgent**: Validation tests
- **TestDesignAgent**: Coverage analysis
- **MockingAgent**: Mock setup

## Workflow Integration

### Complete TDD Workflow
```bash
# 1. Analyze the problem
/bugfix-agents "User registration is broken" auth.py

# 2. Get refactoring plan  
/solid-refactor auth.py

# 3. Implement with TDD
/tdd-implement auth.py --issue "Fix user registration"

# 4. Refactor safely
/tdd-implement auth.py --refactor
```

### Continuous Implementation
```bash
# Start implementation
/tdd-implement feature.py --issue "Add new feature"

# Take a break...

# Continue later
/tdd-implement feature.py --continue

# Keep going until done
/tdd-implement feature.py --continue

# Finally refactor
/tdd-implement feature.py --refactor
```

## Output Files

The command creates/modifies:

1. **Test file**: `test_[filename].py`
   - All generated tests
   - Follows pytest conventions

2. **Implementation file**: Your original file
   - Minimal code additions
   - Preserves existing code

3. **Session file**: `.tdd_session_[timestamp].json`
   - Progress tracking
   - Resume capability

## Best Practices

### 1. Describe Issues Clearly
```bash
# Good
/tdd-implement auth.py --issue "Add password strength validation with min 8 chars"

# Too vague
/tdd-implement auth.py --issue "Fix password"
```

### 2. Review Generated Tests
- Check test quality before implementing
- Modify tests if needed
- Ensure tests actually test the requirement

### 3. Use Auto Mode Wisely
- Good for simple implementations
- Review the code it generates
- Switch to interactive for complex logic

### 4. Refactor Incrementally
- Run tests after each refactoring
- Commit working code frequently
- Use version control

## Common Scenarios

### Adding Validation
```bash
/tdd-implement user.py --issue "Validate phone numbers"
# Generates: test_validate_phone_number()
# Implements: validate_phone_number(phone)
```

### Fixing Bugs
```bash
/tdd-implement calculator.py --issue "Division by zero crashes"
# Generates: test_divide_by_zero_handling()
# Implements: Safe division with error handling
```

### Adding Features
```bash
/tdd-implement api.py --issue "Add pagination to list endpoint"
# Generates: test_list_with_pagination()
# Implements: Pagination logic
```

## Integration with Other Commands

### Before Implementation
```bash
# Get comprehensive analysis
/bugfix-agents "issue" file.py > analysis.md
/solid-refactor file.py > refactoring_plan.md

# Then implement with all insights
/tdd-implement file.py --issue "issue"
```

### After Implementation  
```bash
# Verify improvements
/solid-refactor file.py  # Check improved scores
/bugfix-agents "verify fix" file.py  # Confirm issue resolved
```

## Troubleshooting

### "Test already passes"
- The functionality might already exist
- Check if test is testing the right thing
- Modify test to be more specific

### "Cannot generate implementation"
- Switch to interactive mode
- Implement manually
- Use `--continue` after fixing

### "Tests still failing after implementation"
- Check test expectations
- Verify implementation logic
- Use debugger to trace issue

## Advanced Usage

### Custom Test Templates
Create `.tdd_templates.py` in project root:
```python
TEST_TEMPLATES = {
    'api': '''def test_{name}_api():
    response = client.get('/{name}')
    assert response.status_code == 200''',
    
    'model': '''def test_{name}_model():
    instance = {name}()
    assert instance.id is not None'''
}
```

### Batch Implementation
```bash
# Process multiple files
for file in src/*.py; do
    /tdd-implement "$file" --auto
done
```

## See Also

- [TDD Guide](https://docs.example.com/tdd)
- [Agent System Docs](https://docs.example.com/agents)
- [Refactoring Best Practices](https://docs.example.com/refactoring)

## Changelog

### Version 1.0.0
- Initial release
- Agent integration
- Auto-implementation mode
- Session management
- Refactoring phase