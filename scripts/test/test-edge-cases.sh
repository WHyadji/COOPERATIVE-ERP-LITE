#!/bin/bash

API_URL="http://localhost/api/v1"
TOKEN=$(curl -s -X POST ${API_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{"namaPengguna":"admin","kataSandi":"admin123"}' | jq -r '.data.token')

echo "🧪 Testing Edge Cases and Error Handling"
echo "=========================================="
echo ""

# Test 1: Invalid member ID
echo "1️⃣ Test: Invalid Member ID"
RESPONSE=$(curl -s -X POST ${API_URL}/simpanan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "idAnggota": "00000000-0000-0000-0000-000000000000",
    "tipeSimpanan": "pokok",
    "tanggalTransaksi": "2025-01-01T00:00:00Z",
    "jumlahSetoran": 100000,
    "keterangan": "Test invalid member"
  }')
STATUS=$(echo $RESPONSE | jq -r '.success // false')
if [ "$STATUS" = "false" ]; then
  echo "  ✅ Correctly rejected: $(echo $RESPONSE | jq -r '.message // .error')"
else
  echo "  ❌ Should have rejected invalid member ID"
fi
echo ""

# Test 2: Negative amount
echo "2️⃣ Test: Negative Amount"
MEMBER_ID=$(curl -s ${API_URL}/anggota -H "Authorization: Bearer $TOKEN" | jq -r '.data[0].id')
RESPONSE=$(curl -s -X POST ${API_URL}/simpanan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"idAnggota\": \"$MEMBER_ID\",
    \"tipeSimpanan\": \"pokok\",
    \"tanggalTransaksi\": \"2025-01-01T00:00:00Z\",
    \"jumlahSetoran\": -100000,
    \"keterangan\": \"Test negative amount\"
  }")
STATUS=$(echo $RESPONSE | jq -r '.success // false')
if [ "$STATUS" = "false" ]; then
  echo "  ✅ Correctly rejected: $(echo $RESPONSE | jq -r '.message // .error')"
else
  echo "  ❌ Should have rejected negative amount"
fi
echo ""

# Test 3: Invalid simpanan type
echo "3️⃣ Test: Invalid Simpanan Type"
RESPONSE=$(curl -s -X POST ${API_URL}/simpanan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"idAnggota\": \"$MEMBER_ID\",
    \"tipeSimpanan\": \"invalid_type\",
    \"tanggalTransaksi\": \"2025-01-01T00:00:00Z\",
    \"jumlahSetoran\": 100000,
    \"keterangan\": \"Test invalid type\"
  }")
STATUS=$(echo $RESPONSE | jq -r '.success // false')
if [ "$STATUS" = "false" ]; then
  echo "  ✅ Correctly rejected: $(echo $RESPONSE | jq -r '.message // .error')"
else
  echo "  ❌ Should have rejected invalid simpanan type"
fi
echo ""

# Test 4: Zero amount
echo "4️⃣ Test: Zero Amount"
RESPONSE=$(curl -s -X POST ${API_URL}/simpanan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"idAnggota\": \"$MEMBER_ID\",
    \"tipeSimpanan\": \"pokok\",
    \"tanggalTransaksi\": \"2025-01-01T00:00:00Z\",
    \"jumlahSetoran\": 0,
    \"keterangan\": \"Test zero amount\"
  }")
STATUS=$(echo $RESPONSE | jq -r '.success // false')
if [ "$STATUS" = "false" ]; then
  echo "  ✅ Correctly rejected: $(echo $RESPONSE | jq -r '.message // .error')"
else
  echo "  ❌ Should have rejected zero amount"
fi
echo ""

# Test 5: Missing required fields
echo "5️⃣ Test: Missing Required Field (idAnggota)"
RESPONSE=$(curl -s -X POST ${API_URL}/simpanan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "tipeSimpanan": "pokok",
    "tanggalTransaksi": "2025-01-01T00:00:00Z",
    "jumlahSetoran": 100000,
    "keterangan": "Test missing member"
  }')
STATUS=$(echo $RESPONSE | jq -r '.success // false')
if [ "$STATUS" = "false" ]; then
  echo "  ✅ Correctly rejected: $(echo $RESPONSE | jq -r '.message // .error')"
else
  echo "  ❌ Should have rejected missing idAnggota"
fi
echo ""

# Test 6: Invalid date format
echo "6️⃣ Test: Invalid Date Format"
RESPONSE=$(curl -s -X POST ${API_URL}/simpanan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"idAnggota\": \"$MEMBER_ID\",
    \"tipeSimpanan\": \"pokok\",
    \"tanggalTransaksi\": \"2025-01-01\",
    \"jumlahSetoran\": 100000,
    \"keterangan\": \"Test invalid date\"
  }")
STATUS=$(echo $RESPONSE | jq -r '.success // false')
if [ "$STATUS" = "false" ]; then
  echo "  ✅ Correctly rejected: $(echo $RESPONSE | jq -r '.message // .error')"
else
  echo "  ❌ Should have rejected invalid date format"
fi
echo ""

echo "✅ Edge case testing complete"
