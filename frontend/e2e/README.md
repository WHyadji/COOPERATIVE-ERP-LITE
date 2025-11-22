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

**Test Scenarios:** 35 test cases covering:
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

### 3. `accounting.spec.ts`
**Test Coverage:**
- Chart of Accounts management (CRUD operations)
- Journal Entry creation with double-entry validation
- Transaction listing and filtering
- Account Ledger (Buku Besar) with running balance
- Trial Balance (Neraca Saldo) generation
- Role-based access control (Treasurer, Admin, Cashier)

**Status:** ✅ Complete

**Test Scenarios:** 34 test cases covering:
- ✅ Display chart of accounts with categories
- ✅ Create/update accounts with validation
- ✅ Search and filter accounts by type
- ✅ Prevent deleting accounts with transactions
- ✅ Create balanced journal entries (Debit = Credit)
- ✅ Validate unbalanced journal entries
- ✅ Add/remove multiple journal lines
- ✅ Calculate running totals (Debit/Credit)
- ✅ Auto-generate journal numbers
- ✅ Filter transactions by date range and type
- ✅ Search transactions by description/reference
- ✅ View transaction details with journal lines
- ✅ Export transactions to Excel
- ✅ Display account ledger with running balance
- ✅ Show opening and closing balance
- ✅ Generate trial balance report
- ✅ Verify total debit equals total credit
- ✅ Export trial balance to PDF/Excel
- ✅ Role-based access (Treasurer full access, Cashier no access)

### 4. `reports.spec.ts`
**Test Coverage:**
- Balance Sheet (Neraca) report generation
- Income Statement (Laba Rugi / P&L) report generation
- Cash Flow Statement report
- Member Balance Report with share capital breakdown
- Daily Transaction Summary
- Report export functionality (PDF/Excel)
- Role-based report access control

**Status:** ✅ Complete

**Test Scenarios:** 23 test cases covering:
- ✅ Display and generate Balance Sheet report
- ✅ Verify balance sheet equation (Assets = Liabilities + Equity)
- ✅ Display and generate Income Statement report
- ✅ Calculate net income (Revenue - Expenses)
- ✅ Show revenue breakdown by category
- ✅ Display and generate Cash Flow report
- ✅ Show cash flow sections (Operating, Investing, Financing)
- ✅ Display and generate Member Balance report
- ✅ Show share capital breakdown (Pokok, Wajib, Sukarela)
- ✅ Display and generate Daily Transaction Summary
- ✅ Export balance sheet to PDF
- ✅ Export income statement to PDF
- ✅ Export member balance to Excel
- ✅ Consistent export buttons across all reports
- ✅ Support PDF format for financial reports
- ✅ Support Excel format for data-heavy reports
- ✅ Treasurer and Admin access all reports
- ✅ Cashier has limited report access

### 5. `share-capital.spec.ts`
**Test Coverage:**
- Share capital deposits (Simpanan Pokok, Wajib, Sukarela)
- Member balance viewing and calculations
- Transaction history with filtering
- Member dashboard integration
- Summary statistics
- Role-based access control

**Status:** ✅ Complete

**Test Scenarios:** 27 test cases covering:
- ✅ Display share capital deposit page correctly
- ✅ Record principal deposit (Simpanan Pokok) successfully
- ✅ Record mandatory deposit (Simpanan Wajib) successfully
- ✅ Record voluntary deposit (Simpanan Sukarela) successfully
- ✅ Validate required fields and minimum amounts
- ✅ Auto-generate reference numbers
- ✅ Display member balance list with all deposit types
- ✅ Show individual member balance details
- ✅ Calculate total share capital correctly
- ✅ Export member balance to Excel
- ✅ Display and filter transaction history
- ✅ Filter by deposit type (POKOK, WAJIB, SUKARELA)
- ✅ Filter by date range
- ✅ Search by member name/number
- ✅ View transaction details
- ✅ Export transaction history to Excel
- ✅ Member view own balance on dashboard
- ✅ Member view recent transactions
- ✅ Member cannot record deposits (read-only)
- ✅ Display cooperative-wide summary
- ✅ Show number of active members
- ✅ Display monthly deposit trend
- ✅ Treasurer full access
- ✅ Admin full access
- ✅ Member read-only access
- ✅ Cashier no access to share capital

### 6. `admin.spec.ts`
**Test Coverage:**
- User Management (CRUD operations)
- Role Assignment (Admin, Bendahara, Kasir, Anggota)
- System Settings configuration
- Audit Logs viewing and filtering
- Role-Based Access Control

**Status:** ✅ Complete

**Test Scenarios:** 21 test cases covering:
- ✅ Display user management page correctly
- ✅ Create new user with role assignment
- ✅ Validate required fields
- ✅ Update existing user information
- ✅ Deactivate/activate user accounts
- ✅ Filter users by role (Admin, Bendahara, Kasir, Anggota)
- ✅ Search users by name or email
- ✅ Reset user password
- ✅ Delete user with confirmation
- ✅ Display system settings page
- ✅ Update cooperative information
- ✅ Configure accounting settings
- ✅ Display audit log page
- ✅ Filter audit logs by date range
- ✅ Filter audit logs by user
- ✅ Filter audit logs by action type
- ✅ Export audit logs to Excel
- ✅ Admin full access to admin panel
- ✅ Treasurer NO access to admin panel
- ✅ Cashier NO access to admin panel
- ✅ Member NO access to admin panel

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
| **Accounting - Chart of Accounts** | 9 tests | ✅ Complete |
| **Accounting - Journal Entry** | 8 tests | ✅ Complete |
| **Accounting - Transactions** | 6 tests | ✅ Complete |
| **Accounting - Account Ledger** | 4 tests | ✅ Complete |
| **Accounting - Trial Balance** | 4 tests | ✅ Complete |
| **Accounting - Access Control** | 3 tests | ✅ Complete |
| **Reports - Balance Sheet** | 4 tests | ✅ Complete |
| **Reports - Income Statement** | 5 tests | ✅ Complete |
| **Reports - Cash Flow** | 2 tests | ✅ Complete |
| **Reports - Member Balance** | 4 tests | ✅ Complete |
| **Reports - Daily Summary** | 2 tests | ✅ Complete |
| **Reports - Access & Export** | 6 tests | ✅ Complete |
| **Share Capital - Record Deposits** | 7 tests | ✅ Complete |
| **Share Capital - View Balance** | 4 tests | ✅ Complete |
| **Share Capital - Transaction History** | 6 tests | ✅ Complete |
| **Share Capital - Member Dashboard** | 3 tests | ✅ Complete |
| **Share Capital - Summary** | 3 tests | ✅ Complete |
| **Share Capital - Access Control** | 4 tests | ✅ Complete |

**Total E2E Tests: 131 test cases** (up from 104)

### Next Test Files to Create

1. `admin.spec.ts` - Admin panel, user management, settings (~8 tests)

---

## Resources

- **Playwright Docs**: https://playwright.dev/
- **Best Practices**: https://playwright.dev/docs/best-practices
- **Selectors Guide**: https://playwright.dev/docs/selectors
- **API Reference**: https://playwright.dev/docs/api/class-test

---

**Last Updated:** January 20, 2025
**Maintained By:** Development Team
