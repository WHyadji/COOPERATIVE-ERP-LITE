
---
description: Comprehensive security vulnerability and weakness scanner
argument-hint: [component/layer] (optional focus area)
allowed-tools: grep_search, codebase_search, read_file
---

# Security Vulnerability Detection and Analysis

Perform a comprehensive security-focused scan to identify vulnerabilities, insecure patterns, and security weaknesses across the codebase.

## Security Scanning Categories

### 1. Hardcoded Secrets Detection
Search for exposed credentials and sensitive data:
```javascript
// Patterns to find:
const API_KEY = "sk_live_51H3..." // API keys in code
const password = "admin123" // Hardcoded passwords
const token = "eyJhbGciOiJIUzI1NiI..." // JWT tokens
const connectionString = "mongodb://user:pass@localhost" // DB credentials

// AWS credentials:
AKIA[0-9A-Z]{16} // Access key IDs
// Private keys:
-----BEGIN (RSA|DSA|EC|OPENSSH) PRIVATE KEY-----
```

### 2. Authentication and Authorization Bypasses
Detect security control bypasses:
```javascript
// Auth bypasses:
// if (!isAuthenticated()) return // Commented out auth
if (true || user.isAuthenticated()) // Always true conditions
const SKIP_AUTH = true // Security flags disabled
if (process.env.NODE_ENV === 'development') return next() // Dev overrides

// Authorization issues:
user.isAdmin = true // Hardcoded privileges
req.user = { id: 1, role: 'admin' } // Static user assignment
```

### 3. Insecure Configuration and Defaults
Find insecure settings:
```javascript
// Insecure defaults:
const config = {
  password: "password123", // Default passwords
  cors: { origin: '*' }, // Permissive CORS
  ssl: { rejectUnauthorized: false }, // Disabled SSL verification
  debug: true, // Debug mode in production
  secretKey: "change-this" // Default secrets
}

// Insecure protocols:
const apiUrl = "http://api.example.com" // HTTP in production
```

### 4. Input Validation Vulnerabilities
Check for injection and validation issues:
```javascript
// SQL injection:
db.query(`SELECT * FROM users WHERE id = ${userId}`) // String interpolation
db.query("SELECT * FROM users WHERE name = '" + userName + "'")

// XSS vulnerabilities:
innerHTML = userInput // Direct DOM manipulation
res.send(`<h1>Hello ${userName}</h1>`) // Unescaped output

// Path traversal:
fs.readFile(req.query.file) // Direct file access
app.use('/uploads', express.static(req.params.path))

// Command injection:
exec(`git clone ${repoUrl}`) // Unsanitized command execution
```

### 5. Cryptographic Weaknesses
Identify weak crypto implementations:
```javascript
// Weak algorithms:
crypto.createHash('md5') // MD5 usage
crypto.createHash('sha1') // SHA1 usage
const encrypted = btoa(password) // Base64 for "encryption"

// Weak randomness:
const token = Math.random().toString(36) // Predictable tokens
const sessionId = Date.now() // Predictable session IDs

// Missing salt:
const hash = crypto.createHash('sha256').update(password).digest()
```

### 6. Access Control and Permission Issues
Find authorization gaps:
```javascript
// Missing access control:
app.get('/admin/*', adminHandler) // No auth middleware
router.delete('/users/:id', deleteUser) // No permission check

// Broken object access:
const user = await User.findById(req.params.id) // No ownership check
res.json(getAllUsers()) // Data leakage

// Exposed interfaces:
app.use('/debug', debugRoutes) // Debug endpoints in production
```

### 7. Information Disclosure and Logging Issues
Detect sensitive data exposure:
```javascript
// Sensitive logging:
console.log('User login:', { email, password }) // Password in logs
logger.info('Payment processed:', paymentData) // PII in logs

// Verbose errors:
res.status(500).json({ error: error.stack }) // Stack trace exposure
catch (e) { res.send(e.toString()) } // Error details leaked

// Debug information:
res.json({ user, debug: { query, params } }) // Debug data in response
```

### 8. Dependency and Third-Party Vulnerabilities
Check for insecure dependencies:
```json
// Outdated packages in package.json:
"express": "3.x" // Very old versions
"lodash": "4.17.15" // Known vulnerable versions

// Unsafe package usage:
eval(require('vm').runInThisContext(code)) // Dynamic code execution
require(userInput) // Dynamic require
```

## Severity Classification

### 🔴 CRITICAL (P0) - Immediate Action Required
- Hardcoded production credentials
- Authentication bypasses
- Remote code execution vectors
- SQL injection vulnerabilities
- Exposed admin interfaces

### 🟡 HIGH (P1) - Urgent Security Risk
- Weak cryptographic implementations
- Missing authorization checks
- XSS vulnerabilities
- Sensitive data exposure
- Privilege escalation paths

### 🟠 MEDIUM (P2) - Security Weaknesses
- Missing security headers
- Weak session management
- Information disclosure
- Missing rate limiting
- Insecure defaults

### 🟢 LOW (P3) - Best Practice Violations
- Debug code in production
- Missing security logging
- Outdated dependencies
- Configuration hardening needed

## Analysis Focus

${ARGUMENTS ? `**Targeted security scan**: ${ARGUMENTS}` : '**Full security assessment** across entire codebase and infrastructure'}

## Expected Output

Provide a comprehensive security report with:

### Executive Dashboard
```
🔒 Security Vulnerability Assessment
===================================
🚨 Critical Issues: X (Immediate remediation required)
⚠️  High Risk: X (Fix within 48 hours)
📋 Medium Risk: X (Address in current sprint)
✅ Low Risk: X (Include in security backlog)

🛡️  OWASP Top 10 Coverage:
   A01 Broken Access Control: X issues
   A02 Cryptographic Failures: X issues
   A03 Injection: X issues
   A04 Insecure Design: X issues
   A05 Security Misconfiguration: X issues
   A06 Vulnerable Components: X issues
   A07 Authentication Failures: X issues
   A08 Software Integrity Failures: X issues
   A09 Logging/Monitoring Failures: X issues
   A10 Server-Side Request Forgery: X issues
```

### Critical Vulnerability Details
For each critical/high severity finding:
- **Location**: `file:line` with code context
- **Vulnerability Type**: OWASP category and specific weakness
- **Attack Vector**: How this could be exploited
- **Business Impact**: Potential damage and scope
- **Proof of Concept**: Example exploitation (if safe)
- **Remediation Steps**: Specific fix instructions
- **Code Fix**: Exact secure implementation

### Security Improvement Roadmap

#### 🚀 Immediate Actions (0-24 hours)
- Rotate exposed credentials
- Patch critical vulnerabilities
- Disable dangerous debug features
- Implement emergency access controls

#### 📈 Short-term Fixes (1-7 days)
- Implement proper input validation
- Add authentication/authorization
- Fix cryptographic issues
- Enable security logging

#### 🏗️ Long-term Hardening (1-4 weeks)
- Security architecture improvements
- Comprehensive security testing
- Security training for developers
- Security monitoring implementation

### Compliance and Standards Assessment
- **PCI DSS**: Payment card security violations
- **GDPR**: Data protection compliance gaps
- **HIPAA**: Healthcare data security issues
- **SOC 2**: Security control deficiencies
- **ISO 27001**: Information security gaps

### Security Metrics and KPIs
- **Security Debt**: Estimated hours to remediate
- **Risk Score**: Overall security posture (0-100)
- **Coverage**: % of codebase with security issues
- **Trend**: Security improvement over time

Begin the comprehensive security vulnerability assessment and provide actionable remediation guidance.
