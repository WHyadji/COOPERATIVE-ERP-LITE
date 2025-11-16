---
allowed-tools: Bash(npm:*), Bash(pnpm:*), Bash(yarn:*), Bash(supabase:*), ReadFile, WriteFile
description: Setup development environment with all dependencies and services
---

## Context

- Current directory: !`pwd`
- Node version: !`node --version`
- Package manager: !`which pnpm || which yarn || which npm`
- Supabase status: !`supabase status 2>/dev/null || echo "Supabase CLI not running"`

## Your task

Help me set up the development environment for this project:

1. Check if all required dependencies are installed
2. Install any missing dependencies
3. Set up environment variables if needed
4. Start the development server
5. Provide clear instructions for any manual steps required

Consider:
- Database setup (Supabase)
- Redis/Upstash configuration
- Environment variables
- Any required API keys or services

$ARGUMENTS 