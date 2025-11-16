---
allowed-tools: ReadFile, WriteFile, SearchReplace, CreateFile
description: Implement analytics tracking for user behavior and business metrics
---

## Context

- Analytics setup: @frontend/src/lib/analytics.ts
- Current tracking: !`grep -r "track\\|analytics" frontend/src/components --include="*.tsx" | head -10 || echo "No tracking found"`
- Analytics components: !`ls frontend/src/components/analytics/`

## Your task

Implement analytics tracking for: $ARGUMENTS

Set up:
1. Event tracking for user interactions
2. Page view tracking
3. Conversion tracking
4. Performance metrics
5. Error tracking
6. Custom business metrics
7. A/B testing support
8. User journey mapping

Include:
- Privacy-compliant tracking (GDPR/CCPA)
- Anonymous user tracking
- Session recording considerations
- Real-time dashboards
- Data retention policies
- Export capabilities
- Integration with Google Analytics/Mixpanel
- Custom event properties
- Funnel analysis setup 