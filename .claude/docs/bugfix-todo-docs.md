# /bugfix-todo

Generate comprehensive bug-fixing to-do lists by analyzing code context and creating actionable tasks.

## Overview

The `/bugfix-todo` command helps developers create systematic, context-aware action plans for fixing bugs. It analyzes the surrounding code, identifies potential causes, and generates a prioritized list of tasks with time estimates and verification steps.

## Usage

```bash
/bugfix-todo "<bug_description>" <file_path> [line_number] [options]
```

### Parameters

- `<bug_description>` (required): A description of the bug enclosed in quotes
- `<file_path>` (required): Path to the file containing the bug
- `[line_number]` (optional): The line number where the bug occurs (default: 0)
- `[options]` (optional): Additional command options

### Options

- `--output <filename>`: Specify custom output filename (default: `bugfix_todo_YYYYMMDD_HHMMSS.md`)

## Examples

### Basic usage
```bash
/bugfix-todo "User login fails with null pointer exception" src/auth/login.py
```

### With specific line number
```bash
/bugfix-todo "Array index out of bounds in data processing" src/data/processor.py 145
```

### With custom output file
```bash
/bugfix-todo "Performance issue in search function" src/search/engine.py 230 --output search_bug_fix.md
```

### Complex bug description
```bash
/bugfix-todo "API endpoint returns 500 error when processing large payloads over 10MB" api/handlers/upload.py 89 --output large_payload_bug.md
```

## Features

### 1. Code Context Analysis
- Analyzes ±50 lines around the bug location
- Identifies containing function and class
- Extracts imports and dependencies
- Finds related files and test files
- Detects code patterns (error handling, database operations, etc.)

### 2. Intelligent Bug Analysis
- **Bug Type Classification**: Logic Error, Runtime Error, Performance, Security, Data Corruption, UI/UX, Integration, Concurrency
- **Severity Assessment**: High, Medium, or Low based on bug description and impact
- **Component Identification**: Affected classes, functions, and modules
- **Root Cause Analysis**: Suggests potential causes based on bug type and code patterns

### 3. Task Generation
The command generates tasks in these categories:

#### Investigation
- Bug reproduction steps
- Environment setup
- Log collection

#### Analysis  
- Code flow review
- Variable state inspection
- Dependency checking

#### Testing
- Writing failing test cases
- Edge case identification
- Test-driven development approach

#### Implementation
- Specific fix approaches based on bug type
- Code modification suggestions
- Best practice recommendations

#### Verification
- Unit and integration testing
- Performance validation
- Security checks (when applicable)

#### Documentation
- Code review preparation
- Comment additions
- Commit message drafting

## Output Format

The command generates a markdown file containing:

```markdown
# Bug Fix Todo List

**Bug:** [Description]
**File:** [File path]
**Generated:** [Timestamp]

---

## [Category]

### Task N: [Task description]
**Estimated Time:** [Time estimate]
**Dependencies:** [Previous tasks]

**Steps:**
- [ ] Step 1
- [ ] Step 2
- [ ] Step 3

---

## Summary
- **Total Tasks:** [Count]
- **Estimated Total Time:** [Hours]h [Minutes]min
- **Priority Focus:** [First task]
```

## Bug Type Specific Features

### Logic Errors
- Focuses on conditional logic review
- Suggests edge case handling
- Recommends algorithm verification

### Runtime Errors
- Emphasizes null/type checking
- Proposes error handling improvements
- Suggests input validation

### Performance Issues
- Includes profiling tasks
- Recommends optimization approaches
- Suggests caching strategies

### Security Bugs
- Adds security scanning tasks
- Focuses on input validation
- Includes authentication verification

## Advanced Usage

### Integration with Version Control
```bash
# Generate todo for a bug in a specific commit
git show HEAD:src/module.py > temp.py
/bugfix-todo "Regression from recent commit" temp.py
```

### Batch Processing
```bash
# Process multiple bugs
for bug in bugs/*.txt; do
    desc=$(cat "$bug")
    /bugfix-todo "$desc" src/main.py --output "todos/$(basename $bug .txt).md"
done
```

### Pipeline Integration
```bash
# Use with other tools
/bugfix-todo "$(cat bug_report.txt)" src/app.py | grep "Task 1" | xargs -I {} echo "First step: {}"
```

## Best Practices

1. **Provide Detailed Descriptions**: The more specific your bug description, the better the generated tasks
   ```bash
   # Good
   /bugfix-todo "Login fails with 'TypeError: Cannot read property 'id' of null' when user email contains special characters" auth/login.js 45
   
   # Less helpful
   /bugfix-todo "Login broken" auth/login.js
   ```

2. **Include Line Numbers**: When you know where the bug occurs, always include the line number for better context analysis

3. **Use for Different Bug Types**: The command adapts to various bug types:
   - Crashes and exceptions
   - Performance issues  
   - Security vulnerabilities
   - Data corruption
   - UI/UX problems

4. **Review and Customize**: The generated todo list is a starting point—review and adjust tasks based on your specific needs

5. **Track Progress**: Use the markdown checkboxes to track completion:
   ```markdown
   - [x] Completed task
   - [ ] Pending task
   ```

## Troubleshooting

### Command not found
Ensure the command is properly installed:
```bash
which bugfix-todo
# Should output: /usr/local/bin/bugfix-todo
```

### File not found error
Verify the file path is correct:
```bash
ls -la src/module.py
/bugfix-todo "Bug description" src/module.py
```

### No output generated
Check for Python errors:
```bash
/bugfix-todo "Test bug" file.py 2>&1 | grep -i error
```

## Configuration

The command uses intelligent defaults but can be customized by modifying the source code:

- **Context Lines**: Default ±50 lines (modify in `analyze_code_context`)
- **Task Categories**: Add custom categories in `generate_todo_list`
- **Time Estimates**: Adjust estimates in `TodoItem` creation
- **Bug Patterns**: Extend `code_patterns` dictionary for domain-specific patterns

## See Also

- [Writing Effective Bug Reports](https://docs.example.com/bug-reports)
- [Test-Driven Development Guide](https://docs.example.com/tdd)
- [Code Review Best Practices](https://docs.example.com/code-review)

## Changelog

### Version 1.0.0
- Initial release
- Code context analysis
- Bug categorization
- Task generation with time estimates
- Markdown output format