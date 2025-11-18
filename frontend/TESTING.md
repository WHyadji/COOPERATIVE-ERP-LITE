# Testing Guide

This project uses **Vitest** as the testing framework for frontend tests.

## Quick Start

```bash
# Run all tests
npm test

# Run tests in watch mode (auto-rerun on file changes)
npm run test:watch

# Run tests with UI (browser-based test interface)
npm run test:ui

# Run tests with coverage report
npm run test:coverage
```

## Testing Stack

- **Vitest 4.0** - Fast, modern test runner (Jest-compatible API)
- **@testing-library/react** - React component testing utilities
- **@testing-library/jest-dom** - Custom matchers for DOM assertions
- **@testing-library/user-event** - User interaction simulation
- **jsdom** - Browser environment for tests

## Project Structure

```
frontend/
├── __tests__/              # Test files
│   ├── components/         # Component tests
│   ├── lib/               # Utility/API tests
│   └── example.test.ts    # Example basic tests
├── vitest.config.ts       # Vitest configuration
├── vitest.setup.ts        # Global test setup
└── TESTING.md             # This file
```

## Writing Tests

### Basic Test Example

```typescript
import { describe, it, expect } from 'vitest';

describe('Calculator', () => {
  it('should add two numbers', () => {
    expect(1 + 1).toBe(2);
  });
});
```

### React Component Test Example

```typescript
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import { LoginPage } from '@/app/(auth)/login/page';

describe('LoginPage', () => {
  it('should render login form', () => {
    render(<LoginPage />);

    const emailInput = screen.getByLabelText(/email/i);
    const passwordInput = screen.getByLabelText(/password/i);
    const submitButton = screen.getByRole('button', { name: /login/i });

    expect(emailInput).toBeInTheDocument();
    expect(passwordInput).toBeInTheDocument();
    expect(submitButton).toBeInTheDocument();
  });
});
```

### User Interaction Test Example

```typescript
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';

it('should handle form submission', async () => {
  const user = userEvent.setup();
  const mockSubmit = vi.fn();

  render(<LoginForm onSubmit={mockSubmit} />);

  await user.type(screen.getByLabelText(/email/i), 'test@example.com');
  await user.type(screen.getByLabelText(/password/i), 'password123');
  await user.click(screen.getByRole('button', { name: /login/i }));

  expect(mockSubmit).toHaveBeenCalledWith({
    email: 'test@example.com',
    password: 'password123',
  });
});
```

### API/Utility Test Example

```typescript
import { describe, it, expect, vi } from 'vitest';
import axios from 'axios';
import { memberApi } from '@/lib/api/memberApi';

// Mock axios
vi.mock('axios');

describe('memberApi', () => {
  it('should fetch members successfully', async () => {
    const mockMembers = [
      { id: '1', namaLengkap: 'John Doe' },
      { id: '2', namaLengkap: 'Jane Doe' },
    ];

    vi.mocked(axios.get).mockResolvedValueOnce({
      data: { data: mockMembers },
    });

    const result = await memberApi.getMembers({ page: 1, pageSize: 10 });

    expect(result.data).toEqual(mockMembers);
    expect(axios.get).toHaveBeenCalledWith('/anggota', {
      params: { page: 1, pageSize: 10 },
    });
  });
});
```

## Test Organization

### Naming Conventions

- Test files: `*.test.ts` or `*.test.tsx`
- Place tests in `__tests__/` directory
- Mirror source file structure in tests
- Example:
  ```
  app/(auth)/login/page.tsx
  __tests__/app/(auth)/login/page.test.tsx
  ```

### Test Structure

Use **Arrange-Act-Assert** pattern:

```typescript
it('should update user profile', async () => {
  // Arrange - Set up test data and mocks
  const mockUser = { id: '1', name: 'John' };
  const mockUpdate = vi.fn();

  // Act - Perform the action
  const result = await updateProfile(mockUser.id, { name: 'Jane' });

  // Assert - Verify the outcome
  expect(result.name).toBe('Jane');
  expect(mockUpdate).toHaveBeenCalledTimes(1);
});
```

## Mocking

### Mock Next.js Router

```typescript
// Already configured in vitest.setup.ts
import { useRouter } from 'next/navigation';

it('should navigate on click', () => {
  const router = useRouter();
  // router.push is already mocked
});
```

### Mock API Calls

```typescript
import { vi } from 'vitest';
import axios from 'axios';

vi.mock('axios');

it('should call API', async () => {
  vi.mocked(axios.get).mockResolvedValue({ data: 'success' });
  // Test code
});
```

### Mock Components

```typescript
vi.mock('@/components/Header', () => ({
  Header: () => <div>Mocked Header</div>,
}));
```

## Coverage Configuration

Coverage is configured in `vitest.config.ts`:

```typescript
coverage: {
  provider: 'v8',
  reporter: ['text', 'json', 'html', 'lcov'],
  exclude: [
    'node_modules/',
    'vitest.config.ts',
    '**/*.config.{js,ts}',
    '**/types/**',
  ],
}
```

View coverage report:
```bash
npm run test:coverage
# Opens ./coverage/index.html
```

## Best Practices

1. **Test User Behavior, Not Implementation**
   - ❌ Test internal state changes
   - ✅ Test what user sees and interacts with

2. **Keep Tests Isolated**
   - Each test should run independently
   - Clean up after tests (automatic with afterEach cleanup)

3. **Use Descriptive Test Names**
   - ❌ `it('works', ...)`
   - ✅ `it('should display error message when login fails', ...)`

4. **Test Edge Cases**
   - Empty inputs
   - Error states
   - Loading states
   - Large datasets

5. **Mock External Dependencies**
   - API calls
   - Browser APIs
   - Third-party libraries

## Common Matchers

```typescript
// Equality
expect(value).toBe(expected);           // Strict equality (===)
expect(value).toEqual(expected);        // Deep equality

// Truthiness
expect(value).toBeTruthy();
expect(value).toBeFalsy();
expect(value).toBeNull();
expect(value).toBeUndefined();

// Numbers
expect(value).toBeGreaterThan(3);
expect(value).toBeLessThan(5);

// Strings
expect(string).toMatch(/pattern/);
expect(string).toContain('substring');

// Arrays
expect(array).toContain(item);
expect(array).toHaveLength(3);

// DOM (jest-dom matchers)
expect(element).toBeInTheDocument();
expect(element).toBeVisible();
expect(element).toBeDisabled();
expect(element).toHaveTextContent('text');
expect(element).toHaveAttribute('href', '/link');
```

## Debugging Tests

### Run Single Test File

```bash
npm test __tests__/example.test.ts
```

### Run Tests Matching Pattern

```bash
npm test -- -t "login"  # Runs tests with "login" in name
```

### Use Vitest UI

```bash
npm run test:ui
```

Opens browser with interactive test explorer.

### Debug in VS Code

Add to `.vscode/launch.json`:

```json
{
  "type": "node",
  "request": "launch",
  "name": "Debug Vitest Tests",
  "runtimeExecutable": "npm",
  "runtimeArgs": ["run", "test:watch"],
  "console": "integratedTerminal"
}
```

## CI/CD Integration

Tests run automatically in GitHub Actions:

```yaml
- name: Run tests
  run: npm test -- --coverage
```

## TODO: Tests to Implement

As per PR review findings, these tests are **CRITICAL** before production:

### Priority 1 (Critical - Before Merge)
- [ ] AuthContext login/logout flows
- [ ] API client 401 handling
- [ ] Token refresh mechanism
- [ ] Protected route authorization

### Priority 2 (High - Sprint 1)
- [ ] Member CRUD operations
- [ ] Form validation (Zod schemas)
- [ ] API error handling
- [ ] Multi-tenant data validation

### Priority 3 (Medium - Sprint 2)
- [ ] Component rendering tests
- [ ] User interaction flows
- [ ] Edge case scenarios
- [ ] Performance tests

**Target Coverage**: 50% for MVP, 80% for production

## Resources

- [Vitest Documentation](https://vitest.dev/)
- [Testing Library Docs](https://testing-library.com/docs/react-testing-library/intro/)
- [jest-dom Matchers](https://github.com/testing-library/jest-dom)

---

**Status**: ✅ Vitest setup complete (backbone only)
**Next Step**: Implement actual tests for critical code paths
**CI Status**: Will pass once example tests removed and real tests added
