-- ============================================================================
-- Migration: Implement Row-Level Security (RLS) for Multi-Tenant Isolation
-- Date: 2025-11-20
-- Description: Enable Row-Level Security on all tenant tables to ensure
--              automatic data isolation at the database level
-- ============================================================================

-- CONTEXT:
-- Row-Level Security (RLS) is a PostgreSQL feature that allows you to define
-- policies that restrict which rows users can access. This provides an
-- additional layer of security beyond application-level filtering.
--
-- BENEFITS:
-- 1. Automatic filtering - No risk of developers forgetting WHERE id_koperasi = ?
-- 2. Defense in depth - Works even if application code is compromised
-- 3. Consistent security - Applied to all queries (SELECT, INSERT, UPDATE, DELETE)
-- 4. Audit-friendly - Centralized security rules at database level
--
-- HOW IT WORKS:
-- 1. Application sets session variable: SET app.current_koperasi_id = 'uuid'
-- 2. RLS policies automatically filter all queries by this ID
-- 3. Users can ONLY see/modify data from their cooperative

BEGIN;

-- ============================================================================
-- STEP 1: Create Helper Functions
-- ============================================================================

SELECT '=== Creating RLS Helper Functions ===' as step;

-- Function to get current cooperative ID from session
CREATE OR REPLACE FUNCTION get_current_koperasi_id()
RETURNS UUID AS $$
DECLARE
    koperasi_id UUID;
BEGIN
    -- Try to get from session variable
    BEGIN
        koperasi_id := current_setting('app.current_koperasi_id', TRUE)::UUID;
    EXCEPTION
        WHEN OTHERS THEN
            koperasi_id := NULL;
    END;

    -- If not set, return NULL (which will deny access)
    RETURN koperasi_id;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION get_current_koperasi_id() IS
'Returns the current cooperative ID from session variable. Used by RLS policies.';

-- Function to set current cooperative ID (called by application)
CREATE OR REPLACE FUNCTION set_current_koperasi_id(p_koperasi_id UUID)
RETURNS VOID AS $$
BEGIN
    -- Validate that cooperative exists
    IF NOT EXISTS (SELECT 1 FROM koperasi WHERE id = p_koperasi_id) THEN
        RAISE EXCEPTION 'Invalid cooperative ID: %', p_koperasi_id;
    END IF;

    -- Set session variable
    PERFORM set_config('app.current_koperasi_id', p_koperasi_id::TEXT, FALSE);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION set_current_koperasi_id(UUID) IS
'Sets the current cooperative ID for the session. Must be called by application before queries.';

-- Function to clear current cooperative ID
CREATE OR REPLACE FUNCTION clear_current_koperasi_id()
RETURNS VOID AS $$
BEGIN
    PERFORM set_config('app.current_koperasi_id', '', FALSE);
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

COMMENT ON FUNCTION clear_current_koperasi_id() IS
'Clears the current cooperative ID from session. Call on logout or connection reset.';

SELECT 'Helper functions created successfully' as status;

-- ============================================================================
-- STEP 2: Enable RLS on Tenant Tables
-- ============================================================================

SELECT '=== Enabling RLS on Tenant Tables ===' as step;

-- Enable RLS on all tenant tables
ALTER TABLE anggota ENABLE ROW LEVEL SECURITY;
ALTER TABLE akun ENABLE ROW LEVEL SECURITY;
ALTER TABLE transaksi ENABLE ROW LEVEL SECURITY;
ALTER TABLE baris_transaksi ENABLE ROW LEVEL SECURITY;
ALTER TABLE produk ENABLE ROW LEVEL SECURITY;
ALTER TABLE penjualan ENABLE ROW LEVEL SECURITY;
ALTER TABLE item_penjualan ENABLE ROW LEVEL SECURITY;
ALTER TABLE simpanan ENABLE ROW LEVEL SECURITY;
ALTER TABLE pengguna ENABLE ROW LEVEL SECURITY;

-- Note: koperasi table does NOT have RLS because it's the parent table

SELECT 'RLS enabled on 9 tenant tables' as status;

-- ============================================================================
-- STEP 3: Create RLS Policies for Each Table
-- ============================================================================

SELECT '=== Creating RLS Policies ===' as step;

-- ----------------------------------------------------------------------------
-- ANGGOTA (Members) Policies
-- ----------------------------------------------------------------------------

-- Drop existing policies if any (idempotent)
DROP POLICY IF EXISTS anggota_select_policy ON anggota;
DROP POLICY IF EXISTS anggota_insert_policy ON anggota;
DROP POLICY IF EXISTS anggota_update_policy ON anggota;
DROP POLICY IF EXISTS anggota_delete_policy ON anggota;

-- Policy for SELECT: Users can only see members from their cooperative
CREATE POLICY anggota_select_policy ON anggota
    FOR SELECT
    USING (id_koperasi = get_current_koperasi_id());

-- Policy for INSERT: Users can only insert members to their cooperative
CREATE POLICY anggota_insert_policy ON anggota
    FOR INSERT
    WITH CHECK (id_koperasi = get_current_koperasi_id());

-- Policy for UPDATE: Users can only update members from their cooperative
CREATE POLICY anggota_update_policy ON anggota
    FOR UPDATE
    USING (id_koperasi = get_current_koperasi_id())
    WITH CHECK (id_koperasi = get_current_koperasi_id());

-- Policy for DELETE: Users can only soft-delete members from their cooperative
CREATE POLICY anggota_delete_policy ON anggota
    FOR DELETE
    USING (id_koperasi = get_current_koperasi_id());

COMMENT ON POLICY anggota_select_policy ON anggota IS 'Allow SELECT only for current cooperative';
COMMENT ON POLICY anggota_insert_policy ON anggota IS 'Allow INSERT only to current cooperative';
COMMENT ON POLICY anggota_update_policy ON anggota IS 'Allow UPDATE only for current cooperative';
COMMENT ON POLICY anggota_delete_policy ON anggota IS 'Allow DELETE only for current cooperative';

-- ----------------------------------------------------------------------------
-- AKUN (Chart of Accounts) Policies
-- ----------------------------------------------------------------------------

DROP POLICY IF EXISTS akun_select_policy ON akun;
DROP POLICY IF EXISTS akun_insert_policy ON akun;
DROP POLICY IF EXISTS akun_update_policy ON akun;
DROP POLICY IF EXISTS akun_delete_policy ON akun;

CREATE POLICY akun_select_policy ON akun
    FOR SELECT
    USING (id_koperasi = get_current_koperasi_id());

CREATE POLICY akun_insert_policy ON akun
    FOR INSERT
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY akun_update_policy ON akun
    FOR UPDATE
    USING (id_koperasi = get_current_koperasi_id())
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY akun_delete_policy ON akun
    FOR DELETE
    USING (id_koperasi = get_current_koperasi_id());

COMMENT ON POLICY akun_select_policy ON akun IS 'RLS: Tenant isolation for SELECT';
COMMENT ON POLICY akun_insert_policy ON akun IS 'RLS: Tenant isolation for INSERT';

-- ----------------------------------------------------------------------------
-- TRANSAKSI (Journal Entries) Policies
-- ----------------------------------------------------------------------------

DROP POLICY IF EXISTS transaksi_select_policy ON transaksi;
DROP POLICY IF EXISTS transaksi_insert_policy ON transaksi;
DROP POLICY IF EXISTS transaksi_update_policy ON transaksi;
DROP POLICY IF EXISTS transaksi_delete_policy ON transaksi;

CREATE POLICY transaksi_select_policy ON transaksi
    FOR SELECT
    USING (id_koperasi = get_current_koperasi_id());

CREATE POLICY transaksi_insert_policy ON transaksi
    FOR INSERT
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY transaksi_update_policy ON transaksi
    FOR UPDATE
    USING (id_koperasi = get_current_koperasi_id())
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY transaksi_delete_policy ON transaksi
    FOR DELETE
    USING (id_koperasi = get_current_koperasi_id());

-- ----------------------------------------------------------------------------
-- BARIS_TRANSAKSI (Journal Entry Lines) Policies
-- Note: Filtered through parent transaksi table
-- ----------------------------------------------------------------------------

DROP POLICY IF EXISTS baris_transaksi_select_policy ON baris_transaksi;
DROP POLICY IF EXISTS baris_transaksi_insert_policy ON baris_transaksi;
DROP POLICY IF EXISTS baris_transaksi_update_policy ON baris_transaksi;
DROP POLICY IF EXISTS baris_transaksi_delete_policy ON baris_transaksi;

CREATE POLICY baris_transaksi_select_policy ON baris_transaksi
    FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM transaksi
            WHERE transaksi.id = baris_transaksi.id_transaksi
              AND transaksi.id_koperasi = get_current_koperasi_id()
        )
    );

CREATE POLICY baris_transaksi_insert_policy ON baris_transaksi
    FOR INSERT
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM transaksi
            WHERE transaksi.id = baris_transaksi.id_transaksi
              AND transaksi.id_koperasi = get_current_koperasi_id()
        )
    );

CREATE POLICY baris_transaksi_update_policy ON baris_transaksi
    FOR UPDATE
    USING (
        EXISTS (
            SELECT 1 FROM transaksi
            WHERE transaksi.id = baris_transaksi.id_transaksi
              AND transaksi.id_koperasi = get_current_koperasi_id()
        )
    );

CREATE POLICY baris_transaksi_delete_policy ON baris_transaksi
    FOR DELETE
    USING (
        EXISTS (
            SELECT 1 FROM transaksi
            WHERE transaksi.id = baris_transaksi.id_transaksi
              AND transaksi.id_koperasi = get_current_koperasi_id()
        )
    );

-- ----------------------------------------------------------------------------
-- PRODUK (Products) Policies
-- ----------------------------------------------------------------------------

DROP POLICY IF EXISTS produk_select_policy ON produk;
DROP POLICY IF EXISTS produk_insert_policy ON produk;
DROP POLICY IF EXISTS produk_update_policy ON produk;
DROP POLICY IF EXISTS produk_delete_policy ON produk;

CREATE POLICY produk_select_policy ON produk
    FOR SELECT
    USING (id_koperasi = get_current_koperasi_id());

CREATE POLICY produk_insert_policy ON produk
    FOR INSERT
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY produk_update_policy ON produk
    FOR UPDATE
    USING (id_koperasi = get_current_koperasi_id())
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY produk_delete_policy ON produk
    FOR DELETE
    USING (id_koperasi = get_current_koperasi_id());

-- ----------------------------------------------------------------------------
-- PENJUALAN (Sales) Policies
-- ----------------------------------------------------------------------------

DROP POLICY IF EXISTS penjualan_select_policy ON penjualan;
DROP POLICY IF EXISTS penjualan_insert_policy ON penjualan;
DROP POLICY IF EXISTS penjualan_update_policy ON penjualan;
DROP POLICY IF EXISTS penjualan_delete_policy ON penjualan;

CREATE POLICY penjualan_select_policy ON penjualan
    FOR SELECT
    USING (id_koperasi = get_current_koperasi_id());

CREATE POLICY penjualan_insert_policy ON penjualan
    FOR INSERT
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY penjualan_update_policy ON penjualan
    FOR UPDATE
    USING (id_koperasi = get_current_koperasi_id())
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY penjualan_delete_policy ON penjualan
    FOR DELETE
    USING (id_koperasi = get_current_koperasi_id());

-- ----------------------------------------------------------------------------
-- ITEM_PENJUALAN (Sales Items) Policies
-- Note: Filtered through parent penjualan table
-- ----------------------------------------------------------------------------

DROP POLICY IF EXISTS item_penjualan_select_policy ON item_penjualan;
DROP POLICY IF EXISTS item_penjualan_insert_policy ON item_penjualan;
DROP POLICY IF EXISTS item_penjualan_update_policy ON item_penjualan;
DROP POLICY IF EXISTS item_penjualan_delete_policy ON item_penjualan;

CREATE POLICY item_penjualan_select_policy ON item_penjualan
    FOR SELECT
    USING (
        EXISTS (
            SELECT 1 FROM penjualan
            WHERE penjualan.id = item_penjualan.id_penjualan
              AND penjualan.id_koperasi = get_current_koperasi_id()
        )
    );

CREATE POLICY item_penjualan_insert_policy ON item_penjualan
    FOR INSERT
    WITH CHECK (
        EXISTS (
            SELECT 1 FROM penjualan
            WHERE penjualan.id = item_penjualan.id_penjualan
              AND penjualan.id_koperasi = get_current_koperasi_id()
        )
    );

CREATE POLICY item_penjualan_update_policy ON item_penjualan
    FOR UPDATE
    USING (
        EXISTS (
            SELECT 1 FROM penjualan
            WHERE penjualan.id = item_penjualan.id_penjualan
              AND penjualan.id_koperasi = get_current_koperasi_id()
        )
    );

CREATE POLICY item_penjualan_delete_policy ON item_penjualan
    FOR DELETE
    USING (
        EXISTS (
            SELECT 1 FROM penjualan
            WHERE penjualan.id = item_penjualan.id_penjualan
              AND penjualan.id_koperasi = get_current_koperasi_id()
        )
    );

-- ----------------------------------------------------------------------------
-- SIMPANAN (Deposits) Policies
-- ----------------------------------------------------------------------------

DROP POLICY IF EXISTS simpanan_select_policy ON simpanan;
DROP POLICY IF EXISTS simpanan_insert_policy ON simpanan;
DROP POLICY IF EXISTS simpanan_update_policy ON simpanan;
DROP POLICY IF EXISTS simpanan_delete_policy ON simpanan;

CREATE POLICY simpanan_select_policy ON simpanan
    FOR SELECT
    USING (id_koperasi = get_current_koperasi_id());

CREATE POLICY simpanan_insert_policy ON simpanan
    FOR INSERT
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY simpanan_update_policy ON simpanan
    FOR UPDATE
    USING (id_koperasi = get_current_koperasi_id())
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY simpanan_delete_policy ON simpanan
    FOR DELETE
    USING (id_koperasi = get_current_koperasi_id());

-- ----------------------------------------------------------------------------
-- PENGGUNA (Users) Policies
-- ----------------------------------------------------------------------------

DROP POLICY IF EXISTS pengguna_select_policy ON pengguna;
DROP POLICY IF EXISTS pengguna_insert_policy ON pengguna;
DROP POLICY IF EXISTS pengguna_update_policy ON pengguna;
DROP POLICY IF EXISTS pengguna_delete_policy ON pengguna;

CREATE POLICY pengguna_select_policy ON pengguna
    FOR SELECT
    USING (id_koperasi = get_current_koperasi_id());

CREATE POLICY pengguna_insert_policy ON pengguna
    FOR INSERT
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY pengguna_update_policy ON pengguna
    FOR UPDATE
    USING (id_koperasi = get_current_koperasi_id())
    WITH CHECK (id_koperasi = get_current_koperasi_id());

CREATE POLICY pengguna_delete_policy ON pengguna
    FOR DELETE
    USING (id_koperasi = get_current_koperasi_id());

SELECT 'RLS policies created for all tenant tables' as status;

-- ============================================================================
-- STEP 4: Create Bypass Policy for Superusers (Optional)
-- ============================================================================

-- Allow database superuser to bypass RLS for maintenance
-- Comment out if you want strict enforcement even for superusers

-- ALTER TABLE anggota FORCE ROW LEVEL SECURITY;
-- ALTER TABLE akun FORCE ROW LEVEL SECURITY;
-- etc...

-- Note: By default, table owners and superusers bypass RLS.
-- Use FORCE ROW LEVEL SECURITY to apply policies to them too.

SELECT 'Superuser bypass allowed (default behavior)' as status;

-- ============================================================================
-- STEP 5: Verification
-- ============================================================================

SELECT '=== Verifying RLS Implementation ===' as step;

-- Check which tables have RLS enabled
SELECT
    schemaname,
    tablename,
    CASE
        WHEN rowsecurity THEN '✅ RLS Enabled'
        ELSE '❌ RLS Disabled'
    END as rls_status
FROM pg_tables
WHERE schemaname = 'public'
  AND tablename IN ('anggota', 'akun', 'transaksi', 'baris_transaksi',
                    'produk', 'penjualan', 'item_penjualan', 'simpanan', 'pengguna')
ORDER BY tablename;

-- Count policies per table
SELECT
    tablename,
    COUNT(*) as policy_count
FROM pg_policies
WHERE schemaname = 'public'
GROUP BY tablename
ORDER BY tablename;

COMMIT;

-- ============================================================================
-- USAGE EXAMPLES (Run separately to test)
-- ============================================================================

/*
-- Example 1: Set cooperative context
SELECT set_current_koperasi_id('550e8400-e29b-41d4-a716-446655440001');

-- Example 2: Now all queries are automatically filtered
SELECT * FROM anggota;  -- Returns ONLY members from cooperative 001

-- Example 3: Try to access different cooperative's data (will return empty)
SELECT * FROM anggota WHERE id_koperasi = '550e8400-e29b-41d4-a716-446655440002';
-- Returns 0 rows even if data exists, because RLS blocks it!

-- Example 4: Try to insert to different cooperative (will FAIL)
INSERT INTO anggota (id, id_koperasi, nomor_anggota, ...)
VALUES (gen_random_uuid(), '550e8400-e29b-41d4-a716-446655440002', 'A999', ...);
-- ERROR: new row violates row-level security policy

-- Example 5: Clear context
SELECT clear_current_koperasi_id();
*/

-- ============================================================================
-- SUCCESS MESSAGE
-- ============================================================================

SELECT 'Migration 003 completed successfully!' as status;
SELECT 'Row-Level Security implemented on 9 tenant tables' as message;
SELECT '36 RLS policies created (4 per table: SELECT, INSERT, UPDATE, DELETE)' as detail;
SELECT 'Database now enforces multi-tenant isolation at the row level' as benefit;

-- ============================================================================
-- IMPORTANT NOTES FOR APPLICATION DEVELOPERS
-- ============================================================================

SELECT '=== IMPORTANT: Application Integration Required ===' as notice;
SELECT 'The application MUST call set_current_koperasi_id() before EVERY query!' as requirement_1;
SELECT 'Add this to your database connection initialization or middleware' as requirement_2;
SELECT 'See: backend/internal/middleware/rls_middleware.go (to be created)' as implementation;

-- ============================================================================
-- ROLLBACK INSTRUCTIONS (DANGEROUS - Removes all RLS protection!)
-- ============================================================================

/*
-- WARNING: This removes all Row-Level Security protection!
-- Only use if you need to completely revert RLS implementation.

BEGIN;

-- Drop all policies
DROP POLICY IF EXISTS anggota_select_policy ON anggota;
DROP POLICY IF EXISTS anggota_insert_policy ON anggota;
DROP POLICY IF EXISTS anggota_update_policy ON anggota;
DROP POLICY IF EXISTS anggota_delete_policy ON anggota;
-- ... (repeat for all tables)

-- Disable RLS
ALTER TABLE anggota DISABLE ROW LEVEL SECURITY;
ALTER TABLE akun DISABLE ROW LEVEL SECURITY;
ALTER TABLE transaksi DISABLE ROW LEVEL SECURITY;
ALTER TABLE baris_transaksi DISABLE ROW LEVEL SECURITY;
ALTER TABLE produk DISABLE ROW LEVEL SECURITY;
ALTER TABLE penjualan DISABLE ROW LEVEL SECURITY;
ALTER TABLE item_penjualan DISABLE ROW LEVEL SECURITY;
ALTER TABLE simpanan DISABLE ROW LEVEL SECURITY;
ALTER TABLE pengguna DISABLE ROW LEVEL SECURITY;

-- Drop helper functions
DROP FUNCTION IF EXISTS get_current_koperasi_id();
DROP FUNCTION IF EXISTS set_current_koperasi_id(UUID);
DROP FUNCTION IF EXISTS clear_current_koperasi_id();

COMMIT;
*/
