#!/bin/bash

API_URL="http://localhost/api/v1"

# Login
TOKEN=$(curl -s -X POST ${API_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{"namaPengguna":"admin","kataSandi":"admin123"}' | jq -r '.data.token')

echo "Token: ${TOKEN:0:50}..."
echo ""

# Get first member
MEMBER_ID=$(curl -s -X GET "${API_URL}/anggota?page=1&pageSize=1" \
  -H "Authorization: Bearer $TOKEN" | jq -r '.data[0].id')

echo "Member ID: $MEMBER_ID"
echo ""

# Create deposit
echo "Creating deposit..."
curl -s -X POST ${API_URL}/simpanan \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d "{
    \"idAnggota\": \"$MEMBER_ID\",
    \"tipeSimpanan\": \"pokok\",
    \"tanggalTransaksi\": \"2025-11-18T00:00:00Z\",
    \"jumlahSetoran\": 100000,
    \"keterangan\": \"Test deposit - pokok\"
  }" | jq '.'
