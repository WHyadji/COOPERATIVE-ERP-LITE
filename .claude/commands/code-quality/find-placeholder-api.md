---
description: Comprehensive API placeholder and incomplete implementation scanner
argument-hint: [component/service] (optional focus area)
allowed-tools: codebase_search, grep_search, read_file
---

# API Placeholder and Incomplete Implementation Analysis

Perform a thorough analysis of API-related code to identify incomplete implementations, placeholders, and potential issues across the codebase.

## Scanning Strategy

Search for the following patterns and issues:

### 1. Endpoint Analysis
Find routes and endpoints with:
- Hardcoded/static responses: `return { data: "test" }`
- TODO/FIXME comments in route handlers
- Missing HTTP methods (incomplete CRUD operations)
- Placeholder status codes (always returning 200)
- Mock data responses
- Unimplemented pagination
- Missing query parameter handling
- Incomplete request body validation

### 2. Error Handling Gaps
Identify missing error handling:
```javascript
// Patterns to find:
.catch(() => {})           // Empty catch blocks
.catch(console.log)        // Only logging errors
res.status(200)            // Always success status
// Missing 4xx/5xx responses
// No error response formatting
```

### 3. Authentication/Authorization Placeholders
Check for security gaps:
- Commented out auth middleware: `// requireAuth,`
- Bypassed auth checks: `if (true || isAuthenticated())`
- Hardcoded tokens/API keys in code
- Missing permission checks
- Placeholder user IDs: `userId: "12345"`
- Missing rate limiting implementation
- Disabled CORS or security middleware

### 4. Data Validation Issues
Look for validation gaps:
- Missing input validation on endpoints
- No schema validation (Joi, Yup, etc.)
- Type coercion without proper checks
- Missing sanitization for user inputs
- Unvalidated file uploads
- No request size limits

### 5. API Documentation Gaps
Identify documentation issues:
- Endpoints without proper documentation
- Outdated API documentation
- Missing request/response examples
- Undocumented error codes and responses
- Missing OpenAPI/Swagger annotations

### 6. Database/External Service Placeholders
Find data layer issues:
- Hardcoded database queries
- Missing database transactions
- No connection pooling configuration
- Commented out external API calls
- Mock service responses in production code
- Missing retry logic for external services
- No timeout handling for external calls

### 7. Testing Gaps
Detect testing issues:
- API endpoints without test coverage
- Skipped test cases: `it.skip()`, `describe.skip()`
- Tests with hardcoded assertions
- Missing edge case and error scenario tests
- No integration tests for API flows
- Incomplete mock setups

## Analysis Focus

${ARGUMENTS ? `Focus analysis on: **${ARGUMENTS}**` : 'Analyze all API-related code across the entire codebase'}

## Expected Output Format

Provide a structured report with:

1. **Executive Summary**
   - Total issues found by severity (Critical/High/Medium/Low)
   - Issues categorized by type
   - Overall API implementation completeness percentage

2. **Critical Issues** (Fix immediately)
   - File location and line number
   - Issue type and description
   - Code snippet showing the problem
   - Business impact assessment
   - Recommended fix

3. **Implementation Roadmap**
   - Prioritized list of fixes
   - Dependencies between fixes
   - Quick wins that can be implemented immediately

4. **API Coverage Analysis**
   - Fully implemented endpoints
   - Partially implemented endpoints
   - Not implemented/placeholder endpoints

Start the analysis now and provide actionable recommendations for improving API implementation quality.
