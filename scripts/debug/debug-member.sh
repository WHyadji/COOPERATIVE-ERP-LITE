#!/bin/bash

API_URL="http://localhost/api/v1"

# Login
TOKEN=$(curl -s -X POST ${API_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{"namaPengguna":"admin","kataSandi":"admin123"}' | jq -r '.data.token')

echo "Token: ${TOKEN:0:50}..."
echo ""

# Create member and show full response
echo "Creating member..."
RESPONSE=$(curl -s -X POST ${API_URL}/anggota \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "namaLengkap": "Test Member",
    "nik": "1234567890123456",
    "noTelepon": "08123456789",
    "tanggalBergabung": "2025-01-01T00:00:00Z",
    "alamat": "Jl. Test"
  }')

echo "Full response:"
echo "$RESPONSE" | jq '.'
echo ""
echo "Member ID:"
echo "$RESPONSE" | jq -r '.data.id'
