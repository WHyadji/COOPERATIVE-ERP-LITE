---
allowed-tools: ReadFile, WriteFile, CreateFile
description: Generate a new React component with TypeScript, proper styling, and best practices
---

## Context

- UI components: !`ls frontend/src/components/ui/ | head -10`
- Component patterns: @frontend/src/components/ui/button.tsx
- Design system: @frontend/src/styles/design-system.ts

## Your task

Create a new React component: $ARGUMENTS

Generate:
1. Component file with TypeScript interfaces
2. Proper props definition and defaults
3. Accessibility attributes
4. Responsive design with Tailwind
5. Dark mode support
6. Loading and error states
7. Storybook story (if applicable)
8. Basic unit tests

Follow patterns:
- Use forwardRef for DOM element components
- Implement proper TypeScript generics
- Use cn() for className merging
- Follow naming conventions
- Export types separately
- Add JSDoc comments
- Use composition over inheritance
- Implement proper event handlers
- Consider performance (memo if needed) 