---
description: Comprehensive placeholder and incomplete code scanner
argument-hint: [directory/file] (optional scope)
allowed-tools: grep_search, codebase_search, read_file
---

# Placeholder Code Detection and Analysis

Scan the codebase to identify incomplete implementations, placeholder code, and development artifacts that need attention.

## Search Patterns

Find the following types of placeholder code:

### 1. Comment-Based Placeholders
- **TODO comments**: `TODO`, `@todo`, `// TODO:`, `# TODO:`, `/* TODO */`
- **FIXME comments**: `FIXME`, `@fixme`, `// FIXME:`, `# FIXME:`, `/* FIXME */`
- **Temporary markers**: `TEMP`, `HACK`, `XXX`, `TEMP:`, `HACK:`, `XXX:`
- **Review requests**: `REVIEW`, `NOTE:`, `WARNING:`

### 2. Implementation Placeholders
- **Unimplemented methods**:
  ```javascript
  throw new Error("Not implemented")
  throw new NotImplementedError()
  return // TODO: implement
  ```
- **Placeholder functions**: Functions with pass statements or empty bodies
- **Mock implementations**: Functions returning static/fake data

### 3. Development Artifacts
- **Debug statements**: `console.log()`, `console.debug()`, `print()`, `console.warn()`
- **Empty catch blocks**: Try-catch with empty or minimal error handling
- **Commented code**: Large blocks of commented-out code

### 4. Test/Development Data
- **Hardcoded test values**: Strings like "test", "example", "placeholder", "dummy", "sample"
- **Placeholder URLs**: "example.com", "localhost", "placeholder.com", "test.com"
- **Mock credentials**: "password123", "testuser", "admin", hardcoded tokens
- **Test email addresses**: "test@example.com", "user@test.com"

### 5. Configuration Placeholders
- **Default/placeholder config**: Default API keys, placeholder environment variables
- **Hardcoded environment settings**: Development URLs in production code
- **Missing validation**: No input validation or sanitization

## Analysis Scope

${ARGUMENTS ? `**Focused scan**: ${ARGUMENTS}` : '**Full codebase scan** across all supported file types'}

Target file extensions: `.js`, `.ts`, `.jsx`, `.tsx`, `.py`, `.java`, `.cs`, `.go`, `.rb`, `.php`, `.swift`, `.kt`, `.sql`, `.md`

## Expected Output

Provide a comprehensive report with:

### Summary Dashboard
```
Placeholder Analysis Summary
============================
Total Issues Found: X
By Priority:
  🔴 Critical (Security/Breaking): X
  🟡 High (Functionality Impact): X
  🟠 Medium (Code Quality): X
  🟢 Low (Cleanup/Polish): X

By Category:
  📝 Comment Placeholders: X
  🚧 Unimplemented Code: X
  🐛 Debug Artifacts: X
  🧪 Test Data: X
  ⚙️  Config Issues: X
```

### Detailed Findings
For each placeholder found:
- **File**: `path/to/file.ext:line_number`
- **Type**: Category and specific pattern
- **Code**: Snippet showing the placeholder
- **Priority**: Critical/High/Medium/Low with justification
- **Impact**: What functionality is affected
- **Action**: Specific next steps to resolve

### Implementation Roadmap
1. **Immediate Actions** (Critical/High priority items)
2. **Sprint Backlog** (Medium priority items that affect functionality)
3. **Cleanup Tasks** (Low priority housekeeping items)

### Code Quality Metrics
- **Implementation completeness**: X% of code is production-ready
- **Technical debt**: Number of items requiring refactoring
- **Security concerns**: Number of security-related placeholders

Begin the comprehensive placeholder scan now and prioritize findings by business impact.
