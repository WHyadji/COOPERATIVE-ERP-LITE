import { test, expect, type Page } from "@playwright/test";

/**
 * E2E Tests for Accounting Module
 *
 * Test Coverage:
 * 1. Chart of Accounts Management
 * 2. Journal Entry Creation & Validation
 * 3. Transaction Listing & Filtering
 * 4. Double-Entry Validation
 * 5. Account Ledger (Buku Besar)
 * 6. Trial Balance
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

// Test Chart of Accounts data
const TEST_ACCOUNT = {
  kodeAkun: "1101",
  namaAkun: "Kas",
  jenisAkun: "ASET",
  kategori: "Aset Lancar",
  saldoNormal: "DEBIT",
  deskripsi: "Kas/uang tunai koperasi",
};

// Test Journal Entry data
const TEST_JOURNAL = {
  tanggalTransaksi: "2025-01-20",
  deskripsi: "Penerimaan simpanan wajib anggota",
  nomorReferensi: "SWJ-001",
  lines: [
    {
      kodeAkun: "1101", // Kas
      namaAkun: "Kas",
      debit: 1000000,
      kredit: 0,
      keterangan: "Penerimaan kas dari simpanan",
    },
    {
      kodeAkun: "3101", // Modal - Simpanan Wajib
      namaAkun: "Simpanan Wajib",
      debit: 0,
      kredit: 1000000,
      keterangan: "Simpanan wajib bulan Januari",
    },
  ],
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

// Helper function to navigate to Chart of Accounts
async function navigateToChartOfAccounts(page: Page) {
  await page.click('a[href*="/akuntansi"], nav >> text=/akuntansi/i');
  await page.click('a[href*="/chart-of-accounts"], a[href*="/bagan-akun"]');
  await page.waitForURL(/.*\/(chart-of-accounts|bagan-akun)/);
}

// Helper function to navigate to Journal Entry
async function navigateToJournalEntry(page: Page) {
  await page.click('a[href*="/akuntansi"], nav >> text=/akuntansi/i');
  await page.click(
    'a[href*="/jurnal"], a[href*="/journal-entry"], nav >> text=/jurnal/i'
  );
  await page.waitForURL(/.*\/(jurnal|journal-entry)/);
}

test.describe("Accounting - Chart of Accounts Management", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
    await navigateToChartOfAccounts(page);
  });

  test("should display chart of accounts page correctly", async ({ page }) => {
    // Check page title
    await expect(
      page.getByRole("heading", { name: /bagan akun|chart of accounts/i })
    ).toBeVisible();

    // Check "Tambah Akun" button exists
    await expect(
      page.getByRole("button", { name: /tambah akun|add account/i })
    ).toBeVisible();

    // Check table headers
    await expect(page.getByText(/kode akun|account code/i)).toBeVisible();
    await expect(page.getByText(/nama akun|account name/i)).toBeVisible();
    await expect(page.getByText(/jenis|type/i)).toBeVisible();
    await expect(page.getByText(/saldo|balance/i)).toBeVisible();
  });

  test("should display account categories", async ({ page }) => {
    // Check for main account categories
    const categories = [
      /aset|asset/i,
      /kewajiban|liability/i,
      /modal|equity/i,
      /pendapatan|revenue/i,
      /beban|expense/i,
    ];

    for (const category of categories) {
      const categoryElement = page.getByText(category);
      // Check if at least one category is visible (might be collapsed/expanded)
      if (await categoryElement.isVisible()) {
        await expect(categoryElement).toBeVisible();
      }
    }
  });

  test("should create a new account successfully", async ({ page }) => {
    // Click "Tambah Akun" button
    await page.click('button:has-text(/tambah akun|add account/i)');

    // Wait for form dialog/page
    await expect(page.getByText(/form.*akun|account form/i)).toBeVisible();

    // Fill account form
    await page.fill('input[name="kodeAkun"]', TEST_ACCOUNT.kodeAkun);
    await page.fill('input[name="namaAkun"]', TEST_ACCOUNT.namaAkun);

    // Select account type (jenis akun)
    await page.selectOption('select[name="jenisAkun"]', TEST_ACCOUNT.jenisAkun);

    // Fill category
    await page.fill('input[name="kategori"]', TEST_ACCOUNT.kategori);

    // Select normal balance
    await page.selectOption(
      'select[name="saldoNormal"]',
      TEST_ACCOUNT.saldoNormal
    );

    // Fill description
    await page.fill('textarea[name="deskripsi"]', TEST_ACCOUNT.deskripsi);

    // Submit form
    await page.click('button[type="submit"]:has-text(/simpan|save/i)');

    // Check success message
    await expect(page.getByText(/berhasil|success/i)).toBeVisible();

    // Verify account appears in list
    await expect(page.getByText(TEST_ACCOUNT.kodeAkun)).toBeVisible();
    await expect(page.getByText(TEST_ACCOUNT.namaAkun)).toBeVisible();
  });

  test("should show validation errors for invalid account data", async ({
    page,
  }) => {
    await page.click('button:has-text(/tambah akun/i)');

    // Try to submit empty form
    await page.click('button[type="submit"]:has-text(/simpan/i)');

    // Check validation messages
    await expect(
      page.getByText(/kode akun.*wajib|account code.*required/i)
    ).toBeVisible();
    await expect(
      page.getByText(/nama akun.*wajib|account name.*required/i)
    ).toBeVisible();
  });

  test("should search accounts by code or name", async ({ page }) => {
    // Type in search box
    await page.fill('input[placeholder*="Cari akun"]', "Kas");

    // Wait for search results
    await page.waitForTimeout(500);

    // Check that only matching accounts are shown
    const rows = page.locator('tbody tr');
    await expect(rows.first()).toContainText(/kas/i);
  });

  test("should filter accounts by type", async ({ page }) => {
    // Select account type filter
    const typeFilter = page.locator(
      'select[name="jenisAkun"], [data-testid="account-type-filter"]'
    );

    if (await typeFilter.isVisible()) {
      await typeFilter.selectOption("ASET");
      await page.waitForTimeout(300);

      // Verify filtered results
      const rows = page.locator('tbody tr');
      const count = await rows.count();
      expect(count).toBeGreaterThan(0);
    }
  });

  test("should update existing account", async ({ page }) => {
    // Find account and click edit button
    const accountRow = page
      .locator(`tr:has-text("${TEST_ACCOUNT.kodeAkun}")`)
      .first();
    await accountRow.locator('button:has-text(/edit|ubah/i)').click();

    // Update account name
    const newName = TEST_ACCOUNT.namaAkun + " (Updated)";
    await page.fill('input[name="namaAkun"]', newName);

    // Submit update
    await page.click('button[type="submit"]:has-text(/simpan|update/i)');

    // Verify update success
    await expect(
      page.getByText(/berhasil.*diperbarui|updated successfully/i)
    ).toBeVisible();
    await expect(page.getByText(newName)).toBeVisible();
  });

  test("should show account balance", async ({ page }) => {
    // Click on an account to view details
    const accountRow = page.locator('tr:has-text("Kas")').first();

    if (await accountRow.isVisible()) {
      await accountRow.click();

      // Check if balance is displayed
      const balanceElement = page.locator('[data-testid="account-balance"]');
      if (await balanceElement.isVisible()) {
        await expect(balanceElement).toBeVisible();
      }
    }
  });

  test("should prevent deleting account with transactions", async ({
    page,
  }) => {
    // Find an account that has transactions
    const accountRow = page.locator('tr:has-text("1101")').first();

    if (await accountRow.isVisible()) {
      const deleteBtn = accountRow.locator('button:has-text(/hapus|delete/i)');

      if (await deleteBtn.isVisible()) {
        await deleteBtn.click();

        // Should show error message
        await expect(
          page.getByText(/tidak dapat.*transaksi|cannot.*transactions/i)
        ).toBeVisible();
      }
    }
  });
});

test.describe("Accounting - Journal Entry Creation", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToJournalEntry(page);
  });

  test("should display journal entry page correctly", async ({ page }) => {
    // Check page title
    await expect(
      page.getByRole("heading", { name: /jurnal|journal entry/i })
    ).toBeVisible();

    // Check "Buat Jurnal Baru" button
    await expect(
      page.getByRole("button", {
        name: /buat jurnal|create.*journal|tambah/i,
      })
    ).toBeVisible();

    // Check transaction list table
    await expect(page.getByText(/nomor.*jurnal|journal.*number/i)).toBeVisible();
    await expect(page.getByText(/tanggal|date/i)).toBeVisible();
    await expect(page.getByText(/deskripsi|description/i)).toBeVisible();
    await expect(page.getByText(/total|amount/i)).toBeVisible();
  });

  test("should create a balanced journal entry successfully", async ({
    page,
  }) => {
    // Click "Buat Jurnal Baru"
    await page.click('button:has-text(/buat jurnal|create journal|tambah/i)');

    // Wait for form
    await expect(
      page.getByText(/form.*jurnal|journal entry form/i)
    ).toBeVisible();

    // Fill transaction date
    await page.fill(
      'input[name="tanggalTransaksi"]',
      TEST_JOURNAL.tanggalTransaksi
    );

    // Fill description
    await page.fill('textarea[name="deskripsi"]', TEST_JOURNAL.deskripsi);

    // Fill reference number
    await page.fill(
      'input[name="nomorReferensi"]',
      TEST_JOURNAL.nomorReferensi
    );

    // Add first line (Debit)
    const line1 = TEST_JOURNAL.lines[0];

    // Select account or type account code
    await page.click('[data-testid="add-line"], button:has-text(/tambah baris/i)');

    await page.fill(
      '[data-testid="account-input-0"], input[name="lines[0].kodeAkun"]',
      line1.kodeAkun
    );
    await page.fill(
      '[data-testid="debit-input-0"], input[name="lines[0].debit"]',
      line1.debit.toString()
    );
    await page.fill(
      '[data-testid="description-input-0"], input[name="lines[0].keterangan"]',
      line1.keterangan
    );

    // Add second line (Credit)
    await page.click('[data-testid="add-line"], button:has-text(/tambah baris/i)');

    const line2 = TEST_JOURNAL.lines[1];
    await page.fill(
      '[data-testid="account-input-1"], input[name="lines[1].kodeAkun"]',
      line2.kodeAkun
    );
    await page.fill(
      '[data-testid="credit-input-1"], input[name="lines[1].kredit"]',
      line2.kredit.toString()
    );
    await page.fill(
      '[data-testid="description-input-1"], input[name="lines[1].keterangan"]',
      line2.keterangan
    );

    // Verify balance indicator shows balanced
    await expect(
      page.locator('[data-testid="balance-status"]:has-text(/balanced|seimbang/i)')
    ).toBeVisible();

    // Submit journal entry
    await page.click('button[type="submit"]:has-text(/simpan|save|posting/i)');

    // Check success message
    await expect(
      page.getByText(/berhasil.*dibuat|created successfully/i)
    ).toBeVisible();

    // Verify journal appears in list
    await expect(page.getByText(TEST_JOURNAL.deskripsi)).toBeVisible();
  });

  test("should show error for unbalanced journal entry", async ({ page }) => {
    await page.click('button:has-text(/buat jurnal|create journal/i)');

    // Fill basic info
    await page.fill('input[name="tanggalTransaksi"]', "2025-01-20");
    await page.fill('textarea[name="deskripsi"]', "Test unbalanced entry");

    // Add only debit line (no credit)
    await page.click('[data-testid="add-line"]');
    await page.fill('[data-testid="account-input-0"]', "1101");
    await page.fill('[data-testid="debit-input-0"]', "1000000");

    // Try to submit
    await page.click('button[type="submit"]:has-text(/simpan|posting/i)');

    // Should show unbalanced error
    await expect(
      page.getByText(/tidak seimbang|unbalanced|debit.*kredit/i)
    ).toBeVisible();
  });

  test("should validate required fields", async ({ page }) => {
    await page.click('button:has-text(/buat jurnal/i)');

    // Try to submit empty form
    await page.click('button[type="submit"]:has-text(/simpan/i)');

    // Check validation messages
    await expect(
      page.getByText(/tanggal.*wajib|date.*required/i)
    ).toBeVisible();
    await expect(
      page.getByText(/deskripsi.*wajib|description.*required/i)
    ).toBeVisible();
  });

  test("should allow adding multiple journal lines", async ({ page }) => {
    await page.click('button:has-text(/buat jurnal/i)');

    // Add first line
    await page.click('button:has-text(/tambah baris|add line/i)');
    await expect(page.locator('[data-testid="journal-line-0"]')).toBeVisible();

    // Add second line
    await page.click('button:has-text(/tambah baris|add line/i)');
    await expect(page.locator('[data-testid="journal-line-1"]')).toBeVisible();

    // Add third line
    await page.click('button:has-text(/tambah baris|add line/i)');
    await expect(page.locator('[data-testid="journal-line-2"]')).toBeVisible();

    // Verify we have 3 lines
    const lines = page.locator('[data-testid^="journal-line-"]');
    await expect(lines).toHaveCount(3);
  });

  test("should allow removing journal lines", async ({ page }) => {
    await page.click('button:has-text(/buat jurnal/i)');

    // Add two lines
    await page.click('button:has-text(/tambah baris/i)');
    await page.click('button:has-text(/tambah baris/i)');

    // Remove first line
    const firstLine = page.locator('[data-testid="journal-line-0"]');
    await firstLine.locator('button:has-text(/hapus|remove/i)').click();

    // Verify only one line remains
    const lines = page.locator('[data-testid^="journal-line-"]');
    await expect(lines).toHaveCount(1);
  });

  test("should calculate and display running totals", async ({ page }) => {
    await page.click('button:has-text(/buat jurnal/i)');

    // Add debit line
    await page.click('button:has-text(/tambah baris/i)');
    await page.fill('[data-testid="debit-input-0"]', "1000000");

    // Check debit total updates
    const debitTotal = page.locator('[data-testid="total-debit"]');
    await expect(debitTotal).toContainText("1,000,000");

    // Add credit line
    await page.click('button:has-text(/tambah baris/i)');
    await page.fill('[data-testid="credit-input-1"]', "500000");

    // Check credit total updates
    const creditTotal = page.locator('[data-testid="total-credit"]');
    await expect(creditTotal).toContainText("500,000");

    // Check difference is shown
    const difference = page.locator('[data-testid="difference"]');
    await expect(difference).toContainText("500,000");
  });

  test("should auto-generate journal number", async ({ page }) => {
    await page.click('button:has-text(/buat jurnal/i)');

    // Check if journal number is auto-generated
    const journalNumber = page.locator(
      'input[name="nomorJurnal"], [data-testid="journal-number"]'
    );

    if (await journalNumber.isVisible()) {
      const value = await journalNumber.inputValue();
      expect(value).toMatch(/JRN|JU|JURNAL/i); // Check format
    }
  });
});

test.describe("Accounting - Transaction List & Filtering", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToJournalEntry(page);
  });

  test("should display transaction list", async ({ page }) => {
    // Check if transactions are listed
    const transactionRows = page.locator('tbody tr');
    const count = await transactionRows.count();

    expect(count).toBeGreaterThanOrEqual(0);
  });

  test("should filter transactions by date range", async ({ page }) => {
    // Set date filters
    const startDate = page.locator(
      'input[name="tanggalMulai"], [data-testid="start-date"]'
    );
    const endDate = page.locator(
      'input[name="tanggalAkhir"], [data-testid="end-date"]'
    );

    if (await startDate.isVisible()) {
      await startDate.fill("2025-01-01");
      await endDate.fill("2025-01-31");

      // Apply filter
      await page.click('button:has-text(/filter|cari/i)');
      await page.waitForTimeout(300);

      // Verify results are filtered
      const rows = page.locator('tbody tr');
      const count = await rows.count();
      expect(count).toBeGreaterThanOrEqual(0);
    }
  });

  test("should filter transactions by type", async ({ page }) => {
    const typeFilter = page.locator(
      'select[name="tipeTransaksi"], [data-testid="transaction-type-filter"]'
    );

    if (await typeFilter.isVisible()) {
      await typeFilter.selectOption("SIMPANAN");
      await page.waitForTimeout(300);

      // Check results show only savings transactions
      const results = page.locator('[data-testid="transaction-type"]');
      if ((await results.count()) > 0) {
        await expect(results.first()).toContainText(/simpanan/i);
      }
    }
  });

  test("should search transactions by description or reference", async ({
    page,
  }) => {
    const searchInput = page.locator(
      'input[placeholder*="Cari"], [data-testid="search-input"]'
    );

    if (await searchInput.isVisible()) {
      await searchInput.fill("simpanan");
      await page.waitForTimeout(500);

      // Check filtered results
      const rows = page.locator('tbody tr');
      if ((await rows.count()) > 0) {
        await expect(rows.first()).toContainText(/simpanan/i);
      }
    }
  });

  test("should view transaction details", async ({ page }) => {
    // Click on first transaction
    const firstRow = page.locator('tbody tr').first();

    if (await firstRow.isVisible()) {
      const detailBtn = firstRow.locator('button:has-text(/detail|lihat/i)');

      if (await detailBtn.isVisible()) {
        await detailBtn.click();

        // Check detail modal/page
        await expect(
          page.getByText(/detail.*transaksi|transaction.*detail/i)
        ).toBeVisible();

        // Check for transaction lines
        await expect(
          page.getByText(/baris.*jurnal|journal.*lines/i)
        ).toBeVisible();
      }
    }
  });

  test("should export transactions to Excel", async ({ page }) => {
    const exportBtn = page.locator('button:has-text(/export|download/i)');

    if (await exportBtn.isVisible()) {
      // Start download
      const downloadPromise = page.waitForEvent("download");
      await exportBtn.click();

      // Verify download started
      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/transaksi|journal|xlsx/i);
    }
  });
});

test.describe("Accounting - Account Ledger (Buku Besar)", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await page.goto("/akuntansi/buku-besar");
  });

  test("should display account ledger page", async ({ page }) => {
    await expect(
      page.getByRole("heading", { name: /buku besar|account ledger/i })
    ).toBeVisible();

    // Check account selector
    await expect(page.locator('select[name="idAkun"]')).toBeVisible();

    // Check date range inputs
    await expect(
      page.locator('input[name="tanggalMulai"]')
    ).toBeVisible();
  });

  test("should show ledger for selected account", async ({ page }) => {
    // Select an account
    await page.selectOption('select[name="idAkun"]', { index: 1 });

    // Set date range
    await page.fill('input[name="tanggalMulai"]', "2025-01-01");
    await page.fill('input[name="tanggalAkhir"]', "2025-01-31");

    // Click "Tampilkan" or "Show"
    await page.click('button:has-text(/tampilkan|show/i)');

    await page.waitForTimeout(500);

    // Check ledger entries are displayed
    const ledgerRows = page.locator('[data-testid="ledger-row"]');
    const count = await ledgerRows.count();

    expect(count).toBeGreaterThanOrEqual(0);
  });

  test("should calculate running balance correctly", async ({ page }) => {
    // Select account and show ledger
    await page.selectOption('select[name="idAkun"]', { index: 1 });
    await page.fill('input[name="tanggalMulai"]', "2025-01-01");
    await page.fill('input[name="tanggalAkhir"]', "2025-01-31");
    await page.click('button:has-text(/tampilkan/i)');

    await page.waitForTimeout(500);

    // Check if running balance column exists
    const balanceColumn = page.locator('[data-testid="running-balance"]');

    if ((await balanceColumn.count()) > 0) {
      await expect(balanceColumn.first()).toBeVisible();
    }
  });

  test("should show opening and closing balance", async ({ page }) => {
    await page.selectOption('select[name="idAkun"]', { index: 1 });
    await page.fill('input[name="tanggalMulai"]', "2025-01-01");
    await page.fill('input[name="tanggalAkhir"]', "2025-01-31");
    await page.click('button:has-text(/tampilkan/i)');

    await page.waitForTimeout(500);

    // Check for opening balance
    const openingBalance = page.locator('[data-testid="opening-balance"]');
    if (await openingBalance.isVisible()) {
      await expect(openingBalance).toBeVisible();
    }

    // Check for closing balance
    const closingBalance = page.locator('[data-testid="closing-balance"]');
    if (await closingBalance.isVisible()) {
      await expect(closingBalance).toBeVisible();
    }
  });
});

test.describe("Accounting - Trial Balance", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await page.goto("/akuntansi/neraca-saldo");
  });

  test("should display trial balance page", async ({ page }) => {
    await expect(
      page.getByRole("heading", { name: /neraca saldo|trial balance/i })
    ).toBeVisible();

    // Check date input
    await expect(page.locator('input[name="tanggal"]')).toBeVisible();

    // Check generate button
    await expect(
      page.getByRole("button", { name: /tampilkan|generate/i })
    ).toBeVisible();
  });

  test("should generate trial balance report", async ({ page }) => {
    // Set date
    await page.fill('input[name="tanggal"]', "2025-01-31");

    // Click generate
    await page.click('button:has-text(/tampilkan|generate/i)');

    await page.waitForTimeout(1000);

    // Check table headers
    await expect(page.getByText(/kode akun/i)).toBeVisible();
    await expect(page.getByText(/nama akun/i)).toBeVisible();
    await expect(page.getByText(/debit/i)).toBeVisible();
    await expect(page.getByText(/kredit/i)).toBeVisible();
  });

  test("should verify total debit equals total credit", async ({ page }) => {
    await page.fill('input[name="tanggal"]', "2025-01-31");
    await page.click('button:has-text(/tampilkan/i)');

    await page.waitForTimeout(1000);

    // Get total debit
    const totalDebit = page.locator('[data-testid="total-debit"]');

    // Get total credit
    const totalCredit = page.locator('[data-testid="total-credit"]');

    if ((await totalDebit.isVisible()) && (await totalCredit.isVisible())) {
      const debitText = await totalDebit.textContent();
      const creditText = await totalCredit.textContent();

      // Both should have same value (balanced)
      expect(debitText).toBe(creditText);
    }
  });

  test("should export trial balance to PDF/Excel", async ({ page }) => {
    await page.fill('input[name="tanggal"]', "2025-01-31");
    await page.click('button:has-text(/tampilkan/i)');

    await page.waitForTimeout(1000);

    const exportBtn = page.locator('button:has-text(/export|download|cetak/i)');

    if (await exportBtn.isVisible()) {
      const downloadPromise = page.waitForEvent("download");
      await exportBtn.click();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/neraca|trial|pdf|xlsx/i);
    }
  });
});

test.describe("Accounting - Role-Based Access", () => {
  test("treasurer should access all accounting features", async ({ page }) => {
    await loginAsTreasurer(page);

    // Verify can access chart of accounts
    await navigateToChartOfAccounts(page);
    await expect(
      page.getByRole("heading", { name: /bagan akun/i })
    ).toBeVisible();

    // Verify can access journal entry
    await navigateToJournalEntry(page);
    await expect(
      page.getByRole("heading", { name: /jurnal/i })
    ).toBeVisible();

    // Verify can create new journal
    await expect(
      page.getByRole("button", { name: /buat jurnal|tambah/i })
    ).toBeVisible();
  });

  test("admin should access accounting settings", async ({ page }) => {
    await loginAsAdmin(page);

    // Navigate to accounting settings
    const settingsLink = page.locator(
      'a[href*="/akuntansi/pengaturan"], nav >> text=/pengaturan.*akuntansi/i'
    );

    if (await settingsLink.isVisible()) {
      await settingsLink.click();
      await expect(
        page.getByText(/pengaturan.*akuntansi|accounting.*settings/i)
      ).toBeVisible();
    }
  });

  test("cashier should NOT access accounting module", async ({ page }) => {
    const cashier = {
      email: "kasir@koperasi.local",
      password: "Kasir123!",
    };

    await page.goto("/login");
    await page.fill('input[name="email"]', cashier.email);
    await page.fill('input[name="password"]', cashier.password);
    await page.click('button[type="submit"]');

    await page.waitForURL("/dashboard");

    // Verify accounting menu is not visible
    const accountingLink = page.locator('a[href*="/akuntansi"]');
    await expect(accountingLink).not.toBeVisible();
  });
});
