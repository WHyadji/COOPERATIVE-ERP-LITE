#!/bin/bash

API_URL="http://localhost/api/v1"

# Login
echo "Logging in..."
TOKEN=$(curl -s -X POST ${API_URL}/auth/login \
  -H "Content-Type: application/json" \
  -d '{"namaPengguna":"admin","kataSandi":"admin123"}' | jq -r '.data.token')

echo "Token: ${TOKEN:0:50}..."
echo ""

# Create member
echo "Creating member..."
curl -s -X POST ${API_URL}/anggota \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "namaLengkap": "Budi Santoso",
    "nik": "3201234567890001",
    "noTelepon": "081234567801",
    "tanggalBergabung": "2025-01-01T00:00:00Z",
    "alamat": "Jl. Test No. 123, Jakarta"
  }' | jq '.'
