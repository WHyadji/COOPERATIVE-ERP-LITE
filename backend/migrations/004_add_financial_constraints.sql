-- ============================================================================
-- Migration: Add Financial Data Integrity Constraints
-- Date: 2025-11-20
-- Description: Add CHECK constraints to prevent negative amounts in financial
--              transactions, ensuring accounting data integrity.
-- ============================================================================

-- ISSUE/CONTEXT:
-- Currently, no database-level constraints prevent negative values in financial
-- fields. This relies entirely on application validation, which could be bypassed
-- or have bugs. Database-level constraints provide defense-in-depth.

-- CHANGES:
-- 1. Add CHECK constraints for accounting transaction amounts (debit/credit)
-- 2. Add CHECK constraints for savings deposit amounts
-- 3. Add CHECK constraints for sales payment amounts
-- 4. Prevent negative change amounts in sales transactions

BEGIN;

-- ============================================================================
-- 1. ACCOUNTING TRANSACTION LINE ITEMS (baris_transaksi)
-- ============================================================================

-- Verify current state - check for any negative values
DO $$
DECLARE
    neg_debit_count INTEGER;
    neg_kredit_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO neg_debit_count FROM baris_transaksi WHERE jumlah_debit < 0;
    SELECT COUNT(*) INTO neg_kredit_count FROM baris_transaksi WHERE jumlah_kredit < 0;

    IF neg_debit_count > 0 OR neg_kredit_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % negative debit and % negative kredit values',
                     neg_debit_count, neg_kredit_count;
        RAISE NOTICE 'These will need to be fixed before constraint can be added';
    ELSE
        RAISE NOTICE 'OK: No negative values found in baris_transaksi';
    END IF;
END $$;

-- Add constraints for transaction line items
ALTER TABLE baris_transaksi
    ADD CONSTRAINT chk_baris_debit_positive
    CHECK (jumlah_debit >= 0);

ALTER TABLE baris_transaksi
    ADD CONSTRAINT chk_baris_kredit_positive
    CHECK (jumlah_kredit >= 0);

-- ============================================================================
-- 2. ACCOUNTING TRANSACTIONS (transaksi)
-- ============================================================================

-- Verify current state
DO $$
DECLARE
    neg_debit_count INTEGER;
    neg_kredit_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO neg_debit_count FROM transaksi WHERE total_debit < 0;
    SELECT COUNT(*) INTO neg_kredit_count FROM transaksi WHERE total_kredit < 0;

    IF neg_debit_count > 0 OR neg_kredit_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % negative total_debit and % negative total_kredit values',
                     neg_debit_count, neg_kredit_count;
    ELSE
        RAISE NOTICE 'OK: No negative values found in transaksi';
    END IF;
END $$;

-- Add constraints for transaction totals
ALTER TABLE transaksi
    ADD CONSTRAINT chk_transaksi_debit_positive
    CHECK (total_debit >= 0);

ALTER TABLE transaksi
    ADD CONSTRAINT chk_transaksi_kredit_positive
    CHECK (total_kredit >= 0);

-- ============================================================================
-- 3. SAVINGS DEPOSITS (simpanan)
-- ============================================================================

-- Verify current state
DO $$
DECLARE
    neg_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO neg_count FROM simpanan WHERE jumlah_setoran < 0;

    IF neg_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % negative jumlah_setoran values', neg_count;
    ELSE
        RAISE NOTICE 'OK: No negative values found in simpanan';
    END IF;
END $$;

-- Add constraint for savings deposits (allow zero for potential corrections)
ALTER TABLE simpanan
    ADD CONSTRAINT chk_simpanan_jumlah_positive
    CHECK (jumlah_setoran >= 0);

-- ============================================================================
-- 4. SALES TRANSACTIONS (penjualan)
-- ============================================================================

-- Verify current state
DO $$
DECLARE
    neg_bayar_count INTEGER;
    neg_kembalian_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO neg_bayar_count FROM penjualan WHERE jumlah_bayar < 0;
    SELECT COUNT(*) INTO neg_kembalian_count FROM penjualan WHERE kembalian < 0;

    IF neg_bayar_count > 0 OR neg_kembalian_count > 0 THEN
        RAISE NOTICE 'WARNING: Found % negative jumlah_bayar and % negative kembalian values',
                     neg_bayar_count, neg_kembalian_count;
    ELSE
        RAISE NOTICE 'OK: No negative values found in penjualan';
    END IF;
END $$;

-- Add constraints for sales payments
ALTER TABLE penjualan
    ADD CONSTRAINT chk_penjualan_jumlah_bayar_positive
    CHECK (jumlah_bayar >= 0);

ALTER TABLE penjualan
    ADD CONSTRAINT chk_penjualan_kembalian_positive
    CHECK (kembalian >= 0);

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
    'chk_baris_debit_positive',
    'chk_baris_kredit_positive',
    'chk_transaksi_debit_positive',
    'chk_transaksi_kredit_positive',
    'chk_simpanan_jumlah_positive',
    'chk_penjualan_jumlah_bayar_positive',
    'chk_penjualan_kembalian_positive'
)
ORDER BY table_name, constraint_name;

SELECT 'Migration 004: Financial constraints added successfully' as status;

COMMIT;

-- ============================================================================
-- ROLLBACK INSTRUCTIONS
-- ============================================================================
-- If you need to rollback this migration, run the following:
--
-- BEGIN;
--
-- ALTER TABLE baris_transaksi
--     DROP CONSTRAINT IF EXISTS chk_baris_debit_positive,
--     DROP CONSTRAINT IF EXISTS chk_baris_kredit_positive;
--
-- ALTER TABLE transaksi
--     DROP CONSTRAINT IF EXISTS chk_transaksi_debit_positive,
--     DROP CONSTRAINT IF EXISTS chk_transaksi_kredit_positive;
--
-- ALTER TABLE simpanan
--     DROP CONSTRAINT IF EXISTS chk_simpanan_jumlah_positive;
--
-- ALTER TABLE penjualan
--     DROP CONSTRAINT IF EXISTS chk_penjualan_jumlah_bayar_positive,
--     DROP CONSTRAINT IF EXISTS chk_penjualan_kembalian_positive;
--
-- SELECT 'Migration 004: Rolled back successfully' as status;
--
-- COMMIT;
-- ============================================================================

-- TESTING
-- ============================================================================
-- After applying this migration, test with:
--
-- -- Should FAIL (negative debit)
-- INSERT INTO baris_transaksi (id, id_transaksi, id_akun, jumlah_debit, jumlah_kredit)
-- VALUES (gen_random_uuid(), (SELECT id FROM transaksi LIMIT 1),
--         (SELECT id FROM akun LIMIT 1), -100, 0);
-- -- ERROR: new row violates check constraint "chk_baris_debit_positive"
--
-- -- Should SUCCEED (positive amounts)
-- INSERT INTO baris_transaksi (id, id_transaksi, id_akun, jumlah_debit, jumlah_kredit)
-- VALUES (gen_random_uuid(), (SELECT id FROM transaksi LIMIT 1),
--         (SELECT id FROM akun LIMIT 1), 100, 0);
-- ============================================================================
