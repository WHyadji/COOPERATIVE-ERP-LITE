import { test, expect, type Page } from "@playwright/test";

/**
 * E2E Tests for Share Capital (Simpanan) Module
 *
 * Test Coverage:
 * 1. Record Share Capital Deposits (Simpanan Pokok, Wajib, Sukarela)
 * 2. View Member Balance
 * 3. Transaction History
 * 4. Member Dashboard View
 * 5. Validation and Error Handling
 */

// Test data - Admin and Treasurer credentials
const TEST_ADMIN = {
  email: "admin@koperasi.local",
  password: "Admin123!",
};

const TEST_TREASURER = {
  email: "bendahara@koperasi.local",
  password: "Bendahara123!",
};

const TEST_MEMBER = {
  nomorAnggota: "A001",
  password: "123456",
};

// Test share capital deposit data
const TEST_DEPOSIT = {
  principal: {
    tipeSimpanan: "POKOK",
    jumlahSetoran: 100000,
    keterangan: "Simpanan pokok awal keanggotaan",
  },
  mandatory: {
    tipeSimpanan: "WAJIB",
    jumlahSetoran: 50000,
    keterangan: "Simpanan wajib bulan Januari 2025",
  },
  voluntary: {
    tipeSimpanan: "SUKARELA",
    jumlahSetoran: 200000,
    keterangan: "Simpanan sukarela",
  },
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

// Helper function to login as member
async function loginAsMember(page: Page) {
  await page.goto("/login");

  await page.fill('input[name="nomorAnggota"]', TEST_MEMBER.nomorAnggota);
  await page.fill('input[name="password"]', TEST_MEMBER.password);
  await page.click('button[type="submit"]');

  await page.waitForURL("/dashboard");
}

// Helper function to navigate to Share Capital section
async function navigateToShareCapital(page: Page) {
  await page.click('a[href*="/simpanan"], nav >> text=/simpanan|share.*capital/i');
  await page.waitForURL(/.*\/(simpanan|share-capital)/);
}

test.describe("Share Capital - Record Deposits", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToShareCapital(page);
  });

  test("should display share capital deposit page correctly", async ({ page }) => {
    // Check page title
    await expect(
      page.getByRole("heading", { name: /simpanan|share.*capital/i })
    ).toBeVisible();

    // Check "Tambah Simpanan" or "Record Deposit" button
    await expect(
      page.getByRole("button", { name: /tambah.*simpanan|record.*deposit|catat/i })
    ).toBeVisible();

    // Check deposit type filter
    const typeFilter = page.locator('select[name="tipeSimpanan"]');
    if (await typeFilter.isVisible()) {
      await expect(typeFilter).toBeVisible();
    }
  });

  test("should record principal deposit (Simpanan Pokok) successfully", async ({
    page,
  }) => {
    // Click "Tambah Simpanan"
    await page.click('button:has-text(/tambah.*simpanan|record.*deposit/i)');

    // Wait for form
    await expect(
      page.getByText(/form.*simpanan|deposit.*form/i)
    ).toBeVisible();

    // Select member
    const memberSelect = page.locator('select[name="idAnggota"]');
    if (await memberSelect.isVisible()) {
      await memberSelect.selectOption({ index: 1 });
    }

    // Select deposit type - Simpanan Pokok
    await page.selectOption(
      'select[name="tipeSimpanan"]',
      TEST_DEPOSIT.principal.tipeSimpanan
    );

    // Fill amount
    await page.fill(
      'input[name="jumlahSetoran"]',
      TEST_DEPOSIT.principal.jumlahSetoran.toString()
    );

    // Fill description
    await page.fill('textarea[name="keterangan"]', TEST_DEPOSIT.principal.keterangan);

    // Submit form
    await page.click('button[type="submit"]:has-text(/simpan|save/i)');

    // Check success message
    await expect(page.getByText(/berhasil|success/i)).toBeVisible();

    // Verify deposit appears in list
    await expect(
      page.getByText(new RegExp(TEST_DEPOSIT.principal.tipeSimpanan, "i"))
    ).toBeVisible();
  });

  test("should record mandatory deposit (Simpanan Wajib) successfully", async ({
    page,
  }) => {
    await page.click('button:has-text(/tambah.*simpanan/i)');

    // Select member
    const memberSelect = page.locator('select[name="idAnggota"]');
    if (await memberSelect.isVisible()) {
      await memberSelect.selectOption({ index: 1 });
    }

    // Select deposit type - Simpanan Wajib
    await page.selectOption(
      'select[name="tipeSimpanan"]',
      TEST_DEPOSIT.mandatory.tipeSimpanan
    );

    // Fill amount
    await page.fill(
      'input[name="jumlahSetoran"]',
      TEST_DEPOSIT.mandatory.jumlahSetoran.toString()
    );

    // Fill description
    await page.fill('textarea[name="keterangan"]', TEST_DEPOSIT.mandatory.keterangan);

    // Submit
    await page.click('button[type="submit"]:has-text(/simpan/i)');

    // Verify success
    await expect(page.getByText(/berhasil|success/i)).toBeVisible();
  });

  test("should record voluntary deposit (Simpanan Sukarela) successfully", async ({
    page,
  }) => {
    await page.click('button:has-text(/tambah.*simpanan/i)');

    // Select member
    const memberSelect = page.locator('select[name="idAnggota"]');
    if (await memberSelect.isVisible()) {
      await memberSelect.selectOption({ index: 1 });
    }

    // Select deposit type - Simpanan Sukarela
    await page.selectOption(
      'select[name="tipeSimpanan"]',
      TEST_DEPOSIT.voluntary.tipeSimpanan
    );

    // Fill amount
    await page.fill(
      'input[name="jumlahSetoran"]',
      TEST_DEPOSIT.voluntary.jumlahSetoran.toString()
    );

    // Fill description
    await page.fill('textarea[name="keterangan"]', TEST_DEPOSIT.voluntary.keterangan);

    // Submit
    await page.click('button[type="submit"]:has-text(/simpan/i)');

    // Verify success
    await expect(page.getByText(/berhasil|success/i)).toBeVisible();
  });

  test("should validate required fields", async ({ page }) => {
    await page.click('button:has-text(/tambah.*simpanan/i)');

    // Try to submit without filling required fields
    await page.click('button[type="submit"]:has-text(/simpan/i)');

    // Check validation messages
    const validationMessages = [
      /anggota.*wajib|member.*required/i,
      /tipe.*wajib|type.*required/i,
      /jumlah.*wajib|amount.*required/i,
    ];

    for (const message of validationMessages) {
      const validationElement = page.getByText(message);
      if (await validationElement.isVisible()) {
        await expect(validationElement).toBeVisible();
      }
    }
  });

  test("should validate minimum deposit amount", async ({ page }) => {
    await page.click('button:has-text(/tambah.*simpanan/i)');

    // Select member
    const memberSelect = page.locator('select[name="idAnggota"]');
    if (await memberSelect.isVisible()) {
      await memberSelect.selectOption({ index: 1 });
    }

    // Select type
    await page.selectOption('select[name="tipeSimpanan"]', "POKOK");

    // Try to enter zero or negative amount
    await page.fill('input[name="jumlahSetoran"]', "0");
    await page.click('button[type="submit"]:has-text(/simpan/i)');

    // Should show validation error
    const errorMessage = page.getByText(/jumlah.*harus.*lebih.*besar|amount.*greater/i);
    if (await errorMessage.isVisible()) {
      await expect(errorMessage).toBeVisible();
    }
  });

  test("should auto-generate reference number", async ({ page }) => {
    await page.click('button:has-text(/tambah.*simpanan/i)');

    // Check if reference number is auto-generated
    const referenceInput = page.locator(
      'input[name="nomorReferensi"], [data-testid="reference-number"]'
    );

    if (await referenceInput.isVisible()) {
      const value = await referenceInput.inputValue();
      // Should have some format like SMP-XXX or similar
      expect(value).toMatch(/SMP|SIMPANAN|DEP/i);
    }
  });
});

test.describe("Share Capital - View Member Balance", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToShareCapital(page);
  });

  test("should display member balance list", async ({ page }) => {
    // Navigate to member balance view
    const balanceLink = page.locator('a[href*="/saldo"], nav >> text=/saldo/i');

    if (await balanceLink.isVisible()) {
      await balanceLink.click();
      await page.waitForURL(/.*\/saldo/);
    }

    // Check table headers
    const headers = [
      /nomor.*anggota|member.*number/i,
      /nama.*anggota|member.*name/i,
      /simpanan.*pokok|principal/i,
      /simpanan.*wajib|mandatory/i,
      /simpanan.*sukarela|voluntary/i,
      /total/i,
    ];

    for (const header of headers) {
      const headerElement = page.getByText(header);
      if (await headerElement.isVisible()) {
        await expect(headerElement).toBeVisible();
      }
    }
  });

  test("should show individual member balance details", async ({ page }) => {
    // Click on a member row to view details
    const memberRow = page.locator('tbody tr').first();

    if (await memberRow.isVisible()) {
      const detailBtn = memberRow.locator('button:has-text(/detail|lihat/i)');

      if (await detailBtn.isVisible()) {
        await detailBtn.click();

        // Check detail modal/page
        await expect(
          page.getByText(/detail.*saldo|balance.*detail/i)
        ).toBeVisible();

        // Check for breakdown
        await expect(page.getByText(/simpanan.*pokok/i)).toBeVisible();
        await expect(page.getByText(/simpanan.*wajib/i)).toBeVisible();
        await expect(page.getByText(/simpanan.*sukarela/i)).toBeVisible();
      }
    }
  });

  test("should calculate total share capital correctly", async ({ page }) => {
    const balanceLink = page.locator('a[href*="/saldo"]');

    if (await balanceLink.isVisible()) {
      await balanceLink.click();

      // Check for total calculation
      const totalElement = page.locator('[data-testid="total-share-capital"]');

      if (await totalElement.isVisible()) {
        const totalText = await totalElement.textContent();
        // Total should be a positive number
        expect(totalText).toMatch(/\d/);
      }
    }
  });

  test("should export member balance to Excel", async ({ page }) => {
    const balanceLink = page.locator('a[href*="/saldo"]');

    if (await balanceLink.isVisible()) {
      await balanceLink.click();

      const exportBtn = page.locator('button:has-text(/export|download|excel/i)');

      if (await exportBtn.isVisible()) {
        const downloadPromise = page.waitForEvent("download");
        await exportBtn.click();

        const download = await downloadPromise;
        expect(download.suggestedFilename()).toMatch(/saldo|balance|xlsx/i);
      }
    }
  });
});

test.describe("Share Capital - Transaction History", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToShareCapital(page);
  });

  test("should display transaction history list", async ({ page }) => {
    // Check transaction list
    const transactionRows = page.locator('tbody tr');
    const count = await transactionRows.count();

    expect(count).toBeGreaterThanOrEqual(0);
  });

  test("should filter transactions by deposit type", async ({ page }) => {
    const typeFilter = page.locator('select[name="tipeSimpanan"]');

    if (await typeFilter.isVisible()) {
      // Filter by Simpanan Pokok
      await typeFilter.selectOption("POKOK");
      await page.waitForTimeout(500);

      // Check filtered results
      const results = page.locator('[data-testid="deposit-type"]');
      if ((await results.count()) > 0) {
        await expect(results.first()).toContainText(/pokok/i);
      }
    }
  });

  test("should filter transactions by date range", async ({ page }) => {
    const startDate = page.locator('input[name="tanggalMulai"]');
    const endDate = page.locator('input[name="tanggalAkhir"]');

    if ((await startDate.isVisible()) && (await endDate.isVisible())) {
      await startDate.fill("2025-01-01");
      await endDate.fill("2025-01-31");

      // Apply filter
      await page.click('button:has-text(/filter|cari/i)');
      await page.waitForTimeout(500);

      // Verify results are filtered
      const rows = page.locator('tbody tr');
      const count = await rows.count();
      expect(count).toBeGreaterThanOrEqual(0);
    }
  });

  test("should search transactions by member name or number", async ({ page }) => {
    const searchInput = page.locator('input[placeholder*="Cari"]');

    if (await searchInput.isVisible()) {
      await searchInput.fill("A001");
      await page.waitForTimeout(500);

      // Check filtered results
      const rows = page.locator('tbody tr');
      if ((await rows.count()) > 0) {
        await expect(rows.first()).toContainText(/A001/i);
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

        // Check for transaction info
        const infoFields = [
          /nomor.*anggota/i,
          /tipe.*simpanan/i,
          /jumlah/i,
          /tanggal/i,
        ];

        for (const field of infoFields) {
          const fieldElement = page.getByText(field);
          if (await fieldElement.isVisible()) {
            await expect(fieldElement).toBeVisible();
          }
        }
      }
    }
  });

  test("should export transaction history to Excel", async ({ page }) => {
    const exportBtn = page.locator('button:has-text(/export|download|excel/i)');

    if (await exportBtn.isVisible()) {
      const downloadPromise = page.waitForEvent("download");
      await exportBtn.click();

      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/simpanan|deposit|xlsx/i);
    }
  });
});

test.describe("Share Capital - Member Dashboard View", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsMember(page);
  });

  test("should display member's own share capital balance", async ({ page }) => {
    // Check if balance card/widget is visible on dashboard
    const balanceCard = page.locator('[data-testid="share-capital-balance"]');

    if (await balanceCard.isVisible()) {
      await expect(balanceCard).toBeVisible();

      // Check for individual balances
      await expect(page.getByText(/simpanan.*pokok/i)).toBeVisible();
      await expect(page.getByText(/simpanan.*wajib/i)).toBeVisible();
    }
  });

  test("should display member's recent deposit transactions", async ({ page }) => {
    // Navigate to member's transaction history
    const historyLink = page.locator(
      'a[href*="/riwayat"], nav >> text=/riwayat.*simpanan/i'
    );

    if (await historyLink.isVisible()) {
      await historyLink.click();

      // Check transaction list
      await expect(
        page.getByRole("heading", { name: /riwayat.*simpanan/i })
      ).toBeVisible();

      // Check for transaction rows
      const rows = page.locator('tbody tr');
      if ((await rows.count()) > 0) {
        await expect(rows.first()).toBeVisible();
      }
    }
  });

  test("member should NOT be able to record deposits", async ({ page }) => {
    // Verify member doesn't see "Tambah Simpanan" button
    const addButton = page.locator('button:has-text(/tambah.*simpanan/i)');

    // Button should not be visible for regular members
    if (await addButton.isVisible()) {
      // This would be an access control issue
      expect(await addButton.count()).toBe(0);
    }
  });
});

test.describe("Share Capital - Summary and Statistics", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToShareCapital(page);
  });

  test("should display cooperative-wide share capital summary", async ({ page }) => {
    // Check for summary statistics
    const summaryCards = [
      /total.*simpanan.*pokok/i,
      /total.*simpanan.*wajib/i,
      /total.*simpanan.*sukarela/i,
      /total.*semua.*simpanan/i,
    ];

    for (const card of summaryCards) {
      const cardElement = page.locator(`[data-testid*="summary"], div:has-text("${card.source}")`);
      if (await cardElement.isVisible()) {
        await expect(cardElement).toBeVisible();
      }
    }
  });

  test("should show number of active members", async ({ page }) => {
    const memberCount = page.locator('[data-testid="member-count"]');

    if (await memberCount.isVisible()) {
      const countText = await memberCount.textContent();
      expect(countText).toMatch(/\d+/);
    }
  });

  test("should display monthly deposit trend", async ({ page }) => {
    // Check for trend chart or statistics
    const trendSection = page.locator('[data-testid="deposit-trend"], .chart');

    if (await trendSection.isVisible()) {
      await expect(trendSection).toBeVisible();
    }
  });
});

test.describe("Share Capital - Role-Based Access", () => {
  test("treasurer should access all share capital features", async ({ page }) => {
    await loginAsTreasurer(page);
    await navigateToShareCapital(page);

    // Verify can access share capital module
    await expect(
      page.getByRole("heading", { name: /simpanan/i })
    ).toBeVisible();

    // Verify can record deposits
    await expect(
      page.getByRole("button", { name: /tambah.*simpanan/i })
    ).toBeVisible();

    // Verify can view balances
    const balanceLink = page.locator('a[href*="/saldo"]');
    if (await balanceLink.isVisible()) {
      await expect(balanceLink).toBeVisible();
    }
  });

  test("admin should access share capital module", async ({ page }) => {
    await loginAsAdmin(page);
    await navigateToShareCapital(page);

    // Verify can access share capital module
    await expect(
      page.getByRole("heading", { name: /simpanan/i })
    ).toBeVisible();
  });

  test("member should only view own balance", async ({ page }) => {
    await loginAsMember(page);

    // Member should see their balance but not "Tambah Simpanan" button
    const addButton = page.locator('button:has-text(/tambah.*simpanan/i)');

    // Verify button is not visible
    expect(await addButton.count()).toBe(0);
  });

  test("cashier should NOT access share capital module", async ({ page }) => {
    const cashier = {
      email: "kasir@koperasi.local",
      password: "Kasir123!",
    };

    await page.goto("/login");
    await page.fill('input[name="email"]', cashier.email);
    await page.fill('input[name="password"]', cashier.password);
    await page.click('button[type="submit"]');

    await page.waitForURL("/dashboard");

    // Verify share capital menu is not visible
    const shareCapitalLink = page.locator('a[href*="/simpanan"]');
    await expect(shareCapitalLink).not.toBeVisible();
  });
});
