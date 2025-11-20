-- ============================================================================
-- Migration: Fix Anggota Unique Constraint for Multi-Tenant Support
-- Date: 2025-11-20
-- Description: Change unique constraint on nomor_anggota to be scoped by
--              id_koperasi to allow same member numbers across cooperatives
-- ============================================================================

-- ISSUE:
-- The current unique constraint idx_koperasi_nomor is only on nomor_anggota,
-- which prevents different cooperatives from having members with the same number.
-- This breaks multi-tenant functionality.
--
-- CURRENT: UNIQUE (nomor_anggota)
-- DESIRED: UNIQUE (id_koperasi, nomor_anggota)

BEGIN;

-- Step 1: Check current constraint
SELECT 'Current constraint:' as info;
SELECT
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'anggota'
  AND indexname = 'idx_koperasi_nomor';

-- Step 2: Drop the incorrect unique constraint
DROP INDEX IF EXISTS idx_koperasi_nomor;

-- Step 3: Create the correct unique constraint (scoped by id_koperasi)
CREATE UNIQUE INDEX idx_koperasi_nomor ON anggota (id_koperasi, nomor_anggota);

-- Step 4: Verify new constraint
SELECT 'New constraint:' as info;
SELECT
    schemaname,
    tablename,
    indexname,
    indexdef
FROM pg_indexes
WHERE tablename = 'anggota'
  AND indexname = 'idx_koperasi_nomor';

-- Step 5: Add comment to document the constraint
COMMENT ON INDEX idx_koperasi_nomor IS
'Ensures member numbers are unique within each cooperative (multi-tenant support)';

COMMIT;

-- Verification query
SELECT 'Migration completed successfully' as status;
SELECT 'Members can now have same number across different cooperatives' as note;
