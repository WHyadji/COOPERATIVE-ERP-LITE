---
allowed-tools: ReadFile, WriteFile, Bash(supabase:*), SearchReplace
description: Create and manage Row Level Security policies for Supabase tables
---

## Context

- use supabase MCP
- Auth setup: @frontend/src/utils/supabase/server.ts

## Your task

Create or update RLS policies for: $ARGUMENTS

Help me implement:
1. Analyze the table structure and requirements
2. Create appropriate RLS policies
3. Test policy effectiveness
4. Handle different user roles
5. Ensure data isolation
6. Optimize policy performance

Consider these policy patterns:
- User can only see their own data
- Admin can see all data
- Public read, authenticated write
- Team-based access control
- Soft delete visibility
- Time-based access restrictions

Ensure:
- No security loopholes
- Performance optimization (use indexes)
- Proper error messages
- Policy documentation 