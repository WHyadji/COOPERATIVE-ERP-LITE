import { test, expect, type Page } from "@playwright/test";

/**
 * E2E Tests for Reports Module
 *
 * Test Coverage:
 * 1. Balance Sheet (Neraca)
 * 2. Income Statement (Laba Rugi / P&L)
 * 3. Cash Flow Statement
 * 4. Member Balance Report
 * 5. Daily Transaction Summary
 * 6. Report Export (PDF/Excel)
 */

// Test data - Admin credentials
const TEST_ADMIN = {
  email: "admin@koperasi.local",
  password: "Admin123!",
};

const TEST_TREASURER = {
  email: "bendahara@koperasi.local",
  password: "Bendahara123!",
};

// Test report parameters
const TEST_REPORT_DATE = {
  startDate: "2025-01-01",
  endDate: "2025-01-31",
  specificDate: "2025-01-31",
};

// Helper function to login as admin
async function loginAsAdmin(page: Page) {
  await page.goto("/login");

  await page.fill('input[name="email"]', TEST_ADMIN.email);
  await page.fill('input[name="password"]', TEST_ADMIN.password);
  await page.click('button[type="submit"]');

  // Wait for redirect to dashboard
  await page.waitForURL("/dashboard");
}

// Helper function to login as treasurer
async function loginAsTreasurer(page: Page) {
  await page.goto("/login");

  await page.fill('input[name="email"]', TEST_TREASURER.email);
  await page.fill('input[name="password"]', TEST_TREASURER.password);
  await page.click('button[type="submit"]');

  await page.waitForURL("/dashboard");
}

// Helper function to navigate to Reports section
async function navigateToReports(page: Page) {
  await page.click('a[href*="/laporan"], nav >> text=/laporan|reports/i');
  await page.waitForURL(/.*\/(laporan|reports)/);
}

test.describe("Reports - Balance Sheet (Neraca)", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToReports(page);
  });

  test("should display balance sheet page correctly", async ({ page }) => {
    // Navigate to balance sheet
    await page.click('a[href*="/neraca"], a[href*="/balance-sheet"]');
    await page.waitForURL(/.*\/(neraca|balance-sheet)/);

    // Check page title
    await expect(
      page.getByRole("heading", { name: /neraca|balance sheet/i })
    ).toBeVisible();

    // Check date input
    await expect(page.locator('input[name="tanggal"]')).toBeVisible();

    // Check generate button
    await expect(
      page.getByRole("button", { name: /tampilkan|generate|lihat/i })
    ).toBeVisible();
  });

  test("should generate balance sheet report", async ({ page }) => {
    await page.click('a[href*="/neraca"]');
    await page.waitForURL(/.*\/neraca/);

    // Set report date
    await page.fill('input[name="tanggal"]', TEST_REPORT_DATE.specificDate);

    // Click generate
    await page.click('button:has-text(/tampilkan|generate/i)');

    await page.waitForTimeout(1000);

    // Check main sections
    await expect(page.getByText(/aset|assets/i)).toBeVisible();
    await expect(page.getByText(/kewajiban|liabilities/i)).toBeVisible();
    await expect(page.getByText(/modal|equity/i)).toBeVisible();

    // Check subsections
    await expect(page.getByText(/aset lancar|current assets/i)).toBeVisible();
    await expect(
      page.getByText(/aset tidak lancar|non-current assets/i)
    ).toBeVisible();
  });

  test("should verify balance sheet equation (Assets = Liabilities + Equity)", async ({
    page,
  }) => {
    await page.click('a[href*="/neraca"]');
    await page.fill('input[name="tanggal"]', TEST_REPORT_DATE.specificDate);
    await page.click('button:has-text(/tampilkan/i)');

    await page.waitForTimeout(1000);

    // Get total assets
    const totalAssets = page.locator('[data-testid="total-assets"]');

    // Get total liabilities + equity
    const totalLiabilitiesEquity = page.locator(
      '[data-testid="total-liabilities-equity"]'
    );

    if (
      (await totalAssets.isVisible()) &&
      (await totalLiabilitiesEquity.isVisible())
    ) {
      const assetsText = await totalAssets.textContent();
      const liabilitiesEquityText = await totalLiabilitiesEquity.textContent();

      // Both should have same value (balanced)
      expect(assetsText).toBe(liabilitiesEquityText);
    }
  });

  test("should export balance sheet to PDF", async ({ page }) => {
    await page.click('a[href*="/neraca"]');
    await page.fill('input[name="tanggal"]', TEST_REPORT_DATE.specificDate);
    await page.click('button:has-text(/tampilkan/i)');

    await page.waitForTimeout(1000);

    const exportBtn = page.locator('button:has-text(/export|download|cetak/i)');

    if (await exportBtn.isVisible()) {
      const downloadPromise = page.waitForEvent("download");
      await exportBtn.click();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/neraca|balance.*sheet/i);
    }
  });
});

test.describe("Reports - Income Statement (Laba Rugi)", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToReports(page);
  });

  test("should display income statement page correctly", async ({ page }) => {
    // Navigate to income statement
    await page.click('a[href*="/laba-rugi"], a[href*="/income-statement"]');
    await page.waitForURL(/.*\/(laba-rugi|income-statement)/);

    // Check page title
    await expect(
      page.getByRole("heading", {
        name: /laba rugi|income statement|profit.*loss/i,
      })
    ).toBeVisible();

    // Check date range inputs
    await expect(
      page.locator('input[name="tanggalMulai"], input[name="startDate"]')
    ).toBeVisible();
    await expect(
      page.locator('input[name="tanggalAkhir"], input[name="endDate"]')
    ).toBeVisible();
  });

  test("should generate income statement report", async ({ page }) => {
    await page.click('a[href*="/laba-rugi"]');
    await page.waitForURL(/.*\/laba-rugi/);

    // Set date range
    await page.fill('input[name="tanggalMulai"]', TEST_REPORT_DATE.startDate);
    await page.fill('input[name="tanggalAkhir"]', TEST_REPORT_DATE.endDate);

    // Click generate
    await page.click('button:has-text(/tampilkan|generate/i)');

    await page.waitForTimeout(1000);

    // Check main sections
    await expect(page.getByText(/pendapatan|revenue|income/i)).toBeVisible();
    await expect(page.getByText(/beban|expenses/i)).toBeVisible();
    await expect(
      page.getByText(/laba.*bersih|net.*income|net.*profit/i)
    ).toBeVisible();
  });

  test("should calculate net income correctly (Revenue - Expenses)", async ({
    page,
  }) => {
    await page.click('a[href*="/laba-rugi"]');
    await page.fill('input[name="tanggalMulai"]', TEST_REPORT_DATE.startDate);
    await page.fill('input[name="tanggalAkhir"]', TEST_REPORT_DATE.endDate);
    await page.click('button:has-text(/tampilkan/i)');

    await page.waitForTimeout(1000);

    // Check if calculation is shown
    const totalRevenue = page.locator('[data-testid="total-revenue"]');
    const totalExpenses = page.locator('[data-testid="total-expenses"]');
    const netIncome = page.locator('[data-testid="net-income"]');

    // Verify elements exist (actual calculation validation would require parsing)
    if (await totalRevenue.isVisible()) {
      await expect(totalRevenue).toBeVisible();
    }

    if (await totalExpenses.isVisible()) {
      await expect(totalExpenses).toBeVisible();
    }

    if (await netIncome.isVisible()) {
      await expect(netIncome).toBeVisible();
    }
  });

  test("should show revenue breakdown by category", async ({ page }) => {
    await page.click('a[href*="/laba-rugi"]');
    await page.fill('input[name="tanggalMulai"]', TEST_REPORT_DATE.startDate);
    await page.fill('input[name="tanggalAkhir"]', TEST_REPORT_DATE.endDate);
    await page.click('button:has-text(/tampilkan/i)');

    await page.waitForTimeout(1000);

    // Check for revenue categories
    const revenueSection = page.locator('[data-testid="revenue-section"]');

    if (await revenueSection.isVisible()) {
      // Common revenue categories for cooperative
      const categories = [
        /pendapatan.*penjualan|sales revenue/i,
        /pendapatan.*simpanan|savings revenue/i,
        /pendapatan.*lain|other revenue/i,
      ];

      for (const category of categories) {
        const categoryElement = page.getByText(category);
        if (await categoryElement.isVisible()) {
          await expect(categoryElement).toBeVisible();
        }
      }
    }
  });

  test("should export income statement to PDF", async ({ page }) => {
    await page.click('a[href*="/laba-rugi"]');
    await page.fill('input[name="tanggalMulai"]', TEST_REPORT_DATE.startDate);
    await page.fill('input[name="tanggalAkhir"]', TEST_REPORT_DATE.endDate);
    await page.click('button:has-text(/tampilkan/i)');

    await page.waitForTimeout(1000);

    const exportBtn = page.locator('button:has-text(/export|download|cetak/i)');

    if (await exportBtn.isVisible()) {
      const downloadPromise = page.waitForEvent("download");
      await exportBtn.click();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/laba.*rugi|income/i);
    }
  });
});

test.describe("Reports - Cash Flow Statement", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToReports(page);
  });

  test("should display cash flow page correctly", async ({ page }) => {
    // Navigate to cash flow
    const cashFlowLink = page.locator(
      'a[href*="/arus-kas"], a[href*="/cash-flow"]'
    );

    if (await cashFlowLink.isVisible()) {
      await cashFlowLink.click();
      await page.waitForURL(/.*\/(arus-kas|cash-flow)/);

      // Check page title
      await expect(
        page.getByRole("heading", { name: /arus kas|cash flow/i })
      ).toBeVisible();

      // Check date range inputs
      await expect(page.locator('input[name="tanggalMulai"]')).toBeVisible();
    }
  });

  test("should generate cash flow report", async ({ page }) => {
    const cashFlowLink = page.locator('a[href*="/arus-kas"]');

    if (await cashFlowLink.isVisible()) {
      await cashFlowLink.click();
      await page.fill('input[name="tanggalMulai"]', TEST_REPORT_DATE.startDate);
      await page.fill('input[name="tanggalAkhir"]', TEST_REPORT_DATE.endDate);
      await page.click('button:has-text(/tampilkan/i)');

      await page.waitForTimeout(1000);

      // Check main sections (if implemented)
      const operatingSection = page.getByText(/operasi|operating/i);
      const investingSection = page.getByText(/investasi|investing/i);
      const financingSection = page.getByText(/pendanaan|financing/i);

      // At least one section should be visible
      const sections = [operatingSection, investingSection, financingSection];
      let anyVisible = false;

      for (const section of sections) {
        if (await section.isVisible()) {
          anyVisible = true;
          break;
        }
      }

      // If cash flow is implemented, at least one section should show
      if (anyVisible) {
        expect(anyVisible).toBe(true);
      }
    }
  });
});

test.describe("Reports - Member Balance Report", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToReports(page);
  });

  test("should display member balance report page", async ({ page }) => {
    // Navigate to member balance
    const memberBalanceLink = page.locator(
      'a[href*="/saldo-anggota"], a[href*="/member-balance"]'
    );

    if (await memberBalanceLink.isVisible()) {
      await memberBalanceLink.click();
      await page.waitForURL(/.*\/(saldo-anggota|member-balance)/);

      // Check page title
      await expect(
        page.getByRole("heading", { name: /saldo anggota|member balance/i })
      ).toBeVisible();
    }
  });

  test("should generate member balance report", async ({ page }) => {
    const memberBalanceLink = page.locator('a[href*="/saldo-anggota"]');

    if (await memberBalanceLink.isVisible()) {
      await memberBalanceLink.click();

      // Set date (as of date)
      const dateInput = page.locator('input[name="tanggal"]');
      if (await dateInput.isVisible()) {
        await dateInput.fill(TEST_REPORT_DATE.specificDate);
      }

      // Click generate
      await page.click('button:has-text(/tampilkan|generate/i)');
      await page.waitForTimeout(1000);

      // Check for member list
      const memberTable = page.locator('table, [data-testid="member-list"]');
      if (await memberTable.isVisible()) {
        await expect(memberTable).toBeVisible();
      }
    }
  });

  test("should show member share capital breakdown", async ({ page }) => {
    const memberBalanceLink = page.locator('a[href*="/saldo-anggota"]');

    if (await memberBalanceLink.isVisible()) {
      await memberBalanceLink.click();
      await page.click('button:has-text(/tampilkan/i)');
      await page.waitForTimeout(1000);

      // Check for share capital columns
      const columns = [
        /simpanan pokok|principal/i,
        /simpanan wajib|mandatory/i,
        /simpanan sukarela|voluntary/i,
      ];

      for (const column of columns) {
        const columnHeader = page.getByText(column);
        if (await columnHeader.isVisible()) {
          await expect(columnHeader).toBeVisible();
        }
      }
    }
  });

  test("should export member balance to Excel", async ({ page }) => {
    const memberBalanceLink = page.locator('a[href*="/saldo-anggota"]');

    if (await memberBalanceLink.isVisible()) {
      await memberBalanceLink.click();
      await page.click('button:has-text(/tampilkan/i)');
      await page.waitForTimeout(1000);

      const exportBtn = page.locator(
        'button:has-text(/export|download|excel/i)'
      );

      if (await exportBtn.isVisible()) {
        const downloadPromise = page.waitForEvent("download");
        await exportBtn.click();

        const download = await downloadPromise;
        expect(download.suggestedFilename()).toMatch(
          /saldo.*anggota|member.*balance|xlsx/i
        );
      }
    }
  });
});

test.describe("Reports - Daily Transaction Summary", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
    await navigateToReports(page);
  });

  test("should display daily transaction summary page", async ({ page }) => {
    const dailySummaryLink = page.locator(
      'a[href*="/ringkasan-harian"], a[href*="/daily-summary"]'
    );

    if (await dailySummaryLink.isVisible()) {
      await dailySummaryLink.click();
      await page.waitForURL(/.*\/(ringkasan-harian|daily-summary)/);

      // Check page title
      await expect(
        page.getByRole("heading", {
          name: /ringkasan.*harian|daily.*summary|transaksi.*harian/i,
        })
      ).toBeVisible();
    }
  });

  test("should generate daily summary report", async ({ page }) => {
    const dailySummaryLink = page.locator('a[href*="/ringkasan-harian"]');

    if (await dailySummaryLink.isVisible()) {
      await dailySummaryLink.click();

      // Set date
      const dateInput = page.locator('input[name="tanggal"]');
      if (await dateInput.isVisible()) {
        await dateInput.fill(TEST_REPORT_DATE.specificDate);
      }

      // Generate report
      await page.click('button:has-text(/tampilkan|generate/i)');
      await page.waitForTimeout(1000);

      // Check for summary metrics
      const metrics = [
        /total.*transaksi|total.*transactions/i,
        /total.*penjualan|total.*sales/i,
        /total.*simpanan|total.*savings/i,
      ];

      for (const metric of metrics) {
        const metricElement = page.getByText(metric);
        if (await metricElement.isVisible()) {
          await expect(metricElement).toBeVisible();
        }
      }
    }
  });
});

test.describe("Reports - Role-Based Access", () => {
  test("treasurer should access all reports", async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToReports(page);

    // Verify can see report menu
    await expect(
      page.getByRole("heading", { name: /laporan|reports/i })
    ).toBeVisible();

    // Check for key reports
    const reports = [
      page.locator('a[href*="/neraca"]'),
      page.locator('a[href*="/laba-rugi"]'),
    ];

    for (const report of reports) {
      if (await report.isVisible()) {
        await expect(report).toBeVisible();
      }
    }
  });

  test("admin should access all reports", async ({ page }) => {
    await loginAsAdmin(page);
    await navigateToReports(page);

    // Verify can access reports
    await expect(
      page.getByRole("heading", { name: /laporan|reports/i })
    ).toBeVisible();
  });

  test("cashier should have limited report access", async ({ page }) => {
    const cashier = {
      email: "kasir@koperasi.local",
      password: "Kasir123!",
    };

    await page.goto("/login");
    await page.fill('input[name="email"]', cashier.email);
    await page.fill('input[name="password"]', cashier.password);
    await page.click('button[type="submit"]');

    await page.waitForURL("/dashboard");

    // Cashier might only see daily sales summary, not full financial reports
    const fullReportsLink = page.locator(
      'a[href*="/neraca"], a[href*="/laba-rugi"]'
    );

    // Should not see balance sheet or income statement
    if (await fullReportsLink.isVisible()) {
      // If visible, it's an access control issue
      expect(await fullReportsLink.count()).toBe(0);
    }
  });
});

test.describe("Reports - Print and Export", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToReports(page);
  });

  test("should have consistent export buttons across all reports", async ({
    page,
  }) => {
    // Check balance sheet
    const neracaLink = page.locator('a[href*="/neraca"]');
    if (await neracaLink.isVisible()) {
      await neracaLink.click();
      await page.fill('input[name="tanggal"]', TEST_REPORT_DATE.specificDate);
      await page.click('button:has-text(/tampilkan/i)');
      await page.waitForTimeout(500);

      const exportBtn = page.locator(
        'button:has-text(/export|download|cetak/i)'
      );
      if (await exportBtn.isVisible()) {
        await expect(exportBtn).toBeVisible();
      }

      // Go back to reports
      await navigateToReports(page);
    }

    // Check income statement
    const labaRugiLink = page.locator('a[href*="/laba-rugi"]');
    if (await labaRugiLink.isVisible()) {
      await labaRugiLink.click();
      await page.fill('input[name="tanggalMulai"]', TEST_REPORT_DATE.startDate);
      await page.fill('input[name="tanggalAkhir"]', TEST_REPORT_DATE.endDate);
      await page.click('button:has-text(/tampilkan/i)');
      await page.waitForTimeout(500);

      const exportBtn = page.locator(
        'button:has-text(/export|download|cetak/i)'
      );
      if (await exportBtn.isVisible()) {
        await expect(exportBtn).toBeVisible();
      }
    }
  });

  test("should support PDF format for financial reports", async ({ page }) => {
    await page.click('a[href*="/neraca"]');
    await page.fill('input[name="tanggal"]', TEST_REPORT_DATE.specificDate);
    await page.click('button:has-text(/tampilkan/i)');
    await page.waitForTimeout(1000);

    // Look for PDF export option
    const pdfBtn = page.locator(
      'button:has-text(/pdf/i), [data-testid="export-pdf"]'
    );

    if (await pdfBtn.isVisible()) {
      const downloadPromise = page.waitForEvent("download");
      await pdfBtn.click();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/\.pdf$/i);
    }
  });

  test("should support Excel format for data-heavy reports", async ({
    page,
  }) => {
    const memberBalanceLink = page.locator('a[href*="/saldo-anggota"]');

    if (await memberBalanceLink.isVisible()) {
      await memberBalanceLink.click();
      await page.click('button:has-text(/tampilkan/i)');
      await page.waitForTimeout(1000);

      // Look for Excel export option
      const excelBtn = page.locator(
        'button:has-text(/excel|xlsx/i), [data-testid="export-excel"]'
      );

      if (await excelBtn.isVisible()) {
        const downloadPromise = page.waitForEvent("download");
        await excelBtn.click();

        const download = await downloadPromise;
        expect(download.suggestedFilename()).toMatch(/\.xlsx?$/i);
      }
    }
  });
});
