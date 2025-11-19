-- ============================================================================
-- E2E Test Data Seeding Script (SQL Version)
-- Alternative to the Go seeding script for quick manual setup
-- ============================================================================

-- Clean up existing test data (optional - uncomment if you want fresh start)
-- DELETE FROM simpanan WHERE id_koperasi IN (SELECT id FROM koperasi WHERE nomor_badan_hukum = 'TEST-E2E-001');
-- DELETE FROM saldo_simpanan_anggota WHERE id_koperasi IN (SELECT id FROM koperasi WHERE nomor_badan_hukum = 'TEST-E2E-001');
-- DELETE FROM anggota WHERE id_koperasi IN (SELECT id FROM koperasi WHERE nomor_badan_hukum = 'TEST-E2E-001');
-- DELETE FROM koperasi WHERE nomor_badan_hukum = 'TEST-E2E-001';

-- 1. Create test cooperative
INSERT INTO koperasi (
    id, nama_koperasi, nomor_badan_hukum, tanggal_berdiri,
    alamat, kelurahan, kecamatan, kota_kabupaten, provinsi, kode_pos,
    nomor_telepon, email, website,
    jumlah_anggota, total_aset, simpanan_pokok, simpanan_wajib, simpanan_sukarela,
    aktif, created_at, updated_at
) VALUES (
    gen_random_uuid(),
    'Koperasi Test E2E',
    'TEST-E2E-001',
    '2024-01-01',
    'Jl. Test E2E No. 123',
    'Test Kelurahan',
    'Test Kecamatan',
    'Test City',
    'Test Province',
    '12345',
    '08123456789',
    'test@e2e.com',
    'https://test-e2e.com',
    1,
    5000000,
    1000000,
    500000,
    200000,
    true,
    NOW(),
    NOW()
) ON CONFLICT (nomor_badan_hukum) DO NOTHING;

-- Get the koperasi ID for subsequent inserts
DO $$
DECLARE
    v_koperasi_id UUID;
    v_anggota_id UUID;
BEGIN
    -- Get koperasi ID
    SELECT id INTO v_koperasi_id FROM koperasi WHERE nomor_badan_hukum = 'TEST-E2E-001';

    -- 2. Create test member (A001) with PIN 123456
    -- PIN hash for "123456" using bcrypt (you may need to regenerate this)
    INSERT INTO anggota (
        id, id_koperasi, nomor_anggota, nama_lengkap,
        nik, jenis_kelamin, tempat_lahir, tanggal_lahir,
        alamat, rt, rw, kelurahan, kecamatan, kota_kabupaten, provinsi, kode_pos,
        nomor_telepon, email, pekerjaan,
        tanggal_bergabung, status, pin,
        created_at, updated_at
    ) VALUES (
        gen_random_uuid(),
        v_koperasi_id,
        'A001',
        'Test Member Portal',
        '1234567890123456',
        'L',
        'Jakarta',
        '1990-01-01',
        'Jl. Test Member No. 1',
        '001',
        '002',
        'Test Kelurahan',
        'Test Kecamatan',
        'Jakarta',
        'DKI Jakarta',
        '12345',
        '081234567890',
        'test.member@email.com',
        'Karyawan Swasta',
        '2024-01-01',
        'aktif',
        '$2a$10$YourBcryptHashHere', -- This needs to be generated properly
        NOW(),
        NOW()
    ) ON CONFLICT (nomor_anggota, id_koperasi) DO NOTHING;

    -- Get anggota ID
    SELECT id INTO v_anggota_id FROM anggota WHERE nomor_anggota = 'A001' AND id_koperasi = v_koperasi_id;

    -- 3. Create initial balance
    INSERT INTO saldo_simpanan_anggota (
        id, id_koperasi, id_anggota,
        simpanan_pokok, simpanan_wajib, simpanan_sukarela, total_simpanan,
        created_at, updated_at
    ) VALUES (
        gen_random_uuid(),
        v_koperasi_id,
        v_anggota_id,
        1000000,
        2500000,
        500000,
        4000000,
        NOW(),
        NOW()
    ) ON CONFLICT (id_anggota, id_koperasi) DO NOTHING;

    -- 4. Create sample transactions
    -- Simpanan Pokok
    INSERT INTO simpanan (
        id, id_koperasi, id_anggota, nomor_referensi,
        tipe_simpanan, tipe_transaksi, tanggal_transaksi, jumlah,
        keterangan, metode_pembayaran,
        created_at, updated_at
    ) VALUES (
        gen_random_uuid(), v_koperasi_id, v_anggota_id, 'SP-2024-001',
        'pokok', 'setoran', '2024-01-01 10:00:00', 1000000,
        'Setoran Simpanan Pokok', 'tunai',
        NOW(), NOW()
    ) ON CONFLICT (nomor_referensi, id_koperasi) DO NOTHING;

    -- Simpanan Wajib (monthly)
    INSERT INTO simpanan (
        id, id_koperasi, id_anggota, nomor_referensi,
        tipe_simpanan, tipe_transaksi, tanggal_transaksi, jumlah,
        keterangan, metode_pembayaran,
        created_at, updated_at
    ) VALUES
        (gen_random_uuid(), v_koperasi_id, v_anggota_id, 'SW-2024-001',
         'wajib', 'setoran', '2024-01-15 10:00:00', 500000,
         'Setoran Simpanan Wajib Januari 2024', 'tunai', NOW(), NOW()),
        (gen_random_uuid(), v_koperasi_id, v_anggota_id, 'SW-2024-002',
         'wajib', 'setoran', '2024-02-15 10:00:00', 500000,
         'Setoran Simpanan Wajib Februari 2024', 'tunai', NOW(), NOW()),
        (gen_random_uuid(), v_koperasi_id, v_anggota_id, 'SW-2024-003',
         'wajib', 'setoran', '2024-03-15 10:00:00', 500000,
         'Setoran Simpanan Wajib Maret 2024', 'tunai', NOW(), NOW()),
        (gen_random_uuid(), v_koperasi_id, v_anggota_id, 'SW-2024-004',
         'wajib', 'setoran', '2024-04-15 10:00:00', 500000,
         'Setoran Simpanan Wajib April 2024', 'tunai', NOW(), NOW()),
        (gen_random_uuid(), v_koperasi_id, v_anggota_id, 'SW-2024-005',
         'wajib', 'setoran', '2024-05-15 10:00:00', 500000,
         'Setoran Simpanan Wajib Mei 2024', 'tunai', NOW(), NOW())
    ON CONFLICT (nomor_referensi, id_koperasi) DO NOTHING;

    -- Simpanan Sukarela
    INSERT INTO simpanan (
        id, id_koperasi, id_anggota, nomor_referensi,
        tipe_simpanan, tipe_transaksi, tanggal_transaksi, jumlah,
        keterangan, metode_pembayaran,
        created_at, updated_at
    ) VALUES
        (gen_random_uuid(), v_koperasi_id, v_anggota_id, 'SS-2024-001',
         'sukarela', 'setoran', '2024-02-01 10:00:00', 200000,
         'Setoran Simpanan Sukarela', 'tunai', NOW(), NOW()),
        (gen_random_uuid(), v_koperasi_id, v_anggota_id, 'SS-2024-002',
         'sukarela', 'setoran', '2024-03-20 10:00:00', 300000,
         'Setoran Simpanan Sukarela', 'tunai', NOW(), NOW())
    ON CONFLICT (nomor_referensi, id_koperasi) DO NOTHING;

END $$;

-- Verification queries
SELECT 'Koperasi created:' as status, COUNT(*) as count FROM koperasi WHERE nomor_badan_hukum = 'TEST-E2E-001';
SELECT 'Members created:' as status, COUNT(*) as count FROM anggota WHERE nomor_anggota = 'A001';
SELECT 'Balance created:' as status, COUNT(*) as count FROM saldo_simpanan_anggota
WHERE id_anggota IN (SELECT id FROM anggota WHERE nomor_anggota = 'A001');
SELECT 'Transactions created:' as status, COUNT(*) as count FROM simpanan
WHERE id_anggota IN (SELECT id FROM anggota WHERE nomor_anggota = 'A001');

-- Display test credentials
SELECT '===========================================';
SELECT 'E2E Test Data Seeding Completed';
SELECT '===========================================';
SELECT 'Test Credentials:';
SELECT '  Nomor Anggota: A001';
SELECT '  PIN: 123456';
SELECT '===========================================';
SELECT 'IMPORTANT: The PIN hash in this SQL script is a placeholder.';
SELECT 'Please use the Go seeding script (go run cmd/seed-test-data/main.go)';
SELECT 'to properly hash the PIN.';
SELECT '===========================================';
