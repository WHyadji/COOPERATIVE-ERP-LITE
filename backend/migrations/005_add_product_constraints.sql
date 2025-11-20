-- ============================================================================
-- Migration: Add Product Data Integrity Constraints
-- Date: 2025-11-20
-- Description: Add CHECK constraints to prevent invalid product data such as
--              negative prices, negative stock, and invalid sale item quantities.
-- ============================================================================

-- ISSUE/CONTEXT:
-- Products can have invalid data such as zero/negative prices or negative stock
-- levels. This migration adds database-level validation to ensure product data
-- integrity and prevent common data entry errors.

-- CHANGES:
-- 1. Add CHECK constraints for product prices (must be positive)
-- 2. Add CHECK constraint for product stock (must be non-negative)
-- 3. Add CHECK constraints for sale item quantities and prices
-- 4. Ensure subtotal calculations are consistent

BEGIN;

-- ============================================================================
-- 1. PRODUCT MASTER DATA (produk)
-- ============================================================================

-- Verify current state - check for invalid prices and stock
DO $$
DECLARE
    zero_price_count INTEGER;
    neg_price_count INTEGER;
    zero_buy_price_count INTEGER;
    neg_buy_price_count INTEGER;
    neg_stock_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO zero_price_count FROM produk WHERE harga = 0;
    SELECT COUNT(*) INTO neg_price_count FROM produk WHERE harga < 0;
    SELECT COUNT(*) INTO zero_buy_price_count FROM produk WHERE harga_beli = 0;
    SELECT COUNT(*) INTO neg_buy_price_count FROM produk WHERE harga_beli < 0;
    SELECT COUNT(*) INTO neg_stock_count FROM produk WHERE stok < 0;

    IF neg_price_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % products with negative selling price', neg_price_count;
    END IF;

    IF zero_price_count > 0 THEN
        RAISE NOTICE 'INFO: Found % products with zero selling price (will be prevented)', zero_price_count;
    END IF;

    IF neg_buy_price_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % products with negative buy price', neg_buy_price_count;
    END IF;

    IF neg_stock_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % products with negative stock', neg_stock_count;
    END IF;

    IF neg_price_count = 0 AND zero_price_count = 0 AND
       neg_buy_price_count = 0 AND neg_stock_count = 0 THEN
        RAISE NOTICE 'OK: All product data is valid';
    END IF;
END $$;

-- Add constraint: Selling price must be positive (> 0)
ALTER TABLE produk
    ADD CONSTRAINT chk_produk_harga_positive
    CHECK (harga > 0);

-- Add constraint: Buy price must be positive OR NULL
-- (NULL allowed for products without purchase tracking)
ALTER TABLE produk
    ADD CONSTRAINT chk_produk_harga_beli_positive
    CHECK (harga_beli > 0 OR harga_beli IS NULL);

-- Add constraint: Stock cannot be negative
-- (Use >= 0 to allow zero stock - out of stock products)
ALTER TABLE produk
    ADD CONSTRAINT chk_produk_stok_nonnegative
    CHECK (stok >= 0);

-- Optional: Minimum stock should be non-negative if set
ALTER TABLE produk
    ADD CONSTRAINT chk_produk_stok_minimum_nonnegative
    CHECK (stok_minimum >= 0 OR stok_minimum IS NULL);

-- ============================================================================
-- 2. SALE LINE ITEMS (item_penjualan)
-- ============================================================================

-- Verify current state
DO $$
DECLARE
    zero_qty_count INTEGER;
    neg_qty_count INTEGER;
    zero_price_count INTEGER;
    neg_price_count INTEGER;
    zero_subtotal_count INTEGER;
    neg_subtotal_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO zero_qty_count FROM item_penjualan WHERE kuantitas = 0;
    SELECT COUNT(*) INTO neg_qty_count FROM item_penjualan WHERE kuantitas < 0;
    SELECT COUNT(*) INTO zero_price_count FROM item_penjualan WHERE harga_satuan = 0;
    SELECT COUNT(*) INTO neg_price_count FROM item_penjualan WHERE harga_satuan < 0;
    SELECT COUNT(*) INTO zero_subtotal_count FROM item_penjualan WHERE subtotal = 0;
    SELECT COUNT(*) INTO neg_subtotal_count FROM item_penjualan WHERE subtotal < 0;

    IF neg_qty_count > 0 OR zero_qty_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % negative and % zero quantity values',
                     neg_qty_count, zero_qty_count;
    END IF;

    IF neg_price_count > 0 OR zero_price_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % negative and % zero price values',
                     neg_price_count, zero_price_count;
    END IF;

    IF neg_subtotal_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % negative subtotal values', neg_subtotal_count;
    END IF;

    IF neg_qty_count = 0 AND zero_qty_count = 0 AND
       neg_price_count = 0 AND zero_price_count = 0 AND
       neg_subtotal_count = 0 THEN
        RAISE NOTICE 'OK: All sale item data is valid';
    END IF;
END $$;

-- Add constraint: Quantity must be positive (> 0)
-- Cannot sell 0 or negative quantity
ALTER TABLE item_penjualan
    ADD CONSTRAINT chk_item_kuantitas_positive
    CHECK (kuantitas > 0);

-- Add constraint: Unit price must be positive (> 0)
ALTER TABLE item_penjualan
    ADD CONSTRAINT chk_item_harga_positive
    CHECK (harga_satuan > 0);

-- Add constraint: Subtotal must be non-negative
-- (Subtotal = kuantitas * harga_satuan, should always be positive given above constraints)
ALTER TABLE item_penjualan
    ADD CONSTRAINT chk_item_subtotal_nonnegative
    CHECK (subtotal >= 0);

-- ============================================================================
-- VERIFICATION
-- ============================================================================

-- Verify all constraints were created successfully
SELECT
    conname as constraint_name,
    conrelid::regclass as table_name,
    pg_get_constraintdef(oid) as constraint_definition
FROM pg_constraint
WHERE conname IN (
    'chk_produk_harga_positive',
    'chk_produk_harga_beli_positive',
    'chk_produk_stok_nonnegative',
    'chk_produk_stok_minimum_nonnegative',
    'chk_item_kuantitas_positive',
    'chk_item_harga_positive',
    'chk_item_subtotal_nonnegative'
)
ORDER BY table_name, constraint_name;

SELECT 'Migration 005: Product constraints added successfully' as status;

COMMIT;

-- ============================================================================
-- ROLLBACK INSTRUCTIONS
-- ============================================================================
-- If you need to rollback this migration, run the following:
--
-- BEGIN;
--
-- ALTER TABLE produk
--     DROP CONSTRAINT IF EXISTS chk_produk_harga_positive,
--     DROP CONSTRAINT IF EXISTS chk_produk_harga_beli_positive,
--     DROP CONSTRAINT IF EXISTS chk_produk_stok_nonnegative,
--     DROP CONSTRAINT IF EXISTS chk_produk_stok_minimum_nonnegative;
--
-- ALTER TABLE item_penjualan
--     DROP CONSTRAINT IF EXISTS chk_item_kuantitas_positive,
--     DROP CONSTRAINT IF EXISTS chk_item_harga_positive,
--     DROP CONSTRAINT IF EXISTS chk_item_subtotal_nonnegative;
--
-- SELECT 'Migration 005: Rolled back successfully' as status;
--
-- COMMIT;
-- ============================================================================

-- TESTING
-- ============================================================================
-- After applying this migration, test with:
--
-- -- Should FAIL (zero price)
-- INSERT INTO produk (id, id_koperasi, kode_produk, nama_produk, harga, stok)
-- VALUES (gen_random_uuid(), (SELECT id FROM koperasi LIMIT 1),
--         'TEST001', 'Test Product', 0, 10);
-- -- ERROR: new row violates check constraint "chk_produk_harga_positive"
--
-- -- Should FAIL (negative stock)
-- INSERT INTO produk (id, id_koperasi, kode_produk, nama_produk, harga, stok)
-- VALUES (gen_random_uuid(), (SELECT id FROM koperasi LIMIT 1),
--         'TEST002', 'Test Product', 1000, -5);
-- -- ERROR: new row violates check constraint "chk_produk_stok_nonnegative"
--
-- -- Should SUCCEED (valid data)
-- INSERT INTO produk (id, id_koperasi, kode_produk, nama_produk, harga, stok)
-- VALUES (gen_random_uuid(), (SELECT id FROM koperasi LIMIT 1),
--         'TEST003', 'Test Product', 1000, 0);
-- -- Success (zero stock is allowed - out of stock)
-- ============================================================================
