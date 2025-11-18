#!/bin/bash

# Test Data Creation Script for Cooperative ERP Lite
# Creates test members and simpanan deposits

set -e

API_URL="http://localhost/api/v1"

echo "🔑 Step 1: Login and get token..."
TOKEN=$(curl -s -X POST ${API_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{"namaPengguna":"admin","kataSandi":"admin123"}' | jq -r '.data.token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  echo "❌ Login failed!"
  exit 1
fi

echo "✅ Login successful!"
echo "Token: ${TOKEN:0:50}..."
echo ""

# Function to create member
create_member() {
  local name=$1
  local nik=$2
  local phone=$3

  echo "👤 Creating member: $name..." >&2
  RESPONSE=$(curl -s -X POST ${API_URL}/anggota \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"namaLengkap\": \"$name\",
      \"nik\": \"$nik\",
      \"noTelepon\": \"$phone\",
      \"tanggalBergabung\": \"2025-01-01T00:00:00Z\",
      \"alamat\": \"Jl. Test No. 123, Jakarta\"
    }")

  MEMBER_ID=$(echo $RESPONSE | jq -r '.data.id')
  echo "✅ Member created: $name (ID: $MEMBER_ID)" >&2
  echo "$MEMBER_ID"
}

# Function to create simpanan deposit
create_deposit() {
  local member_id=$1
  local tipe=$2
  local amount=$3
  local timestamp_offset=${4:-0}

  # Use unique timestamps to avoid journal number conflicts
  local timestamp=$(date -u -v+${timestamp_offset}S +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "+${timestamp_offset} seconds" +"%Y-%m-%dT%H:%M:%SZ")

  echo "💰 Creating $tipe deposit for member $member_id: Rp $amount..." >&2
  REF_NUM=$(curl -s -X POST ${API_URL}/simpanan \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $TOKEN" \
    -d "{
      \"idAnggota\": \"$member_id\",
      \"tipeSimpanan\": \"$tipe\",
      \"tanggalTransaksi\": \"$timestamp\",
      \"jumlahSetoran\": $amount,
      \"keterangan\": \"Test deposit - $tipe\"
    }" | jq -r '.data.nomorReferensi')
  echo "✅ Deposit created: $REF_NUM" >&2
}

echo "📊 Step 2: Creating test members..."
echo ""

# Create 5 test members
MEMBER1=$(create_member "Budi Santoso" "3201234567890001" "081234567801")
MEMBER2=$(create_member "Siti Rahayu" "3201234567890002" "081234567802")
MEMBER3=$(create_member "Ahmad Hidayat" "3201234567890003" "081234567803")
MEMBER4=$(create_member "Dewi Lestari" "3201234567890004" "081234567804")
MEMBER5=$(create_member "Eko Prasetyo" "3201234567890005" "081234567805")

echo ""
echo "💰 Step 3: Creating simpanan deposits..."
echo ""

# Simpanan Pokok (one-time, each member)
create_deposit $MEMBER1 "pokok" 100000 0
create_deposit $MEMBER2 "pokok" 100000 1
create_deposit $MEMBER3 "pokok" 100000 2
create_deposit $MEMBER4 "pokok" 100000 3
create_deposit $MEMBER5 "pokok" 100000 4

echo ""

# Simpanan Wajib (monthly, varying amounts)
create_deposit $MEMBER1 "wajib" 50000 5
create_deposit $MEMBER2 "wajib" 75000 6
create_deposit $MEMBER3 "wajib" 50000 7
create_deposit $MEMBER4 "wajib" 100000 8
create_deposit $MEMBER5 "wajib" 50000 9

echo ""

# Simpanan Sukarela (voluntary, some members only)
create_deposit $MEMBER1 "sukarela" 200000 10
create_deposit $MEMBER2 "sukarela" 500000 11
create_deposit $MEMBER4 "sukarela" 300000 12

echo ""
echo "📈 Step 4: Fetching summary..."
echo ""

SUMMARY=$(curl -s -X GET ${API_URL}/simpanan/ringkasan \
  -H "Authorization: Bearer $TOKEN")

echo "Summary:"
echo $SUMMARY | jq '.'

echo ""
echo "📊 Step 5: Fetching balance report..."
echo ""

BALANCE_REPORT=$(curl -s -X GET ${API_URL}/simpanan/laporan-saldo \
  -H "Authorization: Bearer $TOKEN")

echo "Balance Report:"
echo $BALANCE_REPORT | jq '.data[] | {nama: .namaAnggota, pokok: .simpananPokok, wajib: .simpananWajib, sukarela: .simpananSukarela, total: .totalSimpanan}'

echo ""
echo "✅ Test data creation complete!"
echo ""
echo "Expected Totals:"
echo "- Simpanan Pokok: Rp 500,000 (5 members x Rp 100,000)"
echo "- Simpanan Wajib: Rp 325,000"
echo "- Simpanan Sukarela: Rp 1,000,000"
echo "- Total: Rp 1,825,000"
