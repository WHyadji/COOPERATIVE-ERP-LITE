# Production Deployment Security Guide

**Date**: November 19, 2025
**Project**: Cooperative ERP Lite
**Version**: 1.0.0

## Table of Contents

1. [Pre-Deployment Security Checklist](#pre-deployment-security-checklist)
2. [Environment Configuration](#environment-configuration)
3. [Security Middleware Setup](#security-middleware-setup)
4. [Database Security](#database-security)
5. [Network Security](#network-security)
6. [Monitoring & Logging](#monitoring--logging)
7. [Incident Response](#incident-response)
8. [Post-Deployment Verification](#post-deployment-verification)

---

## Pre-Deployment Security Checklist

### Critical (Must Complete Before Production)

- [ ] **JWT Secret Key** - Generate strong, unique secret (256-bit minimum)
- [ ] **Database Password** - Use strong, unique password (not default)
- [ ] **HTTPS/TLS** - SSL certificates installed and configured
- [ ] **CSRF Protection** - Enabled on all state-changing endpoints
- [ ] **Rate Limiting** - Configured for all endpoints
- [ ] **Security Headers** - All headers properly configured
- [ ] **CORS Policy** - Whitelist production domains only
- [ ] **Input Validation** - All endpoints validate input
- [ ] **SQL Injection Protection** - Verified with security tests
- [ ] **XSS Protection** - JSON encoding and security headers enabled

### High Priority (Complete Within First Week)

- [ ] **Login Rate Limiting** - 5 attempts per 5 minutes enabled
- [ ] **JWT Claims Validation** - Strict validation enabled
- [ ] **Security Logging** - Failed login attempts, rate limit hits
- [ ] **Backup Strategy** - Daily automated backups configured
- [ ] **Disaster Recovery** - Recovery procedures documented
- [ ] **Security Monitoring** - Alerts configured for anomalies

### Medium Priority (Complete Within First Month)

- [ ] **Penetration Testing** - Third-party security audit
- [ ] **Security Training** - Team trained on security practices
- [ ] **Incident Response Plan** - Documented and tested
- [ ] **Data Encryption** - Sensitive data encrypted at rest
- [ ] **Access Control Review** - Quarterly review scheduled

---

## Environment Configuration

### 1. Environment Variables (.env)

**NEVER commit .env files to version control!**

```bash
# Production .env template

# Application
APP_ENV=production
APP_PORT=8080
APP_NAME=Cooperative-ERP-Lite

# Database
DB_HOST=<your-cloud-sql-instance>
DB_PORT=5432
DB_NAME=koperasi_erp_prod
DB_USER=<secure-username>
DB_PASSWORD=<STRONG-GENERATED-PASSWORD>
DB_SSLMODE=require

# JWT Configuration
JWT_SECRET=<256-bit-random-hex-string>  # Generate with: openssl rand -hex 32
JWT_EXPIRATION_HOURS=24

# CORS Configuration
CORS_ALLOWED_ORIGINS=https://app.yourcooperative.com,https://www.yourcooperative.com
CORS_ALLOW_CREDENTIALS=true

# Rate Limiting
RATE_LIMIT_REQUESTS_PER_MINUTE=100
LOGIN_RATE_LIMIT_ATTEMPTS=5
LOGIN_RATE_LIMIT_WINDOW_MINUTES=5
LOGIN_LOCKOUT_DURATION_MINUTES=15

# Security
CSRF_ENABLED=true
SECURE_COOKIES=true  # Requires HTTPS

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# Cloud Storage (GCP)
GCS_BUCKET_NAME=<your-bucket-name>
GCS_PROJECT_ID=<your-gcp-project-id>
```

### 2. Generate Secure Secrets

```bash
# Generate JWT secret (256-bit)
openssl rand -hex 32

# Generate database password (strong)
openssl rand -base64 32

# Generate CSRF secret
openssl rand -base64 32
```

### 3. Cloud Run Configuration

```yaml
# cloudbuild.yaml
steps:
  - name: 'gcr.io/cloud-builders/docker'
    args: ['build', '-t', 'gcr.io/$PROJECT_ID/cooperative-erp:$SHORT_SHA', '.']

  - name: 'gcr.io/cloud-builders/docker'
    args: ['push', 'gcr.io/$PROJECT_ID/cooperative-erp:$SHORT_SHA']

  - name: 'gcr.io/google.com/cloudsdktool/cloud-sdk'
    entrypoint: gcloud
    args:
      - 'run'
      - 'deploy'
      - 'cooperative-erp'
      - '--image=gcr.io/$PROJECT_ID/cooperative-erp:$SHORT_SHA'
      - '--platform=managed'
      - '--region=asia-southeast2'
      - '--allow-unauthenticated'
      - '--set-env-vars=APP_ENV=production'
      - '--set-secrets=JWT_SECRET=jwt-secret:latest,DB_PASSWORD=db-password:latest'
      - '--min-instances=1'
      - '--max-instances=10'
      - '--cpu=1'
      - '--memory=512Mi'
      - '--timeout=60s'
```

---

## Security Middleware Setup

### 1. Main Application (cmd/api/main.go)

```go
package main

import (
	"cooperative-erp-lite/internal/middleware"
	"cooperative-erp-lite/internal/routes"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Set Gin mode
	if os.Getenv("APP_ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// --- GLOBAL MIDDLEWARE (ORDER MATTERS!) ---

	// 1. Logger (first, to log everything)
	router.Use(middleware.LoggerMiddleware())

	// 2. Recovery (catch panics)
	router.Use(gin.Recovery())

	// 3. Security Headers (apply to all routes)
	router.Use(middleware.SecurityHeaders())

	// 4. CORS (before authentication)
	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	router.Use(middleware.CORSMiddleware(allowedOrigins))

	// 5. Global Rate Limiting (100 req/min per IP)
	requestsPerMinute, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_REQUESTS_PER_MINUTE"))
	if requestsPerMinute == 0 {
		requestsPerMinute = 100
	}
	router.Use(middleware.RateLimitMiddleware(requestsPerMinute, 1*time.Minute))

	// --- PUBLIC ENDPOINTS ---

	// CSRF token generation (public)
	router.GET("/api/v1/csrf-token", middleware.GenerateCSRFTokenEndpoint)

	// Health check (public, no rate limiting)
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// --- AUTH ENDPOINTS (with login rate limiting) ---

	authGroup := router.Group("/api/v1/auth")
	{
		// Login rate limiting (5 attempts per 5 min, 15 min lockout)
		loginAttempts, _ := strconv.Atoi(os.Getenv("LOGIN_RATE_LIMIT_ATTEMPTS"))
		if loginAttempts == 0 {
			loginAttempts = 5
		}

		authGroup.Use(middleware.LoginRateLimitMiddleware(
			loginAttempts,
			5*time.Minute,
			15*time.Minute,
		))

		// Setup auth routes
		routes.SetupAuthRoutes(authGroup)
	}

	// --- PROTECTED ENDPOINTS ---

	apiGroup := router.Group("/api/v1")
	{
		// Authentication required
		jwtSecret := os.Getenv("JWT_SECRET")
		expirationHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
		if expirationHours == 0 {
			expirationHours = 24
		}

		apiGroup.Use(middleware.AuthMiddleware(utils.NewJWTUtil(jwtSecret, expirationHours)))

		// CSRF Protection (POST, PUT, DELETE)
		if os.Getenv("CSRF_ENABLED") == "true" {
			apiGroup.Use(middleware.CSRFProtection())
		}

		// Setup protected routes
		routes.SetupProtectedRoutes(apiGroup)
	}

	// Start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
```

### 2. Security Headers Middleware

```go
// middleware/security_headers.go
package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders adds security headers to all responses
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Prevent MIME type sniffing
		c.Header("X-Content-Type-Options", "nosniff")

		// Prevent clickjacking
		c.Header("X-Frame-Options", "DENY")

		// Enable XSS protection
		c.Header("X-XSS-Protection", "1; mode=block")

		// Content Security Policy
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:; font-src 'self'; connect-src 'self'; frame-ancestors 'none'")

		// HTTPS only (in production)
		if c.Request.Host != "localhost" && c.Request.Host != "127.0.0.1" {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")
		}

		// Remove server information
		c.Header("Server", "")

		// Referrer Policy
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// Permissions Policy
		c.Header("Permissions-Policy", "geolocation=(), microphone=(), camera=()")

		c.Next()
	}
}
```

### 3. Enable HTTPS Cookies

```go
// middleware/csrf.go - Update cookie settings for production

func GenerateCSRFTokenEndpoint(c *gin.Context) {
	token, err := GenerateCSRFToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate CSRF token",
		})
		return
	}

	// Set as cookie for browser clients
	secure := os.Getenv("SECURE_COOKIES") == "true"  // true in production

	c.SetCookie(
		CSRFTokenCookie,
		token,
		86400, // 24 hours
		"/",
		"",
		secure,  // Only send over HTTPS in production
		true,    // HttpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
```

---

## Database Security

### 1. PostgreSQL Configuration

```sql
-- Create production database user with limited privileges
CREATE USER cooperative_app WITH PASSWORD '<strong-password>';

-- Grant only necessary privileges
GRANT CONNECT ON DATABASE koperasi_erp_prod TO cooperative_app;
GRANT USAGE ON SCHEMA public TO cooperative_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO cooperative_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO cooperative_app;

-- Prevent schema modifications by application user
REVOKE CREATE ON SCHEMA public FROM cooperative_app;

-- Enable SSL/TLS connections only
ALTER SYSTEM SET ssl = on;
ALTER SYSTEM SET ssl_cert_file = '/path/to/server.crt';
ALTER SYSTEM SET ssl_key_file = '/path/to/server.key';

-- Reject non-SSL connections
hostssl  koperasi_erp_prod  cooperative_app  0.0.0.0/0  md5
```

### 2. Connection Pool Settings

```go
// database/connection.go

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),  // "require" in production
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,  // Use prepared statements for all queries
	})

	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
```

### 3. Backup Strategy

```bash
# Daily backup script (run via cron)
#!/bin/bash

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/backups/cooperative-erp"
DB_NAME="koperasi_erp_prod"

# Create backup
pg_dump -h $DB_HOST -U $DB_USER -d $DB_NAME | gzip > $BACKUP_DIR/backup_$DATE.sql.gz

# Upload to Cloud Storage
gsutil cp $BACKUP_DIR/backup_$DATE.sql.gz gs://cooperative-erp-backups/

# Keep only last 30 days locally
find $BACKUP_DIR -name "backup_*.sql.gz" -mtime +30 -delete

# Verify backup
gunzip -t $BACKUP_DIR/backup_$DATE.sql.gz
if [ $? -eq 0 ]; then
    echo "Backup successful: backup_$DATE.sql.gz"
else
    echo "ERROR: Backup verification failed!" | mail -s "Backup Failed" admin@cooperative.com
fi
```

---

## Network Security

### 1. Cloud Run Network Settings

- **VPC Connector**: Use VPC connector for private database access
- **Firewall Rules**: Allow HTTPS (443) only, block HTTP (80)
- **IP Whitelist**: Restrict admin endpoints to known IPs

### 2. Cloud SQL Security

```yaml
# Cloud SQL instance configuration
settings:
  ipConfiguration:
    requireSsl: true
    authorizedNetworks: []
    privateNetwork: projects/<PROJECT_ID>/global/networks/<VPC_NAME>

  backupConfiguration:
    enabled: true
    startTime: "03:00"  # Daily at 3 AM
    pointInTimeRecoveryEnabled: true

  databaseFlags:
    - name: log_connections
      value: "on"
    - name: log_disconnections
      value: "on"
    - name: log_checkpoints
      value: "on"
```

---

## Monitoring & Logging

### 1. Security Event Logging

```go
// middleware/security_logger.go

type SecurityEvent struct {
	Timestamp    time.Time
	EventType    string
	IP           string
	UserAgent    string
	Username     string
	Success      bool
	Details      string
}

func LogSecurityEvent(eventType, ip, userAgent, username string, success bool, details string) {
	event := SecurityEvent{
		Timestamp:    time.Now(),
		EventType:    eventType,
		IP:           ip,
		UserAgent:    userAgent,
		Username:     username,
		Success:      success,
		Details:      details,
	}

	// Log to structured logging system (Cloud Logging, ELK, etc.)
	log.Printf("[SECURITY] %+v", event)

	// Alert on critical events
	if eventType == "BRUTE_FORCE" || eventType == "SQL_INJECTION_ATTEMPT" {
		sendSecurityAlert(event)
	}
}
```

### 2. Metrics to Monitor

- **Authentication**:
  - Failed login attempts per IP
  - Failed login attempts per username
  - Successful logins
  - Account lockouts

- **Rate Limiting**:
  - Rate limit hits per endpoint
  - Rate limit hits per IP
  - Top IPs hitting rate limits

- **API Performance**:
  - Request latency (p50, p95, p99)
  - Error rates (4xx, 5xx)
  - Requests per second

- **Database**:
  - Connection pool usage
  - Query performance
  - Slow query count

### 3. Alerting Rules

```yaml
# Example Prometheus alerting rules
groups:
  - name: security
    interval: 1m
    rules:
      - alert: HighFailedLoginRate
        expr: rate(failed_login_attempts[5m]) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High failed login rate detected"

      - alert: RateLimitExceeded
        expr: rate(rate_limit_hits[5m]) > 100
        for: 2m
        labels:
          severity: warning
        annotations:
          summary: "Unusually high rate limiting activity"

      - alert: DatabaseConnectionPoolExhausted
        expr: db_connection_pool_usage > 0.9
        for: 1m
        labels:
          severity: critical
        annotations:
          summary: "Database connection pool nearly exhausted"
```

---

## Incident Response

### 1. Security Incident Playbook

#### Suspected SQL Injection Attack

1. **Detect**: Alert triggered for SQL errors or suspicious query patterns
2. **Investigate**:
   - Check application logs for the suspicious request
   - Identify the affected endpoint and user
   - Review database logs for any unauthorized data access
3. **Contain**:
   - Block the attacker's IP address
   - Temporarily disable the affected endpoint if necessary
4. **Remediate**:
   - Review and fix the vulnerable code
   - Deploy the fix
   - Run security tests to verify
5. **Review**:
   - Document the incident
   - Update security tests to catch similar issues

#### Account Compromise

1. **Detect**: Unusual activity or user reports unauthorized access
2. **Investigate**:
   - Review login history for the account
   - Check for data access or modifications
   - Identify login locations and IP addresses
3. **Contain**:
   - Immediately invalidate all tokens for the user
   - Lock the account
   - Reset password
4. **Remediate**:
   - Contact the user to verify account ownership
   - Review and restore any modified data
   - Enable 2FA for the account (future enhancement)
5. **Review**:
   - Investigate how the compromise occurred
   - Update security practices if needed

#### DDoS Attack

1. **Detect**: Sudden spike in traffic, rate limiting triggered excessively
2. **Investigate**:
   - Identify attack pattern (IP ranges, user agents, endpoints)
   - Determine if it's distributed or from a few sources
3. **Contain**:
   - Enable Cloud Armor (GCP) or similar DDoS protection
   - Temporarily reduce rate limits
   - Block attacking IP ranges
4. **Remediate**:
   - Scale up resources if legitimate traffic is affected
   - Maintain blocking rules
5. **Review**:
   - Analyze attack patterns
   - Update rate limiting configuration

### 2. Incident Communication

- **Internal**: Notify security team, development team, management
- **External**: Notify affected users (if data breach), regulatory bodies (if required)
- **Timeline**: Document when incident was detected, contained, resolved

---

## Post-Deployment Verification

### 1. Security Test Suite

Run the complete security test suite after deployment:

```bash
# SSH into production-like environment
ssh production-server

# Run security tests
cd /app/backend
go test ./internal/tests/security/... -v

# Expected: All tests should pass
# If any tests fail, investigate and fix before going live
```

### 2. Manual Security Verification

#### Test HTTPS

```bash
# Verify SSL certificate is valid
curl -vI https://api.yourcooperative.com

# Check for security headers
curl -I https://api.yourcooperative.com/health
# Should see: X-Content-Type-Options, X-Frame-Options, etc.
```

#### Test CSRF Protection

```bash
# Without CSRF token (should fail)
curl -X POST https://api.yourcooperative.com/api/v1/anggota \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <valid-token>" \
  -d '{"namaLengkap": "Test"}'
# Expected: 403 Forbidden

# With CSRF token (should succeed)
CSRF_TOKEN=$(curl -s https://api.yourcooperative.com/api/v1/csrf-token | jq -r .token)
curl -X POST https://api.yourcooperative.com/api/v1/anggota \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <valid-token>" \
  -H "X-CSRF-Token: $CSRF_TOKEN" \
  -d '{"namaLengkap": "Test"}'
# Expected: Success
```

#### Test Rate Limiting

```bash
# Trigger rate limit
for i in {1..150}; do
  curl -s https://api.yourcooperative.com/api/v1/anggota \
    -H "Authorization: Bearer <valid-token>"
done
# Expected: 429 Too Many Requests after 100 requests
```

#### Test Login Rate Limiting

```bash
# Trigger login rate limit
for i in {1..6}; do
  curl -X POST https://api.yourcooperative.com/api/v1/auth/login \
    -H "Content-Type: application/json" \
    -d '{"namaPengguna": "test", "kataSandi": "wrong"}'
done
# Expected: 429 Too Many Requests after 5 failed attempts
```

### 3. Penetration Testing

Schedule a professional penetration test:

- **Timing**: Within first month of production
- **Scope**: All API endpoints, authentication, data access
- **Focus Areas**:
  - Authentication bypass
  - Authorization bypass (horizontal/vertical privilege escalation)
  - SQL injection
  - XSS attacks
  - CSRF attacks
  - Business logic flaws
  - Data exposure

---

## Security Maintenance Schedule

### Daily

- Review security logs for anomalies
- Monitor rate limiting metrics
- Check backup success

### Weekly

- Review failed login attempts
- Analyze blocked IPs
- Update security blocklists if needed

### Monthly

- Review and rotate API keys/secrets
- Update dependencies (security patches)
- Review user access and permissions
- Analyze security metrics trends

### Quarterly

- Conduct security training for team
- Review and update incident response plan
- Perform security audit
- Review third-party dependencies for vulnerabilities

### Annually

- Renew SSL certificates
- Conduct comprehensive penetration test
- Review and update security policies
- Evaluate new security threats and mitigations

---

## Emergency Contacts

```
Security Team Lead: security@yourcooperative.com
DevOps Team: devops@yourcooperative.com
On-Call Engineer: +62-XXX-XXXX-XXXX
Management: management@yourcooperative.com

External Security Consultant: [If applicable]
Hosting Provider Support: [GCP Support]
```

---

## Compliance & Regulations

### Indonesian Data Protection

- **UU ITE No. 19/2016**: Electronic Information and Transactions
- **PP No. 71/2019**: Implementation of Electronic Systems and Transactions
- **OJK Regulations**: Financial Services Authority requirements (if applicable)

### Data Residency

- Store data in Indonesian data centers (asia-southeast2)
- Ensure compliance with local data sovereignty laws

### Audit Trail

- Maintain logs for minimum 12 months
- Record all data modifications with user attribution
- Enable PostgreSQL audit logging

---

## Revision History

| Version | Date | Changes | Author |
|---------|------|---------|--------|
| 1.0.0   | 2025-11-19 | Initial production deployment guide | Claude (Security Testing) |

---

## Approval Sign-off

- [ ] **Security Team Approved**: _____________________ Date: _________
- [ ] **DevOps Team Approved**: _____________________ Date: _________
- [ ] **Management Approved**: _____________________ Date: _________

---

**This document is confidential and should only be shared with authorized personnel.**
