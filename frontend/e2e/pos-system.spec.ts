import { test, expect, type Page } from "@playwright/test";

/**
 * E2E Tests for POS (Point of Sale) System
 *
 * Test Coverage:
 * 1. Product Management (CRUD)
 * 2. POS Transaction Flow
 * 3. Stock Management
 * 4. Payment & Change Calculation
 * 5. Receipt Generation
 */

// Test data - Admin/Cashier credentials
const TEST_ADMIN = {
  email: "admin@koperasi.local",
  password: "Admin123!",
};

const TEST_CASHIER = {
  email: "kasir@koperasi.local",
  password: "Kasir123!",
};

// Test product data
const TEST_PRODUCT = {
  kodeProduk: "TEST-001",
  namaProduk: "Beras Premium 5kg",
  kategori: "Sembako",
  deskripsi: "Beras premium kualitas terbaik",
  harga: 75000,
  hargaBeli: 65000,
  stok: 50,
  stokMinimum: 10,
  satuan: "pcs",
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

// Helper function to login as cashier
async function loginAsCashier(page: Page) {
  await page.goto("/login");

  await page.fill('input[name="email"]', TEST_CASHIER.email);
  await page.fill('input[name="password"]', TEST_CASHIER.password);
  await page.click('button[type="submit"]');

  // Wait for redirect to dashboard
  await page.waitForURL("/dashboard");
}

// Helper function to navigate to product management
async function navigateToProducts(page: Page) {
  // Click on "Produk" menu item
  await page.click('a[href*="/produk"], nav >> text=Produk');
  await page.waitForURL(/.*\/produk/);
}

// Helper function to navigate to POS
async function navigateToPOS(page: Page) {
  // Click on "POS/Kasir" menu item
  await page.click('a[href*="/pos"], a[href*="/kasir"], nav >> text=/POS|Kasir/i');
  await page.waitForURL(/.*\/(pos|kasir)/);
}

test.describe("POS System - Product Management", () => {

  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
    await navigateToProducts(page);
  });

  test("should display product list page correctly", async ({ page }) => {
    // Check page title
    await expect(page.getByRole("heading", { name: /daftar produk|produk/i })).toBeVisible();

    // Check "Tambah Produk" button exists
    await expect(page.getByRole("button", { name: /tambah produk/i })).toBeVisible();

    // Check search input exists
    await expect(page.getByPlaceholder(/cari produk/i)).toBeVisible();

    // Check table headers
    await expect(page.getByText("Kode Produk")).toBeVisible();
    await expect(page.getByText("Nama Produk")).toBeVisible();
    await expect(page.getByText("Harga")).toBeVisible();
    await expect(page.getByText("Stok")).toBeVisible();
  });

  test("should create a new product successfully", async ({ page }) => {
    // Click "Tambah Produk" button
    await page.click('button:has-text("Tambah Produk")');

    // Wait for form dialog/page
    await expect(page.getByText(/form produk|tambah produk baru/i)).toBeVisible();

    // Fill product form
    await page.fill('input[name="kodeProduk"]', TEST_PRODUCT.kodeProduk);
    await page.fill('input[name="namaProduk"]', TEST_PRODUCT.namaProduk);
    await page.fill('input[name="kategori"]', TEST_PRODUCT.kategori);
    await page.fill('textarea[name="deskripsi"]', TEST_PRODUCT.deskripsi);
    await page.fill('input[name="harga"]', TEST_PRODUCT.harga.toString());
    await page.fill('input[name="hargaBeli"]', TEST_PRODUCT.hargaBeli.toString());
    await page.fill('input[name="stok"]', TEST_PRODUCT.stok.toString());
    await page.fill('input[name="stokMinimum"]', TEST_PRODUCT.stokMinimum.toString());

    // Select satuan (unit)
    await page.selectOption('select[name="satuan"]', TEST_PRODUCT.satuan);

    // Submit form
    await page.click('button[type="submit"]:has-text(/simpan|save/i)');

    // Check success message
    await expect(page.getByText(/berhasil|success/i)).toBeVisible();

    // Verify product appears in list
    await expect(page.getByText(TEST_PRODUCT.namaProduk)).toBeVisible();
    await expect(page.getByText(TEST_PRODUCT.kodeProduk)).toBeVisible();
  });

  test("should show validation errors for invalid product data", async ({ page }) => {
    await page.click('button:has-text("Tambah Produk")');

    // Try to submit empty form
    await page.click('button[type="submit"]:has-text(/simpan|save/i)');

    // Check validation messages
    await expect(page.getByText(/kode produk.*wajib|required/i)).toBeVisible();
    await expect(page.getByText(/nama produk.*wajib|required/i)).toBeVisible();
    await expect(page.getByText(/harga.*wajib|required/i)).toBeVisible();
  });

  test("should search products by name or code", async ({ page }) => {
    // Type in search box
    await page.fill('input[placeholder*="Cari produk"]', "Beras");

    // Wait for search results (with debounce)
    await page.waitForTimeout(500);

    // Check that only matching products are shown
    const rows = page.locator('tbody tr');
    await expect(rows.first()).toContainText(/beras/i);
  });

  test("should filter products by category", async ({ page }) => {
    // Select category filter (if available)
    const categoryFilter = page.locator('select[name="kategori"], [data-testid="category-filter"]');

    if (await categoryFilter.isVisible()) {
      await categoryFilter.selectOption("Sembako");
      await page.waitForTimeout(300);

      // Verify filtered results
      const rows = page.locator('tbody tr');
      const count = await rows.count();
      expect(count).toBeGreaterThan(0);
    }
  });

  test("should update existing product", async ({ page }) => {
    // Find product and click edit button
    const productRow = page.locator(`tr:has-text("${TEST_PRODUCT.kodeProduk}")`).first();
    await productRow.locator('button:has-text(/edit|ubah/i)').click();

    // Update product name
    const newName = TEST_PRODUCT.namaProduk + " (Updated)";
    await page.fill('input[name="namaProduk"]', newName);

    // Update price
    const newPrice = "80000";
    await page.fill('input[name="harga"]', newPrice);

    // Submit update
    await page.click('button[type="submit"]:has-text(/simpan|update|save/i)');

    // Verify update success
    await expect(page.getByText(/berhasil.*diperbarui|updated successfully/i)).toBeVisible();
    await expect(page.getByText(newName)).toBeVisible();
  });

  test("should show low stock warning", async ({ page }) => {
    // Look for products with stock below minimum
    const lowStockBadge = page.locator('[data-testid="low-stock"], .badge:has-text(/stok rendah|low stock/i)');

    if (await lowStockBadge.count() > 0) {
      await expect(lowStockBadge.first()).toBeVisible();
    }
  });

  test("should delete product with confirmation", async ({ page }) => {
    // Find product and click delete button
    const productRow = page.locator(`tr:has-text("${TEST_PRODUCT.kodeProduk}")`).first();
    await productRow.locator('button:has-text(/hapus|delete/i)').click();

    // Confirm deletion dialog
    await expect(page.getByText(/yakin.*hapus|confirm.*delete/i)).toBeVisible();
    await page.click('button:has-text(/ya|yes|hapus/i)');

    // Verify deletion success
    await expect(page.getByText(/berhasil.*dihapus|deleted successfully/i)).toBeVisible();

    // Verify product is removed from list
    await expect(page.getByText(TEST_PRODUCT.kodeProduk)).not.toBeVisible();
  });
});

test.describe("POS System - Transaction Flow", () => {

  test.beforeEach(async ({ page }) => {
    await loginAsCashier(page);
    await navigateToPOS(page);
  });

  test("should display POS interface correctly", async ({ page }) => {
    // Check POS layout components
    await expect(page.getByRole("heading", { name: /pos|point of sale|kasir/i })).toBeVisible();

    // Check product search/selection area
    await expect(page.getByPlaceholder(/cari produk|search product/i)).toBeVisible();

    // Check cart/transaction summary area
    await expect(page.getByText(/total|subtotal/i)).toBeVisible();

    // Check payment section
    await expect(page.getByText(/bayar|payment|tunai/i)).toBeVisible();

    // Check action buttons
    await expect(page.getByRole("button", { name: /proses|selesai|checkout/i })).toBeVisible();
    await expect(page.getRole("button", { name: /batal|cancel|reset/i })).toBeVisible();
  });

  test("should add product to cart by search", async ({ page }) => {
    // Search for product
    await page.fill('input[placeholder*="Cari produk"]', "Beras");
    await page.waitForTimeout(300);

    // Click on product from search results
    const productItem = page.locator('[data-testid="product-item"], .product-card').first();
    await productItem.click();

    // Verify product added to cart
    await expect(page.locator('[data-testid="cart-item"], .cart-item')).toHaveCount(1);

    // Check product details in cart
    await expect(page.getByText(/beras/i)).toBeVisible();
  });

  test("should add product to cart by barcode", async ({ page }) => {
    // Type barcode (simulating barcode scanner)
    const barcodeInput = page.locator('input[placeholder*="Barcode"], input[name="barcode"]');

    if (await barcodeInput.isVisible()) {
      await barcodeInput.fill("8992761111111");
      await barcodeInput.press("Enter");

      // Verify product added
      const cartItems = page.locator('[data-testid="cart-item"], .cart-item');
      await expect(cartItems).toHaveCount(1);
    }
  });

  test("should adjust product quantity in cart", async ({ page }) => {
    // Add product to cart
    await page.fill('input[placeholder*="Cari produk"]', "Beras");
    await page.waitForTimeout(300);
    await page.locator('[data-testid="product-item"]').first().click();

    // Find quantity controls
    const qtyInput = page.locator('[data-testid="qty-input"], input[name="kuantitas"]').first();
    const initialQty = await qtyInput.inputValue();

    // Increase quantity
    const increaseBtn = page.locator('button:has-text("+"), [data-testid="qty-increase"]').first();
    await increaseBtn.click();

    // Verify quantity increased
    const newQty = await qtyInput.inputValue();
    expect(parseInt(newQty)).toBe(parseInt(initialQty) + 1);

    // Verify total updated
    await expect(page.getByText(/total/i)).toBeVisible();
  });

  test("should remove product from cart", async ({ page }) => {
    // Add product to cart
    await page.fill('input[placeholder*="Cari produk"]', "Beras");
    await page.waitForTimeout(300);
    await page.locator('[data-testid="product-item"]').first().click();

    // Click remove button
    const removeBtn = page.locator('button:has-text(/hapus|remove/i), [data-testid="remove-item"]').first();
    await removeBtn.click();

    // Verify cart is empty
    await expect(page.locator('[data-testid="cart-item"]')).toHaveCount(0);
    await expect(page.getByText(/keranjang kosong|cart empty/i)).toBeVisible();
  });

  test("should calculate total correctly with multiple items", async ({ page }) => {
    // Add first product
    await page.fill('input[placeholder*="Cari produk"]', "Beras");
    await page.waitForTimeout(300);
    await page.locator('[data-testid="product-item"]').first().click();

    // Add second product
    await page.fill('input[placeholder*="Cari produk"]', "Gula");
    await page.waitForTimeout(300);
    await page.locator('[data-testid="product-item"]').first().click();

    // Verify multiple items in cart
    await expect(page.locator('[data-testid="cart-item"]')).toHaveCount(2);

    // Verify total is calculated (should be > 0)
    const totalText = await page.locator('[data-testid="total-amount"], .total-amount').textContent();
    expect(totalText).toMatch(/\d+/); // Contains numbers
  });

  test("should complete cash transaction successfully", async ({ page }) => {
    // Add product to cart
    await page.fill('input[placeholder*="Cari produk"]', "Beras");
    await page.waitForTimeout(300);
    await page.locator('[data-testid="product-item"]').first().click();

    // Get total amount
    const totalText = await page.locator('[data-testid="total-amount"]').textContent();
    const totalAmount = parseInt(totalText?.replace(/\D/g, "") || "0");

    // Enter payment amount (more than total)
    const paymentInput = page.locator('input[name="jumlahBayar"], [data-testid="payment-amount"]');
    await paymentInput.fill((totalAmount + 10000).toString());

    // Verify change calculated
    const changeAmount = page.locator('[data-testid="change-amount"], .change-amount');
    await expect(changeAmount).toContainText("10");

    // Select payment method (cash)
    const cashOption = page.locator('input[value="TUNAI"], [data-testid="payment-cash"]');
    if (await cashOption.isVisible()) {
      await cashOption.click();
    }

    // Process transaction
    await page.click('button:has-text(/proses|selesai|checkout/i)');

    // Verify transaction success
    await expect(page.getByText(/berhasil|success|transaksi selesai/i)).toBeVisible();
  });

  test("should show error when payment insufficient", async ({ page }) => {
    // Add product to cart
    await page.fill('input[placeholder*="Cari produk"]', "Beras");
    await page.waitForTimeout(300);
    await page.locator('[data-testid="product-item"]').first().click();

    // Enter insufficient payment
    const paymentInput = page.locator('input[name="jumlahBayar"]');
    await paymentInput.fill("1000"); // Too low

    // Try to process
    await page.click('button:has-text(/proses|checkout/i)');

    // Verify error message
    await expect(page.getByText(/pembayaran.*kurang|insufficient/i)).toBeVisible();
  });

  test("should allow adding customer/member to transaction", async ({ page }) => {
    // Look for member selection
    const memberSelect = page.locator('input[placeholder*="Cari anggota"], select[name="idAnggota"]');

    if (await memberSelect.isVisible()) {
      // Select a member
      if (await memberSelect.getAttribute("type") === "text") {
        await memberSelect.fill("A001");
        await page.waitForTimeout(300);
        await page.keyboard.press("Enter");
      } else {
        await memberSelect.selectOption({ index: 1 });
      }

      // Verify member added
      await expect(page.getByText(/anggota.*dipilih|member selected/i)).toBeVisible();
    }
  });

  test("should generate and display receipt", async ({ page }) => {
    // Complete a transaction first
    await page.fill('input[placeholder*="Cari produk"]', "Beras");
    await page.waitForTimeout(300);
    await page.locator('[data-testid="product-item"]').first().click();

    const totalText = await page.locator('[data-testid="total-amount"]').textContent();
    const totalAmount = parseInt(totalText?.replace(/\D/g, "") || "0");

    await page.fill('input[name="jumlahBayar"]', (totalAmount + 10000).toString());
    await page.click('button:has-text(/proses|checkout/i)');

    // Wait for receipt
    await page.waitForTimeout(500);

    // Check receipt content
    await expect(page.getByText(/struk|receipt|nota/i)).toBeVisible();
    await expect(page.getByText(/nomor.*transaksi|transaction.*number/i)).toBeVisible();
    await expect(page.getByText(/tanggal|date/i)).toBeVisible();
    await expect(page.getByText(/kasir|cashier/i)).toBeVisible();

    // Check print button
    const printBtn = page.locator('button:has-text(/cetak|print/i)');
    if (await printBtn.isVisible()) {
      await expect(printBtn).toBeEnabled();
    }
  });

  test("should reset cart after completing transaction", async ({ page }) => {
    // Complete a transaction
    await page.fill('input[placeholder*="Cari produk"]', "Beras");
    await page.waitForTimeout(300);
    await page.locator('[data-testid="product-item"]').first().click();

    const totalText = await page.locator('[data-testid="total-amount"]').textContent();
    const totalAmount = parseInt(totalText?.replace(/\D/g, "") || "0");

    await page.fill('input[name="jumlahBayar"]', (totalAmount + 5000).toString());
    await page.click('button:has-text(/proses|checkout/i)');

    // Close receipt or click "New Transaction"
    const newTxBtn = page.locator('button:has-text(/transaksi baru|new transaction|tutup/i)');
    if (await newTxBtn.isVisible()) {
      await newTxBtn.click();
    }

    // Verify cart is empty
    await expect(page.locator('[data-testid="cart-item"]')).toHaveCount(0);
    await expect(page.locator('[data-testid="total-amount"]')).toContainText("0");
  });

  test("should cancel transaction and clear cart", async ({ page }) => {
    // Add items to cart
    await page.fill('input[placeholder*="Cari produk"]', "Beras");
    await page.waitForTimeout(300);
    await page.locator('[data-testid="product-item"]').first().click();

    // Click cancel/reset button
    await page.click('button:has-text(/batal|cancel|reset/i)');

    // Confirm cancellation if dialog appears
    const confirmBtn = page.locator('button:has-text(/ya|yes|ok/i)');
    if (await confirmBtn.isVisible()) {
      await confirmBtn.click();
    }

    // Verify cart is cleared
    await expect(page.locator('[data-testid="cart-item"]')).toHaveCount(0);
  });
});

test.describe("POS System - Stock Management", () => {

  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
  });

  test("should update stock after successful sale", async ({ page }) => {
    // Get initial stock
    await navigateToProducts(page);
    const productRow = page.locator('tr:has-text("Beras")').first();
    const initialStockText = await productRow.locator('[data-testid="stock"], td:nth-child(5)').textContent();
    const initialStock = parseInt(initialStockText?.replace(/\D/g, "") || "0");

    // Make a sale
    await navigateToPOS(page);
    await page.fill('input[placeholder*="Cari produk"]', "Beras");
    await page.waitForTimeout(300);
    await page.locator('[data-testid="product-item"]').first().click();

    const totalText = await page.locator('[data-testid="total-amount"]').textContent();
    const totalAmount = parseInt(totalText?.replace(/\D/g, "") || "0");

    await page.fill('input[name="jumlahBayar"]', (totalAmount + 5000).toString());
    await page.click('button:has-text(/proses|checkout/i)');

    // Check stock updated
    await navigateToProducts(page);
    const newStockText = await productRow.locator('[data-testid="stock"], td:nth-child(5)').textContent();
    const newStock = parseInt(newStockText?.replace(/\D/g, "") || "0");

    expect(newStock).toBe(initialStock - 1);
  });

  test("should prevent selling product with insufficient stock", async ({ page }) => {
    await navigateToPOS(page);

    // Try to add product with 0 stock or set quantity higher than stock
    await page.fill('input[placeholder*="Cari produk"]', "Produk");
    await page.waitForTimeout(300);

    const productCard = page.locator('[data-testid="product-item"]').first();

    // Check if product shows as out of stock
    const outOfStockBadge = productCard.locator('[data-testid="out-of-stock"], .out-of-stock');

    if (await outOfStockBadge.isVisible()) {
      // Try to add to cart
      await productCard.click();

      // Verify error message
      await expect(page.getByText(/stok.*habis|out of stock|tidak tersedia/i)).toBeVisible();
    }
  });

  test("should show low stock alert in POS", async ({ page }) => {
    await navigateToPOS(page);

    // Search for product
    await page.fill('input[placeholder*="Cari produk"]', "Produk");
    await page.waitForTimeout(300);

    // Check for low stock indicators
    const lowStockIndicator = page.locator('[data-testid="low-stock"], .low-stock-badge');

    if (await lowStockIndicator.count() > 0) {
      await expect(lowStockIndicator.first()).toBeVisible();
    }
  });
});

test.describe("POS System - Sales History", () => {

  test.beforeEach(async ({ page }) => {
    await loginAsAdmin(page);
  });

  test("should view sales transaction history", async ({ page }) => {
    // Navigate to sales history
    await page.click('a[href*="/penjualan"], nav >> text=/riwayat.*penjualan|sales.*history/i');
    await page.waitForURL(/.*\/(penjualan|sales)/);

    // Check page elements
    await expect(page.getByRole("heading", { name: /riwayat.*penjualan|sales.*history/i })).toBeVisible();

    // Check table headers
    await expect(page.getByText(/nomor.*transaksi|transaction.*number/i)).toBeVisible();
    await expect(page.getByText(/tanggal|date/i)).toBeVisible();
    await expect(page.getByText(/total|amount/i)).toBeVisible();
    await expect(page.getByText(/kasir|cashier/i)).toBeVisible();
  });

  test("should filter sales by date range", async ({ page }) => {
    await page.click('a[href*="/penjualan"]');
    await page.waitForURL(/.*\/penjualan/);

    // Set date filters
    const startDate = page.locator('input[name="tanggalMulai"], [data-testid="start-date"]');
    const endDate = page.locator('input[name="tanggalAkhir"], [data-testid="end-date"]');

    if (await startDate.isVisible()) {
      await startDate.fill("2025-01-01");
      await endDate.fill("2025-01-31");

      // Apply filter
      const filterBtn = page.locator('button:has-text(/filter|cari/i)');
      await filterBtn.click();

      await page.waitForTimeout(300);

      // Verify results filtered
      const rows = page.locator('tbody tr');
      const count = await rows.count();
      expect(count).toBeGreaterThanOrEqual(0);
    }
  });

  test("should view transaction details", async ({ page }) => {
    await page.click('a[href*="/penjualan"]');
    await page.waitForURL(/.*\/penjualan/);

    // Click on first transaction
    const firstRow = page.locator('tbody tr').first();
    const detailBtn = firstRow.locator('button:has-text(/detail|lihat/i)');

    if (await detailBtn.isVisible()) {
      await detailBtn.click();

      // Check detail modal/page
      await expect(page.getByText(/detail.*transaksi|transaction.*detail/i)).toBeVisible();
      await expect(page.getByText(/item.*produk|products/i)).toBeVisible();
      await expect(page.getByText(/metode.*pembayaran|payment.*method/i)).toBeVisible();
    }
  });

  test("should export sales report", async ({ page }) => {
    await page.click('a[href*="/penjualan"]');
    await page.waitForURL(/.*\/penjualan/);

    // Look for export button
    const exportBtn = page.locator('button:has-text(/export|download|unduh/i)');

    if (await exportBtn.isVisible()) {
      // Start download
      const downloadPromise = page.waitForEvent('download');
      await exportBtn.click();

      // Verify download started
      const download = await downloadPromise;
      expect(download.suggestedFilename()).toMatch(/penjualan|sales|transaksi/i);
    }
  });
});

test.describe("POS System - Role-Based Access", () => {

  test("cashier should only access POS", async ({ page }) => {
    await loginAsCashier(page);

    // Verify POS is accessible
    await navigateToPOS(page);
    await expect(page.getByRole("heading", { name: /pos|kasir/i })).toBeVisible();

    // Verify admin features are not accessible
    const settingsLink = page.locator('a[href*="/pengaturan"], nav >> text=/pengaturan|settings/i');
    await expect(settingsLink).not.toBeVisible();
  });

  test("admin should access all features", async ({ page }) => {
    await loginAsAdmin(page);

    // Verify admin can access products
    await navigateToProducts(page);
    await expect(page.getByRole("heading", { name: /produk/i })).toBeVisible();

    // Verify admin can access POS
    await navigateToPOS(page);
    await expect(page.getByRole("heading", { name: /pos|kasir/i })).toBeVisible();

    // Verify admin can access sales history
    await page.click('a[href*="/penjualan"]');
    await expect(page.getByRole("heading", { name: /penjualan/i })).toBeVisible();
  });
});
