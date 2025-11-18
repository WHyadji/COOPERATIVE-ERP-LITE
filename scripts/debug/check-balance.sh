#!/bin/bash
TOKEN=$(curl -s -X POST http://localhost/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"namaPengguna":"admin","kataSandi":"admin123"}' | jq -r '.data.token')

curl -s -X GET "http://localhost/api/v1/simpanan/laporan-saldo" \
  -H "Authorization: Bearer $TOKEN" | jq '.data'
