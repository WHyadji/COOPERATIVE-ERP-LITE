#!/bin/bash
TOKEN=$(curl -s -X POST http://localhost/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"namaPengguna":"admin","kataSandi":"admin123"}' | jq -r '.data.token')

echo "📊 Summary Endpoint:"
curl -s http://localhost/api/v1/simpanan/ringkasan \
  -H "Authorization: Bearer $TOKEN" | jq '.'

echo ""
echo "📋 Balance Report (first 3 members):"
curl -s http://localhost/api/v1/simpanan/laporan-saldo \
  -H "Authorization: Bearer $TOKEN" | jq '.data[0:3]'
