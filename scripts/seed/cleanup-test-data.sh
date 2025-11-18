#!/bin/bash

# Cleanup test data

set -e

echo "🧹 Cleaning up test data..."

docker compose exec -T postgres psql -U postgres -d koperasi_erp << 'EOF'

-- Hard delete simpanan (including soft deleted ones)
DELETE FROM simpanan;

-- Hard delete baris_transaksi (transaction lines)
DELETE FROM baris_transaksi;

-- Hard delete transaksi (accounting transactions)
DELETE FROM transaksi;

-- Hard delete anggota (members)
DELETE FROM anggota;

-- Show remaining data
SELECT 'Active Members:' as info, COUNT(*) as count FROM anggota WHERE tanggal_dihapus IS NULL
UNION ALL
SELECT 'Active Simpanan:', COUNT(*) FROM simpanan WHERE tanggal_dihapus IS NULL
UNION ALL
SELECT 'Active Transaksi:', COUNT(*) FROM transaksi WHERE tanggal_dihapus IS NULL;

EOF

echo "✅ Test data cleaned up!"
