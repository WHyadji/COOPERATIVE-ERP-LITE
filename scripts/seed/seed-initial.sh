#!/bin/bash

# Initial Seed Data Script
# Creates default cooperative and admin user

set -e

echo "🌱 Seeding initial data..."
echo ""

# Connect to database and create initial data
docker compose exec -T postgres psql -U postgres -d koperasi_erp << 'EOF'

-- Check if koperasi already exists
DO $$
DECLARE
    koperasi_id UUID;
    admin_id UUID;
    admin_exists BOOLEAN;
BEGIN
    -- Generate UUIDs
    koperasi_id := gen_random_uuid();
    admin_id := gen_random_uuid();

    -- Check if admin already exists
    SELECT EXISTS(SELECT 1 FROM pengguna WHERE nama_pengguna = 'admin') INTO admin_exists;

    IF NOT admin_exists THEN
        -- Create default cooperative
        INSERT INTO koperasi (id, nama_koperasi, alamat, no_telepon, email, tanggal_dibuat, tanggal_diperbarui)
        VALUES (
            koperasi_id,
            'Koperasi Demo',
            'Jl. Raya Koperasi No. 1, Jakarta',
            '021-12345678',
            'info@koperasi-demo.co.id',
            NOW(),
            NOW()
        )
        ON CONFLICT (id) DO NOTHING;

        -- Create admin user
        -- Password: admin123 (hashed with bcrypt cost 10)
        INSERT INTO pengguna (id, id_koperasi, nama_pengguna, nama_lengkap, email, kata_sandi_hash, peran, status_aktif, tanggal_dibuat, tanggal_diperbarui)
        VALUES (
            admin_id,
            koperasi_id,
            'admin',
            'Administrator',
            'admin@koperasi-demo.co.id',
            '$2a$10$exdSBdtyRfCcCmjXxftn.OVBLobJxJWeUzQL4CevTlghODaMuiG5W',
            'admin',
            true,
            NOW(),
            NOW()
        )
        ON CONFLICT (id) DO NOTHING;

        RAISE NOTICE '✅ Initial data created successfully!';
        RAISE NOTICE 'Koperasi ID: %', koperasi_id;
        RAISE NOTICE 'Admin ID: %', admin_id;
    ELSE
        RAISE NOTICE 'ℹ️  Admin user already exists. Skipping seed.';
    END IF;
END $$;

-- Show created data
SELECT 'Cooperatives:' as info;
SELECT id, nama_koperasi, email FROM koperasi;

SELECT 'Users:' as info;
SELECT id, nama_pengguna, nama_lengkap, peran FROM pengguna;

EOF

echo ""
echo "✅ Seed complete!"
echo ""
echo "You can now login with:"
echo "Username: admin"
echo "Password: admin123"
