-- ============================================================================
-- Migration: Add Enum Validation Constraints
-- Date: 2025-11-20
-- Description: Add CHECK constraints to validate enum-like fields ensuring
--              only valid status values, roles, and types are stored.
-- ============================================================================

-- ISSUE/CONTEXT:
-- Many fields act as enums (status, role, account type, etc.) but currently
-- have no database-level validation. This allows invalid values to be stored
-- if application validation is bypassed or has bugs.

-- CHANGES:
-- 1. Add status validation for members (anggota)
-- 2. Add role validation for users (pengguna)
-- 3. Add account type validation (akun)
-- 4. Add normal balance validation (akun)
-- 5. Add transaction type validation (transaksi)
-- 6. Add payment method validation (penjualan)
-- 7. Add gender validation (anggota)
-- 8. Add savings type validation (simpanan)

BEGIN;

-- ============================================================================
-- 1. MEMBER STATUS (anggota)
-- ============================================================================

-- Verify current values
DO $$
DECLARE
    invalid_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO invalid_count
    FROM anggota
    WHERE status NOT IN ('AKTIF', 'NONAKTIF', 'KELUAR');

    IF invalid_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % members with invalid status', invalid_count;
    ELSE
        RAISE NOTICE 'OK: All member statuses are valid';
    END IF;
END $$;

-- Add constraint
ALTER TABLE anggota
    ADD CONSTRAINT chk_anggota_status
    CHECK (status IN ('AKTIF', 'NONAKTIF', 'KELUAR'));

-- Gender validation (optional)
ALTER TABLE anggota
    ADD CONSTRAINT chk_anggota_jenis_kelamin
    CHECK (jenis_kelamin IN ('L', 'P') OR jenis_kelamin IS NULL);

-- ============================================================================
-- 2. USER ROLE (pengguna)
-- ============================================================================

-- Verify current values
DO $$
DECLARE
    invalid_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO invalid_count
    FROM pengguna
    WHERE peran NOT IN ('ADMIN', 'BENDAHARA', 'KASIR', 'ANGGOTA');

    IF invalid_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % users with invalid role', invalid_count;
    ELSE
        RAISE NOTICE 'OK: All user roles are valid';
    END IF;
END $$;

-- Add constraint (peran is case-insensitive in queries but stored uppercase)
ALTER TABLE pengguna
    ADD CONSTRAINT chk_pengguna_peran
    CHECK (peran IN ('ADMIN', 'BENDAHARA', 'KASIR', 'ANGGOTA'));

-- ============================================================================
-- 3. ACCOUNT TYPE (akun)
-- ============================================================================

-- Verify current values
DO $$
DECLARE
    invalid_type_count INTEGER;
    invalid_saldo_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO invalid_type_count
    FROM akun
    WHERE tipe_akun NOT IN ('AKTIVA', 'KEWAJIBAN', 'MODAL', 'PENDAPATAN', 'BEBAN');

    SELECT COUNT(*) INTO invalid_saldo_count
    FROM akun
    WHERE normal_saldo NOT IN ('DEBIT', 'KREDIT');

    IF invalid_type_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % accounts with invalid type', invalid_type_count;
    END IF;

    IF invalid_saldo_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % accounts with invalid normal balance', invalid_saldo_count;
    END IF;

    IF invalid_type_count = 0 AND invalid_saldo_count = 0 THEN
        RAISE NOTICE 'OK: All account types and balances are valid';
    END IF;
END $$;

-- Add constraints
ALTER TABLE akun
    ADD CONSTRAINT chk_akun_tipe
    CHECK (tipe_akun IN ('AKTIVA', 'KEWAJIBAN', 'MODAL', 'PENDAPATAN', 'BEBAN'));

ALTER TABLE akun
    ADD CONSTRAINT chk_akun_normal_saldo
    CHECK (normal_saldo IN ('DEBIT', 'KREDIT'));

-- ============================================================================
-- 4. TRANSACTION TYPE (transaksi)
-- ============================================================================

-- Verify current values
DO $$
DECLARE
    invalid_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO invalid_count
    FROM transaksi
    WHERE tipe_transaksi NOT IN ('JURNAL_UMUM', 'SIMPANAN', 'PENJUALAN', 'PEMBELIAN');

    IF invalid_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % transactions with invalid type', invalid_count;
    ELSE
        RAISE NOTICE 'OK: All transaction types are valid';
    END IF;
END $$;

-- Add constraint
ALTER TABLE transaksi
    ADD CONSTRAINT chk_transaksi_tipe
    CHECK (tipe_transaksi IN ('JURNAL_UMUM', 'SIMPANAN', 'PENJUALAN', 'PEMBELIAN'));

-- ============================================================================
-- 5. PAYMENT METHOD (penjualan)
-- ============================================================================

-- Verify current values
DO $$
DECLARE
    invalid_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO invalid_count
    FROM penjualan
    WHERE metode_pembayaran NOT IN ('TUNAI', 'TRANSFER', 'QRIS')
      AND metode_pembayaran IS NOT NULL;

    IF invalid_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % sales with invalid payment method', invalid_count;
    ELSE
        RAISE NOTICE 'OK: All payment methods are valid';
    END IF;
END $$;

-- Add constraint (allow NULL for backward compatibility)
ALTER TABLE penjualan
    ADD CONSTRAINT chk_penjualan_metode
    CHECK (metode_pembayaran IN ('TUNAI', 'TRANSFER', 'QRIS') OR metode_pembayaran IS NULL);

-- ============================================================================
-- 6. SAVINGS TYPE (simpanan)
-- ============================================================================

-- Verify current values
DO $$
DECLARE
    invalid_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO invalid_count
    FROM simpanan
    WHERE tipe_simpanan NOT IN ('POKOK', 'WAJIB', 'SUKARELA');

    IF invalid_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % savings with invalid type', invalid_count;
    ELSE
        RAISE NOTICE 'OK: All savings types are valid';
    END IF;
END $$;

-- Add constraint
ALTER TABLE simpanan
    ADD CONSTRAINT chk_simpanan_tipe
    CHECK (tipe_simpanan IN ('POKOK', 'WAJIB', 'SUKARELA'));

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
    'chk_anggota_status',
    'chk_anggota_jenis_kelamin',
    'chk_pengguna_peran',
    'chk_akun_tipe',
    'chk_akun_normal_saldo',
    'chk_transaksi_tipe',
    'chk_penjualan_metode',
    'chk_simpanan_tipe'
)
ORDER BY table_name, constraint_name;

SELECT 'Migration 007: Enum constraints added successfully' as status;

COMMIT;

-- ============================================================================
-- ROLLBACK INSTRUCTIONS
-- ============================================================================
-- If you need to rollback this migration, run the following:
--
-- BEGIN;
--
-- ALTER TABLE anggota
--     DROP CONSTRAINT IF EXISTS chk_anggota_status,
--     DROP CONSTRAINT IF EXISTS chk_anggota_jenis_kelamin;
--
-- ALTER TABLE pengguna
--     DROP CONSTRAINT IF EXISTS chk_pengguna_peran;
--
-- ALTER TABLE akun
--     DROP CONSTRAINT IF EXISTS chk_akun_tipe,
--     DROP CONSTRAINT IF EXISTS chk_akun_normal_saldo;
--
-- ALTER TABLE transaksi
--     DROP CONSTRAINT IF EXISTS chk_transaksi_tipe;
--
-- ALTER TABLE penjualan
--     DROP CONSTRAINT IF EXISTS chk_penjualan_metode;
--
-- ALTER TABLE simpanan
--     DROP CONSTRAINT IF EXISTS chk_simpanan_tipe;
--
-- SELECT 'Migration 007: Rolled back successfully' as status;
--
-- COMMIT;
-- ============================================================================

-- TESTING
-- ============================================================================
-- After applying this migration, test with:
--
-- -- Should FAIL (invalid status)
-- INSERT INTO anggota (id, id_koperasi, nomor_anggota, nama_lengkap, status)
-- VALUES (gen_random_uuid(), (SELECT id FROM koperasi LIMIT 1),
--         'TEST999', 'Test Member', 'INVALID_STATUS');
-- -- ERROR: new row violates check constraint "chk_anggota_status"
--
-- -- Should SUCCEED (valid status)
-- INSERT INTO anggota (id, id_koperasi, nomor_anggota, nama_lengkap, status)
-- VALUES (gen_random_uuid(), (SELECT id FROM koperasi LIMIT 1),
--         'TEST999', 'Test Member', 'AKTIF');
-- ============================================================================
