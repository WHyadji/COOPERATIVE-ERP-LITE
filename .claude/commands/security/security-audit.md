---
allowed-tools: ReadFile, SearchReplace, Bash(npm:audit), Bash(grep:*)
description: Perform security audit and fix vulnerabilities
---

## Context

- Dependencies audit: !`cd frontend && npm audit 2>&1 | head -20 || echo "No audit results"`
- Auth implementation: @frontend/src/middleware/admin-auth.ts
- Environment variables: !`grep -E "NEXT_PUBLIC_|process.env" frontend/src --include="*.ts*" -r | head -10`
- Security headers: !`grep -E "Content-Security-Policy|X-Frame-Options" frontend/src --include="*.ts*" -r || echo "No security headers found"`

## Your task

Perform security audit for: $ARGUMENTS

Check for:
1. Dependency vulnerabilities
2. SQL injection risks
3. XSS vulnerabilities
4. CSRF protection
5. Authentication bypasses
6. Authorization flaws
7. Sensitive data exposure
8. Insecure direct object references
9. Security misconfigurations
10. Insufficient logging

Implement fixes for:
- Input validation and sanitization
- Proper authentication checks
- Rate limiting on sensitive endpoints
- Secure session management
- HTTPS enforcement
- Security headers
- API key rotation
- Audit logging
- Error message sanitization 