#!/bin/bash

# Seed Chart of Accounts
# Creates essential accounts for simpanan transactions

set -e

echo "📊 Seeding Chart of Accounts..."
echo ""

# Connect to database and create COA
docker compose exec -T postgres psql -U postgres -d koperasi_erp << 'EOF'

-- Get the koperasi ID
DO $$
DECLARE
    koperasi_id UUID;
    kas_id UUID := gen_random_uuid();
    simpanan_pokok_id UUID := gen_random_uuid();
    simpanan_wajib_id UUID := gen_random_uuid();
    simpanan_sukarela_id UUID := gen_random_uuid();
BEGIN
    -- Get koperasi ID
    SELECT id INTO koperasi_id FROM koperasi LIMIT 1;

    -- Insert essential accounts
    -- 1101 - Kas (Cash)
    INSERT INTO akun (id, id_koperasi, kode_akun, nama_akun, tipe_akun, normal_saldo, status_aktif, tanggal_dibuat, tanggal_diperbarui)
    VALUES (
        kas_id,
        koperasi_id,
        '1101',
        'Kas',
        'aset',
        'debit',
        true,
        NOW(),
        NOW()
    )
    ON CONFLICT (id) DO NOTHING;

    -- 3101 - Simpanan Pokok (Principal Share Capital)
    INSERT INTO akun (id, id_koperasi, kode_akun, nama_akun, tipe_akun, normal_saldo, status_aktif, tanggal_dibuat, tanggal_diperbarui)
    VALUES (
        simpanan_pokok_id,
        koperasi_id,
        '3101',
        'Simpanan Pokok Anggota',
        'ekuitas',
        'kredit',
        true,
        NOW(),
        NOW()
    )
    ON CONFLICT (id) DO NOTHING;

    -- 3102 - Simpanan Wajib (Mandatory Share Capital)
    INSERT INTO akun (id, id_koperasi, kode_akun, nama_akun, tipe_akun, normal_saldo, status_aktif, tanggal_dibuat, tanggal_diperbarui)
    VALUES (
        simpanan_wajib_id,
        koperasi_id,
        '3102',
        'Simpanan Wajib Anggota',
        'ekuitas',
        'kredit',
        true,
        NOW(),
        NOW()
    )
    ON CONFLICT (id) DO NOTHING;

    -- 3103 - Simpanan Sukarela (Voluntary Share Capital)
    INSERT INTO akun (id, id_koperasi, kode_akun, nama_akun, tipe_akun, normal_saldo, status_aktif, tanggal_dibuat, tanggal_diperbarui)
    VALUES (
        simpanan_sukarela_id,
        koperasi_id,
        '3103',
        'Simpanan Sukarela Anggota',
        'ekuitas',
        'kredit',
        true,
        NOW(),
        NOW()
    )
    ON CONFLICT (id) DO NOTHING;

    RAISE NOTICE '✅ Chart of Accounts created successfully!';
END $$;

-- Show created accounts
SELECT 'Chart of Accounts:' as info;
SELECT kode_akun, nama_akun, tipe_akun, normal_saldo FROM akun ORDER BY kode_akun;

EOF

echo ""
echo "✅ Chart of Accounts seed complete!"
