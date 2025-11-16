# Supabase Schema Master - Claude Code Slash Command

## Overview

The `/supabase-schema-master` command is a specialized database schema optimization tool designed specifically for Supabase applications. It leverages Supabase's unique features including real-time subscriptions, Row Level Security (RLS), built-in authentication, and PostgREST API optimization. use supabaseMCP to access the database, and do that in there.

## Command Definition

```json
{
  "name": "supabase-schema-master",
  "description": "Supabase-optimized schema mastering tool for maximum performance and feature utilization",
  "version": "1.0.0"
}
```

## Core Expertise Areas

### 1. Real-time Subscriptions Optimization
- **Schema Design**: Tables optimized for real-time subscriptions
- **Payload Optimization**: Minimize data transfer for real-time updates
- **Conflict Resolution**: Design conflict-free replicated data types (CRDTs)
- **Publication Management**: Efficient publication/subscription patterns
- **Performance**: Optimize for high-frequency real-time updates

### 2. Row Level Security (RLS) Design
- **Policy Creation**: Comprehensive RLS policy implementation
- **Performance Optimization**: Index-backed policies for fast evaluation
- **Multi-tenant Strategies**: Efficient tenant isolation patterns
- **User Context**: Optimal user context and JWT claim utilization
- **Security Balance**: Balance security with query performance

### 3. Supabase Auth Integration
- **User Management**: Optimal integration with `auth.users` table
- **Profile Patterns**: Efficient user profile and metadata schemas
- **Role-based Access**: Comprehensive RBAC implementations
- **Social Authentication**: OAuth provider integration patterns
- **Session Management**: Optimized session and token handling

### 4. PostgREST API Optimization
- **Query Efficiency**: Schemas optimized for PostgREST auto-generated APIs
- **Response Optimization**: Computed columns for API responses
- **Relationship Design**: Efficient JOIN patterns for API queries
- **Pagination**: High-performance pagination strategies
- **Bulk Operations**: Optimized bulk insert/update patterns

### 5. Storage & File Management
- **Metadata Schemas**: Efficient file metadata organization
- **Bucket Strategy**: Storage bucket organization patterns
- **Permission Systems**: File-level permission implementations
- **CDN Integration**: Content delivery optimization
- **Media Processing**: Image and video optimization workflows

### 6. Supabase Extensions Utilization
- **pg_cron**: Scheduled task implementations
- **pg_net**: HTTP request handling patterns
- **pgjwt**: JWT token management
- **uuid-ossp**: UUID optimization strategies
- **pg_trgm**: Full-text search implementations

## Command Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `app` | string | No | - | Application type (saas, social, ecommerce, collaborative, chat, analytics) |
| `tables` | string | No | - | Comma-separated list of table names to focus on |
| `optimize` | string | No | balanced | Optimization target (realtime, api, auth, storage, analytics, balanced) |
| `rls` | boolean | No | true | Enable Row Level Security optimization |
| `realtime` | boolean | No | false | Enable real-time subscription optimization |
| `multitenant` | string | No | - | Multi-tenant strategy (row-level, schema-based, hybrid) |
| `auth` | string | No | supabase | Authentication strategy (supabase, custom, social) |
| `action` | string | No | create-schema | Action to perform (create-schema, analyze-performance, optimize-rls, setup-realtime, migrate-with-rls) |
| `environment` | string | No | development | Target environment (development, staging, production) |

## Usage Examples

### Basic Schema Creation

#### Simple SaaS Application
```bash
/supabase-schema-master create-schema --app=saas --tables=organizations,users,projects
```

#### E-commerce with Storage
```bash
/supabase-schema-master create-schema --app=ecommerce --tables=products,orders,reviews --optimize=storage --rls=true
```

### Real-time Applications

#### Chat Application
```bash
/supabase-schema-master create-realtime-schema --app=chat --tables=messages,channels,users --rls=true --realtime=true
```

#### Collaborative Editor
```bash
/supabase-schema-master create-schema --app=collaborative --tables=documents,comments,cursors --realtime=true --optimize=realtime
```

### Multi-tenant Applications

#### Row-level Multi-tenancy
```bash
/supabase-schema-master create-multitenant --isolation=row-level --auth=supabase --app=saas
```

#### Schema-based Multi-tenancy
```bash
/supabase-schema-master create-multitenant --isolation=schema-based --tables=core_tables --environment=production
```

### Performance Analysis

#### Comprehensive Analysis
```bash
/supabase-schema-master analyze-performance --include=rls-policies,realtime,api-performance
```

#### RLS-focused Analysis
```bash
/supabase-schema-master analyze-performance --include=rls-policies --tables=sensitive_data
```

### Migration Operations

#### Migration with RLS
```bash
/supabase-schema-master migrate-with-rls --from-schema=current.sql --tenant-strategy=shared
```

#### Zero-downtime Migration
```bash
/supabase-schema-master migrate-with-rls --from-schema=legacy.sql --requirements=zero-downtime --environment=production
```

## Output Structure

### Standard Schema Output

```sql
-- ================================================
-- Supabase Optimized Schema
-- Generated by: Claude Code Supabase Schema Master
-- Target: High-Performance Supabase Application
-- Date: [Generated Date]
-- ================================================

-- 1. EXTENSIONS & SETUP
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- 2. AUTHENTICATION INTEGRATION
-- Leverage auth.users for user management

-- 3. CORE TABLES WITH RLS
CREATE TABLE public.profiles (
    id UUID REFERENCES auth.users(id) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    avatar_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

ALTER TABLE public.profiles ENABLE ROW LEVEL SECURITY;

-- 4. OPTIMIZED INDEXES
CREATE INDEX CONCURRENTLY idx_profiles_username_gin 
    ON public.profiles USING gin(username gin_trgm_ops);

-- 5. RLS POLICIES
CREATE POLICY "Users can view their own profile" ON public.profiles
    FOR SELECT USING (auth.uid() = id);

CREATE POLICY "Users can update their own profile" ON public.profiles
    FOR UPDATE USING (auth.uid() = id);

-- 6. REALTIME CONFIGURATION
ALTER PUBLICATION supabase_realtime ADD TABLE public.profiles;

-- 7. API OPTIMIZATION
-- Computed columns and views for API efficiency

-- 8. TRIGGERS & FUNCTIONS
CREATE OR REPLACE FUNCTION public.handle_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 9. MONITORING SETUP
-- Performance monitoring and alerting
```

## Specialized Schema Patterns

### Real-time Chat Schema

```bash
/supabase-schema-master create-realtime-schema --app=chat --tables=messages,channels,users,presence --realtime=true
```

**Generated Output:**
```sql
-- Real-time optimized chat schema
CREATE TABLE public.channels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE public.messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id UUID REFERENCES public.channels(id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.users(id),
    content TEXT NOT NULL,
    message_type VARCHAR(20) DEFAULT 'text',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE public.user_presence (
    user_id UUID REFERENCES auth.users(id) PRIMARY KEY,
    status VARCHAR(20) DEFAULT 'offline',
    last_seen TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Real-time subscriptions
ALTER PUBLICATION supabase_realtime ADD TABLE public.channels;
ALTER PUBLICATION supabase_realtime ADD TABLE public.messages;
ALTER PUBLICATION supabase_realtime ADD TABLE public.user_presence;

-- Optimized indexes for real-time queries
CREATE INDEX idx_messages_channel_created ON public.messages(channel_id, created_at DESC);
CREATE INDEX idx_messages_user_created ON public.messages(user_id, created_at DESC);
```

### Multi-tenant SaaS Schema

```bash
/supabase-schema-master create-multitenant --isolation=row-level --app=saas --tables=organizations,projects,tasks
```

**Generated Output:**
```sql
-- Multi-tenant SaaS schema with row-level security
CREATE TABLE public.organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(50) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE public.organization_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES public.organizations(id) ON DELETE CASCADE,
    user_id UUID REFERENCES auth.users(id) ON DELETE CASCADE,
    role VARCHAR(20) DEFAULT 'member',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(organization_id, user_id)
);

CREATE TABLE public.projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    organization_id UUID REFERENCES public.organizations(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    created_by UUID REFERENCES auth.users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- RLS Policies for multi-tenancy
ALTER TABLE public.organizations ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.organization_members ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.projects ENABLE ROW LEVEL SECURITY;

-- Organization access policies
CREATE POLICY "Users can view organizations they belong to" ON public.organizations
    FOR SELECT USING (
        id IN (
            SELECT organization_id 
            FROM public.organization_members 
            WHERE user_id = auth.uid()
        )
    );

-- Project access policies
CREATE POLICY "Users can view projects in their organizations" ON public.projects
    FOR SELECT USING (
        organization_id IN (
            SELECT organization_id 
            FROM public.organization_members 
            WHERE user_id = auth.uid()
        )
    );
```

### E-commerce with Storage Integration

```bash
/supabase-schema-master create-schema --app=ecommerce --optimize=storage --tables=products,orders,reviews
```

**Generated Output:**
```sql
-- E-commerce schema optimized for storage
CREATE TABLE public.product_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    parent_id UUID REFERENCES public.product_categories(id)
);

CREATE TABLE public.products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(200) NOT NULL,
    description TEXT,
    price DECIMAL(10,2) NOT NULL,
    category_id UUID REFERENCES public.product_categories(id),
    image_urls TEXT[], -- Array of storage URLs
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE public.product_images (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID REFERENCES public.products(id) ON DELETE CASCADE,
    storage_path TEXT NOT NULL, -- Supabase storage path
    alt_text VARCHAR(200),
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Storage bucket policies (to be applied in Supabase dashboard)
/*
Bucket: product-images
Public: false
File size limit: 5MB
Allowed MIME types: image/jpeg, image/png, image/webp

RLS Policies for storage:
- Users can upload images for products they own
- All users can view product images
*/

-- Optimized indexes for e-commerce queries
CREATE INDEX idx_products_category_price ON public.products(category_id, price);
CREATE INDEX idx_products_name_gin ON public.products USING gin(name gin_trgm_ops);
CREATE INDEX idx_product_images_product_sort ON public.product_images(product_id, sort_order);
```

## Advanced Optimization Strategies

### Real-time Performance Optimization

```sql
-- Minimize real-time payload sizes
CREATE VIEW public.message_notifications AS
SELECT 
    id,
    channel_id,
    user_id,
    LEFT(content, 100) as preview, -- Truncated content
    created_at
FROM public.messages;

-- Real-time filters to reduce bandwidth
ALTER PUBLICATION supabase_realtime ADD TABLE public.messages 
    WHERE (created_at > NOW() - INTERVAL '1 hour');
```

### RLS Policy Performance

```sql
-- Optimized RLS with proper indexing
CREATE INDEX idx_organization_members_user_org 
    ON public.organization_members(user_id, organization_id);

-- Security definer function for complex policies
CREATE OR REPLACE FUNCTION public.user_organizations(user_uuid UUID)
RETURNS TABLE(organization_id UUID)
LANGUAGE sql
SECURITY DEFINER
STABLE
AS $$
    SELECT om.organization_id
    FROM public.organization_members om
    WHERE om.user_id = user_uuid;
$$;

-- Use function in RLS policy for better performance
CREATE POLICY "Efficient org access" ON public.projects
    FOR SELECT USING (
        organization_id IN (SELECT public.user_organizations(auth.uid()))
    );
```

### API Performance Optimization

```sql
-- Computed columns for API responses
ALTER TABLE public.products 
ADD COLUMN search_vector tsvector 
GENERATED ALWAYS AS (
    to_tsvector('english', name || ' ' || COALESCE(description, ''))
) STORED;

CREATE INDEX idx_products_search ON public.products USING gin(search_vector);

-- Materialized view for complex aggregations
CREATE MATERIALIZED VIEW public.product_stats AS
SELECT 
    p.id,
    p.name,
    p.price,
    COUNT(r.id) as review_count,
    AVG(r.rating) as avg_rating,
    pc.name as category_name
FROM public.products p
LEFT JOIN public.reviews r ON p.id = r.product_id
LEFT JOIN public.product_categories pc ON p.category_id = pc.id
GROUP BY p.id, p.name, p.price, pc.name;

CREATE UNIQUE INDEX idx_product_stats_id ON public.product_stats(id);

-- Refresh strategy
CREATE OR REPLACE FUNCTION public.refresh_product_stats()
RETURNS void AS $$
BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY public.product_stats;
END;
$$ LANGUAGE plpgsql;
```

## Monitoring and Analytics

### Performance Monitoring Setup

```sql
-- Query performance tracking
CREATE TABLE public.query_performance_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    query_type VARCHAR(50),
    execution_time INTERVAL,
    table_name VARCHAR(100),
    user_id UUID REFERENCES auth.users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- RLS policy performance monitoring
CREATE OR REPLACE FUNCTION public.log_slow_policies()
RETURNS trigger AS $$
BEGIN
    -- Log queries that take longer than 100ms
    IF (EXTRACT(EPOCH FROM clock_timestamp() - statement_timestamp()) > 0.1) THEN
        INSERT INTO public.performance_alerts (
            alert_type, 
            message, 
            metadata
        ) VALUES (
            'slow_rls_policy',
            'Slow RLS policy detected',
            jsonb_build_object(
                'table', TG_TABLE_NAME,
                'operation', TG_OP,
                'duration_ms', EXTRACT(EPOCH FROM clock_timestamp() - statement_timestamp()) * 1000
            )
        );
    END IF;
    
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;
```

### Real-time Monitoring

```sql
-- Real-time subscription monitoring
CREATE TABLE public.realtime_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    table_name VARCHAR(100),
    operation VARCHAR(10),
    subscriber_count INTEGER,
    payload_size INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Monitor real-time performance
CREATE OR REPLACE FUNCTION public.track_realtime_performance()
RETURNS trigger AS $$
BEGIN
    INSERT INTO public.realtime_metrics (
        table_name,
        operation,
        payload_size
    ) VALUES (
        TG_TABLE_NAME,
        TG_OP,
        octet_length(row_to_json(COALESCE(NEW, OLD))::text)
    );
    
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;
```

## Security Best Practices

### Comprehensive RLS Implementation

```sql
-- Audit trail for sensitive operations
CREATE TABLE public.audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    table_name VARCHAR(100),
    operation VARCHAR(10),
    old_values JSONB,
    new_values JSONB,
    user_id UUID REFERENCES auth.users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Automatic audit logging
CREATE OR REPLACE FUNCTION public.audit_trigger()
RETURNS trigger AS $$
BEGIN
    INSERT INTO public.audit_log (
        table_name,
        operation,
        old_values,
        new_values,
        user_id
    ) VALUES (
        TG_TABLE_NAME,
        TG_OP,
        CASE WHEN TG_OP IN ('UPDATE', 'DELETE') THEN row_to_json(OLD) END,
        CASE WHEN TG_OP IN ('INSERT', 'UPDATE') THEN row_to_json(NEW) END,
        auth.uid()
    );
    
    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Apply audit trigger to sensitive tables
CREATE TRIGGER audit_sensitive_table
    AFTER INSERT OR UPDATE OR DELETE ON public.sensitive_data
    FOR EACH ROW EXECUTE FUNCTION public.audit_trigger();
```

### Data Encryption Patterns

```sql
-- Encrypted sensitive data
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE public.user_sensitive_data (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES auth.users(id),
    encrypted_ssn TEXT, -- Encrypted with application key
    encrypted_payment_info TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- RLS with encryption
CREATE POLICY "Users can only access their encrypted data" ON public.user_sensitive_data
    FOR ALL USING (user_id = auth.uid());
```

## Installation and Setup

### 1. Create Command File

```bash
# Create the command file
touch ~/.claude-code/commands/supabase-schema-master.json
```

### 2. Add Command Configuration

Copy the complete JSON configuration to the file:

```json
{
  "name": "supabase-schema-master",
  "description": "Supabase-optimized schema mastering tool for maximum performance and feature utilization",
  "version": "1.0.0",
  "prompt": "[Complete prompt from earlier...]",
  "parameters": {
    "app": {
      "description": "Application type (e.g., saas, social, ecommerce, collaborative, chat, analytics)",
      "type": "string",
      "required": false
    },
    "tables": {
      "description": "Comma-separated list of table names to focus on",
      "type": "string",
      "required": false
    },
    "optimize": {
      "description": "Optimization target (realtime, api, auth, storage, analytics, balanced)",
      "type": "string",
      "required": false,
      "default": "balanced"
    },
    "rls": {
      "description": "Enable Row Level Security optimization",
      "type": "boolean",
      "required": false,
      "default": true
    },
    "realtime": {
      "description": "Enable real-time subscription optimization",
      "type": "boolean",
      "required": false,
      "default": false
    },
    "multitenant": {
      "description": "Multi-tenant strategy (row-level, schema-based, hybrid)",
      "type": "string",
      "required": false
    },
    "auth": {
      "description": "Authentication strategy (supabase, custom, social)",
      "type": "string",
      "required": false,
      "default": "supabase"
    },
    "action": {
      "description": "Action to perform (create-schema, analyze-performance, optimize-rls, setup-realtime, migrate-with-rls)",
      "type": "string",
      "required": false,
      "default": "create-schema"
    },
    "environment": {
      "description": "Target environment (development, staging, production)",
      "type": "string",
      "required": false,
      "default": "development"
    }
  },
  "examples": [
    {
      "command": "/supabase-schema-master create-schema --app=saas --multitenant=row-level --rls=true",
      "description": "Create a multi-tenant SaaS schema with RLS"
    },
    {
      "command": "/supabase-schema-master setup-realtime --tables=messages,presence --optimize=realtime",
      "description": "Set up real-time optimized schema for chat application"
    },
    {
      "command": "/supabase-schema-master analyze-performance --include=rls-policies,api-performance",
      "description": "Analyze schema performance including RLS and API optimization"
    },
    {
      "command": "/supabase-schema-master optimize-rls --tables=documents --multitenant=row-level",
      "description": "Optimize RLS policies for multi-tenant document management"
    }
  ]
}
```

### 3. Verify Installation

```bash
# Check if command is available
claude-code --list-commands | grep supabase-schema-master

# Test the command
/supabase-schema-master --help
```

## Integration with Development Workflow

### CI/CD Pipeline Integration

```yaml
# .github/workflows/supabase-schema.yml
name: Supabase Schema Optimization
on:
  pull_request:
    paths:
      - 'supabase/migrations/**'

jobs:
  optimize-schema:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Analyze Schema Performance
        run: |
          claude-code /supabase-schema-master analyze-performance \
            --schema=supabase/migrations/latest.sql \
            --include=rls-policies,realtime,api-performance
      
      - name: Generate Optimized Schema
        run: |
          claude-code /supabase-schema-master create-schema \
            --app=${{ github.event.repository.name }} \
            --environment=production \
            --rls=true
```

### Local Development Setup

```bash
# Add to package.json scripts
{
  "scripts": {
    "db:optimize": "claude-code /supabase-schema-master create-schema --app=myapp --rls=true",
    "db:analyze": "claude-code /supabase-schema-master analyze-performance --include=all",
    "db:migrate": "claude-code /supabase-schema-master migrate-with-rls --from-schema=current.sql",
    "db:realtime": "claude-code /supabase-schema-master setup-realtime --tables=all"
  }
}
```

## Troubleshooting

### Common Issues

#### 1. RLS Policy Performance
```bash
# Analyze slow policies
/supabase-schema-master analyze-performance --include=rls-policies --verbose=true

# Optimize specific policies
/supabase-schema-master optimize-rls --tables=slow_table --performance=high
```

#### 2. Real-time Subscription Issues
```bash
# Debug real-time setup
/supabase-schema-master analyze-performance --include=realtime --debug=true

# Optimize real-time configuration
/supabase-schema-master setup-realtime --tables=problematic_table --optimize=bandwidth
```

#### 3. API Performance Problems
```bash
# Analyze API performance
/supabase-schema-master analyze-performance --include=api-performance

# Create API-optimized views
/supabase-schema-master optimize-api --tables=slow_api_tables
```

## Best Practices Summary

1. **Always Enable RLS**: Use `--rls=true` for production applications
2. **Optimize for Use Case**: Choose appropriate `--optimize` parameter
3. **Multi-tenant Strategy**: Plan tenant isolation early with `--multitenant`
4. **Real-time Considerations**: Only enable real-time for tables that need it
5. **Environment-Specific**: Use different configurations for dev/staging/prod
6. **Regular Analysis**: Run performance analysis regularly
7. **Monitor Metrics**: Set up comprehensive monitoring from day one
8. **Security First**: Implement comprehensive audit trails and encryption
