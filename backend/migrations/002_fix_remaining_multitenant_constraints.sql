-- ============================================================================
-- Migration: Fix Remaining Multi-Tenant Unique Constraints
-- Date: 2025-11-20
-- Description: Fix unique constraints on akun, transaksi, produk, penjualan,
--              and pengguna tables to be scoped by id_koperasi
-- ============================================================================

-- ISSUE:
-- Multiple tables have unique constraints that are NOT scoped by id_koperasi,
-- which prevents different cooperatives from using the same codes/numbers.
--
-- Tables affected:
-- 1. akun - idx_koperasi_kode_akun (kode_akun only)
-- 2. transaksi - idx_koperasi_nomor_jurnal (nomor_jurnal only)
-- 3. produk - idx_koperasi_kode_produk (kode_produk only)
-- 4. penjualan - idx_koperasi_nomor_penjualan (nomor_penjualan only)
-- 5. pengguna - idx_koperasi_username (nama_pengguna only)
--
-- SOLUTION:
-- Change all constraints to be composite: (id_koperasi, column_name)

BEGIN;

-- ============================================================================
-- 1. FIX AKUN TABLE (Chart of Accounts)
-- ============================================================================

SELECT '=== Fixing akun table ===' as step;

-- Show current constraint
SELECT 'Before:' as status, indexdef
FROM pg_indexes
WHERE tablename = 'akun' AND indexname = 'idx_koperasi_kode_akun';

-- Drop incorrect constraint
DROP INDEX IF EXISTS idx_koperasi_kode_akun;

-- Create correct constraint (scoped by id_koperasi)
CREATE UNIQUE INDEX idx_koperasi_kode_akun ON akun (id_koperasi, kode_akun);

-- Add comment
COMMENT ON INDEX idx_koperasi_kode_akun IS
'Ensures account codes are unique within each cooperative (multi-tenant support)';

-- Show new constraint
SELECT 'After:' as status, indexdef
FROM pg_indexes
WHERE tablename = 'akun' AND indexname = 'idx_koperasi_kode_akun';

-- ============================================================================
-- 2. FIX TRANSAKSI TABLE (Journal Entries)
-- ============================================================================

SELECT '=== Fixing transaksi table ===' as step;

-- Show current constraint
SELECT 'Before:' as status, indexdef
FROM pg_indexes
WHERE tablename = 'transaksi' AND indexname = 'idx_koperasi_nomor_jurnal';

-- Drop incorrect constraint
DROP INDEX IF EXISTS idx_koperasi_nomor_jurnal;

-- Create correct constraint (scoped by id_koperasi)
CREATE UNIQUE INDEX idx_koperasi_nomor_jurnal ON transaksi (id_koperasi, nomor_jurnal);

-- Add comment
COMMENT ON INDEX idx_koperasi_nomor_jurnal IS
'Ensures journal numbers are unique within each cooperative (multi-tenant support)';

-- Show new constraint
SELECT 'After:' as status, indexdef
FROM pg_indexes
WHERE tablename = 'transaksi' AND indexname = 'idx_koperasi_nomor_jurnal';

-- ============================================================================
-- 3. FIX PRODUK TABLE (Products)
-- ============================================================================

SELECT '=== Fixing produk table ===' as step;

-- Show current constraint
SELECT 'Before:' as status, indexdef
FROM pg_indexes
WHERE tablename = 'produk' AND indexname = 'idx_koperasi_kode_produk';

-- Drop incorrect constraint
DROP INDEX IF EXISTS idx_koperasi_kode_produk;

-- Create correct constraint (scoped by id_koperasi)
CREATE UNIQUE INDEX idx_koperasi_kode_produk ON produk (id_koperasi, kode_produk);

-- Add comment
COMMENT ON INDEX idx_koperasi_kode_produk IS
'Ensures product codes are unique within each cooperative (multi-tenant support)';

-- Show new constraint
SELECT 'After:' as status, indexdef
FROM pg_indexes
WHERE tablename = 'produk' AND indexname = 'idx_koperasi_kode_produk';

-- ============================================================================
-- 4. FIX PENJUALAN TABLE (Sales)
-- ============================================================================

SELECT '=== Fixing penjualan table ===' as step;

-- Show current constraint
SELECT 'Before:' as status, indexdef
FROM pg_indexes
WHERE tablename = 'penjualan' AND indexname = 'idx_koperasi_nomor_penjualan';

-- Drop incorrect constraint
DROP INDEX IF EXISTS idx_koperasi_nomor_penjualan;

-- Create correct constraint (scoped by id_koperasi)
CREATE UNIQUE INDEX idx_koperasi_nomor_penjualan ON penjualan (id_koperasi, nomor_penjualan);

-- Add comment
COMMENT ON INDEX idx_koperasi_nomor_penjualan IS
'Ensures sales numbers are unique within each cooperative (multi-tenant support)';

-- Show new constraint
SELECT 'After:' as status, indexdef
FROM pg_indexes
WHERE tablename = 'penjualan' AND indexname = 'idx_koperasi_nomor_penjualan';

-- ============================================================================
-- 5. FIX PENGGUNA TABLE (Users)
-- ============================================================================

SELECT '=== Fixing pengguna table ===' as step;

-- Show current constraint
SELECT 'Before:' as status, indexdef
FROM pg_indexes
WHERE tablename = 'pengguna' AND indexname = 'idx_koperasi_username';

-- Drop incorrect constraint
DROP INDEX IF EXISTS idx_koperasi_username;

-- Create correct constraint (scoped by id_koperasi)
CREATE UNIQUE INDEX idx_koperasi_username ON pengguna (id_koperasi, nama_pengguna);

-- Add comment
COMMENT ON INDEX idx_koperasi_username IS
'Ensures usernames are unique within each cooperative (multi-tenant support)';

-- Show new constraint
SELECT 'After:' as status, indexdef
FROM pg_indexes
WHERE tablename = 'pengguna' AND indexname = 'idx_koperasi_username';

-- ============================================================================
-- VERIFICATION
-- ============================================================================

SELECT '=== Migration Summary ===' as section;

-- Count fixed constraints
SELECT
    'Total constraints fixed:' as info,
    COUNT(*) as count
FROM pg_indexes
WHERE tablename IN ('akun', 'transaksi', 'produk', 'penjualan', 'pengguna')
  AND indexname LIKE 'idx_koperasi_%'
  AND indexdef LIKE '%id_koperasi,%';

-- List all fixed constraints
SELECT
    tablename as "Table",
    indexname as "Index Name",
    CASE
        WHEN indexdef LIKE '%id_koperasi,%' THEN '✓ Multi-tenant (Fixed)'
        ELSE '✗ Single tenant (Broken)'
    END as "Status"
FROM pg_indexes
WHERE tablename IN ('akun', 'transaksi', 'produk', 'penjualan', 'pengguna', 'anggota')
  AND indexname LIKE 'idx_koperasi_%'
ORDER BY tablename;

COMMIT;

-- ============================================================================
-- VERIFICATION QUERIES (Run separately after migration)
-- ============================================================================

-- Test 1: Verify we can now have same account codes in different cooperatives
-- Uncomment to test:
/*
BEGIN;

-- Create test cooperative 2
INSERT INTO koperasi (id, nama_koperasi, alamat, no_telepon, email, tahun_buku_mulai, pengaturan)
VALUES (
    '550e8400-e29b-41d4-a716-446655440002',
    'Test Coop 2',
    'Test Address',
    '08123456789',
    'test2@coop.com',
    1,
    '{}'
);

-- Try to create same account code in both cooperatives (should succeed now)
INSERT INTO akun (id, id_koperasi, kode_akun, nama_akun, tipe_akun, normal_saldo)
VALUES
    (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440001', '1010', 'Cash - Coop 1', 'aset', 'debit'),
    (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440002', '1010', 'Cash - Coop 2', 'aset', 'debit');

-- Verify
SELECT id_koperasi, kode_akun, nama_akun FROM akun WHERE kode_akun = '1010';

-- Cleanup
ROLLBACK;
*/

-- ============================================================================
-- SUCCESS MESSAGE
-- ============================================================================

SELECT 'Migration 002 completed successfully!' as status;
SELECT 'All 5 multi-tenant unique constraints have been fixed' as message;
SELECT 'Different cooperatives can now use the same codes/numbers' as note;

-- ============================================================================
-- ROLLBACK INSTRUCTIONS (If needed - DANGEROUS!)
-- ============================================================================
-- WARNING: Rolling back will break multi-tenant support!
-- Only rollback if you need to revert to single-tenant mode.
--
-- BEGIN;
--
-- -- Revert akun
-- DROP INDEX IF EXISTS idx_koperasi_kode_akun;
-- CREATE UNIQUE INDEX idx_koperasi_kode_akun ON akun (kode_akun);
--
-- -- Revert transaksi
-- DROP INDEX IF EXISTS idx_koperasi_nomor_jurnal;
-- CREATE UNIQUE INDEX idx_koperasi_nomor_jurnal ON transaksi (nomor_jurnal);
--
-- -- Revert produk
-- DROP INDEX IF EXISTS idx_koperasi_kode_produk;
-- CREATE UNIQUE INDEX idx_koperasi_kode_produk ON produk (kode_produk);
--
-- -- Revert penjualan
-- DROP INDEX IF EXISTS idx_koperasi_nomor_penjualan;
-- CREATE UNIQUE INDEX idx_koperasi_nomor_penjualan ON penjualan (nomor_penjualan);
--
-- -- Revert pengguna
-- DROP INDEX IF EXISTS idx_koperasi_username;
-- CREATE UNIQUE INDEX idx_koperasi_username ON pengguna (nama_pengguna);
--
-- COMMIT;
-- ============================================================================
