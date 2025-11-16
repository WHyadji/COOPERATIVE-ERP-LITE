# Supabase Schema Reviewer - Claude Code Slash Command

## Overview

The `/supabase-schema-reviewer` command is a comprehensive schema analysis tool that identifies errors, warnings, and bad practices in Supabase database schemas. It leverages deep Supabase platform knowledge to provide actionable recommendations for optimization, security, and best practices.

## Command Definition

```json
{
  "name": "supabase-schema-reviewer",
  "description": "Comprehensive Supabase schema reviewer for error detection and best practice validation",
  "version": "1.0.0",
  "prompt": "You are a Supabase Schema Security and Performance Reviewer with expert knowledge of Supabase platform specifics, PostgreSQL optimization, and database security best practices. Your role is to analyze database schemas, SQL files, and configurations to identify:\n\n## CRITICAL ERROR DETECTION:\n\n### 1. Security Vulnerabilities\n- Missing or inefficient RLS policies\n- Insecure authentication patterns\n- Exposed sensitive data without encryption\n- Weak JWT claim validations\n- Missing audit trails for sensitive operations\n- Improper user role configurations\n- Vulnerable file storage permissions\n- SQL injection risks in dynamic queries\n- Missing rate limiting considerations\n\n### 2. Performance Issues\n- Missing critical indexes\n- Inefficient RLS policy implementations\n- Suboptimal real-time subscription configurations\n- Poor query patterns for PostgREST API\n- Inefficient data types and constraints\n- Missing foreign key relationships\n- Unoptimized pagination strategies\n- Memory-intensive operations\n- Poor partitioning strategies\n\n### 3. Supabase Platform Violations\n- Incorrect usage of Supabase extensions\n- Improper real-time publication setup\n- Inefficient auth.users integration\n- Storage bucket misconfiguration\n- Edge function integration issues\n- Dashboard optimization oversights\n- API endpoint inefficiencies\n- Webhook configuration problems\n\n### 4. Schema Design Anti-patterns\n- Poor normalization strategies\n- Circular dependencies\n- Missing or excessive constraints\n- Improper data modeling\n- Inefficient relationship designs\n- Poor naming conventions\n- Missing documentation\n- Inadequate error handling\n\n## WARNING CATEGORIES:\n\n### 1. Performance Warnings\n- Potential bottlenecks in high-traffic scenarios\n- Suboptimal index strategies\n- Real-time subscription overhead\n- API response size concerns\n- Connection pool limitations\n- Memory usage optimization opportunities\n\n### 2. Security Warnings\n- RLS policy gaps\n- Potential privilege escalation\n- Data exposure risks\n- Audit trail gaps\n- Encryption recommendations\n- Access control improvements\n\n### 3. Maintainability Warnings\n- Complex query patterns\n- Poor schema documentation\n- Hard-to-maintain triggers\n- Complex RLS policies\n- Migration complexity\n- Testing challenges\n\n### 4. Cost Optimization Warnings\n- Storage inefficiencies\n- Bandwidth optimization opportunities\n- Compute resource waste\n- Unnecessary real-time subscriptions\n- Oversized data types\n\n## SUPABASE-SPECIFIC ANALYSIS:\n\n### Real-time Subscriptions Review\n- Analyze real-time publication efficiency\n- Check for unnecessary real-time tables\n- Validate subscription filter patterns\n- Review payload size optimization\n- Check for real-time security policies\n\n### RLS Policy Analysis\n- Validate policy completeness\n- Check for policy performance issues\n- Analyze policy complexity\n- Review multi-tenant isolation\n- Check for policy conflicts\n\n### Authentication Integration Review\n- Analyze auth.users integration patterns\n- Check JWT claim usage\n- Review user metadata strategies\n- Validate social auth configurations\n- Check session management\n\n### Storage Integration Analysis\n- Review storage bucket configurations\n- Check file permission patterns\n- Analyze CDN optimization\n- Review file metadata schemas\n- Check storage security policies\n\n### PostgREST API Optimization\n- Analyze API query efficiency\n- Check for N+1 query problems\n- Review response optimization\n- Validate pagination strategies\n- Check bulk operation patterns\n\n## OUTPUT FORMAT:\n\nProvide comprehensive analysis in this structure:\n\n### 1. EXECUTIVE SUMMARY\n- Critical issues count\n- Security risk level (HIGH/MEDIUM/LOW)\n- Performance impact assessment\n- Compliance status\n- Overall schema health score (0-100)\n\n### 2. CRITICAL ERRORS (🚨)\n- List all critical security and performance issues\n- Provide immediate action items\n- Include code fixes where applicable\n- Estimate impact severity\n\n### 3. WARNINGS (⚠️)\n- Performance optimization opportunities\n- Security improvements\n- Best practice violations\n- Future scalability concerns\n\n### 4. RECOMMENDATIONS (💡)\n- Specific improvement suggestions\n- Alternative approaches\n- Performance optimization strategies\n- Security enhancements\n\n### 5. SUPABASE OPTIMIZATION (🚀)\n- Platform-specific improvements\n- Feature utilization opportunities\n- Integration enhancements\n- Dashboard configuration suggestions\n\n### 6. ACTION PLAN (📋)\n- Prioritized task list\n- Implementation timeline\n- Risk mitigation strategies\n- Testing recommendations\n\n### 7. MONITORING SETUP (📊)\n- Performance monitoring queries\n- Security alerting configurations\n- Health check implementations\n- Maintenance procedures\n\n## ANALYSIS DEPTH LEVELS:\n\n### Quick Scan (--depth=quick)\n- Critical security issues\n- Major performance problems\n- Basic Supabase compliance\n\n### Standard Review (--depth=standard)\n- Comprehensive error detection\n- Performance optimization\n- Security best practices\n- Supabase feature utilization\n\n### Deep Analysis (--depth=deep)\n- Advanced optimization strategies\n- Complex security patterns\n- Scalability planning\n- Cost optimization\n- Future-proofing recommendations\n\n### Enterprise Audit (--depth=enterprise)\n- Compliance validation\n- Advanced security patterns\n- Multi-environment analysis\n- Disaster recovery planning\n- Advanced monitoring setup\n\nWhen analyzing schemas, always:\n1. Consider the specific Supabase context and limitations\n2. Provide actionable, specific recommendations\n3. Include code examples for fixes\n4. Prioritize security and performance issues\n5. Consider multi-tenant implications\n6. Validate real-time subscription efficiency\n7. Check RLS policy completeness and performance\n8. Analyze API optimization opportunities\n9. Review storage and file management patterns\n10. Provide comprehensive monitoring recommendations\n\nRemember: Every recommendation should be Supabase-specific and consider the platform's unique features, limitations, and best practices.",
  "parameters": {
    "schema_file": {
      "description": "Path to schema file(s) to analyze (SQL, migration files, or directory)",
      "type": "string",
      "required": false
    },
    "tables": {
      "description": "Specific tables to focus on (comma-separated)",
      "type": "string",
      "required": false
    },
    "depth": {
      "description": "Analysis depth (quick, standard, deep, enterprise)",
      "type": "string",
      "required": false,
      "default": "standard"
    },
    "focus": {
      "description": "Analysis focus (security, performance, rls, realtime, api, storage, all)",
      "type": "string",
      "required": false,
      "default": "all"
    },
    "environment": {
      "description": "Target environment (development, staging, production)",
      "type": "string",
      "required": false,
      "default": "production"
    },
    "app_type": {
      "description": "Application type for context (saas, ecommerce, social, analytics, chat)",
      "type": "string",
      "required": false
    },
    "compliance": {
      "description": "Compliance requirements (gdpr, hipaa, sox, pci, ccpa)",
      "type": "string",
      "required": false
    },
    "output_format": {
      "description": "Output format (detailed, summary, json, checklist)",
      "type": "string",
      "required": false,
      "default": "detailed"
    },
    "include_fixes": {
      "description": "Include code fixes in output",
      "type": "boolean",
      "required": false,
      "default": true
    },
    "severity_filter": {
      "description": "Minimum severity to report (critical, high, medium, low, info)",
      "type": "string",
      "required": false,
      "default": "medium"
    }
  },
  "examples": [
    {
      "command": "/supabase-schema-reviewer --schema_file=supabase/migrations/ --depth=deep --focus=security",
      "description": "Deep security analysis of migration files"
    },
    {
      "command": "/supabase-schema-reviewer --tables=users,orders --focus=rls --environment=production",
      "description": "RLS policy review for specific tables in production"
    },
    {
      "command": "/supabase-schema-reviewer --depth=enterprise --compliance=gdpr --app_type=saas",
      "description": "Enterprise-level GDPR compliance audit for SaaS application"
    },
    {
      "command": "/supabase-schema-reviewer --focus=performance --output_format=checklist",
      "description": "Performance-focused review with checklist output"
    }
  ]
}
```

## Core Analysis Categories

### 🚨 Critical Errors

#### Security Vulnerabilities
- **Missing RLS Policies**: Tables without row-level security enabled
- **Insecure RLS Patterns**: Policies that allow unauthorized access
- **Exposed Sensitive Data**: PII without encryption or proper access controls
- **JWT Vulnerabilities**: Improper JWT claim validation or usage
- **Storage Misconfigurations**: Public buckets with sensitive data
- **SQL Injection Risks**: Dynamic query construction vulnerabilities

#### Performance Critical Issues
- **Missing Primary Indexes**: Tables without proper primary key indexing
- **Unindexed Foreign Keys**: Foreign key relationships without supporting indexes
- **Inefficient RLS Policies**: Security policies causing table scans
- **Real-time Overload**: Excessive real-time subscriptions
- **API Bottlenecks**: Queries causing N+1 problems in PostgREST

### ⚠️ Warnings

#### Performance Warnings
- **Suboptimal Data Types**: Using VARCHAR instead of TEXT where appropriate
- **Missing Composite Indexes**: Complex queries without supporting indexes
- **Inefficient Triggers**: Heavy operations in trigger functions
- **Poor Pagination**: Missing or inefficient pagination strategies

#### Security Warnings
- **Incomplete Audit Trails**: Missing change tracking for sensitive operations
- **Weak Authentication Patterns**: Suboptimal user management strategies
- **Insufficient Data Classification**: Sensitive data without proper marking

### 💡 Recommendations

#### Supabase Optimization
- **Real-time Efficiency**: Optimize subscription patterns and payloads
- **RLS Performance**: Index optimization for security policies
- **API Optimization**: PostgREST query pattern improvements
- **Storage Optimization**: CDN and file management enhancements

## Usage Examples

### Comprehensive Schema Review

```bash
/supabase-schema-reviewer --schema_file=supabase/migrations/ --depth=deep --focus=all --environment=production
```

**Expected Output:**
```
# 🔍 Supabase Schema Review Report
Generated: 2025-01-15 14:30:00 UTC
Environment: Production
Analysis Depth: Deep

## 📊 EXECUTIVE SUMMARY
- Schema Health Score: 72/100
- Critical Issues: 3
- Security Risk Level: MEDIUM
- Performance Impact: HIGH
- Compliance Status: ⚠️ Partial

## 🚨 CRITICAL ERRORS (3)

### 1. Missing RLS Policy - HIGH SEVERITY
**Table:** `public.user_payments`
**Issue:** Table contains sensitive payment data without RLS protection
**Impact:** Potential data exposure to unauthorized users
**Fix:**
```sql
ALTER TABLE public.user_payments ENABLE ROW LEVEL SECURITY;

CREATE POLICY "Users can only access their payment data" 
ON public.user_payments FOR ALL 
USING (user_id = auth.uid());
```

### 2. Inefficient RLS Policy - HIGH SEVERITY
**Table:** `public.documents`
**Issue:** RLS policy causes full table scan on every query
**Current Policy:**
```sql
-- PROBLEMATIC: No supporting index
CREATE POLICY "org_access" ON public.documents 
FOR SELECT USING (
    organization_id IN (
        SELECT org_id FROM user_orgs WHERE user_id = auth.uid()
    )
);
```
**Fix:**
```sql
-- Add supporting index
CREATE INDEX idx_documents_org_user 
ON public.documents(organization_id, user_id);

-- Optimize policy
CREATE OR REPLACE POLICY "org_access" ON public.documents 
FOR SELECT USING (
    EXISTS (
        SELECT 1 FROM user_orgs uo 
        WHERE uo.org_id = documents.organization_id 
        AND uo.user_id = auth.uid()
    )
);
```

### 3. Real-time Performance Issue - MEDIUM SEVERITY
**Issue:** All tables published to real-time without filtering
**Impact:** Excessive bandwidth usage and client-side processing
**Fix:**
```sql
-- Remove unnecessary publications
ALTER PUBLICATION supabase_realtime DROP TABLE public.audit_logs;
ALTER PUBLICATION supabase_realtime DROP TABLE public.system_config;

-- Add filtered publications for relevant tables
ALTER PUBLICATION supabase_realtime ADD TABLE public.messages 
WHERE (created_at > NOW() - INTERVAL '1 day');
```

## ⚠️ WARNINGS (8)

### Performance Warnings

#### 1. Missing Composite Index
**Table:** `public.orders`
**Query Pattern:** Frequent filtering by user_id AND status
**Recommendation:**
```sql
CREATE INDEX idx_orders_user_status 
ON public.orders(user_id, status, created_at DESC);
```

#### 2. Inefficient Pagination
**API Endpoint:** `/api/products`
**Issue:** Using OFFSET for pagination on large dataset
**Recommendation:**
```sql
-- Add cursor-based pagination support
ALTER TABLE public.products ADD COLUMN cursor_id BIGSERIAL;
CREATE INDEX idx_products_cursor ON public.products(cursor_id);
```

### Security Warnings

#### 3. Incomplete Audit Trail
**Tables:** `public.user_profiles`, `public.organizations`
**Issue:** No change tracking for sensitive operations
**Recommendation:**
```sql
-- Add audit trigger
CREATE OR REPLACE FUNCTION audit_changes()
RETURNS TRIGGER AS $$
BEGIN
    INSERT INTO audit_log (
        table_name, operation, old_data, new_data, user_id
    ) VALUES (
        TG_TABLE_NAME, TG_OP,
        CASE WHEN TG_OP != 'INSERT' THEN to_json(OLD) END,
        CASE WHEN TG_OP != 'DELETE' THEN to_json(NEW) END,
        auth.uid()
    );
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER audit_user_profiles
    AFTER INSERT OR UPDATE OR DELETE ON public.user_profiles
    FOR EACH ROW EXECUTE FUNCTION audit_changes();
```

## 💡 RECOMMENDATIONS

### 1. Implement Comprehensive RLS Strategy
- Enable RLS on all user-data tables
- Create policy templates for common patterns
- Add performance indexes for all RLS policies
- Implement policy testing framework

### 2. Optimize Real-time Subscriptions
- Review all real-time publications
- Implement subscription filtering
- Add real-time payload optimization
- Monitor subscription performance

### 3. Enhance API Performance
- Add composite indexes for common query patterns
- Implement efficient pagination strategies
- Optimize JOIN operations for PostgREST
- Add API response caching where appropriate

## 🚀 SUPABASE OPTIMIZATION

### Real-time Enhancements
```sql
-- Optimize message subscriptions with filtering
CREATE OR REPLACE VIEW public.recent_messages AS
SELECT id, channel_id, user_id, content, created_at
FROM public.messages 
WHERE created_at > NOW() - INTERVAL '24 hours';

ALTER PUBLICATION supabase_realtime ADD TABLE public.recent_messages;
```

### Storage Integration
```sql
-- Add file metadata tracking
CREATE TABLE public.file_metadata (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    storage_path TEXT NOT NULL,
    original_name TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    size_bytes BIGINT NOT NULL,
    uploaded_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- RLS for file access
ALTER TABLE public.file_metadata ENABLE ROW LEVEL SECURITY;
CREATE POLICY "Users can access their files" ON public.file_metadata
FOR ALL USING (uploaded_by = auth.uid());
```

## 📋 ACTION PLAN

### Immediate Actions (Next 1-2 days)
1. ✅ Fix missing RLS policies on sensitive tables
2. ✅ Add critical performance indexes
3. ✅ Optimize inefficient RLS policies
4. ✅ Review and filter real-time publications

### Short Term (Next 1-2 weeks)
1. 🔄 Implement comprehensive audit logging
2. 🔄 Add API performance optimizations
3. 🔄 Enhance storage security patterns
4. 🔄 Set up performance monitoring

### Medium Term (Next month)
1. 📅 Implement advanced security patterns
2. 📅 Add comprehensive testing framework
3. 📅 Optimize for scalability
4. 📅 Enhance documentation

### Long Term (Next quarter)
1. 🎯 Advanced compliance implementation
2. 🎯 Multi-region optimization
3. 🎯 Advanced analytics and monitoring
4. 🎯 Disaster recovery planning

## 📊 MONITORING SETUP

### Performance Monitoring
```sql
-- Query performance tracking
CREATE OR REPLACE VIEW public.slow_queries AS
SELECT 
    query,
    mean_exec_time,
    calls,
    total_exec_time,
    rows,
    100.0 * shared_blks_hit / nullif(shared_blks_hit + shared_blks_read, 0) AS hit_percent
FROM pg_stat_statements 
WHERE mean_exec_time > 100  -- Queries slower than 100ms
ORDER BY mean_exec_time DESC;
```

### Security Monitoring
```sql
-- RLS policy performance monitoring
CREATE OR REPLACE FUNCTION monitor_rls_performance()
RETURNS TABLE(
    table_name TEXT,
    policy_name TEXT,
    avg_execution_time INTERVAL
) AS $$
BEGIN
    -- Implementation for RLS performance monitoring
    RETURN QUERY
    SELECT 
        schemaname::TEXT,
        policyname::TEXT,
        AVG(execution_time)::INTERVAL
    FROM pg_stat_statements_rls
    GROUP BY schemaname, policyname
    HAVING AVG(execution_time) > INTERVAL '50ms';
END;
$$ LANGUAGE plpgsql;
```

### Health Checks
```sql
-- Database health monitoring
CREATE OR REPLACE FUNCTION public.database_health_check()
RETURNS JSON AS $$
DECLARE
    result JSON;
BEGIN
    SELECT json_build_object(
        'timestamp', NOW(),
        'connection_count', (SELECT count(*) FROM pg_stat_activity),
        'slow_queries', (SELECT count(*) FROM public.slow_queries),
        'table_sizes', (
            SELECT json_object_agg(
                schemaname||'.'||tablename,
                pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename))
            )
            FROM pg_tables 
            WHERE schemaname = 'public'
        ),
        'rls_enabled_tables', (
            SELECT json_agg(tablename)
            FROM pg_tables pt
            JOIN pg_class pc ON pc.relname = pt.tablename
            WHERE pt.schemaname = 'public' AND pc.relrowsecurity = true
        )
    ) INTO result;
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;
```
```

### Security-Focused Review

```bash
/supabase-schema-reviewer --focus=security --depth=deep --compliance=gdpr --include_fixes=true
```

### RLS Policy Analysis

```bash
/supabase-schema-reviewer --focus=rls --tables=users,documents,payments --environment=production
```

### Performance Analysis

```bash
/supabase-schema-reviewer --focus=performance --depth=standard --output_format=checklist
```

### Real-time Optimization Review

```bash
/supabase-schema-reviewer --focus=realtime --app_type=chat --depth=deep
```

## Specialized Analysis Features

### 1. Multi-Tenant Security Review
```bash
/supabase-schema-reviewer --focus=security --app_type=saas --depth=enterprise --compliance=sox
```

### 2. E-commerce Platform Analysis
```bash
/supabase-schema-reviewer --app_type=ecommerce --focus=all --compliance=pci --depth=deep
```

### 3. Chat Application Optimization
```bash
/supabase-schema-reviewer --app_type=chat --focus=realtime --depth=deep --tables=messages,channels,presence
```

### 4. Analytics Platform Review
```bash
/supabase-schema-reviewer --app_type=analytics --focus=performance --depth=enterprise
```

## Installation and Setup

### 1. Create Command File
```bash
# Create the command file
touch ~/.claude-code/commands/supabase-schema-reviewer.json
```

### 2. Add Command Configuration
Copy the complete JSON configuration to the file.

### 3. Verify Installation
```bash
# Check if command is available
claude-code --list-commands | grep supabase-schema-reviewer

# Test the command
/supabase-schema-reviewer --help
```

## Integration Examples

### CI/CD Pipeline Integration

```yaml
# .github/workflows/schema-review.yml
name: Supabase Schema Review
on:
  pull_request:
    paths:
      - 'supabase/migrations/**'
      - 'database/**'

jobs:
  schema-review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Security Review
        run: |
          claude-code /supabase-schema-reviewer \
            --schema_file=supabase/migrations/ \
            --focus=security \
            --depth=deep \
            --severity_filter=high \
            --output_format=json > security-report.json
      
      - name: Performance Review
        run: |
          claude-code /supabase-schema-reviewer \
            --schema_file=supabase/migrations/ \
            --focus=performance \
            --depth=standard \
            --output_format=checklist > performance-checklist.md
      
      - name: Upload Reports
        uses: actions/upload-artifact@v3
        with:
          name: schema-review-reports
          path: |
            security-report.json
            performance-checklist.md
```

### Pre-commit Hook Integration

```bash
#!/bin/sh
# .git/hooks/pre-commit

# Run schema review on staged SQL files
STAGED_SQL_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep -E "\.(sql|migration)$")

if [ ! -z "$STAGED_SQL_FILES" ]; then
    echo "🔍 Running Supabase schema review..."
    
    claude-code /supabase-schema-reviewer \
        --schema_file="$STAGED_SQL_FILES" \
        --focus=security \
        --depth=quick \
        --severity_filter=critical \
        --output_format=summary
    
    if [ $? -ne 0 ]; then
        echo "❌ Schema review found critical issues. Commit aborted."
        exit 1
    fi
    
    echo "✅ Schema review passed!"
fi
```

### Development Workflow Integration

```bash
# Add to package.json scripts
{
  "scripts": {
    "db:review": "claude-code /supabase-schema-reviewer --focus=all --depth=standard",
    "db:security": "claude-code /supabase-schema-reviewer --focus=security --depth=deep",
    "db:performance": "claude-code /supabase-schema-reviewer --focus=performance --depth=standard",
    "db:rls": "claude-code /supabase-schema-reviewer --focus=rls --depth=deep",
    "db:compliance": "claude-code /supabase-schema-reviewer --compliance=gdpr --depth=enterprise"
  }
}
```

## Advanced Features

### Custom Rule Configuration

```json
// .supabase-reviewer.config.json
{
  "rules": {
    "security": {
      "enforce_rls": true,
      "require_audit_logs": ["users", "payments", "sensitive_data"],
      "encryption_required": ["ssn", "payment_info", "personal_data"]
    },
    "performance": {
      "max_table_size_without_partitioning": "100GB",
      "require_indexes_for_foreign_keys": true,
      "max_rls_policy_complexity": 3
    },
    "naming": {
      "table_naming_pattern": "^[a-z][a-z0-9_]*[a-z0-9]$",
      "column_naming_pattern": "^[a-z][a-z0-9_]*[a-z0-9]$",
      "index_naming_pattern": "^idx_[a-z0-9_]+$"
    }
  },
  "compliance": {
    "gdpr": {
      "require_data_classification": true,
      "require_deletion_policies": true,
      "require_consent_tracking": true
    }
  }
}
```

### Custom Analysis Plugins

```bash
# Custom security analysis
/supabase-schema-reviewer --config=.supabase-reviewer.config.json --focus=security --custom-rules=security-extra.json
```

## Output Formats

### Detailed Report (Default)
Comprehensive analysis with code examples and fixes.

### Summary Report
```bash
/supabase-schema-reviewer --output_format=summary
```

### JSON Output (for CI/CD)
```bash
/supabase-schema-reviewer --output_format=json --severity_filter=high
```

### Checklist Format
```bash
/supabase-schema-reviewer --output_format=checklist --focus=security
```
