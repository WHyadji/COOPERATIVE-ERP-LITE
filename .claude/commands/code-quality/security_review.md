# Security Review Command

Perform a comprehensive security review of the codebase or specified files.

## Review Scope
$ARGUMENTS

## Instructions

Please conduct a thorough security review focusing on the following areas:

### 1. Authentication & Authorization
- Check for proper authentication mechanisms
- Verify authorization checks at all access points
- Look for hardcoded credentials or API keys
- Review token handling and session management
- Check for proper password policies and storage

### 2. Input Validation & Sanitization
- Identify all user input points
- Check for SQL injection vulnerabilities
- Look for XSS (Cross-Site Scripting) risks
- Verify proper input validation and sanitization
- Check for command injection possibilities
- Review file upload validations

### 3. Data Protection
- Review encryption implementation for sensitive data
- Check for proper HTTPS/TLS usage
- Verify secure data transmission
- Look for exposed sensitive information in logs
- Check for proper data masking/redaction

### 4. Common Vulnerabilities
- CSRF (Cross-Site Request Forgery) protection
- Directory traversal vulnerabilities
- Insecure deserialization
- XML External Entity (XXE) attacks
- Server-Side Request Forgery (SSRF)
- Race conditions
- Buffer overflows (for compiled languages)

### 5. Dependencies & Third-Party Libraries
- Check for known vulnerabilities in dependencies
- Review dependency versions for security patches
- Look for outdated or deprecated packages
- Verify integrity of third-party components

### 6. Configuration & Deployment
- Review security headers implementation
- Check for exposed debug information
- Verify proper error handling (no stack traces to users)
- Look for misconfigured CORS policies
- Check environment variable handling
- Review secrets management

### 7. API Security
- Rate limiting implementation
- API authentication and authorization
- Input validation for API endpoints
- Proper HTTP method restrictions
- API versioning and deprecation handling

### 8. Code Quality & Best Practices
- Check for code comments containing sensitive info
- Review random number generation for cryptographic use
- Verify proper exception handling
- Look for unsafe type conversions
- Check for proper resource cleanup

### 9. Business Logic Vulnerabilities
- Review access control logic
- Check for privilege escalation paths
- Verify transaction integrity
- Look for logic flaws in critical workflows
- Check for time-of-check to time-of-use (TOCTOU) issues

### 10. Logging & Monitoring
- Verify no sensitive data in logs
- Check for proper security event logging
- Review audit trail completeness
- Ensure log injection prevention

## Output Format

Please provide:

1. **Critical Issues** (Must fix immediately)
   - Issue description
   - Location (file/line)
   - Impact assessment
   - Recommended fix

2. **High Priority Issues** (Fix in next release)
   - Issue description
   - Location
   - Risk level
   - Mitigation strategy

3. **Medium Priority Issues** (Schedule for fixing)
   - Issue description
   - Location
   - Potential impact
   - Best practice recommendation

4. **Low Priority / Informational** (Consider improving)
   - Observation
   - Suggestion for improvement

5. **Security Strengths** (What's done well)
   - Positive security implementations noted

6. **Summary Report**
   - Overall security posture assessment
   - Top 3 recommended actions
   - Compliance considerations (if applicable)

## Additional Context

- Focus on practical, exploitable vulnerabilities
- Provide specific code examples when possible
- Include references to security standards (OWASP, CWE) where relevant
- Consider the technology stack and framework-specific vulnerabilities
- If reviewing specific files/functions mentioned in $ARGUMENTS, prioritize those areas

Begin the security review now.