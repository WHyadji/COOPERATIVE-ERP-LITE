---
allowed-tools: ReadFile, WriteFile, SearchReplace, Bash(npm:run)
description: Optimize code for performance, bundle size, and best practices
---

## Context

- Bundle analysis: !`cd frontend && npm run build 2>&1 | grep -E "(First Load|Route)" | head -20 || echo "No build output"`
- TypeScript errors: !`cd frontend && npm run type-check 2>&1 | head -20 || echo "No type errors"`
- Current performance: !`grep -r "use.*Query\\|use.*Mutation" frontend/src --include="*.ts*" | wc -l || echo "0"` React Query hooks

# Code Optimization: $ARGUMENTS

Review and optimize the code at: @$ARGUMENTS

## Review Areas
1. **Code Quality**
   - Check for code smells and anti-patterns
   - Review naming conventions
   - Assess code readability and maintainability
   - Identify overly complex functions

2. **Performance Analysis**
   - Look for performance bottlenecks
   - Check for inefficient algorithms
   - Review memory usage patterns
   - Identify unnecessary computations

3. **Security Review**
   - Check for security vulnerabilities
   - Review input validation
   - Assess authentication/authorization
   - Look for potential injection attacks

4. **Best Practices**
   - Ensure proper error handling
   - Check for appropriate logging
   - Review configuration management
   - Assess testing coverage

5. **Optimization Implementation**
   - Refactor inefficient code
   - Improve algorithm complexity where possible
   - Optimize database queries
   - Reduce memory footprint

6. **Testing**
   - Ensure optimizations don't break functionality
   - Add performance tests if beneficial
   - Verify security improvements

Provide specific recommendations and implement the most impactful optimizations.