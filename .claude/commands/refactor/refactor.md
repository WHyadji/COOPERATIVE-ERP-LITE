---
description: Systematic code refactoring workflow
allowed-tools: Bash(*), Write(*), Search(*)
---

# Refactor: $ARGUMENTS

Refactor the component/module: "$ARGUMENTS"

## Refactoring Steps
1. **Current State Analysis**
   - Understand existing functionality
   - Identify current problems and pain points
   - Document current behavior
   - Ensure comprehensive test coverage exists

2. **Refactoring Strategy**
   - Define clear refactoring goals
   - Plan the new structure/approach
   - Identify potential risks
   - Create step-by-step refactoring plan

3. **Safety Measures**
   - Ensure all existing tests pass
   - Add missing tests for current behavior
   - Create backup/branch if needed
   - Plan rollback strategy

4. **Incremental Refactoring**
   - Make small, safe changes
   - Run tests after each change
   - Maintain functionality throughout process
   - Commit frequently with clear messages

5. **Improvement Implementation**
   - Improve code structure and organization
   - Eliminate code duplication
   - Enhance readability and maintainability
   - Optimize performance where appropriate

6. **Validation**
   - Ensure all tests still pass
   - Verify no regressions introduced
   - Test edge cases thoroughly
   - Update documentation as needed

Think carefully about each step and prioritize safety and maintainability.