---
allowed-tools: ReadFile, WriteFile, CreateFile
description: Create a new API endpoint with proper validation, error handling, and authentication
---

## Context

- API routes: !`find frontend/src/app/api -name "route.ts" | head -10`
- Middleware: !`ls frontend/src/middleware/`
- Validation schemas: @frontend/src/lib/validation/schemas.ts

## Your task

Create a new API endpoint: $ARGUMENTS

Implement:
1. Route handler with proper HTTP methods
2. Request validation using Zod
3. Authentication and authorization checks
4. Rate limiting middleware
5. Error handling and logging
6. Database operations with Supabase
7. Response formatting
8. TypeScript types for request/response

Include:
- Input validation and sanitization
- Proper HTTP status codes
- CORS configuration if needed
- Request logging for debugging
- Performance monitoring
- Cache headers if applicable
- API documentation comments
- Integration with existing services
- Proper error messages (no sensitive data) 