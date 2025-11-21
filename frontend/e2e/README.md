# E2E Testing Guide

This directory contains End-to-End (E2E) tests for the Cooperative ERP system using Playwright.

## Test Files

### 1. `member-portal.spec.ts`
**Test Coverage:**
- Member portal login flow
- Dashboard functionality
- Transaction history
- Profile management
- Balance inquiry

**Status:** ✅ Complete

### 2. `pos-system.spec.ts`
**Test Coverage:**
- Product management (CRUD operations)
- POS transaction flow
- Stock management
- Payment & change calculation
- Receipt generation
- Sales history
- Role-based access control

**Status:** ✅ Complete

**Test Scenarios:** 35+ test cases covering:
- ✅ Product list display
- ✅ Create/update/delete products
- ✅ Product search and filtering
- ✅ Low stock warnings
- ✅ Add products to cart (search & barcode)
- ✅ Adjust quantities
- ✅ Calculate totals with multiple items
- ✅ Cash payment processing
- ✅ Change calculation
- ✅ Insufficient payment validation
- ✅ Customer/member selection
- ✅ Receipt generation
- ✅ Stock updates after sale
- ✅ Out of stock prevention
- ✅ Sales transaction history
- ✅ Date range filtering
- ✅ Transaction detail view
- ✅ Sales report export
- ✅ Cashier vs Admin access control

---

## Prerequisites

### 1. Install Playwright

```bash
cd frontend
npm install
npx playwright install
```

### 2. Backend Running

Ensure backend server is running on port 8080:

```bash
cd backend
go run cmd/api/main.go
```

### 3. Frontend Running

Ensure frontend is running on port 3000:

```bash
cd frontend
npm run dev
```

### 4. Test Data Seeded

Seed test data for E2E tests:

```bash
cd backend
go run cmd/seed-test-data/main.go
```

**Test Credentials Created:**
- **Admin**: `admin@koperasi.local` / `Admin123!`
- **Cashier**: `kasir@koperasi.local` / `Kasir123!`
- **Member**: `A001` / `123456`

---

## Running Tests

### Run All E2E Tests

```bash
cd frontend
npx playwright test
```

### Run Specific Test File

```bash
# Member Portal tests only
npx playwright test e2e/member-portal.spec.ts

# POS System tests only
npx playwright test e2e/pos-system.spec.ts
```

### Run in UI Mode (Interactive)

```bash
npx playwright test --ui
```

### Run in Headed Mode (See Browser)

```bash
npx playwright test --headed
```

### Run Specific Test by Name

```bash
# Run only product management tests
npx playwright test -g "Product Management"

# Run only transaction flow tests
npx playwright test -g "Transaction Flow"
```

### Debug Mode

```bash
npx playwright test --debug
```

---

## Test Reports

### View HTML Report

After tests run:

```bash
npx playwright show-report
```

### Generate Report Manually

```bash
npx playwright test --reporter=html
```

---

## Configuration

E2E tests are configured in `playwright.config.ts`:

```typescript
{
  testDir: "./e2e",
  baseURL: "https://localhost" || process.env.NEXT_PUBLIC_APP_URL,
  timeout: 30000, // 30 seconds per test
  retries: 2, // Retry failed tests twice
  use: {
    screenshot: "only-on-failure",
    video: "retain-on-failure",
    trace: "on-first-retry",
  }
}
```

---

## Writing New E2E Tests

### Test Structure Template

```typescript
import { test, expect, type Page } from "@playwright/test";

// Test data
const TEST_DATA = {
  // Your test data
};

// Helper functions
async function loginHelper(page: Page) {
  // Login logic
}

test.describe("Feature Name", () => {

  test.beforeEach(async ({ page }) => {
    // Setup before each test
    await loginHelper(page);
  });

  test("should do something", async ({ page }) => {
    // Arrange
    await page.goto("/feature");

    // Act
    await page.click('button:has-text("Action")');

    // Assert
    await expect(page.getByText("Success")).toBeVisible();
  });
});
```

### Best Practices

1. **Use Data Attributes**: Prefer `data-testid` selectors:
   ```typescript
   await page.locator('[data-testid="submit-button"]').click();
   ```

2. **Wait for Elements**: Always wait for dynamic content:
   ```typescript
   await page.waitForSelector('[data-testid="result"]');
   await expect(page.getByText("Loaded")).toBeVisible();
   ```

3. **Clean Test Data**: Reset state between tests:
   ```typescript
   test.afterEach(async ({ page }) => {
     // Cleanup
   });
   ```

4. **Use Page Objects**: Extract reusable page interactions:
   ```typescript
   class LoginPage {
     constructor(private page: Page) {}

     async login(email: string, password: string) {
       await this.page.fill('input[name="email"]', email);
       await this.page.fill('input[name="password"]', password);
       await this.page.click('button[type="submit"]');
     }
   }
   ```

---

## CI/CD Integration

### GitHub Actions Example

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '18'

      - name: Install dependencies
        run: cd frontend && npm ci

      - name: Install Playwright
        run: cd frontend && npx playwright install --with-deps

      - name: Start backend
        run: |
          cd backend
          go run cmd/api/main.go &
          sleep 5

      - name: Start frontend
        run: |
          cd frontend
          npm run build
          npm start &
          sleep 10

      - name: Run E2E tests
        run: cd frontend && npx playwright test

      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v3
        with:
          name: playwright-report
          path: frontend/playwright-report/
```

---

## Troubleshooting

### Tests Failing Locally

1. **Check services running**:
   ```bash
   # Backend on :8080
   curl http://localhost:8080/health

   # Frontend on :3000
   curl http://localhost:3000
   ```

2. **Clear test data**:
   ```bash
   cd backend
   go run cmd/seed-test-data/main.go --clean
   go run cmd/seed-test-data/main.go
   ```

3. **Update Playwright**:
   ```bash
   cd frontend
   npm install @playwright/test@latest
   npx playwright install
   ```

### Timeout Errors

Increase timeout in specific tests:

```typescript
test("slow test", async ({ page }) => {
  test.setTimeout(60000); // 60 seconds
  // ... test code
});
```

### Flaky Tests

Add retry logic:

```typescript
test.describe.configure({ retries: 3 });
```

### Screenshot/Video Not Working

Check `playwright.config.ts`:

```typescript
use: {
  screenshot: 'on', // Always take screenshots
  video: 'on',      // Always record video
}
```

---

## Test Coverage

### Current E2E Coverage

| Feature | Tests | Status |
|---------|-------|--------|
| **Member Portal** | 12 tests | ✅ Complete |
| **POS - Product Mgmt** | 8 tests | ✅ Complete |
| **POS - Transactions** | 12 tests | ✅ Complete |
| **POS - Stock** | 3 tests | ✅ Complete |
| **POS - Sales History** | 4 tests | ✅ Complete |
| **POS - Access Control** | 2 tests | ✅ Complete |
| **Accounting** | 0 tests | ⏳ Pending |
| **Reports** | 0 tests | ⏳ Pending |
| **Share Capital** | 0 tests | ⏳ Pending |

**Total E2E Tests: 41 test cases**

### Next Test Files to Create

1. `accounting.spec.ts` - Chart of accounts, journal entries
2. `reports.spec.ts` - Financial reports generation
3. `share-capital.spec.ts` - Member share capital transactions
4. `admin.spec.ts` - Admin panel, user management, settings

---

## Resources

- **Playwright Docs**: https://playwright.dev/
- **Best Practices**: https://playwright.dev/docs/best-practices
- **Selectors Guide**: https://playwright.dev/docs/selectors
- **API Reference**: https://playwright.dev/docs/api/class-test

---

**Last Updated:** January 20, 2025
**Maintained By:** Development Team
