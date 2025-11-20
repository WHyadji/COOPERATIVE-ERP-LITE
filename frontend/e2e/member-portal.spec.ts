import { test, expect, type Page } from "@playwright/test";

// Test data - should match seeded data in backend
const TEST_MEMBER = {
  nomorAnggota: "A001",
  pin: "123456",
  namaLengkap: "Test Member Portal",
};

// Helper function to wait for React hydration
async function waitForHydration(page: Page) {
  // Wait for the page to be fully loaded
  await page.waitForLoadState('networkidle');

  // Wait for the form to be visible
  await page.waitForSelector('form', { state: 'visible' });

  // Give React generous time to hydrate and attach event handlers
  // In production builds, this can take longer than in development
  await page.waitForTimeout(3000);
}

// Helper function to login
async function loginMember(page: Page) {
  await page.goto("/portal/login");

  // Wait for React hydration before interacting with form
  await waitForHydration(page);

  // Fill login form
  await page.fill('input[name="nomorAnggota"]', TEST_MEMBER.nomorAnggota);
  await page.fill('input[name="pin"]', TEST_MEMBER.pin);

  // Submit form
  await page.click('button[type="submit"]');

  // Wait for navigation to dashboard
  await page.waitForURL("/portal");
}

test.describe("Member Portal - Login Flow", () => {
  test("should display login page correctly", async ({ page }) => {
    // Capture console logs and errors
    page.on('console', msg => console.log('BROWSER LOG:', msg.text()));
    page.on('pageerror', err => console.error('BROWSER ERROR:', err.message));

    await page.goto("/portal/login");

    // Check page title and elements
    await expect(page.getByRole("heading", { name: "Portal Anggota" })).toBeVisible();
    await expect(
      page.getByText("Masuk ke portal anggota koperasi")
    ).toBeVisible();

    // Check form fields
    await expect(page.getByLabel("Nomor Anggota")).toBeVisible();
    await expect(page.getByLabel(/PIN.*6 digit/i)).toBeVisible();
    await expect(page.getByRole("button", { name: /masuk/i })).toBeVisible();
  });

  test("should show validation errors for empty fields", async ({ page }) => {
    await page.goto("/portal/login");

    // Wait for React hydration
    await waitForHydration(page);

    // Click submit without filling fields
    await page.click('button[type="submit"]');

    // Wait for validation to trigger and check for validation messages
    await expect(page.getByText("Nomor anggota harus diisi", { exact: false })).toBeVisible();
    await expect(page.getByText(/PIN harus 6 digit/i)).toBeVisible();
  });

  test("should show error for invalid PIN format", async ({ page }) => {
    await page.goto("/portal/login");

    // Wait for React hydration
    await waitForHydration(page);

    await page.fill('input[name="nomorAnggota"]', "A001");
    await page.fill('input[name="pin"]', "12345"); // Only 5 digits

    await page.click('button[type="submit"]');

    await expect(page.getByText(/PIN harus 6 digit/i)).toBeVisible();
  });

  test("should login successfully with valid credentials", async ({ page }) => {
    await loginMember(page);

    // Should redirect to dashboard
    await expect(page).toHaveURL("/portal");

    // Should show welcome message
    await expect(
      page.getByText(`Selamat Datang, ${TEST_MEMBER.namaLengkap}`)
    ).toBeVisible();
  });

  test("should show error for invalid credentials", async ({ page }) => {
    await page.goto("/portal/login");

    // Wait for React hydration
    await waitForHydration(page);

    await page.fill('input[name="nomorAnggota"]', "A001");
    await page.fill('input[name="pin"]', "999999"); // Wrong PIN

    await page.click('button[type="submit"]');

    // Should show error message
    await expect(page.getByText(/Login gagal/i)).toBeVisible();
  });

  test("should toggle PIN visibility", async ({ page }) => {
    await page.goto("/portal/login");

    // Wait for React hydration
    await waitForHydration(page);

    const pinInput = page.locator('input[name="pin"]');

    // Initially should be password type
    await expect(pinInput).toHaveAttribute("type", "password");

    // Click visibility toggle
    await page.click('button[aria-label="toggle PIN visibility"]');

    // Wait for React to update and check type changed to text
    await expect(pinInput).toHaveAttribute("type", "text", { timeout: 2000 });

    // Click again to hide
    await page.click('button[aria-label="toggle PIN visibility"]');
    await expect(pinInput).toHaveAttribute("type", "password", { timeout: 2000 });
  });
});

test.describe("Member Portal - Dashboard", () => {
  test.beforeEach(async ({ page }) => {
    await loginMember(page);
  });

  test("should display dashboard with balance cards", async ({ page }) => {
    // Check welcome message
    await expect(
      page.getByText(`Selamat Datang, ${TEST_MEMBER.namaLengkap}`)
    ).toBeVisible();

    // Check balance cards
    await expect(page.getByText("Total Simpanan")).toBeVisible();
    await expect(page.getByText("Simpanan Pokok")).toBeVisible();
    await expect(page.getByText("Simpanan Wajib")).toBeVisible();
    await expect(page.getByText("Simpanan Sukarela")).toBeVisible();
  });

  test("should display recent transactions", async ({ page }) => {
    // Check recent transactions section
    await expect(page.getByText("Transaksi Terbaru")).toBeVisible();
    await expect(
      page.getByRole("button", { name: "Lihat Semua" })
    ).toBeVisible();
  });

  test("should navigate to balance page", async ({ page }) => {
    await page.click("text=Lihat Saldo Detail");

    await expect(page).toHaveURL("/portal/balance");
    await expect(page.getByText("Saldo Simpanan")).toBeVisible();
  });

  test("should navigate to transactions page", async ({ page }) => {
    await page.click("text=Riwayat Transaksi");

    await expect(page).toHaveURL("/portal/transactions");
    await expect(page.getByText("Riwayat Transaksi")).toBeVisible();
  });
});

test.describe("Member Portal - Balance Page", () => {
  test.beforeEach(async ({ page }) => {
    await loginMember(page);
    await page.goto("/portal/balance");
  });

  test("should display member information", async ({ page }) => {
    await expect(page.getByText("Nomor Anggota")).toBeVisible();
    await expect(page.getByText("Nama Anggota")).toBeVisible();
    await expect(page.getByText(TEST_MEMBER.nomorAnggota)).toBeVisible();
  });

  test("should display all balance types", async ({ page }) => {
    // Check total balance card
    await expect(page.getByText("Total Simpanan")).toBeVisible();

    // Check individual balance cards
    await expect(page.getByText("Simpanan Pokok")).toBeVisible();
    await expect(page.getByText("Simpanan Wajib")).toBeVisible();
    await expect(page.getByText("Simpanan Sukarela")).toBeVisible();
  });

  test("should display balance information", async ({ page }) => {
    // Check information sections
    await expect(page.getByText("Tentang Simpanan Pokok:")).toBeVisible();
    await expect(page.getByText("Tentang Simpanan Wajib:")).toBeVisible();
    await expect(page.getByText("Tentang Simpanan Sukarela:")).toBeVisible();

    // Check important information
    await expect(page.getByText("Informasi Penting")).toBeVisible();
  });
});

test.describe("Member Portal - Transactions Page", () => {
  test.beforeEach(async ({ page }) => {
    await loginMember(page);
    await page.goto("/portal/transactions");
  });

  test("should display transaction filters", async ({ page }) => {
    await expect(page.getByText("Filter")).toBeVisible();
    await expect(page.getByLabel("Tipe Simpanan")).toBeVisible();
    await expect(page.getByLabel("Tanggal Mulai")).toBeVisible();
    await expect(page.getByLabel("Tanggal Akhir")).toBeVisible();
  });

  test("should display transaction summary", async ({ page }) => {
    await expect(page.getByText("Total Transaksi")).toBeVisible();
    await expect(page.getByText("Total Setoran")).toBeVisible();
  });

  test("should filter transactions by type", async ({ page }) => {
    // Select Simpanan Pokok
    await page.click('label:has-text("Tipe Simpanan")');
    await page.click("text=Simpanan Pokok");

    // Click apply filter
    await page.click('button:has-text("Terapkan")');

    // Wait for results
    await page.waitForTimeout(1000);

    // Results should only show Simpanan Pokok
    const chips = await page.$$(
      '[class*="MuiChip"]:has-text("Simpanan Pokok")'
    );
    expect(chips.length).toBeGreaterThan(0);
  });

  test("should reset filters", async ({ page }) => {
    // Apply some filters first
    await page.click('label:has-text("Tipe Simpanan")');
    await page.click("text=Simpanan Wajib");
    await page.click('button:has-text("Terapkan")');

    await page.waitForTimeout(500);

    // Reset filters
    await page.click('button:has-text("Reset")');

    await page.waitForTimeout(500);

    // Filter should be back to default "Semua Tipe"
    const select = page.getByLabel("Tipe Simpanan");
    await expect(select).toHaveValue("all");
  });

  test("should display transaction table", async ({ page }) => {
    // Check table headers
    await expect(page.getByText("No. Referensi")).toBeVisible();
    await expect(page.getByText("Tanggal")).toBeVisible();
    await expect(page.getByText("Tipe Simpanan")).toBeVisible();
    await expect(page.getByText("Keterangan")).toBeVisible();
    await expect(page.getByText("Jumlah")).toBeVisible();
  });
});

test.describe("Member Portal - Profile Page", () => {
  test.beforeEach(async ({ page }) => {
    await loginMember(page);
    await page.goto("/portal/profile");
  });

  test("should display profile header", async ({ page }) => {
    await expect(page.getByText("Profil Saya")).toBeVisible();
    await expect(page.getByText(TEST_MEMBER.namaLengkap)).toBeVisible();
    await expect(
      page.getByText(`Nomor Anggota: ${TEST_MEMBER.nomorAnggota}`)
    ).toBeVisible();
  });

  test("should display personal information section", async ({ page }) => {
    await expect(page.getByText("Informasi Pribadi")).toBeVisible();
    await expect(page.getByLabel("NIK")).toBeVisible();
    await expect(page.getByLabel("Jenis Kelamin")).toBeVisible();
    await expect(page.getByLabel("Tanggal Lahir")).toBeVisible();
  });

  test("should display contact information section", async ({ page }) => {
    await expect(page.getByText("Informasi Kontak")).toBeVisible();
    await expect(page.getByLabel("Nomor Telepon")).toBeVisible();
    await expect(page.getByLabel("Email")).toBeVisible();
  });

  test("should display address section", async ({ page }) => {
    await expect(page.getByText("Alamat")).toBeVisible();
    await expect(page.getByLabel("Alamat Lengkap")).toBeVisible();
    await expect(page.getByLabel("RT")).toBeVisible();
    await expect(page.getByLabel("RW")).toBeVisible();
  });

  test("should enable edit mode", async ({ page }) => {
    // Click edit button
    await page.click('button:has-text("Edit Profil")');

    // Contact fields should be editable
    const phoneInput = page.getByLabel("Nomor Telepon");
    await expect(phoneInput).toBeEnabled();

    // Should show save and cancel buttons
    await expect(
      page.getByRole("button", { name: "Simpan Perubahan" })
    ).toBeVisible();
    await expect(page.getByRole("button", { name: "Batal" })).toBeVisible();
  });

  test("should cancel edit mode", async ({ page }) => {
    // Enter edit mode
    await page.click('button:has-text("Edit Profil")');

    // Click cancel
    await page.click('button:has-text("Batal")');

    // Should show edit button again
    await expect(
      page.getByRole("button", { name: "Edit Profil" })
    ).toBeVisible();

    // Fields should be disabled
    const phoneInput = page.getByLabel("Nomor Telepon");
    await expect(phoneInput).toBeDisabled();
  });
});

test.describe("Member Portal - Navigation", () => {
  test.beforeEach(async ({ page }) => {
    await loginMember(page);
  });

  test("should navigate using sidebar menu", async ({ page }) => {
    // Click Dashboard
    await page.click("text=Dashboard");
    await expect(page).toHaveURL("/portal");

    // Click Saldo Simpanan
    await page.click("text=Saldo Simpanan");
    await expect(page).toHaveURL("/portal/balance");

    // Click Riwayat Transaksi
    await page.click("text=Riwayat Transaksi");
    await expect(page).toHaveURL("/portal/transactions");

    // Click Profil Saya
    await page.click("text=Profil Saya");
    await expect(page).toHaveURL("/portal/profile");
  });

  test("should logout successfully", async ({ page }) => {
    // Click user menu
    await page
      .click('[aria-label="account of current user"]', { timeout: 5000 })
      .catch(() => {
        // If aria-label doesn't work, try clicking the avatar
        return page.click('button:has([class*="MuiAvatar"])');
      });

    // Click logout
    await page.click("text=Keluar");

    // Should redirect to login
    await expect(page).toHaveURL("/portal/login");
  });

  test("should show mobile menu on small screens", async ({ page }) => {
    // Set viewport to mobile
    await page.setViewportSize({ width: 375, height: 667 });

    await page.goto("/portal");

    // Menu button should be visible
    const menuButton = page
      .locator('[aria-label="open drawer"]')
      .or(page.locator('button:has([data-testid="MenuIcon"])'));
    await expect(menuButton).toBeVisible();
  });
});

test.describe("Member Portal - Responsive Design", () => {
  test.beforeEach(async ({ page }) => {
    await loginMember(page);
  });

  test("should be responsive on mobile", async ({ page }) => {
    await page.setViewportSize({ width: 375, height: 667 });

    // Dashboard should still display
    await expect(page.getByText("Total Simpanan")).toBeVisible();

    // Balance cards should stack vertically (grid behavior)
    const cards = await page.$$('[class*="MuiCard"]');
    expect(cards.length).toBeGreaterThan(0);
  });

  test("should be responsive on tablet", async ({ page }) => {
    await page.setViewportSize({ width: 768, height: 1024 });

    await expect(page.getByText("Total Simpanan")).toBeVisible();
  });

  test("should be responsive on desktop", async ({ page }) => {
    await page.setViewportSize({ width: 1920, height: 1080 });

    await expect(page.getByText("Total Simpanan")).toBeVisible();

    // Sidebar should be visible permanently
    await expect(page.getByText("Portal Anggota")).toBeVisible();
  });
});
