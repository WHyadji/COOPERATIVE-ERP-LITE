-- ============================================================================
-- Migration: Add Updated_At Audit Columns
-- Date: 2025-11-20
-- Description: Add tanggal_diubah (updated_at) column to all tables to complete
--              audit trail. Includes automatic trigger to update timestamp on
--              record modification.
-- ============================================================================

-- ISSUE/CONTEXT:
-- Currently, all tables have tanggal_dibuat (created_at) and tanggal_dihapus
-- (deleted_at) for soft deletes, but are missing tanggal_diubah (updated_at).
-- This means we cannot track when records were last modified, which is important
-- for auditing, debugging, and data synchronization.

-- CHANGES:
-- 1. Add tanggal_diubah column to all 10 tables
-- 2. Create trigger function to automatically update tanggal_diubah on UPDATE
-- 3. Apply trigger to all tables
-- 4. Backfill existing records with created_at timestamp (best estimate)

BEGIN;

-- ============================================================================
-- 1. CREATE TRIGGER FUNCTION
-- ============================================================================

-- Drop existing function if it exists (idempotent)
DROP FUNCTION IF EXISTS update_tanggal_diubah() CASCADE;

-- Create function to automatically update tanggal_diubah
CREATE OR REPLACE FUNCTION update_tanggal_diubah()
RETURNS TRIGGER AS $$
BEGIN
    -- Only update if this is an UPDATE operation (not INSERT)
    IF TG_OP = 'UPDATE' THEN
        NEW.tanggal_diubah = CURRENT_TIMESTAMP;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION update_tanggal_diubah() IS
    'Automatically updates tanggal_diubah column to current timestamp on UPDATE operations';

-- ============================================================================
-- 2. ADD COLUMNS TO ALL TABLES
-- ============================================================================

-- Table: koperasi
ALTER TABLE koperasi
    ADD COLUMN IF NOT EXISTS tanggal_diubah TIMESTAMP WITH TIME ZONE;

-- Table: anggota
ALTER TABLE anggota
    ADD COLUMN IF NOT EXISTS tanggal_diubah TIMESTAMP WITH TIME ZONE;

-- Table: akun
ALTER TABLE akun
    ADD COLUMN IF NOT EXISTS tanggal_diubah TIMESTAMP WITH TIME ZONE;

-- Table: transaksi
ALTER TABLE transaksi
    ADD COLUMN IF NOT EXISTS tanggal_diubah TIMESTAMP WITH TIME ZONE;

-- Table: baris_transaksi
ALTER TABLE baris_transaksi
    ADD COLUMN IF NOT EXISTS tanggal_diubah TIMESTAMP WITH TIME ZONE;

-- Table: simpanan
ALTER TABLE simpanan
    ADD COLUMN IF NOT EXISTS tanggal_diubah TIMESTAMP WITH TIME ZONE;

-- Table: produk
ALTER TABLE produk
    ADD COLUMN IF NOT EXISTS tanggal_diubah TIMESTAMP WITH TIME ZONE;

-- Table: penjualan
ALTER TABLE penjualan
    ADD COLUMN IF NOT EXISTS tanggal_diubah TIMESTAMP WITH TIME ZONE;

-- Table: item_penjualan
ALTER TABLE item_penjualan
    ADD COLUMN IF NOT EXISTS tanggal_diubah TIMESTAMP WITH TIME ZONE;

-- Table: pengguna
ALTER TABLE pengguna
    ADD COLUMN IF NOT EXISTS tanggal_diubah TIMESTAMP WITH TIME ZONE;

-- ============================================================================
-- 3. BACKFILL EXISTING RECORDS
-- ============================================================================
-- For existing records, set tanggal_diubah to tanggal_dibuat as a reasonable
-- default (assumes records haven't been updated, which is the best we can do)

UPDATE koperasi SET tanggal_diubah = tanggal_dibuat WHERE tanggal_diubah IS NULL;
UPDATE anggota SET tanggal_diubah = tanggal_dibuat WHERE tanggal_diubah IS NULL;
UPDATE akun SET tanggal_diubah = tanggal_dibuat WHERE tanggal_diubah IS NULL;
UPDATE transaksi SET tanggal_diubah = tanggal_dibuat WHERE tanggal_diubah IS NULL;
UPDATE baris_transaksi SET tanggal_diubah = tanggal_dibuat WHERE tanggal_diubah IS NULL;
UPDATE simpanan SET tanggal_diubah = tanggal_dibuat WHERE tanggal_diubah IS NULL;
UPDATE produk SET tanggal_diubah = tanggal_dibuat WHERE tanggal_diubah IS NULL;
UPDATE penjualan SET tanggal_diubah = tanggal_dibuat WHERE tanggal_diubah IS NULL;
UPDATE item_penjualan SET tanggal_diubah = tanggal_dibuat WHERE tanggal_diubah IS NULL;
UPDATE pengguna SET tanggal_diubah = tanggal_dibuat WHERE tanggal_diubah IS NULL;

-- ============================================================================
-- 4. CREATE TRIGGERS FOR ALL TABLES
-- ============================================================================

-- Drop existing triggers if they exist (idempotent)
DROP TRIGGER IF EXISTS trg_koperasi_update ON koperasi;
DROP TRIGGER IF EXISTS trg_anggota_update ON anggota;
DROP TRIGGER IF EXISTS trg_akun_update ON akun;
DROP TRIGGER IF EXISTS trg_transaksi_update ON transaksi;
DROP TRIGGER IF EXISTS trg_baris_transaksi_update ON baris_transaksi;
DROP TRIGGER IF EXISTS trg_simpanan_update ON simpanan;
DROP TRIGGER IF EXISTS trg_produk_update ON produk;
DROP TRIGGER IF EXISTS trg_penjualan_update ON penjualan;
DROP TRIGGER IF EXISTS trg_item_penjualan_update ON item_penjualan;
DROP TRIGGER IF EXISTS trg_pengguna_update ON pengguna;

-- Create triggers
CREATE TRIGGER trg_koperasi_update
    BEFORE UPDATE ON koperasi
    FOR EACH ROW
    EXECUTE FUNCTION update_tanggal_diubah();

CREATE TRIGGER trg_anggota_update
    BEFORE UPDATE ON anggota
    FOR EACH ROW
    EXECUTE FUNCTION update_tanggal_diubah();

CREATE TRIGGER trg_akun_update
    BEFORE UPDATE ON akun
    FOR EACH ROW
    EXECUTE FUNCTION update_tanggal_diubah();

CREATE TRIGGER trg_transaksi_update
    BEFORE UPDATE ON transaksi
    FOR EACH ROW
    EXECUTE FUNCTION update_tanggal_diubah();

CREATE TRIGGER trg_baris_transaksi_update
    BEFORE UPDATE ON baris_transaksi
    FOR EACH ROW
    EXECUTE FUNCTION update_tanggal_diubah();

CREATE TRIGGER trg_simpanan_update
    BEFORE UPDATE ON simpanan
    FOR EACH ROW
    EXECUTE FUNCTION update_tanggal_diubah();

CREATE TRIGGER trg_produk_update
    BEFORE UPDATE ON produk
    FOR EACH ROW
    EXECUTE FUNCTION update_tanggal_diubah();

CREATE TRIGGER trg_penjualan_update
    BEFORE UPDATE ON penjualan
    FOR EACH ROW
    EXECUTE FUNCTION update_tanggal_diubah();

CREATE TRIGGER trg_item_penjualan_update
    BEFORE UPDATE ON item_penjualan
    FOR EACH ROW
    EXECUTE FUNCTION update_tanggal_diubah();

CREATE TRIGGER trg_pengguna_update
    BEFORE UPDATE ON pengguna
    FOR EACH ROW
    EXECUTE FUNCTION update_tanggal_diubah();

-- ============================================================================
-- VERIFICATION
-- ============================================================================

-- Verify columns were added to all tables
SELECT
    table_name,
    column_name,
    data_type,
    is_nullable
FROM information_schema.columns
WHERE table_schema = 'public'
  AND column_name = 'tanggal_diubah'
  AND table_name IN (
      'koperasi', 'anggota', 'akun', 'transaksi', 'baris_transaksi',
      'simpanan', 'produk', 'penjualan', 'item_penjualan', 'pengguna'
  )
ORDER BY table_name;

-- Verify triggers were created
SELECT
    trigger_name,
    event_object_table as table_name,
    action_timing,
    event_manipulation
FROM information_schema.triggers
WHERE trigger_schema = 'public'
  AND trigger_name LIKE 'trg_%_update'
ORDER BY event_object_table;

-- Check backfill results (count records with updated_at set)
SELECT
    'koperasi' as table_name,
    COUNT(*) as total_records,
    COUNT(tanggal_diubah) as records_with_updated_at
FROM koperasi
UNION ALL
SELECT 'anggota', COUNT(*), COUNT(tanggal_diubah) FROM anggota
UNION ALL
SELECT 'akun', COUNT(*), COUNT(tanggal_diubah) FROM akun
UNION ALL
SELECT 'transaksi', COUNT(*), COUNT(tanggal_diubah) FROM transaksi
UNION ALL
SELECT 'baris_transaksi', COUNT(*), COUNT(tanggal_diubah) FROM baris_transaksi
UNION ALL
SELECT 'simpanan', COUNT(*), COUNT(tanggal_diubah) FROM simpanan
UNION ALL
SELECT 'produk', COUNT(*), COUNT(tanggal_diubah) FROM produk
UNION ALL
SELECT 'penjualan', COUNT(*), COUNT(tanggal_diubah) FROM penjualan
UNION ALL
SELECT 'item_penjualan', COUNT(*), COUNT(tanggal_diubah) FROM item_penjualan
UNION ALL
SELECT 'pengguna', COUNT(*), COUNT(tanggal_diubah) FROM pengguna
ORDER BY table_name;

SELECT 'Migration 006: Updated_at columns and triggers added successfully' as status;

COMMIT;

-- ============================================================================
-- ROLLBACK INSTRUCTIONS
-- ============================================================================
-- If you need to rollback this migration, run the following:
--
-- BEGIN;
--
-- -- Drop triggers
-- DROP TRIGGER IF EXISTS trg_koperasi_update ON koperasi;
-- DROP TRIGGER IF EXISTS trg_anggota_update ON anggota;
-- DROP TRIGGER IF EXISTS trg_akun_update ON akun;
-- DROP TRIGGER IF EXISTS trg_transaksi_update ON transaksi;
-- DROP TRIGGER IF EXISTS trg_baris_transaksi_update ON baris_transaksi;
-- DROP TRIGGER IF EXISTS trg_simpanan_update ON simpanan;
-- DROP TRIGGER IF EXISTS trg_produk_update ON produk;
-- DROP TRIGGER IF EXISTS trg_penjualan_update ON penjualan;
-- DROP TRIGGER IF EXISTS trg_item_penjualan_update ON item_penjualan;
-- DROP TRIGGER IF EXISTS trg_pengguna_update ON pengguna;
--
-- -- Drop trigger function
-- DROP FUNCTION IF EXISTS update_tanggal_diubah() CASCADE;
--
-- -- Drop columns
-- ALTER TABLE koperasi DROP COLUMN IF EXISTS tanggal_diubah;
-- ALTER TABLE anggota DROP COLUMN IF EXISTS tanggal_diubah;
-- ALTER TABLE akun DROP COLUMN IF EXISTS tanggal_diubah;
-- ALTER TABLE transaksi DROP COLUMN IF EXISTS tanggal_diubah;
-- ALTER TABLE baris_transaksi DROP COLUMN IF EXISTS tanggal_diubah;
-- ALTER TABLE simpanan DROP COLUMN IF EXISTS tanggal_diubah;
-- ALTER TABLE produk DROP COLUMN IF EXISTS tanggal_diubah;
-- ALTER TABLE penjualan DROP COLUMN IF EXISTS tanggal_diubah;
-- ALTER TABLE item_penjualan DROP COLUMN IF EXISTS tanggal_diubah;
-- ALTER TABLE pengguna DROP COLUMN IF EXISTS tanggal_diubah;
--
-- SELECT 'Migration 006: Rolled back successfully' as status;
--
-- COMMIT;
-- ============================================================================

-- TESTING
-- ============================================================================
-- After applying this migration, test the trigger:
--
-- -- Get a test record
-- SELECT id, nama_produk, tanggal_dibuat, tanggal_diubah
-- FROM produk LIMIT 1;
--
-- -- Update the record
-- UPDATE produk
-- SET nama_produk = nama_produk || ' (Updated)'
-- WHERE id = '<id_from_above>';
--
-- -- Check that tanggal_diubah was automatically updated
-- SELECT id, nama_produk, tanggal_dibuat, tanggal_diubah
-- FROM produk WHERE id = '<id_from_above>';
--
-- -- tanggal_diubah should now be later than tanggal_dibuat
-- ============================================================================
