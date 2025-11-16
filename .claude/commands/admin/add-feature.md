---
allowed-tools: ReadFile, WriteFile, SearchReplace, CreateFile
description: Add a new feature to the admin panel with proper authentication and permissions
---

## Context

- Admin routes: !`find frontend/src/app/\\(admin\\)/admin -type f -name "*.tsx" | head -20`
- Admin components: !`ls frontend/src/components/admin/`
- Current admin features: @frontend/admin-panel-implementation-guide.md

## Your task

Add a new admin panel feature: $ARGUMENTS

Help me implement:
1. Create the admin route and page component
2. Add navigation menu item
3. Implement proper authentication checks
4. Set up role-based permissions
5. Create necessary API endpoints
6. Add data tables/forms as needed
7. Implement proper loading and error states
8. Add activity logging

Ensure:
- Consistent UI with existing admin panel
- Proper middleware authentication
- Rate limiting on sensitive operations
- Audit trail for admin actions
- Responsive design
- Accessibility compliance 