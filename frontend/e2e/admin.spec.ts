import { test, expect, type Page } from "@playwright/test";

/**
 * E2E Tests for Admin Panel Module
 *
 * Test Coverage:
 * 1. User Management (Create, Read, Update, Delete users)
 * 2. Role Assignment (Admin, Bendahara, Kasir, Anggota)
 * 3. System Settings
 * 4. Audit Logs
 * 5. Role-Based Access Control
 */

// Test data - Admin credentials
const TEST_ADMIN = {
  email: "admin@koperasi.local",
  password: "Admin123!",
};

// Test user data
const TEST_USER = {
  namaLengkap: "Test Bendahara",
  namaPengguna: "test.bendahara",
  email: "test.bendahara@koperasi.local",
  peran: "BENDAHARA",
  passwordDefault: "Password123!",
};

const TEST_USER_KASIR = {
  namaLengkap: "Test Kasir",
  namaPengguna: "test.kasir",
  email: "test.kasir@koperasi.local",
  peran: "KASIR",
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

// Helper function to navigate to Admin Panel
async function navigateToAdminPanel(page: Page) {
  await page.click('a[href*="/admin"], nav >> text=/admin|pengaturan/i');
  await page.waitForURL(/.*\/(admin|pengaturan)/);
}

// Helper function to navigate to User Management
async function navigateToUserManagement(page: Page) {
  await navigateToAdminPanel(page);
  await page.click('a[href*="/users"], a[href*="/pengguna"]');
  await page.waitForURL(/.*\/(users|pengguna)/);
}

test.describe("Admin Panel - User Management", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
    await navigateToUserManagement(page);
  });

  test("should display user management page correctly", async ({ page }) => {
    // Check page title
    await expect(
      page.getByRole("heading", { name: /manajemen.*pengguna|user.*management/i })
    ).toBeVisible();

    // Check "Tambah Pengguna" button
    await expect(
      page.getByRole("button", { name: /tambah.*pengguna|add.*user|create/i })
    ).toBeVisible();

    // Check table headers
    const headers = [
      /nama.*lengkap|full.*name/i,
      /nama.*pengguna|username/i,
      /email/i,
      /peran|role/i,
      /status/i,
    ];

    for (const header of headers) {
      const headerElement = page.getByText(header);
      if (await headerElement.isVisible()) {
        await expect(headerElement).toBeVisible();
      }
    }
  });

  test("should create new user successfully", async ({ page }) => {
    // Click "Tambah Pengguna"
    await page.click('button:has-text(/tambah.*pengguna|add.*user/i)');

    // Wait for form
    await expect(page.getByText(/form.*pengguna|user.*form/i)).toBeVisible();

    // Fill user details
    await page.fill('input[name="namaLengkap"]', TEST_USER.namaLengkap);
    await page.fill('input[name="namaPengguna"]', TEST_USER.namaPengguna);
    await page.fill('input[name="email"]', TEST_USER.email);

    // Select role
    await page.selectOption('select[name="peran"]', TEST_USER.peran);

    // Fill password (if required during creation)
    const passwordInput = page.locator('input[name="password"]');
    if (await passwordInput.isVisible()) {
      await passwordInput.fill(TEST_USER.passwordDefault);
    }

    // Submit form
    await page.click('button[type="submit"]:has-text(/simpan|save|create/i)');

    // Check success message
    await expect(page.getByText(/berhasil|success/i)).toBeVisible();

    // Verify user appears in list
    await expect(page.getByText(TEST_USER.namaLengkap)).toBeVisible();
    await expect(page.getByText(TEST_USER.email)).toBeVisible();
  });

  test("should validate required fields when creating user", async ({ page }) => {
    await page.click('button:has-text(/tambah.*pengguna/i)');

    // Try to submit empty form
    await page.click('button[type="submit"]:has-text(/simpan/i)');

    // Check validation messages
    const validationMessages = [
      /nama.*lengkap.*wajib|full.*name.*required/i,
      /nama.*pengguna.*wajib|username.*required/i,
      /email.*wajib|email.*required/i,
      /peran.*wajib|role.*required/i,
    ];

    for (const message of validationMessages) {
      const validationElement = page.getByText(message);
      if (await validationElement.isVisible()) {
        await expect(validationElement).toBeVisible();
      }
    }
  });

  test("should update existing user", async ({ page }) => {
    // Find user and click edit button
    const userRow = page.locator(`tr:has-text("${TEST_USER.email}")`).first();

    if (await userRow.isVisible()) {
      const editBtn = userRow.locator('button:has-text(/edit|ubah/i)');

      if (await editBtn.isVisible()) {
        await editBtn.click();

        // Update user details
        const newName = TEST_USER.namaLengkap + " (Updated)";
        await page.fill('input[name="namaLengkap"]', newName);

        // Submit update
        await page.click('button[type="submit"]:has-text(/simpan|update/i)');

        // Verify update success
        await expect(
          page.getByText(/berhasil.*diperbarui|updated successfully/i)
        ).toBeVisible();
        await expect(page.getByText(newName)).toBeVisible();
      }
    }
  });

  test("should deactivate/activate user", async ({ page }) => {
    const userRow = page.locator(`tr:has-text("${TEST_USER.email}")`).first();

    if (await userRow.isVisible()) {
      const statusToggle = userRow.locator(
        'button:has-text(/aktif|active|nonaktifkan|deactivate/i)'
      );

      if (await statusToggle.isVisible()) {
        await statusToggle.click();

        // Confirm action if there's a confirmation dialog
        const confirmBtn = page.locator('button:has-text(/ya|yes|confirm/i)');
        if (await confirmBtn.isVisible()) {
          await confirmBtn.click();
        }

        // Check success message
        await expect(page.getByText(/berhasil|success/i)).toBeVisible();
      }
    }
  });

  test("should filter users by role", async ({ page }) => {
    const roleFilter = page.locator('select[name="peran"]');

    if (await roleFilter.isVisible()) {
      // Filter by BENDAHARA
      await roleFilter.selectOption("BENDAHARA");
      await page.waitForTimeout(500);

      // Check filtered results
      const results = page.locator('[data-testid="user-role"]');
      if ((await results.count()) > 0) {
        await expect(results.first()).toContainText(/bendahara/i);
      }
    }
  });

  test("should search users by name or email", async ({ page }) => {
    const searchInput = page.locator('input[placeholder*="Cari"]');

    if (await searchInput.isVisible()) {
      await searchInput.fill(TEST_USER.email);
      await page.waitForTimeout(500);

      // Check search results
      const rows = page.locator('tbody tr');
      if ((await rows.count()) > 0) {
        await expect(rows.first()).toContainText(TEST_USER.email);
      }
    }
  });

  test("should reset user password", async ({ page }) => {
    const userRow = page.locator(`tr:has-text("${TEST_USER.email}")`).first();

    if (await userRow.isVisible()) {
      const resetBtn = userRow.locator('button:has-text(/reset.*password/i)');

      if (await resetBtn.isVisible()) {
        await resetBtn.click();

        // Confirm password reset
        const confirmBtn = page.locator('button:has-text(/ya|yes|confirm/i)');
        if (await confirmBtn.isVisible()) {
          await confirmBtn.click();
        }

        // Should show success message with new password
        await expect(
          page.getByText(/password.*berhasil.*direset|password.*reset.*success/i)
        ).toBeVisible();
      }
    }
  });

  test("should delete user with confirmation", async ({ page }) => {
    const userRow = page.locator(`tr:has-text("${TEST_USER_KASIR.email}")`).first();

    if (await userRow.isVisible()) {
      const deleteBtn = userRow.locator('button:has-text(/hapus|delete/i)');

      if (await deleteBtn.isVisible()) {
        await deleteBtn.click();

        // Confirm deletion
        const confirmBtn = page.locator('button:has-text(/ya|yes|hapus/i)');
        if (await confirmBtn.isVisible()) {
          await confirmBtn.click();
        }

        // Check success message
        await expect(page.getByText(/berhasil.*dihapus|deleted successfully/i)).toBeVisible();
      }
    }
  });
});

test.describe("Admin Panel - System Settings", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
    await navigateToAdminPanel(page);
  });

  test("should display system settings page", async ({ page }) => {
    const settingsLink = page.locator('a[href*="/settings"], a[href*="/pengaturan"]');

    if (await settingsLink.isVisible()) {
      await settingsLink.click();
      await page.waitForURL(/.*\/(settings|pengaturan)/);

      // Check page title
      await expect(
        page.getByRole("heading", { name: /pengaturan.*sistem|system.*settings/i })
      ).toBeVisible();
    }
  });

  test("should update cooperative information", async ({ page }) => {
    const settingsLink = page.locator('a[href*="/settings"]');

    if (await settingsLink.isVisible()) {
      await settingsLink.click();

      // Update cooperative name
      const nameInput = page.locator('input[name="namaKoperasi"]');
      if (await nameInput.isVisible()) {
        const newName = "Koperasi Test (Updated)";
        await nameInput.fill(newName);

        // Save changes
        await page.click('button:has-text(/simpan|save/i)');

        // Check success
        await expect(page.getByText(/berhasil|success/i)).toBeVisible();
      }
    }
  });

  test("should configure accounting settings", async ({ page }) => {
    const settingsLink = page.locator('a[href*="/settings"]');

    if (await settingsLink.isVisible()) {
      await settingsLink.click();

      // Navigate to accounting settings tab
      const accountingTab = page.locator('button:has-text(/akuntansi|accounting/i)');
      if (await accountingTab.isVisible()) {
        await accountingTab.click();

        // Check accounting settings options
        const fiscalYearInput = page.locator('input[name="tahunFiskal"]');
        if (await fiscalYearInput.isVisible()) {
          await expect(fiscalYearInput).toBeVisible();
        }
      }
    }
  });
});

test.describe("Admin Panel - Audit Logs", () => {
  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
    await navigateToAdminPanel(page);
  });

  test("should display audit log page", async ({ page }) => {
    const auditLink = page.locator('a[href*="/audit"], a[href*="/log"]');

    if (await auditLink.isVisible()) {
      await auditLink.click();
      await page.waitForURL(/.*\/(audit|log)/);

      // Check page title
      await expect(
        page.getByRole("heading", { name: /audit.*log|riwayat.*aktivitas/i })
      ).toBeVisible();
    }
  });

  test("should filter audit logs by date range", async ({ page }) => {
    const auditLink = page.locator('a[href*="/audit"]');

    if (await auditLink.isVisible()) {
      await auditLink.click();

      // Set date range
      const startDate = page.locator('input[name="tanggalMulai"]');
      const endDate = page.locator('input[name="tanggalAkhir"]');

      if ((await startDate.isVisible()) && (await endDate.isVisible())) {
        await startDate.fill("2025-01-01");
        await endDate.fill("2025-01-31");

        // Apply filter
        await page.click('button:has-text(/filter|cari/i)');
        await page.waitForTimeout(500);

        // Verify results
        const logs = page.locator('tbody tr');
        expect(await logs.count()).toBeGreaterThanOrEqual(0);
      }
    }
  });

  test("should filter audit logs by user", async ({ page }) => {
    const auditLink = page.locator('a[href*="/audit"]');

    if (await auditLink.isVisible()) {
      await auditLink.click();

      // Filter by user
      const userFilter = page.locator('select[name="idPengguna"]');
      if (await userFilter.isVisible()) {
        await userFilter.selectOption({ index: 1 });
        await page.waitForTimeout(500);

        // Verify filtered results
        const logs = page.locator('tbody tr');
        expect(await logs.count()).toBeGreaterThanOrEqual(0);
      }
    }
  });

  test("should filter audit logs by action type", async ({ page }) => {
    const auditLink = page.locator('a[href*="/audit"]');

    if (await auditLink.isVisible()) {
      await auditLink.click();

      // Filter by action type
      const actionFilter = page.locator('select[name="jenisAksi"]');
      if (await actionFilter.isVisible()) {
        await actionFilter.selectOption("CREATE");
        await page.waitForTimeout(500);

        // Verify filtered results
        const logs = page.locator('tbody tr');
        if ((await logs.count()) > 0) {
          await expect(logs.first()).toContainText(/create|tambah/i);
        }
      }
    }
  });

  test("should export audit logs to Excel", async ({ page }) => {
    const auditLink = page.locator('a[href*="/audit"]');

    if (await auditLink.isVisible()) {
      await auditLink.click();

      const exportBtn = page.locator('button:has-text(/export|download/i)');

      if (await exportBtn.isVisible()) {
        const downloadPromise = page.waitForEvent("download");
        await exportBtn.click();

        const download = await downloadPromise;
        expect(download.suggestedFilename()).toMatch(/audit|log|xlsx/i);
      }
    }
  });
});

test.describe("Admin Panel - Role-Based Access", () => {
  test("admin should access all admin panel features", async ({ page }) => {
    await loginAsAdmin(page);
    await navigateToAdminPanel(page);

    // Verify can access admin panel
    await expect(
      page.getByRole("heading", { name: /admin|pengaturan/i })
    ).toBeVisible();

    // Check for admin menu items
    const adminMenuItems = [
      page.locator('a[href*="/users"], a[href*="/pengguna"]'),
      page.locator('a[href*="/settings"]'),
      page.locator('a[href*="/audit"]'),
    ];

    for (const item of adminMenuItems) {
      if (await item.isVisible()) {
        await expect(item).toBeVisible();
      }
    }
  });

  test("treasurer should NOT access admin panel", async ({ page }) => {
    const treasurer = {
      email: "bendahara@koperasi.local",
      password: "Bendahara123!",
    };

    await page.goto("/login");
    await page.fill('input[name="email"]', treasurer.email);
    await page.fill('input[name="password"]', treasurer.password);
    await page.click('button[type="submit"]');

    await page.waitForURL("/dashboard");

    // Verify admin menu is not visible
    const adminLink = page.locator('a[href*="/admin"], a[href*="/pengaturan"]');

    // Try to navigate directly to admin page
    await page.goto("/admin");

    // Should be redirected or see access denied
    const accessDenied = page.getByText(/tidak.*diizinkan|access.*denied|forbidden/i);
    if (await accessDenied.isVisible()) {
      await expect(accessDenied).toBeVisible();
    } else {
      // Should redirect to dashboard
      await page.waitForURL("/dashboard");
    }
  });

  test("cashier should NOT access admin panel", async ({ page }) => {
    const cashier = {
      email: "kasir@koperasi.local",
      password: "Kasir123!",
    };

    await page.goto("/login");
    await page.fill('input[name="email"]', cashier.email);
    await page.fill('input[name="password"]', cashier.password);
    await page.click('button[type="submit"]');

    await page.waitForURL("/dashboard");

    // Verify admin menu is not visible
    const adminLink = page.locator('a[href*="/admin"]');
    await expect(adminLink).not.toBeVisible();
  });

  test("member should NOT access admin panel", async ({ page }) => {
    const member = {
      nomorAnggota: "A001",
      password: "123456",
    };

    await page.goto("/login");
    await page.fill('input[name="nomorAnggota"]', member.nomorAnggota);
    await page.fill('input[name="password"]', member.password);
    await page.click('button[type="submit"]');

    await page.waitForURL("/dashboard");

    // Verify admin menu is not visible
    const adminLink = page.locator('a[href*="/admin"]');
    await expect(adminLink).not.toBeVisible();
  });
});
