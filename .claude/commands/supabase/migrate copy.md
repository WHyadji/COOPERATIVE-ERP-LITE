---
allowed-tools: Bash(supabase:*), ReadFile, WriteFile, SearchReplace
description: Create and manage Supabase migrations
---

## Context

- use supabase MCP

## Your task

Help me create or manage Supabase database migrations. I need to:

$ARGUMENTS

If no specific request is provided, help me:
1. Review existing migrations
2. Create a new migration based on my requirements
3. Apply migrations to the database
4. Handle any rollback if needed

Consider:
- RLS policies
- Indexes for performance
- Proper foreign key constraints
- Audit columns (created_at, updated_at)
- Soft delete patterns if applicable 