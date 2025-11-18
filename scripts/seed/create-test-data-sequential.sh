#!/bin/bash

# Create test data sequentially with delays to avoid race conditions

set -e

API_URL="http://localhost/api/v1"

echo "🔑 Logging in..."
TOKEN=$(curl -s -X POST ${API_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{"namaPengguna":"admin","kataSandi":"admin123"}' | jq -r '.data.token')

echo "✅ Logged in"
echo ""

# Array of members to create
declare -a MEMBERS=(
  "Budi Santoso:3201234567890001:081234567801"
  "Siti Rahayu:3201234567890002:081234567802"
  "Ahmad Hidayat:3201234567890003:081234567803"
  "Dewi Lestari:3201234567890004:081234567804"
  "Eko Prasetyo:3201234567890005:081234567805"
)

declare -a MEMBER_IDS=()

# Create members
echo "👥 Creating members..."
for member_data in "${MEMBERS[@]}"; do
  IFS=':' read -r nama nik telp <<< "$member_data"

  echo "  Creating: $nama"
  MEMBER_ID=$(curl -s -X POST ${API_URL}/anggota \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"namaLengkap\": \"$nama\",
      \"nik\": \"$nik\",
      \"noTelepon\": \"$telp\",
      \"tanggalBergabung\": \"2025-01-01T00:00:00Z\",
      \"alamat\": \"Jl. Test No. 123, Jakarta\"
    }" | jq -r '.data.id')

  MEMBER_IDS+=("$MEMBER_ID")
  echo "    ✓ ID: $MEMBER_ID"
done

echo ""
echo "💰 Creating deposits (with 1s delay between each)..."

# Simpanan Pokok
echo "  Simpanan Pokok:"
for ID in "${MEMBER_IDS[@]}"; do
  echo "    Creating for $ID..."
  curl -s -X POST ${API_URL}/simpanan \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"idAnggota\": \"$ID\",
      \"tipeSimpanan\": \"pokok\",
      \"tanggalTransaksi\": \"2025-01-10T10:00:00Z\",
      \"jumlahSetoran\": 100000,
      \"keterangan\": \"Simpanan Pokok\"
    }" > /dev/null
  echo "    ✓ Done"
  sleep 1
done

# Simpanan Wajib
echo ""
echo "  Simpanan Wajib:"
AMOUNTS=(50000 75000 50000 100000 50000)
for i in "${!MEMBER_IDS[@]}"; do
  echo "    Creating for ${MEMBER_IDS[$i]}..."
  curl -s -X POST ${API_URL}/simpanan \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"idAnggota\": \"${MEMBER_IDS[$i]}\",
      \"tipeSimpanan\": \"wajib\",
      \"tanggalTransaksi\": \"2025-01-15T10:00:00Z\",
      \"jumlahSetoran\": ${AMOUNTS[$i]},
      \"keterangan\": \"Simpanan Wajib\"
    }" > /dev/null
  echo "    ✓ Done"
  sleep 1
done

# Simpanan Sukarela (only members 0, 1, 3)
echo ""
echo "  Simpanan Sukarela:"
SUKARELA_MEMBERS=(0 1 3)
SUKARELA_AMOUNTS=(200000 500000 300000)
for i in "${!SUKARELA_MEMBERS[@]}"; do
  idx=${SUKARELA_MEMBERS[$i]}
  echo "    Creating for ${MEMBER_IDS[$idx]}..."
  curl -s -X POST ${API_URL}/simpanan \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"idAnggota\": \"${MEMBER_IDS[$idx]}\",
      \"tipeSimpanan\": \"sukarela\",
      \"tanggalTransaksi\": \"2025-01-20T10:00:00Z\",
      \"jumlahSetoran\": ${SUKARELA_AMOUNTS[$i]},
      \"keterangan\": \"Simpanan Sukarela\"
    }" > /dev/null
  echo "    ✓ Done"
  sleep 1
done

echo ""
echo "📊 Fetching summary..."
SUMMARY=$(curl -s -X GET ${API_URL}/simpanan/ringkasan \
  -H "Authorization: Bearer $TOKEN")

echo "$SUMMARY" | jq '.data'

echo ""
echo "✅ Done!"
echo ""
echo "Expected Totals:"
echo "  - Simpanan Pokok: Rp 500,000"
echo "  - Simpanan Wajib: Rp 325,000"
echo "  - Simpanan Sukarela: Rp 1,000,000"
echo "  - Total: Rp 1,825,000"
